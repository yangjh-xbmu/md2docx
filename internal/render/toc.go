package render

import (
	"fmt"

	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// GenerateTOCElements creates the OOXML elements for a Table of Contents.
// Returns a slice of Paragraph elements to be prepended to the document body.
func GenerateTOCElements(toc style.TOCConfig) []any {
	if !toc.Enabled {
		return nil
	}

	var elements []any

	// TOC title paragraph
	if toc.Title != "" {
		titlePara := &ooxml.Paragraph{
			PPr: &ooxml.ParagraphProperties{
				PStyle: &ooxml.PStyle{Val: "Heading1"},
			},
			Runs: []any{
				&ooxml.Run{
					Content: []any{
						&ooxml.Text{Space: "preserve", Value: toc.Title},
					},
				},
			},
		}
		elements = append(elements, titlePara)
	}

	// TOC field code: begin + instrText + separate + end
	// The instrText controls TOC behavior:
	//   \o "1-3" = heading levels 1-3
	//   \h = hyperlinks
	//   \z = hide tab leaders in web view
	//   \u = use applied paragraph outline level
	instrText := fmt.Sprintf(` TOC \o "1-%d" \h \z \u `, toc.Depth)

	tocPara := &ooxml.Paragraph{
		Runs: []any{
			// Begin field
			&ooxml.Run{
				Content: []any{
					&ooxml.FldChar{FldCharType: "begin"},
				},
			},
			// Instruction
			&ooxml.Run{
				Content: []any{
					&ooxml.InstrText{Space: "preserve", Value: instrText},
				},
			},
			// Separate
			&ooxml.Run{
				Content: []any{
					&ooxml.FldChar{FldCharType: "separate"},
				},
			},
			// Placeholder text (Word replaces this on update)
			&ooxml.Run{
				Content: []any{
					&ooxml.Text{Space: "preserve", Value: "请更新目录（Ctrl+A 后 F9）"},
				},
			},
			// End field
			&ooxml.Run{
				Content: []any{
					&ooxml.FldChar{FldCharType: "end"},
				},
			},
		},
	}
	elements = append(elements, tocPara)

	// Page break after TOC
	if toc.PageBreakAfter {
		breakPara := &ooxml.Paragraph{
			Runs: []any{
				&ooxml.Run{
					Content: []any{
						&ooxml.Break{Type: "page"},
					},
				},
			},
		}
		elements = append(elements, breakPara)
	}

	return elements
}
