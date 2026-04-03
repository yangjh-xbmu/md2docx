package render

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func testdataDir(t *testing.T) string {
	t.Helper()
	// Find testdata relative to repo root
	dir, err := os.Getwd()
	require.NoError(t, err)
	// Walk up to find testdata
	for {
		if _, err := os.Stat(filepath.Join(dir, "testdata", "images", "test.png")); err == nil {
			return filepath.Join(dir, "testdata")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find testdata directory")
		}
		dir = parent
	}
}

func TestLoadImage_Basic(t *testing.T) {
	td := testdataDir(t)
	imgStyle := style.ImageStyle{MaxWidthPct: 80, Alignment: "center"}

	img, err := loadImage("images/test.png", td, 11906, 1800, 1800, imgStyle, 1)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, "rIdImg1", img.relID)
	assert.Equal(t, "word/media/image1.png", img.partName)
	assert.Equal(t, "image/png", img.contentType)
	assert.NotEmpty(t, img.data)
	assert.Greater(t, img.widthEmu, int64(0))
	assert.Greater(t, img.heightEmu, int64(0))
}

func TestLoadImage_Scaling(t *testing.T) {
	td := testdataDir(t)
	// Use very small max width to force scaling
	imgStyle := style.ImageStyle{MaxWidthPct: 10, Alignment: "center"}

	img, err := loadImage("images/test.png", td, 11906, 1800, 1800, imgStyle, 1)
	require.NoError(t, err)

	// The image (200x150 px) at 96 DPI = 1905000 x 1428750 EMU
	// Content width = (11906 - 1800 - 1800) * 635 = 5274310 EMU
	// 10% = 527431 EMU, so image should be scaled down
	assert.Less(t, img.widthEmu, int64(1905000), "should be scaled down")
	// Aspect ratio should be preserved
	originalRatio := float64(200) / float64(150)
	actualRatio := float64(img.widthEmu) / float64(img.heightEmu)
	assert.InDelta(t, originalRatio, actualRatio, 0.1, "aspect ratio should be preserved")
}

func TestLoadImage_NotFound(t *testing.T) {
	imgStyle := style.ImageStyle{MaxWidthPct: 80}
	_, err := loadImage("nonexistent.png", "/tmp", 11906, 1800, 1800, imgStyle, 1)
	assert.Error(t, err)
}

func TestBuildImageParagraph_XML(t *testing.T) {
	img := &imageEntry{
		relID:    "rIdImg1",
		widthEmu: 1905000,
		heightEmu: 1428750,
	}
	imgStyle := style.ImageStyle{Alignment: "center"}

	p := buildImageParagraph(img, "test image", imgStyle)
	require.NotNil(t, p)

	// Check paragraph has center alignment
	require.NotNil(t, p.PPr)
	require.NotNil(t, p.PPr.Jc)
	assert.Equal(t, "center", p.PPr.Jc.Val)

	// Marshal to XML and verify structure
	xmlData, err := xml.MarshalIndent(p, "", "  ")
	require.NoError(t, err)

	xmlStr := string(xmlData)
	assert.Contains(t, xmlStr, "w:drawing")
	assert.Contains(t, xmlStr, "wp:inline")
	assert.Contains(t, xmlStr, "a:graphic")
	assert.Contains(t, xmlStr, "pic:pic")
	assert.Contains(t, xmlStr, "r:embed=\"rIdImg1\"")
	assert.Contains(t, xmlStr, "1905000") // width
	assert.Contains(t, xmlStr, "1428750") // height
}

func TestBuildImageParagraph_DrawingStructure(t *testing.T) {
	img := &imageEntry{
		relID:    "rIdImg2",
		widthEmu: 500000,
		heightEmu: 300000,
	}
	imgStyle := style.ImageStyle{}

	p := buildImageParagraph(img, "alt text", imgStyle)
	require.NotNil(t, p)
	require.Len(t, p.Runs, 1)

	run, ok := p.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	require.Len(t, run.Content, 1)

	drawing, ok := run.Content[0].(*ooxml.Drawing)
	require.True(t, ok)
	require.NotNil(t, drawing.Inline)
	assert.Equal(t, "500000", drawing.Inline.Extent.CX)
	assert.Equal(t, "300000", drawing.Inline.Extent.CY)
	assert.Equal(t, "alt text", drawing.Inline.DocPr.Descr)
}
