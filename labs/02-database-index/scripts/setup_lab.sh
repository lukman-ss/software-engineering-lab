#!/bin/bash
# Lab 02: Database Index Foundation
# Setup script to prepare the lab environment

set -euo pipefail

# Database connection parameters (override with environment variables)
DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

# PQ: wrapper so every psql call uses the correct connection params
PQ() {
    PGPASSWORD="$DB_PASSWORD" \
    psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -v ON_ERROR_STOP=1 \
        "$@"
}

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "=== Setting up Lab 02: Database Index ==="

# Check if schema exists
if PQ -tAc "SELECT to_regclass('service')" | grep -q "service"; then
    echo "Schema already exists. Cleaning up..."
    PQ -f "$SCRIPT_DIR/../cleanup.sql"
fi

echo "Creating schema..."
PQ -f "$SCRIPT_DIR/../schema.sql"

echo "Seeding data..."
PQ -f "$SCRIPT_DIR/../seed.sql"

echo "=== Lab 02 Setup Complete ==="
echo ""
echo "Next steps:"
echo "  make lab-02-verify    # Run full experiment suite with validation"
echo "  make lab-02-baseline  # Run baseline query (no indexes)"
echo "  make lab-02-explain  # Run EXPLAIN ANALYZE on sample queries"
echo "  make lab-02-benchmark # Compare execution plans across index strategies"
echo "  make lab-02-clean   # Drop the service table"
echo ""
echo "Or run verification: $LAB_ROOT/labs/02-database-index/scripts/verify_lab.sh"