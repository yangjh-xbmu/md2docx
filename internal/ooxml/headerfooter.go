package ooxml

import (
	"fmt"
	"strings"

	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// GenerateHeaderXML generates word/header1.xml from header config and frontmatter.
func GenerateHeaderXML(hf style.HeaderFooterConfig, meta map[string]any) []byte {
	content := generateHeaderFooterContent(hf, meta)
	if content == "" {
		return nil
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
` + content + `
</w:hdr>`

	return []byte(xml)
}

// GenerateFooterXML generates word/footer1.xml from footer config and frontmatter.
func GenerateFooterXML(hf style.HeaderFooterConfig, meta map[string]any) []byte {
	content := generateHeaderFooterContent(hf, meta)
	if content == "" {
		return nil
	}

	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
` + content + `
</w:ftr>`

	return []byte(xml)
}

// generateHeaderFooterContent builds the paragraph XML for a header or footer.
// Uses tab stops for left/center/right positioning.
func generateHeaderFooterContent(hf style.HeaderFooterConfig, meta map[string]any) string {
	left := resolveTemplate(hf.Left, meta)
	center := resolveTemplate(hf.Center, meta)
	right := resolveTemplate(hf.Right, meta)

	if left == "" && center == "" && right == "" {
		return ""
	}

	fontSize := hf.FontSize
	if fontSize == "" {
		fontSize = "9pt"
	}
	szVal := fontSizeTwips(fontSize)

	// Build paragraph with tab stops for center (4513 twips = center of A4)
	// and right (9026 twips = right margin of A4)
	var runs string

	if left != "" {
		runs += buildRunXML(left, szVal)
	}

	if center != "" {
		if left != "" {
			// Tab to center
			runs += fmt.Sprintf(`  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:tab/>
  </w:r>
`, szVal, szVal)
		}
		runs += buildRunXML(center, szVal)
	}

	if right != "" {
		if left != "" || center != "" {
			// Tab to right
			runs += fmt.Sprintf(`  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:tab/>
  </w:r>
`, szVal, szVal)
		}
		runs += buildRunXML(right, szVal)
	}

	// Determine alignment if only one position is used
	jc := ""
	if left == "" && right == "" && center != "" {
		jc = `<w:jc w:val="center"/>`
	} else if left == "" && center == "" && right != "" {
		jc = `<w:jc w:val="right"/>`
	} else if left != "" && (center != "" || right != "") {
		// Use tab stops for multi-position layout
		jc = `<w:tabs>
        <w:tab w:val="center" w:pos="4513"/>
        <w:tab w:val="right" w:pos="9026"/>
      </w:tabs>`
	}

	return fmt.Sprintf(`  <w:p>
    <w:pPr>
      %s
    </w:pPr>
%s  </w:p>`, jc, runs)
}

// buildRunXML creates run XML, handling {page} as a PAGE field code.
func buildRunXML(text, szVal string) string {
	if text == "{page}" {
		// PAGE field code
		return fmt.Sprintf(`  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:fldChar w:fldCharType="begin"/>
  </w:r>
  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:instrText xml:space="preserve"> PAGE </w:instrText>
  </w:r>
  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:fldChar w:fldCharType="separate"/>
  </w:r>
  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:t>1</w:t>
  </w:r>
  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:fldChar w:fldCharType="end"/>
  </w:r>
`, szVal, szVal, szVal, szVal, szVal, szVal, szVal, szVal, szVal, szVal)
	}

	// Plain text run
	return fmt.Sprintf(`  <w:r>
    <w:rPr><w:sz w:val="%s"/><w:szCs w:val="%s"/></w:rPr>
    <w:t xml:space="preserve">%s</w:t>
  </w:r>
`, szVal, szVal, escapeXML(text))
}

// resolveTemplate replaces {key} placeholders with frontmatter values.
// {page} is preserved for field code generation.
func resolveTemplate(template string, meta map[string]any) string {
	if template == "" {
		return ""
	}

	result := template
	for key, val := range meta {
		placeholder := "{" + key + "}"
		if s, ok := val.(string); ok {
			result = strings.ReplaceAll(result, placeholder, s)
		}
	}
	return result
}

// escapeXML escapes special XML characters in text content.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
