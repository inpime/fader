
vendor:
	go get -v github.com/labstack/echo/... \
		github.com/stretchr/testify/assert \
		github.com/tylerb/graceful \
		github.com/satori/go.uuid \
		github.com/flosch/pongo2 \
		github.com/BurntSushi/toml
.PHONY: vendor 
test:
#	rm -rf _app.db
	GOPATH=${PWD} go test -v \
		-bench=. -benchmem \
		-run=TestAddonsBasic_ContextCurrentFile ./src/api/...
	# GOPATH=${PWD} go test -v \
		-run= ./src/store/...
.PHONY: test

gen:
	-rm -rf ${FADER_DBPATH}
	-rm -rf ${FADER_INITFILE}
	GOPATH=${PWD} go run src/cmd/gen/main.go

run:
	GOOS="darwin" \
	GOARCH="amd64" \
	GOPATH=${PWD}  \
	go run src/cmd/httpserver/main.go 
.PHONY: run-dev 