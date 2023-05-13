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

package template

import (
	"bytes"
	"errors"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	template2 "github.com/labring/sealvm/pkg/utils/template"
)

func EtcHostsTplExecuteToFile(role, outFile string) error {
	tpl := template{}
	data, err := tpl.Get(role)
	if err != nil {
		return err
	}
	tp, ok, _ := template2.TryParse(data)
	if !ok {
		return errors.New("template parse error")
	}

	out := bytes.NewBuffer(nil)
	err = tp.Execute(out, defaultValues.convertMap())
	if err != nil {
		return err
	}
	logger.Debug("execute template %s", out.String())
	return fileutil.WriteFile(outFile, out.Bytes())
}
