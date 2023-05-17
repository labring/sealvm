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

package template

import (
	"github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/yaml"
	"path"
	"runtime"
)

type ValueType string

type Value struct {
	Describe string `json:"describe"`
	Value    string `json:"value"`
}

type Values map[string]Value

func (v Values) convertMap() map[string]string {
	data := make(map[string]string, 0)
	for key, val := range v {
		data[key] = val.Value
	}
	return data
}

var defaultValues = Values{
	"HTTPProxy": Value{
		Describe: "HTTP Proxy address to be used",
		Value:    "192.168.64.1:7890",
	},
	"SocketProxy": Value{
		Describe: "Socket Proxy address to be used",
		Value:    "192.168.64.1:7890",
	},
	"NoProxy": Value{
		Describe: "List of addresses that should bypass the proxy",
		Value:    "",
	},
	"PublicKey": Value{
		Describe: "Path to the public key file",
		Value:    path.Join(file.GetHomeDir(), ".ssh", "id_rsa.pub"),
	},
	"PrivateKey": Value{
		Describe: "Path to the private key file",
		Value:    path.Join(file.GetHomeDir(), ".ssh", "id_rsa"),
	},
	"ARCH": Value{
		Describe: "The architecture of the current machine",
		Value:    runtime.GOARCH,
	},
}

func init() {
	filePath := path.Join(defaultDir, "default.values")
	if !file.IsExist(filePath) {
		_ = yaml.MarshalYamlToFile(filePath, defaultValues)
	} else {
		_ = yaml.UnmarshalYamlFromFile(filePath, &defaultValues)
	}
}
