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

package spdx

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/interlynk-io/sbomgr/pkg/detect"
	spdx_json "github.com/spdx/tools-golang/json"
	spdx_rdf "github.com/spdx/tools-golang/rdfloader"
	spdx_common "github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_3"
	spdx_tv "github.com/spdx/tools-golang/tvloader"
	spdx_yaml "github.com/spdx/tools-golang/yaml"
)

type spdxDoc struct {
	doc *v2_3.Document
}

func loadDoc(s *SpdxModule) (*spdxDoc, error) {
	s.ro.File.Seek(0, io.SeekStart)

	var d *v2_3.Document
	var err error

	switch s.ro.SbomFileFormat {
	case detect.FileFormatJSON:
		d, err = spdx_json.Load2_3(s.ro.File)
	case detect.FileFormatTagValue:
		d, err = spdx_tv.Load2_3(s.ro.File)
	case detect.FileFormatYAML:
		d, err = spdx_yaml.Load2_3(s.ro.File)
	case detect.FileFormatRDF:
		d, err = spdx_rdf.Load2_3(s.ro.File)
	default:
		err = fmt.Errorf("unsupported spdx format %s", string(s.ro.SbomFileFormat))
	}

	if err != nil {
		return nil, err
	}

	doc := &spdxDoc{
		doc: d,
	}

	return doc, nil
}

func (d *spdxDoc) searchPackages(sm *SpdxModule) []int {
	if sm.so.SearchName() != "" {
		return d.matchEngine(sm, sm.so.SearchName(), func(pkg *v2_3.Package) []string {
			return []string{pkg.PackageName}
		})
	}

	if sm.so.SearchCPE() != "" {
		return d.matchEngine(sm, sm.so.SearchCPE(), func(pkg *v2_3.Package) []string {
			if len(pkg.PackageExternalReferences) == 0 {
				return []string{}
			}
			var cpe []string
			for _, ref := range pkg.PackageExternalReferences {
				if ref.RefType == spdx_common.TypeSecurityCPE23Type ||
					ref.RefType == spdx_common.TypeSecurityCPE22Type {
					cpe = append(cpe, ref.Locator)
				}
			}
			return cpe
		})
	}

	if sm.so.SearchPURL() != "" {
		return d.matchEngine(sm, sm.so.SearchPURL(), func(pkg *v2_3.Package) []string {
			if len(pkg.PackageExternalReferences) == 0 {
				return []string{}
			}
			var purl []string
			for _, ref := range pkg.PackageExternalReferences {
				if ref.RefType == spdx_common.TypePackageManagerPURL {
					purl = append(purl, ref.Locator)
				}
			}
			return purl
		})
	}

	if sm.so.SearchHash() != "" {
		return d.matchEngine(sm, sm.so.SearchHash(), func(pkg *v2_3.Package) []string {
			if len(pkg.PackageChecksums) == 0 {
				return []string{}
			}
			var checksums []string
			for _, cs := range pkg.PackageChecksums {
				checksums = append(checksums, cs.Value)
			}
			return checksums
		})
	}
	return []int{}
}

func (d *spdxDoc) matchEngine(sm *SpdxModule, matchCriteria string, mfunc func(*v2_3.Package) []string) []int {
	var regE *regexp.Regexp

	pkgIdx := []int{}

	if sm.so.IsRegularExp() {
		regE = regexp.MustCompile(matchCriteria)
	}

	for i, pkg := range d.doc.Packages {
		mSubject := mfunc(pkg)
		if len(mSubject) == 0 {
			continue
		}
		for _, s := range mSubject {
			if sm.so.IsRegularExp() {
				if sm.so.DoIgnoreCase() {
					if regE.MatchString(strings.ToLower(s)) {
						pkgIdx = append(pkgIdx, i)
					}
				} else {
					if regE.MatchString(s) {
						pkgIdx = append(pkgIdx, i)
					}
				}
			} else {
				if sm.so.DoIgnoreCase() {
					if strings.EqualFold(s, matchCriteria) {
						pkgIdx = append(pkgIdx, i)
					}
				} else {
					if strings.Compare(s, matchCriteria) == 0 {
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
