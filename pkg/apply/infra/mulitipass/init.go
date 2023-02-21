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
	"path"

	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/ssh"
	"github.com/labring/sealvm/pkg/tmpl"
	"github.com/labring/sealvm/pkg/utils/exec"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	"golang.org/x/sync/errgroup"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (r *MultiPassVirtualMachine) DesiredVM() *v1.VirtualMachine {
	return r.Desired
}

func (r *MultiPassVirtualMachine) CurrentVM() *v1.VirtualMachine {
	return r.Current
}

func (r *MultiPassVirtualMachine) Init() {
	logger.Info("Start to create a new infra:", r.Desired.Name)

	pipelines := []func(infra *v1.VirtualMachine){
		r.InitStatus,
		r.ApplyConfig,
		r.CreateVMs,
		r.MountsVMs,
		r.SyncVMs,
		r.PingVms,
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

func GetCloudInitYamlByRole(clusterName, role string) string {

	return path.Join(configs.GetEtcDir(clusterName), fmt.Sprintf("%s.yaml", role))
}

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
	if !fileutil.IsExist(configs.GetDataDir(infra.Name)) {
		_ = fileutil.MkDirs(configs.GetDataDir(infra.Name))
	}
	if err := tmpl.ExecuteNodesToFile(infra.Spec.Proxy, infra.Spec.NoProxy, infra.Spec.SSH.PkFile, infra.Spec.SSH.PublicFile, GetCloudInitYamlByRole(infra.Name, v1.NODE)); err != nil {
		v1.SetConditionError(configCondition, "ConfigNodeGenerateError", err)
		return
	}

	if err := tmpl.ExecuteGolangToFile(infra.Spec.Proxy, infra.Spec.NoProxy, infra.Spec.SSH.PkFile, infra.Spec.SSH.PublicFile, GetCloudInitYamlByRole(infra.Name, v1.GOLANG)); err != nil {
		v1.SetConditionError(configCondition, "ConfigGolangGenerateError", err)
	}
}

func (r *MultiPassVirtualMachine) CreateVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec CreateVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "InitVMs",
		Status:            v12.ConditionTrue,
		Reason:            "VM start",
		Message:           "launch multipass success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

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

func (r *MultiPassVirtualMachine) MountsVMs(infra *v1.VirtualMachine) {
	if !v1.IsConditionsTrue(infra.Status.Conditions) {
		logger.Info("Skip to exec MountsVMs:", r.Desired.Name)
		return
	}
	logger.Info("Start to exec MountsVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "MountsVMs",
		Status:            v12.ConditionTrue,
		Reason:            "MountsVMs to instance",
		Message:           "mount multipass success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	eg, _ := errgroup.WithContext(context.Background())

	for _, host := range infra.Spec.Hosts {
		for i := 0; i < host.Count; i++ {
			dHost := host
			index := i
			eg.Go(func() error {
				if dHost.Mounts != nil {
					for h, m := range dHost.Mounts {
						fn := func() {
							cmd := fmt.Sprintf("multipass mount %s %s:%s", h, strings.GetID(infra.Name, dHost.Role, index), m)
							logger.Info("executing... %s \n", cmd)
							_ = exec.Cmd("bash", "-c", cmd)
						}
						hostBool := len(infra.Status.Hosts) == 0
						if hostBool {
							fn()
						} else {
							if _, ok := infra.Status.Hosts[index].Mounts[h]; !ok {
								fn()
							}
						}
					}
				}
				return nil
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

	var status []v1.VirtualMachineHostStatus
	for _, host := range infra.Spec.Hosts {
		for i := 0; i < host.Count; i++ {
			//retry
			var info *v1.VirtualMachineHostStatus
			var err error
			if e := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				info, err = r.Inspect(infra.Name, host, i)
				if err != nil {
					newInfo, nee := r.InspectByList(infra.Name, host, i)
					if nee != nil {
						return nee
					}
					info = newInfo
				}
				if info.State != "Running" {
					return fmt.Errorf("instance %s is not running", infra.Name)
				}
				return nil
			}); e != nil {
				v1.SetConditionError(configCondition, "VMStatus", fmt.Errorf("vm %s status is not running", strings.GetID(infra.Name, host.Role, i)))
				continue
			}
			status = append(status, *info)
		}
	}

	infra.Status.Hosts = status
}

func (r *MultiPassVirtualMachine) PingVms(infra *v1.VirtualMachine) {
	if !v1.IsConditionsTrue(infra.Status.Conditions) {
		logger.Info("Skip to exec PingVms:", r.Desired.Name)
		return
	}
	logger.Info("Start to exec PingVms:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "PingVms",
		Status:            v12.ConditionTrue,
		Reason:            "VM ssh ping",
		Message:           "multipass instance ssh ping success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)
	client := ssh.NewSSHClient(&infra.Spec.SSH, true)
	var ips []string
	for _, host := range infra.Status.Hosts {
		if host.State != "Running" {
			v1.SetConditionError(configCondition, "VMStatus", fmt.Errorf("vm status is not running"))
			continue
		}
		ips = append(ips, host.IPs[0])
	}
	err := ssh.WaitSSHReady(client, 6, ips...)
	if err != nil {
		logger.Error("ping vms is error: %+v", err)
		return
	}
}

func (r *MultiPassVirtualMachine) CreateVM(infra *v1.VirtualMachine, host *v1.Host, index int) error {
	cfg := GetCloudInitYamlByRole(infra.Name, host.Role)
	debugFlag := ""
	if logger.IsDebugMode() {
		debugFlag = "-vvv"
	}
	if _, err := r.GetById(strings.GetID(infra.Name, host.Role, index)); err != nil {
		cmd := fmt.Sprintf("multipass launch --name %s --cpus %d --mem %dG --disk %dG --cloud-init %s %s %s ", strings.GetID(infra.Name, host.Role, index), host.Resources[v1.CPUKey], host.Resources[v1.MEMKey], host.Resources[v1.DISKKey], cfg, debugFlag, host.Image)
		logger.Info("executing... %s \n", cmd)
		return exec.Cmd("bash", "-c", cmd)
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
