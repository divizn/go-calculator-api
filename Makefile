include .env

.PHONY: docs

build:
	@echo "Building project..."

	@go build -o main cmd/api/main.go

# use air instead
run:
	@echo "Starting dev server..."

	@go run cmd/api/main.gom

db:
	@docker compose up -d


dbc:
	@docker compose down

docs:
	@swag init --dir ./cmd/api,internal/handler,internal/models --output ./docs


goose-up:
	@goose up


goose-down:
	@goose down

start-app:
	@docker compose --profile app up --build

