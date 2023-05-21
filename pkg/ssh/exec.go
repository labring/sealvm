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

package ssh

import (
	"context"
	"fmt"
	v1 "github.com/labring/sealvm/types/api/v1"

	"github.com/labring/sealvm/pkg/utils/logger"

	"golang.org/x/sync/errgroup"
)

type Exec struct {
	vm     *v1.VirtualMachine
	ipList []string
	client Interface
}

func NewExecCmdFromIPs(vm *v1.VirtualMachine, ips []string) (*Exec, error) {
	sshClient := NewSSHClient(&vm.Spec.SSH, true)
	err := WaitSSHReady(sshClient, 6, ips...)
	if err != nil {
		return nil, err
	}
	return &Exec{vm: vm, ipList: ips, client: sshClient}, nil
}

func (e *Exec) RunCmd(cmd string) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, ipAddr := range e.ipList {
		ip := ipAddr
		eg.Go(func() error {
			err := e.client.CmdAsync(ip, cmd)
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

func (e *Exec) RunCopy(srcFilePath, dstFilePath string) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, ipAddr := range e.ipList {
		ip := ipAddr
		eg.Go(func() error {
			err := e.client.Copy(ip, srcFilePath, dstFilePath)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to copy command, err: %v", err)
	}
	logger.Info("transfers files success")
	return nil
}
