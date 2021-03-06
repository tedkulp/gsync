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
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tedkulp/gsync/lib"
)

var cfgFile string
var filelist string
var repo string
var hostname string
var remote string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gsync",
	Short: "For syncing random files to a central git repository",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gsync/config)")
	RootCmd.PersistentFlags().StringVar(&filelist, "filelist", "", "list of watched files (default is $HOME/.gsync/filelist)")
	RootCmd.PersistentFlags().StringVar(&repo, "repo", "", "location of synced git repository (default is $HOME/.gsync/repo)")
	RootCmd.PersistentFlags().StringVar(&remote, "remote", "", "location of synced git repository (default is origin)")
	RootCmd.PersistentFlags().StringVar(&hostname, "hostname", "", "hostname to store this machine's data under (default is "+getHostname()+")")

	viper.BindPFlag("hostname", RootCmd.PersistentFlags().Lookup("hostname"))
	viper.SetDefault("hostname", getHostname())

	viper.BindPFlag("remote", RootCmd.PersistentFlags().Lookup("remote"))
	viper.SetDefault("remote", "origin")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		if lib.IsRoot() {
			viper.AddConfigPath("/etc/gsync")
		} else {
			viper.AddConfigPath(getDataDir())
		}
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if filelist == "" {
		filelist = path.Join(getDataDir(), "filelist")
	}

	if repo == "" {
		repo = path.Join(getDataDir(), "repo")
	}

	if hostname == "" {
		hostname = getHostname()
	}
}

func getDataDir() string {
	if lib.IsRoot() {
		return "/var/lib/gsync"
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return path.Join(home, ".gsync")
	}

}

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return name
}
