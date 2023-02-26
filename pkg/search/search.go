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

import "context"

type SearchParams struct {
	Ctx  context.Context
	Path string

	//optionals
	Name string
	CPE  string
	PURL string

	Exclude bool
	Direct  bool
}

type Option func(sp *SearchParams)

func WithName(name string) Option {
	return func(sp *SearchParams) {
		sp.Name = name
	}
}

func WithCpe(cpe string) Option {
	return func(sp *SearchParams) {
		sp.CPE = cpe
	}
}

func WithPurl(purl string) Option {
	return func(sp *SearchParams) {
		sp.PURL = purl
	}
}

func WithExclude(exclude bool) Option {
	return func(sp *SearchParams) {
		sp.Exclude = exclude
	}
}

func WithDirect(direct bool) Option {
	return func(sp *SearchParams) {
		sp.Direct = direct
	}
}

func NewSearchParams(ctx context.Context, path string, opts ...Option) *SearchParams {
	sp := &SearchParams{
		Ctx:  ctx,
		Path: path,
	}

	for _, opt := range opts {
		opt(sp)
	}

	return sp
}

type SearchResults struct {
	paths []string
}

func (sp SearchParams) Search() (*SearchResults, error) {
	return nil, nil
}
