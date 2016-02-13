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

var forceFlag bool

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			return
		}

		version := nodester.Version{
			Version: args[0],
			Arch:    runtime.GOARCH,
			Os:      runtime.GOOS,
		}

		if !node.Has(version) {
			if !forceFlag {
				fmt.Printf("  %s is not installed\n", version.Version)
				return
			}
			err := install(version)
			if err != nil {
				writeError(err)
			}
		}

		err := NewProcess("  Activating "+args[0], func() error {
			return node.Use(version)
		})

		if err != nil {
			fmt.Printf("    Could not activate %s: %s: \n", args[0], err.Error())
		}

	},
}

func init() {
	RootCmd.AddCommand(useCmd)

	useCmd.Aliases = []string{"u"}

	useCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force installation")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
