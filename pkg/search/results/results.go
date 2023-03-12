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

package results

import "github.com/interlynk-io/sbomgr/pkg/licenses"

type Package struct {
	Name       string                  `json:"name"`
	Version    string                  `json:"version"`
	PURL       string                  `json:"purl,omitempty"`
	CPE        []string                `json:"cpe,omitempty"`
	Direct     bool                    `json:"direct,omitempty"`
	PathToRoot []string                `json:"path_to_root,omitempty"`
	Licenses   []licenses.LicenseStore `json:"license,omitempty"`
}

type File struct {
	Name string `json:"name"`
}

type Result struct {
	Path           string    `json:"path"`
	Format         string    `json:"format"`
	Spec           string    `json:"spec"`
	Error          string    `json:"error,omitempty"`
	ProductName    string    `json:"product_name,omitempty"`
	ProductVersion string    `json:"product_version,omitempty"`
	Packages       []Package `json:"packages,omitempty"`
	Files          []File    `json:"files,omitempty"`
	Matched        bool      `json:"matched,omitempty"`
}

func NewSearchResult(path, format, spec, err string) *Result {
	return &Result{
		Path:   path,
		Format: format,
		Spec:   spec,
		Error:  err,
	}
}

func (sr *Result) AddProductName(name string) {
	sr.ProductName = name
}
func (sr *Result) AddProductVersion(version string) {
	sr.ProductVersion = version
}
func (sr *Result) AddPackage(pkg Package) {
	sr.Packages = append(sr.Packages, pkg)
}
func (sr *Result) AddFile(file File) {
	sr.Files = append(sr.Files, file)
}
