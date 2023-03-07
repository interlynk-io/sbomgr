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
	"strings"

	"github.com/interlynk-io/sbomgr/pkg/licenses"
	"github.com/interlynk-io/sbomgr/pkg/search/options"
	"github.com/interlynk-io/sbomgr/pkg/search/results"

	spdx_common "github.com/spdx/tools-golang/spdx/common"
)

func (s *spdxDoc) constructResults(sm *SpdxModule, pIndices []int) (*results.Result, error) {
	if sm.so.BeQuiet() {
		return &results.Result{
			Path:    sm.ro.CurrentPath,
			Matched: len(pIndices) > 0,
		}, nil
	}
	return nil, nil
}

func (s *spdxDoc) docResults(ro *options.RuntimeOptions, opts options.SearchOptions, pIndices []int, _ []int) *results.Result {
	return &results.Result{
		Path:           ro.CurrentPath,
		Format:         string(ro.SbomFileFormat),
		Spec:           string(ro.SbomSpecType),
		ProductName:    s.docName(),
		ProductVersion: "",
		Packages:       s.pkgResults(pIndices),
		Files:          []results.File{},
	}
}

func (s *spdxDoc) pkgResults(indices []int) []results.Package {
	pkgs := make([]results.Package, len(indices))
	for i, idx := range indices {
		pkgs[i] = results.Package{
			Name:       s.doc.Packages[idx].PackageName,
			Version:    s.doc.Packages[idx].PackageVersion,
			PURL:       s.Purl(idx),
			CPE:        s.CPEs(idx),
			Direct:     s.directDep(idx),
			PathToRoot: s.pathToRoot(idx),
			License:    s.licenses(idx),
		}
	}
	return pkgs
}

func (s *spdxDoc) pathToRoot(index int) []string {
	//_ := s.doc.Packages[index]

	return []string{}
}

func (s *spdxDoc) directDep(index int) bool {
	//_ := s.doc.Packages[index]

	return false
}

func (s *spdxDoc) CPEs(index int) []string {
	urls := []string{}
	pkg := s.doc.Packages[index]
	if len(pkg.PackageExternalReferences) == 0 {
		return urls
	}

	for _, p := range pkg.PackageExternalReferences {
		if p.RefType == spdx_common.TypeSecurityCPE23Type || p.RefType == spdx_common.TypeSecurityCPE22Type {
			urls = append(urls, p.Locator)
		}

	}

	return urls
}

func (s *spdxDoc) Purl(index int) string {
	pkg := s.doc.Packages[index]

	if len(pkg.PackageExternalReferences) == 0 {
		return ""
	}

	for _, p := range pkg.PackageExternalReferences {
		if strings.ToLower(p.RefType) == spdx_common.TypePackageManagerPURL {
			return p.Locator
		}
	}

	return ""
}

func (s *spdxDoc) licenses(index int) []licenses.License {
	lics := []licenses.License{}
	pkg := s.doc.Packages[index]
	checkOtherLics := func(id string) (bool, string) {
		if s.doc.OtherLicenses == nil || len(s.doc.OtherLicenses) <= 0 {
			return false, ""
		}
		for _, l := range s.doc.OtherLicenses {
			if id == l.LicenseIdentifier {
				return true, l.ExtractedText
			}
		}
		return false, ""
	}

	addLicense := func(agg *[]licenses.License, n []licenses.License) {
		*agg = append(*agg, n...)
	}

	present, otherLic := checkOtherLics(pkg.PackageLicenseDeclared)

	if present {
		addLicense(&lics, licenses.NewLicenseFromID(otherLic))
	} else {
		addLicense(&lics, licenses.NewLicenseFromID(pkg.PackageLicenseDeclared))
	}

	present, otherLic = checkOtherLics(pkg.PackageLicenseConcluded)
	if present {
		addLicense(&lics, licenses.NewLicenseFromID(otherLic))
	} else {
		addLicense(&lics, licenses.NewLicenseFromID(pkg.PackageLicenseConcluded))
	}

	removeDups := func(lics []licenses.License) []licenses.License {
		uniqs := []licenses.License{}
		dedup := map[string]bool{}
		for _, l := range lics {
			if _, ok := dedup[l.ShortID()]; !ok {
				uniqs = append(uniqs, l)
				dedup[l.ShortID()] = true
			}
		}
		return uniqs

	}
	return removeDups(lics)
}

func (s *spdxDoc) docName() string {
	return s.doc.DocumentName
}
