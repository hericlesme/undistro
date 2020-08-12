package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completions",
	Long:  "Generate fish completions",
	RunE: func(cmd *cobra.Command, cmdArgs []string) error {
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	},
}

func init() {
	completionCmd.AddCommand(completionFishCmd)
}
