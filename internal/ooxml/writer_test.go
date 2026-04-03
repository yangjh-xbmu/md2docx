package ooxml

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// minimalPackage returns a Package with one paragraph, suitable for WriteDocx tests.
func minimalPackage() *Package {
	return &Package{
		Document: &Document{
			W:  "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
			R:  "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
			WP: "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing",
			Body: Body{
				Elements: []any{
					Paragraph{
						Runs: []any{
							Run{
								Content: []any{
									Text{Value: "Hello World", Space: "preserve"},
								},
							},
						},
					},
				},
			},
		},
		Styles: []byte(`<?xml version="1.0" encoding="UTF-8"?><w:styles/>`),
	}
}

func TestWriteDocx_MinimalPackage(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "test.docx")

	pkg := minimalPackage()
	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	// Verify file exists
	info, err := os.Stat(outPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))

	// Open as zip and verify entries
	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	entries := make(map[string]bool)
	for _, f := range r.File {
		entries[f.Name] = true
	}

	expected := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"word/document.xml",
		"word/styles.xml",
		"word/_rels/document.xml.rels",
	}
	for _, name := range expected {
		assert.True(t, entries[name], "missing zip entry: %s", name)
	}
}

func TestWriteDocx_NoNumbering(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "no-num.docx")

	pkg := minimalPackage()
	pkg.Numbering = nil // explicitly no numbering

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	for _, f := range r.File {
		assert.NotEqual(t, "word/numbering.xml", f.Name, "numbering.xml should not exist when Numbering is nil")
	}
}

func TestWriteDocx_WithNumbering(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "with-num.docx")

	pkg := minimalPackage()
	pkg.Numbering = []byte(`<?xml version="1.0"?><w:numbering/>`)

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	entries := make(map[string]bool)
	for _, f := range r.File {
		entries[f.Name] = true
	}
	assert.True(t, entries["word/numbering.xml"], "numbering.xml should exist")
}

func TestWriteDocx_WithHeaderAndFooter(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "hf.docx")

	pkg := minimalPackage()
	pkg.Header = []byte(`<?xml version="1.0"?><w:hdr/>`)
	pkg.Footer = []byte(`<?xml version="1.0"?><w:ftr/>`)

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	entries := make(map[string]bool)
	for _, f := range r.File {
		entries[f.Name] = true
	}
	assert.True(t, entries["word/header1.xml"], "header1.xml should exist")
	assert.True(t, entries["word/footer1.xml"], "footer1.xml should exist")
}

func TestWriteDocx_WithImages(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "images.docx")

	pkg := minimalPackage()
	pkg.Images = []ImagePart{
		{
			PartName:    "word/media/image1.png",
			ContentType: "image/png",
			Data:        []byte("fake-png-data"),
		},
		{
			PartName:    "word/media/image2.jpeg",
			ContentType: "image/jpeg",
			Data:        []byte("fake-jpeg-data"),
		},
	}

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	entries := make(map[string]bool)
	for _, f := range r.File {
		entries[f.Name] = true
	}
	assert.True(t, entries["word/media/image1.png"], "image1.png should exist")
	assert.True(t, entries["word/media/image2.jpeg"], "image2.jpeg should exist")
}

func TestWriteDocx_OutputDirAutoCreate(t *testing.T) {
	dir := t.TempDir()
	nestedPath := filepath.Join(dir, "deep", "nested", "dir", "output.docx")

	pkg := minimalPackage()
	err := WriteDocx(pkg, nestedPath)
	require.NoError(t, err)

	_, err = os.Stat(nestedPath)
	require.NoError(t, err)
}

func TestWriteDocx_ContentTypesIncludesNumbering(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "ct.docx")

	pkg := minimalPackage()
	pkg.Numbering = []byte(`<?xml version="1.0"?><w:numbering/>`)
	pkg.Header = []byte(`<?xml version="1.0"?><w:hdr/>`)
	pkg.Footer = []byte(`<?xml version="1.0"?><w:ftr/>`)

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	// Read [Content_Types].xml
	var ctContent string
	for _, f := range r.File {
		if f.Name == "[Content_Types].xml" {
			rc, err := f.Open()
			require.NoError(t, err)
			buf := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(buf)
			rc.Close()
			ctContent = string(buf)
			break
		}
	}

	require.NotEmpty(t, ctContent, "[Content_Types].xml should be non-empty")
	assert.Contains(t, ctContent, "numbering.xml")
	assert.Contains(t, ctContent, "header1.xml")
	assert.Contains(t, ctContent, "footer1.xml")
}

func TestWriteDocx_DocumentRelsContainsStyles(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "rels.docx")

	pkg := minimalPackage()
	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	var relsContent string
	for _, f := range r.File {
		if f.Name == "word/_rels/document.xml.rels" {
			rc, err := f.Open()
			require.NoError(t, err)
			buf := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(buf)
			rc.Close()
			relsContent = string(buf)
			break
		}
	}

	require.NotEmpty(t, relsContent)
	assert.Contains(t, relsContent, "styles.xml")
	assert.Contains(t, relsContent, "relationships/styles")
}

func TestWriteDocx_DocumentRelsWithOptionalParts(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "rels-full.docx")

	pkg := minimalPackage()
	pkg.Numbering = []byte(`<?xml version="1.0"?><w:numbering/>`)
	pkg.Header = []byte(`<?xml version="1.0"?><w:hdr/>`)
	pkg.Footer = []byte(`<?xml version="1.0"?><w:ftr/>`)
	pkg.Rels = []Relationship{
		{ID: "rId100", Type: "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image", Target: "media/image1.png"},
	}

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	var relsContent string
	for _, f := range r.File {
		if f.Name == "word/_rels/document.xml.rels" {
			rc, err := f.Open()
			require.NoError(t, err)
			buf := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(buf)
			rc.Close()
			relsContent = string(buf)
			break
		}
	}

	require.NotEmpty(t, relsContent)
	assert.Contains(t, relsContent, "numbering.xml")
	assert.Contains(t, relsContent, "header1.xml")
	assert.Contains(t, relsContent, "footer1.xml")
	assert.Contains(t, relsContent, "rId100")
	assert.Contains(t, relsContent, "media/image1.png")
}

func TestWriteDocx_RootRels(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "root-rels.docx")

	pkg := minimalPackage()
	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	var relsContent string
	for _, f := range r.File {
		if f.Name == "_rels/.rels" {
			rc, err := f.Open()
			require.NoError(t, err)
			buf := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(buf)
			rc.Close()
			relsContent = string(buf)
			break
		}
	}

	require.NotEmpty(t, relsContent)
	assert.Contains(t, relsContent, "officeDocument")
	assert.Contains(t, relsContent, "word/document.xml")
}

func TestWriteDocx_DocumentXMLContainsContent(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "content.docx")

	pkg := minimalPackage()
	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	var docContent string
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			require.NoError(t, err)
			buf := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(buf)
			rc.Close()
			docContent = string(buf)
			break
		}
	}

	require.NotEmpty(t, docContent)
	assert.Contains(t, docContent, "Hello World")
}

func TestWriteDocx_NilStylesOmitsEntry(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "no-styles.docx")

	pkg := minimalPackage()
	pkg.Styles = nil

	err := WriteDocx(pkg, outPath)
	require.NoError(t, err)

	r, err := zip.OpenReader(outPath)
	require.NoError(t, err)
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/styles.xml" {
			t.Fatal("word/styles.xml should not be present when Styles is nil")
		}
	}
}
