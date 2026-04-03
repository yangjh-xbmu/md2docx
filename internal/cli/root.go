package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yangjh-xbmu/md2docx/internal/parser"
	"github.com/yangjh-xbmu/md2docx/internal/render"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

var (
	// Version is set via ldflags at build time.
	Version   = "dev"
	outputOpt string
	styleOpt  string
	noCover   bool
	noTOC     bool
	noNumber  bool
)

var rootCmd = &cobra.Command{
	Use:   "md2docx [file.md]",
	Short: "Markdown to Word converter",
	Long: `md2docx 将 Markdown 文件转换为排版精良的 Word (.docx) 文档。

支持 CJK 双字体、目录生成、标题编号、封面页、页眉页脚、
GFM 表格、图片嵌入等功能。内置多套样式，零配置即可使用。`,
	Example: `  md2docx 论文.md
  md2docx input.md -o output.docx
  md2docx input.md --style academic-cn
  md2docx input.md --no-cover --no-toc`,
	Version: Version,
	Args:    cobra.MaximumNArgs(1),
	RunE:    runConvert,
}

func init() {
	rootCmd.Flags().StringVarP(&outputOpt, "output", "o", "", "输出 .docx 路径（默认: ~/Desktop/<同名>.docx）")
	rootCmd.Flags().StringVarP(&styleOpt, "style", "s", "", "样式名称（默认: default）")
	rootCmd.Flags().BoolVar(&noCover, "no-cover", false, "跳过封面")
	rootCmd.Flags().BoolVar(&noTOC, "no-toc", false, "跳过目录")
	rootCmd.Flags().BoolVar(&noNumber, "no-numbering", false, "跳过标题编号")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runConvert(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	source := args[0]
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("无法读取文件 %s: %w", source, err)
	}

	// Parse frontmatter
	meta, body := parser.ParseFrontmatter(string(content))

	// Determine style
	styleName := styleOpt
	if styleName == "" {
		if s, ok := meta["style"].(string); ok {
			styleName = s
		} else {
			styleName = "default"
		}
	}

	// Load style
	s, err := style.Load(styleName)
	if err != nil {
		available := style.ListAvailable()
		return fmt.Errorf("样式 '%s' 不存在。可用样式: %v", styleName, available)
	}

	// Apply frontmatter overrides
	style.ApplyFrontmatterOverrides(s, meta)

	// Apply CLI overrides
	if noTOC {
		s.TOC.Enabled = false
	}
	if noCover {
		s.Cover.Enabled = false
	}
	if noNumber {
		s.HeadingNumbering.Enabled = false
	}

	// Determine output path
	output := outputOpt
	if output == "" {
		base := filepath.Base(source)
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]
		desktop := filepath.Join(os.Getenv("HOME"), "Desktop")
		if info, err := os.Stat(desktop); err == nil && info.IsDir() {
			output = filepath.Join(desktop, name+".docx")
		} else {
			output = filepath.Join(filepath.Dir(source), name+".docx")
		}
	}

	// Parse markdown
	ast, src := parser.ParseMarkdown(body)

	// Render to docx
	err = render.ToDocx(ast, src, s, meta, output, filepath.Dir(source))
	if err != nil {
		return fmt.Errorf("转换失败: %w", err)
	}

	fmt.Fprintf(os.Stdout, "已导出: %s\n", output)
	return nil
}
