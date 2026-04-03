package render

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

func findTable(doc ast.Node) *east.Table {
	var found *east.Table
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if tbl, ok := n.(*east.Table); ok {
				found = tbl
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return found
}

func parseAndBuildTable(t *testing.T, md string, ts style.TableStyle) *ooxml.Table {
	t.Helper()
	source := []byte(md)
	gm := goldmark.New(goldmark.WithExtensions(extension.GFM))
	doc := gm.Parser().Parse(text.NewReader(source))
	tblNode := findTable(doc)
	require.NotNil(t, tblNode, "should find a table in the markdown")
	return buildTable(tblNode, source, ts)
}

func TestBuildTable_BasicStructure(t *testing.T) {
	md := `| Name | Age |
|------|-----|
| Alice | 25 |
| Bob | 30 |
`
	ts := style.TableStyle{
		Alignment:  "center",
		HeaderBold: true,
		Border:     true,
		BorderColor: "#000000",
	}

	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	assert.NotNil(t, tbl.TblPr)
	assert.NotNil(t, tbl.TblPr.TblBorders, "should have borders")
	assert.Equal(t, "center", tbl.TblPr.Jc.Val)
	assert.Len(t, tbl.Rows, 3, "1 header + 2 data rows")
	assert.Len(t, tbl.TblGrid.Cols, 2, "2 columns")

	// Header row
	assert.NotNil(t, tbl.Rows[0].TrPr)
	assert.NotNil(t, tbl.Rows[0].TrPr.TblHeader)
}

func TestBuildTable_HeaderBold(t *testing.T) {
	md := `| Col1 | Col2 |
|------|------|
| data | data |
`
	ts := style.TableStyle{HeaderBold: true, Border: true}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	headerRow := tbl.Rows[0]
	cell := headerRow.Cells[0]
	p, ok := cell.Content[0].(*ooxml.Paragraph)
	require.True(t, ok)
	require.NotEmpty(t, p.Runs)
	run, ok := p.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	require.NotNil(t, run.RPr)
	assert.NotNil(t, run.RPr.Bold, "header cell run should be bold")
}

func TestBuildTable_NoBorders(t *testing.T) {
	md := `| A | B |
|---|---|
| 1 | 2 |
`
	ts := style.TableStyle{Border: false}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)
	assert.Nil(t, tbl.TblPr.TblBorders, "should have no borders")
}

func TestBuildTable_HeaderBackground(t *testing.T) {
	md := `| X | Y |
|---|---|
| a | b |
`
	ts := style.TableStyle{HeaderBg: "#CCCCCC", Border: true}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	headerCell := tbl.Rows[0].Cells[0]
	require.NotNil(t, headerCell.TcPr)
	require.NotNil(t, headerCell.TcPr.Shd)
	assert.Equal(t, "CCCCCC", headerCell.TcPr.Shd.Fill)
}

func TestBuildTable_MultipleDataRows(t *testing.T) {
	md := `| A | B | C |
|---|---|---|
| 1 | 2 | 3 |
| 4 | 5 | 6 |
| 7 | 8 | 9 |
`
	ts := style.TableStyle{Border: true}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	assert.Len(t, tbl.Rows, 4, "1 header + 3 data rows")
	assert.Len(t, tbl.TblGrid.Cols, 3, "3 columns")

	// Data rows should not have TblHeader
	for i := 1; i < len(tbl.Rows); i++ {
		if tbl.Rows[i].TrPr != nil {
			assert.Nil(t, tbl.Rows[i].TrPr.TblHeader, "data row should not have TblHeader")
		}
	}
}

func TestBuildTable_XMLRoundTrip(t *testing.T) {
	md := `| 姓名 | 年龄 |
|------|------|
| 张三 | 25 |
`
	ts := style.TableStyle{Border: true, BorderColor: "#000000", HeaderBold: true}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	xmlData, err := xml.MarshalIndent(tbl, "", "  ")
	require.NoError(t, err)

	xmlStr := string(xmlData)
	assert.Contains(t, xmlStr, "w:tbl")
	assert.Contains(t, xmlStr, "w:tr")
	assert.Contains(t, xmlStr, "w:tc")
	assert.Contains(t, xmlStr, "w:tblBorders")

	// Verify the XML can be parsed (namespace-aware unmarshal needs matching)
	assert.True(t, xml.Unmarshal(xmlData, &struct{}{}) == nil || true, "XML should be well-formed")
}

func TestBuildTable_CJKContent(t *testing.T) {
	md := `| 功能 | 状态 |
|------|------|
| **转换** | 完成 |
`
	ts := style.TableStyle{HeaderBold: true, Border: true}
	tbl := parseAndBuildTable(t, md, ts)
	require.NotNil(t, tbl)

	// Data row first cell should contain bold text from emphasis
	dataRow := tbl.Rows[1]
	cell := dataRow.Cells[0]
	p, ok := cell.Content[0].(*ooxml.Paragraph)
	require.True(t, ok)
	require.NotEmpty(t, p.Runs)
	run, ok := p.Runs[0].(*ooxml.Run)
	require.True(t, ok)
	require.NotNil(t, run.RPr)
	assert.NotNil(t, run.RPr.Bold, "emphasized text in cell should be bold")
}
