-- Lab 02: Database Index Foundation
-- Experiment 17: Index Strategy Benchmark Harness
--
-- IMPORTANT: This document explains methodology as much as produces executable SQL.
-- Client-visible timing (printing thousands of rows to terminal) is NOT a reliable
-- measure of server query performance. Use EXPLAIN (ANALYZE, BUFFERS) for comparison.

-- ============================================
-- WORKFLOW:
-- 1. Use EXPLAIN (ANALYZE, BUFFERS) to inspect execution plans
--    - Shows row/buffer counts without terminal output overhead
-- 2. Record observations in the comparison table at the end
-- 3. Focus on relative performance between strategies, not absolute timing
-- ============================================

-- ============================================
-- LIMITATIONS NOTE:
-- ============================================
-- WARM-CACHE: Once PostgreSQL has read pages into RAM, queries are faster.
-- We intentionally keep the cache warm for reproducible relative comparisons.
--
-- Other factors affecting measurements:
-- - OS page cache: shared read may be satisfied from kernel cache
-- - Background load: concurrent activity affects timing
-- - Checkpoints: can cause I/O spikes
-- - WAL activity: index maintenance generates WAL
-- - Hardware: CPU, RAM, disk speed vary between systems
-- - PostgreSQL configuration: shared_buffers, work_mem, etc.
--
-- Cold-cache testing requires pg_prewarm or full server restart.
-- Do NOT auto-flush OS caches - results are environment-dependent.
--
-- EXPLAIN (ANALYZE) shows END-TO-END/client timing including network overhead.
-- For pure server timing, consider using EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
-- and measuring system clock time around the query.

-- ============================================
-- SETUP: Ensure clean index state for fair comparison
-- ============================================

DROP INDEX IF EXISTS idx_bench_branch_id;
DROP INDEX IF EXISTS idx_bench_status;
DROP INDEX IF EXISTS idx_bench_service_date;
DROP INDEX IF EXISTS idx_bench_composite_wrong;
DROP INDEX IF EXISTS idx_bench_composite_correct;
DROP INDEX IF EXISTS idx_bench_covering;

ANALYZE service;

-- ============================================
-- SCENARIO A: No Secondary Index (Baseline)
-- ============================================

-- Compare execution plans
-- Full table scan with Sort node
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: Plan = Seq Scan + Sort, Buffers = read _____ hit _____

-- ============================================
-- SCENARIO B: Three Single-Column Indexes
-- ============================================

CREATE INDEX idx_bench_branch_id ON service(branch_id);
CREATE INDEX idx_bench_status ON service(status);
CREATE INDEX idx_bench_service_date ON service(service_date);

ANALYZE service;

-- Planner may use Bitmap Heap Scan + BitmapAnd to combine indexes
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: Plan may be Bitmap Heap Scan, Buffers = read _____ hit _____

DROP INDEX idx_bench_branch_id;
DROP INDEX idx_bench_status;
DROP INDEX idx_bench_service_date;
DROP INDEX IF EXISTS idx_service_status;

-- ============================================
-- SCENARIO C: Wrong Composite Index Order
-- ============================================

CREATE INDEX idx_bench_composite_wrong
    ON service(service_date, status, branch_id);

ANALYZE service;

-- Date is leftmost; equality predicates on branch_id/status cannot be used efficiently
-- PostgreSQL can still use the index for the date range, but must scan
-- a larger portion since leading column doesn't constrain equality
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: May still use index for date range, but inefficient, Buffers = read _____ hit _____

DROP INDEX idx_bench_composite_wrong;

-- ============================================
-- SCENARIO D: Recommended Composite Index
-- ============================================

CREATE INDEX idx_bench_composite_correct
    ON service(branch_id, status, service_date DESC);

ANALYZE service;

-- Equality columns first (branch_id, status), then range column (service_date)
-- This allows tight B-tree range bounding from left to right
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: Index Scan, Sort eliminated, Buffers = read _____ hit _____

DROP INDEX idx_bench_composite_correct;

-- ============================================
-- COMPARISON TABLE (fill in your observations)
-- ============================================
--
-- | Scenario                | Plan                | Buffers (R/H) |
-- |-------------------------|---------------------|-------------|
-- | A: No index             | Seq Scan + Sort     | _____/_____ |
-- | B: Single-column x3     | Bitmap Heap Scan    | _____/_____ |
-- | C: Wrong composite      | Index + Sort        | _____/_____ |
-- | D: Recommended composite| Index Scan (no Sort)| _____/_____ |
--
-- KEY OBSERVATIONS:
-- 1. Seq Scan + Sort is slow for selective queries
-- 2. Bitmap Heap Scan combines multiple indexes efficiently
-- 3. Composite index eliminates both index access and Sort
-- 4. Covering (INCLUDE) indexes belong in separate experiment

-- ============================================
-- NOTES ON PLAN NODES:
-- ============================================
-- - Index Scan: Single index used for finding rows
-- - Bitmap Heap Scan: Multiple indexes combined into bitmap, then heap fetched
-- - Seq Scan: Sequential table read (often correct for high-match queries)
-- - Index Only Scan: All needed columns from index, no heap access
--   (requires visibility map conditions; check "Heap Fetches" in output)
-- - Sort: Explicit sort operation; eliminated when ORDER BY matches index

-- ============================================
-- NOTE: Covering indices comparison
-- ============================================
-- Covering indexes are investigated in 07-covering-index.sql.
-- They enable Index Only Scan but require different query patterns.
-- Comparing SELECT * benchmark against covering-index queries on
-- different columns would not be apples-to-apples.