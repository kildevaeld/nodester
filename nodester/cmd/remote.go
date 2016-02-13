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
	"os"
	"strings"

	"github.com/kildevaeld/nodester"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var detailedOutput bool
var printMax int
var longtermFlag bool

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		r, e := node.ListRemote(nodester.RemoteOptions{
			Lts:                longtermFlag,
			Max:                printMax,
			HostCompatibleOnly: false,
		})

		if e != nil {
			writeError(e)
			return
		}

		if detailedOutput {
			printDetailed(r)
		} else {
			printSimple(r)
		}

	},
}

func init() {
	RootCmd.AddCommand(remoteCmd)

	remoteCmd.Flags().BoolVarP(&detailedOutput, "details", "d", false, "Get info")
	remoteCmd.Flags().BoolVarP(&longtermFlag, "long-term", "l", false, "Get info")
	remoteCmd.Flags().IntVarP(&printMax, "max", "m", 10, "Get info")

	remoteCmd.Aliases = []string{"lsr"}
}

func printSimple(ms nodester.Manifests) {
	versions := make([]string, len(ms))
	for i, m := range ms {
		versions[i] = m.Version
	}
	s := strings.Join(versions, "\t")

	fmt.Println(s)
}

func printDetailed(ms nodester.Manifests) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Version", "Longterm", "Date", "Npm", "V8", "Installed"})
	for _, v := range ms {
		installed := "X"
		if v.Installed {
			installed = "V"
		}
		table.Append([]string{v.Version, fmt.Sprintf("%v", v.Lts), v.Date, v.Npm, v.V8, installed})
	}

	table.Render()
}
