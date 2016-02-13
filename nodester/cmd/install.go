// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"runtime"

	"github.com/kildevaeld/nodester"
	"github.com/spf13/cobra"
)

var sourceFlag bool
var archFlag string
var osFlag string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install a version of node",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Printf("Usage: ")
			return
		}

		for _, varg := range args {
			v := nodester.Version{
				Version: varg,
				Arch:    runtime.GOARCH,
				Os:      runtime.GOOS,
				Source:  sourceFlag,
			}

			fmt.Printf("  Installing %s\n", v.Version)

			if err := install(v); err != nil {
				fmt.Printf("    Error happended: %s\n", err.Error())
			}

		}

		fmt.Println("")

	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Aliases = []string{"i"}
	installCmd.Flags().BoolVarP(&sourceFlag, "source", "s", false, "install and build from source")
	installCmd.Flags().StringVarP(&archFlag, "arch", "a", runtime.GOARCH, "install and build from source")
	installCmd.Flags().StringVarP(&osFlag, "os", "o", runtime.GOOS, "install and build from source")

}

func install(v nodester.Version) error {

	if !node.Has(v) {
		err := NewProgress("  Downloading\t...", func(fn func(str string)) error {
			return node.Download(v, func(p, t int64) {
				fn(fmt.Sprintf("%d/%d kb", p/1024, t/1014))
			})
		})

		if err != nil {
			return err
		}

		err = NewProcess("  Unpacking\t...", func() error {
			return node.Install(v, nil)
		})

		return err
	}
	return nil
}
