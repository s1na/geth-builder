# Makefile for geth-builder

# Go parameters
GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOFMT := $(GO) fmt
GOVET := $(GO) vet
GOMOD := $(GO) mod

# Directories
BINDIR := bin
SRCDIR := .
CONFIGFILE := geth-builder-config.yaml

# Main go file
MAIN := main.go

# Output binary
BINARY_NAME := geth-builder

# Build the project
.PHONY: all
all: clean fmt vet build

# Format the code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Vet the code
.PHONY: vet
vet:
	$(GOVET) ./...

# Build the binary
.PHONY: build
build:
	$(GOBUILD) -o $(BINDIR)/$(BINARY_NAME) $(SRCDIR)/$(MAIN)

# Clean up build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINDIR)/$(BINARY_NAME)
	rm -rf go-ethereum

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run the tool with the default configuration file
.PHONY: run
run: build
	./$(BINDIR)/$(BINARY_NAME) -config=$(CONFIGFILE)

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) tidy

# Display help
.PHONY: help
help:
	@echo "Makefile for geth-builder"
	@echo
	@echo "Usage:"
	@echo "  make all        - Format, vet, and build the project"
	@echo "  make fmt        - Format the code"
	@echo "  make vet        - Vet the code"
	@echo "  make build      - Build the binary"
	@echo "  make clean      - Clean up build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make run        - Run the tool with the default configuration file"
	@echo "  make deps       - Install dependencies"
	@echo "  make help       - Display this help message"