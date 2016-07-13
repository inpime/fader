

BUILD_STAMP = `date +%Y%m%d_%H%M%S`
BUILD_HASH = $(git rev-parse HEAD)
VERSION = dev

GO=go 
GOPATH:=${GOPATH}:/Users/gbv/work
BUILDFLAGS ?= -a --installsuffix cgo -ldflags \
    "-s -X 'github.com/inpime/fader/config.BuildDate=$(BUILD_STAMP)'\
    -X github.com/inpime/fader/config.BuildHash=$(BUILD_HASH)\
    -X github.com/inpime/fader/config.Version=$(VERSION)" \
    fader.go

build: build-linux
.PHONY: build

test:
	$(GO) test -v -bench=. -benchmem -run=TestAppendOrReplace ./...
.PHONY: test

build-linux: prebuild
	@echo Build Linux amd64
	env GOBIN=${PWD}/build/linux_amd64 \
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
        $(GO) install $(BUILDFLAGS)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.$(BUILD_STAMP).zip ${PWD}/build/linux_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.latest.zip ${PWD}/build/linux_amd64/fader
.PHONY: build-linux

build-osx: prebuild
	@echo Build OSX amd64
	env GOBIN=${PWD}/build/osx_amd64 \
        GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
        $(GO) install $(BUILDFLAGS)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.$(BUILD_STAMP).zip ${PWD}/build/osx_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.latest.zip ${PWD}/build/osx_amd64/fader
.PHONY: build-osx


prebuild:
	@echo Pre install
	mkdir -p releases 
	$(GO) get github.com/inpime/fader
.PHONY: prebuild

# for a comfortable development
build-linux-dev: build-linux
	cp ${PWD}/build/linux_amd64/fader ${PWD}/state/app/fader
.PHONY: build-linux-dev

reload-dev: build-linux-dev
	-docker-compose down
	docker-compose build --force-rm
	docker-compose up -d