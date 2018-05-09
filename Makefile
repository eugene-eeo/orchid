PHONY: build

install:
	go get -u github.com/gobuffalo/packr/...

test:
	go test ./...

build: test
	packr build
	packr install .
