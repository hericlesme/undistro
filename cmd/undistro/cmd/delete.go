/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete things in undistro.",
	Long:  `Delete things in undistro.`,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
