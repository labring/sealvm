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
	"errors"
	"fmt"
	"github.com/labring/sealos-vm/pkg/utils/exec"
	"github.com/labring/sealos-vm/pkg/utils/logger"
	v1 "github.com/labring/sealos-vm/types/api/v1"
	errors2 "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
)

func (r *MultiPassVirtualMachine) reconcile() {
	logger.Info("Start to reconcile a new infra:", r.Desired.Name)

	pipelines := []func(infra *v1.VirtualMachine){
		r.InitStatus,
		r.ApplyConfig,
		r.SyncVMs,
		r.PingVms,
		r.FinalStatus,
	}
	if !r.Desired.DeletionTimestamp.IsZero() {
		pipelines = []func(infra *v1.VirtualMachine){
			r.InitStatus,
			r.DeleteVMs,
		}
	}
	for _, fn := range pipelines {
		fn(r.Desired)
	}
	if r.Desired.Status.Phase != v1.PhaseFailed {
		r.Desired.Status.Phase = v1.PhaseSuccess
	}
	logger.Info("succeeded in reconcile, enjoy it!")
}

func (r *MultiPassVirtualMachine) DeleteVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec DeleteVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "DeleteVMs",
		Status:            v12.ConditionTrue,
		Reason:            "Delete VMs",
		Message:           "config has been delete multipass instances",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)

	eg, _ := errgroup.WithContext(context.Background())

	for _, host := range infra.Status.Hosts {
		dHost := host
		eg.Go(func() error {
			return r.DeleteVM(infra, &dHost)
		})
	}
	if err := eg.Wait(); err != nil {
		v1.SetConditionError(configCondition, "DeleteVMsError", err)
		return
	}

}

func (r *MultiPassVirtualMachine) DeleteVM(infra *v1.VirtualMachine, host *v1.VirtualMachineHostStatus) error {
	cmd := fmt.Sprintf("multipass delete -p   %s ", host.ID)
	return exec.Cmd("bash", "-c", cmd)
}

func (r *MultiPassVirtualMachine) Get(name, role string, index int) (string, error) {
	cmd := fmt.Sprintf("multipass info %s-%s-%d --format=json", name, role, index)
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return "", errors.New("not found instance")
	}
	return out, nil
}

func (r *MultiPassVirtualMachine) GetById(name string) (string, error) {
	cmd := fmt.Sprintf("multipass info %s --format=json", name)
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return "", errors.New("not found instance")
	}
	return out, nil
}

func (r *MultiPassVirtualMachine) Inspect(name, role string, index int) (*v1.VirtualMachineHostStatus, error) {
	info, err := r.Get(name, role, index)
	if err != nil {
		return nil, err
	}
	var outStruct map[string]interface{}
	err = json.Unmarshal([]byte(info), &outStruct)
	if err != nil {
		return nil, errors2.Wrap(err, "decode out json from multipass info failed")
	}
	hostStatus := &v1.VirtualMachineHostStatus{
		State:     "",
		Role:      role,
		ID:        fmt.Sprintf("%s-%s-%d", name, role, index),
		IPs:       nil,
		ImageID:   "",
		ImageName: "",
		Capacity:  nil,
		Used:      nil,
		Mounts:    nil,
	}
	capacity := map[string]int{
		"cpu":    2,
		"memory": 4,
		"disk":   50,
	}
	hostStatus.Capacity = capacity
	hostStatus.State, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "state")
	hostStatus.ImageID, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "image_hash")
	hostStatus.ImageName, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "release")
	hostStatus.IPs, _, _ = unstructured.NestedStringSlice(outStruct, "info", hostStatus.ID, "ipv4")
	return hostStatus, nil
}
