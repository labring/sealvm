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
	"errors"
	"fmt"
	"github.com/labring/sealvm/pkg/system"
	"github.com/labring/sealvm/pkg/template"
	"github.com/labring/sealvm/pkg/utils/logger"
	"path/filepath"
	"strings"

	"github.com/labring/sealvm/pkg/apply"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	v1 "github.com/labring/sealvm/types/api/v1"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func newRunCmd() *cobra.Command {
	vm := v1.VirtualMachine{}
	val := template.NewValues()
	var nodes int
	var dev bool
	//var defaultMount = fmt.Sprintf("%s:%s", path.Join(os.Getenv("GOPATH"), "src"), "/root/go/src")
	var defaultImage string
	var mounts []string
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run cloud native vm nodes",
		Long:  `sealvm run --nodes 1  20.04`,
		RunE: func(cmd *cobra.Command, args []string) error {
			applier, err := apply.NewApplierFromArgs(&vm)
			if err != nil {
				return err
			}
			return applier.Apply()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				defaultImage = args[0]
			} else {
				defaultImageLocal, _ := system.Get(system.DefaultImageKey)
				if defaultImageLocal != "" {
					if !fileutil.IsExist(defaultImageLocal) {
						return fmt.Errorf("system config image not set or image file is not exist")
					}
					if !filepath.IsAbs(defaultImageLocal) {
						return errors.New("image local must using abs path")
					}
					defaultImage = fmt.Sprintf("file://%s", defaultImageLocal)
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
			var mountPoints map[string]string
			for _, m := range mounts {
				if strings.TrimSpace(m) != "" {
					points := strings.Split(m, ":")
					if len(points) != 2 {
						return fmt.Errorf("mount args format is error , ex /dd:/ff")
					}
					mountPoints[points[0]] = points[1]
				}
			}

			if err := checkInstall(vm.Spec.Type); err != nil {
				return err
			}
			if dev {

				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.GOLANG,
					Count:  1,
					Mounts: mountPoints,
					Resources: map[string]string{
						v1.CPUKey:  defaultCpuNum,
						v1.DISKKey: defaultDiskGb,
						v1.MEMKey:  defaultMemoryGb,
					},
					Image: defaultImage,
				})
			}
			if nodes != 0 {
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.NODE,
					Count:  nodes,
					Mounts: mountPoints,
					Resources: map[string]string{
						v1.CPUKey:  defaultCpuNum,
						v1.DISKKey: defaultDiskGb,
						v1.MEMKey:  defaultMemoryGb,
					},
					Image: defaultImage,
				})
			}
			vm.Spec.SSH.PublicFile = val.Get("PublicKey")
			if vm.Spec.SSH.PublicFile == "" {
				return fmt.Errorf("public key is required,please set values using 'sealvm values set'")
			}
			vm.Spec.SSH.PkFile = val.Get("PrivateKey")
			if vm.Spec.SSH.PkFile == "" {
				return fmt.Errorf("private key is required,please set values using 'sealvm values set'")
			}
			tpl := template.NewTpl()
			for _, r := range vm.GetRoles() {
				_, err := tpl.Get(r)
				if err != nil {
					return fmt.Errorf("template role %s is not exist", r)
				}
			}
			return nil
		},
	}
	runCmd.Flags().StringVar(&vm.Spec.SSH.PkPasswd, "pk-passwd", "", "passphrase for decrypting a PEM encoded private key")
	runCmd.Flags().StringVarP(&vm.Spec.Type, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	runCmd.Flags().StringVar(&vm.Name, "name", "default", "name of cluster to applied init action")

	runCmd.Flags().IntVarP(&nodes, "nodes", "n", 0, "number of nodes")
	runCmd.Flags().StringSliceVarP(&mounts, "mounts", "m", []string{}, "mounts for vm")
	runCmd.Flags().BoolVarP(&dev, "dev", "d", false, "number of dev")
	//runCmd.Flags().StringVarP(&src, "dev-mounts", "s", defaultMount, "gopath src dir")
	//runCmd.Flags().IntVarP(&defaultCpuNum, "default-node-cpu", "c", 2, "default vcpu num per node. ")
	//runCmd.Flags().IntVarP(&defaultMemoryGb, "default-node-mem", "m", 4, "default mem size per node. （GB） ")
	//runCmd.Flags().IntVarP(&defaultDiskGb, "default-node-disk", "k", 50, "default disk size per node. （GB）")
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
