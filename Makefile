test:
	go test ./...

build: test
	go build
	go install .
