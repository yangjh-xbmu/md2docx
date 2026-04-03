package render

import (
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// buildTable converts a goldmark Table AST node to an OOXML Table element.
func buildTable(table *east.Table, source []byte, ts style.TableStyle) *ooxml.Table {
	// Collect all rows with header flag
	type rowData struct {
		cells    [][]ast.Node // each cell's child nodes
		isHeader bool
	}
	var rows []rowData
	colCount := 0

	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		isHeader := false
		if _, ok := child.(*east.TableHeader); ok {
			isHeader = true
		}
		var cells [][]ast.Node
		for cell := child.FirstChild(); cell != nil; cell = cell.NextSibling() {
			var children []ast.Node
			for c := cell.FirstChild(); c != nil; c = c.NextSibling() {
				children = append(children, c)
			}
			cells = append(cells, children)
		}
		if len(cells) > colCount {
			colCount = len(cells)
		}
		rows = append(rows, rowData{cells: cells, isHeader: isHeader})
	}

	if len(rows) == 0 || colCount == 0 {
		return nil
	}

	// Build table properties
	borderColor := "000000"
	if ts.BorderColor != "" {
		borderColor = ooxml.ColorHex(ts.BorderColor)
	}

	tblPr := &ooxml.TableProperties{
		TblW: &ooxml.TableWidth{W: "0", Type: "auto"},
	}

	if ts.Alignment != "" {
		tblPr.Jc = &ooxml.Justification{Val: ooxml.AlignmentVal(ts.Alignment)}
	}

	if ts.Border {
		border := &ooxml.BorderVal{Val: "single", Sz: "4", Space: "0", Color: borderColor}
		tblPr.TblBorders = &ooxml.TableBorders{
			Top: border, Left: border, Bottom: border, Right: border,
			InsideH: border, InsideV: border,
		}
	}

	tblPr.TblLook = &ooxml.TableLook{
		Val: "04A0", FirstRow: "1", LastRow: "0",
		FirstColumn: "1", LastColumn: "0", NoHBand: "0", NoVBand: "1",
	}

	// Build grid columns (equal width)
	var gridCols []ooxml.GridCol
	for i := 0; i < colCount; i++ {
		gridCols = append(gridCols, ooxml.GridCol{W: "0"})
	}

	// Build rows
	var tblRows []ooxml.TableRow
	for _, row := range rows {
		tr := ooxml.TableRow{}
		if row.isHeader {
			tr.TrPr = &ooxml.TableRowProperties{TblHeader: &ooxml.Empty{}}
		}

		for _, cellNodes := range row.cells {
			tc := ooxml.TableCell{
				TcPr: &ooxml.TableCellProperties{
					TcW: &ooxml.TableWidth{W: "0", Type: "auto"},
				},
			}

			// Header row shading
			if row.isHeader && ts.HeaderBg != "" {
				tc.TcPr.Shd = &ooxml.Shading{
					Val:  "clear",
					Fill: ooxml.ColorHex(ts.HeaderBg),
				}
			}

			// Build paragraph with cell text
			p := buildCellParagraph(cellNodes, source, row.isHeader, ts)
			tc.Content = []any{p}
			tr.Cells = append(tr.Cells, tc)
		}

		// Pad missing cells
		for len(tr.Cells) < colCount {
			tc := ooxml.TableCell{
				TcPr: &ooxml.TableCellProperties{
					TcW: &ooxml.TableWidth{W: "0", Type: "auto"},
				},
				Content: []any{&ooxml.Paragraph{}},
			}
			tr.Cells = append(tr.Cells, tc)
		}

		tblRows = append(tblRows, tr)
	}

	return &ooxml.Table{
		TblPr:   tblPr,
		TblGrid: &ooxml.TableGrid{Cols: gridCols},
		Rows:    tblRows,
	}
}

// buildCellParagraph creates a paragraph from cell AST nodes.
func buildCellParagraph(nodes []ast.Node, source []byte, isHeader bool, ts style.TableStyle) *ooxml.Paragraph {
	var runs []any

	for _, node := range nodes {
		switch n := node.(type) {
		case *ast.Text:
			rpr := cellRunProps(isHeader, ts)
			runs = append(runs, &ooxml.Run{
				RPr:     rpr,
				Content: []any{&ooxml.Text{Space: "preserve", Value: string(n.Segment.Value(source))}},
			})
		case *ast.Emphasis:
			rpr := cellRunProps(isHeader, ts)
			if rpr == nil {
				rpr = &ooxml.RunProperties{}
			}
			if n.Level == 1 {
				rpr.Italic = &ooxml.Empty{}
			} else {
				rpr.Bold = &ooxml.Empty{}
			}
			text := extractText(n, source)
			runs = append(runs, &ooxml.Run{
				RPr:     rpr,
				Content: []any{&ooxml.Text{Space: "preserve", Value: text}},
			})
		case *ast.CodeSpan:
			rpr := cellRunProps(isHeader, ts)
			if rpr == nil {
				rpr = &ooxml.RunProperties{}
			}
			rpr.RStyle = &ooxml.SVal{Val: "InlineCode"}
			runs = append(runs, &ooxml.Run{
				RPr:     rpr,
				Content: []any{&ooxml.Text{Space: "preserve", Value: string(n.Text(source))}},
			})
		default:
			// Fallback: extract raw text
			text := extractText(node, source)
			if text != "" {
				rpr := cellRunProps(isHeader, ts)
				runs = append(runs, &ooxml.Run{
					RPr:     rpr,
					Content: []any{&ooxml.Text{Space: "preserve", Value: text}},
				})
			}
		}
	}

	ppr := &ooxml.ParagraphProperties{}
	if ts.FontSize != "" {
		// Table font size is applied via run properties, not paragraph
	}

	return &ooxml.Paragraph{PPr: ppr, Runs: runs}
}

func cellRunProps(isHeader bool, ts style.TableStyle) *ooxml.RunProperties {
	var rpr *ooxml.RunProperties

	needProps := isHeader && (ts.HeaderBold || ts.HeaderColor != "") || ts.FontSize != ""
	if !needProps {
		return nil
	}

	rpr = &ooxml.RunProperties{}
	if isHeader && ts.HeaderBold {
		rpr.Bold = &ooxml.Empty{}
	}
	if isHeader && ts.HeaderColor != "" {
		rpr.Color = &ooxml.SVal{Val: ooxml.ColorHex(ts.HeaderColor)}
	}
	if ts.FontSize != "" {
		sz := ooxml.FontSizeHalfPoints(ts.FontSize)
		rpr.Sz = &ooxml.SzVal{Val: sz}
		rpr.SzCs = &ooxml.SzVal{Val: sz}
	}
	return rpr
}

// extractText recursively extracts text from an AST node.
func extractText(node ast.Node, source []byte) string {
	var text string
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			text += string(t.Segment.Value(source))
		} else {
			text += extractText(c, source)
		}
	}
	return text
}
