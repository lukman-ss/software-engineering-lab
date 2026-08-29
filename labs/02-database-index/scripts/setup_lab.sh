#!/bin/bash
# Lab 02: Database Index Foundation
# Setup script to prepare the lab environment

set -e

DB_NAME="${DB_NAME:-se_lab}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

PQ="PGPASSWORD='$DB_PASSWORD' psql -h '$DB_HOST' -p '$DB_PORT' -U '$DB_USER' -d '$DB_NAME' -v ON_ERROR_STOP=1"

echo "=== Setting up Lab 02: Database Index ==="

# Check if schema exists
if $PQ -tAc "SELECT to_regclass('service')" | grep -q "service"; then
    echo "Schema already exists. Cleaning up..."
    $PQ -f "$(dirname "$0")/../cleanup.sql"
fi

echo "Creating schema..."
$PQ -f "$(dirname "$0")/../schema.sql"

echo "Seeding data..."
$PQ -f "$(dirname "$0")/../seed.sql"

echo "=== Lab 02 Setup Complete ==="
echo "Run verification with: make lab-02-verify"
echo "Or run experiments with: make lab-02-baseline"