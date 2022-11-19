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

package strings

import (
	"fmt"
	"strings"
)

func GetID(name, role string, index int) string {
	return fmt.Sprintf("%s-%s-%d", name, role, index)
}

func GetHostV1FromAliasName(aliasName string) (name, role, index string) {
	all := strings.Split(aliasName, "-")
	if len(all) == 3 {
		return all[0], all[1], all[2]
	}
	return
}
