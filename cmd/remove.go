// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

	"github.com/spf13/cobra"
	"github.com/tedkulp/gsync/lib"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a file from the watch list",
	Long: `This removes an file from the list of watched files.

It will remove the file from the watch list, the git repository and
cause a commit. At the moment, it does not cause a push since this
will happen the next time update is run.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		removed, err := lib.RemoveLine(filelist, args[0])

		if err != nil {
			fmt.Println(err)
			return
		}

		if !removed && err == nil {
			fmt.Println(args[0] + " is not being watched")
			return
		}

		if removed && err == nil {
			lib.GitRemove(args[0], hostname, repo)

			if lib.GitHasChangesToCommit(repo) {
				lib.GitCommit(hostname, repo)
			}

			fmt.Println("Removed: " + args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
