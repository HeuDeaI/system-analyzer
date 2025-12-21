# Makefile for the System Analyzer project

.PHONY: run build clean

# Go parameters
GO = go
GO_FLAGS = -ldflags="-s -w"

# Application parameters
APP_NAME = system-analyzer
CMD_PATH = ./cmd/analyzer

# Suppress linker warning on macOS
ifeq ($(shell uname), Darwin)
	CGO_LDFLAGS = "-Wl,-w"
endif

# Default target
all: build

# Run the application
run: 
	@echo "Running the application..."
	@$(GO) run $(CMD_PATH)

# Build the application
build:
	@echo "Building the application..."
	@CGO_LDFLAGS=$(CGO_LDFLAGS) $(GO) build $(GO_FLAGS) -o $(APP_NAME) $(CMD_PATH)

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
