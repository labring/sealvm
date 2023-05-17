/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/labring/sealvm/pkg/process"
	"github.com/spf13/cobra"
)

// listCmd represents the list command

func newListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			i, err := process.NewInterfaceFromName(name)
			if err != nil {
				return err
			}
			return i.List()
		},
	}
	listCmd.Flags().StringVarP(&name, "name", "n", "default", "name of cluster to applied init action")
	return listCmd
}
