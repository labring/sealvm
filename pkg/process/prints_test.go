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

package process

import (
	v1 "github.com/labring/sealvm/types/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_printVMs(t *testing.T) {
	type args struct {
		vm *v1.VirtualMachine
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				vm: &v1.VirtualMachine{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       v1.VirtualMachineSpec{},
					Status: v1.VirtualMachineStatus{
						Hosts: []v1.VirtualMachineHostStatus{
							{
								State:     "Running",
								ID:        "xxxx",
								Role:      "node",
								IPs:       []string{"127.0.0.1"},
								ImageName: "Ubuntu 11.01",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := printVMs(tt.args.vm); (err != nil) != tt.wantErr {
				t.Errorf("printVMs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
