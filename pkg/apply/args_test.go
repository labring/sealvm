/*
Copyright 2023 cuisongliu@qq.com.

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

package apply

import (
	"reflect"
	"testing"
)

func TestParseMounts(t *testing.T) {
	type args struct {
		mounts []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				mounts: []string{
					"node@/data:/data",
					"master@/data1:/data1",
				},
			},
			wantErr: false,
			want: map[string]map[string]string{
				"node": {
					"/data": "/data",
				},
				"master": {
					"/data1": "/data1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMounts(tt.args.mounts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMounts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNodes(t *testing.T) {
	type args struct {
		nodes string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				nodes: "node1:1,node2:2",
			},
			wantErr: false,
			want: map[string]int{
				"node1": 1,
				"node2": 2,
			},
		},
		{
			name: "test-fales",
			args: args{
				nodes: "node1:ff,node2:ee",
			},
			wantErr: true,
		},
		{
			name: "test-fales",
			args: args{
				nodes: "node1:1,node2:22:33",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNodes(tt.args.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
