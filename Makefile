
.PHONY:	 build test

build:
	go build -o .out/clilocker

test:
	go test ./...

