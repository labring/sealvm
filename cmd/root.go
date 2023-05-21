/*
Copyright © 2022 cuisongliu@qq.com

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
	"github.com/labring/sealvm/pkg/system"
	"github.com/labring/sealvm/pkg/template"
	"os"
	"path"
	"runtime"

	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	debug          bool
	name           string
	clusterRootDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sealvm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onBootOnDie)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logger")
	rootCmd.PersistentFlags().StringVar(&clusterRootDir, "cluster-root", path.Join(file.GetHomeDir(), ".sealvm"), "cluster root directory")

	groups := templates.CommandGroups{
		{
			Message: "VM Management Commands:",
			Commands: []*cobra.Command{
				//newApplyCmd(),
				newRunCmd(),
				newResetCmd(),
				newInspectCmd(),
				newListCmd(),
			},
		},
		{
			Message: "Remote Operation Commands:",
			Commands: []*cobra.Command{
				newActionCmd(),
			},
		},
		{
			Message: "System Management Commands:",
			Commands: []*cobra.Command{
				newInstallCmd(),
				system.NewConfigCmd(),
				template.NewTemplateCmd(),
				template.NewValuesCmd(),
			},
		},
	}
	groups.Add(rootCmd)
	var filters []string
	templates.ActsAsRootCommand(rootCmd, filters, groups...)
}

func onBootOnDie() {
	if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		logger.Fatal("only support darwin and windows")
	}
	configs.DefaultClusterRootfsDir = clusterRootDir
	var rootDirs = []string{
		path.Join(clusterRootDir, "logs"),
		path.Join(clusterRootDir, "data"),
		path.Join(clusterRootDir, "etc"),
	}
	if err := file.MkDirs(rootDirs...); err != nil {
		logger.Error(err)
		panic(1)
	}
	logger.CfgConsoleAndFileLogger(debug, path.Join(clusterRootDir, "logs"), "sealvm", false)
}
