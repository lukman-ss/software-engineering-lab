-- Lab 02: Database Index Foundation
-- Experiment: Production-Safe Index Creation (`CONCURRENTLY`)

-- CONTEXT:
-- Adding a new index to a large production table is an OPERATIONAL concern.
-- It can block writes for the duration of the build, or put heavy load on the disk.

================================
-- PART 1: The Dangers of `CREATE INDEX`
================================

-- Standard `CREATE INDEX` takes an ACCESS EXCLUSIVE LOCK briefly, then a SHARE lock
-- while building. The SHARE lock BLOCKS all writes (INSERT/UPDATE/DELETE) to the table
-- until the index build completes!
--
-- On a 100M+ row production table, an index build can take minutes or hours.
-- During this time, your application is essentially frozen for writes to that table.

-- DON'T RUN THIS ON PRODUCTION (it's instant on small local data):
-- CREATE INDEX idx_blocking ON service(branch_id);

================================
-- PART 2: The Solution - `CREATE INDEX CONCURRENTLY`
================================

-- `CREATE INDEX CONCURRENTLY` builds the index in two phases, doing multiple scans
-- of the table and allowing CONCURRENT reads and writes during the build.

-- 1. Phase 1: Scan table and build index, recording any concurrent writes.
-- 2. Phase 2: Re-scan table to apply the changes that happened during Phase 1.

-- TRADE-OFFS:
-- 1. Build takes significantly longer (and uses more resources).
-- 2. If it fails (due to conflicts or disk space), the index is left INVALID.
--    You MUST DROP it and retry.
-- 3. Cannot run inside a TRANSACTION block.

-- The "safer" syntax for production:
DROP INDEX IF EXISTS idx_service_concurrent_branch;
CREATE INDEX CONCURRENTLY idx_service_concurrent_branch
    ON service(branch_id);

-- Verify validity
SELECT
    indexrelid::regclass AS index_name,
    indisvalid AS is_valid,
    indisready AS is_ready
FROM pg_index
WHERE indexrelid = 'idx_service_concurrent_branch'::regclass;

================================
-- PART 3: Operational Checklist for Index Creation
================================

-- Before adding an index in production:

-- 1. Measure table size
SELECT
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS total_indexes;

-- 2. Estimate index size (estimate only; actual depends on data)
--    A B-tree index on a `bigint` or `int` typically takes ~30-50% of the table size.

-- 3. Check available disk space
--    Use OS commands like `df -h` on the database server.

-- 4. Choose an appropriate load window (low-traffic time)
--    CONCURRENTLY still uses I/O and CPU resources.

-- 5. Run with extended statement_timeout to prevent timeouts
--    Example: SET statement_timeout = '0';
--    Or set it in the connection string / pool config.

-- 6. Monitor progress with `pg_stat_progress_create_index`
SELECT * FROM pg_stat_progress_create_index;

-- 7. Monitor locks to ensure no long-running conflicting queries
SELECT
    locktype,
    relation::regclass,
    mode,
    granted
FROM pg_locks
WHERE relation = 'service'::regclass;

================================
-- PART 4: Handling Failed Builds (Invalid Indexes)
================================

-- If a CONCURRENTLY build fails, the index is INVALID.
-- NEVER use an invalid index; it just wastes space.

-- Check for invalid indexes
SELECT
    indexrelid::regclass AS invalid_index
FROM pg_index
WHERE NOT indisvalid
  AND indrelid = 'service'::regclass;

-- Drop invalid indexes before retrying
-- DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;

-- Retry the build
-- CREATE INDEX CONCURRENTLY idx_service_concurrent_branch ON service(branch_id);

-- Verify after build
-- SELECT indisvalid FROM pg_index WHERE indexrelid = 'idx_service_concurrent_branch'::regclass;

================================
-- PART 5: Verify Query Plan AFTER Index Creation
================================

-- Run the query the index was meant for:
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 5;

-- Did the planner actually pick the new index?
-- If not, check:
-- - Does the query match the index's leftmost prefix?
-- - Is the table too small (planner prefers Seq Scan for small tables)?
-- - Are statistics up to date? (`ANALYZE`)

================================
-- CLEANUP
================================

DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;