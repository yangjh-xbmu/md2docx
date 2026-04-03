package render

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// DocBuilder accumulates document elements during AST walking.
type DocBuilder struct {
	Style    *style.Style
	Meta     map[string]any
	BaseDir  string // directory of source .md file for resolving relative paths
	elements []any
	curRuns  []any
	inPara   bool
	curPPr   *ooxml.ParagraphProperties
	// List tracking
	listDepth  int
	listType   string // "bullet" or "ordered"
	// Heading numbering
	numberer *HeadingNumberer
}

// ToDocx converts a goldmark AST to a .docx file.
func ToDocx(doc ast.Node, source []byte, s *style.Style, meta map[string]any, outputPath string, baseDir string) error {
	b := &DocBuilder{
		Style:    s,
		Meta:     meta,
		BaseDir:  baseDir,
		numberer: NewHeadingNumberer(s.HeadingNumbering),
	}

	// Walk the AST
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		return b.renderNode(n, source, entering)
	})
	if err != nil {
		return fmt.Errorf("AST 遍历失败: %w", err)
	}

	// Build OOXML package
	pkg := b.buildPackage()

	return ooxml.WriteDocx(pkg, outputPath)
}

func (b *DocBuilder) renderNode(n ast.Node, source []byte, entering bool) (ast.WalkStatus, error) {
	switch node := n.(type) {
	case *ast.Document:
		return ast.WalkContinue, nil

	case *ast.Heading:
		if entering {
			styleID := fmt.Sprintf("Heading%d", node.Level)
			b.startParagraph(&ooxml.ParagraphProperties{
				PStyle: &ooxml.PStyle{Val: styleID},
			})
			// Add heading numbering prefix
			if prefix := b.numberer.FormatNumber(node.Level); prefix != "" {
				b.addTextRun(prefix+" ", nil)
			}
		} else {
			b.endParagraph()
		}

	case *ast.Paragraph:
		if entering {
			// Don't start a paragraph if we're inside a list item (list item handles it)
			if n.Parent() != nil && n.Parent().Kind() == ast.KindListItem {
				return ast.WalkContinue, nil
			}
			b.startParagraph(&ooxml.ParagraphProperties{
				PStyle: &ooxml.PStyle{Val: "Normal"},
			})
		} else {
			if n.Parent() != nil && n.Parent().Kind() == ast.KindListItem {
				return ast.WalkContinue, nil
			}
			b.endParagraph()
		}

	case *ast.List:
		if entering {
			b.listDepth++
			if node.IsOrdered() {
				b.listType = "ordered"
			} else {
				b.listType = "bullet"
			}
		} else {
			b.listDepth--
			if b.listDepth == 0 {
				b.listType = ""
			}
		}

	case *ast.ListItem:
		if entering {
			numID := "1" // bullet
			if b.listType == "ordered" {
				numID = "2"
			}
			ilvl := fmt.Sprintf("%d", b.listDepth-1)
			b.startParagraph(&ooxml.ParagraphProperties{
				NumPr: &ooxml.NumberingPr{
					Ilvl:  &ooxml.IntVal{Val: ilvl},
					NumId: &ooxml.IntVal{Val: numID},
				},
			})
		} else {
			b.endParagraph()
		}

	case *ast.Blockquote:
		// Handled by child paragraphs with Quote style
		if entering {
			// Override next paragraph style to Quote
		}

	case *ast.FencedCodeBlock, *ast.CodeBlock:
		if entering {
			b.startParagraph(&ooxml.ParagraphProperties{
				PStyle: &ooxml.PStyle{Val: "Code"},
			})
			// Collect all lines
			lines := node.(ast.Node).Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				text := string(line.Value(source))
				if i > 0 {
					b.addRun(nil, &ooxml.Break{})
				}
				b.addTextRun(text, nil)
			}
			b.endParagraph()
			return ast.WalkSkipChildren, nil
		}

	case *ast.ThematicBreak:
		if entering {
			// Add a horizontal rule as a paragraph with bottom border
			b.startParagraph(&ooxml.ParagraphProperties{})
			b.endParagraph()
		}

	// Inline elements
	case *ast.Text:
		if entering {
			text := string(node.Segment.Value(source))
			b.addTextRun(text, nil)
			if node.SoftLineBreak() {
				b.addTextRun(" ", nil)
			}
			if node.HardLineBreak() {
				b.addRun(nil, &ooxml.Break{})
			}
		}

	case *ast.Emphasis:
		if entering {
			// Emphasis level 1 = italic, level 2 = bold
			rpr := &ooxml.RunProperties{}
			if node.Level == 1 {
				rpr.Italic = &ooxml.Empty{}
			} else {
				rpr.Bold = &ooxml.Empty{}
			}
			b.pushRunProps(rpr)
		} else {
			b.popRunProps()
		}

	case *ast.CodeSpan:
		if entering {
			rpr := &ooxml.RunProperties{
				RStyle: &ooxml.SVal{Val: "InlineCode"},
			}
			text := string(node.Text(source))
			b.addTextRun(text, rpr)
			return ast.WalkSkipChildren, nil
		}

	case *ast.Link:
		if entering {
			// For now, render link text with Hyperlink style
			rpr := &ooxml.RunProperties{
				RStyle: &ooxml.SVal{Val: "Hyperlink"},
			}
			b.pushRunProps(rpr)
		} else {
			b.popRunProps()
		}

	case *ast.AutoLink:
		if entering {
			text := string(node.URL(source))
			rpr := &ooxml.RunProperties{
				RStyle: &ooxml.SVal{Val: "Hyperlink"},
			}
			b.addTextRun(text, rpr)
			return ast.WalkSkipChildren, nil
		}

	case *east.Table:
		if entering {
			b.renderTable(node, source)
			return ast.WalkSkipChildren, nil
		}

	default:
		// Skip unknown nodes
	}

	return ast.WalkContinue, nil
}

// Run property stack for nested inline formatting
var runPropsStack []*ooxml.RunProperties

func (b *DocBuilder) pushRunProps(rpr *ooxml.RunProperties) {
	runPropsStack = append(runPropsStack, rpr)
}

func (b *DocBuilder) popRunProps() {
	if len(runPropsStack) > 0 {
		runPropsStack = runPropsStack[:len(runPropsStack)-1]
	}
}

func (b *DocBuilder) currentRunProps() *ooxml.RunProperties {
	if len(runPropsStack) == 0 {
		return nil
	}
	// Merge all props in stack
	merged := &ooxml.RunProperties{}
	for _, rpr := range runPropsStack {
		if rpr.Bold != nil {
			merged.Bold = rpr.Bold
		}
		if rpr.Italic != nil {
			merged.Italic = rpr.Italic
		}
		if rpr.Underline != nil {
			merged.Underline = rpr.Underline
		}
		if rpr.RStyle != nil {
			merged.RStyle = rpr.RStyle
		}
		if rpr.RFonts != nil {
			merged.RFonts = rpr.RFonts
		}
		if rpr.Color != nil {
			merged.Color = rpr.Color
		}
	}
	return merged
}

func (b *DocBuilder) startParagraph(ppr *ooxml.ParagraphProperties) {
	if b.inPara {
		b.endParagraph()
	}
	b.inPara = true
	b.curPPr = ppr
	b.curRuns = nil
}

func (b *DocBuilder) endParagraph() {
	if !b.inPara {
		return
	}
	p := &ooxml.Paragraph{
		PPr:  b.curPPr,
		Runs: b.curRuns,
	}
	b.elements = append(b.elements, p)
	b.inPara = false
	b.curPPr = nil
	b.curRuns = nil
}

func (b *DocBuilder) addTextRun(text string, rpr *ooxml.RunProperties) {
	if rpr == nil {
		rpr = b.currentRunProps()
	}
	t := &ooxml.Text{
		Space: "preserve",
		Value: text,
	}
	run := &ooxml.Run{
		RPr:     rpr,
		Content: []any{t},
	}
	b.curRuns = append(b.curRuns, run)
}

func (b *DocBuilder) addRun(rpr *ooxml.RunProperties, content ...any) {
	run := &ooxml.Run{
		RPr:     rpr,
		Content: content,
	}
	b.curRuns = append(b.curRuns, run)
}

func (b *DocBuilder) buildPackage() *ooxml.Package {
	// Prepend TOC elements if enabled (TOC comes after cover)
	if tocElems := GenerateTOCElements(b.Style.TOC); len(tocElems) > 0 {
		b.elements = append(tocElems, b.elements...)
	}

	// Prepend cover elements if enabled (cover comes first)
	if coverElems := GenerateCoverElements(b.Style.Cover, b.Meta); len(coverElems) > 0 {
		b.elements = append(coverElems, b.elements...)
	}

	// Page setup
	pageW, pageH := ooxml.PageSizeTwips(b.Style.Page.Size)
	orient := ""
	if b.Style.Page.Orientation == "landscape" {
		orient = "landscape"
		pageW, pageH = pageH, pageW
	}

	sectPr := &ooxml.SectionProperties{
		PgSz: &ooxml.PageSize{
			W:      pageW,
			H:      pageH,
			Orient: orient,
		},
		PgMar: &ooxml.PageMargins{
			Top:    ooxml.MeasureToTwips(b.Style.Page.Margin.Top),
			Bottom: ooxml.MeasureToTwips(b.Style.Page.Margin.Bottom),
			Left:   ooxml.MeasureToTwips(b.Style.Page.Margin.Left),
			Right:  ooxml.MeasureToTwips(b.Style.Page.Margin.Right),
			Header: ooxml.MeasureToTwips(b.Style.Page.Margin.Header),
			Footer: ooxml.MeasureToTwips(b.Style.Page.Margin.Footer),
		},
	}

	pkg := &ooxml.Package{
		Styles: ooxml.GenerateStylesXML(b.Style),
	}

	// Generate numbering.xml for lists
	pkg.Numbering = ooxml.GenerateNumberingXML(b.Style)

	// Generate header/footer XML and wire up references
	// Relationship IDs must match what writeDocRels produces.
	// rId1 = styles, rId2 = numbering (if present), then header, then footer.
	nextID := 2
	if pkg.Numbering != nil {
		nextID++
	}

	headerXML := ooxml.GenerateHeaderXML(b.Style.Header, b.Meta)
	if headerXML != nil {
		pkg.Header = headerXML
		sectPr.HeaderRef = &ooxml.HeaderFooterRef{
			Type: "default",
			RID:  fmt.Sprintf("rId%d", nextID),
		}
		nextID++
	}

	footerXML := ooxml.GenerateFooterXML(b.Style.Footer, b.Meta)
	if footerXML != nil {
		pkg.Footer = footerXML
		sectPr.FooterRef = &ooxml.HeaderFooterRef{
			Type: "default",
			RID:  fmt.Sprintf("rId%d", nextID),
		}
	}

	doc := &ooxml.Document{
		W:  "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		R:  "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
		WP: "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing",
		Body: ooxml.Body{
			Elements: b.elements,
			SectPr:   sectPr,
		},
	}
	pkg.Document = doc

	return pkg
}

func (b *DocBuilder) renderTable(table *east.Table, source []byte) {
	// Collect rows: goldmark Table has TableHeader and TableRow children
	var rows [][]string
	isHeaderRow := true
	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		switch row := child.(type) {
		case *east.TableHeader:
			cells := collectTableCells(row, source)
			if len(cells) > 0 {
				rows = append(rows, cells)
			}
		case *east.TableRow:
			_ = row
			cells := collectTableCells(child, source)
			rows = append(rows, cells)
		}
		isHeaderRow = false
	}
	_ = isHeaderRow

	if len(rows) == 0 {
		return
	}

	// Generate table XML as a raw paragraph (simplified)
	// For now, render as paragraphs with tab-separated cells
	for i, row := range rows {
		ppr := &ooxml.ParagraphProperties{
			PStyle: &ooxml.PStyle{Val: "Normal"},
		}
		b.startParagraph(ppr)
		for j, cell := range row {
			if j > 0 {
				b.addTextRun("\t", nil)
			}
			rpr := (*ooxml.RunProperties)(nil)
			if i == 0 {
				rpr = &ooxml.RunProperties{Bold: &ooxml.Empty{}}
			}
			b.addTextRun(cell, rpr)
		}
		b.endParagraph()
	}
}

func collectTableCells(row ast.Node, source []byte) []string {
	var cells []string
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		text := ""
		for c := cell.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				text += string(t.Segment.Value(source))
			}
		}
		cells = append(cells, text)
	}
	return cells
}
