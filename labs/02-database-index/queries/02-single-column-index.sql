-- Lab 02: Database Index Foundation
-- Experiment: Three independent single-column indexes

-- Create individual indexes for each predicate column
-- These may be combined by PostgreSQL using Bitmap operations

-- Drop first if exists (safe to run multiple times)
DROP INDEX IF EXISTS idx_service_branch_id;
DROP INDEX IF EXISTS idx_service_status;
DROP INDEX IF EXISTS idx_service_service_date;

-- Create the three single-column indexes
CREATE INDEX idx_service_branch_id
    ON service(branch_id);

CREATE INDEX idx_service_status
    ON service(status);

CREATE INDEX idx_service_service_date
    ON service(service_date);

-- Run the main query and observe what plan PostgreSQL chooses
-- It may use:
-- - One index (whichever is most selective)
-- - Multiple indexes (BitmapAnd)
-- - Sequential scan (if indexes not helpful)
-- - Bitmap Heap Scan (combining bitmap from multiple indexes)
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record your observations:
-- Plan type: Sequential Scan / Index Scan / Bitmap Heap Scan
-- Which indexes were used? _____
-- Rows examined: _____
-- Rows returned: _____
-- Execution Time (ms): _____
-- Buffers: read _____, hit _____
-- Sort needed? Yes / No

-----------------------------------

-- Understanding PostgreSQL's index combination

-- PostgreSQL CAN combine indexes via BitmapAnd when multiple are useful:
-- 1. BitmapIndexScan on idx_service_branch_id
-- 2. BitmapIndexScan on idx_service_status
-- 3. BitmapIndexScan on idx_service_service_date
-- 4. Bitmap Heap Scan to fetch actual rows

-- But PostgreSQL may also:
-- - Use just one index if others are not selective
-- - Choose Seq Scan if table is small or most rows match
-- - Not use index at all for ORDER BY

-- The planner chooses the cheapest plan based on statistics

-----------------------------------

-- Compare with baseline
-- Baseline execution time: _____ ms
-- This execution time: _____ ms
-- Improvement: _____ %

-- Key insight: Multiple single-column indexes
-- DO NOT automatically mean better performance
-- PostgreSQL uses what's cheapest based on selectivity

-----------------------------------

-- Cleanup: Drop these indexes to prepare for composite index test
DROP INDEX IF EXISTS idx_service_branch_id;
DROP INDEX IF EXISTS idx_service_status;
DROP INDEX IF EXISTS idx_service_service_date;