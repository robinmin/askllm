.PHONY: all build test lint clean

BINARY_NAME=bin/askllm

all: test lint build

build:
	go build -o $(BINARY_NAME) ./cmd/askllm

test:
	go test -v ./...

lint:
	golangci-lint run
	errcheck ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

release:
	goreleaser release --snapshot --rm-dist