package ooxml

import (
	"fmt"

	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// GenerateStylesXML generates word/styles.xml from a Style config.
func GenerateStylesXML(s *style.Style) []byte {
	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:docDefaults>
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:ascii="` + s.Fonts.Latin + `" w:hAnsi="` + s.Fonts.Latin + `" w:eastAsia="` + s.Fonts.CJK + `" w:cs="` + s.Fonts.Latin + `"/>
        <w:sz w:val="` + fontSizeTwips(s.Styles.Body.FontSize) + `"/>
        <w:szCs w:val="` + fontSizeTwips(s.Styles.Body.FontSize) + `"/>
        <w:lang w:val="en-US" w:eastAsia="zh-CN"/>
      </w:rPr>
    </w:rPrDefault>
    <w:pPrDefault>
      <w:pPr>
        <w:spacing w:line="` + lineSpacingVal(s.Styles.Body.LineSpacing) + `" w:lineRule="auto"/>
      </w:pPr>
    </w:pPrDefault>
  </w:docDefaults>`

	// Normal style
	xml += generateParagraphStyleXML("Normal", "Normal", &s.Styles.Body, s)

	// Heading styles
	headings := []struct {
		id    string
		name  string
		ps    *style.ParagraphStyle
	}{
		{"Heading1", "heading 1", &s.Styles.Heading1},
		{"Heading2", "heading 2", &s.Styles.Heading2},
		{"Heading3", "heading 3", &s.Styles.Heading3},
		{"Heading4", "heading 4", &s.Styles.Heading4},
		{"Heading5", "heading 5", &s.Styles.Heading5},
		{"Heading6", "heading 6", &s.Styles.Heading6},
	}
	for _, h := range headings {
		xml += generateHeadingStyleXML(h.id, h.name, h.ps, s)
	}

	// Quote style
	xml += generateParagraphStyleXML("Quote", "Quote", &s.Styles.Quote, s)

	// Code style
	xml += generateCodeStyleXML(s)

	// Hyperlink character style
	xml += `
  <w:style w:type="character" w:styleId="Hyperlink">
    <w:name w:val="Hyperlink"/>
    <w:rPr>
      <w:color w:val="0563C1"/>
      <w:u w:val="single"/>
    </w:rPr>
  </w:style>`

	xml += `
</w:styles>`

	return []byte(xml)
}

func generateParagraphStyleXML(id, name string, ps *style.ParagraphStyle, s *style.Style) string {
	fontCJK := ps.FontCJK
	if fontCJK == "" {
		fontCJK = s.Fonts.CJK
	}
	fontLatin := ps.FontLatin
	if fontLatin == "" {
		fontLatin = s.Fonts.Latin
	}

	result := fmt.Sprintf(`
  <w:style w:type="paragraph" w:styleId="%s">
    <w:name w:val="%s"/>
    <w:pPr>`, id, name)

	if ps.Alignment != "" {
		result += fmt.Sprintf(`
      <w:jc w:val="%s"/>`, alignmentVal(ps.Alignment))
	}

	spaceBefore := twipsFromPt(ps.SpaceBefore)
	spaceAfter := twipsFromPt(ps.SpaceAfter)
	if spaceBefore != "" || spaceAfter != "" || ps.LineSpacing > 0 {
		result += `
      <w:spacing`
		if spaceBefore != "" {
			result += fmt.Sprintf(` w:before="%s"`, spaceBefore)
		}
		if spaceAfter != "" {
			result += fmt.Sprintf(` w:after="%s"`, spaceAfter)
		}
		if ps.LineSpacing > 0 {
			result += fmt.Sprintf(` w:line="%s" w:lineRule="auto"`, lineSpacingVal(ps.LineSpacing))
		}
		result += `/>`
	}

	if ps.FirstLineIndent != "" || ps.LeftIndent != "" || ps.RightIndent != "" {
		result += `
      <w:ind`
		if ps.FirstLineIndent != "" {
			result += fmt.Sprintf(` w:firstLine="%s"`, indentTwips(ps.FirstLineIndent, ps.FontSize))
		}
		if ps.LeftIndent != "" {
			result += fmt.Sprintf(` w:left="%s"`, measureToTwips(ps.LeftIndent))
		}
		if ps.RightIndent != "" {
			result += fmt.Sprintf(` w:right="%s"`, measureToTwips(ps.RightIndent))
		}
		result += `/>`
	}

	if ps.KeepWithNext {
		result += `
      <w:keepNext/>`
	}
	if ps.PageBreakBefore {
		result += `
      <w:pageBreakBefore/>`
	}

	result += `
    </w:pPr>
    <w:rPr>
      <w:rFonts w:ascii="` + fontLatin + `" w:hAnsi="` + fontLatin + `" w:eastAsia="` + fontCJK + `"/>`

	if ps.FontSize != "" {
		sz := fontSizeTwips(ps.FontSize)
		result += fmt.Sprintf(`
      <w:sz w:val="%s"/>
      <w:szCs w:val="%s"/>`, sz, sz)
	}
	if ps.Bold {
		result += `
      <w:b/>`
	}
	if ps.Italic {
		result += `
      <w:i/>`
	}
	if ps.Color != "" {
		result += fmt.Sprintf(`
      <w:color w:val="%s"/>`, colorHex(ps.Color))
	}

	result += `
    </w:rPr>
  </w:style>`

	return result
}

func generateHeadingStyleXML(id, name string, ps *style.ParagraphStyle, s *style.Style) string {
	fontCJK := ps.FontCJK
	if fontCJK == "" {
		fontCJK = "黑体"
	}
	fontLatin := ps.FontLatin
	if fontLatin == "" {
		fontLatin = "Arial"
	}

	// Temporarily set fonts for generation
	orig := *ps
	ps.FontCJK = fontCJK
	ps.FontLatin = fontLatin
	result := generateParagraphStyleXML(id, name, ps, s)
	*ps = orig
	return result
}

func generateCodeStyleXML(s *style.Style) string {
	ps := &s.Styles.CodeBlock
	fontMono := ps.FontLatin
	if fontMono == "" {
		fontMono = s.Fonts.Mono
	}
	fontCJK := ps.FontCJK
	if fontCJK == "" {
		fontCJK = "等线"
	}

	return fmt.Sprintf(`
  <w:style w:type="paragraph" w:styleId="Code">
    <w:name w:val="Code"/>
    <w:pPr>
      <w:spacing w:line="%s" w:lineRule="auto"/>
    </w:pPr>
    <w:rPr>
      <w:rFonts w:ascii="%s" w:hAnsi="%s" w:eastAsia="%s"/>
      <w:sz w:val="%s"/>
      <w:szCs w:val="%s"/>
    </w:rPr>
  </w:style>
  <w:style w:type="character" w:styleId="InlineCode">
    <w:name w:val="Inline Code"/>
    <w:rPr>
      <w:rFonts w:ascii="%s" w:hAnsi="%s" w:eastAsia="%s"/>
      <w:sz w:val="%s"/>
      <w:szCs w:val="%s"/>
    </w:rPr>
  </w:style>`,
		lineSpacingVal(ps.LineSpacing),
		fontMono, fontMono, fontCJK,
		fontSizeTwips(ps.FontSize), fontSizeTwips(ps.FontSize),
		fontMono, fontMono, fontCJK,
		fontSizeTwips(ps.FontSize), fontSizeTwips(ps.FontSize))
}
