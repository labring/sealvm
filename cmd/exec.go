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

var roles string
var ips, hostnames []string

var exampleExec = `
exec to default cluster: default
	sealvm exec "cat /etc/hosts"
specify the cluster name(If there is only one cluster in the $HOME/.sealvm directory, it should be applied. ):
    sealvm exec -c my-cluster "cat /etc/hosts"
set role label to exec cmd:
    sealvm exec -c my-cluster -r master,node "cat /etc/hosts"
set ips to exec cmd:
    sealvm exec -c my-cluster --ips 172.16.1.38 "cat /etc/hosts"
`

func newExecCmd() *cobra.Command {
	var execCmd = &cobra.Command{
		Use:     "exec",
		Short:   "Execute shell command or script on specified nodes",
		Example: exampleExec,
		Args:    cobra.ExactArgs(1),
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
			return i.Exec(exec, args[0])
		},
	}
	execCmd.Flags().StringVarP(&name, "name", "n", "default", "name of cluster to applied exec action")
	execCmd.Flags().StringVarP(&roles, "roles", "r", "", "run command on nodes with role")
	execCmd.Flags().StringSliceVar(&ips, "ips", []string{}, "run command on nodes with ip address")
	execCmd.Flags().StringSliceVar(&hostnames, "hostnames", []string{}, "run command on nodes with hostnames")

	return execCmd
}

func init() {
	rootCmd.AddCommand(newExecCmd())
}
