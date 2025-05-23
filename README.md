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

# `sbomgr`: SBOM Grep :mag: - Search through SBOMs

[![Go Reference](https://pkg.go.dev/badge/github.com/interlynk-io/sbomgr.svg)](https://pkg.go.dev/github.com/interlynk-io/sbomgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/interlynk-io/sbomgr)](https://goreportcard.com/report/github.com/interlynk-io/sbomgr)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/interlynk-io/sbomgr/badge)](https://securityscorecards.dev/viewer/?uri=github.com/interlynk-io/sbomgr)
![GitHub all releases](https://img.shields.io/github/downloads/interlynk-io/sbomgr/total)

`sbomgr` is a grep like command line utility to help search the SBOM repository based on criteria like the name, checksum, CPE, and PURL.

```sh
go install github.com/interlynk-io/sbomgr@latest
```

other installations [options](#installation)

# SBOM Platform - Free Community Tier

Our SBOM Automation Platform has a free community tier that provides a comprehensive solution to manage SBOMs (Software Bill of Materials) effortlessly. From centralized SBOM storage, built-in SBOM editor, continuous vulnerability mapping and assessment, and support for organizational policies, all while ensuring compliance and enhancing software supply chain security using integrated SBOM quality scores. The community tier is ideal for small teams. Learn more [here](https://www.interlynk.io/community-tier) or [Sign up](https://app.interlynk.io/auth)

# SBOM Card

[![SBOMCard](https://api.interlynk.io/api/v1/badges?type=hcard&project_group_id=e8e2ba0c-3d04-4a2e-9b37-dca774bd08bd)](https://app.interlynk.io/customer/products?id=e8e2ba0c-3d04-4a2e-9b37-dca774bd08bd&signed_url_params=eyJfcmFpbHMiOnsibWVzc2FnZSI6IklqSmtaakkyTkRRMUxXSTBaR0V0TkdJME9TMWhPVFpqTFRBd09UZGtZMlptTWpabU9TST0iLCJleHAiOm51bGwsInB1ciI6InNoYXJlX2x5bmsvc2hhcmVfbHluayJ9fQ==--6d74d14e40d6676522b1c529d44e4a320f05bcf3d42121e61e1275a1297a3453)

# Basic usage

Search for packages with exact name matching "abbrev".

```sh
sbomgr packages -N 'abbrev' <sbom file or dir>
```

Search for packages with regexp name matching "log4"

```sh
sbomgr packages -EN 'log4' <sbom file or dir>
```

Search for packages in air gapped environment for name matching "log4"

```sh
export INTERLYNK_DISABLE_VERSION_CHECK=true sbomgr packages -EN 'log4' <sbom file or dir>
```

# Features

- SBOM format agnostic and currently supports searching through SPDX and CycloneDX.
- Blazing Fast :rocket:
- Output search results as [jsonl](https://jsonlines.org/).
- Supports RE2 [regular expressions](https://github.com/google/re2/wiki/Syntax)

# Use cases

`sbomgr` can answer some of the most common SBOM use cases by searching an SBOM file or SBOM repository.

## How many SBOM and packages exist in the repository?

```sh
➜ sbomgr packages -c ~/data/sbom-repo/docker-images
sbom_files_matched: 86
packages_matched: 33556
```

## Are there packages with `zlib` in the name?

```sh
➜ sbomgr packages -cEN 'zlib' ~/data/sbom-repo/docker-images
sbom_files_matched: 71
packages_matched: 145
```

## Are there packages with a given checksum?

```sh
➜ sbomgr packages -c -H '5c260231de4f62ee26888776190b4c3fda6cbe14' ~/data/sbom-repo/docker-images
sbom_files_matched: 2
packages_matched: 2
```

## Create a json report of packages with .zip files

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

## Create a json report of all licenses included in an sbom

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

## During CI check if a malicious package is present??

```sh
➜  sbomgr packages -qN 'abbrev' ~/tmp/app.spdx.json
➜  echo $?
0
➜  sbomgr packages -qN 'abbrev-random' ~/tmp/app.spdx.json
➜  echo $?
1
```

## extract data using user-defined output

```sh
sbomgr packages -O 'toolv,tooln,pkgn,pkgv' ~/tmp/app.spdx.json
2.0.88	Microsoft.SBOMTool	Coordinated Packages                 	229170
2.0.88	Microsoft.SBOMTool	chalk                                	2.4.2
2.0.88	Microsoft.SBOMTool	async-settle                         	1.0.0
```

## Using containerized sbomgr

```sh
$docker run [volume-maps] ghcr.io/interlynk-io/sbomgr [command] [options]
```

Example

```sh
$docker run -v ~/interlynk/sbomlc/:/app/sbomlc ghcr.io/interlynk-io/sbomgr packages -c /app/sbomlc
```

```
Unable to find image 'ghcr.io/interlynk-io/sbomgr:latest' locally
latest: Pulling from interlynk-io/sbomgr
479c7812d0ff: Already exists
5b3064dc8fe2: Already exists
Digest: sha256:d359b7e6e2b870542500dc00967ca2c5a4e78c8f1658b5c6dbdc8330effe38f8
Status: Downloaded newer image for ghcr.io/interlynk-io/sbomgr:latest

A new version of sbomgr is available v0.0.6.

Matching file count: 3153
Matching package count: 716953
```

# Search flags

## Packages

This section explains the flags relevant to the packages search feature.
The packages search takes only a single argument, either a file or a directory. There are man flags which can be specified to control its behaviour.

## _Match Criteria_

---

- `-N` or `--name` used for package/component name search.
- `-C` or `--cpe` used for package/component cpe search.
- `-P` or `--purl` used for pacakge/component purl search.
- `-H` or `--checksum` used for package/component checksum value search.

all of these match criteria are exclusive to each other.

## _Patter Matching_

---

- `-E` or `--extended-regexp` flag can be used to indicate if the match criteria is a regular expression. Syntax supported is https://github.com/google/re2/wiki/Syntax.

## _Matching Control_

---

- `-i` or `--ignore-case` case insensitive matching.

## _Output Control_

---

- `-l` or `--license` this includes the license of the package/component in the output.
- `-q` or `--quiet` this suppresses all output of the tool, the return value of the tool is 0 indicating success, if it finds the search criteria.
- `--no-filename` removes the filename from the output.
- `-j` or `--jsonl` outputs the search results in [jsonl](https://jsonlines.org/).
- `-p` or `--print-errors` includes errors encoundered during searching. Default is to ignore them.
- `-O` or `--output-format` user-defined output format. Options are listed below
  - `filen` - filepath
  - `tooln` - tool with which sbom was generated, only prints the first one
  - `toolv` - tool version
  - `docn` - sbom document name
  - `docv` - sbom document version
  - `cpe` - package cpe, only prints the first one, indicates how many cpe's exists.
  - `purl` - package purl
  - `pkgn` - package name
  - `pkgv` - package version
  - `pkgl` - package licenses
  - `specn` - spec of the sbom document, spdx or cdx.
  - `chkn` - checksum name
  - `chkv` - checksum value
  - `repo` - repository url
  - `direct` - package is a direct dependency

## _Stats Control_

---

- `-c` or `--count` suppresses the normal output and print matching counts of sbom filenames and packages.

## _Directory Control_

---

- `-r` or `--recurse` when set, recursively scans all sub directories.

## _Spec Control_

---

- `--spdx` searches only files which are SPDX.
- `--cdx` searches only files which are CycloneDX.

# Future work

- Search using files.
- Search using tool metadata.
- Search using CVE-ID.
- Search only direct dependencies.
- Search until a specified depth.
- Provide a list of malicious packages

# SBOM Samples

- A sample set of SBOM is present in the [samples](https://github.com/interlynk-io/sbomgr/tree/main/samples) directory above.
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

# Other SBOM Open Source tools

- [SBOM Assembler](https://github.com/interlynk-io/sbomasm) - A tool for conditional edits and merging of SBOMs
- [SBOM Seamless Transfer](https://github.com/interlynk-io/sbommv) - A primary tool to transfer SBOM's between different systems.
- [SBOM Quality Score](https://github.com/interlynk-io/sbomqs) - A tool for evaluating the quality and compliance of SBOMs
- [SBOM Explorer](https://github.com/interlynk-io/sbomex) - A tool for discovering and downloading SBOM from a public SBOM repository
- [SBOM Benchmark](https://www.sbombenchmark.dev) is a repository of SBOM and quality score for most popular containers and repositories

# Contact

We appreciate all feedback. The best ways to get in touch with us:

- :phone: [Live Chat](https://www.interlynk.io/#hs-chat-open)
- 📫 [Email Us](mailto:hello@interlynk.io)
- 🐛 [Report a bug or enhancement](https://github.com/interlynk-io/sbomex/issues)
- :x: [Follow us on X](https://twitter.com/InterlynkIo)

# Stargazers

If you like this project, please support us by starring it.

[![Stargazers](https://starchart.cc/interlynk-io/sbomgr.svg)](https://starchart.cc/interlynk-io/sbomgr)
