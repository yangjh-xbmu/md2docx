package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// GenerateCoverElements creates the OOXML elements for a cover page.
// Returns a slice of Paragraph elements to be prepended to the document body.
func GenerateCoverElements(cover style.CoverConfig, meta map[string]any) []any {
	if !cover.Enabled || len(cover.Layout) == 0 {
		return nil
	}

	var elements []any

	for _, elem := range cover.Layout {
		switch elem.Type {
		case "spacer":
			p := buildSpacer(elem.Height)
			elements = append(elements, p)
		case "text":
			text := resolveText(elem.Source, meta)
			if text == "" {
				continue
			}
			p := buildCoverText(text, elem)
			elements = append(elements, p)
		}
	}

	if len(elements) == 0 {
		return nil
	}

	if cover.PageBreakAfter {
		elements = append(elements, &ooxml.Paragraph{
			Runs: []any{
				&ooxml.Run{
					Content: []any{
						&ooxml.Break{Type: "page"},
					},
				},
			},
		})
	}

	return elements
}

// resolveText resolves a cover element source to a text value.
// "literal:xxx" returns "xxx" directly; otherwise looks up meta[source].
func resolveText(source string, meta map[string]any) string {
	if strings.HasPrefix(source, "literal:") {
		return strings.TrimPrefix(source, "literal:")
	}
	if meta == nil {
		return ""
	}
	v, ok := meta[source]
	if !ok {
		return ""
	}
	if t, ok := v.(time.Time); ok {
		return t.Format("2006-01-02")
	}
	return fmt.Sprintf("%v", v)
}

func buildSpacer(height string) *ooxml.Paragraph {
	return &ooxml.Paragraph{
		PPr: &ooxml.ParagraphProperties{
			Spacing: &ooxml.Spacing{
				Before: ptToTwips(height),
			},
		},
	}
}

func buildCoverText(text string, elem style.CoverElement) *ooxml.Paragraph {
	rpr := &ooxml.RunProperties{}

	if elem.Bold {
		rpr.Bold = &ooxml.Empty{}
	}

	if elem.FontSize != "" {
		halfPts := fontSizeHalfPts(elem.FontSize)
		rpr.Sz = &ooxml.SzVal{Val: halfPts}
		rpr.SzCs = &ooxml.SzVal{Val: halfPts}
	}

	if elem.FontCJK != "" {
		if rpr.RFonts == nil {
			rpr.RFonts = &ooxml.RunFonts{}
		}
		rpr.RFonts.EastAsia = elem.FontCJK
	}

	ppr := &ooxml.ParagraphProperties{}
	if elem.Alignment != "" {
		ppr.Jc = &ooxml.Justification{Val: elem.Alignment}
	}

	return &ooxml.Paragraph{
		PPr: ppr,
		Runs: []any{
			&ooxml.Run{
				RPr: rpr,
				Content: []any{
					&ooxml.Text{Space: "preserve", Value: text},
				},
			},
		},
	}
}

// ptToTwips converts "120pt" to twips string. 1pt = 20 twips.
func ptToTwips(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "pt") {
		val, err := strconv.ParseFloat(strings.TrimSuffix(s, "pt"), 64)
		if err != nil {
			return "0"
		}
		return strconv.Itoa(int(math.Round(val * 20)))
	}
	return ooxml.MeasureToTwips(s)
}

// fontSizeHalfPts converts "28pt" to half-points string (e.g. "56").
func fontSizeHalfPts(size string) string {
	size = strings.TrimSpace(size)
	pt, err := strconv.ParseFloat(strings.TrimSuffix(size, "pt"), 64)
	if err != nil {
		return "24" // default 12pt
	}
	return strconv.Itoa(int(pt * 2))
}
