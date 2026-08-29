#!/bin/bash
# Lab 02: Database Index Foundation
# Automated verification harness
# Validates structural preconditions and executes all experiments

set -euo pipefail

DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

# PQ: wrapper so every psql call uses the correct connection params
PQ() { PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 "$@"; }

# run_sql: helper to execute a SQL file, preserving stderr for visibility
run_sql() {
    local file="$1"

    echo "[RUN] $file"

    if ! PQ -f "$file" >/dev/null; then
        echo "[FAIL] $file"
        return 1
    fi

    echo "[PASS] $file"
}

echo "=== Verifying Lab 02 ==="

# ============================================
# Helper: Reset to canonical 500k dataset
# ============================================
reset_dataset() {
    echo "[RESET] Restoring canonical 500,000-row dataset..."
    PQ -f labs/02-database-index/cleanup.sql >/dev/null
    PQ -f labs/02-database-index/schema.sql >/dev/null
    PQ -f labs/02-database-index/seed.sql >/dev/null

    local row_count
    row_count=$(PQ -Atc "SELECT COUNT(*) FROM service;")
    if [ "$row_count" -ne 500000 ]; then
        echo "[FAIL] Reset failed: Row count $row_count != 500000"
        exit 1
    fi
    echo "[PASS] Reset OK: $row_count rows"
}

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
# 2. Initial setup: canonical 500k dataset
# ============================================
reset_dataset

# ============================================
# Validation: branch 2 quality checks
# Prevents regression to old correlated seed
# ============================================
echo "[VALIDATE] Branch 2 quality checks..."

BR2_STATUSES=$(PQ -Atc "SELECT COUNT(DISTINCT status) FROM service WHERE branch_id = 2;")
if [ "$BR2_STATUSES" -le 1 ]; then
    echo "[FAIL] Branch 2 has only $BR2_STATUSES distinct status(es)"
    exit 1
fi
echo "[PASS] Branch 2 has $BR2_STATUSES distinct statuses"

BR2_FINISHED_PCT=$(PQ -Atc "SELECT round(COUNT(*) FILTER (WHERE status = 'FINISHED') * 100.0 / COUNT(*), 2) FROM service WHERE branch_id = 2;")
if [ "$BR2_FINISHED_PCT" = "100.00" ]; then
    echo "[FAIL] Branch 2 has 100% FINISHED — status not independent"
    exit 1
fi
echo "[PASS] Branch 2 FINISHED: $BR2_FINISHED_PCT%"

# ============================================
# 3. Validate Status Distribution (exact counts)
# ============================================
echo "[VALIDATE] Status distribution..."

FINISHED_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE status = 'FINISHED';")
CANCELLED_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE status = 'CANCELLED';")
IN_PROGRESS_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE status = 'IN_PROGRESS';")
WAITING_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE status = 'WAITING';")
PENDING_REFUND_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE status = 'PENDING_REFUND';")

echo "  FINISHED:       $FINISHED_COUNT (expected 350000)"
echo "  CANCELLED:      $CANCELLED_COUNT (expected 100000)"
echo "  IN_PROGRESS:    $IN_PROGRESS_COUNT (expected 25000)"
echo "  WAITING:        $WAITING_COUNT (expected 24500)"
echo "  PENDING_REFUND: $PENDING_REFUND_COUNT (expected 500)"

[ "$FINISHED_COUNT" -eq 350000 ] || { echo "[FAIL] FINISHED count mismatch: $FINISHED_COUNT"; exit 1; }
[ "$CANCELLED_COUNT" -eq 100000 ] || { echo "[FAIL] CANCELLED count mismatch: $CANCELLED_COUNT"; exit 1; }
[ "$IN_PROGRESS_COUNT" -eq 25000 ] || { echo "[FAIL] IN_PROGRESS count mismatch: $IN_PROGRESS_COUNT"; exit 1; }
[ "$WAITING_COUNT" -eq 24500 ]    || { echo "[FAIL] WAITING count mismatch: $WAITING_COUNT"; exit 1; }
[ "$PENDING_REFUND_COUNT" -eq 500 ] || { echo "[FAIL] PENDING_REFUND count mismatch: $PENDING_REFUND_COUNT"; exit 1; }
echo "[PASS] Status distribution correct"

# ============================================
# 4. Validate Branch Distribution (exact counts)
# ============================================
echo "[VALIDATE] Branch distribution..."

BR1_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 1;")
BR2_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 2;")
BR3_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 3;")
BR4_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 4;")
BR5_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 5;")
BR6_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service WHERE branch_id = 6;")

echo "  Branch 1: $BR1_COUNT (expected 75000)"
echo "  Branch 2: $BR2_COUNT (expected 125000)"
echo "  Branch 3: $BR3_COUNT (expected 100000)"
echo "  Branch 4: $BR4_COUNT (expected 75000)"
echo "  Branch 5: $BR5_COUNT (expected 75000)"
echo "  Branch 6: $BR6_COUNT (expected 50000)"

[ "$BR1_COUNT" -eq 75000 ] || { echo "[FAIL] Branch 1 count mismatch: $BR1_COUNT"; exit 1; }
[ "$BR2_COUNT" -eq 125000 ] || { echo "[FAIL] Branch 2 count mismatch: $BR2_COUNT"; exit 1; }
[ "$BR3_COUNT" -eq 100000 ] || { echo "[FAIL] Branch 3 count mismatch: $BR3_COUNT"; exit 1; }
[ "$BR4_COUNT" -eq 75000 ] || { echo "[FAIL] Branch 4 count mismatch: $BR4_COUNT"; exit 1; }
[ "$BR5_COUNT" -eq 75000 ] || { echo "[FAIL] Branch 5 count mismatch: $BR5_COUNT"; exit 1; }
[ "$BR6_COUNT" -eq 50000 ] || { echo "[FAIL] Branch 6 count mismatch: $BR6_COUNT"; exit 1; }
echo "[PASS] Branch distribution correct"

# ============================================
# 5. ANALYZE
# ============================================
echo "[RUN] ANALYZE service..."
PQ -c "ANALYZE service;" >/dev/null
echo "[PASS] ANALYZE completed"

# ============================================
# 6. Verify baseline has no experiment indexes
# ============================================
echo "[VALIDATE] Baseline has no query-supporting secondary indexes..."

INDEXES=$(PQ -Atc "
SELECT indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename = 'service'
ORDER BY indexname;
")
echo "Existing indexes:"
echo "$INDEXES"

FORBIDDEN=$(PQ -Atc "
SELECT count(*) FROM pg_class
WHERE relname ~ '^(idx_service_branch|idx_service_status|idx_service_service_date|idx_bench|idx_order|idx_cover|idx_dashboard|idx_in_progress)'
  AND relkind = 'i';
")

if [ "$FORBIDDEN" -gt 0 ]; then
    echo "[FAIL] Found unexpected experiment indexes: $FORBIDDEN"
    exit 1
fi
echo "[PASS] No experiment indexes in baseline"

# ============================================
# 7. Execute All Experiments
# ============================================

run_sql labs/02-database-index/queries/01-baseline.sql

echo "[CLEANUP] Dropping stray experiment indexes..."
PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_status, idx_service_service_date;" >/dev/null

run_sql labs/02-database-index/queries/02-single-column-index.sql
PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_status, idx_service_service_date;" >/dev/null

run_sql labs/02-database-index/queries/03-composite-index.sql
PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date;" >/dev/null

run_sql labs/02-database-index/queries/04-column-order-experiment.sql
PQ -c "DROP INDEX IF EXISTS idx_service_a_branch_status_date, idx_service_b_date_branch_status, idx_service_c_status_date_branch;" >/dev/null

run_sql labs/02-database-index/queries/05-low-cardinality-selectivity.sql
PQ -c "DROP INDEX IF EXISTS idx_service_status;" >/dev/null

run_sql labs/02-database-index/queries/06-order-by-limit.sql
PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date_desc;" >/dev/null

run_sql labs/02-database-index/queries/07-covering-index.sql
PQ -c "DROP INDEX IF EXISTS idx_service_covering, idx_service_dashboard;" >/dev/null

run_sql labs/02-database-index/queries/08-write-cost.sql
PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date, idx_service_status_single, idx_service_service_date, idx_service_customer;" >/dev/null

# Mandatory reset: 08-write-cost left the table in an altered state
reset_dataset

# Validate reset before continuing
ROW_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service;")
[ "$ROW_COUNT" -eq 500000 ] || { echo "[FAIL] Dataset not restored before experiment 09 (got $ROW_COUNT rows)"; exit 1; }

run_sql labs/02-database-index/queries/09-storage-cost.sql
PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_branch_status_date, idx_service_status, idx_service_service_date, idx_service_date, idx_service_dashboard, idx_service_branch_status_date_desc;" >/dev/null

run_sql labs/02-database-index/queries/10-index-audit.sql

run_sql labs/02-database-index/queries/11-redundant-indexes.sql
PQ -c "DROP INDEX IF EXISTS idx_service_red_1, idx_service_red_2, idx_service_red_3;" >/dev/null

run_sql labs/02-database-index/queries/12-partial-index.sql
PQ -c "DROP INDEX IF EXISTS idx_service_in_progress_branch_date, idx_service_full_status_branch_date;" >/dev/null

run_sql labs/02-database-index/queries/13-functions-on-indexes.sql
PQ -c "DROP INDEX IF EXISTS idx_service_service_date, idx_service_year_expr;" >/dev/null

# 14-statistics-and-analyze inserts 50k rows; reset before and clean up after
reset_dataset
run_sql labs/02-database-index/queries/14-statistics-and-analyze.sql
reset_dataset

PQ -c "DROP INDEX IF EXISTS idx_service_status;" >/dev/null
run_sql labs/02-database-index/queries/15-seqscan-is-correct.sql
PQ -c "DROP INDEX IF EXISTS idx_service_status;" >/dev/null

run_sql labs/02-database-index/queries/16-production-safe-index.sql

run_sql labs/02-database-index/queries/17-benchmark.sql

# ============================================
# 8. Final cleanup
# ============================================
echo "[CLEANUP] Running cleanup.sql..."
PQ -f labs/02-database-index/cleanup.sql >/dev/null

EXISTS=$(PQ -Atc "SELECT to_regclass('service');")
if [ -n "$EXISTS" ]; then
    echo "[FAIL] Table still exists after cleanup"
    exit 1
fi
echo "[PASS] Cleanup completed"

# ============================================
# 9. Re-seed for future use
# ============================================
echo "[RUN] Re-running setup..."
PQ -f labs/02-database-index/schema.sql >/dev/null
PQ -f labs/02-database-index/seed.sql >/dev/null

FINAL_COUNT=$(PQ -Atc "SELECT COUNT(*) FROM service;")
if [ "$FINAL_COUNT" -ne 500000 ]; then
    echo "[FAIL] Final row count $FINAL_COUNT != 500000"
    exit 1
fi
echo "[PASS] Final row count: $FINAL_COUNT"

echo ""
echo "========================================="
echo "  Lab 02 Verification PASSED"
echo "========================================="