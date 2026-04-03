package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yangjh-xbmu/md2docx/internal/merge"
	"github.com/yangjh-xbmu/md2docx/internal/parser"
	"github.com/yangjh-xbmu/md2docx/internal/render"
	"github.com/yangjh-xbmu/md2docx/internal/style"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <files...>",
	Short: "合并多个 Markdown 文件为一个 Word 文档",
	Long: `合并多个 Markdown 文件为一个 Word 文档。

支持两种输入方式：
  1. 直接指定多个 .md 文件
  2. 指定一个 contents.yaml 配置文件

示例：
  md2docx merge ch1.md ch2.md ch3.md -o book.docx
  md2docx merge contents.yaml -o book.docx
  md2docx merge "docs/*.md" --style academic-cn`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&outputOpt, "output", "o", "", "输出 .docx 路径")
	mergeCmd.Flags().StringVarP(&styleOpt, "style", "s", "", "样式名称（默认: default）")
	mergeCmd.Flags().BoolVar(&noCover, "no-cover", false, "跳过封面")
	mergeCmd.Flags().BoolVar(&noTOC, "no-toc", false, "跳过目录")
	mergeCmd.Flags().BoolVar(&noNumber, "no-numbering", false, "跳过标题编号")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	var entries []merge.Entry

	// Determine input mode: contents.yaml or file list/globs
	if len(args) == 1 && merge.IsContentsYAML(args[0]) {
		var err error
		entries, err = merge.ParseContentsYAML(args[0])
		if err != nil {
			return err
		}
	} else {
		// Check if args are globs or direct files
		var paths []string
		for _, arg := range args {
			// Try glob expansion
			matches, err := filepath.Glob(arg)
			if err != nil {
				return fmt.Errorf("无效的模式 '%s': %w", arg, err)
			}
			if len(matches) == 0 {
				// Treat as literal path (will error later if not found)
				paths = append(paths, arg)
			} else {
				paths = append(paths, matches...)
			}
		}

		// Deduplicate and filter .md
		seen := make(map[string]bool)
		for _, p := range paths {
			abs, _ := filepath.Abs(p)
			if seen[abs] {
				continue
			}
			seen[abs] = true
			entries = append(entries, merge.Entry{Path: p, HeadingOffset: 0})
		}
	}

	if len(entries) == 0 {
		return fmt.Errorf("没有找到要合并的文件")
	}

	// Merge markdown content
	merged, err := merge.MergeFiles(entries)
	if err != nil {
		return err
	}

	// Parse frontmatter from merged content
	meta, body := parser.ParseFrontmatter(merged)

	// Determine style
	styleName := styleOpt
	if styleName == "" {
		if s, ok := meta["style"].(string); ok {
			styleName = s
		} else {
			styleName = "default"
		}
	}

	s, err := style.Load(styleName)
	if err != nil {
		available := style.ListAvailable()
		return fmt.Errorf("样式 '%s' 不存在。可用样式: %v", styleName, available)
	}

	style.ApplyFrontmatterOverrides(s, meta)

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
		desktop := filepath.Join(os.Getenv("HOME"), "Desktop")
		if info, err := os.Stat(desktop); err == nil && info.IsDir() {
			output = filepath.Join(desktop, "merged.docx")
		} else {
			output = "merged.docx"
		}
	}

	// Determine base dir for relative paths (use first entry's dir)
	baseDir := filepath.Dir(entries[0].Path)

	// Parse and render
	ast, src := parser.ParseMarkdown(body)
	if err := render.ToDocx(ast, src, s, meta, output, baseDir); err != nil {
		return fmt.Errorf("转换失败: %w", err)
	}

	fmt.Fprintf(os.Stdout, "已合并 %d 个文件: %s\n", len(entries), output)
	return nil
}
