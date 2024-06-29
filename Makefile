.PHONY: all build test lint clean

BINARY_NAME=bin/askllm
GOFILES=$(shell find . -name "*.go")
VERSION := $(shell grep -Eo 'VERSION = "(.*)"' internal/config/config.go | sed -E 's/.*"(.*)".*/\1/')

all: test lint build

build: ## Build the binary
	$(info ******************** Build the binary ******************** $(VERSION))
	go build -o $(BINARY_NAME) ./cmd/askllm

run: ## Run the application
	$(info ******************** Run the application ******************** $(VERSION))
	go run ./cmd/askllm/main.go "hello, llm"
#	go run ./cmd/askllm/main.go -e chatgpt -m gpt-3.5-turbo "hello, llm"
#	go run ./cmd/askllm/main.go -e gemini -m gemini-1.5-pro "hello, llm"
#	go run ./cmd/askllm/main.go -e claude -m claude-3-sonnet-20240229 "hello, llm"

test: ## Run tests
	go test -v ./...

lint: ## Lint the code
	$(info ******************** Lint the code ******************** $(VERSION))
#	@errcheck -ignoretests ./...
	@errcheck ./...
	@go vet ./...
	@golangci-lint run ./...

fmt: ## Format all code
	$(info ******************** Format all code ******************** $(VERSION))
	@test -z $(shell gofmt -l $(GOFILES)) || (gofmt -d $(GOFILES); exit 1)

clean: ## Remove object files and binary
	$(info ******************** Remove object files and binary ******************** $(VERSION))
	go clean
	rm -f $(BINARY_NAME)

release: ## Build and release a new version
	$(info ******************** Build and release a new version ******************** $(VERSION))
	goreleaser release --snapshot --clean

check: ## Check if the release is valid
	$(info ******************** Check if the release is valid ******************** $(VERSION))
	goreleaser check

tag: ## Create a new tag
	$(info ******************** Create a new tag ******************** $(VERSION))
	git tag v$(VERSION) && git push origin v$(VERSION)

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)
