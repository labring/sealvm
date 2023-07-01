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

package vm

import (
	"github.com/labring/sealvm/pkg/apply/runtime"
	"github.com/labring/sealvm/pkg/configs"
	v1 "github.com/labring/sealvm/types/api/v1"
)

type VirtualMachine struct {
	Desired  *v1.VirtualMachine
	Current  *v1.VirtualMachine
	Config   configs.Interface
	DiffFunc runtime.Diff
	Interface
}

type Interface interface {
	CreateVM(infra *v1.VirtualMachine, host *v1.Host, index int) error
	DeleteVM(infra *v1.VirtualMachine, host *v1.VirtualMachineHostStatus) error
	Get(name, role string, index int) (string, error)
	List() (string, error)
	GetById(name string) (string, error)
	Inspect(name string, role v1.Host, index int) (*v1.VirtualMachineHostStatus, error)
}
