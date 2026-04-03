package style

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// EmbeddedStyles is set by the main package which has access to the styles directory.
// It must provide a "styles/" subtree containing .yaml files.
var EmbeddedStyles fs.FS //nolint:gochecknoglobals

// Load loads a style by name. Searches user dir first, then embedded.
func Load(name string) (*Style, error) {
	// 1. User styles directory
	userDir := filepath.Join(os.Getenv("HOME"), ".md2docx", "styles")
	userPath := filepath.Join(userDir, name+".yaml")
	if data, err := os.ReadFile(userPath); err == nil {
		return parseStyle(data)
	}

	// 2. Embedded styles
	data, err := fs.ReadFile(EmbeddedStyles, "styles/"+name+".yaml")
	if err != nil {
		return nil, fmt.Errorf("样式 '%s' 不存在", name)
	}
	return parseStyle(data)
}

// ListAvailable returns all available style names.
func ListAvailable() []string {
	var names []string
	seen := make(map[string]bool)

	// Embedded
	entries, err := fs.ReadDir(EmbeddedStyles, "styles")
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
				n := strings.TrimSuffix(e.Name(), ".yaml")
				names = append(names, n)
				seen[n] = true
			}
		}
	}

	// User directory
	userDir := filepath.Join(os.Getenv("HOME"), ".md2docx", "styles")
	entries2, err := os.ReadDir(userDir)
	if err == nil {
		for _, e := range entries2 {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
				n := strings.TrimSuffix(e.Name(), ".yaml")
				if !seen[n] {
					names = append(names, n)
				}
			}
		}
	}

	return names
}

func parseStyle(data []byte) (*Style, error) {
	s := &Style{}
	if err := yaml.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("样式解析失败: %w", err)
	}
	applyDefaults(s)
	return s, nil
}

// ApplyFrontmatterOverrides applies frontmatter fields to override style settings.
func ApplyFrontmatterOverrides(s *Style, meta map[string]any) {
	if v, ok := meta["toc"]; ok {
		if b, ok := v.(bool); ok {
			s.TOC.Enabled = b
		}
	}
	if v, ok := meta["heading_numbering"]; ok {
		if b, ok := v.(bool); ok {
			s.HeadingNumbering.Enabled = b
		}
	}
	if v, ok := meta["cover"]; ok {
		if b, ok := v.(bool); ok {
			s.Cover.Enabled = b
		}
	}
	if v, ok := meta["page_size"]; ok {
		if ps, ok := v.(string); ok {
			s.Page.Size = ps
		}
	}
	if v, ok := meta["header_left"]; ok {
		if hl, ok := v.(string); ok {
			s.Header.Left = hl
		}
	}
	if v, ok := meta["header_right"]; ok {
		if hr, ok := v.(string); ok {
			s.Header.Right = hr
		}
	}
	if v, ok := meta["footer_center"]; ok {
		if fc, ok := v.(string); ok {
			s.Footer.Center = fc
		}
	}
}

func applyDefaults(s *Style) {
	if s.Page.Size == "" {
		s.Page.Size = "A4"
	}
	if s.Page.Orientation == "" {
		s.Page.Orientation = "portrait"
	}
	if s.Page.Margin.Top == "" {
		s.Page.Margin.Top = "25.4mm"
	}
	if s.Page.Margin.Bottom == "" {
		s.Page.Margin.Bottom = "25.4mm"
	}
	if s.Page.Margin.Left == "" {
		s.Page.Margin.Left = "31.8mm"
	}
	if s.Page.Margin.Right == "" {
		s.Page.Margin.Right = "31.8mm"
	}
	if s.Page.Margin.Header == "" {
		s.Page.Margin.Header = "15mm"
	}
	if s.Page.Margin.Footer == "" {
		s.Page.Margin.Footer = "15mm"
	}
	if s.Fonts.Latin == "" {
		s.Fonts.Latin = "Times New Roman"
	}
	if s.Fonts.CJK == "" {
		s.Fonts.CJK = "宋体"
	}
	if s.Fonts.Mono == "" {
		s.Fonts.Mono = "Courier New"
	}
	if s.TOC.Title == "" {
		s.TOC.Title = "目录"
	}
	if s.TOC.Depth == 0 {
		s.TOC.Depth = 3
	}

	// Default heading styles (pageBreakBefore not set here; defaults to false via Go zero value)
	setHeadingDefaults(&s.Styles.Heading1, "22pt", "center", "24pt", "12pt")
	setHeadingDefaults(&s.Styles.Heading2, "16pt", "left", "12pt", "6pt")
	setHeadingDefaults(&s.Styles.Heading3, "14pt", "left", "6pt", "3pt")
	setHeadingDefaults(&s.Styles.Heading4, "12pt", "left", "6pt", "3pt")
	setHeadingDefaults(&s.Styles.Heading5, "11pt", "left", "3pt", "3pt")
	setHeadingDefaults(&s.Styles.Heading6, "10.5pt", "left", "3pt", "3pt")

	if s.Styles.Body.FontSize == "" {
		s.Styles.Body.FontSize = "12pt"
	}
	if s.Styles.Body.Alignment == "" {
		s.Styles.Body.Alignment = "justify"
	}
	if s.Styles.Body.LineSpacing == 0 {
		s.Styles.Body.LineSpacing = 1.5
	}
	if s.Styles.Body.FirstLineIndent == "" {
		s.Styles.Body.FirstLineIndent = "2em"
	}

	if s.Styles.Quote.FontSize == "" {
		s.Styles.Quote.FontSize = "10.5pt"
	}
	if s.Styles.Quote.FontCJK == "" {
		s.Styles.Quote.FontCJK = "楷体"
	}
	if s.Styles.Quote.Alignment == "" {
		s.Styles.Quote.Alignment = "justify"
	}
	if s.Styles.Quote.LeftIndent == "" {
		s.Styles.Quote.LeftIndent = "2cm"
	}
	if s.Styles.Quote.RightIndent == "" {
		s.Styles.Quote.RightIndent = "2cm"
	}

	if s.Styles.CodeBlock.FontSize == "" {
		s.Styles.CodeBlock.FontSize = "9pt"
	}
	if s.Styles.CodeBlock.FontLatin == "" {
		s.Styles.CodeBlock.FontLatin = "Courier New"
	}

	if s.Styles.Image.MaxWidthPct == 0 {
		s.Styles.Image.MaxWidthPct = 80
	}

	if s.Styles.Table.Border {
		// already set
	} else {
		s.Styles.Table.Border = true
	}
	if s.Styles.Table.HeaderBold {
		// already set
	} else {
		s.Styles.Table.HeaderBold = true
	}

	if s.List.BulletChars == nil {
		s.List.BulletChars = []string{"●", "○", "■"}
	}
	if s.List.OrderedFormat == "" {
		s.List.OrderedFormat = "{n}."
	}
	if s.List.IndentPerLevel == "" {
		s.List.IndentPerLevel = "0.75cm"
	}
}

func setHeadingDefaults(ps *ParagraphStyle, fontSize, align, before, after string) {
	if ps.FontSize == "" {
		ps.FontSize = fontSize
	}
	if !ps.Bold {
		ps.Bold = true
	}
	if ps.Alignment == "" {
		ps.Alignment = align
	}
	if ps.SpaceBefore == "" {
		ps.SpaceBefore = before
	}
	if ps.SpaceAfter == "" {
		ps.SpaceAfter = after
	}
	if !ps.KeepWithNext {
		ps.KeepWithNext = true
	}
}
