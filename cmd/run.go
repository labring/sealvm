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
	"fmt"
	"github.com/cuisongliu/sealos-dev/pkg/apply"
	fileutil "github.com/cuisongliu/sealos-dev/pkg/utils/file"
	v1 "github.com/cuisongliu/sealos-dev/types/api/v1"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// runCmd represents the run command
func newRunCmd() *cobra.Command {
	vm := v1.VirtualMachine{}
	var masters, nodes, registry int
	var dev bool
	var src string
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			applier, err := apply.NewApplierFromArgs(&vm)
			if err != nil {
				return err
			}
			return applier.Apply()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if dev {
				if src == "" {
					return fmt.Errorf("src must be set")
				}
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.DEV,
					Count:  1,
					Mounts: map[string]string{src: "/root/go/src"},
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			if registry != 0 {
				if registry != 1 {
					return fmt.Errorf("registry must be 1")
				}
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.REGISTRY,
					Count:  1,
					Mounts: map[string]string{},
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			if masters != 0 {
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.MASTER,
					Count:  masters,
					Mounts: map[string]string{},
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			if nodes != 0 {
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.NODE,
					Count:  nodes,
					Mounts: map[string]string{},
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			return nil
		},
	}
	runCmd.Flags().StringVarP(&vm.Spec.SSH.PkFile, "pk", "i", path.Join(fileutil.GetHomeDir(), ".ssh", "id_rsa"), "selects a file from which the identity (private key) for public key authentication is read")
	runCmd.Flags().StringVar(&vm.Spec.SSH.PkPasswd, "pk-passwd", "", "passphrase for decrypting a PEM encoded private key")
	runCmd.Flags().StringVarP(&vm.Spec.SSH.PublicFile, "pub", "p", path.Join(fileutil.GetHomeDir(), ".ssh", "id_rsa.pub"), "selects a file from which the identity (public key) for public key authentication is read")
	runCmd.Flags().StringVarP(&vm.Spec.Type, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	runCmd.Flags().StringVarP(&vm.Name, "name", "n", "default", "name of cluster to applied init action")
	runCmd.Flags().IntVarP(&masters, "masters", "m", 1, "number of masters")
	runCmd.Flags().IntVarP(&nodes, "nodes", "w", 1, "number of nodes")
	runCmd.Flags().IntVarP(&registry, "registry", "r", 1, "number of registry")
	runCmd.Flags().BoolVarP(&dev, "dev", "d", true, "number of dev")
	runCmd.Flags().StringVarP(&src, "dev-mounts", "s", path.Join(os.Getenv("GOPATH"), "src"), "gopath src dir")
	return runCmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
