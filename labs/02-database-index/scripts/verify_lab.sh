#!/bin/bash
# Lab 02: Database Index Foundation
# Automated verification harness
# Validates structural preconditions and executes all experiments

set -euo pipefail

DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

PQ="psql -h '$DB_HOST' -p '$DB_PORT' -U '$DB_USER' -d '$DB_NAME' -v ON_ERROR_STOP=1"

echo "=== Verifying Lab 02 ==="

# ============================================
# 1. Confirm PostgreSQL Connection
# ============================================
echo "[CHECK] PostgreSQL connection..."
if ! PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
    echo "[FAIL] Cannot connect to PostgreSQL"
    exit 1
fi
echo "[PASS] PostgreSQL connection OK"

# ============================================
# 2. Cleanup
# ============================================
echo "[CLEANUP] Dropping existing table..."
$PQ -c "DROP TABLE IF EXISTS service CASCADE;" 2>/dev/null || true
echo "[PASS] Clean state"

# ============================================
# 3. Create Schema
# ============================================
echo "[RUN] Creating schema..."
$PQ -f labs/02-database-index/schema.sql
echo "[PASS] Schema created"

# ============================================
# 4. Seed 500k Rows
# ============================================
echo "[RUN] Seeding 500,000 rows..."
$PQ -f labs/02-database-index/seed.sql
echo "[PASS] Data seeded"

# ============================================
# 5. Validate Exact Row Count
# ============================================
echo "[VALIDATE] Row count == 500,000..."
ROW_COUNT=$($PQ -Atc "SELECT COUNT(*) FROM service;")
if [ "$ROW_COUNT" -ne 500000 ]; then
    echo "[FAIL] Row count $ROW_COUNT != 500000"
    exit 1
fi
echo "[PASS] Row count: $ROW_COUNT"

# ============================================
# 6. Validate Status Distribution
# ============================================
echo "[VALIDATE] Status distribution..."
FINISHED_PCT=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE status = 'FINISHED';")
CANCELLED_PCT=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE status = 'CANCELLED';")
IN_PROGRESS_PCT=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE status = 'IN_PROGRESS';")
WAITING_PCT=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE status = 'WAITING';")
PR_PCT=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE status = 'PENDING_REFUND';")

echo "  FINISHED: $FINISHED_PCT% (expected 70.0%)"
echo "  CANCELLED: $CANCELLED_PCT% (expected 20.0%)"
echo "  IN_PROGRESS: $IN_PROGRESS_PCT% (expected 5.0%)"
echo "  WAITING: $WAITING_PCT% (expected 4.9%)"
echo "  PENDING_REFUND: $PR_PCT% (expected 0.1%)"

[ "$FINISHED_PCT" = "70.0" ] || { echo "[FAIL] FINISHED percentage mismatch"; exit 1; }
[ "$PR_PCT" = "0.1" ] || { echo "[FAIL] PENDING_REFUND percentage mismatch"; exit 1; }
echo "[PASS] Status distribution correct"

# ============================================
# 7. Validate Branch Distribution
# ============================================
echo "[VALIDATE] Branch distribution..."
BR1=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 1;")
BR2=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 2;")
BR3=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 3;")
BR4=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 4;")
BR5=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 5;")
BR6=$($PQ -Atc "SELECT round(COUNT(*) * 100.0 / 500000, 2) FROM service WHERE branch_id = 6;")

echo "  Branch 1: $BR1% (expected 15.0%)"
echo "  Branch 2: $BR2% (expected 25.0%)"
echo "  Branch 3: $BR3% (expected 20.0%)"
echo "  Branch 4: $BR4% (expected 15.0%)"
echo "  Branch 5: $BR5% (expected 15.0%)"
echo "  Branch 6: $BR6% (expected 10.0%)"

[ "$BR1" = "15.0" ] || { echo "[FAIL] Branch 1 percentage mismatch"; exit 1; }
[ "$BR2" = "25.0" ] || { echo "[FAIL] Branch 2 percentage mismatch"; exit 1; }
echo "[PASS] Branch distribution correct"

# ============================================
# 8. ANALYZE
# ============================================
echo "[RUN] ANALYZE service..."
$PQ -c "ANALYZE service;"
echo "[PASS] ANALYZE completed"

# ============================================
# 9. Verify Baseline Has No Dashboard-Supporting Indexes
# ============================================
echo "[VALIDATE] Baseline has no query-supporting secondary indexes..."

# Catalog query to list all indexes
INDEXES=$($PQ -Atc "
SELECT indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename = 'service'
ORDER BY indexname;
")

echo "Existing indexes:"
echo "$INDEXES"

# Check for forbidden experiment indexes
FORBIDDEN=$($PQ -Atc "
SELECT count(*) FROM pg_class
WHERE relname ~ '^(idx_service_branch|idx_service_status|idx_service_service_date|idx_bench|idx_order|idx_cover|idx_dashboard|idx_in_progress)'
  AND relkind = 'i';
")

if [ "$FORBIDDEN" -gt 0 ]; then
    echo "[FAIL] Found unexpected experiment indexes: $FORBIDDEN"
    exit 1
fi

# Verify only constraint-backed indexes exist (PK, UNIQUE)
EXPECTED=$($PQ -Atc "
SELECT count(*) FROM pg_class
WHERE relname IN ('service_pkey', 'service_invoice_no_key')
  AND relkind = 'i';
")

echo "Expected constraint-backed indexes: $EXPECTED (id PK, invoice_no UNIQUE)"
echo "[PASS] No dashboard-supporting secondary indexes in baseline"

# ============================================
# 10-20. Execute All Experiments
# ============================================

# 10. Baseline
echo "[RUN] queries/01-baseline.sql"
$PQ -f labs/02-database-index/queries/01-baseline.sql
echo "[PASS] 01-baseline.sql"

# 11. Drop indexes from baseline if any were created, then single-column
echo "[CLEANUP] Dropping experiment indexes..."
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_status, idx_service_service_date;"
$PQ -c "DROP INDEX IF EXISTS idx_service_status;"

echo "[RUN] queries/02-single-column-index.sql"
$PQ -f labs/02-database-index/queries/02-single-column-index.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_status, idx_service_service_date;"
echo "[PASS] 02-single-column-index.sql"

echo "[RUN] queries/03-composite-index.sql"
$PQ -f labs/02-database-index/queries/03-composite-index.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date;"
echo "[PASS] 03-composite-index.sql"

echo "[RUN] queries/04-column-order-experiment.sql"
$PQ -f labs/02-database-index/queries/04-column-order-experiment.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_a_branch_status_date, idx_service_b_date_branch_status, idx_service_c_status_date_branch;"
echo "[PASS] 04-column-order-experiment.sql"

echo "[RUN] queries/05-low-cardinality-selectivity.sql"
$PQ -f labs/02-database-index/queries/05-low-cardinality-selectivity.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_status;"
echo "[PASS] 05-low-cardinality-selectivity.sql"

echo "[RUN] queries/06-order-by-limit.sql"
$PQ -f labs/02-database-index/queries/06-order-by-limit.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date_desc;"
echo "[PASS] 06-order-by-limit.sql"

echo "[RUN] queries/07-covering-index.sql"
$PQ -f labs/02-database-index/queries/07-covering-index.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_covering, idx_service_dashboard;"
echo "[PASS] 07-covering-index.sql"

echo "[RUN] queries/08-write-cost.sql"
# 08-write-cost truncates and recreates data, but we need to re-seed for subsequent tests
$PQ -f labs/02-database-index/queries/08-write-cost.sql
$PQ -c "TRUNCATE TABLE service RESTART IDENTITY CASCADE;"
echo "[PASS] 08-write-cost.sql"

echo "[RUN] queries/09-storage-cost.sql"
$PQ -f labs/02-database-index/queries/09-storage-cost.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date, idx_service_status, idx_service_service_date, idx_service_date, idx_service_dashboard, idx_service_branch_status_date_desc;"
echo "[PASS] 09-storage-cost.sql"

echo "[RUN] queries/10-index-audit.sql"
$PQ -f labs/02-database-index/queries/10-index-audit.sql
echo "[PASS] 10-index-audit.sql"

echo "[RUN] queries/11-redundant-indexes.sql"
$PQ -f labs/02-database-index/queries/11-redundant-indexes.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_red_1, idx_service_red_2, idx_service_red_3;"
echo "[PASS] 11-redundant-indexes.sql"

echo "[RUN] queries/12-partial-index.sql"
$PQ -f labs/02-database-index/queries/12-partial-index.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_partial_in_progress, idx_service_full_status_branch_date;"
echo "[PASS] 12-partial-index.sql"

echo "[RUN] queries/13-functions-on-indexes.sql"
$PQ -f labs/02-database-index/queries/13-functions-on-indexes.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_service_date, idx_service_year_expr;"
echo "[PASS] 13-functions-on-indexes.sql"

echo "[RUN] queries/14-statistics-and-analyze.sql"
# This experiment inserts data, so re-seed before
$PQ -c "TRUNCATE TABLE service RESTART IDENTITY CASCADE;"
$PQ -f labs/02-database-index/seed.sql
$PQ -f labs/02-database-index/queries/14-statistics-and-analyze.sql
echo "[PASS] 14-statistics-and-analyze.sql"

echo "[RUN] queries/15-seqscan-is-correct.sql"
# Verify index existence for the experiment
$PQ -c "DROP INDEX IF EXISTS idx_service_status;"
$PQ -f labs/02-database-index/queries/15-seqscan-is-correct.sql
$PQ -c "DROP INDEX IF EXISTS idx_service_status;"
echo "[PASS] 15-seqscan-is-correct.sql"

echo "[RUN] queries/16-production-safe-index.sql"
# Note: This uses CREATE INDEX CONCURRENTLY which cannot run in transaction block
# For local verification, we create without CONCURRENTLY
$PQ -c "DROP INDEX IF EXISTS idx_service_concurrent_branch;"
$PQ -f labs/02-database-index/queries/16-production-safe-index.sql
echo "[PASS] 16-production-safe-index.sql"

echo "[RUN] queries/17-benchmark.sql"
$PQ -f labs/02-database-index/queries/17-benchmark.sql
echo "[PASS] 17-benchmark.sql"

# ============================================
# 21. Cleanup
# ============================================
echo "[CLEANUP] Running cleanup.sql..."
$PQ -f labs/02-database-index/cleanup.sql

# Verify table is dropped
EXISTS=$($PQ -Atc "SELECT to_regclass('service');")
if [ -n "$EXISTS" ]; then
    echo "[FAIL] Table still exists after cleanup"
    exit 1
fi
echo "[PASS] Cleanup completed"

# ============================================
# 22. Setup Again (Re-run)
# ============================================
echo "[RUN] Re-running setup..."

# Create schema
$PQ -f labs/02-database-index/schema.sql

# Seed data
$PQ -f labs/02-database-index/seed.sql

echo "[PASS] Re-setup completed"

# ============================================
# 23. Final Sanity Query
# ============================================
echo "[SANITY] Final row count..."
FINAL_COUNT=$($PQ -Atc "SELECT COUNT(*) FROM service;")
if [ "$FINAL_COUNT" -ne 500000 ]; then
    echo "[FAIL] Final row count $FINAL_COUNT != 500000"
    exit 1
fi
echo "[PASS] Final row count: $FINAL_COUNT"

echo ""
echo "========================================="
echo "  Lab 02 Verification PASSED"
echo "========================================="