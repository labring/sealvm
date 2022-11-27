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

package install

import (
	"path"
	"runtime"

	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/exec"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/progress"
)

type multipass struct{}

func (i *multipass) Install() error {
	macURL := "https://github.com/canonical/multipass/releases/download/v1.8.1/multipass-1.8.1+mac-Darwin.pkg"
	newMacURL := "https://github.com/canonical/multipass/releases/download/v1.8.1/multipass-1.8.1+mac-Darwin.pkg"
	winURL := "https://github.com/canonical/multipass/releases/download/v1.8.0/multipass-1.8.0+win-win64.exe"
	dirName := path.Join(configs.DefaultRootfsDir(), "multipass")
	_ = fileutil.MkDirs(dirName)
	if runtime.GOOS == "darwin" {
		if "arm64" == runtime.GOARCH {
			macURL = newMacURL
		}
		if !AutoDownload {
			logger.Info("please download multipass from %s", macURL)
			return nil
		}
		fileName := path.Join(dirName, "multipass.pkg")
		if !fileutil.IsExist(fileName) {
			err := progress.Download(macURL, fileName)
			if err != nil {
				return err
			}
			logger.Info("your multipass pkg download success")
		}
		return exec.Cmd("open", fileName)
	}
	if runtime.GOOS == "windows" {
		if !AutoDownload {
			logger.Info("please download multipass from %s", winURL)
			return nil
		}
		fileName := path.Join(dirName, "multipass.exe")
		if !fileutil.IsExist(fileName) {
			err := progress.Download(winURL, fileName)
			if err != nil {
				return err
			}
			logger.Info("your multipass exe download success")
		}
		return exec.Cmd("open", fileName)
	}
	return nil
}

func (i *multipass) IsInstall() bool {
	_, ok := exec.CheckCmdIsExist("multipass")
	return ok
}
