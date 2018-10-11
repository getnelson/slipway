# Slipway

[![Build Status](https://travis-ci.org/getnelson/slipway.svg?branch=master)](https://travis-ci.org/getnelson/slipway)
[![Latest Release](https://img.shields.io/github/release/getnelson/slipway.svg)](https://github.com/getnelson/slipway/releases)

Slipway provides a small, native binary that creates Github releases and the associated metadata needed for the [Nelson](https://github.com/getnelson/nelson) deployment system.

## Instalation

If you just want to use `slipway`, then run the following:

```
curl -GqL https://raw.githubusercontent.com/getnelson/slipway/master/scripts/install | bash
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

# specify the github domain, tag, and a properties file for github credentials.
# the repo slug will automatically be read from TRAVIS_REPO_SLUG
slipway release -x github.yourcompany.com -t 2.0.0 -c /path/to/credentials

# specify the github domain, tag, repo slug and input directory
slipway release -x github.yourcompany.com -t 2.0.0 -r getnelson/howdy -d `pwd`/target

```

Create a [Github deployment](https://developer.github.com/v3/repos/deployments/), encoding any `*.deployable.yml` files into the `payload` field (as directed by Github API):

```
# specify the branch to use for the deployment (repo is infered from `TRAVIS_REPO_SLUG`):
slipway deploy --ref master

# specify the tag to use for the deployment (repo is infered from `TRAVIS_REPO_SLUG`):
slipway deploy -t 2.0.0

# specify the repo and an exact SHA to use for the deployment:
slipway deploy --ref fdb7da2ab3b2cd172e86c1af9adefa3523f6d65b  -r getnelson/howdy

# specify the github domain, tag, repo slug and input directory
slipway deploy -x github.yourcompany.com -t 2.0.0 -r getnelson/howdy -d `pwd`/target
```

### Authentication

In the event you wish to read Github credentials for **Slipway** from a file, the format of that file must be something like this:

```
github.login=username
github.token=XXXXXXXXXXXXXXXXXXXXXXXXX
```

This is a [classic Java properties format](https://www.mkyong.com/java/java-properties-file-examples/) - essentialy key=value pairs delimited by the equals sign. Whilst this functionality is supported, the authors recommend that credentials are instead retrieved from the shell environment, instead of being persisted to a plaintext file.

## Using with Travis

As **Slipway** is a native binary, instalation is super simple and using it in Travis fits in with the regular bash-style script execution used by Travis. Here's an example `.travis.yml` (be sure to update the link to the latest version):

```
install:
  - curl -GqL https://raw.githubusercontent.com/getnelson/slipway/master/scripts/install | bash

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

Be sure that you have [Go installed](https://golang.org/doc/install); this can be achieved on OSX with Homebew:

```
brew install go
```

Next, checkout the Slipway repository into your `$GOPATH`, with something like this:

```
cd $GOPATH && \
git clone git@github.com:getnelson/slipway.git github.com/getnelson/slipway
```

Next, install the tools `slipway` needs to build:

```
make tools
make deps.install
```

Generate the protobuf (as a one-time operation) and then build:

```
make generate && \
make build && \
make test
```

The slipway binary will then be available in `./bin/slipway`


