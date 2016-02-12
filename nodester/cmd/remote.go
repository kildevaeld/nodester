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
	"strings"

	"github.com/kildevaeld/nodester"
	"github.com/spf13/cobra"
)

var detailedOutput bool
var printMax int
var longtermFlag bool

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		r, e := node.ListRemote(nodester.RemoteOptions{
			Lts: longtermFlag,
			Max: printMax,
		})

		if e != nil {
			writeError(e)
			return
		}

		if detailedOutput {
			return
		}

		printSimple(r)
	},
}

func init() {
	RootCmd.AddCommand(remoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// remoteCmd.PersistentFlags().String("foo", "", "A help for foo")
	remoteCmd.Flags().BoolVarP(&detailedOutput, "details", "d", false, "Get info")
	remoteCmd.Flags().BoolVarP(&longtermFlag, "long-term", "l", false, "Get info")
	remoteCmd.Flags().IntVarP(&printMax, "max", "m", 10, "Get info")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// remoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func printSimple(ms nodester.Manifests) {
	versions := make([]string, len(ms))
	for i, m := range ms {
		versions[i] = m.Version
	}
	s := strings.Join(versions, "\t")

	fmt.Println(s)
}
