-- Lab 02: Database Index Foundation
-- Experiment: Composite index matching WHERE + ORDER BY

-- The composite index candidate:
-- Equality predicates first: branch_id, status
-- Range predicate after: service_date
-- DESC for ORDER BY to potentially eliminate Sort

-- Drop first if exists (safe to run multiple times)
DROP INDEX IF EXISTS idx_service_branch_status_date;

-- Create composite index with proper column order
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date DESC);

-- Important PostgreSQL 16 multicolumn B-tree behavior:
-- ===========================================
-- For an index on (a, b, c), constraints on leading columns
-- determine how much of the B-tree scan range can be bounded:
--
-- WHERE a = ?                     → constrains index range efficiently
-- WHERE a = ? AND b = ?           → constrains further
-- WHERE a = ? AND b = ? AND c > ? → bounds a tight range
--
-- WHERE b = ? or WHERE c = ?
--   PostgreSQL may technically use the index, but without a
--   constraint on the leftmost column, it may need to scan a
--   large fraction or all of the index. The planner will
--   therefore often prefer Seq Scan.
--
-- WHERE b = ? AND c = ?
--   Same reasoning applies.
--
-- PostgreSQL 16: no skip-scan optimization (that arrives in PG 18).
-- ===========================================

-- Run the main query
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record your observations:
-- Plan type: Sequential Scan / Index Scan / Bitmap Heap Scan / Index Only Scan
-- Was Sort eliminated? Yes / No
-- Rows examined: _____
-- Rows returned: _____
-- Execution Time (ms): _____
-- Buffers: read _____, hit _____

-- ===========================================
-- Understanding the column order in composite index

-- Why (branch_id, status, service_date DESC) is a good candidate:

-- 1. Equality predicates first: branch_id = 2
--    - Index can use equality condition efficiently
--    - Narrows search to specific branch

-- 2. Equality predicates second: status = 'FINISHED'
--    - Further narrows to finished services
--    - Still equality - index can use

-- 3. Range predicate last: service_date BETWEEN ...
--    - Index can use for range scan
--    - ORDER BY service_date DESC can use same index in reverse

-- The heuristic: equality < range < ORDER BY

-- BUT remember: this is a heuristic, not absolute law
-- PostgreSQL may choose differently based on statistics

-- ===========================================
-- Can PostgreSQL satisfy ORDER BY using backward scan?

-- Look in your plan for:
-- "Index Scan using idx_service_branch_status_date"
-- or "Index Only Scan using idx_service_branch_status_date"
-- followed by "Sort" being ABSENT or marked as unnecessary

-- Check: does the output come out in correct order?
-- The index stores (branch_id, status, service_date DESC)
-- So for branch_id=stable, status=stable, date DESC comes naturally

-- ===========================================
-- BACKWARD SCAN FACT (important distinction):

-- PostgreSQL B-tree indexes support both forward and backward scans.
-- With equality constraints on branch_id and status, an ASC date key
-- may also satisfy descending date order via backward scan.
--
-- Therefore, service_date DESC in the index definition
-- is not mandatory simply because the query uses ORDER BY ... DESC.
--
-- Explicit ASC/DESC definitions become more relevant for
-- mixed-order multicolumn sort requirements (e.g., ORDER BY a ASC, b DESC).

-- ===========================================
-- Compare all three testing approaches

-- Baseline (no useful index):
-- Execution Time: _____ ms
-- Rows Examined: ALL _____
-- Sort: _____

-- Three single-column indexes:
-- Execution Time: _____ ms
-- Rows Examined: _____
-- Sort: _____
-- Plan type: _____

-- One composite index:
-- Execution Time: _____ ms
-- Rows Examined: _____
-- Sort: _____
-- Plan type: _____

-- Winner: _____

-- ===========================================
-- Key PostgreSQL multicolumn B-tree facts

-- 1. Index supports queries on leftmost prefix
--    Index on (a, b, c) supports WHERE a, WHERE a AND b, WHERE a AND b AND c
--    WHERE b, WHERE c, WHERE b AND c may use index but scan large portion
--    Planner often prefers Seq Scan for these cases

-- 2. Index can help with ORDER BY if:
--    ORDER BY matches index column order
--    Index provides data in required direction

-- 3. Backward scan capability:
--    PostgreSQL can scan index backward for DESC ordering
--    No explicit Sort needed if ORDER BY uses index

-- 4. Index-only scan possible:
--    If all required columns are in index (covering index)
--    Can avoid heap tuple fetches when visibility map conditions allow it.
--    Check "Heap Fetches" in EXPLAIN ANALYZE output to verify.

-- ===========================================
-- Cleanup options (comment out if using for homework)

-- Option A: Keep composite index for queries
-- CREATE INDEX idx_service_branch_status_date
--     ON service(branch_id, status, service_date DESC);

-- Option B: Drop to restore baseline
-- DROP INDEX idx_service_branch_status_date;