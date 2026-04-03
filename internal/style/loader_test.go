package style

import (
	"embed"
	"io/fs"
	"strings"
	"testing"
)

//go:embed all:testdata
var testStylesRaw embed.FS

func init() {
	sub, err := fs.Sub(testStylesRaw, "testdata")
	if err != nil {
		panic(err)
	}
	EmbeddedStyles = sub
}

func TestLoadAcademicCN(t *testing.T) {
	s, err := Load("academic-cn")
	if err != nil {
		t.Fatalf("Load academic-cn: %v", err)
	}

	// Meta
	if s.Meta.ID != "academic-cn" {
		t.Errorf("Meta.ID = %q, want %q", s.Meta.ID, "academic-cn")
	}
	if s.Meta.Locale != "zh-CN" {
		t.Errorf("Meta.Locale = %q, want %q", s.Meta.Locale, "zh-CN")
	}

	// Heading1: 黑体, centered, page break before
	if s.Styles.Heading1.FontCJK != "黑体" {
		t.Errorf("Heading1.FontCJK = %q, want %q", s.Styles.Heading1.FontCJK, "黑体")
	}
	if s.Styles.Heading1.FontLatin != "Arial" {
		t.Errorf("Heading1.FontLatin = %q, want %q", s.Styles.Heading1.FontLatin, "Arial")
	}
	if s.Styles.Heading1.FontSize != "22pt" {
		t.Errorf("Heading1.FontSize = %q, want %q", s.Styles.Heading1.FontSize, "22pt")
	}
	if !s.Styles.Heading1.Bold {
		t.Error("Heading1.Bold should be true")
	}
	if s.Styles.Heading1.Alignment != "center" {
		t.Errorf("Heading1.Alignment = %q, want %q", s.Styles.Heading1.Alignment, "center")
	}
	if !s.Styles.Heading1.PageBreakBefore {
		t.Error("Heading1.PageBreakBefore should be true for academic-cn")
	}

	// Body: 宋体, 12pt, justified, 1.5 spacing, 2em indent
	if s.Styles.Body.FontCJK != "宋体" {
		t.Errorf("Body.FontCJK = %q, want %q", s.Styles.Body.FontCJK, "宋体")
	}
	if s.Styles.Body.Alignment != "justify" {
		t.Errorf("Body.Alignment = %q, want %q", s.Styles.Body.Alignment, "justify")
	}
	if s.Styles.Body.LineSpacing != 1.5 {
		t.Errorf("Body.LineSpacing = %v, want %v", s.Styles.Body.LineSpacing, 1.5)
	}
	if s.Styles.Body.FirstLineIndent != "2em" {
		t.Errorf("Body.FirstLineIndent = %q, want %q", s.Styles.Body.FirstLineIndent, "2em")
	}

	// Features enabled
	if !s.HeadingNumbering.Enabled {
		t.Error("HeadingNumbering should be enabled for academic-cn")
	}
	if !s.TOC.Enabled {
		t.Error("TOC should be enabled for academic-cn")
	}
	if !s.Cover.Enabled {
		t.Error("Cover should be enabled for academic-cn")
	}

	// Heading numbering formats
	if f, ok := s.HeadingNumbering.Formats[1]; !ok || !strings.Contains(f, "{1}") {
		t.Errorf("HeadingNumbering.Formats[1] = %q, want format containing {1}", f)
	}
}

func TestLoadSimple(t *testing.T) {
	s, err := Load("simple")
	if err != nil {
		t.Fatalf("Load simple: %v", err)
	}

	if s.Meta.ID != "simple" {
		t.Errorf("Meta.ID = %q, want %q", s.Meta.ID, "simple")
	}

	// Features disabled
	if s.HeadingNumbering.Enabled {
		t.Error("HeadingNumbering should be disabled for simple")
	}
	if s.TOC.Enabled {
		t.Error("TOC should be disabled for simple")
	}
	if s.Cover.Enabled {
		t.Error("Cover should be disabled for simple")
	}

	// Should still have sensible defaults from applyDefaults
	if s.Styles.Body.FontSize != "12pt" {
		t.Errorf("Body.FontSize = %q, want %q (from defaults)", s.Styles.Body.FontSize, "12pt")
	}

	// Heading1 should NOT have page break before
	if s.Styles.Heading1.PageBreakBefore {
		t.Error("Heading1.PageBreakBefore should be false for simple")
	}

	// Simple uses smaller heading font
	if s.Styles.Heading1.FontSize != "18pt" {
		t.Errorf("Heading1.FontSize = %q, want %q", s.Styles.Heading1.FontSize, "18pt")
	}
}

func TestListAvailableIncludesAll(t *testing.T) {
	names := ListAvailable()
	want := map[string]bool{"default": false, "academic-cn": false, "simple": false}
	for _, n := range names {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("ListAvailable() missing %q", name)
		}
	}
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("nonexistent-style-xyz")
	if err == nil {
		t.Error("Load nonexistent style should return error")
	}
}
