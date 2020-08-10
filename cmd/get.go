/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get info from a management or a workload cluster",
	Long:  `Get info from a management or a workload cluster`,
}

func init() {
	RootCmd.AddCommand(getCmd)
}
