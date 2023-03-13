// Copyright 2023 Interlynk.io
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
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sbomgr",
	Short: "a modern tool to search sboms",
	Long: `
sbomgr is a modern tool to search sboms. It is designed to be fast 
and easy to use. It is a command line tool that can be used to search
sboms for packages and files. 
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//Pattern
	rootCmd.PersistentFlags().BoolP("extended-regexp", "E", false, "intrepret filters as regular expressions, https://github.com/google/re2/wiki/Syntax")

	//Matching Control
	rootCmd.PersistentFlags().BoolP("ignore-case", "i", false, "ignore case distinctions in filters, lowers the package/file criterias")
	rootCmd.PersistentFlags().BoolP("direct-deps", "d", false, "search direct dependencies only, default is to search all packages/files")
	rootCmd.PersistentFlags().MarkHidden("direct-deps")

	//Output Control
	rootCmd.PersistentFlags().BoolP("license", "l", false, "output with licenses")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress normal output, exits with 0 if match found")
	rootCmd.PersistentFlags().BoolP("no-filename", "", false, "output with no filename")
	rootCmd.PersistentFlags().BoolP("jsonl", "j", false, "results in json-lines format https://jsonlines.org/")
	rootCmd.PersistentFlags().BoolP("print-errors", "p", false, "include errors in output")

	//stats Control
	rootCmd.PersistentFlags().BoolP("count", "c", false, "suppress normal output, print count of matching packages with pattern")
	rootCmd.PersistentFlags().BoolP("stats", "s", false, "suppress normal output, print stats of matching packages/files")
	rootCmd.MarkFlagsMutuallyExclusive("count", "stats")
	rootCmd.PersistentFlags().MarkHidden("stats")

	//Directory Control
	rootCmd.PersistentFlags().BoolP("recurse", "r", false, "recurse into subdirectories")

	//Spec Control
	rootCmd.PersistentFlags().BoolP("spdx", "", false, "limit searches to spdx sboms")
	rootCmd.PersistentFlags().BoolP("cdx", "", false, "limit searches to cdx sboms")
	rootCmd.MarkFlagsMutuallyExclusive("spdx", "cdx")

	//Resource Control
	rootCmd.PersistentFlags().IntP("cpus", "", 0, "restrict number of cpus, default is all")
}
