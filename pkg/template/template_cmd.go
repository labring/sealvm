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
	"errors"
	"fmt"
	"github.com/labring/sealvm/pkg/utils/confirm"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	var templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Display or change template config for sealvm.",
		Long:  "Template is cloud-init config, See https://cloudinit.readthedocs.io.",
	}
	templateCmd.AddCommand(newGetCmd())
	templateCmd.AddCommand(newSetCmd())
	templateCmd.AddCommand(newListCmd())
	templateCmd.AddCommand(newDefaultCmd())
	templateCmd.AddCommand(newRestCmd())
	return templateCmd
}

func newGetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:     "get <role>",
		Short:   "Print the value of a given template role",
		Args:    cobra.ExactArgs(1),
		Example: `sealvm template get key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := template{}
			data, err := t.Get(args[0])
			if err != nil {
				logger.Error(err.Error())
				return nil
			}
			fmt.Println(data)
			return nil
		},
	}
	return getCmd
}

func newSetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:     "set",
		Short:   "Update template with a path for the given role",
		Args:    cobra.NoArgs,
		Example: `sealvm template set`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := template{}
			_ = t.Set()
			return nil
		},
	}
	return getCmd
}

func newListCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:     "list",
		Short:   "Print a list of template",
		Args:    cobra.NoArgs,
		Example: `sealvm template list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := template{}
			logger.Info("Listing all templates:")
			t.List()
			return nil
		},
	}
	return getCmd
}

func newRestCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:     "reset",
		Short:   "Clean all template",
		Args:    cobra.NoArgs,
		Example: `sealvm template reset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := template{}
			prompt := "are you sure to reset all template?"
			cancelledMsg := "you have canceled to reset all template"
			yes, err := confirm.Confirm(prompt, cancelledMsg)
			if err != nil {
				return err
			}
			if !yes {
				return errors.New("cancelled")
			}
			logger.Info("Resetting all values...")
			t.Reset()
			logger.Info("Listing all templates:")
			t.List()
			return nil
		},
	}
	return getCmd
}

func newDefaultCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "default",
		Short: "Print default template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := template{}
			tpl := t.Default()
			fmt.Println(tpl)
			return nil
		},
	}
	return getCmd
}
