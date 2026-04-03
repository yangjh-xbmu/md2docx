package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func TestGenerateTOCElements_Disabled(t *testing.T) {
	elems := GenerateTOCElements(style.TOCConfig{Enabled: false})
	assert.Nil(t, elems)
}

func TestGenerateTOCElements_Enabled(t *testing.T) {
	elems := GenerateTOCElements(style.TOCConfig{
		Enabled:        true,
		Title:          "目录",
		Depth:          3,
		PageBreakAfter: true,
	})

	// Should have: title paragraph + TOC field paragraph + page break paragraph
	assert.Len(t, elems, 3)

	// First element: title paragraph
	titlePara, ok := elems[0].(*ooxml.Paragraph)
	assert.True(t, ok)
	assert.Equal(t, "Heading1", titlePara.PPr.PStyle.Val)

	// Second element: TOC field paragraph with fldChar
	tocPara, ok := elems[1].(*ooxml.Paragraph)
	assert.True(t, ok)
	assert.Len(t, tocPara.Runs, 5) // begin, instrText, separate, placeholder, end

	// Verify fldChar begin
	beginRun, ok := tocPara.Runs[0].(*ooxml.Run)
	assert.True(t, ok)
	fldChar, ok := beginRun.Content[0].(*ooxml.FldChar)
	assert.True(t, ok)
	assert.Equal(t, "begin", fldChar.FldCharType)

	// Verify instrText contains TOC command
	instrRun, ok := tocPara.Runs[1].(*ooxml.Run)
	assert.True(t, ok)
	instrText, ok := instrRun.Content[0].(*ooxml.InstrText)
	assert.True(t, ok)
	assert.Contains(t, instrText.Value, `TOC \o "1-3"`)

	// Third element: page break
	breakPara, ok := elems[2].(*ooxml.Paragraph)
	assert.True(t, ok)
	breakRun, ok := breakPara.Runs[0].(*ooxml.Run)
	assert.True(t, ok)
	br, ok := breakRun.Content[0].(*ooxml.Break)
	assert.True(t, ok)
	assert.Equal(t, "page", br.Type)
}

func TestGenerateTOCElements_NoPageBreak(t *testing.T) {
	elems := GenerateTOCElements(style.TOCConfig{
		Enabled:        true,
		Title:          "目录",
		Depth:          2,
		PageBreakAfter: false,
	})

	// Should have: title + TOC field (no page break)
	assert.Len(t, elems, 2)
}

func TestGenerateTOCElements_NoTitle(t *testing.T) {
	elems := GenerateTOCElements(style.TOCConfig{
		Enabled: true,
		Title:   "",
		Depth:   3,
	})

	// Should have: only TOC field paragraph
	assert.Len(t, elems, 1)
}
