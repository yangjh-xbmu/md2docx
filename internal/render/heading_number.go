package render

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// HeadingNumberer tracks heading counters and generates numbering prefixes.
type HeadingNumberer struct {
	enabled  bool
	formats  map[int]string
	counters [6]int // counters for levels 1-6
}

// NewHeadingNumberer creates a numbering tracker from style config.
func NewHeadingNumberer(hn style.HeadingNumbering) *HeadingNumberer {
	return &HeadingNumberer{
		enabled: hn.Enabled,
		formats: hn.Formats,
	}
}

// FormatNumber increments the counter for the given level (1-based)
// and returns the formatted numbering prefix. Returns "" if disabled.
func (h *HeadingNumberer) FormatNumber(level int) string {
	if !h.enabled || level < 1 || level > 6 {
		return ""
	}

	// Increment current level counter
	h.counters[level-1]++

	// Reset all deeper level counters
	for i := level; i < 6; i++ {
		h.counters[i] = 0
	}

	format, ok := h.formats[level]
	if !ok {
		return ""
	}

	return h.applyFormat(format, level)
}

// applyFormat replaces {1}, {2}, etc. with actual counter values.
// Also supports {1:zh} for Chinese numbering.
func (h *HeadingNumberer) applyFormat(format string, level int) string {
	result := format
	for i := 1; i <= level; i++ {
		// Chinese format: {1:zh}
		zhPlaceholder := fmt.Sprintf("{%d:zh}", i)
		if strings.Contains(result, zhPlaceholder) {
			result = strings.ReplaceAll(result, zhPlaceholder, toChineseNumber(h.counters[i-1]))
		}
		// Standard format: {1}
		placeholder := fmt.Sprintf("{%d}", i)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%d", h.counters[i-1]))
	}
	return result
}

// HasExistingNumber checks if the heading text already starts with a numbering
// pattern that matches the generated prefix, to avoid duplication.
var existingNumberRe = regexp.MustCompile(`^[\d]+(?:\.[\d]+)*\.?\s+`)

// StripExistingNumber removes existing numbering from heading text if present.
func StripExistingNumber(text string) string {
	return existingNumberRe.ReplaceAllString(text, "")
}

// toChineseNumber converts an integer to Chinese character representation.
func toChineseNumber(n int) string {
	chars := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	if n <= 10 {
		return chars[n]
	}
	if n < 20 {
		return "十" + chars[n-10]
	}
	if n < 100 {
		tens := n / 10
		ones := n % 10
		if ones == 0 {
			return chars[tens] + "十"
		}
		return chars[tens] + "十" + chars[ones]
	}
	return fmt.Sprintf("%d", n)
}
