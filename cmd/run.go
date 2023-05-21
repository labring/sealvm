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
	"github.com/labring/sealvm/pkg/system"
	"github.com/labring/sealvm/pkg/template"
	"github.com/labring/sealvm/pkg/utils/confirm"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
	"strings"

	"github.com/labring/sealvm/pkg/apply"
	v1 "github.com/labring/sealvm/types/api/v1"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func newRunCmd() *cobra.Command {
	vm := v1.VirtualMachine{}
	val := template.NewValues()
	var nodes string
	//var defaultMount = fmt.Sprintf("%s:%s", path.Join(os.Getenv("GOPATH"), "src"), "/root/go/src")
	var defaultImage string
	var runCmd = &cobra.Command{
		Use:     "run",
		Short:   "Run cloud native vm nodes",
		Example: `sealvm run -n node:2,master:1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			applier, err := apply.NewApplierFromArgs(&vm)
			if err != nil {
				return err
			}
			return applier.Apply()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			system.List()
			template.NewTpl().List()
			template.NewValues().List()
			if len(args) != 0 {
				defaultImage = args[0]
			} else {
				newDefaultImage, err := apply.GetDefaultImage()
				if err != nil {
					return err
				}
				if newDefaultImage != "" {
					defaultImage = newDefaultImage
				}
			}
			//var defaultCpuNum int
			//var defaultDiskGb int
			//var defaultMemoryGb int
			defaultCpuNum, _ := system.Get(system.DefaultCPUKey)
			defaultDiskGb, _ := system.Get(system.DefaultDISKKey)
			defaultMemoryGb, _ := system.Get(system.DefaultMemKey)
			logger.Debug("default cpu number is %s", defaultCpuNum)
			logger.Debug("default disk number is %s", defaultDiskGb)
			logger.Debug("default disk memory is %s", defaultMemoryGb)
			if strings.Contains(vm.Name, "-") {
				return fmt.Errorf("your cluster name contains chart '-' ")
			}
			if err := checkInstall(vm.Spec.Type); err != nil {
				return err
			}

			nodeMap, err := apply.ParseNodes(nodes)
			if err != nil {
				return errors.WithMessage(err, "parse nodes error")
			}

			for n, node := range nodeMap {
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:  n,
					Count: node,
					Resources: map[string]string{
						v1.CPUKey:  defaultCpuNum,
						v1.DISKKey: defaultDiskGb,
						v1.MEMKey:  defaultMemoryGb,
					},
					Image: defaultImage,
				})
			}
			vm.Spec.SSH.PublicFile = val.Get("PublicKey")
			vm.Spec.SSH.PkFile = val.Get("PrivateKey")
			data, err := yaml.Marshal(&vm)
			if err != nil {
				return err
			}
			logger.Debug("vm yaml is %s", string(data))
			if yes, err := confirm.Confirm("Are you sure to run this command?", "you have canceled to crate these vms !"); err != nil {
				return err
			} else {
				if !yes {
					return errors.New("cancelled")
				}
			}
			err = apply.ValidateTemplate(&vm)
			if err != nil {
				return err
			}
			return nil
		},
	}
	runCmd.Flags().StringVar(&vm.Spec.SSH.PkPasswd, "pk-passwd", "", "passphrase for decrypting a PEM encoded private key")
	runCmd.Flags().StringVarP(&vm.Spec.Type, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	runCmd.Flags().StringVar(&vm.Name, "name", "default", "name of cluster to applied init action")
	runCmd.Flags().StringVarP(&nodes, "nodes", "n", "", "number of nodes, eg: node:1,node2:2")
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
