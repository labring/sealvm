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
	"github.com/labring/sealvm/pkg/ssh"
	v1 "github.com/labring/sealvm/types/api/v1"
)

type mulitipass struct {
	vm *v1.VirtualMachine
}

func (mp *mulitipass) List() error {
	return printVMs(mp.vm)
}

func (mp *mulitipass) Exec(exec ssh.Exec, shell string) error {
	return exec.RunCmd(shell)
}

func (*mulitipass) Transfer(exec ssh.Exec, src, dstFilePath string) error {
	return exec.RunCopy(src, dstFilePath)
}
func (mp *mulitipass) Inspect(name string) {
	inspectHostname(mp.vm, name)
}

func (mp *mulitipass) VMInfo() *v1.VirtualMachine {
	return mp.vm
}
