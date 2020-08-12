package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completions",
	Long:  "Generate zsh completions",
	RunE: func(cmd *cobra.Command, cmdArgs []string) error {
		return cmd.Root().GenZshCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionZshCmd)
}
