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

package search

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/interlynk-io/sbomgr/pkg/search/results"
	tw "github.com/olekukonko/tablewriter"
)

func outputQuiet(r *results.Result, nr *SearchParams) (int, error) {
	if r.Matched {
		return 1, nil
	}
	return 0, fmt.Errorf("no match found")
}
func outputJsonl(r *results.Result, nr *SearchParams) (int, error) {
	matchedPkgCount := len(r.Packages)

	if nr.DoCount() {
		if r.Matched {
			return matchedPkgCount, nil
		}
		return matchedPkgCount, fmt.Errorf("no match found")
	}
	if r.Error != "" && nr.DoPrintErrors() {
		return matchedPkgCount, fmt.Errorf("{\"path\":%s, \"error\": %s}", r.Path, r.Error)
	}

	// Looks like we have no packages that match the search criteria
	if len(r.Packages) == 0 {
		return matchedPkgCount, fmt.Errorf("no match found")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return matchedPkgCount, fmt.Errorf("error marshalling json: %w", err)
	}
	fmt.Println(string(b))
	return matchedPkgCount, nil
}
func outputBasic(r *results.Result, nr *SearchParams) (int, error) {
	matchedPkgCount := len(r.Packages)
	if nr.DoCount() {
		if r.Matched {
			return matchedPkgCount, nil
		}
		return matchedPkgCount, fmt.Errorf("no match found")
	}

	if r.Error != "" && nr.DoPrintErrors() {
		return matchedPkgCount, fmt.Errorf("path:%s error: %s", r.Path, r.Error)
	}
	// Looks like we have no packages that match the search criteria
	if len(r.Packages) == 0 {
		return matchedPkgCount, fmt.Errorf("no match found")
	}

	data := [][]string{}

	for _, pkg := range r.Packages {
		p := []string{}
		if !nr.DoFilename() {
			p = append(p, r.Path)
		}
		p = append(p, r.ProductName, r.ProductVersion, pkg.Name, pkg.Version, pkg.PURL)
		if nr.DoLicense() {
			var b []string
			for _, l := range pkg.Licenses {
				b = append(b, l.Name())
			}
			p = append(p, strings.Join(b, ","))
		}
		data = append(data, p)
	}

	table := tw.NewWriter(os.Stdout)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tw.ALIGN_RIGHT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()

	return matchedPkgCount, nil
}

type outputJson struct {
	SbomFilesMatched int `json:"sbom_files_matched"`
	PackagesMatched  int `json:"packages_matched"`
}

func handleFinalOutput(nr *SearchParams, matched []int, outputErrs []error) error {
	matchedSbomFilesCount := 0
	matchedItems := 0

	for _, pkgCnt := range matched {
		matchedItems += pkgCnt
	}

	for _, err := range outputErrs {
		if err == nil {
			matchedSbomFilesCount++
		}
	}

	if nr.DoCount() && nr.DoJson() {
		j := outputJson{
			SbomFilesMatched: matchedSbomFilesCount,
			PackagesMatched:  matchedItems,
		}
		b, _ := json.Marshal(j)
		fmt.Println(string(b))
		return nil
	}

	if nr.DoCount() && (!nr.DoJson() || !nr.BeQuiet()) {
		fmt.Printf("sbom_files_matched: %d\n", matchedSbomFilesCount)
		fmt.Printf("packages_matched: %d\n", matchedItems)
		return nil
	}

	if matchedSbomFilesCount == 0 {
		return fmt.Errorf("no matches found")
	}

	return nil
}
