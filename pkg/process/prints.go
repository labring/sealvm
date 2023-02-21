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

package process

import (
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	"github.com/modood/table"
)

func printVMs(vm *v1.VirtualMachine) error {
	type printTable struct {
		Name  string
		State string
		Role  string
		Ipv4  []string
		Image string
	}
	tables := make([]printTable, 0)
	for _, h := range vm.Spec.Hosts {
		if h.Count > 0 {
			for i := 0; i < h.Count; i++ {
				status := vm.GetHostStatusByRoleIndex(h.Role, i)
				if status == nil {
					tables = append(tables, printTable{
						Name:  strings.GetID(vm.Name, h.Role, i),
						State: "UNKNOWN",
						Role:  h.Role,
					})
				} else {
					tables = append(tables, printTable{
						Ipv4:  status.IPs,
						Name:  status.ID,
						Image: status.ImageName,
						State: status.State,
						Role:  h.Role,
					})
				}
			}
		}
	}
	table.OutputA(tables)
	return nil
}

func inspectHostname(vm *v1.VirtualMachine, hostname string) {
	type printTable struct {
		Name string
		Info any
	}
	tables := make([]printTable, 0)
	if len(vm.Status.Hosts) > 0 {
		for _, h := range vm.Status.Hosts {
			if h.ID == hostname {
				tables = append(tables, printTable{
					Name: "Name",
					Info: h.ID,
				})

				tables = append(tables, printTable{
					Name: "State",
					Info: h.State,
				})

				tables = append(tables, printTable{
					Name: "IPv4",
					Info: h.IPs,
				})

				tables = append(tables, printTable{
					Name: "Release",
					Info: h.ImageName,
				})

				tables = append(tables, printTable{
					Name: "Mounts",
					Info: h.Mounts,
				})
				tables = append(tables, printTable{
					Name: "Capacity",
					Info: h.Capacity,
				})
				tables = append(tables, printTable{
					Name: "Used",
					Info: h.Used,
				})
			}
		}
	}
	table.OutputA(tables)
}
