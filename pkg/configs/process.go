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

package configs

import (
	"bytes"
	"errors"
	fileutil "github.com/labring/sealos-vm/pkg/utils/file"
	"github.com/labring/sealos-vm/pkg/utils/logger"
	v1 "github.com/labring/sealos-vm/types/api/v1"
	"github.com/mitchellh/go-homedir"
	"k8s.io/apimachinery/pkg/util/yaml"
	"path"
)

var ErrVirtualMachineFileNotExists = errors.New("the vm file is not exist")
var ErrTypeNotFound = errors.New("no type structure was found")

type PreProcessor interface {
	Process() error
}

var DefaultClusterRootfsDir = ""

func DefaultRootfsDir() string {
	if DefaultClusterRootfsDir != "" {
		return DefaultClusterRootfsDir
	}
	home, _ := homedir.Dir()
	return path.Join(home, ".sealos-vm")
}

func GetDataDir(clusterName string) string {
	return path.Join(DefaultRootfsDir(), "data", clusterName)
}

func GetEtcDir(clusterName string) string {
	return path.Join(DefaultRootfsDir(), "etc", clusterName)
}

func VirtualMachineFilePath(clusterName string) string {
	return path.Join(GetDataDir(clusterName), "VirtualMachineFile")
}

func (c *VirtualMachineFile) Process() (err error) {
	if !fileutil.IsExist(VirtualMachineFilePath(c.name)) {
		return ErrVirtualMachineFileNotExists
	}
	c.once.Do(func() {
		err = func() error {
			//for i := range c.customEnvs {
			//	kv := strings.SplitN(c.customEnvs[i], "=", 2)
			//	if len(kv) == 2 {
			//		_ = os.Setenv(kv[0], kv[1])
			//	}
			//}
			clusterFileData, err := c.loadVirtualMachineFile()
			if err != nil {
				return err
			}
			return c.decode(clusterFileData)
		}()
	})
	return
}

func (c *VirtualMachineFile) loadVirtualMachineFile() ([]byte, error) {
	body, err := fileutil.ReadAll(VirtualMachineFilePath(c.name))
	if err != nil {
		return nil, err
	}
	logger.Debug("loadVirtualMachineFile body: %+v", string(body))
	out := bytes.NewBuffer(nil)
	out.Write(body)
	return out.Bytes(), nil
}
func (c *VirtualMachineFile) decode(data []byte) error {
	for _, fn := range []func([]byte) error{
		c.DecodeVirtualMachine,
	} {
		if err := fn(data); err != nil && err != ErrTypeNotFound {
			return err
		}
	}
	return nil
}

func (c *VirtualMachineFile) DecodeVirtualMachine(data []byte) error {
	vm, err := GetVirtualMachineFromDataCompatV1(data)
	if err != nil {
		return err
	}
	c.VirtualMachine = vm
	return nil
}

func GetVirtualMachineFromDataCompatV1(data []byte) (*v1.VirtualMachine, error) {
	var vm *v1.VirtualMachine
	err := yaml.Unmarshal(data, &vm)
	if err != nil {
		return nil, err
	}
	vm.TypeMeta.APIVersion = v1.GroupVersion.String()
	vm.TypeMeta.Kind = "VirtualMachine"
	return vm, nil
}
