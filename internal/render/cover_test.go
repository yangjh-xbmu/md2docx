package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func TestGenerateCoverElements_Disabled(t *testing.T) {
	elems := GenerateCoverElements(style.CoverConfig{Enabled: false}, nil)
	assert.Nil(t, elems)
}

func TestGenerateCoverElements_NoLayout(t *testing.T) {
	elems := GenerateCoverElements(style.CoverConfig{Enabled: true, Layout: nil}, nil)
	assert.Nil(t, elems)
}

func TestGenerateCoverElements_BasicLayout(t *testing.T) {
	cfg := style.CoverConfig{
		Enabled:        true,
		PageBreakAfter: true,
		Layout: []style.CoverElement{
			{Type: "spacer", Height: "120pt"},
			{Type: "text", Source: "title", FontCJK: "黑体", FontSize: "28pt", Bold: true, Alignment: "center"},
			{Type: "text", Source: "author", FontSize: "16pt", Alignment: "center"},
			{Type: "text", Source: "date", FontSize: "14pt", Alignment: "center"},
		},
	}
	meta := map[string]any{
		"title":  "测试论文",
		"author": "张三",
		"date":   "2026-04-03",
	}

	elems := GenerateCoverElements(cfg, meta)
	require.NotNil(t, elems)

	// spacer + 3 text + page break = 5 elements
	assert.Len(t, elems, 5)

	// First: spacer (empty paragraph with spacing)
	spacer, ok := elems[0].(*ooxml.Paragraph)
	require.True(t, ok)
	assert.NotNil(t, spacer.PPr.Spacing)
	assert.Equal(t, "2400", spacer.PPr.Spacing.Before) // 120pt = 120*20 = 2400 twips

	// Second: title paragraph
	titlePara, ok := elems[1].(*ooxml.Paragraph)
	require.True(t, ok)
	assert.Equal(t, "center", titlePara.PPr.Jc.Val)
	require.Len(t, titlePara.Runs, 1)
	titleRun, ok := titlePara.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	assert.NotNil(t, titleRun.RPr.Bold)
	assert.Equal(t, "黑体", titleRun.RPr.RFonts.EastAsia)
	assert.Equal(t, "56", titleRun.RPr.Sz.Val) // 28pt * 2 = 56 half-points
	titleText, ok := titleRun.Content[0].(*ooxml.Text)
	require.True(t, ok)
	assert.Equal(t, "测试论文", titleText.Value)

	// Third: author paragraph
	authorPara, ok := elems[2].(*ooxml.Paragraph)
	require.True(t, ok)
	authorRun, ok := authorPara.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	authorText, ok := authorRun.Content[0].(*ooxml.Text)
	require.True(t, ok)
	assert.Equal(t, "张三", authorText.Value)

	// Last: page break
	breakPara, ok := elems[4].(*ooxml.Paragraph)
	require.True(t, ok)
	breakRun, ok := breakPara.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	br, ok := breakRun.Content[0].(*ooxml.Break)
	require.True(t, ok)
	assert.Equal(t, "page", br.Type)
}

func TestGenerateCoverElements_NoPageBreak(t *testing.T) {
	cfg := style.CoverConfig{
		Enabled:        true,
		PageBreakAfter: false,
		Layout: []style.CoverElement{
			{Type: "text", Source: "title", FontSize: "28pt", Alignment: "center"},
		},
	}
	meta := map[string]any{"title": "Test"}

	elems := GenerateCoverElements(cfg, meta)
	// Only the title paragraph, no page break
	assert.Len(t, elems, 1)
}

func TestGenerateCoverElements_MissingMetaField(t *testing.T) {
	cfg := style.CoverConfig{
		Enabled: true,
		Layout: []style.CoverElement{
			{Type: "text", Source: "title", FontSize: "16pt", Alignment: "center"},
			{Type: "text", Source: "institution", FontSize: "14pt", Alignment: "center"},
		},
	}
	meta := map[string]any{"title": "Test"}

	elems := GenerateCoverElements(cfg, meta)
	// title only, institution skipped (missing in meta)
	assert.Len(t, elems, 1)
}

func TestGenerateCoverElements_LiteralText(t *testing.T) {
	cfg := style.CoverConfig{
		Enabled: true,
		Layout: []style.CoverElement{
			{Type: "text", Source: "literal:机密文档", FontSize: "12pt", Alignment: "right"},
		},
	}

	elems := GenerateCoverElements(cfg, nil)
	require.Len(t, elems, 1)
	para, ok := elems[0].(*ooxml.Paragraph)
	require.True(t, ok)
	run, ok := para.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	text, ok := run.Content[0].(*ooxml.Text)
	require.True(t, ok)
	assert.Equal(t, "机密文档", text.Value)
	assert.Equal(t, "right", para.PPr.Jc.Val)
}
