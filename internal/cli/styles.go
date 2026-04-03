package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yangjh-xbmu/md2docx/internal/style"
	"gopkg.in/yaml.v3"
)

var stylesCmd = &cobra.Command{
	Use:   "styles",
	Short: "样式管理",
}

var stylesListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有可用样式",
	RunE: func(cmd *cobra.Command, args []string) error {
		names := style.ListAvailable()
		for _, name := range names {
			s, err := style.Load(name)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %-15s %s\n", name, s.Meta.Description)
		}
		return nil
	},
}

var stylesShowCmd = &cobra.Command{
	Use:   "show [style-name]",
	Short: "显示样式详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		s, err := style.Load(name)
		if err != nil {
			return fmt.Errorf("样式 '%s' 不存在。可用样式: %v", name, style.ListAvailable())
		}

		data, err := yaml.Marshal(s)
		if err != nil {
			return fmt.Errorf("序列化样式失败: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "# 样式: %s\n\n%s", name, string(data))
		return nil
	},
}

func init() {
	stylesCmd.AddCommand(stylesListCmd)
	stylesCmd.AddCommand(stylesShowCmd)
	rootCmd.AddCommand(stylesCmd)
}
