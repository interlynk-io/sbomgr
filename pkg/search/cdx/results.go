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
package cdx

import (
	cydx "github.com/CycloneDX/cyclonedx-go"
	"github.com/interlynk-io/sbomgr/pkg/licenses"
	"github.com/interlynk-io/sbomgr/pkg/search/results"
)

func (doc *cdxDoc) constructResults(pIndices []int) (*results.Result, error) {
	if doc.opts.BeQuiet() {
		return &results.Result{
			Path:    doc.ro.CurrentPath,
			Matched: len(pIndices) > 0,
		}, nil
	}

	var docName, docVersion string

	if doc.doc.Metadata != nil && doc.doc.Metadata.Component != nil {
		docName = doc.doc.Metadata.Component.Name
		docVersion = doc.doc.Metadata.Component.Version
	}

	result := &results.Result{
		Path:           doc.ro.CurrentPath,
		Format:         string(doc.ro.SbomFileFormat),
		Spec:           string(doc.ro.SbomSpecType),
		ProductName:    docName,
		ProductVersion: docVersion,
		Packages:       doc.pkgResults(pIndices),
		Files:          []results.File{},
		Matched:        len(pIndices) > 0,
	}

	if doc.doc.Metadata != nil && doc.doc.Metadata.Tools != nil {
		tools := doc.doc.Metadata.Tools
		if tools.Tools != nil && len(*tools.Tools) > 0 {
			tool := (*tools.Tools)[0]
			result.ToolName = tool.Name
			result.ToolVersion = tool.Version
		} else if tools.Components != nil && len(*tools.Components) > 0 {
			tool := (*tools.Components)[0]
			result.ToolName = tool.Name
			result.ToolVersion = tool.Version
		} else if tools.Services != nil && len(*tools.Services) > 0 {
			tool := (*tools.Services)[0]
			result.ToolName = tool.Name
			result.ToolVersion = tool.Version
		}
	}

	return result, nil
}

func (doc *cdxDoc) pkgResults(pIndices []int) []results.Package {
	var pkgs []results.Package

	for _, val := range pIndices {
		comp := doc.allComps[val]

		res := results.Package{
			Name:    comp.Name,
			Version: comp.Version,
		}

		if len(comp.PackageURL) > 0 {
			res.PURL = comp.PackageURL
		}

		if len(comp.CPE) > 0 {
			res.CPE = []string{comp.CPE}
		}

		if comp.ExternalReferences != nil && len(*comp.ExternalReferences) > 0 {
			for _, er := range *comp.ExternalReferences {
				if er.Type == cydx.ERTypeVCS {
					res.Repository = er.URL
					break
				}
			}
		}

		if comp.Hashes != nil {
			for _, c := range *comp.Hashes {
				res.Checksums = append(res.Checksums, results.Checksum{
					Algorithm: string(c.Algorithm),
					Value:     string(c.Value),
				})
			}
		}

		if doc.opts.DoLicense() {
			ls := doc.licenses(comp)
			for _, lic := range ls {
				res.Licenses = append(res.Licenses, licenses.LicenseStore{
					Nm: lic.Name(),
					Ss: lic.ShortID(),
					Ds: lic.Deprecated(),
				})
			}
		}

		pkgs = append(pkgs, res)
	}

	return pkgs
}

func (c *cdxDoc) licenses(comp *cydx.Component) []licenses.License {
	lics := []licenses.License{}

	addLicense := func(agg *[]licenses.License, n []licenses.License) {
		*agg = append(*agg, n...)
	}

	if comp.Licenses == nil {
		return []licenses.License{}
	}

	for _, cl := range *comp.Licenses {
		if cl.Expression != "" {
			addLicense(&lics, licenses.NewLicenseFromID(cl.Expression))
		} else if cl.License != nil {
			if cl.License.ID != "" {
				addLicense(&lics, licenses.NewLicenseFromID(cl.License.ID))
			} else {
				addLicense(&lics, []licenses.License{licenses.LicenseObjectByName(cl.License.Name)})
			}
		}
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
