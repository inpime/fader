
_GOPATH=${GOPATH}:${PWD}

vendor:
	go get -v github.com/labstack/echo/... \
		github.com/stretchr/testify/assert \
		github.com/tylerb/graceful \
		github.com/satori/go.uuid \
		github.com/flosch/pongo2 \
		github.com/BurntSushi/toml
.PHONY: vendor 
test: vendor
	GOPATH=${_GOPATH} go test -v \
		-bench=. -benchmem \
		-run=. ./src/...
.PHONY: test

run-dev:
	GOOS="darwin" \
	GOARCH="amd64" \
	GOPATH=${_GOPATH}  \
	go run src/cmd/httpserver/main.go
.PHONY: run-dev 