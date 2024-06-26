# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@tailwindcss -i input.css -o web/css/app.css 
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v - 1

cover:
	@echo "Testing..."
	@go test -p 1 -v  -covermode=count -coverprofile=cover.out ./... -coverpkg ./... 
	@go tool cover -html=cover.out

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main cover.*

# Live Reload
watch:
	@echo "Watching..."
	@air

.PHONY: all build run test clean
