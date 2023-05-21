/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/labring/sealvm/pkg/process"
	"github.com/spf13/cobra"
)

func newInspectCmd() *cobra.Command {
	var inspectCmd = &cobra.Command{
		Use:   "inspect",
		Short: "inspect the vm node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname := args[0]
			i, err := process.NewInterfaceFromName(name)
			if err != nil {
				return err
			}
			i.Inspect(hostname)
			return nil
		},
	}
	inspectCmd.Flags().StringVarP(&name, "name", "n", "default", "name of cluster to applied init action")
	return inspectCmd
}
