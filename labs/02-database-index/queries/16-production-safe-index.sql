-- Lab 02: Database Index Foundation
-- Experiment: Production-Safe Index Creation (`CONCURRENTLY`)

-- CONTEXT:
-- Adding a new index to a large production table is an OPERATIONAL concern.
-- It can block writes for the duration of the build, or put heavy load on the disk.

-- ============================================
-- PART 1: The Dangers of Standard `CREATE INDEX`
-- ============================================

-- A normal CREATE INDEX permits reads but blocks WRITES (INSERT/UPDATE/DELETE)
-- on the indexed table until the build finishes.
--
-- On a 100M+ row production table, an index build can take minutes or hours.
-- During this time, your application is frozen for writes to that table.

-- DON'T RUN THIS ON PRODUCTION (it's instant on small local data):
-- CREATE INDEX idx_blocking ON service(branch_id);

-- ============================================
-- PART 2: The Solution - `CREATE INDEX CONCURRENTLY`
-- ============================================

-- `CREATE INDEX CONCURRENTLY` builds the index in multiple phases:
-- 1. First scan builds index, recording concurrent writes
-- 2. Second scan applies changes that happened during Phase 1

-- TRADE-OFFS:
-- 1. Build takes longer and uses more resources.
-- 2. If it fails (conflicts or disk space), index is left INVALID.
--    You MUST drop it and retry (DROP INDEX CONCURRENTLY).
-- 3. Cannot run inside a TRANSACTION block.
-- 4. Monitors old snapshots; can wait for long-running transactions.
-- 5. Still consumes CPU and I/O.

-- The operational syntax for production:
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

-- ============================================
-- PART 3: Operational Checklist for Index Creation
-- ============================================

-- Before adding an index in production, consider:
-- 1. Table size
SELECT
    pg_size_pretty(pg_relation_size('service')) AS table_size,
    pg_size_pretty(pg_indexes_size('service')) AS total_indexes;

-- 2. Available disk space
--    PostgreSQL needs space for:
--    - New index
--    - WAL for the index creation
--    - Temporary files (during sorting, if needed)
--    Use OS: df -h on the database server

-- 3. Write traffic level
--    High write volume = longer build time under lock with standard CREATE INDEX.
--    CONCURRENTLY helps but still uses resources.

-- 4. Maintenance window considerations
--    CONCURRENTLY still impacts performance; not truly "non-blocking".

-- 5. Replication/WAL impact
--    Index creation generates WAL; remote replicas replay all of it.

-- 6. Run with extended statement_timeout:
--    SET statement_timeout = '0';

-- 7. Monitor progress:
SELECT * FROM pg_stat_progress_create_index
WHERE pid = pg_backend_pid();

-- 8. Monitor locks for blocking queries:
SELECT locktype, relation::regclass, mode, granted
FROM pg_locks
WHERE relation = 'service'::regclass;

-- ============================================
-- PART 4: Handling Failed Builds (Invalid Indexes)
-- ============================================

-- If CONCURRENTLY build fails, index is INVALID.
-- NEVER use an invalid index; drop it immediately.

SELECT
    indexrelid::regclass AS invalid_index
FROM pg_index
WHERE NOT indisvalid
  AND indrelid = 'service'::regclass;

-- DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;

-- ============================================
-- PART 5: Verify Query Plan AFTER Index Creation
-- ============================================

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 5;

-- Check: Did planner use the new index?
-- If not, verify:
-- - Query matches index's leftmost prefix
-- - Statistics current (`ANALYZE`)
-- - Table not too small (Seq Scan may be cheaper)

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;