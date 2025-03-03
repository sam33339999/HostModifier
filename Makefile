# Makefile for HostModifier

# Define the Go module path
MODULE_PATH := github.com/sam33339999/HostModifier

# Define the binary name
BINARY_NAME := hostmodifier

# Define the build target
.PHONY: build
build:
	go build -o $(BINARY_NAME) .

# Define the run target
.PHONY: run
run: build
	./$(BINARY_NAME)

# Define the deps target
.PHONY: deps
deps:
	go mod tidy
