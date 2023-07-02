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
	"github.com/labring/sealvm/pkg/utils/exec"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
)

func newMultiPassAction(client *ssh.Exec) Interface {
	return &multiPassAction{
		client: client,
	}
}

type multiPassAction struct {
	client *ssh.Exec
}

func (m *multiPassAction) Exec(names []string, data v1.ActionData) error {
	if data.ActionExec == "" {
		return nil
	}
	logger.Debug("names %+v,exec %s", names, data.ActionExec)
	return m.client.RunCmd(data.ActionExec)
}
func (m *multiPassAction) Copy(names []string, data v1.ActionData) error {
	if data.ActionCopy == nil {
		return nil
	}
	if data.ActionCopy.Source == "" || data.ActionCopy.Target == "" {
		return fmt.Errorf("copy data is empty source or target")
	}
	logger.Debug("names %+v,copy from %s to %s", names, data.ActionCopy.Source, data.ActionCopy.Target)
	return m.client.RunCopy(data.ActionCopy.Source, data.ActionCopy.Target)
}

func (m *multiPassAction) MountOnce(name, src, target string) error {
	cmd := fmt.Sprintf("multipass mount %s %s:%s", src, name, target)
	logger.Info("executing... %s \n", cmd)
	return exec.Cmd("bash", "-c", cmd)
}

func (m *multiPassAction) UnMountOnce(name, target string) error {
	cmd := fmt.Sprintf("multipass unmount  %s:%s", name, target)
	logger.Info("executing... %s \n", cmd)
	return exec.Cmd("bash", "-c", cmd)
}
