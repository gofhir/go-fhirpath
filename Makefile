# Makefile for your FHIRPath project

.PHONY: all test race coverage clean lint

# Default task
all: test

# Run tests with verbose output
test:
	go test -v ./...

# Run tests with race detection
race:
	go test -v -race ./...

# Run tests with coverage report and generate HTML
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "🔥 Coverage report generated at coverage.html"

# Run linter (optional)
lint:
	go vet $(shell go list ./... | grep -v "parser/grammar")

# Clean temporary files
clean:
	rm -f coverage.out coverage.html