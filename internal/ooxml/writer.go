package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// Package holds all parts of a docx file.
type Package struct {
	Document     *Document
	Styles       []byte // pre-rendered styles.xml
	Numbering    []byte // pre-rendered numbering.xml (optional)
	Header       []byte // pre-rendered header1.xml (optional)
	Footer       []byte // pre-rendered footer1.xml (optional)
	Images       []ImagePart
	Rels         []Relationship
}

type ImagePart struct {
	PartName    string // e.g. "word/media/image1.png"
	ContentType string // e.g. "image/png"
	Data        []byte
}

type Relationship struct {
	ID     string
	Type   string
	Target string
}

// WriteDocx writes the package to a .docx file.
func WriteDocx(pkg *Package, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// [Content_Types].xml
	if err := writeContentTypes(w, pkg); err != nil {
		return err
	}

	// _rels/.rels
	if err := writeRootRels(w); err != nil {
		return err
	}

	// word/document.xml
	docData, err := xml.MarshalIndent(pkg.Document, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 document.xml 失败: %w", err)
	}
	if err := writeZipFile(w, "word/document.xml", append([]byte(xml.Header), docData...)); err != nil {
		return err
	}

	// word/styles.xml
	if pkg.Styles != nil {
		if err := writeZipFile(w, "word/styles.xml", pkg.Styles); err != nil {
			return err
		}
	}

	// word/numbering.xml
	if pkg.Numbering != nil {
		if err := writeZipFile(w, "word/numbering.xml", pkg.Numbering); err != nil {
			return err
		}
	}

	// word/header1.xml
	if pkg.Header != nil {
		if err := writeZipFile(w, "word/header1.xml", pkg.Header); err != nil {
			return err
		}
	}

	// word/footer1.xml
	if pkg.Footer != nil {
		if err := writeZipFile(w, "word/footer1.xml", pkg.Footer); err != nil {
			return err
		}
	}

	// word/media/*
	for _, img := range pkg.Images {
		if err := writeZipFile(w, img.PartName, img.Data); err != nil {
			return err
		}
	}

	// word/_rels/document.xml.rels
	if err := writeDocRels(w, pkg); err != nil {
		return err
	}

	return nil
}

func writeZipFile(w *zip.Writer, name string, data []byte) error {
	fw, err := w.Create(name)
	if err != nil {
		return fmt.Errorf("创建 zip 条目 %s 失败: %w", name, err)
	}
	_, err = fw.Write(data)
	return err
}

func writeContentTypes(w *zip.Writer, pkg *Package) error {
	ct := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Default Extension="png" ContentType="image/png"/>
  <Default Extension="jpeg" ContentType="image/jpeg"/>
  <Default Extension="jpg" ContentType="image/jpeg"/>
  <Default Extension="gif" ContentType="image/gif"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>`

	if pkg.Numbering != nil {
		ct += `
  <Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>`
	}
	if pkg.Header != nil {
		ct += `
  <Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>`
	}
	if pkg.Footer != nil {
		ct += `
  <Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>`
	}
	ct += `
</Types>`

	return writeZipFile(w, "[Content_Types].xml", []byte(ct))
}

func writeRootRels(w *zip.Writer) error {
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	return writeZipFile(w, "_rels/.rels", []byte(rels))
}

func writeDocRels(w *zip.Writer, pkg *Package) error {
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`

	nextID := 2

	if pkg.Numbering != nil {
		rels += fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/>`, nextID)
		nextID++
	}

	if pkg.Header != nil {
		rels += fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>`, nextID)
		nextID++
	}

	if pkg.Footer != nil {
		rels += fmt.Sprintf(`
  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>`, nextID)
		nextID++
	}

	for _, rel := range pkg.Rels {
		rels += fmt.Sprintf(`
  <Relationship Id="%s" Type="%s" Target="%s"/>`, rel.ID, rel.Type, rel.Target)
	}

	rels += `
</Relationships>`

	return writeZipFile(w, "word/_rels/document.xml.rels", []byte(rels))
}
