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

package infra

import (
	"fmt"
	"github.com/labring/sealvm/pkg/apply/runtime"
	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	"github.com/labring/sealvm/pkg/utils/yaml"
	v1 "github.com/labring/sealvm/types/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
)

type Interface interface {
	Init()
	Reconcile(diff runtime.Diff)
	DesiredVM() *v1.VirtualMachine
	CurrentVM() *v1.VirtualMachine
}

type driver struct {
	Infra Interface
}

func (c *driver) Apply() error {
	if c.Infra.DesiredVM().CreationTimestamp.IsZero() {
		c.Infra.Init()
		c.Infra.DesiredVM().CreationTimestamp = metav1.Now()
	} else {
		c.Infra.Reconcile(DiffVirtualMachine)
	}
	return c.updateCRStatus()

}

func DiffVirtualMachine(old, new *v1.VirtualMachine) (add, delete []string) {
	oldSpec := sets.NewString()
	for _, h := range old.Spec.Hosts {
		for i := 0; i < h.Count; i++ {
			oldSpec.Insert(strings.GetID(old.Name, h.Role, i))
		}
	}
	newSpec := sets.NewString()
	for _, h := range new.Spec.Hosts {
		for i := 0; i < h.Count; i++ {
			newSpec.Insert(strings.GetID(new.Name, h.Role, i))
		}
	}
	addSets := newSpec.Difference(oldSpec)
	deleteSets := oldSpec.Difference(newSpec)
	return addSets.List(), deleteSets.List()
}

func (c *driver) getWriteBackObjects() []interface{} {
	obj := []interface{}{c.Infra.DesiredVM()}
	//if configs := c.ClusterFile.GetConfigs(); len(configs) > 0 {
	//	for i := range configs {
	//		obj = append(obj, configs[i])
	//	}
	//}
	return obj
}

// todo: atomic updating status after each installation for better reconcile?
// todo: set up signal handler
func (c *driver) updateCRStatus() error {
	if !c.Infra.DesiredVM().DeletionTimestamp.IsZero() {
		t := metav1.Now()
		cfPath := configs.VirtualMachineFilePath(c.Infra.DesiredVM().Name)
		target := fmt.Sprintf("%s.%d", cfPath, t.Unix())
		logger.Debug("write reset vm file to local: %s", target)
		if err := yaml.MarshalYamlToFile(cfPath, c.getWriteBackObjects()...); err != nil {
			logger.Error("failed to store vm file: %v", err)
		}
		_ = os.Rename(cfPath, target)
		return nil
	}
	infraPath := configs.VirtualMachineFilePath(c.Infra.DesiredVM().Name)
	logger.Debug("write cluster file to local storage: %s", infraPath)
	return yaml.MarshalYamlToFile(infraPath, c.getWriteBackObjects()...)
}
