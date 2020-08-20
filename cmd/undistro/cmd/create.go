/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create things in undistro",
	Long:  `Create things in undistro`,
}

func init() {
	RootCmd.AddCommand(createCmd)
}
