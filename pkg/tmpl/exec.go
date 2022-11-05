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

package tmpl

import (
	"bytes"
	fileutil "github.com/cuisongliu/sealos-dev/pkg/utils/file"
	"github.com/cuisongliu/sealos-dev/pkg/utils/logger"
	"html/template"
	"os"
	"runtime"
)

func (tpl Tpl) Execute(data map[string]interface{}) (string, error) {
	t := template.New("tpl_mulipass")
	t, err := t.Parse(string(tpl))
	if err != nil {
		return "", err
	}
	t = template.Must(t, err)
	out := bytes.NewBuffer(nil)
	err = t.Execute(out, data)
	if err != nil {
		logger.Error("tpl_mulipass template can not excute %s", err)
		return "", err
	}
	return out.String(), nil
}

func ExecuteNodesToFile(proxy, noProxy, privateKeyFile, publicKeyFile, file string) error {
	tpl := NodesTpl
	privateKeyFileRead, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}
	publicKeyFileRead, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}
	if proxy == "" {
		proxy = "192.168.64.1:7890"
	}
	if noProxy == "" {
		noProxy = "192.168.0.0/16"
	}
	data := map[string]interface{}{
		"Proxy":      proxy,
		"NoProxy":    noProxy,
		"PrivateKey": string(privateKeyFileRead),
		"PublicKey":  string(publicKeyFileRead),
	}
	outString, err := tpl.Execute(data)
	if err != nil {
		return err
	}
	return fileutil.WriteFile(file, []byte(outString))
}

func ExecuteGolangToFile(proxy, noProxy, privateKeyFile, publicKeyFile, file string) error {
	tpl := GolangTpl
	privateKeyFileRead, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}
	publicKeyFileRead, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}
	if proxy == "" {
		proxy = "192.168.64.1:7890"
	}
	if noProxy == "" {
		noProxy = "192.168.0.0/16"
	}
	data := map[string]interface{}{
		"Proxy":      proxy,
		"NoProxy":    noProxy,
		"PrivateKey": string(privateKeyFileRead),
		"PublicKey":  string(publicKeyFileRead),
		"ARCH":       runtime.GOARCH,
	}
	outString, err := tpl.Execute(data)
	if err != nil {
		return err
	}
	return fileutil.WriteFile(file, []byte(outString))
}
