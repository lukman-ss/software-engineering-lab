-- Lab 02: Database Index Foundation
-- Experiment: When Seq Scan IS Correct
-- Destroys the misconception: "Seq Scan = Bad, Index Scan = Good"

-- ============================================
-- PRIMARY LESSON
-- ============================================
-- The goal is NOT to force index usage.
-- The goal is to MINIMIZE TOTAL QUERY COST for the real workload.
--
-- Seq Scan is not inherently bad.
-- Index Scan is not inherently good.
-- The planner chooses based on cost estimates derived from statistics.

-- ============================================
-- PART 1: Ensure Index Exists
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
CREATE INDEX idx_service_status ON service(status);
ANALYZE service;

-- Sanity check: total rows and distribution
SELECT count(*) AS total_rows FROM service;

SELECT status, count(*) AS cnt,
       round(count(*) * 100.0 / (SELECT count(*) FROM service), 2) AS percentage
FROM service GROUP BY status ORDER BY cnt DESC;

-- ============================================
-- PART 2: High Match Fraction Query (70%)
-- Seq Scan is expected/likely under this generated distribution
-- ============================================

-- status = 'FINISHED' matches ~350,000 of 500,000 rows (70% match fraction).
-- High match fraction → not very selective → index traversal likely more
-- expensive than a sequential table scan.
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';

-- EXPECTED / LIKELY: Seq Scan.
-- REASON: ~70% of heap pages must be visited regardless.
-- Reading the index and then fetching most heap pages would cost more,
-- not less. The planner's cost model reflects this.

-- ============================================
-- PART 3: Low Match Fraction Query (0.1%)
-- An index-based plan is expected/likely for the rare predicate
-- ============================================

-- status = 'PENDING_REFUND' matches ~500 of 500,000 rows (0.1% match fraction).
-- Low match fraction → highly selective → index traversal much cheaper.
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'PENDING_REFUND';

-- EXPECTED / LIKELY: Index Scan or Bitmap Index Scan.
-- REASON: The index narrows the search to ~500 rows immediately.
-- Far fewer heap pages are touched compared to a full table scan.

-- ============================================
-- PART 4: Planner-Switch Demonstration (Educational Only)
-- ============================================

-- WARNING: Planner switches are diagnostic tools, NOT production tuning.
-- They adjust the planner's cost factors, they do NOT force a specific
-- access path. The planner always picks the cheapest available plan
-- under the adjusted cost model; it may still choose unexpectedly.

-- Baseline: confirm the planner's natural choice for the high-match query.
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';

-- Wrap SET LOCAL in a transaction so the setting cannot leak to subsequent
-- queries. ROLLBACK restores the session to its previous state.
BEGIN;

SET LOCAL enable_seqscan = off;
-- With sequential plans discouraged, the planner looks for an index-based
-- alternative. It may still choose Bitmap Heap Scan rather than a plain
-- Index Scan, and in rare cases may re-evaluate Seq Scan if that remains
-- the only viable option.
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';

ROLLBACK;
-- The enable_seqscan setting is now fully restored.

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
