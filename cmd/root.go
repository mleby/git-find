// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/martinlebeda/mlgotools/scripttools"
	"gopkg.in/pipe.v2"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-find",
	Args:  cobra.ExactArgs(1), // TODO Lebeda - more then 1 arg
	Short: "find keyword in got log.",
	Long:  `Find keyword in got log and list all commits, tags and branches contain it.`,
	// TODO Lebeda - usage
	Run: func(cmd *cobra.Command, args []string) {

		// find commits
		p := pipe.Line(
			pipe.Exec("git", "--no-pager", "log", "--all", "--oneline", "-i", "--grep="+args[0]),
			//	pipe.Filter(func(line []byte) bool {return bytes.Contains(line, []byte("MCV")) }),
		)
		split := scripttools.GetOutputLines(p)

		var commits []string

		if len(split) > 0 {
			fmt.Println("commits:")
			for _, logLine := range split {
				fmt.Println(logLine)
				logSplit := strings.Split(logLine, " ")
				commits = append(commits, strings.TrimSpace(logSplit[0]))
			}
		} else {
			fmt.Println("no commit found")
		}

		if len(commits) > 0 {
			findTags(commits)
			findBranches(commits, false)
			findBranches(commits, true)
		}

	},
}

// find tags
func findTags(commits []string) {
	var tags []string
	for _, commit := range commits {
		//fmt.Println("git", "tag", "--contains", commit)
		p := pipe.Line(
			pipe.Exec("git", "--no-pager", "tag", "--contains", commit),
			//	pipe.Filter(func(line []byte) bool {return bytes.Contains(line, []byte("MCV")) }),
		)
		outputLines := scripttools.GetOutputLines(p)
		tags = append(tags, outputLines...)
	}
	if len(tags) > 0 {
		//sort.Slice(people, func(i, j int) bool {
		//	return people[i].Age > people[j].Age
		//})
		sort.Strings(tags)
		// TODO Lebeda - unique
		fmt.Println("\ntags:")
		for _, tag := range tags {
			fmt.Println(" ", tag)
		}

	} else {
		fmt.Println("\nno tag found")
	}
}

func findBranches(commits []string, remote bool) {
	var branchType string
	var branches []string

	for _, commit := range commits {
		var exec pipe.Pipe
		if remote {
			exec = pipe.Exec("git", "--no-pager", "branch", "-r", "--contains", commit)
			branchType = "remote"
		} else {
			exec = pipe.Exec("git", "--no-pager", "branch", "--contains", commit)
			branchType = "local"
		}
		p := pipe.Line(
			exec,
			//	pipe.Filter(func(line []byte) bool {return bytes.Contains(line, []byte("MCV")) }),
		)
		outputLines := scripttools.GetOutputLines(p)
		branches = append(branches, outputLines...)
	}
	if len(branches) > 0 {
		sort.Strings(branches)
		branches = scripttools.RemoveDuplicates(branches)
		fmt.Println("\n" + branchType + " branches:")
		for _, tag := range branches {
			fmt.Println(tag)
		}

	} else {
		fmt.Println("\nno " + branchType + " branch found")
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-find.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO Lebeda - commits
	// TODO Lebeda - ignore RC
	// TODO Lebeda - sort by numeric version padded from left/right
	// TODO Lebeda - only tags/local branch/remote branch
	// TODO Lebeda - run in specific directory
	// TODO Lebeda - -i, --regexp-ignore-case // Match the regular expression limiting patterns without regard to letter case.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// Search config in home directory with name ".git-find" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigName(".git-find")
	//}
	//
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
}
