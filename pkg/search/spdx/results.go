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
	"github.com/interlynk-io/sbomgr/pkg/search/results"

	spdx_common "github.com/spdx/tools-golang/spdx/common"
)

func (s *spdxDoc) constructResults(pIndices []int) (*results.Result, error) {
	if s.opts.BeQuiet() {
		return &results.Result{
			Path:    s.ro.CurrentPath,
			Matched: len(pIndices) > 0,
		}, nil
	}

	tools := s.toolInfo()

	result := &results.Result{
		Path:           s.ro.CurrentPath,
		Format:         string(s.ro.SbomFileFormat),
		Spec:           string(s.ro.SbomSpecType),
		ProductName:    s.docName(),
		ProductVersion: s.docVersion(),
		Packages:       s.pkgResults(pIndices),
		Files:          []results.File{},
		Matched:        len(pIndices) > 0,
	}

	// pick the first for now, later we probably need to return all
	if len(tools) > 0 {
		result.ToolName = tools[0].name
		result.ToolVersion = tools[0].version
	}

	return result, nil
}

func (s *spdxDoc) pkgResults(indices []int) []results.Package {
	pkgs := make([]results.Package, len(indices))
	for i, idx := range indices {
		pkgs[i] = results.Package{
			Name:    s.doc.Packages[idx].PackageName,
			Version: s.doc.Packages[idx].PackageVersion,
		}

		purls := s.Purl(idx)
		if len(purls) > 0 {
			pkgs[i].PURL = purls
		}

		cpes := s.CPEs(idx)
		if len(cpes) > 0 {
			pkgs[i].CPE = cpes
		}

		for _, c := range s.doc.Packages[idx].PackageChecksums {
			pkgs[i].Checksums = append(pkgs[i].Checksums, results.Checksum{
				Algorithm: string(c.Algorithm),
				Value:     string(c.Value),
			})
		}

		if s.opts.DoLicense() {
			ls := s.licenses(idx)
			for _, lic := range ls {
				pkgs[i].Licenses = append(pkgs[i].Licenses, licenses.LicenseStore{
					Nm: lic.Name(),
					Ss: lic.ShortID(),
					Ds: lic.Deprecated(),
				})
			}
		}
	}
	return pkgs
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

	removeDups := func(ll []licenses.License) []licenses.License {
		uniqs := []licenses.License{}
		dedup := map[string]bool{}
		for _, l := range ll {
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

func (s *spdxDoc) docVersion() string {
	return s.doc.DocumentNamespace
}

type tool struct {
	name    string
	version string
}

func (s *spdxDoc) toolInfo() []tool {
	tools := []tool{}

	if s.doc.CreationInfo == nil {
		return tools
	}

	//https://spdx.github.io/spdx-spec/v2.3/document-creation-information/#68-creator-field
	//spdx2.3 spec says If the SPDX document was created using a software tool,
	//indicate the name and version for that tool
	extractVersion := func(name string) (string, string) {
		//check if the version is a single word, i.e no spaces
		if strings.Contains(name, " ") {
			return name, ""
		}
		//check if name has - in it
		tool, ver, ok := strings.Cut(name, "-")

		if !ok {
			return name, ""
		}
		return tool, ver
	}

	for _, c := range s.doc.CreationInfo.Creators {
		ctType := strings.ToLower(c.CreatorType)
		if ctType != "tool" {
			continue
		}
		t := tool{}
		t.name, t.version = extractVersion(c.Creator)
		tools = append(tools, t)
	}

	return tools
}
