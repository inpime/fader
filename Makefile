
# https://github.com/travis-ci/worker/blob/master/Makefile#L8
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M:%S%z`
BUILD_HASH = `git rev-parse HEAD 2>/dev/null || echo "???"`

ifeq ($(BUILD_HASH),)
	BUILD_HASH := ${TRAVIS_COMMIT}
endif
VERSION = v0.1.1
# BUILD_NUMBER = ${TRAVIS_BUILD_NUMBER}

GO=go
GOPATH:=${GOPATH}:/Users/gbv/work
BUILDFLAGS ?= -a --installsuffix cgo -ldflags \
    "-s -X 'github.com/inpime/fader/config.BuildDate=$(BUILD_DATE)'\
    -X github.com/inpime/fader/config.BuildHash=$(BUILD_HASH)\
    -X github.com/inpime/fader/config.Version=$(VERSION)" \
    fader.go

build: build-linux
.PHONY: build

test: prebuild
	$(GO) test -v -bench=. -benchmem -run=TestAppendOrReplace ./...
.PHONY: test

build-linux: prebuild
	@echo Build Linux amd64
	env GOBIN=${PWD}/build/linux_amd64 \
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
        $(GO) install $(BUILDFLAGS)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.$(BUILD_DATE).zip ${PWD}/build/linux_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.latest.zip ${PWD}/build/linux_amd64/fader
.PHONY: build-linux

build-osx: prebuild
	@echo Build OSX amd64
	env GOBIN=${PWD}/build/osx_amd64 \
        GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
        $(GO) install $(BUILDFLAGS)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.$(BUILD_DATE).zip ${PWD}/build/osx_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.latest.zip ${PWD}/build/osx_amd64/fader
.PHONY: build-osx


prebuild:
	@echo Pre install
	mkdir -p releases
	-$(GO) get github.com/inpime/fader
.PHONY: prebuild

# for a comfortable development
build-linux-dev: build-linux
	cp ${PWD}/build/linux_amd64/fader ${PWD}/state/app/fader
.PHONY: build-linux-dev

reload-dev: build-linux-dev
	-docker-compose down
	docker-compose build --force-rm
	docker-compose up
