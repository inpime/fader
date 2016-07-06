GO=go
ENTIRYPOINTFILES ?= src/fader.go src/config.go src/fixtures.go src/fixturex.go
GOPATH:=${GOPATH}:${PWD}

RUNARGS ?= -workspace="./_workspace/" \
	-address="192.168.1.36:3322" \
	-es_address="https://es.idheap.com" \
	-es_index="fader" \
	-session_secret="secure-key" \
	-mode="info"

STAMP = `date +%Y%m%d_%H%M%S`

build: build-linux
	 
run:
	$(GO) run $(ENTIRYPOINTFILES) $(RUNARGS)
test:
	$(GO) test -v -bench=. -benchmem -run=. ./src/api/...
testx:
	$(GO) test -v -run=TestAppendOrReplace ./src/...

prebuild:
	$(GO) get github.com/Sirupsen/logrus \
		github.com/stretchr/testify \
		github.com/boltdb/bolt \
		github.com/inpime/dbox \
		gopkg.in/olivere/elastic.v3 \
		github.com/labstack/echo/... \
		github.com/yosssi/boltstore/reaper \ 
		github.com/yosssi/boltstore/store \ 
		github.com/flosch/pongo2 \ 
		github.com/flosch/pongo2-addons \
		github.com/gebv/echo-session \ 
		github.com/gorilla/sessions \
		github.com/BurntSushi/toml \
		github.com/levigross/grequests \
		github.com/dyatlov/go-oembed/... \
		github.com/dyatlov/go-url2oembed/... \
		github.com/jeffail/tunny \ 
		github.com/lionelbarrow/braintree-go \
		gopkg.in/go-playground/validator.v8

	mkdir -p releases

.PHONY: prebuild
	
build-linux: prebuild
	@echo Build Linux amd64
	env GOBIN=${PWD}/build/linux_amd64 GOOS=linux GOARCH=amd64 $(GO) install $(ENTIRYPOINTFILES)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.$(STAMP).zip ${PWD}/build/linux_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).linux_amd64.latest.zip ${PWD}/build/linux_amd64/fader
.PHONY: build-linux

build-osx: prebuild
	@echo Build OSX amd64
	env GOBIN=${PWD}/build/osx_amd64 GOOS=darwin GOARCH=amd64 $(GO) install $(ENTIRYPOINTFILES)
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.$(STAMP).zip ${PWD}/build/osx_amd64/fader
	zip -j ${PWD}/releases/fader.go$(TRAVIS_GO_VERSION).osx_amd64.latest.zip ${PWD}/build/osx_amd64/fader
.PHONY: build-osx

build-linux-dev: build-linux
	cp ${PWD}/build/linux_amd64/fader ${PWD}/state/app/fader
.PHONY: build-linux-dev

# helper
reload-dev: build-linux-dev
	-docker-compose down
	docker-compose build
	docker-compose up -d

	