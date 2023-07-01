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

package system

import (
	"fmt"
	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/yaml"
	"path"
)

type envSystemConfig struct{}

func Get(key string) (string, error) {
	return globalConfig.getValue(key)
}

func Set(key, value string) error {
	return globalConfig.setValue(key, value)
}

func List() {
	list := ConfigOptions()
	for _, v := range list {
		data, _ := Get(v.Key)
		println(fmt.Sprintf("%s=%s", v.Key, data))
	}
}

func toDefaultYaml() {
	list := ConfigOptions()
	data := make(map[string]string)
	for _, v := range list {
		data[v.Key] = v.DefaultValue
	}
	_ = yaml.MarshalYamlToFile(path.Join(defaultDir, configFile), data)
}

func fromDefaultYaml() {
	data := make(map[string]string)
	_ = yaml.UnmarshalYamlFromFile(path.Join(defaultDir, configFile), &data)

	for i, v := range ConfigOptions() {
		if val, ok := data[v.Key]; ok {
			ConfigOptions()[i].DefaultValue = val
		}
	}
}

var globalConfig *envSystemConfig

func init() {
	globalConfig = &envSystemConfig{}
}

type ConfigOption struct {
	Key          string
	Description  string
	DefaultValue string
}

var configOptions = []ConfigOption{
	{
		Key:          DefaultCPUKey,
		Description:  "sealvm default cpu cores",
		DefaultValue: "2",
	},
	{
		Key:          DefaultMemKey,
		Description:  "sealvm default memory size. unit is GB",
		DefaultValue: "4",
	},
	{
		Key:          DefaultDISKKey,
		Description:  "sealvm default disk size. unit is GB",
		DefaultValue: "50",
	},
	{
		Key:          DefaultImageKey,
		Description:  "sealvm default image local",
		DefaultValue: "",
	},
	{
		Key:          DefaultProvider,
		Description:  "sealvm default provider",
		DefaultValue: "multipass",
	},
}

const (
	DefaultCPUKey   = "default_cpu"
	DefaultMemKey   = "default_mem"
	DefaultDISKKey  = "default_disk"
	DefaultImageKey = "default_image"
	DefaultProvider = "default_provider"
)

var defaultDir = path.Join(configs.DefaultRootfsDir(), "etc")

const configFile = "default.cfg"

func init() {
	filePath := path.Join(defaultDir, configFile)
	if !file.IsExist(filePath) {
		toDefaultYaml()
	} else {
		fromDefaultYaml()
	}
}

func (*envSystemConfig) getValue(key string) (string, error) {
	for _, option := range configOptions {
		if option.Key == key {
			return option.DefaultValue, nil
		}
	}
	return "", fmt.Errorf("not found config key %s", key)
}

func (*envSystemConfig) setValue(key, value string) error {
	for i, option := range configOptions {
		if option.Key == key {
			configOptions[i].DefaultValue = value
		}
	}
	toDefaultYaml()
	return nil
}

func ConfigOptions() []ConfigOption {
	return configOptions
}
