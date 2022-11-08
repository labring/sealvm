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

package apply

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/labring/sealosvm/pkg/apply/infra"
	"github.com/labring/sealosvm/pkg/apply/runtime"
	"github.com/labring/sealosvm/pkg/configs"
	"github.com/labring/sealosvm/pkg/utils/logger"
	v1 "github.com/labring/sealosvm/types/api/v1"
)

func NewApplierFromArgs(args *v1.VirtualMachine) (runtime.Interface, error) {
	name := args.Name
	cf := configs.NewVirtualMachineFile(name)
	err := cf.Process()
	if err != nil && err != configs.ErrVirtualMachineFileNotExists {
		return nil, err
	}
	i := cf.GetVirtualMachine()
	if i == nil {
		logger.Debug("current VirtualMachine is nil")
		i = initVirtualMachine(name)
	}

	if err := mergo.Merge(i, args); err != nil {
		return nil, fmt.Errorf("merge: %v", err)
	}

	return infra.NewDefaultVirtualMachine(i, cf)
}

func initVirtualMachine(clusterName string) *v1.VirtualMachine {
	i := &v1.VirtualMachine{}
	i.Name = clusterName
	i.Kind = "VirtualMachine"
	i.APIVersion = v1.GroupVersion.String()
	return i
}
