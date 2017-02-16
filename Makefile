
TRAVIS_BUILD_NUMBER ?= 9999
NEXUS_CREDENTIALS ?= "anon:precioussecrets"

BINARY_NAME=slipway
BINARY_FEATURE_VERSION=0.2
BINARY_VERSION=${BINARY_FEATURE_VERSION}.${TRAVIS_BUILD_NUMBER}
IMAGE_TAG=docker.oncue.verizon.net/units/consort-${BINARY_FEATURE_VERSION}:${BINARY_VERSION}
TGZ_NAME=${BINARY_NAME}-${TARGET_PLATFORM}-${TARGET_ARCH}-${BINARY_VERSION}.tgz
# if not set, then we're doing local development
# as this will be set by the travis matrix for realz
TARGET_PLATFORM ?= darwin
TARGET_ARCH ?= amd64

all: package

watch:
	fswatch

format:
	gofmt -l -w src/

# compile for linux, but this binary is going to immedietly be stuffed
# into an alpine linux image. if someone wants to build this thing for
# themselves, then they can simply do: `gb build`
compile: format
	GOOS=${TARGET_PLATFORM} GOARCH=${TARGET_ARCH} CGO_ENABLED=0 gb build -ldflags "-X main.globalBuildVersion=${BINARY_VERSION}"

test: compile
	gb test -v

package: test
	mkdir target && \
	mv bin/${BINARY_NAME}-${TARGET_PLATFORM}-amd64 ./${BINARY_NAME} && \
	tar -zcvf ${TGZ_NAME} ${BINARY_NAME} && \
	rm ${BINARY_NAME} && \
	mv ${TGZ_NAME} target/${TGZ_NAME}

clean:
	rm -rf bin && \
	rm -rf pkg && \
	rm -rf target

publish:
	tar -c -z -f ${TGZ_NAME} -C bin consort && \
	curl -v -u ${NEXUS_CREDENTIALS} --upload-file ${TGZ_NAME} http://nexus.oncue.verizon.net/nexus/content/repositories/releases/verizon/inf/${BINARY_NAME}/${BINARY_VERSION}/${TGZ_NAME}
