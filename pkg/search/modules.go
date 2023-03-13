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

	"github.com/interlynk-io/sbomgr/pkg/detect"
	"github.com/interlynk-io/sbomgr/pkg/search/cdx"
	"github.com/interlynk-io/sbomgr/pkg/search/options"
	"github.com/interlynk-io/sbomgr/pkg/search/results"
	"github.com/interlynk-io/sbomgr/pkg/search/spdx"
)

type SearchModules interface {
	Search(*options.RuntimeOptions, options.SearchOptions) (*results.Result, error)
}

var search_mods = map[detect.SBOMSpecFormat]SearchModules{}

func init() {
	_ = registerSearchMod(detect.SBOMSpecCDX, &cdx.CdxModule{})
	_ = registerSearchMod(detect.SBOMSpecSPDX, &spdx.SpdxModule{})
}

func registerSearchMod(format detect.SBOMSpecFormat, mod SearchModules) error {
	if _, ok := search_mods[format]; ok {
		return fmt.Errorf("the format is being overwritten %s", format)
	}
	search_mods[format] = mod
	return nil
}
