/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/labring/sealvm/pkg/apply"
	v1 "github.com/labring/sealvm/types/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
func newResetCmd() *cobra.Command {
	vm := v1.VirtualMachine{}
	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			applier, err := apply.NewApplierFromArgs(&vm)
			if err != nil {
				return err
			}
			return applier.Apply()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkInstall(vm.Spec.Type); err != nil {
				return err
			}
			t := metav1.Now()
			vm.DeletionTimestamp = &t
			return nil
		},
	}
	resetCmd.Flags().StringVarP(&vm.Name, "name", "n", "default", "name of cluster to applied init action")
	resetCmd.Flags().StringVarP(&vm.Spec.Type, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	return resetCmd
}
func init() {
	rootCmd.AddCommand(newResetCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
