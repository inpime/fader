build:
	gofmt -w ./src
	GOPATH=${GOPATH}:${PWD} go build -o bin/fader src/fader.go src/config.go src/fixtures.go 
run:
	GOPATH=${GOPATH}:${PWD} go run src/fader.go src/config.go src/fixtures.go -workspace="./_workspace/" -address="192.168.1.36:3322"
test:
	GOPATH=${GOPATH}:${PWD} go test -v -bench=. -benchmem -run=. ./src/...
testx:
	GOPATH=${GOPATH}:${PWD} go test -v -run=TestImporter_simple ./src/...
	