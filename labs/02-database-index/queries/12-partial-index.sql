-- Lab 02: Database Index Foundation
-- Experiment: Partial Indexes for Subset Optimization
-- Teaches when and why partial indexes are useful

-- ============================================
-- PARTIAL INDEX DEFINITION:
-- ============================================
-- CREATE INDEX name ON table(columns) WHERE predicate;
-- Only stores rows that satisfy the predicate. Smaller, faster to maintain.
--
-- Partial indexes are most beneficial when:
-- 1. Queries frequently target a stable, predictable predicate
-- 2. Indexed subset is significantly smaller than full table
-- 3. Query predicate implies the partial-index predicate
--    (The query's WHERE clause restricts to a subset covered by the index)

-- ============================================
-- PART 1: Full Index vs Partial Index Comparison
-- ============================================
-- Use IN_PROGRESS (~5% of rows) because it's a small, stable subset

DROP INDEX IF EXISTS idx_service_in_progress_branch_date;
DROP INDEX IF EXISTS idx_service_full_status_branch_date;

-- Full index: covers all status values
CREATE INDEX idx_service_full_status_branch_date
    ON service(status, branch_id, service_date DESC);

-- Partial index: ONLY IN_PROGRESS rows (~5% of table)
-- This index is much smaller since it excludes 95% of rows
CREATE INDEX idx_service_in_progress_branch_date
    ON service(branch_id, service_date DESC)
    WHERE status = 'IN_PROGRESS';

-- Compare sizes
SELECT
    pg_size_pretty(pg_relation_size('idx_service_full_status_branch_date')) AS full_index_size,
    pg_size_pretty(pg_relation_size('idx_service_in_progress_branch_date')) AS partial_index_size;

-- ============================================
-- PART 2: Query That CAN Use Partial Index
-- ============================================
-- Query for IN_PROGRESS status matches the partial index predicate

EXPLAIN (ANALYZE, BUFFERS)
SELECT id,
       customer_id,
       mechanic_id,
       service_date,
       invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'IN_PROGRESS'
ORDER BY service_date DESC
LIMIT 20;

-- Look for: "Index Scan using idx_service_in_progress_branch_date"
-- The index is used because:
-- 1. branch_id = 2 can use leading column
-- 2. status = 'IN_PROGRESS' matches the partial predicate exactly
-- 3. ORDER BY service_date DESC satisfied by index order

-- ============================================
-- PART 3: Query That CANNOT Use Partial Index
-- ============================================
-- Query for FINISHED status does NOT match the partial predicate

EXPLAIN (ANALYZE, BUFFERS)
SELECT id,
       customer_id,
       mechanic_id,
       service_date,
       invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- This query does NOT use idx_service_in_progress_branch_date!
-- Reason: The query's WHERE status = 'FINISHED' does NOT imply/overlap
-- with the partial index's WHERE status = 'IN_PROGRESS' predicate.
--
-- PostgreSQL's predicate implication checking requires the query
-- predicate to be a subset of or equal to the partial predicate.
-- FINISHED and IN_PROGRESS are mutually exclusive, so no implication.

-- Planner will use idx_service_full_status_branch_date instead
-- or fall back to Seq Scan if full index is less efficient

-- ============================================
-- PART 4: Write Overhead Comparison
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

-- ============================================
-- PART 5: Planner Limitation Note
-- ============================================
-- For parameterized queries (prepared statements), the planner cannot
-- prove at planning time that a parameter will always match the
-- partial index predicate. This can prevent use of the partial index
-- or cause inconsistent plan choices.
--
-- Example: A prepared query with $1 = status parameter:
--   SELECT ... WHERE status = $1 AND branch_id = 2
-- PostgreSQL cannot prove at plan time that $1 = 'IN_PROGRESS'
-- Therefore it may not use the partial index reliably.
--
-- For such queries, consider using the full index instead.

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_in_progress_branch_date;
DROP INDEX IF EXISTS idx_service_full_status_branch_date;