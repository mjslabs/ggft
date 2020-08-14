# Go Generate from Template (ggft)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmjslabs%2Fggft.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmjslabs%2Fggft?ref=badge_shield)


ggft is a tool for managing templates and generating new projects based on them.
Commands are available to download, list, and delete templates. This leads to
less typing of boilerplate or IDE magic when starting a new project, and helping
to enable consistency on teams, or if you need to manage and use git repos on
many computers.

## Installation & usage

To install ggft, use the standard `go install` process.

```bash
git clone https://github.com/mjslabs/ggft && cd ggft
go install
```

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

## Caveats

ggft will try to initialize a config structure in `~/.ggft` and use this for all
template storage. This is not yet configurable.  

This is a brand new project, and as such is missing many features and not yet
heavily tested. Contributions welcome!


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmjslabs%2Fggft.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmjslabs%2Fggft?ref=badge_large)