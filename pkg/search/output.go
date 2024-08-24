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

var OutputFormatOptions = map[string]bool{
	"filen": true, //filename
	"tooln": true, //toolname
	"toolv": true, //toolversion
	"docn":  true, //doc name
	"docv":  true, //doc version
	"cpe":   true, //cpe
	"purl":  true, //purl
	"pkgn":  true, //component name
	"pkgv":  true, //component version
	"pkgl":  true, //component license
	"specn": true, //spec name
	"chkn":  true, //checksum name
	"chkv":  true, //checksum version
	"repo":  true, //repository
}

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

	newPackages := []results.Package{}

	for _, p := range r.Packages {
		newP := results.Package{}
		for _, f := range nr.Formats {
			switch f {
			case "cpe":
				newP.CPE = p.CPE
			case "purl":
				newP.PURL = p.PURL
			case "pkgn":
				newP.Name = p.Name
			case "pkgv":
				newP.Version = p.Version
			case "pkgl":
				newP.Licenses = p.Licenses
			case "chkn":
				newP.Checksums = p.Checksums
			case "chkv":
				newP.Checksums = p.Checksums
			case "repo":
				newP.Repository = p.Repository
			}
		}
		newPackages = append(newPackages, newP)
	}

	r.Packages = newPackages

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

	isEmpty := func(s []string) bool {
		for _, v := range s {
			if strings.Trim(v, " ") != "" {
				return false
			}
		}
		return true
	}

	data := [][]string{}

	for idx, pkg := range r.Packages {
		p := []string{}

		if nr.HasOutputFormats() {
			noOfChecksums := len(pkg.Checksums)
			if noOfChecksums > 0 {
				for id := range pkg.Checksums {
					data = append(data, customOutput(idx, id, r, nr))
				}
			} else {
				p = customOutput(idx, -1, r, nr)
			}
		} else {
			if !nr.DoFilename() {
				p = append(p, r.Path)
			}
			p = append(p, r.ProductName, r.ProductVersion, pkg.Name, pkg.Version)
			if nr.DoLicense() {
				var b []string
				for _, l := range pkg.Licenses {
					b = append(b, l.Name())
				}
				p = append(p, strings.Join(b, ","))
			}
		}

		if !isEmpty(p) {
			data = append(data, p)
		}
	}

	table := tw.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tw.ALIGN_LEFT)
	table.SetAlignment(tw.ALIGN_LEFT)

	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()

	return matchedPkgCount, nil
}

type outputJson struct {
	SbomFilesMatched int `json:"files_matched"`
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
		fmt.Printf("Matching file count: %d\n", matchedSbomFilesCount)
		fmt.Printf("Matching package count: %d\n", matchedItems)
		return nil
	}

	if matchedSbomFilesCount == 0 {
		return fmt.Errorf("no matches found")
	}

	return nil
}

func customOutput(idx int, chkIdx int, r *results.Result, nr *SearchParams) []string {
	p := []string{}
	pkg := r.Packages[idx]

	for _, f := range nr.Formats {
		switch f {
		case "filen":
			p = append(p, r.Path)
		case "tooln":
			if len(r.ToolName) > 0 {
				p = append(p, r.ToolName)
			} else {
				p = append(p, "[NOTOOL]")
			}
		case "toolv":
			if len(r.ToolVersion) > 0 {
				p = append(p, r.ToolVersion)
			} else {
				p = append(p, "[NOTOOLVER]")
			}
		case "docn":
			if len(r.ProductName) > 0 {
				p = append(p, r.ProductName)
			} else {
				p = append(p, "[NODOC]")
			}
		case "docv":
			if len(r.ProductVersion) > 0 {
				p = append(p, r.ProductVersion)
			} else {
				p = append(p, "[NODOCVER]")
			}
		case "cpe":
			if len(pkg.CPE) > 0 {
				cpef := fmt.Sprintf("%s[%d more]", pkg.CPE[0], len(pkg.CPE))
				p = append(p, cpef)
			} else {
				p = append(p, "[NOCPE]")
			}
		case "purl":
			if len(pkg.PURL) > 0 {
				p = append(p, pkg.PURL)
			} else {
				p = append(p, "[NOPURL]")
			}
		case "pkgn":
			if len(pkg.Name) > 0 {
				p = append(p, pkg.Name)
			} else {
				p = append(p, "[NOPKGNAME]")
			}
		case "pkgv":
			if len(pkg.Version) > 0 {
				p = append(p, pkg.Version)
			} else {
				p = append(p, "[NOPKGVER]")
			}
		case "pkgl":
			var b []string
			for _, l := range pkg.Licenses {
				b = append(b, l.Name())
			}
			p = append(p, strings.Join(b, ","))
		case "specn":
			p = append(p, r.Spec)
		case "chkn":
			if chkIdx >= 0 {
				p = append(p, pkg.Checksums[chkIdx].Algorithm)
			} else {
				p = append(p, "[NOCHKN]")
			}
		case "chkv":
			if chkIdx >= 0 {
				p = append(p, pkg.Checksums[chkIdx].Value)
			} else {
				p = append(p, "[NOCHKV]")
			}
		case "repo":
			if chkIdx >= 0 {
				p = append(p, pkg.Repository)
			} else {
				p = append(p, "[NOREPO]")
			}
		}
	}

	return p

}
