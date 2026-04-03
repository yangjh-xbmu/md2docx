package merge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffsetHeadings(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		offset int
		want   string
	}{
		{
			name:   "no offset",
			input:  "# Title\ntext\n## Sub",
			offset: 0,
			want:   "# Title\ntext\n## Sub",
		},
		{
			name:   "offset by 1",
			input:  "# Title\ntext\n## Sub",
			offset: 1,
			want:   "## Title\ntext\n### Sub",
		},
		{
			name:   "cap at level 6",
			input:  "##### H5\n###### H6",
			offset: 3,
			want:   "###### H5\n###### H6",
		},
		{
			name:   "non-heading hash preserved",
			input:  "#hashtag\n# Real heading",
			offset: 1,
			want:   "#hashtag\n## Real heading",
		},
		{
			name:   "empty string",
			input:  "",
			offset: 1,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OffsetHeadings(tt.input, tt.offset)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFM   string
		wantBody string
	}{
		{
			name:     "no frontmatter",
			input:    "# Title\nBody",
			wantFM:   "",
			wantBody: "# Title\nBody",
		},
		{
			name:     "with frontmatter",
			input:    "---\ntitle: Test\n---\n# Title",
			wantFM:   "---\ntitle: Test\n---",
			wantBody: "# Title",
		},
		{
			name:     "unclosed frontmatter",
			input:    "---\ntitle: Test\n# Title",
			wantFM:   "",
			wantBody: "---\ntitle: Test\n# Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body := splitFrontmatter(tt.input)
			assert.Equal(t, tt.wantFM, fm)
			assert.Equal(t, tt.wantBody, body)
		})
	}
}

func TestMergeFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	ch1 := filepath.Join(dir, "ch1.md")
	ch2 := filepath.Join(dir, "ch2.md")
	ch3 := filepath.Join(dir, "ch3.md")

	require.NoError(t, os.WriteFile(ch1, []byte("---\ntitle: Book\n---\n# Chapter 1\n\nContent 1"), 0644))
	require.NoError(t, os.WriteFile(ch2, []byte("---\ntitle: Ignored\n---\n# Chapter 2\n\nContent 2"), 0644))
	require.NoError(t, os.WriteFile(ch3, []byte("# Chapter 3\n\nContent 3"), 0644))

	entries := []Entry{
		{Path: ch1, HeadingOffset: 0},
		{Path: ch2, HeadingOffset: 0},
		{Path: ch3, HeadingOffset: 0},
	}

	result, err := MergeFiles(entries)
	require.NoError(t, err)

	// First file's frontmatter preserved
	assert.Contains(t, result, "---\ntitle: Book\n---")
	// Second file's frontmatter stripped
	assert.NotContains(t, result, "title: Ignored")
	// All content present in order
	assert.Contains(t, result, "# Chapter 1")
	assert.Contains(t, result, "# Chapter 2")
	assert.Contains(t, result, "# Chapter 3")
	assert.Contains(t, result, "Content 1")
	assert.Contains(t, result, "Content 2")
	assert.Contains(t, result, "Content 3")

	// Verify order
	idx1 := indexOf(result, "Chapter 1")
	idx2 := indexOf(result, "Chapter 2")
	idx3 := indexOf(result, "Chapter 3")
	assert.Less(t, idx1, idx2)
	assert.Less(t, idx2, idx3)
}

func TestMergeFilesWithOffset(t *testing.T) {
	dir := t.TempDir()

	ch1 := filepath.Join(dir, "ch1.md")
	ch2 := filepath.Join(dir, "ch2.md")

	require.NoError(t, os.WriteFile(ch1, []byte("# Main Title\n\nIntro"), 0644))
	require.NoError(t, os.WriteFile(ch2, []byte("# Sub Section\n\n## Detail"), 0644))

	entries := []Entry{
		{Path: ch1, HeadingOffset: 0},
		{Path: ch2, HeadingOffset: 1},
	}

	result, err := MergeFiles(entries)
	require.NoError(t, err)

	assert.Contains(t, result, "# Main Title")
	assert.Contains(t, result, "## Sub Section")
	assert.Contains(t, result, "### Detail")
}

func TestMergeFilesEmpty(t *testing.T) {
	_, err := MergeFiles(nil)
	assert.Error(t, err)
}

func TestMergeFilesMissing(t *testing.T) {
	entries := []Entry{{Path: "/nonexistent/file.md"}}
	_, err := MergeFiles(entries)
	assert.Error(t, err)
}

func TestParseContentsYAML(t *testing.T) {
	dir := t.TempDir()

	// Create referenced md files
	require.NoError(t, os.WriteFile(filepath.Join(dir, "intro.md"), []byte("# Intro"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ch1.md"), []byte("# Ch1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ch2.md"), []byte("# Ch2"), 0644))

	yamlContent := `heading_offset: 0
files:
  - intro.md
  - path: ch1.md
    heading_offset: 1
  - path: ch2.md
`
	yamlPath := filepath.Join(dir, "contents.yaml")
	require.NoError(t, os.WriteFile(yamlPath, []byte(yamlContent), 0644))

	entries, err := ParseContentsYAML(yamlPath)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	assert.Equal(t, filepath.Join(dir, "intro.md"), entries[0].Path)
	assert.Equal(t, 0, entries[0].HeadingOffset)

	assert.Equal(t, filepath.Join(dir, "ch1.md"), entries[1].Path)
	assert.Equal(t, 1, entries[1].HeadingOffset)

	assert.Equal(t, filepath.Join(dir, "ch2.md"), entries[2].Path)
	assert.Equal(t, 0, entries[2].HeadingOffset)
}

func TestResolveGlobs(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("A"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.md"), []byte("B"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "c.txt"), []byte("C"), 0644))

	pattern := filepath.Join(dir, "*.md")
	result, err := ResolveGlobs([]string{pattern})
	require.NoError(t, err)

	assert.Len(t, result, 2)
	assert.Contains(t, result[0], "a.md")
	assert.Contains(t, result[1], "b.md")
}

func TestResolveGlobsDedup(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("A"), 0644))

	pattern := filepath.Join(dir, "*.md")
	result, err := ResolveGlobs([]string{pattern, pattern})
	require.NoError(t, err)

	assert.Len(t, result, 1)
}

func TestIsContentsYAML(t *testing.T) {
	assert.True(t, IsContentsYAML("contents.yaml"))
	assert.True(t, IsContentsYAML("contents.yml"))
	assert.True(t, IsContentsYAML("path/to/contents.YAML"))
	assert.False(t, IsContentsYAML("file.md"))
	assert.False(t, IsContentsYAML("file.txt"))
}

func indexOf(s, substr string) int {
	for i := range s {
		if len(s[i:]) >= len(substr) && s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
