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

package infra

import (
	"fmt"
	"github.com/cuisongliu/sealos-dev/pkg/configs"
	"github.com/cuisongliu/sealos-dev/pkg/infra/mulitipass"
	v1 "github.com/cuisongliu/sealos-dev/types/api/v1"
)

type Interface interface {
	Apply() error
}

func NewDefaultVirtualMachine(infra *v1.VirtualMachine, cf configs.Interface) (Interface, error) {
	if infra.Spec.Type != v1.MultipassType {
		return nil, fmt.Errorf("infra type %s is not supported", infra.Spec.Type)
	}
	return newMultiPassVirtualMachine(infra, cf)
}

func newMultiPassVirtualMachine(infra *v1.VirtualMachine, cf configs.Interface) (Interface, error) {
	if infra.Name == "" {
		return nil, fmt.Errorf("infra name cannot be empty")
	}
	if cf == nil {
		cf = configs.NewVirtualMachineFile(infra.Name)
	}
	err := cf.Process()
	if !infra.CreationTimestamp.IsZero() && err != nil {
		return nil, err
	}

	return &mulitipass.MultiPassVirtualMachine{
		Desired: infra,
		Current: cf.GetVirtualMachine(),
		Config:  cf,
	}, nil
}
