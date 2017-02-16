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

## Using with Travis

As **Slipway** is a native binary, instalation is super simple and using it in Travis fits in with the regular bash-style script execution used by Travis. Here's an example `.travis.yml` (be sure to update the link to the latest version):

```
install:
  - wget http://nexus.oncue.verizon.net/nexus/content/repositories/releases/verizon/inf/slipway/0.1.7/slipway-linux-amd64-0.1.7.tar.gz -O slipway.tar.gz
  - tar xvf slipway.tar.gz
  - mv slipway $HOME/slipway

script:
  - // do your build stuff
  # This assumes you are output docker images for internal consumption,
  # but essentially do whatever you need to in order to generate deployables
  # for each container you want to output from this repository
  - docker images | grep docker.oncue.verizon.net | awk '{print $1 ":" $2}' | docker gen
  - |
    if [ $TRAVIS_PULL_REQUEST = 'false' ];
      git tag $RELEASE_VERSION &&
      git push --tags origin &&
      slipway release -x github.oncue.verizon.net -t $RELEASE_VERSION -d `pwd`
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

