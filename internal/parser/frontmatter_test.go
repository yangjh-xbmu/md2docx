package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMeta map[string]any
		wantBody string
	}{
		{
			name:     "no frontmatter",
			input:    "# Hello\n\nWorld",
			wantMeta: map[string]any{},
			wantBody: "# Hello\n\nWorld",
		},
		{
			name:  "with frontmatter",
			input: "---\ntitle: Test\nauthor: Me\n---\n# Hello",
			wantMeta: map[string]any{
				"title":  "Test",
				"author": "Me",
			},
			wantBody: "# Hello",
		},
		{
			name:     "unclosed frontmatter",
			input:    "---\ntitle: Test\n# Hello",
			wantMeta: map[string]any{},
			wantBody: "---\ntitle: Test\n# Hello",
		},
		{
			name:     "empty content",
			input:    "",
			wantMeta: map[string]any{},
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, body := ParseFrontmatter(tt.input)
			assert.Equal(t, tt.wantMeta, meta)
			assert.Equal(t, tt.wantBody, body)
		})
	}
}
