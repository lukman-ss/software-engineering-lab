-- Lab 02: Database Index Foundation
-- Experiment: Functions on Indexed Columns (Expression Indexes)
-- Schema: service_date is of type DATE

-- ============================================
-- KEY CONCEPT
-- ============================================
-- Wrapping an indexed column in a function creates a new computed value.
-- PostgreSQL can only use a plain B-tree index when the predicate
-- references the indexed expression directly (same column, no transform).
-- If the predicate applies a function, the index over the raw column
-- cannot satisfy it with a bounded Index Cond on the raw DATE key.
--
-- NOT every expression predicate forces a specific plan type.
-- The planner may still choose a full index scan, a Bitmap Heap Scan,
-- or Sequential Scan depending on cost estimates.
--
-- The planner inspects whether the indexed expression matches what it needs.
-- If there is an expression index whose definition matches the predicate
-- expression, PostgreSQL CAN use that expression index.

-- ============================================
-- PART 1: Standard index on service_date
-- ============================================

DROP INDEX IF EXISTS idx_service_service_date;
CREATE INDEX idx_service_service_date ON service(service_date);
ANALYZE service;

-- ============================================
-- PART 2: Expression predicate (non-SARGable against plain index)
-- ============================================

-- This predicate applies EXTRACT() to service_date.
-- The plain service_date index cannot use
-- EXTRACT(YEAR FROM service_date) = 2026
-- as a bounded Index Cond on the raw DATE key.
-- The expression must be evaluated for candidate entries/rows because
-- the raw B-tree key cannot bound the requested year.
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE EXTRACT(YEAR FROM service_date) = 2026;

-- Observe the plan:
-- - Plain index provides no direct bound on year 2026
-- - Planner may choose Index Scan, Bitmap Heap Scan, or Seq Scan
--   based on cost estimates for this 500k-row table

-- ============================================
-- PART 3: SARGable rewrite using a range predicate
-- ============================================

-- Equivalent intent, rewritten to match the raw service_date column.
-- The index can now bound the scan to the date range directly.
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE service_date >= DATE '2026-01-01'
  AND service_date <  DATE '2027-01-01';

-- Observe: Index Scan or Bitmap Index Scan?
-- "Index Cond" should show the date range.
-- Compare shared buffers read vs Part 2.

-- SARGable = Search Argument Able.
-- A predicate is SARGable when the storage engine can use an index to
-- satisfy it directly.  A range on the raw column is SARGable;
-- EXTRACT(...) on the column is not SARGable against a plain index.

-- ============================================
-- PART 4: Expression index matching the function predicate
-- ============================================

-- When rewriting the query is impractical, an expression index stores the
-- pre-computed function result, making the EXTRACT predicate SARGable.

DROP INDEX IF EXISTS idx_service_year_expr;
CREATE INDEX idx_service_year_expr
    ON service ((EXTRACT(YEAR FROM service_date)));
ANALYZE service;

-- Re-run the expression predicate.
-- The index definition now matches the predicate.
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE EXTRACT(YEAR FROM service_date) = 2026;

-- Observe: does the plan show "Index Scan using idx_service_year_expr"?
-- PostgreSQL matches the parsed/normalized query expression against the
-- indexed expression. The query must contain an expression PostgreSQL
-- can recognize as matching the expression index definition.

-- ============================================
-- PART 5: Trade-off comparison
-- ============================================

-- Option A — Rewrite query to SARGable range predicate
--   Pros: Uses existing standard index. No extra index to maintain.
--   Cons: Requires changing application code.

-- Option B — Create expression index
--   Pros: Existing non-SARGable query gains index access; no app change.
--   Cons: One more index to maintain on every write.
--         The expression index is specific to that exact function form.
--         If the expression changes (e.g. date_part vs EXTRACT), the index
--         no longer matches.

-- General guidance: prefer the SARGable query rewrite.
-- Use an expression index only when rewriting the query is impossible or
-- carries unacceptable risk.

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_service_date;
DROP INDEX IF EXISTS idx_service_year_expr;