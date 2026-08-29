-- Lab 02: Database Index Foundation
-- Baseline experiment: measuring query BEFORE any index optimization

-- Run EXPLAIN ANALYZE to see actual execution
-- BUFFERS shows page cache activity (shared hit vs shared read)
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record your observations:
-- rows examined: _____
-- Rows returned: _____
-- Rows Removed by Filter: _____
-- Explicit Sort needed? YES / NO
-- Execution Time (ms): _____
-- Planning Time (ms): _____
-- Buffers: _____ reads, _____ hits
--
-- Plan nodes observed:
-- 1. Seq Scan / Index Scan / Bitmap Heap Scan: _____
-- 2. Sort / Sort Method: _____
--
-- Estimated vs Actual rows:
-- - Estimated rows: _____
-- - Actual rows: _____
--
-- What does this tell you about the default behavior?
-- The database is doing a ____________ which ____________ (is/is not) efficient
-- because ____________.

-- ============================================
-- Understanding the output
-- ============================================

-- Seq Scan: PostgreSQL reads every row in the table
-- -- Not automatically bad; good for small tables or many rows match
-- -- Cost shown in EXPLAIN (not actual time)
-- -- actual rows and rows removed by filter show the reality

-- "rows=X" in the plan output (without "actual") is a planner estimate.
-- "actual rows=X" shows how many rows the node produced per loop.

-- For a simple Seq Scan, approximate total rows examined:
--   (actual rows + "Rows Removed by Filter") * loops

-- Rows Removed by Filter:
-- -- Rows the node examined but rejected by its filter condition
-- -- For Seq Scan: shows how many rows were scanned but filtered out

-- Sort: Explicit operation when ORDER BY can't use index
-- -- Sort Method shows memory/disk usage
-- -- Could be eliminated with proper index

-- Cost: Planner's estimate (not wall clock time)
-- -- startup cost: time before first row
-- -- total cost: estimated work

-- Buffer meaning:
-- -- shared read = PostgreSQL read block into shared buffers
-- -- (Note: OS page cache may satisfy the read, so shared read does
--   not necessarily mean physical disk I/O)
-- -- shared hit = block was already in PostgreSQL shared buffers
-- -- Fewer reads is better cache usage

-- Planning Time: Query planning duration
-- -- Planning time usually depends more on query complexity, available
--   paths, partitions, indexes, and planner work than on table rows,
--   because planning does not execute the table scan.

-- Execution Time: Actual wall clock time (ANALYZE only)
-- -- Varies by hardware, cache, row count

-- ============================================
-- Exercise: Identify these metrics from your plan
-- ============================================

-- Q1: How many total rows were examined?
-- Look for: "actual rows=X" + "Rows Removed by Filter" under Seq Scan
--   Approximate total: (actual rows + Rows Removed by Filter) * loops
-- Answer: _____

-- Q2: How many rows actually matched all predicates?
-- Look for: "actual rows=X" in the Seq Scan node
-- Answer: _____

-- Q3: Was there an explicit Sort?
-- Look for: "Sort" in plan or "Sort Method"
-- Yes / No

-- Q4: What is the total execution time?
-- Look for: "Execution Time: X.XXXX ms"
-- Answer: _____ ms

-- Q5: What are estimated vs actual rows?
-- Estimated: _____
-- Actual: _____

-- Q6: How many blocks were already in PostgreSQL shared buffers,
--      and how many blocks PostgreSQL had to read into shared buffers?
-- shared hit: _____ blocks
-- shared read: _____ blocks
--
-- Note: shared read does not prove physical disk I/O because the OS
-- page cache may satisfy the read.