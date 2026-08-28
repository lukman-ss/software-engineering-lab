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
-- --------------------------
-- Rows examined: _____
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

-----------------------------------

-- Understanding the output

-- Seq Scan: PostgreSQL reads every row in the table
-- -- Not automatically bad; good for small tables or many rows match
-- -- Cost shown in EXPLAIN (not actual time)
-- -- actual rows and rows removed by filter show the reality

-- Sort: Explicit operation when ORDER BY can't use index
-- -- Sort Method shows memory/disk usage
-- -- Could be eliminated with proper index

-- Cost: Planner's estimate (not wall clock time)
-- -- startup cost: time before first row
-- -- total cost: estimated work

-- Rows: Planner's estimate (may be wrong - check actual rows!)
-- -- Mismatch = stale statistics

-- Buffers:
-- -- shared read = disk reads
-- -- shared hit = already in cache
-- -- Fewer reads = better cache usage

-- Planning Time: Query planning duration
-- -- Not affected by table size much

-- Execution Time: Actual wall clock time (ANALYZE only)
-- -- Varies by hardware, cache, row count

-----------------------------------

-- Exercise: Identify these metrics from your plan

-- Q1: How many total rows were scanned?
-- Look for: "rows=X" under Seq Scan
-- Answer: _____

-- Q2: How many rows actually matched?
-- Look for: rows X, COUNT
-- Answer: _____

-- Q3: Was there an explicit Sort?
-- Look for: "Sort" in plan or "Sort Method"
-- Yes / No

-- Q4: What is the total execution time?
-- Look for: "Execution Time: X.XXXX ms"
-- Answer: _____ ms

-- Q5: What are estimated vs actual rows for index usage?
-- Estimated: _____
-- Actual: _____

-- Q6: How much data was read from disk vs cache?
-- shared read: _____ blocks
-- shared hit: _____ blocks
-- Total blocks read: _____