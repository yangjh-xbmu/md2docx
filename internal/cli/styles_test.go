package cli

import (
	"bytes"
	"embed"
	"io/fs"
	"testing"

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

func TestStylesListCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"styles", "list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("styles list: %v", err)
	}

	output := buf.String()
	for _, name := range []string{"default", "academic-cn", "simple"} {
		if !bytes.Contains([]byte(output), []byte(name)) {
			t.Errorf("styles list output missing %q", name)
		}
	}
}

func TestStylesShowCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"styles", "show", "academic-cn"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("styles show academic-cn: %v", err)
	}

	output := buf.String()
	// Should contain style name and key properties
	if !bytes.Contains([]byte(output), []byte("academic-cn")) {
		t.Error("styles show output missing style ID")
	}
	if !bytes.Contains([]byte(output), []byte("黑体")) {
		t.Error("styles show output missing 黑体 font")
	}
}

func TestStylesShowNonexistent(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"styles", "show", "nonexistent"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("styles show nonexistent should return error")
	}
}
