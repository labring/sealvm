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
	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/ssh"
	"github.com/labring/sealvm/pkg/system"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/util/errors"
	"os"
	"path"
)

type Interface interface {
	MountOnce(name, src, target string) error
	UnMountOnce(name, target string) error
	Copy(names []string, data v1.ActionData) error
	Exec(names []string, data v1.ActionData) error
}

type action struct {
	vm        *v1.VirtualMachine
	nameAndIp map[string]string
	client    *ssh.Exec
	Interface
}

func (m *action) Apply(action *v1.Action) error {
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
	defaultProvider, _ := system.Get(system.DefaultProvider)
	var ii Interface
	switch defaultProvider {
	case v1.MultipassType:
		execClient, err = ssh.NewExecCmdFromIPs(m.vm, ips)
		if err != nil {
			return err
		}
		m.client = execClient
		ii = newMultiPassAction(m.client)
	case v1.OrbType:
		ii = newOrbAction()
	default:
		return fmt.Errorf("action not support type: %s", defaultProvider)
	}
	m.Interface = ii
	fns := []func(names []string, data v1.ActionData) error{
		m.Mount,
		m.UnMount,
		m.Exec,
		m.Copy,
		m.CopyContent,
	}
	for _, data := range action.Spec.Data {
		for _, fn := range fns {
			fnErr := fn(names, data)
			if fnErr != nil {
				err = fnErr
				return err
			}
		}
	}
	action.Status.Phase = v1.ActionPhaseComplete
	return nil
}

func (m *action) Mount(names []string, data v1.ActionData) error {
	if data.ActionMount == nil {
		return nil
	}
	if data.ActionMount.Source == "" || data.ActionMount.Target == "" {
		return fmt.Errorf("mount data is empty source or target")
	}
	eg, _ := errgroup.WithContext(context.Background())

	for _, name := range names {
		name := name
		eg.Go(func() error {
			logger.Debug("mount %s %s:%s", data.ActionMount.Source, name, data.ActionMount.Target)
			return m.MountOnce(name, data.ActionMount.Source, data.ActionMount.Target)
		})
	}
	return eg.Wait()
}
func (m *action) UnMount(names []string, data v1.ActionData) error {
	if data.ActionUmount == "" {
		return nil
	}
	eg, _ := errgroup.WithContext(context.Background())

	for _, name := range names {
		name := name
		eg.Go(func() error {
			logger.Debug("unmount %s:%s", name, data.ActionUmount)
			return m.UnMountOnce(name, data.ActionUmount)
		})
	}
	return eg.Wait()
}

func (m *action) CopyContent(names []string, data v1.ActionData) error {
	if data.ActionCopyContent == nil {
		return nil
	}
	if data.ActionCopyContent.Target == "" {
		return fmt.Errorf("copy data is empty target")
	}
	tmpDir := path.Join(configs.DefaultRootfsDir(), "tmp")
	_ = os.MkdirAll(tmpDir, 0755)
	newDir, _ := fileutil.MkTmpdir(tmpDir)
	defer func() {
		_ = os.RemoveAll(newDir)
	}()
	newFile := path.Join(newDir, "action-generator.sh")
	_ = fileutil.WriteFile(newFile, []byte(data.ActionCopyContent.Content))
	logger.Debug("copy content to %s", data.ActionCopyContent.Target)
	newData := v1.ActionData{
		ActionCopy: &v1.SourceAndTarget{
			Source: newFile,
			Target: data.ActionCopyContent.Target,
		},
	}
	return m.Copy(names, newData)
}
