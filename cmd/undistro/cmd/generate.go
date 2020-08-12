/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate yaml using undistro yaml processor.",
	Long:  `Generate yaml using undistro yaml processor.`,
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
