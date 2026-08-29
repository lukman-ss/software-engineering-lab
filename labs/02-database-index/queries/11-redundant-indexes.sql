-- Lab 02: Database Index Foundation
-- Experiment: Redundant and Overlapping Indexes
-- Teaches workload-driven index cleanup and prefix redundancy analysis

-- KEY CONCEPTS:
-- 1. Duplicate index: Identical normalized definition (key columns, INCLUDE, predicate). Always redundant.
-- 2. Overlapping index: Indexes sharing leading columns (e.g., (a) and (a, b)).
-- 3. Prefix index: A shorter index that is a leftmost prefix of a longer composite index.
-- 4. Workload-specific index: Kept because usage statistics justify it despite overlap.

-- ============================================
-- PART 1: Exact Duplicate Detection (by normalized definition)
-- ============================================

-- Find pairs of indexes with identical key columns, INCLUDE, and predicate.
-- Group by normalized definition to detect duplicates (same definition, different names).
-- CRITICAL: Before dropping, verify neither index supports a constraint:
-- - PRIMARY KEY index
-- - UNIQUE constraint index
-- - EXCLUSION constraint index

SELECT
    i1.indexrelid::regclass AS index_a,
    i2.indexrelid::regclass AS index_b,
    i1.indkey::text AS key_columns,
    i2.indpred AS predicate
FROM pg_index i1
JOIN pg_index i2 ON i1.indexrelid::regclass::text < i2.indexrelid::regclass::text
WHERE i1.indrelid = 'service'::regclass
  AND i1.indkey = i2.indkey
  AND i1.indpred = i2.indpred
  AND (
      -- Same INCLUDE columns (for partial matches, compare indoption for now)
      ARRAY(
          SELECT attnum
          FROM pg_attribute
          WHERE attrelid = i1.indrelid
            AND attnum = ANY(i1.indkey)
          ORDER BY attnum
      ) = ARRAY(
          SELECT attnum
          FROM pg_attribute
          WHERE attrelid = i2.indrelid
            AND attnum = ANY(i2.indkey)
          ORDER BY attnum
      )
  )
ORDER BY index_a, index_b;

-- ============================================
-- PART 2: Overlapping / Prefix Redundancy Scenario
-- ============================================

DROP INDEX IF EXISTS idx_service_red_1;
DROP INDEX IF EXISTS idx_service_red_2;
DROP INDEX IF EXISTS idx_service_red_3;

-- Index 1: Single column (prefix of Index 3)
CREATE INDEX idx_service_red_1 ON service(branch_id);

-- Index 2: Two columns (prefix of Index 3)
CREATE INDEX idx_service_red_2 ON service(branch_id, status);

-- Index 3: Three columns (covers queries that use Index 1 or 2)
CREATE INDEX idx_service_red_3 ON service(branch_id, status, service_date DESC);

-- Analyze storage size of each overlapping index
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_indexes
WHERE tablename = 'service'
  AND indexname LIKE 'idx_service_red_%';

-- ============================================
-- PART 3: Which queries use which index?
-- ============================================

-- Query A: WHERE branch_id = 2
-- Can use idx_service_red_1, idx_service_red_2, or idx_service_red_3
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 LIMIT 10;

-- Query B: WHERE branch_id = 2 AND status = 'FINISHED'
-- Can use idx_service_red_2 or idx_service_red_3
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED' LIMIT 10;

-- Query C: Full query with ORDER BY
-- Only idx_service_red_3 provides index ordering
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- ============================================
-- PART 4: The Trade-off (When to keep vs drop prefix indexes)
-- ============================================

-- CRITICAL DECISION BEFORE DROPPING:
-- Check pg_stat_user_indexes for actual usage patterns.
-- A smaller prefix index might be kept if:
--   - It serves high-frequency simple queries
--   - It reduces average query latency
--   - Write reduction outweighs read benefits of larger index

-- Example decision framework:
SELECT
    indexname,
    idx_scan AS times_used,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND tablename = 'service'
  AND indexname LIKE 'idx_service_red_%'
ORDER BY idx_scan DESC;

-- Decision: If idx_service_red_1 has idx_scan = 0 and rarely used,
-- it can be dropped. If it serves frequent dashboard queries at
-- microsecond latency, keep it despite overlap with idx_service_red_3.

-- Cleanup
DROP INDEX IF EXISTS idx_service_red_1;
DROP INDEX IF EXISTS idx_service_red_2;
DROP INDEX IF EXISTS idx_service_red_3;