# Makefile

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# App related variables
BINARY_NAME=loggy
DATABASE_DIR=$(HOME)/.local/share/loggy
CONFIG_DIR=$(HOME)/.config/loggy

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

# Help
## logctl: Compiles the log viewer command line tool.
logctl: cmd/logctl/*.go
	@echo "Building $*..."
	@GOBIN=$(GOBIN) go build -o $(GOBIN)/logctl ./cmd/logctl
	@echo "$* build complete!"

## loggy: Compiles the loggy application.
loggy: cmd/loggy/*.go
	@echo "Building $*..."
	@GOBIN=$(GOBIN) go build -o $(GOBIN)/loggy ./cmd/loggy
	@echo "$* build complete!"

## run: Builds and runs the application.
run: build
	@echo "Running..."
	@$(GOBIN)/$(BINARY_NAME)

## install: Creates necessary directories and copies the default configuration file if it does not exist.
install: build
	@echo "Creating necessary directories..."
	@mkdir -p $(DATABASE_DIR)
	@mkdir -p $(CONFIG_DIR)
	@echo "Copying default config file if it does not exist..."
	@if [ ! -f $(CONFIG_DIR)/config.json ]; then \
		cp config.json $(CONFIG_DIR)/config.json; \
	fi
	@echo "Installation complete!"

## clean: Removes the binary and cleans up the build.
clean:
	@echo "Cleaning..."
	@GOBIN=$(GOBIN) go clean
	@rm -f $(GOBIN)/$(BINARY_NAME)
	@echo "Clean complete!"

## test: Runs the Go tests.
test:
	@echo "Testing..."
	@go test -v
	@echo "Testing complete!"


## help: Displays help for each target (this message).
help:
	@echo "Usage: ${YELLOW}make${NC} ${GREEN}<target>${NC}"
	@echo "\nTargets\n"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //g' -e 's/:/|/' | awk 'BEGIN {FS = "|";} {printf "\033[0;32m%-30s\033[0;33m %s\033[0m\n", $$1, $$2;}'


.PHONY: run install clean test help $(GOBIN)/%

