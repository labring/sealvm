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

package mulitipass

import (
	"context"
	"fmt"
	"github.com/labring/sealos-vm/pkg/configs"
	"github.com/labring/sealos-vm/pkg/tmpl"
	"github.com/labring/sealos-vm/pkg/utils/exec"
	"github.com/labring/sealos-vm/pkg/utils/logger"
	v1 "github.com/labring/sealos-vm/types/api/v1"
	"golang.org/x/sync/errgroup"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
)

func (r *MultiPassVirtualMachine) init() {
	logger.Info("Start to create a new infra:", r.Desired.Name)

	pipelines := []func(infra *v1.VirtualMachine){
		r.InitStatus,
		r.ApplyConfig,
		r.ApplyVMs,
		//r.TransferSSHKey,
		r.SyncVMs,
		r.FinalStatus,
	}

	for _, fn := range pipelines {
		fn(r.Desired)
	}
	if r.Desired.Status.Phase != v1.PhaseFailed {
		r.Desired.Status.Phase = v1.PhaseSuccess
	}
	logger.Info("succeeded in creating a new infra, enjoy it!")

}

var defaultNodesYaml = path.Join(configs.DefaultRootfsDir(), "etc", "nodes.yaml")
var defaultGolangYaml = path.Join(configs.DefaultRootfsDir(), "etc", "golang.yaml")

func (r *MultiPassVirtualMachine) InitStatus(infra *v1.VirtualMachine) {
	logger.Info("Start to exec InitStatus:", r.Desired.Name)
	var initializedCondition = &v1.Condition{
		Type:              "Initialized",
		Status:            v12.ConditionTrue,
		Reason:            "Initialized",
		Message:           "infra has been initialized",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, initializedCondition)
	infra.Status.Phase = v1.PhaseInProcess
}

func (r *MultiPassVirtualMachine) ApplyConfig(infra *v1.VirtualMachine) {
	logger.Info("Start to exec ApplyConfig:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "Config",
		Status:            v12.ConditionTrue,
		Reason:            "Config Generated",
		Message:           "config has been generated to launch multipass",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	if err := tmpl.ExecuteNodesToFile(infra.Spec.Proxy, infra.Spec.NoProxy, infra.Spec.SSH.PkFile, infra.Spec.SSH.PublicFile, defaultNodesYaml); err != nil {
		v1.SetConditionError(configCondition, "ConfigNodeGenerateError", err)
		return
	}

	if err := tmpl.ExecuteGolangToFile(infra.Spec.Proxy, infra.Spec.NoProxy, infra.Spec.SSH.PkFile, infra.Spec.SSH.PublicFile, defaultGolangYaml); err != nil {
		v1.SetConditionError(configCondition, "ConfigGolangGenerateError", err)
	}
}

func (r *MultiPassVirtualMachine) ApplyVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec ApplyVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "VMs",
		Status:            v12.ConditionTrue,
		Reason:            "VM start",
		Message:           "launch multipass success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	//sshClient, sshErr := ssh.NewSSHByVirtualMachine(infra, true)
	//if sshErr != nil {
	//	v1.SetConditionError(configCondition, "GetSSH", sshErr)
	//	return
	//}

	eg, _ := errgroup.WithContext(context.Background())

	for _, host := range infra.Spec.Hosts {
		for i := 0; i < host.Count; i++ {
			dHost := host
			index := i
			eg.Go(func() error {
				return r.CreateVM(infra, &dHost, index)
			})
		}

	}
	if err := eg.Wait(); err != nil {
		v1.SetConditionError(configCondition, "CreateVMError", err)
		return
	}
}

func (r *MultiPassVirtualMachine) TransferSSHKey(infra *v1.VirtualMachine) {
	if !v1.IsConditionsTrue(infra.Status.Conditions) {
		logger.Info("Skip to exec TransferSSHKey:", r.Desired.Name)
		return
	}
	logger.Info("Start to exec TransferSSHKey:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "TransferSSHKey",
		Status:            v12.ConditionTrue,
		Reason:            "TransferSSHKey to instance",
		Message:           "transfer multipass success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	//sshClient, sshErr := ssh.NewSSHByVirtualMachine(infra, true)
	//if sshErr != nil {
	//	v1.SetConditionError(configCondition, "GetSSH", sshErr)
	//	return
	//}

	eg, _ := errgroup.WithContext(context.Background())

	for _, host := range infra.Spec.Hosts {
		for i := 0; i < host.Count; i++ {
			dHost := host
			index := i
			eg.Go(func() error {
				return r.Transfer(infra, &dHost, index)
			})
		}
	}
	if err := eg.Wait(); err != nil {
		v1.SetConditionError(configCondition, "TransferSSHKeyError", err)
		return
	}
}

func (r *MultiPassVirtualMachine) SyncVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec SyncVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "SyncVMs",
		Status:            v12.ConditionTrue,
		Reason:            "VM status sync",
		Message:           "multipass instance sync success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	//sshClient, sshErr := ssh.NewSSHByVirtualMachine(infra, true)
	//if sshErr != nil {
	//	v1.SetConditionError(configCondition, "GetSSH", sshErr)
	//	return
	//}
	var status []v1.VirtualMachineHostStatus
	for _, host := range infra.Spec.Hosts {
		for i := 0; i < host.Count; i++ {
			info, err := r.Inspect(infra.Name, host.Role, i)
			if err != nil {
				v1.SetConditionError(configCondition, "GetVM", err)
				continue
			}
			if info.State != "Running" {
				v1.SetConditionError(configCondition, "VMStatus", fmt.Errorf("vm status is not running"))
			}
			status = append(status, *info)
		}

	}
	infra.Status.Hosts = status
}

func (r *MultiPassVirtualMachine) CreateVM(infra *v1.VirtualMachine, host *v1.Host, index int) error {
	cfg := defaultNodesYaml
	if v1.DEV == host.Role {
		cfg = defaultGolangYaml
	}
	cmd := fmt.Sprintf("multipass launch --name %s-%s-%d --cpus %d --mem %dG --disk %dG --cloud-init %s", infra.Name, host.Role, index, host.Resources[v1.CPUKey], host.Resources[v1.MEMKey], host.Resources[v1.DISKKey], cfg)
	return exec.Cmd("bash", "-c", cmd)
}

func (r *MultiPassVirtualMachine) Transfer(infra *v1.VirtualMachine, host *v1.Host, index int) error {
	var cmds []string
	cmds = append(cmds, fmt.Sprintf("multipass transfer %s %s-%s-%d:/home/ubuntu/", infra.Spec.SSH.PkFile, infra.Name, host.Role, index))
	cmds = append(cmds, fmt.Sprintf("multipass transfer %s %s-%s-%d:/home/ubuntu/", infra.Spec.SSH.PublicFile, infra.Name, host.Role, index))
	cmds = append(cmds, fmt.Sprintf("multipass exec %s-%s-%d -- sudo -u root rm -rf ~/.ssh && mkdir ~/.ssh &&  cp -r /home/ubuntu/id_rsa ~/.ssh/id_rsa && cp -r /home/ubuntu/id_rsa.pub ~/.ssh/id_rsa.pub && chmod 600 ~/.ssh/id_rsa && cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys", infra.Name, host.Role, index))
	for _, cmd := range cmds {
		if err := exec.Cmd("bash", "-c", cmd); err != nil {
			return err
		}
	}
	return nil
}

func (r *MultiPassVirtualMachine) FinalStatus(infra *v1.VirtualMachine) {
	condition := &v1.Condition{
		Type:              "Ready",
		Status:            v12.ConditionTrue,
		LastHeartbeatTime: metav1.Now(),
		Reason:            "Ready",
		Message:           "MultiPass is available now",
	}
	defer r.saveCondition(infra, condition)

	if !v1.IsConditionsTrue(infra.Status.Conditions) {
		condition.LastHeartbeatTime = metav1.Now()
		condition.Status = v12.ConditionFalse
		condition.Reason = "NotReady"
		condition.Message = "MultiPass is not available now"
		infra.Status.Phase = v1.PhaseFailed
	} else {
		infra.Status.Phase = v1.PhaseSuccess
	}
}

// Language: go
func (r *MultiPassVirtualMachine) saveCondition(infra *v1.VirtualMachine, condition *v1.Condition) {
	if !v1.IsConditionTrue(infra.Status.Conditions, *condition) {
		infra.Status.Conditions = v1.UpdateCondition(infra.Status.Conditions, *condition)
	}
}
