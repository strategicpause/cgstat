default: release
build:
	@echo "Building project"
	GOOS=linux go build -o bin/cgstat main.go

fmt:
	@echo "Running go fmt"
	go fmt

lint:
	@if [ ! -d /tmp/golangci-lint ]; then \
		echo "Installing golangci-lint"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.50.1; \
		mkdir -p /tmp/golangci-lint/; \
		mv ./bin/golangci-lint /tmp/golangci-lint/golangci-lint; \
	fi; \
	/tmp/golangci-lint/golangci-lint run ./... --issues-exit-code=1 \

tidy:
	echo "Running go mod tidy"
	go mod tidy

vulncheck:
	@if [ ! `command -v govulncheck` ]; then \
  		echo "Installing govulncheck"; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi; \
	GOOS=linux govulncheck ./...

install: build
	cp ./bin/cgstat /usr/bin/local

release: fmt lint tidy vulncheck build