-- Lab 02: Database Index Foundation
-- Experiment: When Seq Scan IS Correct
-- Destroys the misconception: "Seq Scan = Bad, Index Scan = Good"

-- ============================================
-- PRIMARY LESSON:
-- ============================================
-- The goal is NOT to force index usage.
-- The goal is to MINIMIZE TOTAL QUERY COST for the real workload.
--
-- Seq Scan is not inherently bad.
-- Index Scan is not inherently good.
-- The best choice depends on selectivity, table size, caches, etc.

-- ============================================
-- PART 1: Ensure Index Exists
-- ============================================

DROP INDEX IF EXISTS idx_service_status;
CREATE INDEX idx_service_status ON service(status);

-- Sanity check: total rows in the table
SELECT count(*) AS total_rows FROM service;

SELECT status, count(*) AS cnt,
       round(count(*) * 100.0 / (SELECT count(*) FROM service), 2) AS percentage
FROM service GROUP BY status ORDER BY cnt DESC;

-- ============================================
-- PART 2: High Match Fraction Query (70%)
-- Seq Scan should win
-- ============================================

-- Query: Fetch most FINISHED services (approximately 70% of rows)
-- High match fraction -> not very selective -> Seq Scan wins
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';

-- EXPECTATION: PostgreSQL should choose Seq Scan!
-- REASON: Returning ~70% of the table means reading most heap pages anyway.
--         Traversing an index would add more I/O without benefit.

-- ============================================
-- PART 3: Low Match Fraction Query (0.1%)
-- Index Scan should win
-- ============================================

-- Query: Find rare status (PENDING_REFUND approximately 0.1% of rows)
-- Low match fraction -> highly selective -> Index Scan wins
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'PENDING_REFUND';

-- EXPECTATION: PostgreSQL should choose Index Scan!
-- REASON: Very few rows match. Index lookup is far cheaper than
--         scanning entire table.

-- ============================================
-- PART 4: Forced Index Comparison (Educational Only)
-- ============================================

-- WARNING: Planner switches do NOT force a specific access path.
-- They strongly discourage the planner's usual choice.
--
-- enable_seqscan = off does not force a particular index.
-- It strongly discourages sequential plans.
--
-- Planner switches are diagnostic/educational tools,
-- not a production tuning strategy.

-- Verify Seq Scan is chosen for high-match query normally
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';

-- With enable_seqscan = off, PostgreSQL will try Index Scan,
-- but may still use Bitmap Heap Scan if that's cheaper, or fall back
-- to Seq Scan if no usable index exists.
SET LOCAL enable_seqscan = off;
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT count(*) FROM service WHERE status = 'FINISHED';
RESET enable_seqscan;

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_status;