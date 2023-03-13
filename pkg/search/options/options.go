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

package options

import (
	"context"
	"io"

	"github.com/interlynk-io/sbomgr/pkg/detect"
)

type SearchOptions interface {
	IsRegularExp() bool
	Context() context.Context
	RootPath() string
	DoIgnoreCase() bool
	DoInverseMatch() bool
	DoDirectDeps() bool
	DoLicense() bool
	BeQuiet() bool
	DoFilename() bool
	DoJson() bool
	DoPrintErrors() bool
	DoCount() bool
	DoStats() bool
	DoRecurse() bool
	SpdxOnly() bool
	CdxOnly() bool
	CpuToUse() int
	SearchName() string
	SearchCPE() string
	SearchPURL() string
	SearchHash() string
}

type RuntimeOptions struct {
	CurrentPath    string
	SbomSpecType   detect.SBOMSpecFormat
	SbomFileFormat detect.FileFormat
	File           io.ReadSeeker
}

func NewRuntimeOptions() *RuntimeOptions {
	return &RuntimeOptions{}
}
