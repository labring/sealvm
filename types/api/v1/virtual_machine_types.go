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
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/labring/sealvm/pkg/utils/iputils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MultipassType = "Multipass"

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	Hosts []Host `json:"hosts,omitempty"`
	SSH   SSH    `json:"ssh"`
	Type  string `json:"provider,omitempty"`
}

type SSH struct {
	PublicFile string `json:"publicFile,omitempty"`
	PkFile     string `json:"pkFile,omitempty"`
	PkPasswd   string `json:"pkPasswd,omitempty"`
}

type Host struct {
	Role   string            `json:"roles,omitempty"`
	Count  int               `json:"count,omitempty"`
	Mounts map[string]string `json:"mounts,omitempty"`
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

func (c *Host) GetRoles() string {
	return c.Role
}

func (c *VirtualMachine) GetIPSByRole(role string) []string {
	var hosts []string
	for _, host := range c.Status.Hosts {
		if role == host.Role {
			hosts = append(hosts, host.IPs...)
		}
	}
	return hosts
}

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

func (c *VirtualMachine) GetSSH() SSH {
	return c.Spec.SSH
}

func (c *VirtualMachine) GetALLIPList() []string {
	ips := make([]string, 0)
	for _, r := range c.GetRoles() {
		ips = append(ips, c.GetIPSByRole(r)...)
	}
	return ips
}

func (c *VirtualMachine) GetMaster0IP() string {
	if len(c.Spec.Hosts) == 0 {
		return ""
	}
	if len(c.Status.Hosts[0].IPs) == 0 {
		return ""
	}
	return iputils.GetHostIP(c.Status.Hosts[0].IPs[0])
}

func (c *VirtualMachine) GetMaster0IPAPIServer() string {
	master0 := c.GetMaster0IP()
	return fmt.Sprintf("https://%s:6443", master0)
}

func (c *VirtualMachine) GetRolesByIP(ip string) string {
	for _, host := range c.Status.Hosts {
		if In(ip, host.IPs) {
			return host.Role
		}
	}
	return ""
}

func In(key string, slice []string) bool {
	for _, s := range slice {
		if key == s {
			return true
		}
	}
	return false
}
