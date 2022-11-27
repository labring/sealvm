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
	"reflect"
	"testing"

	v1 "github.com/labring/sealvm/types/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDiffVirtualMachine(t *testing.T) {
	type args struct {
		old *v1.VirtualMachine
		new *v1.VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		wantAdd    []string
		wantDelete []string
	}{
		{
			name: "default",
			args: args{
				old: &v1.VirtualMachine{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
					Spec: v1.VirtualMachineSpec{
						Hosts: []v1.Host{
							{
								Role:  "master",
								Count: 1,
							},
							{
								Role:  "node",
								Count: 1,
							},
						},
					},
					Status: v1.VirtualMachineStatus{
						Hosts: []v1.VirtualMachineHostStatus{
							{
								Role:  "master",
								Index: 0,
							},
							{
								Role:  "node",
								Index: 0,
							},
							{
								Role:  "node",
								Index: 1,
							},
						},
					},
				},
				new: &v1.VirtualMachine{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
					Spec: v1.VirtualMachineSpec{
						Hosts: []v1.Host{
							{
								Role:  "master",
								Count: 3,
							},
						},
					},
				},
			},
			wantAdd:    []string{"default-master-1", "default-master-2"},
			wantDelete: []string{"default-node-0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdd, gotDelete := DiffVirtualMachine(tt.args.old, tt.args.new)
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("DiffVirtualMachine() gotAdd = %v, want %v", gotAdd, tt.wantAdd)
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("DiffVirtualMachine() gotDelete = %v, want %v", gotDelete, tt.wantDelete)
			}
		})
	}
}
