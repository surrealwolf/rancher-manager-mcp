.PHONY: help build test test-integration test-coverage lint fmt clean docker-build docker-run docker-test install verify-token

# Variables
BINARY_NAME=rancher-mcp
VERIFY_TOKEN_BINARY=verify-token
BIN_DIR=bin
CMD_DIR=cmd
DOCKER_IMAGE=rancher-mcp
DOCKER_TAG=latest
HARBOR_REGISTRY=harbor.dataknife.net
HARBOR_PROJECT=library
HARBOR_IMAGE=$(HARBOR_REGISTRY)/$(HARBOR_PROJECT)/$(DOCKER_IMAGE)
GO_VERSION=1.23

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build the main binary
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✓ Built $(BIN_DIR)/$(BINARY_NAME)$(NC)"

build-verify-token: ## Build the verify-token tool
	@echo "$(GREEN)Building $(VERIFY_TOKEN_BINARY)...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(VERIFY_TOKEN_BINARY) ./$(CMD_DIR)/verify-token
	@echo "$(GREEN)✓ Built $(BIN_DIR)/$(VERIFY_TOKEN_BINARY)$(NC)"

build-all: build build-verify-token ## Build all binaries

install: ## Install binaries to GOPATH/bin
	@echo "$(GREEN)Installing binaries...$(NC)"
	@go install ./$(CMD_DIR)
	@go install ./$(CMD_DIR)/verify-token
	@echo "$(GREEN)✓ Installed binaries$(NC)"

test: ## Run unit tests
	@echo "$(GREEN)Running unit tests...$(NC)"
	@go test -v ./... || true

test-integration: ## Run integration tests (requires RANCHER_URL and RANCHER_TOKEN)
	@echo "$(GREEN)Running integration tests...$(NC)"
	@if [ -z "$$RANCHER_URL" ] || [ -z "$$RANCHER_TOKEN" ]; then \
		echo "$(YELLOW)⚠ Skipping integration tests: RANCHER_URL and RANCHER_TOKEN must be set$(NC)"; \
		exit 0; \
	fi
	@go test -v ./internal/client/...

test-all: test test-integration ## Run all tests

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./... || true
	@if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"; \
	else \
		echo "$(YELLOW)⚠ No coverage data generated$(NC)"; \
	fi

test-coverage-integration: ## Run integration tests with coverage
	@echo "$(GREEN)Running integration tests with coverage...$(NC)"
	@if [ -z "$$RANCHER_URL" ] || [ -z "$$RANCHER_TOKEN" ]; then \
		echo "$(YELLOW)⚠ Skipping integration tests: RANCHER_URL and RANCHER_TOKEN must be set$(NC)"; \
		exit 0; \
	fi
	@go test -v -coverprofile=coverage-integration.out ./internal/client/...
	@go tool cover -html=coverage-integration.out -o coverage-integration.html
	@echo "$(GREEN)✓ Integration coverage report generated: coverage-integration.html$(NC)"

lint: ## Run linters
	@echo "$(GREEN)Running linters...$(NC)"
	@go fmt ./...
	@go vet ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

fmt: lint ## Format code (alias for lint)

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@rm -f coverage-integration.out coverage-integration.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(HARBOR_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"
	@echo "$(GREEN)✓ Docker image tagged: $(HARBOR_IMAGE):$(DOCKER_TAG)$(NC)"

docker-login: ## Login to Harbor registry
	@echo "$(GREEN)Logging into Harbor registry...$(NC)"
	@if [ -z "$$HARBOR_USERNAME" ] || [ -z "$$HARBOR_PASSWORD" ]; then \
		echo "$(YELLOW)⚠ Error: HARBOR_USERNAME and HARBOR_PASSWORD must be set$(NC)"; \
		exit 1; \
	fi
	@echo "$$HARBOR_PASSWORD" | docker login $(HARBOR_REGISTRY) -u "$$HARBOR_USERNAME" --password-stdin
	@echo "$(GREEN)✓ Logged into $(HARBOR_REGISTRY)$(NC)"

docker-pull: ## Pull Docker image from Harbor
	@echo "$(GREEN)Pulling Docker image from Harbor...$(NC)"
	@docker pull $(HARBOR_IMAGE):$(DOCKER_TAG) || echo "$(YELLOW)⚠ Image not found in registry$(NC)"
	@echo "$(GREEN)✓ Pulled $(HARBOR_IMAGE):$(DOCKER_TAG)$(NC)"

docker-push: ## Push Docker image to Harbor
	@echo "$(GREEN)Pushing Docker image to Harbor...$(NC)"
	@if [ -z "$$HARBOR_USERNAME" ] || [ -z "$$HARBOR_PASSWORD" ]; then \
		echo "$(YELLOW)⚠ Error: HARBOR_USERNAME and HARBOR_PASSWORD must be set$(NC)"; \
		exit 1; \
	fi
	@echo "$$HARBOR_PASSWORD" | docker login $(HARBOR_REGISTRY) -u "$$HARBOR_USERNAME" --password-stdin
	@docker push $(HARBOR_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)✓ Pushed $(HARBOR_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container (stdio mode)
	@echo "$(GREEN)Running Docker container (stdio mode)...$(NC)"
	@docker run --rm -i \
		-e RANCHER_URL=$${RANCHER_URL} \
		-e RANCHER_TOKEN=$${RANCHER_TOKEN} \
		-e RANCHER_INSECURE_SKIP_VERIFY=$${RANCHER_INSECURE_SKIP_VERIFY:-false} \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		--transport stdio

docker-run-http: ## Run Docker container (HTTP mode on port 8080)
	@echo "$(GREEN)Running Docker container (HTTP mode)...$(NC)"
	@docker run --rm -d \
		-p 8080:8080 \
		-e RANCHER_URL=$${RANCHER_URL} \
		-e RANCHER_TOKEN=$${RANCHER_TOKEN} \
		-e RANCHER_INSECURE_SKIP_VERIFY=$${RANCHER_INSECURE_SKIP_VERIFY:-false} \
		--name $(BINARY_NAME) \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		--transport http --http-addr :8080
	@echo "$(GREEN)✓ Container running on http://localhost:8080$(NC)"
	@echo "$(YELLOW)Stop with: docker stop $(BINARY_NAME)$(NC)"

docker-test: ## Test Docker image
	@echo "$(GREEN)Testing Docker image...$(NC)"
	@docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) --help
	@docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) ls -la /app/rancher-mcp
	@docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) ls -la /app/verify-token
	@echo "$(GREEN)✓ Docker image tests passed$(NC)"

docker-clean: ## Remove Docker image
	@echo "$(GREEN)Removing Docker image...$(NC)"
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@echo "$(GREEN)✓ Docker image removed$(NC)"

verify-token: build-verify-token ## Build and run verify-token tool
	@echo "$(GREEN)Verifying Rancher token...$(NC)"
	@if [ -z "$$RANCHER_URL" ] || [ -z "$$RANCHER_TOKEN" ]; then \
		echo "$(YELLOW)Error: RANCHER_URL and RANCHER_TOKEN must be set$(NC)"; \
		exit 1; \
	fi
	@$(BIN_DIR)/$(VERIFY_TOKEN_BINARY)

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

deps-update: ## Update dependencies
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

check: lint test ## Run lint and tests (CI check)

ci: check build-all ## Run full CI checks (lint, test, build)

dev: build ## Build and prepare for development
	@echo "$(GREEN)Development environment ready$(NC)"
	@echo "$(YELLOW)Set environment variables:$(NC)"
	@echo "  export RANCHER_URL=https://your-rancher-server"
	@echo "  export RANCHER_TOKEN=token-XXXXX:YYYYY"
	@echo "  export RANCHER_INSECURE_SKIP_VERIFY=false"
	@echo ""
	@echo "$(YELLOW)Run the server:$(NC)"
	@echo "  ./$(BIN_DIR)/$(BINARY_NAME) --transport stdio"

.DEFAULT_GOAL := help
