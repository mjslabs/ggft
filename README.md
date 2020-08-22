# Go Generate from Template (ggft)

[![Build Status][travis-badge]][travis]
[![Go Report Card][goreport-badge]][goreport]
[![Test Coverage][coverage]][codeclimate]
[![FOSSA Status][fossa-badge]][fossa]

[travis-badge]: https://travis-ci.org/mjslabs/ggft.svg?branch=master
[travis]: https://travis-ci.org/mjslabs/ggft
[goreport-badge]: https://goreportcard.com/badge/github.com/mjslabs/ggft
[goreport]: https://goreportcard.com/report/github.com/mjslabs/ggft
[coverage]: https://api.codeclimate.com/v1/badges/9ed016c632217a454cac/test_coverage
[codeclimate]: https://codeclimate.com/github/mjslabs/ggft/test_coverage
[fossa-badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmjslabs%2Fggft.svg?type=shield
[fossa]: https://app.fossa.com/projects/git%2Bgithub.com%2Fmjslabs%2Fggft?ref=badge_shield

ggft is a tool for managing templates and generating new projects based on them.
Commands are available to download, list, and delete templates. This leads to
less typing of boilerplate or IDE magic when starting a new project, and helping
to enable consistency on teams, or if you need to manage and use git repos on
many computers.

## Installation

To install ggft, the quickest way is to use the standard `go get` process:

```bash
go get github.com/mjslabs/ggft
```

Alternatively, if you'd like to use custom build directives, add features, or
install a versioned binary, use the clone and make method:

```bash
git clone https://github.com/mjslabs/ggft.git && cd ggft
make install
```

## Usage

A common pattern to use with ggft is to download a template once, and then
deploy it many times. For example:

```bash
ggft get https://github.com/mjslabs/ggft-microservice ms
ggft new ms my-new-project
ggft new ms my-other-new-project
```

This would download the `ggft-microservice` repo to a local cache directory
(defaults to `~/.ggft/templates`), and then creates two new directories using
the `ggft-microservice` directory as a starting point. Assuming the
`ggft-microservice` repo uses
[Go template syntax](https://golang.org/pkg/text/template/) in its files, the
user will be prompted for values for those template variables.

### Advanced

`ggft new` has a few flags that allow you to customize how files are processed.
There are options to treat files as regular (non-template) files, trimming file
suffixes, and more. See `ggft new -h` for more info.

## Creating a template repo

ggft template repos are simply git repos that have files containing Go template
syntax and special file suffixes. Here's an example with inline comments:

```text
example
├── Makefile       <-- Makefile specific to this project
├── README.md      <-- README for the template itself. Will be ignored on parse
├── README.md.ggft <-- Will be turned into README.md when parsed by ggft
└── myproject.c    <-- A file using Go template syntax to be parsed by ggft
```

You can see other than README, we simply have a repo of files that we want to
use to seed our new project directory when we run `ggft`. For more info on how
the README parsing works, see the help documentation on the `-t` and `-i` flags
of `ggft new`.

## Caveats

ggft will try to initialize a config structure in `~/.ggft` and use this for all
template storage. This is not yet configurable.  

This is a brand new project, and as such is missing many features and not yet
heavily tested. Contributions welcome!
