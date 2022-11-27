/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

package cmd

import (
	"fmt"

	"github.com/labring/sealvm/pkg/install"
	"github.com/labring/sealvm/pkg/utils/logger"
	v1 "github.com/labring/sealvm/types/api/v1"

	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var vmType string
	var installer install.Interface
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install vm tools",
		Run: func(cmd *cobra.Command, args []string) {
			if installer.IsInstall() {
				logger.Info("kubernetes is installed")
				return
			}
			err := installer.Install()
			if err != nil {
				logger.Error(err)
			}
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			installer = install.NewInstaller(vmType)
			if installer == nil {
				return fmt.Errorf("vm type %s not support", vmType)
			}
			return nil
		},
	}
	installCmd.Flags().StringVarP(&vmType, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	installCmd.Flags().BoolVarP(&install.AutoDownload, "download", "d", true, "auto download vm tools online")
	return installCmd
}

func init() {
	rootCmd.AddCommand(newInstallCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkInstall(vmType string) error {
	installer := install.NewInstaller(vmType)
	if installer == nil {
		return fmt.Errorf("vm type %s not support", vmType)
	}
	if !installer.IsInstall() {
		return fmt.Errorf("vm tools %s is not installed, please use `sealvm install` retry", vmType)
	}
	return nil
}
