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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tedkulp/gsync/lib"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <files>",
	Short: "Adds new files to the watch list",
	Long: `This adds new files to the list of watched files.

It does not add it to the git repository directly, as this will happen
the next time update is run.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Loop once and check all files
		for _, arg := range args {
			filename, err := filepath.Abs(arg)
			if err != nil {
				fmt.Println(arg + " is not a valid file. Exiting.")
				os.Exit(2)
			}

			if info, err := os.Stat(filename); err == nil && info.IsDir() {
				fmt.Println(arg + " is a directory. Exiting.")
				os.Exit(2)
			}
		}

		// Loop again and actually add them
		for _, arg := range args {
			filename, err := filepath.Abs(arg)
			added, err := lib.AddLine(filelist, filename)
			if added && err == nil {
				fmt.Println("Added: " + filename)
			} else if !added && err == nil {
				fmt.Println(filename + " already exists in file list. Skipping.")
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
