# Slipway

[![Build Status](https://travis-ci.org/Verizon/slipway.svg?branch=master)](https://travis-ci.org/Verizon/slipway)
[![Latest Release](https://img.shields.io/github/release/verizon/slipway.svg)](https://github.com/Verizon/slipway/releases)

Slipway provides a small, native binary that creates Github releases and the associated metadata needed for the [Nelson](https://github.com/verizon/nelson) deployment system.

## Instalation

If you just want to use `slipway`, then run the following:

```
curl -GqL https://raw.githubusercontent.com/Verizon/slipway/master/scripts/install | bash
```

This script will download and install the latest version and put it on your `$PATH`. We do not endorse piping scripts from the wire to `bash`, and you should read the script before executing the command. It will:

1. Fetch the latest version from Github
2. Verify the SHA1 sum
3. Extract the tarball
4. Copy `slipway` to `/usr/local/bin/slipway`

It is safe to rerun this script to keep `slipway` current. If you have the source code checked out locally, you need only execute: `scripts/install` to install the latest version of `slipway`.

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
slipway release -x github.yourcompany.com -t 2.0.0 -d `pwd`/target

# specify the github domain, tag, repo slug and input directory
slipway release -x github.yourcompany.com -t 2.0.0 -r tim/sbt-release-sandbox -d `pwd`/target

```

## Using with Travis

As **Slipway** is a native binary, instalation is super simple and using it in Travis fits in with the regular bash-style script execution used by Travis. Here's an example `.travis.yml` (be sure to update the link to the latest version):

```
install:
  - curl -GqL https://raw.githubusercontent.com/Verizon/slipway/master/scripts/install | bash

script:
  - // do your build stuff
  # This assumes you are output docker images for internal consumption,
  # but essentially do whatever you need to in order to generate deployables
  # for each container you want to output from this repository
  - docker images | grep docker.yourcompany.com | awk '{print $1 ":" $2}' | slipway gen
  - |
    if [ $TRAVIS_PULL_REQUEST = 'false' ]; then
      git tag $RELEASE_VERSION &&
      git push --tags origin &&
      slipway release -x github.yourcompany.com -t $RELEASE_VERSION -d `pwd`
    fi

env:
  global:
    - RELEASE_VERSION="0.1.$TRAVIS_BUILD_NUMBER"

```

That's all there is to it.

## Development

1. `brew install go` - install the Go programming language:
1. `make devel`
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `make watch`

This should give continous compilation without the tedious need to constantly restart `gb build`

