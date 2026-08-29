-- Lab 02: Database Index Foundation
-- Experiment: Production-Safe Index Creation (`CONCURRENTLY`)

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

-- CREATE INDEX CONCURRENTLY uses multiple transaction phases and performs
-- two table scans. The index is initially entered into PostgreSQL catalogs
-- as invalid.
--
-- Before its scan phases, PostgreSQL can wait for transactions that could
-- interfere with the build. After the second scan, it can also wait for
-- old snapshots before marking the index valid.
--
-- This additional work is why CONCURRENTLY normally takes longer than a
-- regular CREATE INDEX, while allowing ordinary INSERT/UPDATE/DELETE traffic
-- to continue.
--
-- TRADE-OFFS:
-- 1. Build takes longer and uses more resources than a plain CREATE INDEX.
-- 2. If the build fails, the index may be left in an INVALID state.
--    PostgreSQL ignores an invalid index for normal query planning because
--    it may be incomplete, but the index can still create update/maintenance
--    overhead. Investigate and remediate the failure.
--    Depending on the situation, remediation can include dropping and
--    recreating the index or using an appropriate REINDEX strategy.
-- 3. Cannot run inside a TRANSACTION block.
-- 4. May wait for long-running transactions or old snapshots.
-- 5. Consumes CPU, I/O, and WAL resources.

-- The operational syntax for production:
DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;
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

-- 7. Monitor progress (see PART 4)
--    Run the query from a SECOND session while the CREATE INDEX is running

-- 8. Monitor locks for blocking queries:
SELECT locktype, relation::regclass, mode, granted
FROM pg_locks
WHERE relation = 'service'::regclass;

-- ============================================
-- PART 4: Monitoring Progress (Two-Session Pattern)
-- ============================================

-- IMPORTANT: pg_stat_progress_create_index must be queried from a
-- DIFFERENT session than the one running CREATE INDEX CONCURRENTLY.
-- Running it in the same session shows the already-completed build.

-- Session A: Run in one terminal
--   psql -d se_lab -c "CREATE INDEX CONCURRENTLY idx_service_concurrent_branch_test ON service(branch_id);"

-- Session B: Run in another terminal while Session A is active
-- This demonstrates the real-time progress:
SELECT
    pid,
    datname,
    relid::regclass AS table_name,
    index_relid::regclass AS index_name,
    command,
    phase,
    lockers_total,
    lockers_done,
    blocks_total,
    blocks_done
FROM pg_stat_progress_create_index
WHERE relid = 'service'::regclass::oid;

-- ============================================
-- PART 5: Handling Failed Builds (Invalid Indexes)
-- ============================================

-- If CONCURRENTLY build fails, check for invalid indexes:
SELECT
    indexrelid::regclass AS invalid_index
FROM pg_index
WHERE NOT indisvalid
  AND indrelid = 'service'::regclass;

-- Remediation options:
-- - DROP INDEX CONCURRENTLY if the index is incomplete
-- - REINDEX to rebuild in-place (handles some failure cases internally)
-- - Recreate with fresh CREATE INDEX CONCURRENTLY

-- DROP INDEX CONCURRENTLY IF EXISTS idx_service_concurrent_branch;

-- ============================================
-- PART 6: Verify Query Plan AFTER Index Creation
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