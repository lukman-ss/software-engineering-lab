-- Lab 02: Database Index Foundation
-- Experiment 17: Plan Comparison Harness
--
-- PURPOSE: Compare execution plans across four index strategies against the
-- same canonical query and dataset. Record plan type, buffers, rows, Sort
-- presence, planning time, and execution time for each scenario.
--
-- This is a plan comparison, not a repeated statistical benchmark.
-- One EXPLAIN (ANALYZE, BUFFERS) per scenario is run.

-- ============================================
-- INSTRUMENTATION NOTE
-- ============================================
-- EXPLAIN ANALYZE executes the query on the PostgreSQL server and reports
-- planning/execution instrumentation (plan nodes, rows, buffers, timings).
-- It does not measure terminal rendering of the result set in the same way
-- as running SELECT and printing all rows to a client.
-- The instrumentation itself adds overhead — do not treat execution time
-- figures as identical to production query latency.
--
-- TIMING OFF disables per-node timing instrumentation but still executes
-- the query and reports planning time and total execution time.
-- Use it when per-node overhead would distort measurements.
--
-- WARM-CACHE CAUTION:
-- Once PostgreSQL has read pages into shared_buffers, re-runs are faster.
-- We keep the cache warm here for reproducible relative comparisons.
-- Cold-cache testing requires pg_prewarm or a full server restart — not
-- done here as results are environment-dependent.

-- ============================================
-- HOW TO READ THE PLAN NODES
-- ============================================
-- Index Scan:
--   Navigates the B-tree index, then fetches each matching heap tuple.
--
-- Bitmap Heap Scan:
--   Fetches heap pages/tuples using a bitmap built in a prior step.
--   The bitmap may come from a single Bitmap Index Scan, BitmapAnd
--   (intersection of multiple bitmaps), or BitmapOr (union).
--   It is not exclusively a "multiple index" operation.
--
-- Seq Scan:
--   Reads the table sequentially. Often the correct choice when a large
--   fraction of rows match — not an indicator of a missing index.
--
-- Index Only Scan:
--   All needed columns are in the index; heap access is minimized.
--   Requires the visibility map to confirm row visibility (check Heap Fetches).
--
-- Sort:
--   Explicit sort step. Absent when the chosen index already supplies the
--   ORDER BY order (forward or backward scan of the B-tree).

-- ============================================
-- SETUP: Ensure clean index state for fair comparison
-- ============================================

DROP INDEX IF EXISTS idx_bench_branch_id;
DROP INDEX IF EXISTS idx_bench_status;
DROP INDEX IF EXISTS idx_bench_service_date;
DROP INDEX IF EXISTS idx_bench_composite_wrong;
DROP INDEX IF EXISTS idx_bench_composite_correct;
DROP INDEX IF EXISTS idx_bench_covering;

ANALYZE service;

-- ============================================
-- SCENARIO A: No Secondary Index (Baseline)
-- ============================================

-- Expected/likely: Seq Scan + Sort.
-- Without an index the planner reads the whole table and then sorts.
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record actual observed plan: _____
-- Buffers: shared read _____ / shared hit _____

-- ============================================
-- SCENARIO B: Three Single-Column Indexes
-- ============================================

CREATE INDEX idx_bench_branch_id    ON service(branch_id);
CREATE INDEX idx_bench_status       ON service(status);
CREATE INDEX idx_bench_service_date ON service(service_date);

ANALYZE service;

-- Expected/likely: Bitmap Heap Scan driven by one or more Bitmap Index Scans.
-- The planner may combine bitmaps from the three indexes via BitmapAnd, or
-- may choose just the most selective single index.
-- Actual plan depends on cost estimates from current statistics.
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record actual observed plan: _____
-- Buffers: shared read _____ / shared hit _____

DROP INDEX idx_bench_branch_id;
DROP INDEX idx_bench_status;
DROP INDEX idx_bench_service_date;
DROP INDEX IF EXISTS idx_service_status;

-- ============================================
-- SCENARIO C: Composite Index — range column first (suboptimal order)
-- ============================================

CREATE INDEX idx_bench_composite_wrong
    ON service(service_date, status, branch_id);

ANALYZE service;

-- Expected/likely: service_date, as the leading range column, determines
-- the B-tree range that must be scanned. status and branch_id can still be
-- checked using index entries and may appear in Index Cond, but because they
-- occur to the right of the first range-constrained key, they generally
-- do not further reduce the portion of the index that has to be scanned.
-- A Sort node may still be present depending on the planner's choice.
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record actual observed plan: _____
-- Buffers: shared read _____ / shared hit _____
-- (Whether predicates appear as Index Cond or Filter must be observed
-- from the actual EXPLAIN output, not predicted.)

DROP INDEX idx_bench_composite_wrong;

-- ============================================
-- SCENARIO D: Recommended Composite Index
-- ============================================

CREATE INDEX idx_bench_composite_correct
    ON service(branch_id, status, service_date DESC);

ANALYZE service;

-- Expected/likely: the recommended composite index can tightly bound the
-- B-tree range on all three columns (branch_id equality → status equality
-- → service_date range) and may also supply the ORDER BY order, eliminating
-- an explicit Sort.  The planner may still choose a different plan if cost
-- estimates favour it.
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record actual observed plan: _____
-- Buffers: shared read _____ / shared hit _____
-- Sort node present? _____

DROP INDEX idx_bench_composite_correct;

-- ============================================
-- COMPARISON TABLE (fill in after running all scenarios)
-- ============================================
--
-- | Scenario                | Expected/likely plan     | Actual observed plan | Buffers (R/H) | Sort? |
-- |-------------------------|--------------------------|----------------------|---------------|-------|
-- | A: No index             | Seq Scan + Sort          | _____                | _____/_____   | _     |
-- | B: Single-column x3     | Bitmap Heap Scan         | _____                | _____/_____   | _     |
-- | C: Wrong composite      | _____ (observed)         | _____                | _____/_____   | _     |
-- | D: Recommended composite| Index Scan (tight range) | _____                | _____/_____   | _     |
--
-- The "expected/likely plan" column is based on cost reasoning and this
-- dataset's distribution.  The planner's actual choice may differ.
-- Always trust the actual observed plan over the prediction.
--
-- KEY THINGS TO COMPARE:
-- 1. Shared buffers read: fewer = better index range bounding
-- 2. "Rows Removed by Filter" in Index Scan nodes: non-zero = predicate
--    was not an Index Cond (index did not fully bound that predicate)
-- 3. Sort node: present = index did not supply ORDER BY order
-- 4. Actual rows vs estimated rows: large gap = stale statistics

-- ============================================
-- NOTE: Covering indexes
-- ============================================
-- Covering indexes (INCLUDE columns → Index Only Scan) are covered
-- in 07-covering-index.sql.  Mixing SELECT * benchmark results with
-- covering-index queries on different column sets is not apples-to-apples.
