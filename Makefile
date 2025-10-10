.PHONY: help build test run clean docker-build docker-up docker-down lint coverage

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o server ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out

run: ## Run the application locally
	@echo "Running application..."
	@go run ./cmd/server

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f server coverage.out coverage.html
	@go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t myapp:latest .

docker-up: ## Start Docker Compose services
	@echo "Starting services..."
	@docker-compose up -d

docker-down: ## Stop Docker Compose services
	@echo "Stopping services..."
	@docker-compose down

docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

lint: ## Run linter
	@echo "Running linter..."
	@go fmt ./...
	@go vet ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

mock-generate: ## Generate mocks
	@echo "Generating mocks..."
	@go generate ./...

helm-install: ## Install with Helm
	@echo "Installing with Helm..."
	@helm install myapp ./helm/myapp

helm-upgrade: ## Upgrade with Helm
	@echo "Upgrading with Helm..."
	@helm upgrade myapp ./helm/myapp

helm-uninstall: ## Uninstall Helm release
	@echo "Uninstalling Helm release..."
	@helm uninstall myapp
