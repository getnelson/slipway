# Slipway

[![Build Status](https://travis.oncue.verizon.net/iptv/slipway.svg?token=Lp2ZVD96vfT8T599xRfV)](https://travis.oncue.verizon.net/iptv/slipway)

Slipway provides a small, native binary that creates Github releases and the associated metadata needed for the [Nelson](https://github.oncue.verizon.net/pages/iptv/nelson) deployment system. 

## Instalation

Download the latest release of slipway [from the internal nexus](http://nexus.oncue.verizon.net/nexus/content/groups/internal/verizon/inf/slipway/), and place the binary on your `$PATH`

## Usage

Generate metadata for a given container:

```
# will assume you want the .deployable.yml generated in the current working directory
slipway gen your.docker.com/foo/bar:1.2.3

# optionally specify an output directory
slipway gen -d /path/to/dir your.docker.com/foo/bar:1.2.3
```

Cut a release with an optional set of deployables (note, for use with Nelson, you *need* the `.deployable.yml` files):

```
# release a tag for a repository hosted on github.com
# this release has zero release assets
slipway release -t 2.0.0

# release a tag for a repository hosted on github.com
# read the release assets from `pwd`/target
slipway release -t 2.0.0 -d `pwd`/target

# specify the github domain, tag, and input directory.
# the repo slug will automatically be read from TRAVIS_REPO_SLUG
slipway release -x github.oncue.verizon.net -t 2.0.0 -d `pwd`/target

# specify the github domain, tag, repo slug and input directory
slipway release -x github.oncue.verizon.net -t 2.0.0 -r tim/sbt-release-sandbox -d `pwd`/target

```

## Development

1. `brew install go` - install the Go programming language:
1. `go get github.com/constabulary/gb/...` - install the `gb` build tool
1. `go get github.com/codeskyblue/fswatch` - install `fswatch` so we can do continous compilation
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`

