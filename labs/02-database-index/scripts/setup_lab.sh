#!/bin/bash
# Lab 02: Database Index Foundation
# Setup script to prepare the lab environment

set -e

DB_NAME="${DB_NAME:-software_engineer_lab}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

PQ="psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"

echo "=== Setting up Lab 02: Database Index ==="

# Check if schema exists
if $PQ -c "SELECT to_regclass('service')" | grep -q "service"; then
    echo "Schema already exists. Cleaning up..."
    $PQ -f "$(dirname "$0")/../cleanup.sql"
fi

echo "Creating schema..."
$PQ -f "$(dirname "$0")/../schema.sql"

echo "Seeding data..."
$PQ -f "$(dirname "$0")/../seed.sql"

echo "=== Lab 02 Setup Complete ==="
echo "Run queries with: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f queries/case1_queries.sql"