
#-------------------
# Variables
#-------------------

PROGRAM_NAME=slipway

TRAVIS_BUILD_NUMBER ?= dev
CLI_FEATURE_VERSION ?= 1.0
CLI_VERSION ?= ${CLI_FEATURE_VERSION}.${TRAVIS_BUILD_NUMBER}
# if not set, then we're doing local development
# as this will be set by the travis matrix for realz
TARGET_PLATFORM ?= darwin
TARGET_ARCH ?= amd64
TAR_NAME = ${PROGRAM_NAME}-${TARGET_PLATFORM}-${TARGET_ARCH}-${CLI_VERSION}.tar.gz

SHELL 	:= /bin/bash
BINDIR	:= bin
PKG 		:= github.com/envoyproxy/go-control-plane
GOFILES	 = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GODIRS	 = $(shell go list -f '{{.Dir}}' ./... \
						| grep -vFf <(go list -f '{{.Dir}}' ./vendor/...))

release: format test package

compile: format
	GOOS=${TARGET_PLATFORM} GOARCH=amd64 CGO_ENABLED=0 gb build -ldflags "-X main.globalBuildVersion=${CLI_VERSION}"

watch:
	fswatch

# test: compile
# 	gb test -v

package: test
	mkdir -p target && \
	mv bin/${PROGRAM_NAME}-${TARGET_PLATFORM}-amd64 ./${PROGRAM_NAME} && \
	tar -zcvf ${TAR_NAME} ${PROGRAM_NAME} && \
	rm ${PROGRAM_NAME} && \
	sha1sum ${TAR_NAME} > ${TAR_NAME}.sha1 && \
	shasum -c ${TAR_NAME}.sha1 && \
	mv ${TAR_NAME} target/${TAR_NAME} && \
	mv ${TAR_NAME}.sha1 target/${TAR_NAME}.sha1

format:
	go fmt src/github.com/getnelson/${PROGRAM_NAME}/*.go

tar:
	echo ${TAR_NAME}

.PHONY: build
build: vendor
	@echo "--> building"
	@go build ./...

.PHONY: clean
clean:
	@echo "--> cleaning compiled objects and binaries"
	@rm -rf $(BINDIR)/*

.PHONY: test
test: vendor
	@echo "--> running unit tests"
	@go test ./pkg/...

.PHONY: lint
lint: tools.golint
	@echo "--> checking code style with 'golint' tool"
	@echo $(GODIRS) | xargs -n 1 golint

#-------------------
#-- code generaion
#-------------------

generate:
	@echo "--> generating pb.go files"
	$(SHELL) build/generate_protos.sh

#------------------
#-- dependencies
#------------------
.PHONY: depend.update depend.install

depend.update: tools.glide
	@echo "--> updating dependencies from glide.yaml"
	@glide update

depend.install: tools.glide
	@echo "--> installing dependencies from glide.lock "
	@glide install

vendor:
	@echo "--> installing dependencies from glide.lock "
	@glide install

#-------------------
#-- tools
#-------------------

tools: tools.gb tools.glide tools.golint tools.fswatch

tools.gb:
	@command -v gb >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing gb"; \
		go get -u go get github.com/constabulary/gb/...; \
	fi

tools.fswatch:
	@command -v fswatch >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing fswatch"; \
		go get -u github.com/codeskyblue/fswatch; \
	fi

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get -u golang.org/x/lint/golint; \
	fi

tools.glide:
	@command -v glide >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing glide"; \
		curl https://glide.sh/get | sh; \
	fi
