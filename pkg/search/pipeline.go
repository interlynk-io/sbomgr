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
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/interlynk-io/sbomgr/pkg/search/options"
	"github.com/interlynk-io/sbomgr/pkg/search/results"
)

func fetchFiles(ps *pipeSetup) <-chan string {
	outc := make(chan string)

	go func() {
		defer close(outc)
		rs, err := os.Stat(ps.sParams.Path)
		if err != nil {
			log.Fatal(err)
		}

		if rs.Mode().IsRegular() {
			fmt.Printf("fullPath: %s\n", ps.sParams.Path)
			outc <- ps.sParams.Path
			return
		}

		files, err := ioutil.ReadDir(ps.sParams.Path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if ps.sParams.Ctx.Err() != nil {
				return
			}

			if !file.Mode().IsRegular() {
				continue
			}

			if file.Size() <= 0 {
				continue
			}

			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			fullPath := filepath.Join(ps.sParams.Path, file.Name())
			fmt.Printf("fullPath: %s\n", fullPath)
			outc <- fullPath
		}
	}()

	return outc

}

func fetchFilesRecursive(ps *pipeSetup) <-chan string {
	outc := make(chan string)

	go func() {
		defer close(outc)
		err := filepath.WalkDir(ps.sParams.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if ps.sParams.Ctx.Err() != nil {
				return filepath.SkipAll
			}

			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case outc <- path:
			case <-ps.sParams.Ctx.Done():
				return filepath.SkipAll
			}

			return nil
		})

		if err != nil {
			fmt.Printf("error walking path: %v\n", err)
		}
	}()
	return outc
}

func outputResults(ps *pipeSetup, sr <-chan *results.Result, errc <-chan error) []error {
	var outErr []error
	var cnt int = 0
out:
	for {
		select {
		case <-ps.sParams.Ctx.Done():
			break out
		case result, ok := <-sr:
			if !ok {
				break out
			}
			cnt += 1
			e := ps.outputFunc(result, ps.sParams)
			outErr = append(outErr, e)
		case err := <-errc:
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		}
	}
	fmt.Printf("output results processed: %d\n", cnt)
	return outErr
}

func stepSearch(ps *pipeSetup, inPathc <-chan string) (<-chan *results.Result, <-chan error) {
	outc := make(chan *results.Result)
	errc := make(chan error)
	var wg sync.WaitGroup

	go func() {
		defer close(outc)
		defer close(errc)

	outerloop:
		for {
			select {
			case <-ps.sParams.Ctx.Done():
				break outerloop
			case path, ok := <-inPathc:
				if !ok {
					break outerloop
				}
				wg.Add(1)
				go func(p string, sp *SearchParams) {
					defer wg.Done()
					select {
					case <-ps.sParams.Ctx.Done():
						return
					case outc <- ps.searchFunc(p, sp):
					}
				}(path, ps.sParams)
			}
		}
		wg.Wait()
	}()
	return outc, errc
}

func merge(ctx context.Context, err ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)
	monitor := func(c <-chan error) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(err))
	for _, e := range err {
		go monitor(e)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

type pipeSetup struct {
	sParams        *SearchParams
	searchFunc     func(string, options.SearchOptions) *results.Result
	outputFunc     func(*results.Result, *SearchParams) error
	fetchFilesFunc func(*pipeSetup) <-chan string
}

func runPipeline(sp *pipeSetup) []error {
	pathsChan := sp.fetchFilesFunc(sp)
	resultsChan, errorChan := stepSearch(sp, pathsChan)
	allErrorChan := merge(sp.sParams.Ctx, errorChan)
	return outputResults(sp, resultsChan, allErrorChan)
}
