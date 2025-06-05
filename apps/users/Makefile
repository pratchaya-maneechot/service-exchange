# Makefile
.PHONY: all generate build run clean

all: generate build dev start

generate:
	@echo "Running generate scripts..."
	@sh ./scripts/generate-proto.sh
	@sh ./scripts/generate-wire.sh

build:
	@echo "Building the application..."
	@go build -o bin/main ./cmd/main.go

dev:
	@go run cmd/main.go

start:
	@echo "Running the application..."
	@./bin/main

clean:
	@echo "Cleaning up..."
	@rm -rf bin/