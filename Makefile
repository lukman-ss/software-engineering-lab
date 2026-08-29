.PHONY: all run test test-race lint fmt vet clean \
	infra-up infra-down migrate-up migrate-down \
	lab-02-setup lab-02-seed lab-02-baseline lab-02-indexes \
	lab-02-explain lab-02-benchmark lab-02-clean lab-02-verify

DB_NAME ?= se_lab
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_PSQL := psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -v ON_ERROR_STOP=1
PSQL := PGPASSWORD=$(DB_PASSWORD) $(DB_PSQL)

all: vet test

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

# ==================== Lab 02: Database Index ====================
# Requires a running PostgreSQL (use: make infra-up)

lab-02-setup:
	@echo "=== Setting up Lab 02: Database Index ==="
	@$(PSQL) -f labs/02-database-index/cleanup.sql
	@$(PSQL) -f labs/02-database-index/schema.sql
	@$(PSQL) -f labs/02-database-index/seed.sql
	@echo "=== Lab 02 Setup Complete ==="

lab-02-seed:
	@echo "=== Seeding Lab 02 data ==="
	@$(PSQL) -f labs/02-database-index/seed.sql
	@echo "=== Seeding Complete ==="

lab-02-baseline:
	@echo "=== Lab 02: Baseline query (no indexes) ==="
	@$(PSQL) -f labs/02-database-index/queries/01-baseline.sql
	@echo "=== Baseline Complete ==="

lab-02-indexes:
	@echo "=== Lab 02: Create indexes ==="
	@$(PSQL) -f labs/02-database-index/indexes/create_indexes.sql
	@echo "=== Indexes Created ==="

lab-02-explain:
	@echo "=== Lab 02: EXPLAIN ANALYZE ==="
	@$(PSQL) -f labs/02-database-index/queries/explain_queries.sql
	@echo "=== Explain Complete ==="

lab-02-benchmark:
	@echo "=== Lab 02: Benchmark ==="
	@$(PSQL) -f labs/02-database-index/queries/17-benchmark.sql
	@echo "=== Benchmark Complete ==="

lab-02-clean:
	@echo "=== Cleaning up Lab 02 ==="
	@$(PSQL) -f labs/02-database-index/cleanup.sql
	@echo "=== Cleanup Complete ==="

lab-02-verify:
	@echo "=== Verifying Lab 02 ==="
	@./labs/02-database-index/scripts/verify_lab.sh
	@echo "=== Verification Complete ==="