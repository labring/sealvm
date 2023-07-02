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

package vm

import (
	"fmt"
	"github.com/labring/sealvm/pkg/utils/logger"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type File struct {
	Content     string `yaml:"content"`
	Path        string `yaml:"path"`
	Permissions string `yaml:"permissions"`
}

type Config struct {
	WriteFiles []File   `yaml:"write_files"`
	RunCmd     []string `yaml:"runcmd"`
}

func cloudInit(fileName string) *Config {
	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error("yamlFile.Get err   #%v ", err)
		return nil
	}
	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logger.Error("Unmarshal cloud init config: %v", err)
		return nil
	}
	return &config
}

func (c *Config) toScript() string {
	sb := strings.Builder{}
	sb.WriteString("#!/bin/bash\n")
	for _, f := range c.WriteFiles {
		sb.WriteString("cat <<EOF > ")
		sb.WriteString(f.Path)
		sb.WriteString("\n")
		sb.WriteString(f.Content)
		sb.WriteString("\nEOF\n")
		sb.WriteString(fmt.Sprintf("chmod %s %s\n", f.Permissions, f.Path))
	}
	sb.WriteString("apt-get update -y\n")
	sb.WriteString("apt-get install -y openssh-client\n")
	sb.WriteString("apt-get install -y openssh-server\n")
	sb.WriteString("mkdir -p ~/.ssh\n")
	for _, cmd := range c.RunCmd {
		if strings.Contains(cmd, "/etc/cloud") {
			continue
		}
		if strings.Contains(cmd, "/var/lib/cloud") {
			continue
		}
		sb.WriteString(cmd)
		sb.WriteString("\n")
	}
	return sb.String()
}
