-- Lab 02: Database Index Foundation
-- Experiment: Index Storage Cost Analysis
-- Demonstrates that indexes consume real storage resources

-- Key concept: An index is NOT "free" - it consumes disk space
-- Storage cost is part of the read/write trade-off

-- Clear existing indexes to start fresh
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status;
DROP INDEX IF EXISTS idx_service_service_date;
DROP INDEX IF EXISTS idx_service_date_branch;
DROP INDEX IF EXISTS idx_service_dashboard_covering;

-- ================================
-- PART 1: Create table only (no secondary indexes)
-- ================================

-- Table size without secondary indexes
SELECT
    'Table only (PK only)' AS description,
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS indexes_size,
    pg_size_pretty(pg_total_relation_size('service')) AS total_size;

-- ================================
-- PART 2: Add single-column indexes
-- ================================

CREATE INDEX idx_service_branch_id ON service(branch_id);
CREATE INDEX idx_service_status ON service(status);
CREATE INDEX idx_service_service_date ON service(service_date DESC);

SELECT
    'Three single-column indexes' AS description,
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS indexes_size,
    pg_size_pretty(pg_total_relation_size('service')) AS total_size;

-- List individual index sizes
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_indexes
WHERE tablename = 'service'
ORDER BY pg_relation_size(indexname::regclass) DESC;

-- ================================
-- PART 3: Add composite index
-- ================================

DROP INDEX idx_service_branch_id;
DROP INDEX idx_service_status;
DROP INDEX idx_service_service_date;

CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date DESC);

SELECT
    'One composite index' AS description,
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS indexes_size,
    pg_size_pretty(pg_total_relation_size('service')) AS total_size;

-- ================================
-- PART 4: Add covering index with INCLUDE
-- ================================

DROP INDEX idx_service_branch_status_date;

CREATE INDEX idx_service_dashboard
    ON service(branch_id, status, service_date DESC)
    INCLUDE (customer_id, mechanic_id, invoice_no);

SELECT
    'Covering index (with INCLUDE)' AS description,
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS indexes_size,
    pg_size_pretty(pg_total_relation_size('service')) AS total_size;

-- Compare individual index sizes
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size,
    pg_relation_size(indexname::regclass) AS bytes
FROM pg_indexes
WHERE tablename = 'service'
ORDER BY pg_relation_size(indexname::regclass) DESC;

-- ================================
-- PART 5: Storage summary
-- ================================

-- Comprehensive view of all index storage costs
SELECT
    'service' AS relation,
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS indexes_size,
    pg_size_pretty(pg_total_relation_size('service')) AS total_size,
    round(100.0 * pg_indexes_size('service') / pg_total_relation_size('service'), 2) AS index_percent,
    pg_size_pretty(pg_relation_size('service') * 0.01) AS one_percent_storage
FROM generate_series(1);

-- Cleanup
DROP INDEX IF EXISTS idx_service_dashboard;