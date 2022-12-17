// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configs

import (
	"sync"

	v1 "github.com/labring/sealvm/types/api/v1"
)

type VirtualMachineFile struct {
	name           string
	VirtualMachine *v1.VirtualMachine
	once           sync.Once
}

type Interface interface {
	PreProcessor
	GetVirtualMachine() *v1.VirtualMachine
}

func (c *VirtualMachineFile) GetVirtualMachine() *v1.VirtualMachine {
	return c.VirtualMachine
}

type OptionFunc func(*VirtualMachineFile)

func NewVirtualMachineFile(name string, opts ...OptionFunc) Interface {
	cf := &VirtualMachineFile{
		name: name,
	}
	for _, opt := range opts {
		opt(cf)
	}
	return cf
}
