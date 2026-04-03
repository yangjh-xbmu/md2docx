package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func TestHeadingNumberer_Disabled(t *testing.T) {
	hn := NewHeadingNumberer(style.HeadingNumbering{Enabled: false})
	assert.Equal(t, "", hn.FormatNumber(1))
}

func TestHeadingNumberer_BasicNumbering(t *testing.T) {
	hn := NewHeadingNumberer(style.HeadingNumbering{
		Enabled: true,
		Formats: map[int]string{
			1: "{1}",
			2: "{1}.{2}",
			3: "{1}.{2}.{3}",
		},
	})

	assert.Equal(t, "1", hn.FormatNumber(1))
	assert.Equal(t, "1.1", hn.FormatNumber(2))
	assert.Equal(t, "1.2", hn.FormatNumber(2))
	assert.Equal(t, "1.2.1", hn.FormatNumber(3))
	assert.Equal(t, "2", hn.FormatNumber(1))
	assert.Equal(t, "2.1", hn.FormatNumber(2))
}

func TestHeadingNumberer_ChineseFormat(t *testing.T) {
	hn := NewHeadingNumberer(style.HeadingNumbering{
		Enabled: true,
		Formats: map[int]string{
			1: "第{1:zh}章",
		},
	})

	assert.Equal(t, "第一章", hn.FormatNumber(1))
	assert.Equal(t, "第二章", hn.FormatNumber(1))
	assert.Equal(t, "第三章", hn.FormatNumber(1))
}

func TestHeadingNumberer_NoFormatForLevel(t *testing.T) {
	hn := NewHeadingNumberer(style.HeadingNumbering{
		Enabled: true,
		Formats: map[int]string{
			1: "{1}",
		},
	})

	assert.Equal(t, "1", hn.FormatNumber(1))
	assert.Equal(t, "", hn.FormatNumber(2)) // no format for level 2
}

func TestHeadingNumberer_ResetDeeperCounters(t *testing.T) {
	hn := NewHeadingNumberer(style.HeadingNumbering{
		Enabled: true,
		Formats: map[int]string{
			1: "{1}",
			2: "{1}.{2}",
		},
	})

	hn.FormatNumber(1) // 1
	hn.FormatNumber(2) // 1.1
	hn.FormatNumber(2) // 1.2
	hn.FormatNumber(1) // 2 (resets level 2 counter)
	assert.Equal(t, "2.1", hn.FormatNumber(2)) // should be 2.1, not 2.3
}

func TestStripExistingNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.1 Title", "Title"},
		{"1.2.3 Deep Title", "Deep Title"},
		{"1 Top Level", "Top Level"},
		{"No number here", "No number here"},
		{"1. Ordered", "Ordered"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, StripExistingNumber(tt.input))
	}
}

func TestToChineseNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1, "一"},
		{5, "五"},
		{10, "十"},
		{11, "十一"},
		{20, "二十"},
		{25, "二十五"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, toChineseNumber(tt.input))
	}
}
