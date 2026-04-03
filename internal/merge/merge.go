package merge

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Entry represents a single file in a merge operation.
type Entry struct {
	Path          string
	HeadingOffset int
}

// ContentsFile represents a contents.yaml structure.
type ContentsFile struct {
	HeadingOffset int   `yaml:"heading_offset"`
	Files         []any `yaml:"files"`
}

// ParseContentsYAML reads a contents.yaml and returns merge entries.
// Paths are resolved relative to the directory containing the YAML file.
func ParseContentsYAML(yamlPath string) ([]Entry, error) {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 %s: %w", yamlPath, err)
	}

	var cf ContentsFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", yamlPath, err)
	}

	baseDir := filepath.Dir(yamlPath)
	var entries []Entry

	for _, item := range cf.Files {
		switch v := item.(type) {
		case string:
			entries = append(entries, Entry{
				Path:          filepath.Join(baseDir, v),
				HeadingOffset: cf.HeadingOffset,
			})
		case map[string]any:
			p, _ := v["path"].(string)
			if p == "" {
				continue
			}
			offset := cf.HeadingOffset
			if ho, ok := v["heading_offset"]; ok {
				switch n := ho.(type) {
				case int:
					offset = n
				case float64:
					offset = int(n)
				}
			}
			entries = append(entries, Entry{
				Path:          filepath.Join(baseDir, p),
				HeadingOffset: offset,
			})
		}
	}

	return entries, nil
}

// ResolveGlobs expands glob patterns and returns sorted, deduplicated .md file paths.
func ResolveGlobs(patterns []string) ([]string, error) {
	seen := make(map[string]bool)
	var result []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("glob 模式无效 '%s': %w", pattern, err)
		}
		for _, m := range matches {
			if !strings.HasSuffix(strings.ToLower(m), ".md") {
				continue
			}
			abs, err := filepath.Abs(m)
			if err != nil {
				continue
			}
			if !seen[abs] {
				seen[abs] = true
				result = append(result, m)
			}
		}
	}

	sort.Strings(result)
	return result, nil
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s`)

// OffsetHeadings increases heading levels by offset.
// Headings are capped at level 6.
func OffsetHeadings(content string, offset int) string {
	if offset <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		m := headingRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		level := min(len(m[1])+offset, 6)
		lines[i] = strings.Repeat("#", level) + line[len(m[1]):]
	}
	return strings.Join(lines, "\n")
}

// MergeFiles reads all entries, applies heading offsets, and returns combined markdown.
// The first file's frontmatter is preserved; subsequent files' frontmatter is stripped.
func MergeFiles(entries []Entry) (string, error) {
	if len(entries) == 0 {
		return "", fmt.Errorf("没有要合并的文件")
	}

	var parts []string
	for i, e := range entries {
		data, err := os.ReadFile(e.Path)
		if err != nil {
			return "", fmt.Errorf("无法读取 %s: %w", e.Path, err)
		}
		content := string(data)

		if i == 0 {
			// Keep first file as-is (including frontmatter), then offset body
			fm, body := splitFrontmatter(content)
			body = OffsetHeadings(body, e.HeadingOffset)
			if fm != "" {
				parts = append(parts, fm+"\n"+body)
			} else {
				parts = append(parts, body)
			}
		} else {
			// Strip frontmatter from subsequent files
			_, body := splitFrontmatter(content)
			body = OffsetHeadings(body, e.HeadingOffset)
			parts = append(parts, body)
		}
	}

	return strings.Join(parts, "\n\n") + "\n", nil
}

// splitFrontmatter separates YAML frontmatter from body.
// Returns (frontmatter_block_including_delimiters, body).
func splitFrontmatter(content string) (string, string) {
	if !strings.HasPrefix(content, "---") {
		return "", content
	}

	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", content
	}

	fmEnd := 3 + idx + 4 // "---" + rest[:idx] + "\n---"
	fm := content[:fmEnd]
	body := content[fmEnd:]
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}
	return fm, body
}

// IsContentsYAML checks if a path looks like a contents.yaml file.
func IsContentsYAML(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}
