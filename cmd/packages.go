/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/interlynk-io/sbomgr/pkg/logger"
	"github.com/interlynk-io/sbomgr/pkg/search"
	"github.com/spf13/cobra"
)

type UserCmd struct {
	//Pattern Flags
	basicRegexp bool

	//Matching Control
	ignoreCase bool
	directDeps bool

	//Output Control
	license     bool
	quiet       bool
	filename    bool
	json        bool
	printErrors bool

	//stats Control
	count bool
	stats bool

	//Directory Control
	recurse bool

	//Spec Control
	spdx bool
	cdx  bool

	//Resource Control
	cpus int

	//Search Control
	name string
	cpe  string
	purl string
	hash string

	//Search Path
	path string

	//Output Format
	formats []string
}

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:          "packages",
	Short:        "search over packages in sboms",
	Long:         ``,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := logger.WithLoggerAndCancel(context.Background())
		defer stop()

		uCmd := toUserCmd(cmd, args[0])

		if err := validateFlags(uCmd); err != nil {
			return err
		}
		sp := toSearchParams(ctx, uCmd)
		return search.Search(sp)
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)

	packagesCmd.Flags().StringP("name", "N", "", "filter packages by name")
	packagesCmd.Flags().StringP("cpe", "C", "", "filter packages by cpe")
	packagesCmd.Flags().StringP("purl", "P", "", "filter packages by purl")
	packagesCmd.Flags().StringP("checksum", "H", "", "filter packages by checksum")
	packagesCmd.MarkFlagsMutuallyExclusive("cpe", "purl", "name", "checksum")

	packagesCmd.Flags().StringP("output-format", "O", "", "user-defined output format, comma separated list of columns. https://github.com/interlynk-io/sbomgr#output-control ")
}

func validatePath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

func validateFlags(cmd *UserCmd) error {
	if err := validatePath(cmd.path); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if cmd.basicRegexp {
		if _, err := regexp.Compile(cmd.name); err != nil {
			return fmt.Errorf("invalid regular expression: %w", err)
		}

		if _, err := regexp.Compile(cmd.cpe); err != nil {
			return fmt.Errorf("invalid regular expression: %w", err)
		}

		if _, err := regexp.Compile(cmd.purl); err != nil {
			return fmt.Errorf("invalid regular expression: %w", err)
		}

		if _, err := regexp.Compile(cmd.hash); err != nil {
			return fmt.Errorf("invalid regular expression: %w", err)
		}
	}

	if cmd.formats != nil {
		for _, f := range cmd.formats {
			_, ok := search.OutputFormatOptions[f]
			if !ok {
				return fmt.Errorf("invalid output format: %s", f)
			}
		}
	}
	return nil
}

func toSearchParams(ctx context.Context, cmd *UserCmd) *search.SearchParams {
	sp := &search.SearchParams{}

	sp.Ctx = ctx
	sp.Path = cmd.path

	sp.Regexp = cmd.basicRegexp
	sp.IgnoreCase = cmd.ignoreCase
	sp.DirectDeps = cmd.directDeps

	sp.License = cmd.license
	sp.Quiet = cmd.quiet
	sp.Filename = cmd.filename
	sp.Json = cmd.json
	sp.PrintErrors = cmd.printErrors

	sp.Count = cmd.count
	sp.Stats = cmd.stats

	sp.Recurse = cmd.recurse

	sp.Spdx = cmd.spdx
	sp.Cdx = cmd.cdx

	sp.Cpus = cmd.cpus

	sp.Name = cmd.name
	sp.CPE = cmd.cpe
	sp.PURL = cmd.purl
	sp.Hash = cmd.hash

	sp.Formats = cmd.formats

	return sp
}

func toUserCmd(cmd *cobra.Command, path string) *UserCmd {
	cUser := &UserCmd{}
	basicRegexp, _ := cmd.Flags().GetBool("extended-regexp")
	ignoreCase, _ := cmd.Flags().GetBool("ignore-case")
	directDeps, _ := cmd.Flags().GetBool("direct")
	license, _ := cmd.Flags().GetBool("license")
	quiet, _ := cmd.Flags().GetBool("quiet")
	filename, _ := cmd.Flags().GetBool("no-filename")
	json, _ := cmd.Flags().GetBool("jsonl")
	printErrors, _ := cmd.Flags().GetBool("print-errors")
	count, _ := cmd.Flags().GetBool("count")
	stats, _ := cmd.Flags().GetBool("stats")
	recurse, _ := cmd.Flags().GetBool("recurse")
	spdx, _ := cmd.Flags().GetBool("spdx")
	cdx, _ := cmd.Flags().GetBool("cdx")
	cpus, _ := cmd.Flags().GetInt("cpus")
	name, _ := cmd.Flags().GetString("name")
	cpe, _ := cmd.Flags().GetString("cpe")
	purl, _ := cmd.Flags().GetString("purl")
	hash, _ := cmd.Flags().GetString("checksum")
	formats, _ := cmd.Flags().GetString("output-format")

	cUser.basicRegexp = basicRegexp
	cUser.ignoreCase = ignoreCase
	cUser.directDeps = directDeps
	cUser.license = license
	cUser.quiet = quiet
	cUser.filename = filename
	cUser.json = json
	cUser.printErrors = printErrors
	cUser.count = count
	cUser.stats = stats
	cUser.recurse = recurse
	cUser.spdx = spdx
	cUser.cdx = cdx
	cUser.cpus = cpus
	cUser.name = name
	cUser.cpe = cpe
	cUser.purl = purl
	cUser.hash = hash

	cUser.path = path

	sanitize := func(s string, sep string) []string {
		splitStrings := strings.Split(s, sep)
		for i, str := range splitStrings {
			splitStrings[i] = strings.ToLower(strings.TrimSpace(str))
		}
		return splitStrings
	}

	if formats != "" {
		cUser.formats = sanitize(formats, ",")
	}
	return cUser
}
