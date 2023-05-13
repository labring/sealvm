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
	"github.com/labring/sealvm/pkg/utils/confirm"
	"github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/yaml"
	"github.com/modood/table"
	"path"
	"sort"
)

func NewValues() *values {
	return &values{}
}

type values struct{}

func (*values) List() {
	prints := make([]valuesPrint, 0)
	for key, p := range defaultValues {
		prints = append(prints, valuesPrint{
			Key:      key,
			Value:    p.Value,
			Describe: p.Describe,
		})
	}
	sort.Slice(prints, func(i, j int) bool {
		return prints[i].Key < prints[j].Key
	})
	table.OutputA(prints)
}

func (*values) Get(key string) string {
	if val, ok := defaultValues[key]; ok {
		return val.Value
	}
	return ""
}

func (*values) Default() {
	filePath := path.Join(defaultDir, "default.values")
	_ = file.CleanFiles(filePath)
}
func (*values) Set() error {
	key, err := confirm.Input("Please input a template values key", func(input string) error {
		return nil
	})
	if err != nil {
		return err
	}

	value, err := confirm.Input("Please input a template values value:", func(input string) error {
		return nil
	})
	if err != nil {
		return err
	}
	describe, err := confirm.Input("Please input a template values describe:", func(input string) error {
		return nil
	})
	if err != nil {
		return err
	}
	old := defaultValues[key]
	if describe == "" {
		describe = old.Describe
	}
	defaultValues[key] = Value{
		Describe: describe,
		Value:    value,
	}
	filePath := path.Join(defaultDir, "default.values")
	_ = yaml.MarshalYamlToFile(filePath, defaultValues)
	logger.Info("Set template values success")
	return nil
}

type valuesPrint struct {
	Key      string
	Value    string
	Describe string
	Sort     int `table:"-"`
}
