-- Lab 02: Database Index Foundation
-- Experiment: Low Cardinality vs Selectivity
-- Corrects the oversimplified "don't index low-cardinality columns" rule

-- KEY INSIGHT: Cardinality (distinct values) is different from Match fraction
-- - Cardinality = number of distinct values in a column
-- - Match fraction (selectivity) = fraction of rows a predicate filters out (0.0 to 1.0)

-- Note on pg_stats.n_distinct:
-- n_distinct > 0 = estimated number of distinct values
-- n_distinct < 0 = negative fraction: |n_distinct| = distinct_rows / total_rows
-- Example: n_distinct = -1 means nearly one distinct value per row (high uniqueness)
-- Example: n_distinct = -0.5 means about half as many distinct values as rows

-- An index on a low-cardinality column CAN be useful when:
-- 1. The predicate selects a small fraction (highly selective)
-- 2. Index is combined with other indexes (BitmapAnd)
-- 3. The column is used together with highly selective columns

-- Check statistics for the status column
SELECT
    attname,
    n_distinct,
    array_to_string(most_common_vals, ', ') AS most_common_vals,
    array_to_string(most_common_freqs, ', ') AS most_common_freqs
FROM pg_stats
WHERE tablename = 'service'
AND attname = 'status';

-- The rule "don't index low-cardinality columns" is WRONG because:
-- Scenario A: status = 'FINISHED' (75% of rows) -> not very selective -> Seq Scan may win
-- Scenario B: status = 'PENDING_REFUND' (0.1% of rows) -> highly selective -> Index very helpful!

-- ============================================
-- SCENARIO A: FINISHED = ~75% of rows
-- Not very selective predicate
-- ============================================

SELECT status, COUNT(*) * 100.0 / (SELECT COUNT(*) FROM service) AS pct
FROM service
GROUP BY status
ORDER BY pct DESC;

-- Run the query
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'FINISHED';

-- Typical outcome: Seq Scan wins because index would scan 75% of rows
-- This is the "index not helpful" case

-- Why might index still be used?
-- If combined with other filters in a composite index!

-- ============================================
-- SCENARIO B: PENDING_REFUND = ~0.1% of rows
-- Highly selective predicate (good for index!)
-- ============================================

SELECT status, COUNT(*) FROM service WHERE status = 'PENDING_REFUND';

EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'PENDING_REFUND';

-- Typical outcome: Index Scan wins because only ~500 rows need checking

-- ============================================
-- COMBINED INDEX DEMONSTRATION
-- How low-cardinality can still help
-- ============================================

-- Create composite index for our original query pattern
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date);

-- Now the status filter becomes useful:
-- Important: 0.75 * 0.20 is an APPROXIMATION based on statistical independence.
-- In this synthetic data, FINISHED and branch_id are approximately independent,
-- but in production they may be correlated.
-- The actual selectivity could differ from this naive multiplication.

EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Clean up
DROP INDEX IF EXISTS idx_service_branch_status_date;