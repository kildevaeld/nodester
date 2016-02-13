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

		v := nodester.Version{
			Version: args[0],
			Arch:    runtime.GOARCH,
			Os:      runtime.GOOS,
			Source:  sourceFlag,
		}
		fmt.Printf("  Installing %s\n", v.Version)
		if !node.Has(v) {
			err := NewProgress("  Downloading\t...", func(fn func(str string)) error {
				return node.Download(v, func(p, t int64) {
					fn(fmt.Sprintf("%d/%d kb", p/1024, t/1014))
				})
			})

			if err != nil {
				writeError(err)
			}

			err = NewProgress("  Installing\t...", func(fn func(str string)) error {
				return node.Install(v, func(s nodester.Step) {
					fn(s.String())
				})
			})

			fmt.Println("")

			if err != nil {
				writeError(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Aliases = []string{"i"}
	installCmd.Flags().BoolVarP(&sourceFlag, "source", "s", false, "install and build from source")
	installCmd.Flags().StringVarP(&archFlag, "arch", "a", runtime.GOARCH, "install and build from source")
	installCmd.Flags().StringVarP(&osFlag, "os", "o", runtime.GOOS, "install and build from source")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func install(v nodester.Version) error {
	fmt.Printf("  Installing %s\n", v.Version)
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

		fmt.Println("")

		return err
	}
	return nil
}
