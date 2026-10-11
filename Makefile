APP := familytree

.PHONY: help build run migrate seed createadmin test fmt vet lint tidy docker-up docker-down docker-logs

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary into ./bin
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

run: ## Run the server locally
	go run ./cmd/server

migrate: ## Apply database migrations then exit
	go run ./cmd/server -migrate

seed: ## Load demo data (needs psql + DATABASE_URL)
	psql "$(DATABASE_URL)" -f scripts/seed_demo.sql

createadmin: ## Create/update admin user: make createadmin USER=admin PASS=secret
	go run ./cmd/createadmin -username "$(USER)" -password "$(PASS)"

test: ## Run all tests
	go test ./...

fmt: ## Format all Go files
	gofmt -w .

vet: ## Static analysis
	go vet ./...

lint: ## Run golangci-lint (if installed)
	golangci-lint run

tidy: ## Sync go.mod / go.sum
	go mod tidy

docker-up: ## Build and start the docker stack
	docker compose up -d --build

docker-down: ## Stop the docker stack
	docker compose down

docker-logs: ## Follow app logs
	docker compose logs -f app
