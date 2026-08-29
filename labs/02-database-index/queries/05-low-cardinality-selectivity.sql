-- Lab 02: Database Index Foundation
-- Experiment: Low Cardinality vs Selectivity
-- Corrects the oversimplified "don't index low-cardinality columns" rule

-- KEY INSIGHT: Cardinality (distinct values) is different from Match fraction
-- ============================================
-- Cardianity = number of distinct values
-- Match fraction = rows matching predicate / total rows
-- ============================================
--
-- Cardinality of `status`: 5 distinct values (low cardinality)
-- But match fraction varies dramatically:
-- - FINISHED: ~70% match fraction (not very selective)
-- - PENDING_REFUND: ~0.1% match fraction (highly selective)
--
-- Index usefulness depends on match fraction, not cardinality alone.

-- Note on pg_stats.n_distinct:
-- n_distinct > 0 = estimated number of distinct values
-- n_distinct < 0 = negative fraction: |n_distinct| = distinct_rows / total_rows
-- Example: n_distinct = -1 means approximately one distinct value per row
-- Example: n_distinct = -0.5 means approximately half as many distinct values as total rows

-- An index on a low-cardinality column CAN be useful when:
-- 1. The predicate selects a small fraction (highly selective)
-- 2. Index is combined with other indexes (BitmapAnd)
-- 3. The column is used together with highly selective columns

-- ============================================
-- SETUP: Create single-column index on status
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
CREATE INDEX idx_service_status ON service(status);

-- ============================================
-- SCENARIO A: FINISHED = ~70% of rows
-- Not very selective predicate
-- ============================================

SELECT status, round(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM service), 2) AS percentage
FROM service
GROUP BY status
ORDER BY percentage DESC;

-- Run the query with single-column index
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'FINISHED';

-- Typical outcome: Seq Scan wins because index would scan ~70% of rows
-- This is the "index not helpful" case due to not very selective predicate
-- Even though status has low cardinality (5 values), FINISHED = 70% match fraction
-- makes index traversal more expensive than sequential scan

-- ============================================
-- SCENARIO B: PENDING_REFUND = ~0.1% of rows
-- Highly selective predicate (good for index!)
-- ============================================

SELECT status FROM service WHERE status = 'PENDING_REFUND';

EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'PENDING_REFUND';

-- Typical outcome: Index Scan wins because only ~500 rows match
-- Low cardinality of status + highly selective predicate (0.1% match fraction)
-- = index is very attractive

-- ============================================
-- KEY LESSON: Same index, different predicates
-- ============================================
-- The SAME physical index idx_service_status on the same low-cardinality column
-- can produce DIFFERENT planner choices because:
-- - status = 'FINISHED'  → 70% match → Seq Scan cheaper
-- - status = 'PENDING_REFUND' → 0.1% match → Index Scan cheaper

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_status;