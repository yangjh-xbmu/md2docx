package main

import (
	"embed"
	"os"

	"github.com/yangjh-xbmu/md2docx/internal/cli"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

//go:embed all:styles
var embeddedStyles embed.FS

func main() {
	style.EmbeddedStyles = embeddedStyles
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
