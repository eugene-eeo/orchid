PHONY: build

install:
	go get -u github.com/gobuffalo/packr/...
	dep ensure
	make release

test:
	go test ./...

build: test
	packr build

release:
	packr -z
	go build -ldflags '-s -w'
	go install -ldflags '-s -w'
	packr clean
