-- Lab 02: Database Index Foundation
-- Experiment: Index Usage Audit
-- Teaches how to inspect index usage statistics in PostgreSQL

-- WARNING: Do NOT simply drop indexes with idx_scan = 0!
-- Zero usage can be misleading due to statistics reset, restarts, or short observation periods

-- ============================================
-- Key Statistics to Monitor (pg_stat_user_indexes)
-- ============================================

-- idx_scan: Number of index scans initiated. This is incremented
-- each time the planner chooses to use this index for a query.

-- idx_tup_read: Number of index entries returned by scans.
-- Includes rows found via Bitmap Index Scans (bitmap is built from multiple indexes).

-- idx_tup_fetch: Number of live heap tuples fetched by "Index Scans" (not Index Only Scans).
-- Note: Bitmap Heap Scans do NOT increment per-index idx_tup_fetch - they fetch
-- rows in a separate pass after combining bitmaps.

-- IMPORTANT: idx_tup_fetch / idx_tup_read is NOT a generic "efficiency percentage".
-- It measures the ratio of heap tuples fetched via Index Scans only.
-- Bitmap Scans can have very different behavior.

-- ============================================
-- Index Structure Reference
-- ============================================

SELECT
    indexrelid::regclass AS index_name,
    indisunique AS is_unique,
    indisprimary AS is_primary,
    indkey::int[] AS column_numbers,
    indpred AS partial_predicate
FROM pg_index
WHERE indrelid = 'service'::regclass;

-- ============================================
-- Current Index Usage Statistics
-- ============================================

SELECT
    indexname,
    idx_scan AS "Index Scans",
    idx_tup_read AS "Tuples Read",
    idx_tup_fetch AS "Tuples Fetched",
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND tablename = 'service'
ORDER BY idx_scan DESC;

-- ============================================
-- Audit Checklist for Candidate Indexes
-- ============================================

-- 1. Identify low-usage indexes
SELECT
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND tablename = 'service'
ORDER BY idx_scan ASC NULLS LAST;

-- WARNING: Before dropping an index, verify:
--
-- ✓ Statistics period is long enough for your workload
-- ✓ No recent server restart would have reset counters to 0
-- ✓ Query patterns are stable (not seasonal changes)
-- ✓ Index might support rare operational queries
-- ✓ Index backs a constraint that requires it:
--   - PRIMARY KEY (automatically created, cannot drop)
--   - UNIQUE (required for constraint, cannot drop)
--   - EXCLUSION constraint (if applicable)
-- ✓ Foreign key indexes: a FK column typically has an index on the
--   referencing table for performance, but it is NOT strictly required.
-- ✓ Multiple indexes sharing same column might be useful for different queries

-- ============================================
-- Indexes Supporting Constraints
-- ============================================

-- PRIMARY KEY and UNIQUE indexes are constraint-backed
SELECT
    c.conname AS constraint_name,
    c.contype,
    i.indexrelid::regclass AS index_name
FROM pg_constraint c
JOIN pg_index i ON c.conindid = i.indexrelid
WHERE c.conrelid = 'service'::regclass;

-- ============================================
-- Tables Without Any Index Usage
-- (Investigate before dropping!)
-- ============================================

SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND tablename = 'service'
  AND idx_scan = 0
ORDER BY pg_relation_size(indexname::regclass) DESC;