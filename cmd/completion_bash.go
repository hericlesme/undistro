package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completions",
	Long:  "Generate bash completions",
	RunE: func(cmd *cobra.Command, cmdArgs []string) error {
		return cmd.Root().GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
}
