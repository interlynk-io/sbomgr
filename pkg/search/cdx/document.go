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
	"fmt"
	"regexp"
	"strings"

	cydx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
	"github.com/interlynk-io/sbomgr/pkg/detect"
	"github.com/interlynk-io/sbomgr/pkg/search/options"
)

type cdxDoc struct {
	doc      *cydx.BOM
	ro       *options.RuntimeOptions
	opts     options.SearchOptions
	allComps []*cydx.Component
}

func loadDoc(ro *options.RuntimeOptions, opts options.SearchOptions) (*cdxDoc, error) {
	var err error
	var bom *cydx.BOM

	switch ro.SbomFileFormat {
	case detect.FileFormatJSON:
		bom = new(cydx.BOM)
		decoder := cydx.NewBOMDecoder(ro.File, cydx.BOMFileFormatJSON)
		if err = decoder.Decode(bom); err != nil {
			return nil, err
		}
	case detect.FileFormatXML:
		bom = new(cydx.BOM)
		decoder := cydx.NewBOMDecoder(ro.File, cydx.BOMFileFormatXML)
		if err = decoder.Decode(bom); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported cdx file format: %s", string(ro.SbomFileFormat))
	}

	doc := &cdxDoc{
		doc:      bom,
		ro:       ro,
		opts:     opts,
		allComps: extractAllComponents(bom),
	}
	return doc, nil
}

func extractAllComponents(bom *cydx.BOM) []*cydx.Component {
	var all_comps []*cydx.Component

	comps := map[string]*cydx.Component{}

	if bom.Metadata != nil && bom.Metadata.Component != nil {
		walkComponents(&[]cydx.Component{*bom.Metadata.Component}, comps)
	}

	if bom.Components != nil {
		walkComponents(bom.Components, comps)
	}

	for _, v := range comps {
		all_comps = append(all_comps, v)
	}
	return all_comps
}

func walkComponents(comps *[]cydx.Component, store map[string]*cydx.Component) {
	if comps == nil {
		return
	}
	for i, c := range *comps {
		if c.Components != nil {
			walkComponents(c.Components, store)
		}
		if _, ok := store[compID(&c)]; ok {
			//already present no need to re add it.
			continue
		}
		store[compID(&c)] = &(*comps)[i]
	}
}

func compID(comp *cydx.Component) string {
	if comp.BOMRef != "" {
		return comp.BOMRef
	}

	if comp.PackageURL != "" {
		return comp.PackageURL
	}

	// A component with no BOMREF or PackageURL is a bad component, so we generate a UUID for it.
	// This is a temporary solution until we can figure out how to handle this case.
	id := uuid.New()
	return id.String()
}

func (d *cdxDoc) searchPackages() []int {
	if d.opts.SearchName() != "" {
		return d.matchEngine(d.opts.SearchName(), func(pkg *cydx.Component) []string {
			return []string{pkg.Name}
		})
	}

	if d.opts.SearchCPE() != "" {
		return d.matchEngine(d.opts.SearchCPE(), func(pkg *cydx.Component) []string {
			return []string{pkg.CPE}
		})
	}

	if d.opts.SearchPURL() != "" {
		return d.matchEngine(d.opts.SearchPURL(), func(pkg *cydx.Component) []string {
			return []string{pkg.PackageURL}
		})
	}

	if d.opts.SearchHash() != "" {
		return d.matchEngine(d.opts.SearchHash(), func(pkg *cydx.Component) []string {
			var checksums []string
			if pkg.Hashes != nil {
				for _, hash := range *pkg.Hashes {
					checksums = append(checksums, hash.Value)
				}
			}
			return checksums
		})
	}

	allPkgs := []int{}
	for i, _ := range d.allComps {
		allPkgs = append(allPkgs, i)
	}
	return allPkgs
}

func (doc *cdxDoc) matchEngine(mCriteria string, mfunc func(*cydx.Component) []string) []int {
	var regE *regexp.Regexp

	pkgIdx := []int{}

	if doc.opts.IsRegularExp() {
		regE = regexp.MustCompile(mCriteria)
	}

	for i, comp := range doc.allComps {
		mSubject := mfunc(comp)
		if len(mSubject) == 0 {
			continue
		}
		for _, s := range mSubject {
			if doc.opts.IsRegularExp() {
				if doc.opts.DoIgnoreCase() {
					if regE.MatchString(strings.ToLower(s)) {
						pkgIdx = append(pkgIdx, i)
					}
				} else {
					if regE.MatchString(s) {
						pkgIdx = append(pkgIdx, i)
					}
				}
			} else {
				if doc.opts.DoIgnoreCase() {
					if strings.EqualFold(s, mCriteria) {
						pkgIdx = append(pkgIdx, i)
					}
				} else {
					if strings.Compare(s, mCriteria) == 0 {
						pkgIdx = append(pkgIdx, i)
					}
				}
			}
		}
	}
	uniqInts := func(in []int) []int {
		keys := make(map[int]bool)
		list := []int{}
		for _, entry := range in {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
		return list
	}
	return uniqInts(pkgIdx)
}
