.PHONY: all run test test-race lint fmt vet clean \
	infra-up infra-down migrate-up migrate-down \
	lab-02-setup lab-02-seed lab-02-baseline lab-02-indexes \
	lab-02-explain lab-02-benchmark lab-02-clean lab-02-verify \
	lab-04-test lab-04-test-race lab-04-vet lab-04-demo lab-04-integration \
	lab-05-test lab-05-test-race lab-05-vet lab-05-fmt lab-05-integration \
	lab-07-test lab-07-test-race lab-07-vet lab-07-fmt lab-07-demo \
	lab-08-test lab-08-test-race lab-08-vet \
	lab-09-test lab-09-test-race lab-09-vet \
	lab-14-test lab-14-test-race lab-14-vet \
	lab-15-test lab-15-test-race lab-15-vet \
	lab-16-test lab-16-test-race lab-16-vet

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
test: lab-04-test lab-05-test
	@echo "Running unit tests..."
	@go test ./... -v

test-race: lab-04-test-race lab-05-test-race
	@echo "Running tests with race detector..."
	@go test ./... -race -v

lint: lab-04-vet lab-05-vet
	@echo "Running linter..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@cd labs/04-caching && go fmt ./...
	@cd labs/05-race-condition && go fmt ./...

vet: lab-04-vet lab-05-vet
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

# ==================== Lab 04: Caching ====================

lab-04-test:
	@echo "=== Testing Lab 04 (Unit Tests) ==="
	@cd labs/04-caching && go test -count=1 ./...

lab-04-test-race:
	@echo "=== Testing Lab 04 (Race Detector) ==="
	@cd labs/04-caching && go test -race -count=1 ./...

lab-04-vet:
	@echo "=== Vet Lab 04 ==="
	@cd labs/04-caching && go vet ./...

lab-04-demo:
	@echo "=== Running Lab 04 Demo ==="
	@cd labs/04-caching && go run ./cmd/demo -scenario=cache-aside

lab-04-integration:
	@echo "=== Testing Lab 04 (Integration) ==="
	@cd labs/04-caching && go test -tags=integration -count=1 ./...

# ==================== Lab 05: Race Condition ====================

lab-05-test:
	@echo "=== Testing Lab 05: Race Condition ==="
	@cd labs/05-race-condition && go test -count=1 ./...

lab-05-test-race:
	@echo "=== Testing Lab 05 (Race Detector) ==="
	@cd labs/05-race-condition && go test -race -count=1 ./...

lab-05-vet:
	@echo "=== Vet Lab 05 ==="
	@cd labs/05-race-condition && go vet ./...

lab-05-fmt:
	@echo "=== Formatting Lab 05 ==="
	@cd labs/05-race-condition && go fmt ./...

lab-05-integration:
	@echo "=== Testing Lab 05 (Integration) ==="
	@cd labs/05-race-condition && go test -tags=integration -count=1 ./...

# ==================== Lab 07: Observability ====================

lab-07-test:
	@echo "=== Testing Lab 07: Observability ==="
	@cd labs/07-observability && go test -v -count=1 ./...

lab-07-test-race:
	@echo "=== Testing Lab 07 (Race Detector) ==="
	@cd labs/07-observability && go test -race -v -count=1 ./...

lab-07-vet:
	@echo "=== Vet Lab 07 ==="
	@cd labs/07-observability && go vet ./...

lab-07-fmt:
	@echo "=== Formatting Lab 07 ==="
	@cd labs/07-observability && go fmt ./...

lab-07-demo:
	@echo "=== Running Lab 07 Demo ==="
	@go run ./labs/07-observability/cmd/demo

# ==================== Lab 08: Database Isolation Level ====================

lab-08-test:
	@echo "=== Testing Lab 08: Database Isolation Level ==="
	@cd labs/08-database-isolation-level && go test -v -count=1 ./...

lab-08-test-race:
	@echo "=== Testing Lab 08 (Race Detector) ==="
	@cd labs/08-database-isolation-level && go test -race -v -count=1 ./...

lab-08-vet:
	@echo "=== Vet Lab 08 ==="
	@cd labs/08-database-isolation-level && go vet ./...

# ==================== Lab 14: Outbox Pattern ====================

lab-14-test:
	@echo "=== Testing Lab 14: Outbox Pattern ==="
	@go test -v ./labs/14-outbox-pattern/...

lab-14-test-race:
	@echo "=== Testing Lab 14 (Race Detector) ==="
	@go test -race -v ./labs/14-outbox-pattern/...

lab-14-vet:
	@echo "=== Vet Lab 14 ==="
	@go vet ./labs/14-outbox-pattern/...

# ==================== Lab 15: Retry ====================

lab-15-test:
	@echo "=== Testing Lab 15: Retry ==="
	@go test -v ./labs/15-retry/...

lab-15-test-race:
	@echo "=== Testing Lab 15 (Race Detector) ==="
	@go test -race -v ./labs/15-retry/...

lab-15-vet:
	@echo "=== Vet Lab 15 ==="
	@go vet ./labs/15-retry/...

# ==================== Lab 09: Code Review ====================

lab-09-test:
	@echo "=== Testing Lab 09: Code Review ==="
	@cd labs/09-code-review && go test -v -count=1 ./...

lab-09-test-race:
	@echo "=== Testing Lab 09 (Race Detector) ==="
	@cd labs/09-code-review && go test -race -v -count=1 ./...

lab-09-vet:
	@echo "=== Vet Lab 09 ==="
	@cd labs/09-code-review && go vet ./...

# ==================== Lab 16: Circuit Breaker ====================

lab-16-test:
	@echo "=== Testing Lab 16: Circuit Breaker ==="
	@cd labs/16-circuit-breaker && go test -v -count=1 ./...

lab-16-test-race:
	@echo "=== Testing Lab 16 (Race Detector) ==="
	@cd labs/16-circuit-breaker && go test -race -v -count=1 ./...

lab-16-vet:
	@echo "=== Vet Lab 16 ==="
	@cd labs/16-circuit-breaker && go vet ./...
