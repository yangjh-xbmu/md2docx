package render

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yangjh-xbmu/md2docx/internal/parser"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func loadDefaultStyle(t *testing.T) *style.Style {
	t.Helper()
	s, err := style.Load("default")
	require.NoError(t, err, "should load default style")
	return s
}

func renderToFile(t *testing.T, body string, s *style.Style, meta map[string]any) string {
	t.Helper()
	tmp := t.TempDir()
	out := filepath.Join(tmp, "output.docx")
	ast, src := parser.ParseMarkdown(body)
	err := ToDocx(ast, src, s, meta, out, tmp)
	require.NoError(t, err, "ToDocx should succeed")
	return out
}

func assertValidDocx(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err, "output file should exist")
	assert.Greater(t, info.Size(), int64(0), "file should not be empty")

	r, err := zip.OpenReader(path)
	require.NoError(t, err, "output should be a valid zip/docx")
	defer r.Close()

	var hasDocXML bool
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			hasDocXML = true
			break
		}
	}
	assert.True(t, hasDocXML, "docx should contain word/document.xml")
}

func TestToDocx_BasicParagraph(t *testing.T) {
	s := loadDefaultStyle(t)
	out := renderToFile(t, "Hello world", s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_HeadingLevels(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "# H1\n\n## H2\n\n### H3\n\n#### H4\n\n##### H5\n\n###### H6\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_EmphasisAndBold(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "This has **bold** and *italic* and ***both***.\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_UnorderedList(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "- item one\n- item two\n- item three\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_OrderedList(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "1. first\n2. second\n3. third\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_CodeBlock(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_InlineCode(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "Use `fmt.Println` to print.\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_Link(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "Visit [Google](https://www.google.com) for more.\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_ThematicBreak(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "Before\n\n---\n\nAfter\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_EmptyContent(t *testing.T) {
	s := loadDefaultStyle(t)
	out := renderToFile(t, "", s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_CJKContent(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "# 中文标题\n\n这是一段中文正文。包含**粗体**和*斜体*。\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_WithCover(t *testing.T) {
	s, err := style.Load("academic-cn")
	require.NoError(t, err)
	meta := map[string]any{
		"title":  "Test Paper",
		"author": "Author",
		"date":   "2026-01-01",
	}
	out := renderToFile(t, "# Intro\n\nBody.\n", s, meta)
	assertValidDocx(t, out)
}

func TestToDocx_WithTOC(t *testing.T) {
	s := loadDefaultStyle(t)
	s.TOC.Enabled = true
	s.TOC.Title = "Contents"
	s.TOC.Depth = 3
	body := "# Chapter 1\n\n## Section 1.1\n\nText.\n\n# Chapter 2\n\nMore text.\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_WithNumbering(t *testing.T) {
	s := loadDefaultStyle(t)
	s.HeadingNumbering.Enabled = true
	s.HeadingNumbering.Formats = map[int]string{
		1: "{1}",
		2: "{1}.{2}",
		3: "{1}.{2}.{3}",
	}
	body := "# First\n\n## Sub A\n\n## Sub B\n\n# Second\n\n## Sub C\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_MixedContent(t *testing.T) {
	s := loadDefaultStyle(t)
	body := `# Title

A paragraph with **bold**, *italic*, and ` + "`code`" + `.

## Lists

- bullet 1
- bullet 2

1. ordered 1
2. ordered 2

## Code

` + "```python\nprint('hello')\n```" + `

---

End.
`
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_TableContent(t *testing.T) {
	s := loadDefaultStyle(t)
	body := `| Name | Value |
|------|-------|
| A    | 1     |
| B    | 2     |
`
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_LandscapeOrientation(t *testing.T) {
	s := loadDefaultStyle(t)
	s.Page.Orientation = "landscape"
	out := renderToFile(t, "# Landscape\n\nContent.\n", s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_HeaderFooter(t *testing.T) {
	s := loadDefaultStyle(t)
	s.Header.Left = "Header Left"
	s.Header.Right = "{title}"
	s.Footer.Center = "Page {page}"
	meta := map[string]any{"title": "My Doc"}
	out := renderToFile(t, "# Test\n\nContent.\n", s, meta)
	assertValidDocx(t, out)

	// Verify header/footer XML parts exist in the zip
	r, err := zip.OpenReader(out)
	require.NoError(t, err)
	defer r.Close()

	var hasHeader, hasFooter bool
	for _, f := range r.File {
		if f.Name == "word/header1.xml" {
			hasHeader = true
		}
		if f.Name == "word/footer1.xml" {
			hasFooter = true
		}
	}
	assert.True(t, hasHeader, "docx should contain header XML")
	assert.True(t, hasFooter, "docx should contain footer XML")
}

func TestToDocx_DisabledFeatures(t *testing.T) {
	s := loadDefaultStyle(t)
	s.TOC.Enabled = false
	s.Cover.Enabled = false
	s.HeadingNumbering.Enabled = false
	out := renderToFile(t, "# Simple\n\nJust text.\n", s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_NestedList(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "- level 1\n  - level 2\n    - level 3\n- back to 1\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_AutoLink(t *testing.T) {
	s := loadDefaultStyle(t)
	body := "Check https://example.com for details.\n"
	out := renderToFile(t, body, s, nil)
	assertValidDocx(t, out)
}

func TestToDocx_MultipleStyles(t *testing.T) {
	for _, styleName := range []string{"default", "academic-cn", "simple"} {
		t.Run(styleName, func(t *testing.T) {
			s, err := style.Load(styleName)
			require.NoError(t, err)
			out := renderToFile(t, "# Title\n\nParagraph.\n", s, nil)
			assertValidDocx(t, out)
		})
	}
}
