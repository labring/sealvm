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
	"context"
	"fmt"
	"github.com/labring/sealvm/pkg/ssh"
	"github.com/labring/sealvm/pkg/utils/exec"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
	"golang.org/x/sync/errgroup"
)

type multiPassAction struct {
	vm        *v1.VirtualMachine
	nameAndIp map[string]string
	client    *ssh.Exec
}

func (m *multiPassAction) Apply(action *v1.Action) error {
	action.Status.Phase = v1.ActionPhaseInProcess
	var err error
	defer func() {
		if err != nil {
			action.Status.Phase = v1.ActionPhaseFailed
			action.Status.Message = err.Error()
		}
	}()
	names, nameAndIPs := getNameAndIPs(action, m.vm)
	m.nameAndIp = nameAndIPs
	if len(names) == 0 {
		logger.Warn("lookup names is empty")
		return nil
	}
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
	for _, data := range action.Spec.Data {
		if data.ActionMount != nil {
			if err = m.Mount(names, data); err != nil {
				return err
			}
		}
		if data.ActionUmount != "" {
			if err = m.UnMount(names, data); err != nil {
				return err
			}
		}
		if data.ActionExec != "" {
			if err = m.Exec(names, data); err != nil {
				return err
			}
		}
		if data.ActionCopy != nil {
			if err = m.Copy(names, data); err != nil {
				return err
			}
		}
		if data.ActionCopyContent != nil {
			if err = m.CopyContent(names, data); err != nil {
				return err
			}
		}
	}
	action.Status.Phase = v1.ActionPhaseComplete
	return nil
}

func (m *multiPassAction) Mount(names []string, data v1.ActionData) error {
	if data.ActionMount == nil {
		logger.Warn("mount data is nil")
		return nil
	}
	if data.ActionMount.Source == "" || data.ActionMount.Target == "" {
		logger.Warn("mount data is empty source or target")
		return nil
	}
	eg, _ := errgroup.WithContext(context.Background())

	for _, name := range names {
		name := name
		eg.Go(func() error {
			return m.mount(name, data.ActionMount.Source, data.ActionMount.Target)
		})
	}
	return eg.Wait()
}
func (m *multiPassAction) UnMount(names []string, data v1.ActionData) error {
	if data.ActionUmount == "" {
		logger.Warn("unmount data is nil")
		return nil
	}
	eg, _ := errgroup.WithContext(context.Background())

	for _, name := range names {
		name := name
		eg.Go(func() error {
			return m.unmount(name, data.ActionUmount)
		})
	}
	return eg.Wait()
}
func (*multiPassAction) Exec(names []string, data v1.ActionData) error {
	if data.ActionExec == "" {
		logger.Warn("exec data is nil")
		return nil
	}
	return nil
}
func (*multiPassAction) Copy(names []string, data v1.ActionData) error {
	if data.ActionCopy == nil {
		logger.Warn("copy data is nil")
		return nil
	}
	return nil
}
func (*multiPassAction) CopyContent(names []string, data v1.ActionData) error {
	if data.ActionCopyContent == nil {
		logger.Warn("copy content data is nil")
		return nil
	}
	return nil
}
func (m *multiPassAction) mount(name, src, target string) error {
	cmd := fmt.Sprintf("multipass mount %s %s:%s", src, name, target)
	logger.Info("executing... %s \n", cmd)
	return exec.Cmd("bash", "-c", cmd)
}

func (m *multiPassAction) unmount(name, target string) error {
	cmd := fmt.Sprintf("multipass unmount  %s:%s", name, target)
	logger.Info("executing... %s \n", cmd)
	return exec.Cmd("bash", "-c", cmd)
}
