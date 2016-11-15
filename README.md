# Slipway

[![Build Status](https://travis.oncue.verizon.net/iptv/slipway.svg?token=Lp2ZVD96vfT8T599xRfV)](https://travis.oncue.verizon.net/iptv/slipway)

Slipway provides a small, native binary that creates Github releases and the associated metadata needed for the [Nelson](https://github.oncue.verizon.net/pages/iptv/nelson) deployment system.

## Usage

Generate metadata for a given container:

```
slipway deployable your.docker.com/foo/bar:1.2.3

# optionally specify an output directory
slipway deployable -d /path/to/dir your.docker.com/foo/bar:1.2.3
```

Cut a release with a set of deployables:

```
# specify the github domain
slipway release -x github.oncue.verizon.net --auth $GITHUB_TOKEN

# optionally specify an input directory
slipway release -d /path/to/dir
```

## Development

1. `brew install go` - install the Go programming language:
1. `go get github.com/constabulary/gb/...` - install the `gb` build tool
1. `go get github.com/codeskyblue/fswatch` - install `fswatch` so we can do continous compilation
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`

