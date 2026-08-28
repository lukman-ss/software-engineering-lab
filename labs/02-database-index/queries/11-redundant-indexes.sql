-- Lab 02: Database Index Foundation
-- Experiment: Redundant and Overlapping Indexes
-- Teaches workload-driven index cleanup and prefix redundancy

-- KEY CONCEPTS:
-- 1. Duplicate index: Identical definition on the same table (always redundant, safe to drop).
-- 2. Overlapping index: Indexes sharing leading columns (e.g., (a) and (a, b)).
-- 3. Potentially redundant index: A prefix index that may be completely covered by a longer composite index.
-- 4. Necessary workload-specific index: A shorter prefix index kept specifically because it is much smaller,
--    quicker to scan, and cheaper to maintain for high-frequency simple queries.

================================
-- PART 1: Exact Duplicate Detection
================================

-- Query to find exact duplicate indexes in PostgreSQL
SELECT
    indrelid::regclass AS table_name,
    indexrelid::regclass AS index_name,
    indkey::text AS column_mapping,
    indclass::text AS op_classes
FROM pg_index
WHERE indrelid = 'service'::regclass
GROUP BY table_name, index_name, column_mapping, op_classes
HAVING count(*) > 1;

================================
-- PART 2: Overlapping / Prefix Redundancy Scenario
================================

-- Imagine a table with these three indexes:
DROP INDEX IF EXISTS idx_service_red_1;
DROP INDEX IF EXISTS idx_service_red_2;
DROP INDEX IF EXISTS idx_service_red_3;

-- Index 1: Single column
CREATE INDEX idx_service_red_1 ON service(branch_id);

-- Index 2: Two columns (starts with branch_id)
CREATE INDEX idx_service_red_2 ON service(branch_id, status);

-- Index 3: Three columns (starts with branch_id, status)
CREATE INDEX idx_service_red_3 ON service(branch_id, status, service_date DESC);

-- Analyze storage size of each overlapping index
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_indexes
WHERE tablename = 'service'
  AND indexname LIKE 'idx_service_red_%';

================================
-- PART 3: Which queries use which index?
================================

-- Query A: WHERE branch_id = 2
-- Can be satisfied by idx_service_red_1, idx_service_red_2, or idx_service_red_3!
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 LIMIT 10;

-- Query B: WHERE branch_id = 2 AND status = 'FINISHED'
-- Can be satisfied by idx_service_red_2 or idx_service_red_3! (idx_service_red_1 is not enough)
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED' LIMIT 10;

-- Query C: WHERE branch_id = 2 AND status = 'FINISHED' AND service_date BETWEEN ...
-- Only idx_service_red_3 satisfies all three efficiently!
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31';

================================
-- PART 4: The Trade-off (When to keep vs drop prefix indexes)
================================

-- Can we drop idx_service_red_1 and idx_service_red_2?
-- Technically, idx_service_red_3 covers all queries that idx_service_red_1 and idx_service_red_2 handle.
-- BUT consider:
-- 1. Size: idx_service_red_1 is much smaller than idx_service_red_3. A scan on idx_service_red_1
--    reads fewer pages than scanning the left prefix of idx_service_red_3 if high-frequency simple queries run.
-- 2. Write overhead: Maintaining 3 indexes costs more INSERT/UPDATE performance than maintaining 1.
-- 3. Conclusion: Do not blindly drop prefix indexes without checking usage statistics (`pg_stat_user_indexes`)
--    and query frequency.

-- Cleanup
DROP INDEX IF EXISTS idx_service_red_1;
DROP INDEX IF EXISTS idx_service_red_2;
DROP INDEX IF EXISTS idx_service_red_3;