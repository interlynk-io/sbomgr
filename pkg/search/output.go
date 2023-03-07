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
	return nil
}
func outputBasic(r *results.Result, nr *SearchParams) error {
	return nil
}
