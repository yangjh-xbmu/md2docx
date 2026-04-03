package cli

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper: reset global flags that cobra persists between tests.
func resetFlags() {
	outputOpt = ""
	styleOpt = ""
	noCover = false
	noTOC = false
	noNumber = false
}

func writeTempMD(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "input.md")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestConvertBasicMarkdown(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Hello\n\nA paragraph.\n")
	out := filepath.Join(tmp, "out.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out})

	err := rootCmd.Execute()
	require.NoError(t, err, "convert should succeed")

	// output file must exist
	info, err := os.Stat(out)
	require.NoError(t, err, "output file should exist")
	assert.Greater(t, info.Size(), int64(0))

	// must be a valid zip (docx is a zip)
	r, err := zip.OpenReader(out)
	require.NoError(t, err, "output should be a valid zip/docx")
	defer r.Close()

	// must contain document.xml
	var hasDoc bool
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			hasDoc = true
			break
		}
	}
	assert.True(t, hasDoc, "docx should contain word/document.xml")
}

func TestConvertWithStyle(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n\nBody text.\n")
	out := filepath.Join(tmp, "styled.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--style", "academic-cn"})

	err := rootCmd.Execute()
	require.NoError(t, err, "convert with academic-cn style should succeed")

	_, err = os.Stat(out)
	require.NoError(t, err, "output file should exist")
}

func TestConvertNonexistentFile(t *testing.T) {
	resetFlags()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"/tmp/does_not_exist_abc123.md", "-o", "/tmp/nope.docx"})

	err := rootCmd.Execute()
	assert.Error(t, err, "should error on nonexistent input")
}

func TestConvertNonexistentStyle(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n")
	out := filepath.Join(tmp, "out.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--style", "nonexistent"})

	err := rootCmd.Execute()
	require.Error(t, err, "should error on nonexistent style")
	assert.Contains(t, err.Error(), "不存在")
}

func TestConvertNoArgs(t *testing.T) {
	resetFlags()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()
	// No args shows help, should not error
	assert.NoError(t, err, "no args should print help without error")
}

func TestConvertWithNoCoverFlag(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n\nContent.\n")
	out := filepath.Join(tmp, "no_cover.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--no-cover"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err, "output should exist")
}

func TestConvertWithNoTOCFlag(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n\nContent.\n")
	out := filepath.Join(tmp, "no_toc.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--no-toc"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err, "output should exist")
}

func TestConvertWithNoNumberingFlag(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n\n## Sub\n\nContent.\n")
	out := filepath.Join(tmp, "no_num.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--no-numbering"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err, "output should exist")
}

func TestConvertAllFlagsCombined(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	md := writeTempMD(t, tmp, "# Title\n\n## Sub\n\nContent.\n")
	out := filepath.Join(tmp, "all_flags.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out, "--no-cover", "--no-toc", "--no-numbering"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	r, err := zip.OpenReader(out)
	require.NoError(t, err)
	defer r.Close()
}

func TestConvertDefaultOutputPath(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	// Create md file named "sample.md"
	mdPath := filepath.Join(tmp, "sample.md")
	require.NoError(t, os.WriteFile(mdPath, []byte("# Test\n"), 0o644))

	// Don't specify -o; the code will try ~/Desktop/sample.docx or same dir
	// We can't easily test the Desktop path, but we can verify it doesn't error
	// by overriding HOME to tmp so it falls back to source dir
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{mdPath})

	err := rootCmd.Execute()
	require.NoError(t, err)

	// Should create sample.docx somewhere; check both possible locations
	desktopPath := filepath.Join(tmp, "Desktop", "sample.docx")
	sameDirPath := filepath.Join(tmp, "sample.docx")
	_, err1 := os.Stat(desktopPath)
	_, err2 := os.Stat(sameDirPath)
	assert.True(t, err1 == nil || err2 == nil, "output docx should exist in Desktop or source dir")
}

func TestConvertWithFrontmatterStyle(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	content := `---
title: Test Doc
style: simple
---
# Hello

World.
`
	md := writeTempMD(t, tmp, content)
	out := filepath.Join(tmp, "fm.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err)
}

// --- Merge subcommand tests ---

func TestMergeBasic(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()

	// Create two small md files
	f1 := filepath.Join(tmp, "ch1.md")
	f2 := filepath.Join(tmp, "ch2.md")
	require.NoError(t, os.WriteFile(f1, []byte("# Chapter 1\n\nContent one.\n"), 0o644))
	require.NoError(t, os.WriteFile(f2, []byte("# Chapter 2\n\nContent two.\n"), 0o644))

	out := filepath.Join(tmp, "merged.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", f1, f2, "-o", out})

	err := rootCmd.Execute()
	require.NoError(t, err, "merge should succeed")

	r, err := zip.OpenReader(out)
	require.NoError(t, err, "merged output should be valid docx")
	defer r.Close()
}

func TestMergeWithStyle(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()

	f1 := filepath.Join(tmp, "a.md")
	require.NoError(t, os.WriteFile(f1, []byte("# A\n\nText.\n"), 0o644))
	out := filepath.Join(tmp, "merged_styled.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", f1, "-o", out, "--style", "simple"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err)
}

func TestMergeWithFlags(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()

	f1 := filepath.Join(tmp, "x.md")
	require.NoError(t, os.WriteFile(f1, []byte("# X\n\nHello.\n"), 0o644))
	out := filepath.Join(tmp, "merged_flags.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", f1, "-o", out, "--no-cover", "--no-toc", "--no-numbering"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err)
}

func TestMergeNonexistentFile(t *testing.T) {
	resetFlags()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", "/tmp/no_such_file_zzz.md", "-o", "/tmp/nope.docx"})

	err := rootCmd.Execute()
	assert.Error(t, err, "merge with nonexistent file should error")
}

func TestMergeNonexistentStyle(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "m.md")
	require.NoError(t, os.WriteFile(f1, []byte("# M\n"), 0o644))
	out := filepath.Join(tmp, "out.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", f1, "-o", out, "--style", "nonexistent"})

	err := rootCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

func TestMergeContentsYAML(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()

	// Create md files
	f1 := filepath.Join(tmp, "part1.md")
	f2 := filepath.Join(tmp, "part2.md")
	require.NoError(t, os.WriteFile(f1, []byte("# Part 1\n\nContent.\n"), 0o644))
	require.NoError(t, os.WriteFile(f2, []byte("# Part 2\n\nMore.\n"), 0o644))

	// Create contents.yaml
	contentsYAML := "files:\n  - part1.md\n  - part2.md\n"
	yamlPath := filepath.Join(tmp, "contents.yaml")
	require.NoError(t, os.WriteFile(yamlPath, []byte(contentsYAML), 0o644))

	out := filepath.Join(tmp, "from_yaml.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", yamlPath, "-o", out})

	err := rootCmd.Execute()
	require.NoError(t, err)
	_, err = os.Stat(out)
	assert.NoError(t, err)
}

func TestMergeDefaultOutput(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "doc.md")
	require.NoError(t, os.WriteFile(f1, []byte("# Doc\n\nText.\n"), 0o644))

	// Set HOME to temp so default output path goes somewhere predictable
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", f1})

	err := rootCmd.Execute()
	require.NoError(t, err)

	// Check merged.docx exists somewhere
	desktopPath := filepath.Join(tmp, "Desktop", "merged.docx")
	fallbackPath := filepath.Join("merged.docx")
	_, err1 := os.Stat(desktopPath)
	_, err2 := os.Stat(fallbackPath)
	assert.True(t, err1 == nil || err2 == nil, "merged.docx should exist")
}

func TestConvertRichContent(t *testing.T) {
	resetFlags()
	tmp := t.TempDir()
	content := `# Main Title

## Section One

A paragraph with **bold** and *italic* text.

- bullet one
- bullet two

1. ordered one
2. ordered two

### Code Example

` + "```go\nfmt.Println(\"hello\")\n```" + `

---

End of document.
`
	md := writeTempMD(t, tmp, content)
	out := filepath.Join(tmp, "rich.docx")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{md, "-o", out})

	err := rootCmd.Execute()
	require.NoError(t, err)

	r, err := zip.OpenReader(out)
	require.NoError(t, err)
	defer r.Close()

	var hasStyles bool
	for _, f := range r.File {
		if f.Name == "word/styles.xml" {
			hasStyles = true
			break
		}
	}
	assert.True(t, hasStyles, "docx should contain word/styles.xml")
}
