/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/interlynk-io/sbomgr/pkg/logger"
	"github.com/interlynk-io/sbomgr/pkg/reporter"
	"github.com/interlynk-io/sbomgr/pkg/search"
	"github.com/spf13/cobra"
)

var (
	name    string
	exclude bool
	cpe     string
	purl    string
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "search for packages",
	Long:  `search thru all your packages`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := logger.WithLogger(context.Background())

		if validatePath(args[0]) != nil {
			return fmt.Errorf("invalid path: %s", args[0])
		}
		fmt.Println("Print:" + strings.Join(args, " "))

		direct, licenses, stats, exclude := extractFlags(cmd)

		sp := search.NewSearchParams(ctx, args[0],
			search.WithName(name), search.WithCpe(cpe),
			search.WithPurl(purl), search.WithExclude(exclude),
			search.WithDirect(direct))

		results, error := sp.Search()
		if error != nil {
			return error
		}

		nr := reporter.NewReport(ctx, results,
			reporter.WithLicenses(licenses),
			reporter.WithStats(stats),
			reporter.WithDirect(direct))

		nr.Report()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)

	packagesCmd.Flags().BoolVarP(&exclude, "exclude", "e", false, "Exclude packages")
	packagesCmd.Flags().StringVarP(&name, "name", "n", "", "filter packages by name")
	packagesCmd.Flags().StringVarP(&cpe, "cpe", "c", "", "filter packages by cpe")
	packagesCmd.Flags().StringVarP(&purl, "purl", "p", "", "filter packages by purl")
	packagesCmd.MarkFlagsMutuallyExclusive("cpe", "purl", "name")
}

func validatePath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

func extractFlags(cmd *cobra.Command) (bool, bool, bool, bool) {
	direct, _ := cmd.Flags().GetBool("direct")
	licenses, _ := cmd.Flags().GetBool("licenses")
	stats, _ := cmd.Flags().GetBool("stats")
	return direct, licenses, stats, exclude
}
