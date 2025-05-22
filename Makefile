PKG := ./...

.PHONY: all test lint fmt clean

all: test

# Run tests with race detector and coverage
test:
	@echo "Running tests..."
	go test -race -v -covermode=atomic ${PKG}

# Run golangci-lint excluding test files
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --fix --tests=false -v

# Format code: basic fmt, organize imports, strict formatting, grouped imports
fmt:
	golangci-lint fmt

# Clean temporary files and cache
clean:
	@echo "Cleaning..."
	go clean
