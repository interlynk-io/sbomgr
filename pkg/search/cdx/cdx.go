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

package cdx

import (
	"github.com/interlynk-io/sbomgr/pkg/search/options"
	"github.com/interlynk-io/sbomgr/pkg/search/results"
)

type CdxModule struct {
	ro *options.RuntimeOptions
	so options.SearchOptions
}

func (c *CdxModule) SetRuntimeOptions(ro *options.RuntimeOptions) {
	c.ro = ro
}

func (c *CdxModule) SetSearchOptions(so options.SearchOptions) {
	c.so = so
}

func (c *CdxModule) Search() (*results.Result, error) {
	return &results.Result{
		Path:           c.ro.CurrentPath,
		Format:         string(c.ro.SbomFileFormat),
		Spec:           string(c.ro.SbomSpecType),
		ProductName:    "cdxtest",
		ProductVersion: "1.0",
		Packages:       []results.Package{},
		Files:          []results.File{},
	}, nil
}
