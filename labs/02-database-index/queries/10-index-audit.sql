-- Lab 02: Database Index Foundation
-- Experiment: Index Usage Audit
-- Teaches how to inspect index usage statistics in PostgreSQL

-- Do NOT simply drop indexes with idx_scan = 0!
-- Zero usage can be misleading due to statistics reset, restarts, or short observation periods

================================
-- Key Statistics to Monitor
================================

-- From pg_stat_user_indexes:
-- idx_scan      : How many times has index been used?
-- idx_tup_read  : How many table rows has index have looked at?
-- idx_tup_fetch : How many table rows has index successfully fetched?
--
-- If idx_tup_read is much larger than idx_tup_fetch, index was used for filtering
-- but many rows were scanned that didn't match.

-- Table structure reference
SELECT
    indexrelid::regclass AS index_name,
    indisunique AS is_unique,
    indisprimary AS is_primary,
    indkey::int[] AS column_numbers
FROM pg_index
WHERE indrelid = 'service'::regclass;

================================
-- Current Index Usage Statistics
================================

SELECT
    indexname,
    idx_scan AS "Index Scans",
    idx_tup_read AS "Tuple Reads",
    idx_tup_fetch AS "Tuple Fetches",
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public' AND tablename = 'service'
ORDER BY idx_scan DESC;

================================
-- Audit Checklist for Candidate Indexes
================================

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
-- ✓ Index might enforce constraints (e.g., UNIQUE, FOREIGN KEY)
-- ✓ Multiple indexes sharing same column might be useful

================================
-- Detailed Per-Index Analysis
================================

SELECT
    i.indexrelid::regclass AS index_name,
    i.indkey::text AS indexed_columns,
    a.idx_scan,
    a.idx_tup_read,
    a.idx_tup_fetch,
    -- Efficiency ratio: how selective is this index?
    CASE
        WHEN a.idx_tup_read > 0 THEN
            round(a.idx_tup_fetch::numeric / a.idx_tup_read::numeric * 100, 2)
        ELSE 0
    END AS "Fetch Efficiency %"
FROM pg_index i
LEFT JOIN pg_stat_user_indexes a ON i.indexrelid = a.indexrelid
WHERE i.indrelid = 'service'::regclass
ORDER BY a.idx_scan DESC NULLS LAST;

================================
-- Tables Without Any Index Usage
-- (Investigate before dropping!)

SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND tablename = 'service'
  AND idx_scan = 0
ORDER BY pg_relation_size(indexname::regclass) DESC;