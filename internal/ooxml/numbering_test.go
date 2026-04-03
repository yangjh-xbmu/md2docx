package ooxml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func TestGenerateNumberingXML_RootElement(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	assert.Contains(t, xml, `<w:numbering xmlns:w=`)
	assert.Contains(t, xml, `</w:numbering>`)
}

func TestGenerateNumberingXML_BulletAbstractNum(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `w:abstractNumId="0"`)
	// Three levels
	assert.Contains(t, xml, `w:ilvl="0"`)
	assert.Contains(t, xml, `w:ilvl="1"`)
	assert.Contains(t, xml, `w:ilvl="2"`)
	// Bullet format
	bulletCount := strings.Count(xml, `w:numFmt w:val="bullet"`)
	require.Equal(t, 3, bulletCount, "expected 3 bullet levels in abstractNum 0")
}

func TestGenerateNumberingXML_OrderedAbstractNum(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `w:abstractNumId="1"`)
	assert.Contains(t, xml, `w:numFmt w:val="decimal"`)
	assert.Contains(t, xml, `w:numFmt w:val="lowerLetter"`)
	assert.Contains(t, xml, `w:numFmt w:val="lowerRoman"`)
}

func TestGenerateNumberingXML_NumberingInstances(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	// numId 1 -> bullet (abstractNumId 0)
	assert.Contains(t, xml, `w:numId="1"`)
	// numId 2 -> ordered (abstractNumId 1)
	assert.Contains(t, xml, `w:numId="2"`)
}

func TestGenerateNumberingXML_DefaultBullets(t *testing.T) {
	s := &style.Style{} // No BulletChars set
	xml := string(GenerateNumberingXML(s))

	// Default bullets: ●, ○, ■
	assert.Contains(t, xml, `w:lvlText w:val="●"`)
	assert.Contains(t, xml, `w:lvlText w:val="○"`)
	assert.Contains(t, xml, `w:lvlText w:val="■"`)
}

func TestGenerateNumberingXML_CustomBullets(t *testing.T) {
	s := &style.Style{
		List: style.ListConfig{
			BulletChars: []string{"▸", "▹", "▪"},
		},
	}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `w:lvlText w:val="▸"`)
	assert.Contains(t, xml, `w:lvlText w:val="▹"`)
	assert.Contains(t, xml, `w:lvlText w:val="▪"`)
	// Should NOT contain default bullets
	assert.NotContains(t, xml, `w:lvlText w:val="●"`)
}

func TestGenerateNumberingXML_CustomBulletsPadding(t *testing.T) {
	// Only 1 bullet char provided; should be padded to 3 levels
	s := &style.Style{
		List: style.ListConfig{
			BulletChars: []string{"★"},
		},
	}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `w:lvlText w:val="★"`)
	// Padded with ●
	bulletCount := strings.Count(xml, `w:lvlText w:val="●"`)
	assert.Equal(t, 2, bulletCount, "expected 2 padded bullet levels")
}

func TestGenerateNumberingXML_Indentation(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	// Level 0: left=720, hanging=360
	assert.Contains(t, xml, `w:left="720" w:hanging="360"`)
	// Level 1: left=1440
	assert.Contains(t, xml, `w:left="1440" w:hanging="360"`)
	// Level 2: left=2160
	assert.Contains(t, xml, `w:left="2160" w:hanging="360"`)
}

func TestGenerateNumberingXML_OrderedLevelText(t *testing.T) {
	s := &style.Style{}
	xml := string(GenerateNumberingXML(s))

	assert.Contains(t, xml, `w:lvlText w:val="%1."`)
	assert.Contains(t, xml, `w:lvlText w:val="%2."`)
	assert.Contains(t, xml, `w:lvlText w:val="%3."`)
}

func TestMeasureToTwipsExported(t *testing.T) {
	// Test the exported wrapper
	assert.Equal(t, "567", MeasureToTwips("1cm"))
	assert.Equal(t, "1440", MeasureToTwips("1in"))
	assert.Equal(t, "", MeasureToTwips(""))
}
