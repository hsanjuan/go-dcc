all: build

deps:
	go get -u github.com/stianeikeland/go-rpio
build: deps
	go build ./...
test: deps
	go test -v ./...
coverage:
	./coverage.sh
install: deps
	echo "Installing dccpi"
	go install ./dccpi

.PHONY=all deps build test install coverage
