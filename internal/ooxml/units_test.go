package ooxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFontSizeTwips(t *testing.T) {
	assert.Equal(t, "24", fontSizeTwips("12pt"))
	assert.Equal(t, "44", fontSizeTwips("22pt"))
	assert.Equal(t, "24", fontSizeTwips(""))
}

func TestLineSpacingVal(t *testing.T) {
	assert.Equal(t, "360", lineSpacingVal(1.5))
	assert.Equal(t, "240", lineSpacingVal(1.0))
	assert.Equal(t, "480", lineSpacingVal(2.0))
}

func TestMeasureToTwips(t *testing.T) {
	assert.Equal(t, "567", measureToTwips("1cm"))
	assert.Equal(t, "1134", measureToTwips("2cm"))
	assert.Equal(t, "1440", measureToTwips("1in"))
	assert.Equal(t, "240", measureToTwips("12pt"))
	assert.Equal(t, "", measureToTwips(""))
}

func TestPageSizeTwips(t *testing.T) {
	w, h := PageSizeTwips("A4")
	assert.Equal(t, "11906", w)
	assert.Equal(t, "16838", h)

	w, h = PageSizeTwips("Letter")
	assert.Equal(t, "12240", w)
	assert.Equal(t, "15840", h)
}

func TestColorHex(t *testing.T) {
	assert.Equal(t, "FF0000", colorHex("#FF0000"))
	assert.Equal(t, "000000", colorHex("000000"))
}
