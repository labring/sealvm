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

package apply

import (
	"errors"
	"fmt"
	"github.com/labring/sealvm/pkg/system"
	"github.com/labring/sealvm/pkg/template"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
	"strconv"
	"strings"
)

func GetDefaultImage() (string, error) {
	defaultImageLocal, _ := system.Get(system.DefaultImageKey)
	if defaultImageLocal != "" {
		return defaultImageLocal, nil
	}
	return "", nil
}

func ValidateTemplate(vm *v1.VirtualMachine) error {
	if vm.Spec.SSH.PublicFile == "" {
		return fmt.Errorf("public key is required,please set values using 'sealvm values set'")
	}
	if vm.Spec.SSH.PkFile == "" {
		return fmt.Errorf("private key is required,please set values using 'sealvm values set'")
	}
	tpl := template.NewTpl()
	logger.Debug("current vm roles", vm.GetRoles())
	for _, r := range vm.GetRoles() {
		_, err := tpl.Get(r)
		if err != nil {
			return fmt.Errorf("template role %s is not exist", r)
		}
	}
	return nil
}

// node:3,master:5
func ParseNodes(nodes string) (map[string]int, error) {
	if nodes == "" {
		return nil, errors.New("nodes is required")
	}
	nodeMap := make(map[string]int)
	for _, node := range strings.Split(nodes, ",") {
		nodeArr := strings.Split(node, ":")
		if len(nodeArr) != 2 {
			return nil, errors.New("nodes format is wrong")
		}
		nodeInt, _ := strconv.Atoi(nodeArr[1])
		logger.Debug("node", nodeArr[0], "count", nodeInt)
		nodeMap[nodeArr[0]] = nodeInt
	}
	return nodeMap, nil
}

// node@/xx/xx:/ff/ff
func ParseMounts(mounts []string) (map[string]map[string]string, error) {
	if len(mounts) == 0 {
		return map[string]map[string]string{}, nil
	}
	mountMap := make(map[string]map[string]string)
	for _, mount := range mounts {
		mountArr := strings.Split(mount, "@")
		if len(mountArr) != 2 {
			return nil, errors.New("mount format is wrong")
		}
		mountRole := mountArr[0]
		mountPoint := strings.Split(mountArr[1], ":")
		if len(mountPoint) != 2 {
			return nil, errors.New("mount format is wrong")
		}
		if _, ok := mountMap[mountRole]; !ok {
			mountMap[mountRole] = map[string]string{mountPoint[0]: mountPoint[1]}
		} else {
			mountMap[mountRole][mountPoint[0]] = mountPoint[1]
		}
	}
	return mountMap, nil
}
