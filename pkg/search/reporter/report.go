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

package reporter

import (
	"context"
)

type Report struct {
	Ctx context.Context

	//optionals
	Direct   bool
	Stats    bool
	Licenses bool
}

type Option func(r *Report)

func WithDirect(direct bool) Option {
	return func(r *Report) {
		r.Direct = direct
	}
}

func WithStats(stats bool) Option {
	return func(r *Report) {
		r.Stats = stats
	}
}

func WithLicenses(licenses bool) Option {
	return func(r *Report) {
		r.Licenses = licenses
	}
}

func NewReport(ctx context.Context, opts ...Option) *Report {
	r := &Report{
		Ctx: ctx,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}
