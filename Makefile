.PHONY: all run test test-race lint fmt vet clean infra-up infra-down migrate-up migrate-down

# Build and run API server
run: fmt vet
	@echo "Starting API server..."
	@go run ./cmd/api

# Run tests
test:
	@echo "Running unit tests..."
	@go test ./... -v

test-race:
	@echo "Running tests with race detector..."
	@go test ./... -race -v

lint:
	@echo "Running linter..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -f coverage.out

infra-up:
	@echo "Starting infrastructure..."
	@docker compose up -d

infra-down:
	@echo "Stopping infrastructure..."
	@docker compose down

migrate-up:
	@echo "Running migrations up..."
	@go run ./scripts/migrate main up

migrate-down:
	@echo "Running migrations down..."
	@go run ./scripts/migrate main down