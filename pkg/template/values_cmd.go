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
	"github.com/labring/sealvm/pkg/utils/confirm"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/spf13/cobra"
)

func NewValuesCmd() *cobra.Command {
	var templateCmd = &cobra.Command{
		Use:   "values",
		Short: "Display or change template values for sealvm.",
		Long: `This program provides a set of additional template functions:

See https://github.com/labring/sealvm/blob/main/pkg/tmpl/funcs.yaml
`,
	}
	templateCmd.AddCommand(newValuesSetCmd())
	templateCmd.AddCommand(newValuesListCmd())
	templateCmd.AddCommand(newValuesDefaultCmd())
	return templateCmd
}

func newValuesSetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "set",
		Short: "Update values for template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := values{}
			return t.Set()
		},
	}
	return getCmd
}

func newValuesListCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "list",
		Short: "Print a list of template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := values{}
			logger.Info("Listing all values:")
			t.List()
			return nil
		},
	}
	return getCmd
}

func newValuesDefaultCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "default",
		Short: "Set and display default values",
		Long: `This command will reset all values to their defaults and then display them. 
The default values are pre-defined and they include HTTPProxy, SocketProxy, NoProxy, PublicKey, PrivateKey, and ARCH.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := values{}
			prompt := "are you sure to set all values to default?"
			cancelledMsg := "you have canceled to set all values to default"
			yes, err := confirm.Confirm(prompt, cancelledMsg)
			if err != nil {
				return err
			}
			if !yes {
				return errors.New("cancelled")
			}
			logger.Info("Resetting values to defaults...")
			t.Default()
			logger.Info("Listing all default values:")
			t.List()
			logger.Info("Operation completed successfully.")
			return nil
		},
	}
	return getCmd
}
