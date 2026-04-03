package style

// Style is the top-level YAML style definition.
type Style struct {
	Meta             MetaConfig           `yaml:"meta"`
	Page             PageConfig           `yaml:"page"`
	Fonts            FontsConfig          `yaml:"fonts"`
	Styles           StylesConfig         `yaml:"styles"`
	HeadingNumbering HeadingNumbering     `yaml:"heading_numbering"`
	TOC              TOCConfig            `yaml:"toc"`
	Cover            CoverConfig          `yaml:"cover"`
	Header           HeaderFooterConfig   `yaml:"header"`
	Footer           HeaderFooterConfig   `yaml:"footer"`
	List             ListConfig           `yaml:"list"`
}

type MetaConfig struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Locale      string `yaml:"locale"`
}

type PageConfig struct {
	Size        string       `yaml:"size"`
	Width       string       `yaml:"width"`
	Height      string       `yaml:"height"`
	Orientation string       `yaml:"orientation"`
	Margin      MarginConfig `yaml:"margin"`
}

type MarginConfig struct {
	Top    string `yaml:"top"`
	Bottom string `yaml:"bottom"`
	Left   string `yaml:"left"`
	Right  string `yaml:"right"`
	Header string `yaml:"header"`
	Footer string `yaml:"footer"`
}

type FontsConfig struct {
	Latin string `yaml:"latin"`
	CJK   string `yaml:"cjk"`
	Mono  string `yaml:"mono"`
}

type StylesConfig struct {
	Heading1  ParagraphStyle `yaml:"heading1"`
	Heading2  ParagraphStyle `yaml:"heading2"`
	Heading3  ParagraphStyle `yaml:"heading3"`
	Heading4  ParagraphStyle `yaml:"heading4"`
	Heading5  ParagraphStyle `yaml:"heading5"`
	Heading6  ParagraphStyle `yaml:"heading6"`
	Body      ParagraphStyle `yaml:"body"`
	Quote     ParagraphStyle `yaml:"quote"`
	CodeBlock ParagraphStyle `yaml:"code_block"`
	Caption   ParagraphStyle `yaml:"caption"`
	Table     TableStyle     `yaml:"table"`
	Image     ImageStyle     `yaml:"image"`
}

type ParagraphStyle struct {
	FontCJK          string  `yaml:"font_cjk"`
	FontLatin        string  `yaml:"font_latin"`
	FontSize         string  `yaml:"font_size"`
	Bold             bool    `yaml:"bold"`
	Italic           bool    `yaml:"italic"`
	Color            string  `yaml:"color"`
	Alignment        string  `yaml:"alignment"`
	LineSpacing      float64 `yaml:"line_spacing"`
	SpaceBefore      string  `yaml:"space_before"`
	SpaceAfter       string  `yaml:"space_after"`
	FirstLineIndent  string  `yaml:"first_line_indent"`
	LeftIndent       string  `yaml:"left_indent"`
	RightIndent      string  `yaml:"right_indent"`
	PageBreakBefore  bool    `yaml:"page_break_before"`
	KeepWithNext     bool    `yaml:"keep_with_next"`
}

type TableStyle struct {
	FontSize    string `yaml:"font_size"`
	Alignment   string `yaml:"alignment"`
	HeaderBold  bool   `yaml:"header_bold"`
	HeaderBg    string `yaml:"header_bg"`
	HeaderColor string `yaml:"header_color"`
	Border      bool   `yaml:"border"`
	BorderColor string `yaml:"border_color"`
}

type ImageStyle struct {
	MaxWidthPct int    `yaml:"max_width_pct"`
	Alignment   string `yaml:"alignment"`
}

type HeadingNumbering struct {
	Enabled bool              `yaml:"enabled"`
	Formats map[int]string    `yaml:"formats"`
}

type TOCConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Title          string `yaml:"title"`
	Depth          int    `yaml:"depth"`
	PageBreakAfter bool   `yaml:"page_break_after"`
}

type CoverConfig struct {
	Enabled        bool          `yaml:"enabled"`
	PageBreakAfter bool          `yaml:"page_break_after"`
	Layout         []CoverElement `yaml:"layout"`
}

type CoverElement struct {
	Type      string `yaml:"type"`      // text, image, spacer
	Source    string `yaml:"source"`    // frontmatter field name or literal
	Height    string `yaml:"height"`    // for spacer
	FontCJK   string `yaml:"font_cjk"`
	FontSize  string `yaml:"font_size"`
	Bold      bool   `yaml:"bold"`
	Alignment string `yaml:"alignment"`
}

type HeaderFooterConfig struct {
	Left     string `yaml:"left"`
	Right    string `yaml:"right"`
	Center   string `yaml:"center"`
	FontSize string `yaml:"font_size"`
}

type ListConfig struct {
	BulletChars     []string `yaml:"bullet_chars"`
	OrderedFormat   string   `yaml:"ordered_format"`
	IndentPerLevel  string   `yaml:"indent_per_level"`
}
