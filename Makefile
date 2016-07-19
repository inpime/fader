
# https://github.com/travis-ci/worker/blob/master/Makefile#L8
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M:%S%z`
BUILD_HASH = `git rev-parse HEAD 2>/dev/null || echo "???"`

ifeq ($(BUILD_HASH),)
	BUILD_HASH := ${TRAVIS_COMMIT}
endif
VERSION = v0.1.1
# BUILD_NUMBER = ${TRAVIS_BUILD_NUMBER}
platforms = linux
GOPATH := ${GOPATH}

BUILDFLAGS ?= -a --installsuffix cgo -ldflags \
    "-s -X 'github.com/inpime/fader/config.BuildDate=$(BUILD_DATE)'\
    -X github.com/inpime/fader/config.BuildHash=$(BUILD_HASH)\
    -X github.com/inpime/fader/config.Version=$(VERSION)" \
    fader.go

test: vendor vendor-test
	#docker run --rm -v "${PWD}":/usr/src/github.com/inpime/fader \
	#	-w /usr/src/github.com/inpime/fader \
	#	golang:1.6 go test -v -bench=. -benchmem -run=. ./...
	go test -v -race -bench=. -benchmem -run=. ./...
.PHONY: test

build-via-docker:
	@echo Docker build Linux amd64
	docker run --rm -v "${PWD}":/usr/src/github.com/inpime/fader \
		-w /usr/src/github.com/inpime/fader \
		golang:1.6 make build
.PHONY: build-via-docker

build: vendor
	@echo Build
	@$(foreach GOOS,$(platforms), $(foreach GOARCH,amd64, \
		env GOPATH=${GOPATH}:/usr \
				CGO_ENABLED=0 \
				go build -o build/${GOOS}_${GOARCH}/fader -v $(BUILDFLAGS) \
	;))
.PHONY: build

build-withrace: vendor
	@echo Build
	@$(foreach GOOS,$(platforms), $(foreach GOARCH,amd64, \
		env GOPATH=${GOPATH}:/usr \
				CGO_ENABLED=1 \
				go build -o build/${GOOS}_${GOARCH}/fader -v -race $(BUILDFLAGS) \
	;))
.PHONY: build-withrace

release: build
	@mkdir -p releases
	@$(foreach GOOS,$(platforms), $(foreach GOARCH,amd64, \
		zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).${GOOS}_${GOARCH}.$(BUILD_DATE).zip ${PWD}/build/${GOOS}_${GOARCH}/fader \
	;))
.PHONY: release

vendor-test:
	go get -v github.com/stretchr/testify/assert
.PHONY: vendor-test

vendor:
	@echo Installation vendor
	go get -v github.com/BurntSushi/toml \
		github.com/Sirupsen/logrus \
		github.com/boltdb/bolt \
		github.com/flosch/pongo2 \
		github.com/flosch/pongo2-addons \
		github.com/echo-contrib/sessions \
		github.com/gorilla/sessions \
		github.com/inpime/dbox \
		github.com/labstack/echo \
		github.com/labstack/echo/engine/standard \
		github.com/labstack/echo/middleware \
		github.com/lionelbarrow/braintree-go \
		github.com/mailgun/mailgun-go \
		github.com/yosssi/boltstore/reaper \
		github.com/yosssi/boltstore/store \
		gopkg.in/go-playground/validator.v8 \
		gopkg.in/olivere/elastic.v3
.PHONY: vendor

# for a comfortable development
reload-dev: build-via-docker
# reload-dev:
	# env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o build/linux_amd64/fader -v -race $(BUILDFLAGS)
	
	cp ${PWD}/build/linux_amd64/fader ${PWD}/state/app/fader
	-docker-compose down
	docker-compose build --force-rm
	docker-compose up -d
.PHONY: reload-dev