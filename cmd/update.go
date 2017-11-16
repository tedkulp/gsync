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
	"github.com/spf13/viper"
	"github.com/tedkulp/gsync/lib"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Syncs the watched files with git and pushes",
	Long: `Syncs the watched files with git and pushes

This does the following:
- Copies all watched files to the git repository
- Creates a commit if there are changes to commit
- Pushes to the git remote`,
	Aliases: []string{"cron", "sync"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := lib.ReadLines(filelist)
		if err != nil {
			fmt.Println(err)
			return
		}

		hostname := viper.GetString("hostname")
		remote := viper.GetString("remote")

		for _, line := range list {
			lib.CopyFileToRepo(line, hostname, repo)
			lib.GitAdd(line, hostname, repo)
		}

		if lib.GitHasChangesToCommit(repo) {
			if lib.GetHasRemote(repo, remote) {
				lib.GitPull(repo, remote)
			}

			lib.GitCommit(hostname, repo)

			if lib.GetHasRemote(repo, remote) {
				lib.GitPush(repo, remote)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
