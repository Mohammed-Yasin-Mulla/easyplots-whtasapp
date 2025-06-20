.PHONY: run build clean test fmt vet

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Install dependencies
deps:
	go mod tidy
	go mod download

# Help
help:
	@echo "Available commands:"
	@echo "  run    - Run the application"
	@echo "  build  - Build the application"
	@echo "  clean  - Clean build artifacts"
	@echo "  test   - Run tests"
	@echo "  fmt    - Format code"
	@echo "  vet    - Vet code"
	@echo "  check  - Run all checks (fmt, vet, test)"
	@echo "  deps   - Install and tidy dependencies" 