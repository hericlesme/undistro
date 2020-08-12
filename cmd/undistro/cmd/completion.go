package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generators for shell completions",
	Long:  "Generators for shell completions",
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
