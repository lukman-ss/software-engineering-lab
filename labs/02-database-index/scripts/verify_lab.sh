#!/bin/bash
# Lab 02: Database Index Foundation
# Automated structural verification script

set -e

DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

PQ="PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -A"

echo "=== Verifying Lab 02 Structure and Execution ==="

# 1. Verify schema creates successfully
echo "[1/8] Verifying schema creation..."
$PQ -f labs/02-database-index/schema.sql > /dev/null
echo "✓ Schema created successfully"

# 2. Verify seed executes successfully
echo "[2/8] Verifying data seeding (~500k rows)..."
$PQ -f labs/02-database-index/seed.sql > /dev/null
echo "✓ Data seeded successfully"

# 3. Verify expected table exists
echo "[3/8] Verifying table existence..."
TABLE_EXISTS=$($PQ -c "SELECT to_regclass('service');")
if [ "$TABLE_EXISTS" != "service" ]; then
    echo "❌ Table 'service' does not exist!"
    exit 1
fi
echo "✓ Table 'service' exists"

# 4. Verify expected columns exist
echo "[4/8] Verifying table columns..."
COLUMNS_COUNT=$($PQ -c "SELECT count(*) FROM information_schema.columns WHERE table_name = 'service' AND column_name IN ('id', 'branch_id', 'customer_id', 'mechanic_id', 'status', 'service_date', 'invoice_no', 'created_at');")
if [ "$COLUMNS_COUNT" -ne 8 ]; then
    echo "❌ Expected 8 columns, found $COLUMNS_COUNT"
    exit 1
fi
echo "✓ All 8 expected columns present"

# 5. Verify row count threshold is reached (>= 100,000)
echo "[5/8] Verifying row count threshold (>= 100,000)..."
ROW_COUNT=$($PQ -c "SELECT count(*) FROM service;")
echo "   Found $ROW_COUNT rows"
if [ "$ROW_COUNT" -lt 100000 ]; then
    echo "❌ Row count $ROW_COUNT is below 100,000 threshold"
    exit 1
fi
echo "✓ Row count threshold reached"

# 6. Verify index creation works
echo "[6/8] Verifying index creation..."
$PQ -c "CREATE INDEX idx_verify_composite ON service(branch_id, status, service_date DESC);" > /dev/null
echo "✓ Index created successfully"

# 7. Verify SQL experiment files parse and execute (EXPLAIN check)
echo "[7/8] Verifying query execution and EXPLAIN..."
$PQ -c "EXPLAIN ANALYZE SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED' LIMIT 5;" > /dev/null
echo "✓ Query execution and EXPLAIN verified"

# 8. Verify cleanup works
echo "[8/8] Verifying cleanup..."
$PQ -f labs/02-database-index/cleanup.sql > /dev/null
TABLE_EXISTS_AFTER=$($PQ -c "SELECT to_regclass('service');")
if [ ! -z "$TABLE_EXISTS_AFTER" ] && [ "$TABLE_EXISTS_AFTER" != "" ]; then
    echo "❌ Cleanup failed to drop service table!"
    exit 1
fi
echo "✓ Cleanup executed successfully"

echo "=== Lab 02 Verification Passed Successfully ==="