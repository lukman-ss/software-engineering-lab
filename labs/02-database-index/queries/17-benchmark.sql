-- Lab 02: Database Index Foundation
-- Experiment 17: Index Strategy Benchmark Harness
--
-- IMPORTANT: This document explains methodology as much as produces executable SQL.
-- Client-visible timing (printing thousands of rows to terminal) is NOT a reliable
-- measure of server query performance. Use EXPLAIN (ANALYZE, BUFFERS) for comparison.

-- ============================================
-- WORKFLOW:
-- 1. Use EXPLAIN (ANALYZE, BUFFERS, TIMING OFF) to inspect execution plans
--    - Shows row/buffer counts without terminal output overhead
-- 2. Record observations in the comparison table at the end
-- 3. Focus on relative performance between strategies, not absolute timing
--
-- LIMITATIONS NOTE:
-- - WARM-CACHE: Once PostgreSQL has read pages into RAM, queries are faster.
-- - We intentionally keep the cache warm for reproducible relative comparisons.
-- - COLD-CACHE testing requires `pg_prewarm` or full server restart.

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
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
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
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
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

-- ============================================
-- SCENARIO C: Wrong Composite Index Order
-- ============================================

CREATE INDEX idx_bench_composite_wrong
    ON service(service_date, status, branch_id);

ANALYZE service;

-- Date is leftmost, cannot efficiently use equality predicates on branch_id/status
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
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
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: Index Scan, Sort eliminated, Buffers = read _____ hit _____

DROP INDEX idx_bench_composite_correct;

-- ============================================
-- SCENARIO E: Covering Index (INCLUDE)
-- ============================================
-- NOTE: This is a DIFFERENT workload comparison.
-- A covering index enables Index-Only Scan by storing all needed columns.
-- This measures index-only vs heap access, NOT apples-to-apples with SELECT * scenarios.

CREATE INDEX idx_bench_covering
    ON service(branch_id, status, service_date DESC)
    INCLUDE (customer_id, mechanic_id, invoice_no);

ANALYZE service;

-- Projection query: only indexed columns needed
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF)
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Record: Index Only Scan (no heap access), Buffers = read _____ hit _____

DROP INDEX idx_bench_covering;

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
-- | E: Covering (INCLUDE)   | Index Only Scan     | _____/_____ |
--
-- KEY OBSERVATIONS:
-- 1. Seq Scan + Sort is slow for selective queries
-- 2. Bitmap Heap Scan combines multiple indexes efficiently
-- 3. Composite index eliminates both index access and Sort
-- 4. Covering index avoids heap access entirely
--
-- REMEMBER: Raw query timing via \timing is affected by:
-- - Network latency (psql to database)
-- - Terminal rendering
-- - pg_backend_pid() scheduling
-- Use EXPLAIN ANALYZE for server-side comparison.