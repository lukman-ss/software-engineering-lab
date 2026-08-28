-- Lab 02: Database Index Foundation
-- EXPLAIN statements for query analysis

-- 1. Basic EXPLAIN (shows planned execution, no actual execution)
EXPLAIN
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- 2. EXPLAIN ANALYZE (executes the query and shows actual runtime)
EXPLAIN ANALYZE
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- 3. EXPLAIN ANALYZE with buffer information
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- 4. Check the index we created
EXPLAIN
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- 5. Test with different query patterns to understand index usage
-- Query by status only (tests index selectivity)
EXPLAIN
SELECT *
FROM service
WHERE status = 'FINISHED'
ORDER BY service_date DESC;

-- Query by branch only
EXPLAIN
SELECT *
FROM service
WHERE branch_id = 2
ORDER BY service_date DESC;