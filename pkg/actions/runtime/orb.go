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
	"fmt"
	"github.com/labring/sealvm/pkg/ssh"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
	"k8s.io/apimachinery/pkg/util/errors"
)

type orbAction struct {
	multiPassAction
}

func (m *orbAction) Apply(action *v1.Action) error {
	action.Status.Phase = v1.ActionPhaseInProcess
	var err error
	defer func() {
		if err != nil {
			action.Status.Phase = v1.ActionPhaseFailed
			switch err.(type) {
			case errors.Aggregate:
				action.Status.Message = err.(errors.Aggregate).Error()
			default:
				action.Status.Message = err.Error()
			}

		}
	}()
	names, nameAndIPs := getNameAndIPs(action, m.vm)
	m.nameAndIp = nameAndIPs
	if len(names) == 0 {
		logger.Warn("lookup names is empty")
		return nil
	}
	logger.Info("lookup names: %v", nameAndIPs)
	ips := make([]string, 0)
	for _, name := range names {
		if _, ok := nameAndIPs[name]; !ok {
			return fmt.Errorf("name %s not found", name)
		}
		ips = append(ips, nameAndIPs[name])
	}
	var execClient *ssh.Exec
	execClient, err = ssh.NewExecCmdFromIPs(m.vm, ips)
	if err != nil {
		return err
	}
	m.client = execClient
	fns := []func(names []string, data v1.ActionData) error{
		m.Mount,
		m.UnMount,
		m.Exec,
		m.Copy,
		m.CopyContent,
	}
	errArr := make([]error, 0)
	for _, data := range action.Spec.Data {
		for _, fn := range fns {
			fnErr := fn(names, data)
			if fnErr != nil {
				errArr = append(errArr, fnErr)
				break
			}
		}
	}
	if len(errArr) > 0 {
		err = errors.NewAggregate(errArr)
		return err
	}
	action.Status.Phase = v1.ActionPhaseComplete
	return nil
}

func (m *orbAction) mount(name, src, target string) error {
	logger.Warn("orb does not need to support mount, it is mounted in the root directory by default")
	return nil
}

func (m *orbAction) unmount(name, target string) error {
	logger.Warn("orb does not need to support mount, it is mounted in the root directory by default")
	return nil
}
