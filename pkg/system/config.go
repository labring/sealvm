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
	"os"
	"strings"
)

type envSystemConfig struct{}

func Get(key string) (string, error) {
	return globalConfig.getValueOrDefault(key)
}

func Set(key, value string) error {
	return globalConfig.setValue(key, value)
}

var globalConfig *envSystemConfig

func init() {
	globalConfig = &envSystemConfig{}
}

type ConfigOption struct {
	Key           string
	Description   string
	DefaultValue  string
	OSEnv         string
	AllowedValues []string
}

var configOptions = []ConfigOption{
	{
		Key:          DefaultCPUKey,
		OSEnv:        "SEALVM_DEFAULT_CPU",
		Description:  "sealvm default cpu cores",
		DefaultValue: "2",
	},
	{
		Key:          DefaultMemKey,
		OSEnv:        "SEALVM_DEFAULT_MEM",
		Description:  "sealvm default memory size. unit is GB",
		DefaultValue: "4",
	},
	{
		Key:          DefaultDISKKey,
		OSEnv:        "SEALVM_DEFAULT_DISK",
		Description:  "sealvm default disk size. unit is GB",
		DefaultValue: "50",
	},
	{
		Key:          DefaultImageKey,
		OSEnv:        "SEALVM_DEFAULT_IMAGE",
		Description:  "sealvm default image local",
		DefaultValue: "",
	},
}

const (
	DefaultCPUKey   = "default_cpu"
	DefaultMemKey   = "default_mem"
	DefaultDISKKey  = "default_disk"
	DefaultImageKey = "default_image"
)

func (*envSystemConfig) getValueOrDefault(key string) (string, error) {
	for _, option := range configOptions {
		if option.Key == key {
			if option.OSEnv == "" {
				option.OSEnv = strings.ReplaceAll(strings.ToUpper("sealvm"+"_"+option.Key), "-", "_")
			}
			if value, ok := os.LookupEnv(option.OSEnv); ok {
				return value, nil
			}
			return option.DefaultValue, nil
		}
	}
	return "", fmt.Errorf("not found config key %s", key)
}

func (*envSystemConfig) setValue(key, value string) error {
	for _, option := range configOptions {
		if option.Key == key {
			if option.OSEnv == "" {
				return fmt.Errorf("not support set key %s, env not set", key)
			}
			if option.AllowedValues != nil {
				for _, allowedValue := range option.AllowedValues {
					if allowedValue == value {
						return os.Setenv(option.OSEnv, value)
					}
				}
				return fmt.Errorf("value %s is not allowed for key %s", value, key)
			}
		}
	}
	return nil
}

func ConfigOptions() []ConfigOption {
	return configOptions
}
