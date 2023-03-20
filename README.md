<!--
 Copyright 2023 Interlynk.io
 
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 
     http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

# `sbomgr`: Search SBOMs

[![Go Reference](https://pkg.go.dev/badge/github.com/interlynk-io/sbomgr.svg)](https://pkg.go.dev/github.com/interlynk-io/sbomgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/interlynk-io/sbomgr)](https://goreportcard.com/report/github.com/interlynk-io/sbomgr)

`sbomgr` helps you search sboms based on criteria like name, checksum, cpe and purl. 

```sh
go install github.com/interlynk-io/sbomgr@latest
```
other installations [options](#installation)

# Basic usage
Search for packages with exact name match "abbrev". 
```sh
sbomgr packages -N 'abbrev' <sbom file or dir>
```

Search for packages with regexp name match "log4"
```sh
sbomgr packages -EN 'log4' <sbom file or dir>
```
# Features
- SBOM format agnostic, we support  SPDX and CycloneDX. 
- Fast. 
- Output search results as [jsonl](https://jsonlines.org/).
- Supports [regex](https://github.com/google/re2/wiki/Syntax)


# Use cases 
Lets say you have a folder/s of sboms. We wanted to answer the following questions

#### how many packages & sboms files exists?
```sh
➜ sbomgr packages -c ~/data/sbom-repo/docker-images
sbom_files_matched: 86
packages_matched: 33556
```
#### Are there any packages with zlib in the name?
```sh
➜ sbomgr packages -cEN 'zlib' ~/data/sbom-repo/docker-images
sbom_files_matched: 71
packages_matched: 145
```
#### Are there any packages with a known shasum 
```sh
➜ sbomgr packages -c -H '5c260231de4f62ee26888776190b4c3fda6cbe14' ~/data/sbom-repo/docker-images
sbom_files_matched: 2
packages_matched: 2
```

#### Are there any packages with .zip files and output it in json?
```sh
➜ sbomgr packages -jrE -N '\.zip$' ~/data/ | jq .
{
  "path": "/home/riteshno/data/spdx-trivy-circleci_clojure-sha256:d8944a6b1bec524314cf4889c104b302036690070a5353b64bb9d11b330e8c76.json",
  "format": "json",
  "spec": "spdx",
  "product_name": "circleci/clojure@sha256:d8944a6b1bec524314cf4889c104b302036690070a5353b64bb9d11b330e8c76",
  "packages": [
    {
      "name": "org.clojure:data.zip",
      "version": "0.1.3",
      "purl": "pkg:maven/org.clojure/data.zip@0.1.3"
    }
  ],
  "matched": true
}
```

#### List all packages of an sbom with their licenses
```sh
➜ sbomgr packages -jl ~/data/some-sboms/julia.spdx | jq .
{
  "path": "/home/riteshno/data/some-sboms/julia.spdx",
  "format": "tag-value",
  "spec": "spdx",
  "product_name": "julia-spdx",
  "packages": [
    {
      "name": "Julia",
      "version": "1.8.0-DEV",
      "license": [
        {
          "name": "MIT License",
          "short": "MIT"
        }
      ]
    },
```


# Future work
- Search using files
- Search using tool metadata
- Search using CVE-ID
- Search only direct dependencies 

# SBOM Samples
- A sample set of SBOM is present in the [samples](https://github.com/interlynk-io/sbomgr/tree/main/samples) directory above
- [SBOM Benchmark](https://www.sbombenchmark.dev) is a repository of SBOM and quality score for most popular containers and repositories
- [SBOM Explorer](https://github.com/interlynk-io/sbomex) is a command line utility to search and pull SBOMs

# Installation

## Using Prebuilt binaries

```console
https://github.com/interlynk-io/sbomgr/releases
```

## Using Homebrew
```console
brew tap interlynk-io/interlynk
brew install sbomgr
```

## Using Go install

```console
go install github.com/interlynk-io/sbomgr@latest
```

## Using repo

This approach involves cloning the repo and building it.

1. Clone the repo `git clone git@github.com:interlynk-io/sbomgr.git`
2. `cd` into `sbomgr` folder
3. make build
4. To test if the build was successful run the following command `./build/sbomgr version`


# Contributions
We look forward to your contributions, below are a few guidelines on how to submit them

- Fork the repo
- Create your feature/bug branch (`git checkout -b feature/new-feature`)
- Commit your changes (`git commit -am "awesome new feature"`)
- Push your changes (`git push origin feature/new-feature`)
- Create a new pull-request

# Contact
We appreciate all feedback, the best way to get in touch with us
- hello@interlynk.io
- github.com/interlynk-io/sbomgr/issues
- https://twitter.com/InterlynkIo