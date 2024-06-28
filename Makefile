.PHONY: all build test lint clean

BINARY_NAME=bin/askllm

all: test lint build

build: ## Build the binary
	go build -o $(BINARY_NAME) ./cmd/askllm

test: ## Run tests
	go test -v ./...

lint: ## Lint the code
	golangci-lint run
	errcheck ./...

clean: ## Remove object files and binary
	go clean
	rm -f $(BINARY_NAME)

release: ## Build and release a new version
	goreleaser release --snapshot --rm-dist

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)
