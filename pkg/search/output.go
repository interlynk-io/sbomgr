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

	"github.com/interlynk-io/sbomgr/pkg/search/results"
)

func outputQuiet(r *results.Result, nr *SearchParams) error {
	if r.Matched {
		return nil
	}
	return fmt.Errorf("no match found")
}
func outputJsonl(r *results.Result, nr *SearchParams) error {
	if r.Error != "" && nr.DoPrintErrors() {
		return fmt.Errorf("{\"path\":%s, \"error\": %s}", r.Path, r.Error)
	}

	// Looks like we have no packages that match the search criteria
	if len(r.Packages) == 0 {
		return fmt.Errorf("no match found")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}
	fmt.Println(string(b))
	return nil
}
func outputBasic(r *results.Result, nr *SearchParams) error {
	if r.Error != "" && nr.DoPrintErrors() {
		return fmt.Errorf("path:%s error: %s", r.Path, r.Error)
	}
	// Looks like we have no packages that match the search criteria
	if len(r.Packages) == 0 {
		return fmt.Errorf("no match found")
	}

	for _, pkg := range r.Packages {
		fmt.Println(r.Path, r.ProductName, r.ProductVersion, pkg.Name, pkg.Version, pkg.PURL)
	}
	return nil
}

func handleFinalOutput(nr *SearchParams, outputErrs []error) error {
	matchedCount := 0

	for _, err := range outputErrs {
		if err == nil {
			matchedCount++
		}
	}

	if nr.DoCount() && nr.DoJson() {
		fmt.Printf("{\"count\": %d}\n", matchedCount)
		return nil
	}

	if nr.DoCount() && (!nr.DoJson() || !nr.BeQuiet()) {
		fmt.Printf("count: %d\n", matchedCount)
		return nil
	}

	if matchedCount == 0 {
		return fmt.Errorf("no matches found")
	}

	return nil
}
