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
	"bytes"
	"errors"
	"fmt"
	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/confirm"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	template2 "github.com/labring/sealvm/pkg/utils/template"
	"github.com/modood/table"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const templateSuffix = ".tmpl"

var defaultDir = path.Join(configs.DefaultRootfsDir(), "etc")

func NewTpl() *template {
	return &template{}
}

type template struct{}

func (*template) Get(role string) (string, error) {
	filePath := path.Join(defaultDir, fmt.Sprintf("%s%s", role, templateSuffix))
	if !fileutil.IsExist(filePath) {
		return "", fmt.Errorf("the role %s is not exist", role)
	}
	data, _ := fileutil.ReadAll(filePath)
	return string(data), nil
}

func (*template) Default() string {
	return `write_files:
- content: |
    export https_proxy=http://{{ .HTTPProxy }} http_proxy=http://{{ .HTTPProxy }} all_proxy=socks5://{{ .SocketProxy }}
    export no_proxy=localhost,127.0.0.1,localaddress,.localdomain.com,apiserver.cluster.local,{{ .NoProxy }}
    echo -e "已开启代理"
  path: /usr/bin/proxy_on
  permissions: '0755'
- content: |
    unset http_proxy
    unset https_proxy
    unset ftp_proxy
    unset rsync_proxy
    echo -e "已关闭代理"
  path: /usr/bin/proxy_off
  permissions: '0755'
runcmd:
  - echo "{{ $result := readFile .PublicKey }}{{ b64enc $result }}" | base64 -d >> /root/.ssh/authorized_keys
  - echo "{{ $result := readFile .PrivateKey }}{{ b64enc $result }}" | base64 -d >> /root/.ssh/id_rsa
  - chmod 600 /root/.ssh/id_rsa
  - sed -i "/update_etc_hosts/c \ - ['update_etc_hosts', 'once-per-instance']" /etc/cloud/cloud.cfg
  - touch /var/lib/cloud/instance/sem/config_update_etc_hosts`
}

func (*template) Set() error {
	role, err := confirm.Input("Please input a role (more than 3 characters)", func(input string) error {
		if len(input) < 3 {
			return errors.New("the role name must have more than 3 characters")
		}
		return nil
	})
	if err != nil {
		return err
	}

	p, err := confirm.Input(fmt.Sprintf("Please input an absolute and existing file path for the %s template:", role), func(input string) error {
		if !filepath.IsAbs(input) {
			return errors.New("the input must be an absolute path")
		}
		if !fileutil.IsExist(input) {
			return errors.New("the input file path must exist")
		}
		if !fileutil.IsFile(input) {
			return errors.New("the input file path must is file")
		}
		inputData, _ := fileutil.ReadAll(input)
		// 首先，使用 text/template 来解析模板
		tmpl, ok, _ := template2.TryParse(string(inputData))
		if !ok {
			return fmt.Errorf("failed to parse the template")
		}
		// 然后，生成 cloud-init 文件
		var cloudInit bytes.Buffer
		err = tmpl.Execute(&cloudInit, defaultValues.convertMap())
		if err != nil {
			return fmt.Errorf("failed to generate the cloud-init file: %+v", err)
		}
		var yamlData interface{}
		err = yaml.Unmarshal(cloudInit.Bytes(), &yamlData)
		if err != nil {
			return fmt.Errorf("the cloud-init file is not in valid YAML format: %+v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	filePath := path.Join(defaultDir, fmt.Sprintf("%s%s", role, templateSuffix))
	_ = fileutil.Copy(p, filePath)
	logger.Info("Sync template role=%s config success.", role)
	return nil
}

func (*template) Reset() {
	paths, err := getTemplateFiles(defaultDir)
	if err != nil {
		logger.Error("get template files error: %+v", err)
		return
	}
	logger.Info("reset template files: %+v", paths)
	_ = fileutil.CleanFiles(paths...)
}

func (*template) List() {
	paths, err := getTemplateFiles(defaultDir)
	if err != nil {
		logger.Error("get template files error: %+v", err)
		return
	}
	prints := make([]templatePrint, 0)
	for _, p := range paths {
		fName := fileutil.Filename(p)
		if strings.HasSuffix(fName, templateSuffix) {
			role := strings.ReplaceAll(fName, templateSuffix, "")
			prints = append(prints, templatePrint{
				Role:    role,
				AbsPath: p,
			})
		}

	}
	if len(prints) == 0 {
		logger.Info("Print template list lens is 0.")
		return
	}
	table.OutputA(prints)
}

func getTemplateFiles(p string) (paths []string, err error) {
	_, err = os.Stat(p)
	if err != nil {
		return
	}

	err = filepath.Walk(p, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if p == path.Join(defaultDir, info.Name()) && strings.HasSuffix(info.Name(), templateSuffix) {
			paths = append(paths, p)
		}
		return err
	})
	return paths, err
}

type templatePrint struct {
	Role    string
	AbsPath string
}

func init() {
	_ = fileutil.MkDirs(defaultDir)
}
