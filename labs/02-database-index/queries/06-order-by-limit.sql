-- Lab 02: Database Index Foundation
-- Experiment: ORDER BY + LIMIT Optimization
-- Critical production pattern for dashboards and pagination

-- PostgreSQL 16 compatibility note:
-- This repository targets PostgreSQL 16. PostgreSQL 18 introduced B-tree skip-scan
-- optimization, which is NOT covered here.

-- ============================================
-- The pattern: WHERE + ORDER BY + LIMIT
-- Key insight: Proper index can avoid explicit Sort AND stop early
-- ============================================

-- Create index that can satisfy ORDER BY
DROP INDEX IF EXISTS idx_service_branch_status_date_desc;
CREATE INDEX idx_service_branch_status_date_desc
    ON service(branch_id, status, service_date DESC);
ANALYZE service;

-- ============================================
-- QUERY: Dashboard-style - latest 20 finished services for branch 2
-- ============================================

-- With the appropriate index in place:
-- PLAN ANALYSIS GUIDANCE:
-- Look for: "Index Scan" or "Index Only Scan"
-- Check: Was Sort present or absent?
-- Look for: "Index Cond" columns = predicates the index resolved directly
-- Look for: "Filter" columns = predicates still requiring heap checks
-- Record: Rows Examined, shared hit/read, actual rows returned
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- ============================================
-- COMPARE: Without appropriate index
-- ============================================

-- Drop the index to simulate missing index
DROP INDEX IF EXISTS idx_service_branch_status_date_desc;
ANALYZE service;

-- Run same query without the index
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- Record:
-- - Plan type (Seq Scan vs Index Scan)
-- - Sort present? Yes/No
-- - Rows examined: _____
-- - Shared buffers: read _____, hit _____
-- - Execution time: _____ ms

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_branch_status_date_desc;