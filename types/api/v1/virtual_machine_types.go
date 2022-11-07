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
	"github.com/labring/sealos-vm/pkg/utils/iputils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MultipassType = "Multipass"

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	Hosts   []Host `json:"hosts,omitempty"`
	SSH     SSH    `json:"ssh"`
	Type    string `json:"provider,omitempty"`
	Proxy   string
	NoProxy string
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
	Resources map[string]int `json:"resources,omitempty"`
	// ecs.t5-lc1m2.large
	Image string `json:"image,omitempty"`
}

//Name:           sealos-dev
//State:          Running
//IPv4:           192.168.64.23
//Release:        Ubuntu 22.04.1 LTS
//Image hash:     4d8d5b95082e (Ubuntu 22.04 LTS)
//Load:           0.00 0.00 0.00
//Disk usage:     3.9G out of 96.7G
//Memory usage:   247.5M out of 3.8G
//Mounts:         /Users/cuisongliu/Workspaces/go/src/github.com => /root/go/src/github.com
//UID map: 0:default
//GID map: 0:default

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
	State     string            `json:"state"`
	Role      string            `json:"roles"`
	ID        string            `json:"ID,omitempty"`
	IPs       []string          `json:"IPs,omitempty"`
	ImageID   string            `json:"imageID,omitempty"`
	ImageName string            `json:"imageName,omitempty"`
	Capacity  map[string]int    `json:"capacity"`
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
	MASTER   = "master"
	NODE     = "node"
	REGISTRY = "registry"
	DEV      = "dev"

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

func (c *VirtualMachine) GetSSH() SSH {
	return c.Spec.SSH
}

func (c *VirtualMachine) GetMasterIPList() []string {
	return iputils.GetHostIPs(c.GetIPSByRole(MASTER))
}

func (c *VirtualMachine) GetNodeIPList() []string {
	return iputils.GetHostIPs(c.GetIPSByRole(NODE))
}

func (c *VirtualMachine) GetRegistryIPList() []string {
	return iputils.GetHostIPs(c.GetIPSByRole(REGISTRY))
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
