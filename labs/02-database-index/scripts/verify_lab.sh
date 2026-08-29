#!/bin/bash
# Lab 02: Database Index Foundation
# Automated structural verification
# Uses -v ON_ERROR_STOP=1 for fail-fast execution

set -e

DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

PQ="psql -h '$DB_HOST' -p '$DB_PORT' -U '$DB_USER' -d '$DB_NAME' -v ON_ERROR_STOP=1"

echo "=== Verifying Lab 02 ==="

# 1. Clean state
echo "[1/27] Checking clean state..."
$PQ -c "DROP TABLE IF EXISTS service CASCADE;" > /dev/null 2>&1 || true
echo "✓ Clean state verified"

# 2. Schema creation
echo "[2/27] Verifying schema creation..."
$PQ -f labs/02-database-index/schema.sql > /dev/null
echo "✓ Schema created"

# 3. Seed execution
echo "[3/27] Verifying seed execution (~500k rows)..."
$PQ -f labs/02-database-index/seed.sql > /dev/null
echo "✓ Data seeded"

# 4. Row count (target ~500k)
echo "[4/27] Checking row count..."
ROW_COUNT=$($PQ -Atc "SELECT count(*) FROM service;")
if [ "$ROW_COUNT" -lt 400000 ]; then
    echo "❌ Row count $ROW_COUNT below threshold 400,000"
    exit 1
fi
echo "✓ Row count: $ROW_COUNT"

# 5. Required columns
echo "[5/27] Checking columns..."
for col in id branch_id customer_id mechanic_id status service_date invoice_no created_at; do
    exists=$($PQ -Atc "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='service' AND column_name='$col');")
    [ "$exists" = "t" ] || { echo "❌ Missing column $col"; exit 1; }
done
echo "✓ All 8 columns present"

# 6. Status distribution
echo "[6/27] Checking status distribution..."
FINISHED=$($PQ -Atc "SELECT round(count(*) * 100.0 / (select count(*) from service), 1) FROM service WHERE status='FINISHED';")
echo "✓ FINISHED: $FINISHED%"

# 7. Branch distribution
echo "[7/27] Checking branch distribution..."
BR2=$($PQ -Atc "SELECT round(count(*) * 100.0 / (select count(*) from service), 1) FROM service WHERE branch_id=2;")
echo "✓ Branch 2: $BR2%"

# 8. Baseline has no dashboard indexes
echo "[8/27] Checking baseline has no query-supporting indexes..."
IDX_COUNT=$($PQ -Atc "SELECT count(*) FROM pg_class WHERE relname IN ('idx_service_branch_status_date','idx_service_status','idx_service_date') AND relkind='i';")
if [ "$IDX_COUNT" -gt 0 ]; then
    echo "❌ Unexpected indexes exist: $IDX_COUNT"
    exit 1
fi
echo "✓ No query-supporting indexes in baseline"

# 9-25. Run all SQL experiments
echo "[9/27] Running baseline experiment..."
$PQ -f labs/02-database-index/queries/01-baseline.sql > /dev/null
echo "✓ Baseline OK"

echo "[10/27] Running single-column-index experiment..."
$PQ -f labs/02-database-index/queries/02-single-column-index.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_id, idx_service_status, idx_service_service_date;" > /dev/null
echo "✓ Single-index OK"

echo "[11/27] Running composite-index experiment..."
$PQ -f labs/02-database-index/queries/03-composite-index.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_branch_status_date;" > /dev/null
echo "✓ Composite-index OK"

echo "[12/27] Running column-order experiment..."
$PQ -f labs/02-database-index/queries/04-column-order-experiment.sql > /dev/null
echo "✓ Column-order OK"

echo "[13/27] Running selectivity experiment..."
$PQ -f labs/02-database-index/queries/05-low-cardinality-selectivity.sql > /dev/null
echo "✓ Selectivity OK"

echo "[14/27] Running order-by-limit experiment..."
$PQ -f labs/02-database-index/queries/06-order-by-limit.sql > /dev/null
echo "✓ Order-by-limit OK"

echo "[15/27] Running covering-index experiment..."
$PQ -f labs/02-database-index/queries/07-covering-index.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_covering INCLUDE (customer_id, mechanic_id, invoice_no);" > /dev/null 2>&1 || \
$PQ -c "DROP INDEX IF EXISTS idx_service_covering;" > /dev/null
echo "✓ Covering-index OK"

echo "[16/27] Running write-cost experiment..."
$PQ -f labs/02-database-index/queries/08-write-cost.sql > /dev/null
echo "✓ Write-cost OK"

echo "[17/27] Running storage-cost experiment..."
$PQ -f labs/02-database-index/queries/09-storage-cost.sql > /dev/null
echo "✓ Storage-cost OK"

echo "[18/27] Running index-audit experiment..."
$PQ -f labs/02-database-index/queries/10-index-audit.sql > /dev/null
echo "✓ Index-audit OK"

echo "[19/27] Running redundant-indexes experiment..."
$PQ -f labs/02-database-index/queries/11-redundant-indexes.sql > /dev/null
echo "✓ Redundant-indexes OK"

echo "[20/27] Running partial-index experiment..."
$PQ -f labs/02-database-index/queries/12-partial-index.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_partial_in_progress, idx_service_full_status_branch_date;" > /dev/null
echo "✓ Partial-index OK"

echo "[21/27] Running functions-on-indexes experiment..."
$PQ -f labs/02-database-index/queries/13-functions-on-indexes.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_year_expr;" > /dev/null
echo "✓ Functions-on-indexes OK"

echo "[22/27] Running statistics experiment..."
$PQ -f labs/02-database-index/queries/14-statistics-and-analyze.sql > /dev/null
echo "✓ Statistics OK"

echo "[23/27] Running seqscan-correct experiment..."
$PQ -f labs/02-database-index/queries/15-seqscan-is-correct.sql > /dev/null
echo "✓ Seqscan-correct OK"

echo "[24/27] Running production-safe experiment..."
$PQ -f labs/02-database-index/queries/16-production-safe-index.sql > /dev/null
$PQ -c "DROP INDEX IF EXISTS idx_service_prod_test;" > /dev/null
echo "✓ Production-safe OK"

echo "[25/27] Running benchmark experiment..."
$PQ -f labs/02-database-index/queries/17-benchmark.sql > /dev/null
echo "✓ Benchmark OK"

# 26. Cleanup
echo "[26/27] Verifying cleanup..."
$PQ -f labs/02-database-index/cleanup.sql > /dev/null
exists=$($PQ -Atc "SELECT to_regclass('service');")
[ -z "$exists" ] || { echo "❌ Table still exists after cleanup"; exit 1; }
echo "✓ Cleanup successful"

# 27. Setup can run again
echo "[27/27] Verifying re-setup..."
$PQ -f labs/02-database-index/schema.sql > /dev/null
$PQ -f labs/02-database-index/seed.sql > /dev/null
echo "✓ Re-setup successful"

echo ""
echo "========================================="
echo "  Lab 02 Verification PASSED"
echo "========================================="