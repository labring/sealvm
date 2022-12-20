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
	"errors"
	"github.com/labring/sealvm/pkg/configs"
	v1 "github.com/labring/sealvm/types/api/v1"
)

func NewInterfaceFromName(name string) (Interface, error) {
	cf := configs.NewVirtualMachineFile(name)
	err := cf.Process()
	if err != nil {
		return nil, err
	}
	i := cf.GetVirtualMachine()
	switch i.Spec.Type {
	case v1.MultipassType:
		return &mulitipass{vm: i}, nil
	default:
		return nil, errors.New("not support type:" + i.Spec.Type)
	}
}
