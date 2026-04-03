package render

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/yangjh-xbmu/md2docx/internal/ooxml"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

// imageEntry holds data for one embedded image.
type imageEntry struct {
	relID       string
	partName    string
	contentType string
	data        []byte
	widthEmu    int64
	heightEmu   int64
}

// loadImage reads an image file, determines its size, and returns an imageEntry.
// It scales the image so its width does not exceed maxWidthPct% of the page content width.
func loadImage(src string, baseDir string, pageWidthTwips int, marginLeftTwips int, marginRightTwips int, imgStyle style.ImageStyle, imgIndex int) (*imageEntry, error) {
	// Resolve path relative to the markdown file
	imgPath := src
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(baseDir, imgPath)
	}

	data, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, fmt.Errorf("读取图片失败 %s: %w", src, err)
	}

	// Detect content type from extension
	ext := strings.ToLower(filepath.Ext(imgPath))
	contentType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".png":
		contentType = "image/png"
	}

	// Get image dimensions
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("打开图片失败 %s: %w", src, err)
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("解析图片尺寸失败 %s: %w", src, err)
	}

	widthPx := cfg.Width
	heightPx := cfg.Height
	if widthPx == 0 || heightPx == 0 {
		widthPx = 400
		heightPx = 300
	}

	// Calculate max width in EMU
	// Page content width = page width - left margin - right margin (all in twips)
	contentWidthTwips := pageWidthTwips - marginLeftTwips - marginRightTwips
	if contentWidthTwips <= 0 {
		contentWidthTwips = 9000 // fallback ~15.9cm
	}

	maxPct := imgStyle.MaxWidthPct
	if maxPct <= 0 || maxPct > 100 {
		maxPct = 80
	}

	// 1 twip = 914400/1440 = 635 EMU
	contentWidthEmu := int64(contentWidthTwips) * 635
	maxWidthEmu := contentWidthEmu * int64(maxPct) / 100

	// Convert pixel dimensions to EMU (96 DPI: 1px = 9525 EMU)
	widthEmu := ooxml.EmuFromPx(widthPx)
	heightEmu := ooxml.EmuFromPx(heightPx)

	// Scale down if exceeds max width
	if widthEmu > maxWidthEmu {
		ratio := float64(maxWidthEmu) / float64(widthEmu)
		widthEmu = maxWidthEmu
		heightEmu = int64(float64(heightEmu) * ratio)
	}

	partName := fmt.Sprintf("word/media/image%d%s", imgIndex, ext)
	relID := fmt.Sprintf("rIdImg%d", imgIndex)

	return &imageEntry{
		relID:       relID,
		partName:    partName,
		contentType: contentType,
		data:        data,
		widthEmu:    widthEmu,
		heightEmu:   heightEmu,
	}, nil
}

// buildImageParagraph creates a paragraph containing an inline image drawing.
func buildImageParagraph(img *imageEntry, alt string, imgStyle style.ImageStyle) *ooxml.Paragraph {
	ppr := &ooxml.ParagraphProperties{}
	if imgStyle.Alignment != "" {
		ppr.Jc = &ooxml.Justification{Val: ooxml.AlignmentVal(imgStyle.Alignment)}
	}

	drawing := &ooxml.Drawing{
		Inline: &ooxml.InlineDrawing{
			DistT: "0", DistB: "0", DistL: "0", DistR: "0",
			Extent: &ooxml.Extent{
				CX: ooxml.FormatEmu(img.widthEmu),
				CY: ooxml.FormatEmu(img.heightEmu),
			},
			DocPr: &ooxml.DocPr{
				ID:    img.relID,
				Name:  "Image",
				Descr: alt,
			},
			Graphic: &ooxml.Graphic{
				A: "http://schemas.openxmlformats.org/drawingml/2006/main",
				GraphicData: &ooxml.GraphicData{
					URI: "http://schemas.openxmlformats.org/drawingml/2006/picture",
					Pic: &ooxml.Picture{
						PicNS: "http://schemas.openxmlformats.org/drawingml/2006/picture",
						NvPicPr: &ooxml.NvPicPr{
							CNvPr:    &ooxml.DocPr{ID: "0", Name: "Image"},
							CNvPicPr: &ooxml.Empty{},
						},
						BlipFill: &ooxml.BlipFill{
							Blip:    &ooxml.Blip{Embed: img.relID},
							Stretch: &ooxml.Stretch{FillRect: &ooxml.Empty{}},
						},
						SpPr: &ooxml.ShapeProperties{
							Xfrm: &ooxml.Transform2D{
								Ext: &ooxml.Extent{
									CX: ooxml.FormatEmu(img.widthEmu),
									CY: ooxml.FormatEmu(img.heightEmu),
								},
							},
							PrstGeom: &ooxml.PresetGeometry{Prst: "rect"},
						},
					},
				},
			},
		},
	}

	run := &ooxml.Run{
		Content: []any{drawing},
	}

	return &ooxml.Paragraph{
		PPr:  ppr,
		Runs: []any{run},
	}
}
