.PHONY: help build test clean swagger run

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	go build -v -o server ./cmd/server

test: ## Run tests
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -f server coverage.out coverage.html
	rm -rf docs

swagger: ## Generate Swagger documentation
	@command -v swag >/dev/null 2>&1 || { echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	swag init -g cmd/server/main.go --output docs --parseDependency --parseInternal

run: swagger ## Run the application
	go run ./cmd/server

tidy: ## Tidy go.mod
	go mod tidy

fmt: ## Format code
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed"; exit 1; }
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t myapp:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 myapp:latest

migrate-up: ## Run database migrations
	@command -v migrate >/dev/null 2>&1 || { echo "migrate CLI not installed. Install from https://github.com/golang-migrate/migrate"; exit 1; }
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Rollback last migration
	@command -v migrate >/dev/null 2>&1 || { echo "migrate CLI not installed. Install from https://github.com/golang-migrate/migrate"; exit 1; }
	migrate -path migrations -database "${DATABASE_URL}" down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@command -v migrate >/dev/null 2>&1 || { echo "migrate CLI not installed. Install from https://github.com/golang-migrate/migrate"; exit 1; }
	@test -n "$(NAME)" || { echo "NAME is required. Usage: make migrate-create NAME=migration_name"; exit 1; }
	migrate create -ext sql -dir migrations -seq $(NAME)

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@command -v migrate >/dev/null 2>&1 || { echo "migrate CLI not installed. Install from https://github.com/golang-migrate/migrate"; exit 1; }
	@test -n "$(VERSION)" || { echo "VERSION is required. Usage: make migrate-force VERSION=1"; exit 1; }
	migrate -path migrations -database "${DATABASE_URL}" force $(VERSION)

