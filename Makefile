build:
	@echo "Building project..."

	@go build -o main cmd/api/main.go

# use air instead
run:
	@echo "Starting dev server..."

	@go run cmd/api/main.gom