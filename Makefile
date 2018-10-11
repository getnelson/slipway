
#-------------------
# Variables
#-------------------

PROGRAM_NAME=slipway

TRAVIS_BUILD_NUMBER ?= 999999
CLI_FEATURE_VERSION ?= 1.0
CLI_VERSION ?= ${CLI_FEATURE_VERSION}.${TRAVIS_BUILD_NUMBER}
# if not set, then we're doing local development
# as this will be set by the travis matrix for realz
TARGET_PLATFORM ?= darwin
TARGET_ARCH ?= amd64
BINARY_NAME := ${PROGRAM_NAME}-${TARGET_PLATFORM}-${TARGET_ARCH}-${CLI_VERSION}
TAR_NAME = ${BINARY_NAME}.tar.gz

SHELL 	:= /bin/bash
BINDIR	:= bin
PKG 		:= github.com/getnelson/${PROGRAM_NAME}
GOFILES	 = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GODIRS	 = $(shell go list -f '{{.Dir}}' ./... \
						| grep -vFf <(go list -f '{{.Dir}}' ./vendor/...))

.PHONY: release
release: format test package

.PHONY: package
package: test build
	mkdir -p target && \
	mv bin/${BINARY_NAME} ./${PROGRAM_NAME} && \
	tar -zcvf ${TAR_NAME} ${PROGRAM_NAME} && \
	rm ${PROGRAM_NAME} && \
	sha1sum ${TAR_NAME} > ${TAR_NAME}.sha1 && \
	shasum -c ${TAR_NAME}.sha1 && \
	mv ${TAR_NAME} target/${TAR_NAME} && \
	mv ${TAR_NAME}.sha1 target/${TAR_NAME}.sha1

.PHONY: build
build: vendor
	@echo "--> building"
	GOOS=${TARGET_PLATFORM} \
	GOARCH=${TARGET_ARCH} \
	CGO_ENABLED=0 \
	GOBIN=$(BINDIR) \
	go build \
	-v \
	-ldflags "-X main.globalBuildVersion=${CLI_VERSION}" \
	-o ${BINDIR}/${BINARY_NAME} \
	./cmd

.PHONY: watch
watch:
	@echo "--> watching for changed files"
	@fswatch

.PHONY: clean
clean:
	@echo "--> cleaning compiled objects and binaries"
	@rm -rf $(BINDIR)/*

.PHONY: test
test: vendor
	@echo "--> running unit tests"
	@go test ./cmd/...

.PHONY: format
format: tools.goimports
	@echo "--> formatting code with 'goimports' tool"
	@goimports -local $(PKG) -w -l $(GOFILES)

.PHONY: lint
lint: tools.golint
	@echo "--> checking code style with 'golint' tool"
	@echo $(GODIRS) | xargs -n 1 golint

#-------------------
#-- code generaion
#-------------------

generate: $(BINDIR)/gogofast $(BINDIR)/validate
	@echo "--> generating pb.go files"
	$(SHELL) scripts/generate-protos

#------------------
#-- dependencies
#------------------

.PHONY: deps.update deps.install

deps.update: tools.glide
	@echo "--> updating dependencies from glide.yaml"
	@glide update

deps.install: tools.glide
	@echo "--> installing dependencies from glide.lock "
	@glide install

vendor:
	@echo "--> installing dependencies from glide.lock "
	@glide install

#-------------------
#-- tools
#-------------------

tools: tools.glide tools.golint tools.fswatch tools.goimports

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

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing goimports"; \
		go get golang.org/x/tools/cmd/goimports; \
	fi

$(BINDIR)/gogofast: vendor
	@echo "--> building $@"
	@go build -o $@ vendor/github.com/gogo/protobuf/protoc-gen-gogofast/main.go

$(BINDIR)/validate: vendor
	@echo "--> building $@"
	@go build -o $@ vendor/github.com/lyft/protoc-gen-validate/main.go
