package parser

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFrontmatter extracts YAML frontmatter from markdown content.
// Returns metadata map and the body without frontmatter.
func ParseFrontmatter(content string) (map[string]any, string) {
	if !strings.HasPrefix(content, "---") {
		return make(map[string]any), content
	}

	// Find closing ---
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return make(map[string]any), content
	}

	yamlStr := rest[:idx]
	body := rest[idx+4:] // skip \n---
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	meta := make(map[string]any)
	if err := yaml.Unmarshal([]byte(yamlStr), &meta); err != nil {
		return make(map[string]any), content
	}

	return meta, body
}
