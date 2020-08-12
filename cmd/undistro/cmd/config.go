/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display provider configuration and templates to create workload clusters.",
	Long:  `Display provider configuration and templates to create workload clusters.`,
}

func init() {
	RootCmd.AddCommand(configCmd)
}
