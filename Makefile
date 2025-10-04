.PHONY: test fmt lint run build clean

# Test targets
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run --skip-dirs=node_modules,vendor

# Run the application
run:
	go run main.go

# Build the application
build:
	go build -o bin/api-key-generator main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with race detection
test-race:
	go test -race -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
