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
	"context"
)

type SearchParams struct {
	//required
	Ctx  context.Context
	Path string

	//Patterns
	Regexp bool

	//Matching Control
	IgnoreCase  bool
	InvertMatch bool
	DirectDeps  bool

	//Output Control
	License     bool
	Quiet       bool
	Filename    bool
	Json        bool
	PrintErrors bool

	//stats Control
	Count bool
	Stats bool

	//Directory Control
	Recurse bool

	//Spec Control
	Spdx bool
	Cdx  bool

	//Resource Control
	Cpus int

	//Search Control
	Name string
	CPE  string
	PURL string
	Hash string
}

func (s *SearchParams) IsRegularExp() bool {
	return s.Regexp
}
func (s *SearchParams) Context() context.Context {
	return s.Ctx
}
func (s *SearchParams) RootPath() string {
	return s.Path
}

func (s *SearchParams) DoIgnoreCase() bool {
	return s.IgnoreCase
}

func (s *SearchParams) DoInverseMatch() bool {
	return s.InvertMatch
}

func (s *SearchParams) DoDirectDeps() bool {
	return s.DirectDeps
}

func (s *SearchParams) DoLicense() bool {
	return s.License
}

func (s *SearchParams) BeQuiet() bool {
	return s.Quiet
}

func (s *SearchParams) DoFilename() bool {
	return s.Filename
}

func (s *SearchParams) DoJson() bool {
	return s.Json
}

func (s *SearchParams) DoPrintErrors() bool {
	return s.PrintErrors
}

func (s *SearchParams) DoCount() bool {
	return s.Count
}

func (s *SearchParams) DoStats() bool {
	return s.Stats
}

func (s *SearchParams) DoRecurse() bool {
	return s.Recurse
}

func (s *SearchParams) SpdxOnly() bool {
	return s.Spdx
}

func (s *SearchParams) CdxOnly() bool {
	return s.Cdx
}

func (s *SearchParams) CpuToUse() int {
	return s.Cpus
}

func (s *SearchParams) SearchName() string {
	return s.Name
}

func (s *SearchParams) SearchCPE() string {
	return s.CPE
}

func (s *SearchParams) SearchPURL() string {
	return s.PURL
}

func (s *SearchParams) SearchHash() string {
	return s.Hash
}
