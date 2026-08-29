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
-- cannot satisfy it — the planner must evaluate the function for every row.
--
-- NOT every function on every indexed column automatically forces a Seq Scan.
-- The planner inspects whether the indexed expression matches what it needs.
-- If there is an expression index whose definition matches the predicate
-- expression exactly, PostgreSQL CAN use that expression index.

-- ============================================
-- PART 1: Standard index on service_date
-- ============================================

DROP INDEX IF EXISTS idx_service_service_date;
CREATE INDEX idx_service_service_date ON service(service_date);
ANALYZE service;

-- ============================================
-- PART 2: Expression predicate (non-SARGable)
-- ============================================

-- This predicate applies EXTRACT() to service_date.
-- The index stores raw DATE values; it has no entry for year 2026 as
-- a searchable key.  The planner must call EXTRACT for every row.
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE EXTRACT(YEAR FROM service_date) = 2026;

-- Observe: is this a Seq Scan or Bitmap Heap Scan?
-- Note "Rows Removed by Filter" — the index provides no help.

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
-- The index definition now matches the predicate exactly.
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE EXTRACT(YEAR FROM service_date) = 2026;

-- Observe: does the plan now show "Index Scan using idx_service_year_expr"?
-- The expression in the CREATE INDEX must match the predicate expression
-- character-for-character as PostgreSQL normalises it.  A mismatch (e.g.
-- EXTRACT(year ...) vs EXTRACT(YEAR ...)) can prevent matching.

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
