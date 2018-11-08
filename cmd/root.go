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

var optCommits, optNoRc bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-find",
	Args:  cobra.MinimumNArgs(1),
	Short: "find keyword in git log.",
	Long:  `Find keyword in git log and list all commits, tags and branches contain it.`,
	Example: `  git-find KEY1 KEY2 ...
  git-find AA-35210 AA-35211

  git-find -c commit1 commit2 ...
  git-find -c d7c2924b17 e9ac8dd7bd`,

	// TODO Lebeda - Version Flag
	Run: func(cmd *cobra.Command, args []string) {

		// find commits
		var commits []string
		if optCommits {
			commits = scripttools.RemoveDuplicates(args)
		} else {
			commits = getCommits(args)
		}

		noPattern := ""
		if optNoRc {
			noPattern = "rc"
		}

		// print tags and branches contain this commits
		if len(commits) > 0 {
			findTags(commits, noPattern)
			findBranches(commits, false)
			findBranches(commits, true)
		}

	},
}

func getCommits(patterns []string) []string {
	var logLines []string

	for _, value := range patterns {
		p := pipe.Exec("git", "--no-pager", "log", "--all", "--oneline", "-i", "--grep="+value)
		logLines = append(logLines, scripttools.GetOutputLines(p)...)
	}

	scripttools.RemoveDuplicates(logLines)

	var commits []string
	if len(logLines) > 0 {
		scripttools.Header("commits:")
		for _, logLine := range logLines {
			fmt.Println(logLine)
			logSplit := strings.Split(logLine, " ")
			commits = append(commits, strings.TrimSpace(logSplit[0]))
		}
	} else {
		fmt.Println("no commit found")
	}

	return scripttools.RemoveDuplicates(commits)
}

// find tags
func findTags(commits []string, noPattern string) {
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
		tags = scripttools.RemoveDuplicates(tags)
		if noPattern != "" {
			tags = scripttools.RemovePattern(tags, noPattern)
		}
		scripttools.Header("\ntags:")
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
		scripttools.Header("\n" + branchType + " branches:")
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
	rootCmd.Flags().BoolVarP(&optCommits, "commits", "c", false, "use commits instead find pattern")
	rootCmd.Flags().BoolVarP(&optNoRc, "ignore-rc", "r", false, "ignore tags with 'rc' in name")
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
