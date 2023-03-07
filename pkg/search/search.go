// Copyright 2023 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package search

import (
	"os"

	"github.com/interlynk-io/sbomgr/pkg/detect"
	"github.com/interlynk-io/sbomgr/pkg/search/options"
	"github.com/interlynk-io/sbomgr/pkg/search/results"
)

func searchFunc(path string, o options.SearchOptions) *results.Result {
	f, err := os.Open(path)
	if err != nil {
		return &results.Result{
			Path:  path,
			Error: err.Error(),
		}
	}
	defer f.Close()

	sbomSpecFormat, fileFormat, err := detect.Detect(f)
	if err != nil {
		return &results.Result{
			Path:  path,
			Error: err.Error(),
		}
	}

	if o.SpdxOnly() {
		if sbomSpecFormat != detect.SBOMSpecSPDX {
			return &results.Result{
				Path:  path,
				Error: "not an SPDX document",
			}
		}
	}

	if o.CdxOnly() {
		if sbomSpecFormat != detect.SBOMSpecCDX {
			return &results.Result{
				Path:  path,
				Error: "not a CycloneDX document",
			}
		}
	}

	ro := options.NewRuntimeOptions()
	ro.CurrentPath = path
	ro.SbomSpecType = sbomSpecFormat
	ro.SbomFileFormat = fileFormat
	ro.File = f

	sm := search_mods[sbomSpecFormat]
	sm.SetRuntimeOptions(ro)
	sm.SetSearchOptions(o)

	sr, err := sm.Search()
	if err != nil {
		return &results.Result{
			Path:  path,
			Error: err.Error(),
		}
	}
	return sr
}
