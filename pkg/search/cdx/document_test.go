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
	"strings"
	"testing"

	cydx "github.com/CycloneDX/cyclonedx-go"
)

func decodeBOM(t *testing.T, doc string) *cydx.BOM {
	t.Helper()
	bom := new(cydx.BOM)
	decoder := cydx.NewBOMDecoder(strings.NewReader(doc), cydx.BOMFileFormatJSON)
	if err := decoder.Decode(bom); err != nil {
		t.Fatal(err)
	}
	return bom
}

// A dependency entry is allowed to carry a ref and no dependsOn, which is how a
// component with no dependencies of its own is expressed.
const noDependsOn = `{
  "bomFormat": "CycloneDX",
  "specVersion": "1.5",
  "version": 1,
  "metadata": {
    "component": {"type": "application", "name": "root", "bom-ref": "root-ref"}
  },
  "components": [
    {"type": "library", "name": "lib", "bom-ref": "lib-ref"}
  ],
  "dependencies": [
    {"ref": "root-ref"}
  ]
}`

const withDependsOn = `{
  "bomFormat": "CycloneDX",
  "specVersion": "1.5",
  "version": 1,
  "metadata": {
    "component": {"type": "application", "name": "root", "bom-ref": "root-ref"}
  },
  "components": [
    {"type": "library", "name": "lib", "bom-ref": "lib-ref"}
  ],
  "dependencies": [
    {"ref": "root-ref", "dependsOn": ["lib-ref"]}
  ]
}`

func TestDirectComps(t *testing.T) {
	t.Run("primary component with no dependsOn", func(t *testing.T) {
		got := directComps(decodeBOM(t, noDependsOn))
		if len(got) != 0 {
			t.Errorf("expected no direct components, got %v", got)
		}
	})

	t.Run("primary component with dependsOn", func(t *testing.T) {
		got := directComps(decodeBOM(t, withDependsOn))
		if !got["lib-ref"] {
			t.Errorf("expected lib-ref to be a direct component, got %v", got)
		}
	})
}
