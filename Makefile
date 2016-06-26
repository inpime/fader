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
	go test -v -bench=. -benchmem -run=. ./src/api/...
testx:
	go test -v -run=TestImporter_simple ./src/...

prebuild:
	go get github.com/Sirupsen/logrus
	go get github.com/stretchr/testify
	
	go get github.com/boltdb/bolt
	go get github.com/inpime/dbox
	go get gopkg.in/olivere/elastic.v3

	go get github.com/labstack/echo
	go get github.com/yosssi/boltstore/reaper
	go get github.com/yosssi/boltstore/store

	go get github.com/flosch/pongo2
	go get github.com/flosch/pongo2-addons

	go get github.com/gebv/echo-session
	go get github.com/gorilla/sessions

	go get github.com/BurntSushi/toml

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