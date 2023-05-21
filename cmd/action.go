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

package cmd

import (
	"github.com/labring/sealvm/pkg/actions"
	"github.com/spf13/cobra"
)

func newActionCmd() *cobra.Command {
	var file string
	var printDefault bool
	var actionCmd = &cobra.Command{
		Use:  "action",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if printDefault {
				return actions.PrintDefault()
			}
			return actions.Do(name, file)
		},
	}
	actionCmd.Flags().StringVarP(&name, "name", "n", "default", "name of cluster to applied init action")
	actionCmd.Flags().StringVarP(&file, "file", "f", "", "file to apply action")
	actionCmd.Flags().BoolVarP(&printDefault, "print-default", "p", false, "print default action")
	return actionCmd
}
