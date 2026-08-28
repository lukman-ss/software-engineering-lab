-- Lab 02: Database Index Foundation
-- Experiment 17: Repeatable Benchmark Harness
-- Compares actual measured times across multiple index strategies

-- WORKFLOW:
-- 1. Run each EXPLAIN with `(ANALYZE, BUFFERS)` to inspect plan
-- 3. Run raw timing query multiple times (e.g., 5 times)
-- 4. Record MIN / MAX / AVG manually

-- LIMITATIONS:
-- - WARM-CACHE: Once PostgreSQL has read a table into RAM, subsequent queries
--   will be much faster than the first run.
-- - COLD-CACHE: We deliberately do NOT drop OS caches between runs.
--   That is unsafe on shared systems / production servers.
-- - So all numbers reflect a warm-cache scenario by default.
-- - The point is RELATIVE performance between strategies, not absolute timing.

================================
-- SETUP: ensure clean index state for fair comparison
================================

DROP INDEX IF EXISTS idx_bench_branch_id;
DROP INDEX IF EXISTS idx_bench_status;
DROP INDEX IF EXISTS idx_bench_service_date;
DROP INDEX IF EXISTS idx_bench_composite_wrong;
DROP INDEX IF EXISTS idx_bench_composite_correct;
DROP INDEX IF EXISTS idx_bench_covering;

ANALYZE service;

================================
-- SCENARIO A: No Secondary Index (Baseline)
================================

-- Inspect plan
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Time 5 raw runs (warm cache)
\timing on
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
\timing off

-- Record:
-- Min: _____ ms  Max: _____ ms  Avg: _____
-- Plan: _____  Buffers: read _____ hit _____

================================
-- SCENARIO B: Three Single-Column Indexes
================================

CREATE INDEX idx_bench_branch_id ON service(branch_id);
CREATE INDEX idx_bench_status ON service(status);
CREATE INDEX idx_bench_service_date ON service(service_date);

ANALYZE service;

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

\timing on
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
\timing off

-- Record:
-- Min: _____ ms  Max: _____ ms  Avg: _____
-- Plan: _____  Buffers: read _____ hit _____

================================
-- SCENARIO C: Wrong Composite Index Order
================================

DROP INDEX idx_bench_branch_id;
DROP INDEX idx_bench_status;
DROP INDEX idx_bench_service_date;

CREATE INDEX idx_bench_composite_wrong
    ON service(service_date, status, branch_id);

ANALYZE service;

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

\timing on
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
\timing off

-- Record:
-- Min: _____ ms  Max: _____ ms  Avg: _____
-- Plan: _____  Buffers: read _____ hit _____

================================
-- SCENARIO D: Recommended Composite Index
================================

DROP INDEX idx_bench_composite_wrong;

CREATE INDEX idx_bench_composite_correct
    ON service(branch_id, status, service_date DESC);

ANALYZE service;

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

\timing on
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
\timing off

-- Record:
-- Min: _____ ms  Max: _____ ms  Avg: _____
-- Plan: _____  Buffers: read _____ hit _____

================================
-- SCENARIO E: Covering Index (INCLUDE)
================================

DROP INDEX idx_bench_composite_correct;

CREATE INDEX idx_bench_covering
    ON service(branch_id, status, service_date DESC)
    INCLUDE (customer_id, mechanic_id, invoice_no);

ANALYZE service;

EXPLAIN (ANALYZE, BUFFERS)
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

\timing on
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service WHERE branch_id = 2 AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
\timing off

-- Record:
-- Min: _____ ms  Max: _____ ms  Avg: _____
-- Plan: _____  Buffers: read _____ hit _____

-- ============================================
-- COMPARISON TABLE (fill in your numbers)
-- ============================================
--
-- | Scenario                | Plan                  | Avg ms | Buffers |
-- |-------------------------|-----------------------|--------|---------|
-- | A: No index             |                       |        |         |
-- | B: Single-column x3     |                       |        |         |
-- | C: Wrong composite      |                       |        |         |
-- | D: Recommended composite|                       |        |         |
-- | E: Covering (INCLUDE)   |                       |        |         |

================================
-- CLEANUP
================================

DROP INDEX IF EXISTS idx_bench_branch_id;
DROP INDEX IF EXISTS idx_bench_status;
DROP INDEX IF EXISTS idx_bench_service_date;
DROP INDEX IF EXISTS idx_bench_composite_wrong;
DROP INDEX IF EXISTS idx_bench_composite_correct;
DROP INDEX IF EXISTS idx_bench_covering;