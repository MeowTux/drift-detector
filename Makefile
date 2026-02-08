# Drift Detector Makefile
# Author: MeowTux

.PHONY: help build build-all install test clean run docker-build docker-run

# Variables
BINARY_NAME=drift-detector
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildDate=${BUILD_TIME}"

# Default target
help:
	@echo "Drift Detector - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-all      Build for all platforms"
	@echo "  make install        Install binary to /usr/local/bin"
	@echo "  make test           Run tests"
	@echo "  make clean          Remove build artifacts"
	@echo "  make run            Build and run"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-run     Run in Docker"
	@echo ""

# Build for current platform
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .
	@echo "Build complete: ./${BINARY_NAME}"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64
	
	# Linux ARM64 (for Raspberry Pi, Android)
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-arm64
	
	# Linux ARM (32-bit)
	GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BINARY_NAME}-linux-arm
	
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-amd64
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-windows-amd64.exe
	
	@echo "All builds complete!"
	@ls -lh ${BINARY_NAME}-*

# Install to system
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	sudo cp ${BINARY_NAME} /usr/local/bin/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "Tests complete!"

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}-*
	rm -f coverage.txt coverage.html
	@echo "Clean complete!"

# Build and run
run: build
	./${BINARY_NAME}

# Initialize config
init: build
	./${BINARY_NAME} init

# Run drift detection
detect: build
	./${BINARY_NAME} detect

# Watch mode
watch: build
	./${BINARY_NAME} detect --watch --interval 5m

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t meowtux/drift-detector:${VERSION} -t meowtux/drift-detector:latest .
	@echo "Docker image built!"

# Docker run
docker-run:
	docker run --rm \
		-v $(PWD)/config:/app/config \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL} \
		meowtux/drift-detector:latest detect

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Update dependencies
deps:
	go mod tidy
	go mod download

# Security scan
security:
	gosec ./...

.DEFAULT_GOAL := help
