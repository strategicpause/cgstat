ROOT := $(shell pwd)

build:
	@echo "Building project"
	go build -o bin/cgstat main.go

fmt:
	@echo "Running go fmt"
	go fmt

lint:
	@if [ ! -d /tmp/golangci-lint ]; then \
		echo "Installing golangci-lint"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.45.2; \
		mkdir -p /tmp/golangci-lint/; \
		mv ./bin/golangci-lint /tmp/golangci-lint/golangci-lint; \
	fi; \
	/tmp/golangci-lint/golangci-lint run ./... --issues-exit-code=1 \

tidy:
	"Running go mod tidy"
	go mod tidy

release: fmt lint tidy build