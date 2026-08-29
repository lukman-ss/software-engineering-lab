-- Lab 02: Database Index Foundation
-- Experiment: Composite Index Column Order
-- Target: PostgreSQL 16

-- ============================================
-- BACKGROUND: How PostgreSQL 16 uses a B-tree on (a, b, c)
-- ============================================

-- A B-tree index on (branch_id, status, service_date) sorts entries
-- first by branch_id, then by status within each branch_id value,
-- then by service_date within each (branch_id, status) pair.

-- The planner can navigate this structure efficiently ONLY when it can
-- bound the scan range starting from the leftmost column.

-- Leading equality predicate:
--   WHERE branch_id = 2
--   → enters the subtree for branch_id = 2; efficiently bounded.

-- Two leading equality predicates:
--   WHERE branch_id = 2 AND status = 'FINISHED'
--   → enters the subtree for (branch_id=2, status='FINISHED'); narrower bound.

-- Equality + equality + range on leading three columns:
--   WHERE branch_id = 2 AND status = 'FINISHED'
--     AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
--   → tight bounded scan; very few index pages read.

-- Non-leading predicates:
--   WHERE status = 'FINISHED'
--   WHERE service_date BETWEEN ...
--   WHERE status = 'FINISHED' AND service_date BETWEEN ...
--
--   These are NOT impossible uses of the index.  PostgreSQL 16 can
--   technically access the index for them, but without a constraint on
--   branch_id the scan must cover a large or complete range of index
--   pages.  The planner compares the estimated cost of that large index
--   scan against a sequential table scan and may prefer:
--     - Seq Scan
--     - a different, more selective index
--     - a bitmap plan
--   depending on cost estimates from statistics.
--
--   Do NOT assume the planner always rejects these cases.
--   Check the actual execution plan.

-- PostgreSQL 18 note:
--   B-tree skip-scan (introduced in PG 18) can jump over leading-column
--   values to satisfy non-leading predicates more efficiently.
--   That optimization does NOT exist in PostgreSQL 16.

-- ============================================
-- SETUP: Three indexes with different column orders
-- ============================================

DROP INDEX IF EXISTS idx_service_a_branch_status_date;
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
DROP INDEX IF EXISTS idx_service_c_status_date_branch;

-- Index A: leading equality columns, then range column
--   Optimal order for: WHERE branch_id = ? AND status = ? AND service_date BETWEEN ...
CREATE INDEX idx_service_a_branch_status_date
    ON service(branch_id, status, service_date);

-- Index B: range column first
--   Leading column is a high-cardinality date; branch_id and status are non-leading
CREATE INDEX idx_service_b_date_branch_status
    ON service(service_date, branch_id, status);

-- Index C: status first, then date, then branch_id
--   Leading column is low-cardinality status
CREATE INDEX idx_service_c_status_date_branch
    ON service(status, service_date, branch_id);

ANALYZE service;

-- ============================================
-- HOW TO READ THE OUTPUT BELOW
-- ============================================
-- For each EXPLAIN block, note:
--   1. Which index name appears in the plan (Index Scan / Bitmap Index Scan)?
--   2. "Index Cond" lines  → predicates the index resolved directly
--   3. "Filter" lines      → predicates applied AFTER index access (heap rows examined)
--   4. "Rows Removed by Filter" → work done that an index didn't prevent
--   5. "Buffers: shared hit/read" → pages touched (fewer = better)
--   6. "Sort" node present? → the index did not supply the ORDER BY order
--   7. actual rows vs estimated rows (plan quality indicator)

-- ============================================
-- QUERY 1: Optimal for Index A
-- All three columns bounded; tight range
-- ============================================

-- Question: which index does the planner select?
-- Does it use all three Index Cond columns or does one appear as Filter?
-- How many pages (shared hit) does each plan touch?

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Expected: Index A (branch_id, status, service_date) with a tight index range.
-- branch_id = 2 AND status = 'FINISHED' become Index Cond;
-- service_date BETWEEN ... also becomes Index Cond (range).
-- With all three columns as Index Cond and DESC in the query matching
-- a backward scan, a Sort node may be absent.

-- ============================================
-- QUERY 2: Leading equality only (branch_id)
-- ============================================

-- With only branch_id constrained:
--   Index A: leading column — efficient; scans branch_id=2 subtree (~25% of rows)
--   Index B: branch_id is column 2 — non-leading; scan may be large
--   Index C: branch_id is column 3 — non-leading; scan may be large

-- The planner will pick the cheapest option.
-- It might choose Index A, Seq Scan, or a bitmap plan.
-- It will NOT necessarily choose Index B or C.

EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*) FROM service WHERE branch_id = 2;

-- Observe: which index was selected (if any)?
-- Was it Index A, B, C, or Seq Scan?
-- How many pages were touched?

-- ============================================
-- QUERY 3: Non-leading predicate only (status)
-- ============================================

-- status is:
--   column 2 in Index A (non-leading without branch_id)
--   column 2 in Index B (non-leading)
--   column 1 in Index C (leading — most relevant)

-- With only status = 'FINISHED' (70% of rows):
--   The match fraction is large.  Even Index C (which leads on status)
--   must scan roughly 70% of its entries.
--   The planner may prefer Seq Scan over any index.

EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*) FROM service WHERE status = 'FINISHED';

-- Observe: Seq Scan or Index Scan?
-- If Seq Scan: note that 70% match fraction makes index traversal more
--   expensive than sequential table access.
-- If Index C is chosen: compare shared hit count against Seq Scan cost.

-- ============================================
-- QUERY 4: Non-leading predicate only (service_date range)
-- ============================================

-- service_date is:
--   column 3 in Index A (non-leading)
--   column 1 in Index B (leading — most relevant)
--   column 2 in Index C (non-leading)

-- One month of data out of 730 days ≈ 4.1% match fraction.
-- Index B leads on service_date, so it can bound this range tightly.

EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*) FROM service
WHERE service_date BETWEEN '2026-01-01' AND '2026-01-31';

-- Observe: which index was selected?
-- Was it Index B (service_date leading)?
-- Alternatively: Bitmap Heap Scan or Seq Scan?
-- How many pages did it read vs a Seq Scan would read?

-- ============================================
-- QUERY 5: Two non-leading predicates (status + service_date)
-- No branch_id constraint
-- ============================================

-- For Index A: branch_id is missing; status and service_date are non-leading
-- For Index B: branch_id and status are non-leading after service_date
-- For Index C: status is leading, service_date is next

-- With status = 'PENDING_REFUND' (0.1% match fraction) + date range:
-- The combined predicate is highly selective.

EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*) FROM service
WHERE status = 'PENDING_REFUND'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31';

-- Observe: which index, if any, was selected?
-- What appears as Index Cond vs Filter?
-- A bitmap plan or Seq Scan is also possible.

-- ============================================
-- QUERY 6: ORDER BY — can the index avoid a Sort?
-- ============================================

-- ORDER BY service_date DESC with WHERE branch_id = 2 AND status = 'FINISHED':
-- Index A stores (branch_id, status, service_date ASC).
-- Within branch_id=2 and status='FINISHED', entries are ordered by service_date.
-- PostgreSQL can satisfy DESC by scanning the index backward — no Sort node needed.

EXPLAIN (ANALYZE, BUFFERS)
SELECT id, service_date, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- Observe: Is there a Sort node?
-- "Index Scan Backward" or "Index Scan" with Backward in the plan?
-- Rows examined before LIMIT stops it?

-- ============================================
-- SUMMARY: Fill in your observations
-- ============================================

-- Query | Predicate columns        | Index selected | Index Cond cols | Filter cols | Sort? | Pages read
-- ------|--------------------------|----------------|-----------------|-------------|-------|----------
-- Q1    | branch_id+status+date    | ___________    | ___________     | ___________  | _____ | _____
-- Q2    | branch_id only           | ___________    | ___________     | ___________  | _____ | _____
-- Q3    | status only (70%)        | ___________    | ___________     | ___________  | _____ | _____
-- Q4    | service_date range only  | ___________    | ___________     | ___________  | _____ | _____
-- Q5    | status+date (selective)  | ___________    | ___________     | ___________  | _____ | _____
-- Q6    | branch_id+status+LIMIT   | ___________    | ___________     | ___________  | _____ | _____

-- Key findings to derive from your observations:
--
-- 1. Does Index A always win? Or does the planner sometimes choose B or C?
--
-- 2. For Q3 (status = 'FINISHED', 70% match): did the planner avoid all indexes?
--    Why or why not?
--
-- 3. For Q4 (date range ~4%): did Index B (date leading) get selected?
--    What does this tell you about when a non-equality leading column is useful?
--
-- 4. For Q5 (PENDING_REFUND + date): very selective combined predicate.
--    Which index helped? What appears as Index Cond vs Filter?
--
-- 5. For Q6: was a Sort node absent? What does that tell you about
--    the cost benefit of matching index order to ORDER BY?

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_a_branch_status_date;
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
DROP INDEX IF EXISTS idx_service_c_status_date_branch;
