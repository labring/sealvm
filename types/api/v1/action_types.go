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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ActionSpec defines the desired state of Action
type ActionSpec struct {
	//master master-0
	//node node-0
	//registry registry-0
	//dev dev-0
	//host
	On   []string `json:"on,omitempty"`
	Type string   `json:"type,omitempty"`
	Cmd  []string `json:"cmd,omitempty"`
}

type ActionPhase string

const (
	ActionPhaseFailed    ActionPhase = "Failed"
	ActionPhaseComplete  ActionPhase = "Complete"
	ActionPhaseInProcess ActionPhase = "InProcess"
)

// ActionStatus defines the observed state of Action
type ActionStatus struct {
	Phase      ActionPhase `json:"phase,omitempty"`
	Conditions []Condition `json:"conditions,omitempty" `
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Action is the Schema for the action API
type Action struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActionSpec   `json:"spec,omitempty"`
	Status ActionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActionList contains a list of Action
type ActionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Action `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Action{}, &ActionList{})
}
