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
	"github.com/labring/sealvm/pkg/utils/exec"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"
	"golang.org/x/sync/errgroup"
	"strings"
)

func newOrbAction() Interface {
	return &orbAction{}
}

type orbAction struct {
	multiPassAction
}

func (m *orbAction) MountOnce(name, src, target string) error {
	logger.Warn("orb does not need to support mount, it is mounted in the root directory by default")
	return nil
}

func (m *orbAction) UnMountOnce(name, target string) error {
	logger.Warn("orb does not need to support mount, it is mounted in the root directory by default")
	return nil
}

func (m *orbAction) Exec(names []string, data v1.ActionData) error {
	if data.ActionExec == "" {
		return nil
	}
	logger.Debug("names %+v,exec %s", names, data.ActionExec)

	for _, name := range names {
		for _, cmd := range strings.Split(data.ActionExec, "\n") {
			if strings.TrimSpace(cmd) == "" {
				continue
			}
			err := exec.Cmd("/bin/bash", "-c", fmt.Sprintf("ssh root@%s@orb \"%s\"", name, cmd))
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (m *orbAction) Copy(names []string, data v1.ActionData) error {
	if data.ActionCopy == nil {
		return nil
	}
	if data.ActionCopy.Source == "" || data.ActionCopy.Target == "" {
		return fmt.Errorf("copy data is empty source or target")
	}
	logger.Debug("names %+v,copy from %s to %s", names, data.ActionCopy.Source, data.ActionCopy.Target)
	eg, _ := errgroup.WithContext(context.Background())
	for _, name := range names {
		name := name
		eg.Go(func() error {
			err := exec.Cmd("/bin/bash", "-c", fmt.Sprintf("scp %s root@%s@orb:%s", data.ActionCopy.Source, name, data.ActionCopy.Target))
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to exec command, err: %v", err)
	}
	return nil
}
