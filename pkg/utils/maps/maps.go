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

package maps

import "strings"

func StringToMap(data string, spilt string) map[string]string {
	list := strings.Split(data, spilt)
	return ListToMap(list)
}

func ListToMap(data []string) map[string]string {
	m := make(map[string]string)
	for _, l := range data {
		if l != "" {
			kv := strings.Split(l, ":")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
	}
	return m
}
