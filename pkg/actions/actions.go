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

package actions

import (
	"fmt"
	"github.com/labring/sealvm/pkg/actions/runtime"
	"github.com/labring/sealvm/pkg/utils/confirm"
	"github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	yutil "github.com/labring/sealvm/pkg/utils/yaml"
	v1 "github.com/labring/sealvm/types/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Do(name, p string) error {
	if !file.IsExist(p) {
		return fmt.Errorf("file %s not exist", p)
	}
	data, err := file.ReadAll(p)
	if err != nil {
		return err
	}
	logger.Info("action yamls: %s", string(data))
	if yes, err := confirm.Confirm("Are you sure to run this command?", "you have canceled to exec action !"); err != nil {
		return err
	} else {
		if !yes {
			return fmt.Errorf("you have canceled to exec action ")
		}
	}
	yamls := yutil.ToYalms(string(data))
	actions := make([]v1.Action, 0)
	for _, y := range yamls {
		action := v1.Action{}
		err = yaml.Unmarshal([]byte(y), &action)
		if err != nil {
			logger.Warn("unmarshal action error: %v", err)
			continue
		}
		actions = append(actions, action)
	}
	r, err := runtime.NewAction(name)
	if err != nil {
		return err
	}
	errArr := make([]error, 0)
	newActions := make([]any, 0)
	for index, action := range actions {
		if err = r.Apply(&action); err != nil {
			logger.Error("apply action %d error: %v", index, err)
			errArr = append(errArr, err)
		}
		newActions = append(newActions, action)
	}

	outActionfile, _ := yutil.MarshalYamlConfigs(newActions...)

	logger.Info("outActionfile: %s", string(outActionfile))

	if len(errArr) > 0 {
		logger.Error("apply actions error: %v", errors.NewAggregate(errArr).Error())
		return nil
	}
	return nil
}

func PrintDefault() error {
	actions := make([]any, 0)
	action := v1.Action{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Action",
			APIVersion: v1.GroupVersion.String(),
		},
		Spec: v1.ActionSpec{
			Ons: []v1.ActionOn{
				{
					Role:    "node",
					Indexes: nil,
				},
				{
					Role:    "master",
					Indexes: []int32{0, 1},
				},
			},
			Data: []v1.ActionData{
				{
					ActionMount: &v1.SourceAndTarget{
						Source: "/source",
						Target: "/target",
					},
				},
				{
					ActionUmount: "/target",
				},
				{
					ActionExec: "ls -l /",
				},
				{
					ActionCopy: &v1.SourceAndTarget{
						Source: "/source",
						Target: "/target",
					},
				},
				{
					ActionCopyContent: &v1.ContentAndTarget{
						Content: "write code\ndfff",
						Target:  "/target",
					},
				},
			},
		},
	}
	actions = append(actions, action)
	outActionfile, _ := yutil.MarshalYamlConfigs(actions...)
	println(string(outActionfile))
	return nil
}
