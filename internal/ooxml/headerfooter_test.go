package ooxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func TestGenerateHeaderXML_Empty(t *testing.T) {
	hf := style.HeaderFooterConfig{FontSize: "9pt"}
	result := GenerateHeaderXML(hf, nil)
	assert.Nil(t, result)
}

func TestGenerateHeaderXML_RightTitle(t *testing.T) {
	hf := style.HeaderFooterConfig{
		Right:    "{title}",
		FontSize: "9pt",
	}
	meta := map[string]any{"title": "测试文档"}
	result := GenerateHeaderXML(hf, meta)
	assert.NotNil(t, result)
	xml := string(result)
	assert.Contains(t, xml, "w:hdr")
	assert.Contains(t, xml, "测试文档")
	assert.Contains(t, xml, `w:jc w:val="right"`)
}

func TestGenerateFooterXML_CenterPage(t *testing.T) {
	hf := style.HeaderFooterConfig{
		Center:   "{page}",
		FontSize: "9pt",
	}
	result := GenerateFooterXML(hf, nil)
	assert.NotNil(t, result)
	xml := string(result)
	assert.Contains(t, xml, "w:ftr")
	assert.Contains(t, xml, "PAGE")
	assert.Contains(t, xml, `fldCharType="begin"`)
	assert.Contains(t, xml, `fldCharType="end"`)
	assert.Contains(t, xml, `w:jc w:val="center"`)
}

func TestGenerateHeaderXML_LeftAndRight(t *testing.T) {
	hf := style.HeaderFooterConfig{
		Left:     "文档标题",
		Right:    "{page}",
		FontSize: "9pt",
	}
	result := GenerateHeaderXML(hf, nil)
	assert.NotNil(t, result)
	xml := string(result)
	assert.Contains(t, xml, "文档标题")
	assert.Contains(t, xml, "PAGE")
	// Should have tab stops for multi-position
	assert.Contains(t, xml, "w:tabs")
}

func TestResolveTemplate(t *testing.T) {
	meta := map[string]any{
		"title":  "My Title",
		"author": "Author Name",
	}

	assert.Equal(t, "My Title", resolveTemplate("{title}", meta))
	assert.Equal(t, "By Author Name", resolveTemplate("By {author}", meta))
	assert.Equal(t, "{page}", resolveTemplate("{page}", meta)) // not in meta, preserved
	assert.Equal(t, "", resolveTemplate("", meta))
}

func TestEscapeXML(t *testing.T) {
	assert.Equal(t, "&amp;&lt;&gt;&quot;", escapeXML("&<>\""))
	assert.Equal(t, "normal text", escapeXML("normal text"))
}
