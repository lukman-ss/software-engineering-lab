-- Lab 02: Database Index Foundation
-- Experiment: Low Cardinality vs Selectivity
-- Corrects the oversimplified "don't index low-cardinality columns" rule

-- ============================================
-- KEY TERMINOLOGY (used consistently in this lab)
-- ============================================
-- cardinality      = number of distinct values in a column
-- match fraction    = matching_rows / total_rows
-- small match fraction   = highly selective predicate
-- large match fraction   = not very selective predicate
--
-- NOTE: match fraction is NOT "fraction filtered out".
-- When 70% of rows match, the match fraction is 0.70 (not 0.30).

-- Deterministic dataset facts (see seed.sql):
--   status distribution: FINISHED 70.0%, CANCELLED 20.0%, IN_PROGRESS 5.0%,
--                        WAITING 4.9%, PENDING_REFUND 0.1%
--   status cardinality  = 5 distinct values (low cardinality)
--   branch 2 fraction    = 25.0% of all rows

-- pg_stats.n_distinct semantics (PostgreSQL):
--   n_distinct > 0  → estimated absolute number of distinct values
--   n_distinct < 0  → negative ratio; |n_distinct| = distinct_rows / total_rows
--   n_distinct = -1   → approximately one distinct value per row
--   n_distinct = -0.5 → approximately half as many distinct values as rows

-- ============================================
-- SETUP: Build the SAME index for both comparisons
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
CREATE INDEX idx_service_status ON service(status);
ANALYZE service;

-- ============================================
-- SCENARIO A: status = 'FINISHED' (70% match fraction)
-- Large match fraction → not very selective
-- ============================================

-- Confirm the actual distribution in this dataset:
SELECT status, round(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM service), 2) AS percentage
FROM service
GROUP BY status
ORDER BY percentage DESC;

-- Now inspect the plan with the SAME physical index in place:
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE status = 'FINISHED';

-- EXPECTED / LIKELY under this generated distribution:
-- Seq Scan is typically cheaper than an index scan because ~350,000 rows
-- (70%) match.  The planner would have to traverse most of the index and
-- then visit most heap pages anyway.
-- This is NOT because `status` has low cardinality (5 values) — it is
-- because the match fraction for this specific predicate is large.

-- ============================================
-- SCENARIO B: status = 'PENDING_REFUND' (0.1% match fraction)
-- Small match fraction → highly selective
-- ============================================

-- The same index, the same column, but a rare value:
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE status = 'PENDING_REFUND';

-- EXPECTED / LIKELY under this generated distribution:
-- Index Scan is typically cheaper. Only ~500 rows match (0.1% match
-- fraction), so the index lookup visits far fewer pages than a Seq Scan.
-- Low cardinality + highly selective predicate = index is attractive.

-- ============================================
-- CORE LESSON: Same table, same column, same physical index,
-- different match fraction → potentially different plan
-- ============================================
--   status = 'FINISHED'  → 70.0% match → Seq Scan likely cheaper
--   status = 'PENDING_REFUND' → 0.1% match → Index Scan likely cheaper
--
-- The decision is driven by match fraction, not by the bare cardinality
-- of the status column.

-- ============================================
-- INSPECT n_distinct FOR status
-- ============================================

SELECT attname, n_distinct, most_common_vals, most_common_freqs
FROM pg_stats
WHERE tablename = 'service' AND attname = 'status';

-- With exactly 5 distinct status values, expect n_distinct ≈ 5 (positive).

-- ============================================
-- COMBINED PREDICATES: Do NOT assume independence
-- ============================================
-- It is tempting to multiply match fractions:
--   status = 'FINISHED' (0.70) × branch_id = 2 (0.25) = 0.175
-- That arithmetic assumes statistical independence.
-- This deterministic dataset derives branch_id and status from RELATED
-- modulo ranges of the same gs counter, so they are NOT independent.
-- Measure the ACTUAL joint distribution instead:

SELECT branch_id, status, count(*)
FROM service
WHERE status IN ('FINISHED', 'PENDING_REFUND')
GROUP BY branch_id, status
ORDER BY branch_id, status;

-- Use the real counts above to reason about combined-predicate selectivity.
-- Never assume 0.70 × 0.25 as a fact for this dataset.

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
