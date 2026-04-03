package ooxml

import "encoding/xml"

// Document represents word/document.xml
type Document struct {
	XMLName xml.Name `xml:"w:document"`
	W       string   `xml:"xmlns:w,attr"`
	R       string   `xml:"xmlns:r,attr"`
	WP      string   `xml:"xmlns:wp,attr"`
	Body    Body     `xml:"w:body"`
}

type Body struct {
	Elements []any `xml:",any"`
	SectPr   *SectionProperties `xml:"w:sectPr,omitempty"`
}

type SectionProperties struct {
	PgSz     *PageSize    `xml:"w:pgSz,omitempty"`
	PgMar    *PageMargins `xml:"w:pgMar,omitempty"`
	HeaderRef *HeaderFooterRef `xml:"w:headerReference,omitempty"`
	FooterRef *HeaderFooterRef `xml:"w:footerReference,omitempty"`
}

type PageSize struct {
	W    string `xml:"w:w,attr"`
	H    string `xml:"w:h,attr"`
	Orient string `xml:"w:orient,attr,omitempty"`
}

type PageMargins struct {
	Top    string `xml:"w:top,attr"`
	Bottom string `xml:"w:bottom,attr"`
	Left   string `xml:"w:left,attr"`
	Right  string `xml:"w:right,attr"`
	Header string `xml:"w:header,attr"`
	Footer string `xml:"w:footer,attr"`
}

type HeaderFooterRef struct {
	Type string `xml:"w:type,attr"`
	RID  string `xml:"r:id,attr"`
}

// Paragraph represents w:p
type Paragraph struct {
	XMLName xml.Name              `xml:"w:p"`
	PPr     *ParagraphProperties  `xml:"w:pPr,omitempty"`
	Runs    []any                 `xml:",any"`
}

type ParagraphProperties struct {
	PStyle       *PStyle        `xml:"w:pStyle,omitempty"`
	Jc           *Justification `xml:"w:jc,omitempty"`
	Spacing      *Spacing       `xml:"w:spacing,omitempty"`
	Ind          *Indentation   `xml:"w:ind,omitempty"`
	KeepNext     *Empty         `xml:"w:keepNext,omitempty"`
	PageBreakBefore *Empty      `xml:"w:pageBreakBefore,omitempty"`
	NumPr        *NumberingPr   `xml:"w:numPr,omitempty"`
}

type PStyle struct {
	Val string `xml:"w:val,attr"`
}

type Justification struct {
	Val string `xml:"w:val,attr"`
}

type Spacing struct {
	Before  string `xml:"w:before,attr,omitempty"`
	After   string `xml:"w:after,attr,omitempty"`
	Line    string `xml:"w:line,attr,omitempty"`
	LineRule string `xml:"w:lineRule,attr,omitempty"`
}

type Indentation struct {
	FirstLine   string `xml:"w:firstLine,attr,omitempty"`
	Left        string `xml:"w:left,attr,omitempty"`
	Right       string `xml:"w:right,attr,omitempty"`
	Hanging     string `xml:"w:hanging,attr,omitempty"`
}

type NumberingPr struct {
	Ilvl  *IntVal `xml:"w:ilvl,omitempty"`
	NumId *IntVal `xml:"w:numId,omitempty"`
}

type IntVal struct {
	Val string `xml:"w:val,attr"`
}

type Empty struct{}

// Run represents w:r
type Run struct {
	XMLName xml.Name        `xml:"w:r"`
	RPr     *RunProperties  `xml:"w:rPr,omitempty"`
	Content []any           `xml:",any"`
}

type RunProperties struct {
	RFonts    *RunFonts `xml:"w:rFonts,omitempty"`
	Bold      *Empty    `xml:"w:b,omitempty"`
	Italic    *Empty    `xml:"w:i,omitempty"`
	Underline *UVal     `xml:"w:u,omitempty"`
	Sz        *SzVal    `xml:"w:sz,omitempty"`
	SzCs      *SzVal    `xml:"w:szCs,omitempty"`
	Color     *SVal     `xml:"w:color,omitempty"`
	RStyle    *SVal     `xml:"w:rStyle,omitempty"`
}

type RunFonts struct {
	ASCII    string `xml:"w:ascii,attr,omitempty"`
	HAnsi    string `xml:"w:hAnsi,attr,omitempty"`
	EastAsia string `xml:"w:eastAsia,attr,omitempty"`
	CS       string `xml:"w:cs,attr,omitempty"`
}

type SzVal struct {
	Val string `xml:"w:val,attr"`
}

type SVal struct {
	Val string `xml:"w:val,attr"`
}

type UVal struct {
	Val string `xml:"w:val,attr"`
}

type Text struct {
	XMLName xml.Name `xml:"w:t"`
	Space   string   `xml:"xml:space,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Break struct {
	XMLName xml.Name `xml:"w:br"`
	Type    string   `xml:"w:type,attr,omitempty"`
}

type Hyperlink struct {
	XMLName xml.Name `xml:"w:hyperlink"`
	RID     string   `xml:"r:id,attr,omitempty"`
	Runs    []Run    `xml:",any"`
}

// FldChar for field codes (TOC, PAGE, etc.)
type FldChar struct {
	XMLName     xml.Name `xml:"w:fldChar"`
	FldCharType string   `xml:"w:fldCharType,attr"`
}

type InstrText struct {
	XMLName xml.Name `xml:"w:instrText"`
	Space   string   `xml:"xml:space,attr,omitempty"`
	Value   string   `xml:",chardata"`
}
