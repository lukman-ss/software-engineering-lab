-- Lab 02: Database Index Foundation
-- Experiment: Low Cardinality vs Selectivity
-- Corrects the oversimplified "don't index low-cardinality columns" rule

-- KEY INSIGHT: Cardinality (distinct values) is different from Selectivity (fraction matched)
-- - Cardinality = number of distinct values in a column
-- - Selectivity = fraction of rows a predicate filters out (0.0 to 1.0)

-- An index on a low-cardinality column CAN be useful when:
-- 1. The predicate selects a small fraction (high selectivity)
-- 2. Index is combined with other indexes (BitmapAnd)
-- 3. The column is used together with high-selectivity columns

-- Check statistics for the status column
SELECT
    attname,
    n_distinct,
    array_to_string(most_common_vals, ', ') AS most_common_vals,
    array_to_string(most_common_freqs, ', ') AS most_common_freqs
FROM pg_stats
WHERE tablename = 'service'
AND attname = 'status';

-- n_distinct = -1 means negative value = count of distinct values
-- n_distinct = -4 means absolute value = count of distinct values
-- Positive value = fraction (selectivity) of distinct values

-- The rule "don't index low-cardinality columns" is WRONG because:
-- Scenario A: status = 'FINISHED' (75% of rows) -> Low selectivity -> Index may not help
-- Scenario B: status = 'PENDING_REFUND' (0.1% of rows) -> High selectivity -> Index very helpful!

================================
-- SCENARIO A: FINISHED = ~75% of rows
-- Low selectivity predicate
-- PostgreSQL may prefer Seq Scan
================================

-- First, verify actual percentages
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

================================
-- SCENARIO B: PENDING_REFUND = ~0.1% of rows
-- High selectivity predicate (good for index!)
-- Same index, different value
================================

-- Verify the rare status distribution
SELECT status, COUNT(*) FROM service WHERE status = 'PENDING_REFUND';

-- Run the query
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'PENDING_REFUND';

-- Typical outcome: Index Scan wins because only ~500 rows need checking

-- This proves: same index, different predicate selectivity = different plans!

================================
-- COMBINED INDEX DEMONSTRATION
-- How low-cardinality can still help
================================

-- Create composite index for our original query pattern
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date);

-- Now the status filter becomes useful:
-- Even though status='FINISHED' matches 75% of table,
-- when combined with branch_id=2 (20% of table),
-- the effective selectivity is 0.75 * 0.20 = 15%

EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Clean up
DROP INDEX IF EXISTS idx_service_branch_status_date;