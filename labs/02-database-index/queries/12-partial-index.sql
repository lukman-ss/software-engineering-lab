-- Lab 02: Database Index Foundation
-- Experiment: Partial Indexes for Subset Optimization
-- Teaches when and why partial indexes are useful

-- PARTIAL INDEX DEFINITION:
-- CREATE INDEX name ON table(columns) WHERE predicate;
-- Only stores rows that satisfy the predicate. Smaller, faster to maintain.

-- ============================================
-- PART 1: Create Partial Index for IN_PROGRESS
-- ============================================

-- Scenario: Dashboard frequently shows "in-progress" services
-- IN_PROGRESS is only ~5% of rows, so a partial index is very small and efficient

DROP INDEX IF EXISTS idx_service_partial_in_progress;
DROP INDEX IF EXISTS idx_service_full_status_branch_date;

-- Partial index: only index rows where status = 'IN_PROGRESS'
CREATE INDEX idx_service_partial_in_progress
    ON service(branch_id, service_date DESC)
    WHERE status = 'IN_PROGRESS';

-- Full index for comparison (covers all statuses)
CREATE INDEX idx_service_full_status_branch_date
    ON service(status, branch_id, service_date DESC);

-- Compare sizes
SELECT
    'Partial Index (IN_PROGRESS only)' AS index_type,
    pg_size_pretty(pg_relation_size('idx_service_partial_in_progress')) AS size
UNION ALL
SELECT
    'Full Index (all rows)',
    pg_size_pretty(pg_relation_size('idx_service_full_status_branch_date'));

-- ============================================
-- PART 2: Verify Partial Index Can Be Used
-- ============================================

-- Query that matches the partial index predicate exactly
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'IN_PROGRESS'
ORDER BY service_date DESC
LIMIT 20;

-- Look for: "Index Scan using idx_service_partial_in_progress"
-- The index is used because the WHERE clause matches the partial predicate

-- ============================================
-- PART 3: Query That CANNOT Use Partial Index
-- ============================================

-- Query for FINISHED status - partial index only covers IN_PROGRESS!
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- This query does NOT use idx_service_partial_in_progress!
-- PostgreSQL must either:
-- 1. Do a Seq Scan, or
-- 2. Use a different index if available (idx_service_full_status_branch_date)

-- ============================================
-- PART 4: Compare Full Index Plan for FINISHED
-- ============================================

-- Same query, but with full index available
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- The full index can help, but it's larger to maintain

-- ============================================
-- PART 5: Write Overhead Implications
-- ============================================

-- INSERT a FINISHED service: does NOT go into partial index (faster!)
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
VALUES (2, 999, 100, 'FINISHED', '2026-07-01', 'PARTIAL-TEST', NOW());

-- INSERT an IN_PROGRESS service: goes into partial index
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
VALUES (2, 998, 100, 'IN_PROGRESS', '2026-07-01', 'PARTIAL-TEST-2', NOW());

-- This means:
-- - INSERT of FINISHED row: does NOT touch partial index (faster!)
-- - INSERT of IN_PROGRESS row: must update partial index
-- Partial indexes are MORE efficient for writes of non-matching rows!

-- Cleanup
DELETE FROM service WHERE invoice_no LIKE 'PARTIAL-%';
DROP INDEX IF EXISTS idx_service_partial_in_progress;
DROP INDEX IF EXISTS idx_service_full_status_branch_date;