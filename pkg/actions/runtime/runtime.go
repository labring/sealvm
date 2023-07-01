/*
Copyright 2023 cuisongliu@qq.com.

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

package runtime

import (
	"errors"
	"github.com/labring/sealvm/pkg/process"
	"github.com/labring/sealvm/pkg/system"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Runtime interface {
	Apply(action *v1.Action) error
}

func getNameAndIPs(action *v1.Action, vm *v1.VirtualMachine) ([]string, map[string]string) {
	if action == nil {
		return nil, nil
	}
	data := make(map[string]string)
	defer func() {
		logger.Debug("get name and ips: %v", data)
	}()
	names := sets.NewString()
	for _, on := range action.Spec.Ons {
		if len(on.Indexes) == 0 {
			h := vm.GetHostByRole(on.Role)
			if h != nil {
				for i := 0; i < h.Count; i++ {
					names.Insert(strings.GetID(vm.Name, on.Role, i))
				}
			}
		} else {
			for _, i := range on.Indexes {
				names.Insert(strings.GetID(vm.Name, on.Role, int(i)))
			}
		}
	}

	if len(names.List()) == 0 {
		return nil, nil
	}

	for _, name := range names.List() {
		status := vm.GetHostStatusByName(name)
		if status != nil {
			if len(status.IPs) == 0 {
				continue
			}
			data[name] = status.IPs[0]
		}
	}
	return names.List(), data
}

func NewAction(name string) (Runtime, error) {
	i, err := process.NewInterfaceFromName(name)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if i.VMInfo() != nil {
		defaultProvider, _ := system.Get(system.DefaultProvider)
		switch defaultProvider {
		case v1.MultipassType:
			return &multiPassAction{vm: i.VMInfo()}, nil
		default:
			return nil, errors.New("action not support type:" + defaultProvider)
		}
	}
	return nil, errors.New("load vm config error")
}
