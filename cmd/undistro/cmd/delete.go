/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete workload clusters.",
	Long:  `Delete workload clusters.`,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
