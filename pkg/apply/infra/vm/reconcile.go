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

package vm

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/labring/sealvm/pkg/apply/runtime"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	"golang.org/x/sync/errgroup"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *VirtualMachine) Reconcile(diff runtime.Diff) {
	logger.Info("Start to reconcile a new infra:", r.Desired.Name)
	r.DiffFunc = diff
	pipelines := []func(infra *v1.VirtualMachine){
		r.InitStatus,
		r.ApplyConfig,
		r.ApplyVMs,
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
func (r *VirtualMachine) ApplyVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec ApplyVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "ApplyVMs",
		Status:            v12.ConditionTrue,
		Reason:            "VM apply",
		Message:           "apply multipass success",
		LastHeartbeatTime: metav1.Now(),
	}
	defer r.saveCondition(infra, configCondition)
	addHostNames, deleteHostNames := r.DiffFunc(r.Current, r.Desired)

	eg, _ := errgroup.WithContext(context.Background())

	for _, host := range addHostNames {
		h := host
		eg.Go(func() error {
			_, role, index := strings.GetHostV1FromAliasName(h)
			hostObj := infra.GetHostByRole(role)
			if hostObj != nil {
				indexInt, err := strconv.Atoi(index)
				if err == nil {
					time.Sleep(time.Duration(indexInt) * time.Millisecond * 100)
					return r.CreateVM(infra, hostObj, indexInt)
				}
				return err
			}
			return fmt.Errorf("not found host from role: %s", role)
		})
	}

	for _, host := range deleteHostNames {
		h := host
		eg.Go(func() error {
			_, role, index := strings.GetHostV1FromAliasName(h)
			indexInt, err := strconv.Atoi(index)
			if err == nil {
				hostStatus := r.Current.GetHostStatusByRoleIndex(role, indexInt)
				if hostStatus != nil {
					return r.DeleteVM(r.Current, hostStatus)
				}
				return fmt.Errorf("not found host status from role: %s, index: %d", role, indexInt)
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		v1.SetConditionError(configCondition, "ApplyVMsError", err)
		return
	}
}

func (r *VirtualMachine) DeleteVMs(infra *v1.VirtualMachine) {
	logger.Info("Start to exec DeleteVMs:", r.Desired.Name)
	var configCondition = &v1.Condition{
		Type:              "DeleteVMs",
		Status:            v12.ConditionTrue,
		Reason:            "Delete VMs",
		Message:           "config has been delete local vm instances",
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
