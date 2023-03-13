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

func Search(s *SearchParams) error {
	ps := &pipeSetup{}
	ps.sParams = s
	ps.searchFunc = searchFunc

	if s.DoJson() {
		ps.outputFunc = outputJsonl
	} else if s.BeQuiet() {
		ps.outputFunc = outputQuiet
	} else {
		ps.outputFunc = outputBasic
	}

	if s.DoRecurse() {
		ps.fetchFilesFunc = fetchFilesRecursive
	} else {
		ps.fetchFilesFunc = fetchFiles
	}

	errs := runPipeline(ps)

	if ps.sParams.BeQuiet() {
		//In quiet mode, if we find even a single match, we return success
		for _, err := range errs {
			if err == nil {
				return nil
			}
		}
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
