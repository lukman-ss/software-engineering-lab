-- Lab 02: Database Index Foundation
-- Experiment: Partial Indexes for Subset Optimization
-- Teaches when and why partial indexes are useful

-- PARTIAL INDEX DEFINITION:
-- CREATE INDEX name ON table(columns) WHERE predicate;
-- Only stores rows that satisfy the predicate. Smaller, faster to maintain.

================================
-- PART 1: Create Partial Index for Active/FINISHED Services
================================

-- Scenario: Dashboard frequently shows "active or recently finished" services
-- A partial index for FINISHED status can be very effective

DROP INDEX IF EXISTS idx_service_partial_finished;

-- Partial index: only index rows where status = 'FINISHED'
CREATE INDEX idx_service_partial_finished
    ON service(branch_id, service_date DESC)
    WHERE status = 'FINISHED';

-- Compare size with a full index for the same columns
CREATE INDEX idx_service_full_branch_date
    ON service(branch_id, service_date DESC);

SELECT
    'Partial Index (FINISHED only)' AS index_type,
    pg_size_pretty(pg_relation_size('idx_service_partial_finished')) AS size
UNION ALL
SELECT
    'Full Index (all rows)',
    pg_size_pretty(pg_relation_size('idx_service_full_branch_date'));

================================
-- PART 2: Verify Partial Index Can Be Used
================================

-- Query that matches the partial index predicate exactly
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- Look for: "Index Scan using idx_service_partial_finished"
-- The index is used because the WHERE clause matches the partial predicate

================================
-- PART 3: Query That CANNOT Use Partial Index
================================

-- Query for CANCELLED status - partial index only covers FINISHED!
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'CANCELLED'
ORDER BY service_date DESC
LIMIT 20;

-- This query does NOT use idx_service_partial_finished!
-- PostgreSQL must either:
-- 1. Do a Seq Scan, or
-- 2. Use a different index if available

-- This demonstrates: Partial indexes are NOT universal.

================================
-- PART 4: Compare Full Index Plan
================================

-- Same query, but with full index available
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'CANCELLED'
ORDER BY service_date DESC
LIMIT 20;

-- Full index can still help, but may not be as efficient as partial for FINISHED

================================
-- PART 5: Write Overhead Implications
================================

-- INSERT a FINISHED service: goes into partial index
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
VALUES (2, 999, 100, 'FINISHED', '2026-07-01', 'PARTIAL-TEST', NOW());

-- INSERT a CANCELLED service: does NOT go into partial index (cheaper!)
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
VALUES (2, 998, 100, 'CANCELLED', '2026-07-01', 'PARTIAL-TEST-2', NOW());

-- This means:
-- - INSERT of FINISHED row: must update partial index + any other indexes
-- - INSERT of CANCELLED row: does NOT touch partial index (faster!)
-- Partial indexes are often MORE efficient for writes of non-matching rows!

-- Cleanup
DELETE FROM service WHERE invoice_no LIKE 'PARTIAL-%';
DROP INDEX IF EXISTS idx_service_partial_finished;
DROP INDEX IF EXISTS idx_service_full_branch_date;