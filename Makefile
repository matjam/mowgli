.PHONY: help build test test-go test-js build-js clean install-js

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-js: ## Build the JavaScript/TypeScript project
	@echo "Building JavaScript/TypeScript project..."
	cd js && yarn build

build: build-js ## Build all projects

test-go: ## Run Go tests
	@echo "Running Go tests..."
	go test -v ./...

test-js: ## Run JavaScript/TypeScript tests
	@echo "Running JavaScript/TypeScript tests..."
	cd js && yarn test

test: test-go test-js ## Run all tests (Go and JavaScript/TypeScript)

install-js: ## Install JavaScript dependencies
	@echo "Installing JavaScript dependencies..."
	cd js && yarn install

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf js/dist
	rm -rf js/node_modules/.cache
	go clean -cache -testcache

clean-all: clean ## Clean all artifacts including node_modules
	@echo "Cleaning all artifacts..."
	rm -rf js/node_modules
	rm -rf js/yarn.lock

.DEFAULT_GOAL := help

