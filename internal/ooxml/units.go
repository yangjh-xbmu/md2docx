package ooxml

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// OOXML uses twips (1/20 pt) for most measurements.
// 1 pt = 20 twips
// 1 cm = 567 twips
// 1 mm = 56.7 twips
// 1 inch = 1440 twips
// 1 em ≈ font size in pt (approximate)
// Font sizes in OOXML are in half-points (1 pt = 2 half-points)

// fontSizeTwips converts "12pt" to half-points string (e.g. "24").
func fontSizeTwips(size string) string {
	if size == "" {
		return "24" // default 12pt
	}
	pt := parsePt(size)
	return strconv.Itoa(int(pt * 2))
}

// lineSpacingVal converts a multiplier (e.g. 1.5) to OOXML line spacing (240 * multiplier).
func lineSpacingVal(multiplier float64) string {
	if multiplier <= 0 {
		multiplier = 1.0
	}
	return strconv.Itoa(int(math.Round(240 * multiplier)))
}

// twipsFromPt converts "12pt" to twips string.
func twipsFromPt(s string) string {
	if s == "" {
		return ""
	}
	pt := parsePt(s)
	return strconv.Itoa(int(math.Round(pt * 20)))
}

// measureToTwips converts "2cm", "25.4mm", "1in" to twips string.
func measureToTwips(s string) string {
	if s == "" {
		return ""
	}
	s = strings.TrimSpace(s)

	if strings.HasSuffix(s, "cm") {
		val := parseFloat(strings.TrimSuffix(s, "cm"))
		return strconv.Itoa(int(math.Round(val * 567)))
	}
	if strings.HasSuffix(s, "mm") {
		val := parseFloat(strings.TrimSuffix(s, "mm"))
		return strconv.Itoa(int(math.Round(val * 56.7)))
	}
	if strings.HasSuffix(s, "in") {
		val := parseFloat(strings.TrimSuffix(s, "in"))
		return strconv.Itoa(int(math.Round(val * 1440)))
	}
	if strings.HasSuffix(s, "pt") {
		val := parseFloat(strings.TrimSuffix(s, "pt"))
		return strconv.Itoa(int(math.Round(val * 20)))
	}

	// Assume twips if no unit
	return s
}

// indentTwips handles "2em" (relative to font size) or absolute measurements.
func indentTwips(s string, fontSize string) string {
	if s == "" {
		return ""
	}
	if strings.HasSuffix(s, "em") {
		emVal := parseFloat(strings.TrimSuffix(s, "em"))
		ptSize := parsePt(fontSize)
		if ptSize <= 0 {
			ptSize = 12 // default
		}
		twips := emVal * ptSize * 20
		return strconv.Itoa(int(math.Round(twips)))
	}
	return measureToTwips(s)
}

// alignmentVal converts friendly names to OOXML values.
func alignmentVal(s string) string {
	switch strings.ToLower(s) {
	case "justify", "both":
		return "both"
	case "center":
		return "center"
	case "right":
		return "right"
	default:
		return "left"
	}
}

// colorHex strips # from hex color.
func colorHex(s string) string {
	return strings.TrimPrefix(s, "#")
}

// PageSizeTwips returns (width, height) in twips for a named page size.
func PageSizeTwips(size string) (string, string) {
	switch strings.ToUpper(size) {
	case "LETTER":
		return "12240", "15840" // 8.5 x 11 in
	case "LEGAL":
		return "12240", "20160" // 8.5 x 14 in
	case "B5":
		return "10318", "14572" // 182 x 257 mm
	case "A3":
		return "16838", "23812" // 297 x 420 mm
	default: // A4
		return "11906", "16838" // 210 x 297 mm
	}
}

// ColorHex strips # from hex color (exported).
func ColorHex(s string) string {
	return colorHex(s)
}

// AlignmentVal converts friendly names to OOXML values (exported).
func AlignmentVal(s string) string {
	return alignmentVal(s)
}

// FontSizeHalfPoints converts "12pt" to half-points string (exported).
func FontSizeHalfPoints(size string) string {
	return fontSizeTwips(size)
}

func parsePt(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "pt")
	return parseFloat(s)
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return v
}

// EmuFromPx converts pixels to EMU (English Metric Units).
// 1 inch = 914400 EMU, assuming 96 DPI: 1 px = 9525 EMU
func EmuFromPx(px int) int64 {
	return int64(px) * 9525
}

// FormatEmu formats an EMU value as string.
func FormatEmu(emu int64) string {
	return fmt.Sprintf("%d", emu)
}
