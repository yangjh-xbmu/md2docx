package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// ParseMarkdown parses markdown body into goldmark AST.
// Returns the AST root node and the source bytes.
func ParseMarkdown(body string) (ast.Node, []byte) {
	src := []byte(body)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // tables, strikethrough, autolinks, task lists
		),
		goldmark.WithParserOptions(
			gparser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader)
	return doc, src
}
