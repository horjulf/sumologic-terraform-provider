.DEFAULT_GOAL := build
GOOS          := $(shell go env GOOS)
GOARCH        := $(shell go env GOARCH)
PLUGIN_PATH   := ${HOME}/.terraform.d/plugins/${GOOS}_${GOARCH}
PLUGIN_NAME   := terraform-provider-sumologic
DIST_PATH     := dist
BUILD_PATH    := ${DIST_PATH}/${GOOS}_${GOARCH}
GO_PACKAGES   := $(shell go list -mod vendor ./...)
GO_FILES      := $(shell find . -name vendor -prune -or -type f -name '*.go' -print)
VERSION_PATH  := VERSION
VERSION       := $(shell cat ${VERSION_PATH})

.PHONY: all
all: test build

.PHONY: test
test: test-all

.PHONY: test-all
test-all:
	@TF_ACC=1 go test -mod vendor -v -race ${GO_PACKAGES}

${BUILD_PATH}/${PLUGIN_NAME}: go.sum ${GO_FILES}
	mkdir -p ${BUILD_PATH}; \
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod vendor \
		-ldflags "-s -w" \
		-o ${BUILD_PATH}/${PLUGIN_NAME} \
		.

.PHONY: build
build: ${BUILD_PATH}/${PLUGIN_NAME}

.PHONY: pack
pack: build
	# Compress
	upx -q ${BUILD_PATH}/${PLUGIN_NAME}
	# Test
	upx -t ${BUILD_PATH}/${PLUGIN_NAME}

.PHONY: install
install: build
	mkdir -p ${PLUGIN_PATH}; \
	rm -f ${PLUGIN_PATH}/${PLUGIN_NAME}_* ${PLUGIN_PATH}/${PLUGIN_NAME}; \
	install -m 0755 ${BUILD_PATH}/${PLUGIN_NAME} ${PLUGIN_PATH}/${PLUGIN_NAME}_v${VERSION}

.PHONY: clean
clean:
	rm -rf ${DIST_PATH}

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: vendor_update
vendor_update:
	go get -u ./...
	${MAKE} vendor
