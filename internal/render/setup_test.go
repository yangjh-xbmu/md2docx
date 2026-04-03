package render

import (
	"embed"
	"io/fs"

	"github.com/yangjh-xbmu/md2docx/internal/style"
)

//go:embed testdata
var testStylesRaw embed.FS

func init() {
	sub, err := fs.Sub(testStylesRaw, "testdata")
	if err != nil {
		panic(err)
	}
	style.EmbeddedStyles = sub
}
