# Slideway

[![Build Status](https://travis.oncue.verizon.net/iptv/slideway.svg?token=Lp2ZVD96vfT8T599xRfV)](https://travis.oncue.verizon.net/iptv/slideway)

Slideway provides a small, native binary that creates Github releases and the associated metadata needed for the [Nelson](https://github.oncue.verizon.net/pages/iptv/nelson) deployment system.

## Development

1. `brew install go` - install the Go programming language:
1. `go get github.com/constabulary/gb/...` - install the `gb` build tool
1. `go get github.com/codeskyblue/fswatch` - install `fswatch` so we can do continous compilation
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`

