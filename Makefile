
vendor:
	curl https://glide.sh/get | sh
	cd src && glide up
.PHONY: vendor 
test:
	rm -rf _app.db
	GOPATH=${PWD} go test -v \
		-run= ./src/test/...
	GOPATH=${PWD} go test -v \
		-bench=. -benchmem \
		-run= ./src/api/...
	GOPATH=${PWD} go test -v \
		-run= ./src/interfaces/...
	GOPATH=${PWD} go test -v \
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