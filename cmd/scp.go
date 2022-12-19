/*
Copyright 2022 cuisongliu@qq.com.

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
	"github.com/labring/sealvm/pkg/process"
	"github.com/labring/sealvm/pkg/ssh"
	"github.com/labring/sealvm/pkg/utils/strings"
	"github.com/spf13/cobra"
)

// Shared with exec.go
// var roles string
// var clusterName string
// var ips []string

const exampleScp = `
copy file to default cluster: default
	sealvm scp "/root/aa.txt" "/root/dd.txt"
specify the cluster name(If there is only one cluster in the $HOME/.sealvm directory, it should be applied. ):
    sealvm scp -c my-cluster "/root/aa.txt" "/root/dd.txt"
set role label to copy file:
    sealvm scp -c my-cluster -r master,node "cat /etc/hosts"
set ips to copy file:
    sealvm scp -c my-cluster --ips 172.16.1.38  "/root/aa.txt" "/root/dd.txt"
`

func newScpCmd() *cobra.Command {
	var scpCmd = &cobra.Command{
		Use:     "scp",
		Short:   "Copy file to remote on specified nodes",
		Example: exampleScp,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			i, err := process.NewInterfaceFromName(name)
			if err != nil {
				return err
			}
			if len(hostnames) > 0 {
				for _, h := range i.VMInfo().Status.Hosts {
					if len(h.IPs[0]) > 0 && strings.In(h.ID, hostnames) {
						if !strings.In(h.IPs[0], ips) {
							ips = append(ips, h.IPs[0])
						}
					}
				}
			}
			var exec ssh.Exec
			if len(ips) > 0 {
				exec, err = ssh.NewExecCmdFromIPs(i.VMInfo(), ips)
				if err != nil {
					return err
				}

			}
			exec, err = ssh.NewExecCmdFromRoles(i.VMInfo(), roles)
			if err != nil {
				return err
			}
			return exec.RunCopy(args[0], args[1])
		},
	}
	scpCmd.Flags().StringVarP(&name, "name", "n", "default", "name of cluster to applied exec action")
	scpCmd.Flags().StringVarP(&roles, "roles", "r", "", "copy file to nodes with role")
	scpCmd.Flags().StringSliceVar(&ips, "ips", []string{}, "copy file to nodes with ip address")
	scpCmd.Flags().StringSliceVar(&hostnames, "hostnames", []string{}, "copy file to nodes with hostnames")
	return scpCmd
}

func init() {
	rootCmd.AddCommand(newScpCmd())
}
