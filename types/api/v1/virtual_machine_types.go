// Copyright © 2022 cuisongliu@qq.com Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MultipassType = "multipass"
const OrbType = "orb"

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	Hosts []Host `json:"hosts,omitempty"`
	SSH   SSH    `json:"ssh"`
}

type SSH struct {
	PublicFile string `json:"publicFile,omitempty"`
	PkFile     string `json:"pkFile,omitempty"`
	PkPasswd   string `json:"pkPasswd,omitempty"`
}

type Host struct {
	Role  string `json:"roles,omitempty"`
	Count int    `json:"count,omitempty"`
	// key values resources.
	// cpu: 2
	// memory: 4
	// other resources like GPU
	Resources map[string]string `json:"resources,omitempty"`
	// ecs.t5-lc1m2.large
	Image string `json:"image,omitempty"`
}

type Phase string

const (
	PhaseFailed    Phase = "Failed"
	PhaseSuccess   Phase = "Success"
	PhaseInProcess Phase = "InProcess"
)

// VirtualMachineStatus defines the observed state of VirtualMachine
type VirtualMachineStatus struct {
	Phase      Phase                      `json:"phase,omitempty"`
	Hosts      []VirtualMachineHostStatus `json:"hosts"`
	Conditions []Condition                `json:"conditions,omitempty" `
}

type Condition struct {
	Type              string             `json:"type"`
	Status            v1.ConditionStatus `json:"status"`
	LastHeartbeatTime metav1.Time        `json:"lastHeartbeatTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

func (s *VirtualMachineHostStatus) IsRunning() bool {
	if s.State == "Running" {
		return true
	}
	if s.State == "running" {
		return true
	}
	return false
}

type VirtualMachineHostStatus struct {
	State string `json:"state"`
	Role  string `json:"roles"`
	ID    string `json:"ID,omitempty"`

	//当前主机的所有IP，可能包括公开或者私有的IP
	IPs       []string          `json:"IPs,omitempty"`
	ImageID   string            `json:"imageID,omitempty"`
	ImageName string            `json:"imageName,omitempty"`
	Capacity  map[string]string `json:"capacity"`
	Used      map[string]string `json:"used"`
	Mounts    map[string]string `json:"mounts,omitempty"`
	Index     int               `json:"index,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualMachine is the Schema for the infra API
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualMachineList contains a list of VirtualMachine
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}

// Language: go

var (
	CPUKey  = "cpu"
	MEMKey  = "memory"
	DISKKey = "disk"
)

func (c *VirtualMachine) GetRoles() []string {
	roles := sets.NewString()
	for _, host := range c.Spec.Hosts {
		roles = roles.Insert(host.Role)
	}
	return roles.List()
}

func (c *VirtualMachine) GetHostByRole(role string) *Host {
	for _, host := range c.Spec.Hosts {
		if role == host.Role {
			return &host
		}
	}
	return nil
}

func (c *VirtualMachine) GetHostStatusByRoleIndex(role string, index int) *VirtualMachineHostStatus {
	for _, host := range c.Status.Hosts {
		if role == host.Role && index == host.Index {
			return &host
		}
	}
	return nil
}

func (c *VirtualMachine) GetHostStatusByName(name string) *VirtualMachineHostStatus {
	for _, host := range c.Status.Hosts {
		if host.ID == name {
			return &host
		}
	}
	return nil
}
