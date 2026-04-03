package ooxml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// newTestStyle returns a style.Style with defaults populated,
// matching the behaviour of the unexported applyDefaults.
func newTestStyle() *style.Style {
	return &style.Style{
		Fonts: style.FontsConfig{
			Latin: "Times New Roman",
			CJK:   "宋体",
			Mono:  "Courier New",
		},
		Styles: style.StylesConfig{
			Body: style.ParagraphStyle{
				FontSize:        "12pt",
				Alignment:       "justify",
				LineSpacing:     1.5,
				FirstLineIndent: "2em",
			},
			Heading1: style.ParagraphStyle{
				FontSize:  "22pt",
				Alignment: "center",
				Bold:      true,
			},
			Heading2: style.ParagraphStyle{
				FontSize:  "16pt",
				Alignment: "left",
				Bold:      true,
			},
			Heading3: style.ParagraphStyle{
				FontSize:  "14pt",
				Alignment: "left",
				Bold:      true,
			},
			Heading4: style.ParagraphStyle{
				FontSize:  "12pt",
				Alignment: "left",
				Bold:      true,
			},
			Heading5: style.ParagraphStyle{
				FontSize:  "11pt",
				Alignment: "left",
				Bold:      true,
			},
			Heading6: style.ParagraphStyle{
				FontSize:  "10.5pt",
				Alignment: "left",
				Bold:      true,
			},
			Quote: style.ParagraphStyle{
				FontSize:    "10.5pt",
				FontCJK:     "楷体",
				Alignment:   "justify",
				LeftIndent:  "2cm",
				RightIndent: "2cm",
			},
			CodeBlock: style.ParagraphStyle{
				FontSize:    "9pt",
				LineSpacing: 1.0,
			},
		},
	}
}

func TestGenerateStylesXML_RootElement(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	assert.Contains(t, xml, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	assert.Contains(t, xml, `<w:styles xmlns:w=`)
	assert.Contains(t, xml, `</w:styles>`)
}

func TestGenerateStylesXML_DocDefaults(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	assert.Contains(t, xml, `<w:docDefaults>`)
	assert.Contains(t, xml, `<w:rPrDefault>`)
	assert.Contains(t, xml, `<w:pPrDefault>`)
	// Default fonts
	assert.Contains(t, xml, `w:ascii="Times New Roman"`)
	assert.Contains(t, xml, `w:hAnsi="Times New Roman"`)
	assert.Contains(t, xml, `w:eastAsia="宋体"`)
}

func TestGenerateStylesXML_NormalStyle(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	assert.Contains(t, xml, `w:styleId="Normal"`)
	assert.Contains(t, xml, `<w:name w:val="Normal"/>`)
}

func TestGenerateStylesXML_AllHeadingStyles(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	headings := []struct {
		id   string
		name string
	}{
		{"Heading1", "heading 1"},
		{"Heading2", "heading 2"},
		{"Heading3", "heading 3"},
		{"Heading4", "heading 4"},
		{"Heading5", "heading 5"},
		{"Heading6", "heading 6"},
	}
	for _, h := range headings {
		assert.Contains(t, xml, `w:styleId="`+h.id+`"`, "missing style %s", h.id)
		assert.Contains(t, xml, `<w:name w:val="`+h.name+`"/>`, "missing name %s", h.name)
	}
}

func TestGenerateStylesXML_QuoteCodeInlineCodeHyperlink(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	assert.Contains(t, xml, `w:styleId="Quote"`)
	assert.Contains(t, xml, `w:styleId="Code"`)
	assert.Contains(t, xml, `w:styleId="InlineCode"`)
	assert.Contains(t, xml, `w:styleId="Hyperlink"`)
	// Hyperlink specifics
	assert.Contains(t, xml, `<w:color w:val="0563C1"/>`)
	assert.Contains(t, xml, `<w:u w:val="single"/>`)
}

func TestGenerateStylesXML_CJKDualFonts(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	// w:eastAsia must appear multiple times (docDefaults + each style)
	count := strings.Count(xml, `w:eastAsia=`)
	require.GreaterOrEqual(t, count, 2, "expected multiple w:eastAsia attributes for CJK dual fonts")
}

func TestGenerateStylesXML_FontSizeHalfPoints(t *testing.T) {
	s := newTestStyle()
	// Body 12pt -> half-points = 24
	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `<w:sz w:val="24"/>`)
	assert.Contains(t, xml, `<w:szCs w:val="24"/>`)

	// Heading1 22pt -> half-points = 44
	assert.Contains(t, xml, `<w:sz w:val="44"/>`)
}

func TestGenerateStylesXML_HeadingDefaultFonts(t *testing.T) {
	// When heading has no FontCJK/FontLatin set,
	// generateHeadingStyleXML defaults to 黑体/Arial
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	// Must contain 黑体 (heading default CJK) and Arial (heading default Latin)
	assert.Contains(t, xml, "黑体")
	assert.Contains(t, xml, "Arial")
}

func TestGenerateStylesXML_CustomHeadingFont(t *testing.T) {
	s := newTestStyle()
	s.Styles.Heading1.FontCJK = "微软雅黑"
	s.Styles.Heading1.FontLatin = "Helvetica"

	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `w:eastAsia="微软雅黑"`)
	assert.Contains(t, xml, `w:ascii="Helvetica"`)
}

func TestGenerateStylesXML_BodyAlignment(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	// justify -> "both" in OOXML
	assert.Contains(t, xml, `<w:jc w:val="both"/>`)
}

func TestGenerateStylesXML_CodeBlockMonoFont(t *testing.T) {
	s := newTestStyle()
	xml := string(GenerateStylesXML(s))

	// Code style uses Mono font
	assert.Contains(t, xml, `w:ascii="Courier New"`)
	// Code font size 9pt -> 18 half-points
	assert.Contains(t, xml, `<w:sz w:val="18"/>`)
}

func TestGenerateStylesXML_BoldAndItalic(t *testing.T) {
	s := newTestStyle()
	s.Styles.Body.Italic = true

	xml := string(GenerateStylesXML(s))
	// Heading1 has Bold=true
	assert.Contains(t, xml, `<w:b/>`)
	// Body now has Italic
	assert.Contains(t, xml, `<w:i/>`)
}

func TestGenerateStylesXML_Color(t *testing.T) {
	s := newTestStyle()
	s.Styles.Body.Color = "#333333"

	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `<w:color w:val="333333"/>`)
}

func TestGenerateStylesXML_SpacingAndIndent(t *testing.T) {
	s := newTestStyle()
	s.Styles.Heading1.SpaceBefore = "24pt"
	s.Styles.Heading1.SpaceAfter = "12pt"

	xml := string(GenerateStylesXML(s))
	// 24pt = 480 twips
	assert.Contains(t, xml, `w:before="480"`)
	// 12pt = 240 twips
	assert.Contains(t, xml, `w:after="240"`)
}

func TestGenerateStylesXML_KeepWithNext(t *testing.T) {
	s := newTestStyle()
	s.Styles.Heading1.KeepWithNext = true

	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `<w:keepNext/>`)
}

func TestGenerateStylesXML_PageBreakBefore(t *testing.T) {
	s := newTestStyle()
	s.Styles.Heading1.PageBreakBefore = true

	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `<w:pageBreakBefore/>`)
}

func TestGenerateStylesXML_LineSpacing(t *testing.T) {
	s := newTestStyle()
	// Body line spacing 1.5 -> 360
	xml := string(GenerateStylesXML(s))
	assert.Contains(t, xml, `w:line="360"`)
}
