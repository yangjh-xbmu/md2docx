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
	HeaderRef *HeaderFooterRef `xml:"w:headerReference,omitempty"`
	FooterRef *HeaderFooterRef `xml:"w:footerReference,omitempty"`
	PgSz      *PageSize        `xml:"w:pgSz,omitempty"`
	PgMar     *PageMargins     `xml:"w:pgMar,omitempty"`
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

// Table represents w:tbl
type Table struct {
	XMLName xml.Name         `xml:"w:tbl"`
	TblPr   *TableProperties `xml:"w:tblPr,omitempty"`
	TblGrid *TableGrid       `xml:"w:tblGrid,omitempty"`
	Rows    []TableRow       `xml:",any"`
}

type TableProperties struct {
	TblStyle *SVal          `xml:"w:tblStyle,omitempty"`
	TblW     *TableWidth    `xml:"w:tblW,omitempty"`
	Jc       *Justification `xml:"w:jc,omitempty"`
	TblBorders *TableBorders `xml:"w:tblBorders,omitempty"`
	TblLook  *TableLook     `xml:"w:tblLook,omitempty"`
}

type TableWidth struct {
	W    string `xml:"w:w,attr"`
	Type string `xml:"w:type,attr"`
}

type TableGrid struct {
	Cols []GridCol `xml:"w:gridCol"`
}

type GridCol struct {
	W string `xml:"w:w,attr"`
}

type TableBorders struct {
	Top     *BorderVal `xml:"w:top,omitempty"`
	Left    *BorderVal `xml:"w:left,omitempty"`
	Bottom  *BorderVal `xml:"w:bottom,omitempty"`
	Right   *BorderVal `xml:"w:right,omitempty"`
	InsideH *BorderVal `xml:"w:insideH,omitempty"`
	InsideV *BorderVal `xml:"w:insideV,omitempty"`
}

type BorderVal struct {
	Val   string `xml:"w:val,attr"`
	Sz    string `xml:"w:sz,attr,omitempty"`
	Space string `xml:"w:space,attr,omitempty"`
	Color string `xml:"w:color,attr,omitempty"`
}

type TableLook struct {
	Val          string `xml:"w:val,attr"`
	FirstRow     string `xml:"w:firstRow,attr,omitempty"`
	LastRow      string `xml:"w:lastRow,attr,omitempty"`
	FirstColumn  string `xml:"w:firstColumn,attr,omitempty"`
	LastColumn   string `xml:"w:lastColumn,attr,omitempty"`
	NoHBand      string `xml:"w:noHBand,attr,omitempty"`
	NoVBand      string `xml:"w:noVBand,attr,omitempty"`
}

// TableRow represents w:tr
type TableRow struct {
	XMLName xml.Name           `xml:"w:tr"`
	TrPr    *TableRowProperties `xml:"w:trPr,omitempty"`
	Cells   []TableCell        `xml:",any"`
}

type TableRowProperties struct {
	TblHeader *Empty `xml:"w:tblHeader,omitempty"`
}

// TableCell represents w:tc
type TableCell struct {
	XMLName xml.Name             `xml:"w:tc"`
	TcPr    *TableCellProperties `xml:"w:tcPr,omitempty"`
	Content []any                `xml:",any"`
}

type TableCellProperties struct {
	TcW   *TableWidth `xml:"w:tcW,omitempty"`
	Shd   *Shading    `xml:"w:shd,omitempty"`
	VAlign *SVal      `xml:"w:vAlign,omitempty"`
}

type Shading struct {
	Val   string `xml:"w:val,attr"`
	Color string `xml:"w:color,attr,omitempty"`
	Fill  string `xml:"w:fill,attr,omitempty"`
}

// Drawing represents w:drawing for inline images
type Drawing struct {
	XMLName xml.Name      `xml:"w:drawing"`
	Inline  *InlineDrawing `xml:"wp:inline,omitempty"`
}

type InlineDrawing struct {
	DistT  string `xml:"distT,attr"`
	DistB  string `xml:"distB,attr"`
	DistL  string `xml:"distL,attr"`
	DistR  string `xml:"distR,attr"`
	Extent *Extent `xml:"wp:extent"`
	DocPr  *DocPr  `xml:"wp:docPr"`
	Graphic *Graphic `xml:"a:graphic"`
}

type Extent struct {
	CX string `xml:"cx,attr"`
	CY string `xml:"cy,attr"`
}

type DocPr struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Descr string `xml:"descr,attr,omitempty"`
}

type Graphic struct {
	XMLName     xml.Name     `xml:"a:graphic"`
	A           string       `xml:"xmlns:a,attr"`
	GraphicData *GraphicData `xml:"a:graphicData"`
}

type GraphicData struct {
	URI string `xml:"uri,attr"`
	Pic *Picture `xml:"pic:pic"`
}

type Picture struct {
	XMLName  xml.Name  `xml:"pic:pic"`
	PicNS    string    `xml:"xmlns:pic,attr"`
	NvPicPr  *NvPicPr  `xml:"pic:nvPicPr"`
	BlipFill *BlipFill `xml:"pic:blipFill"`
	SpPr     *ShapeProperties `xml:"pic:spPr"`
}

type NvPicPr struct {
	CNvPr    *DocPr  `xml:"pic:cNvPr"`
	CNvPicPr *Empty  `xml:"pic:cNvPicPr"`
}

type BlipFill struct {
	Blip    *Blip    `xml:"a:blip"`
	Stretch *Stretch `xml:"a:stretch"`
}

type Blip struct {
	Embed string `xml:"r:embed,attr"`
}

type Stretch struct {
	FillRect *Empty `xml:"a:fillRect"`
}

type ShapeProperties struct {
	Xfrm *Transform2D `xml:"a:xfrm,omitempty"`
	PrstGeom *PresetGeometry `xml:"a:prstGeom,omitempty"`
}

type Transform2D struct {
	Ext *Extent `xml:"a:ext"`
}

type PresetGeometry struct {
	Prst string `xml:"prst,attr"`
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
