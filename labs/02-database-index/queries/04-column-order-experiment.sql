-- Lab 02: Database Index Foundation
-- Experiment: Composite Index Column Order
-- Demonstrates why column order matters in multicolumn B-tree indexes

-- PostgreSQL 16 compatibility note:
-- This repository targets PostgreSQL 16. PostgreSQL 18 introduced B-tree skip-scan
-- optimization, which is NOT covered here. All behavior described below is for PG16.

-- ============================================
-- SETUP: Create three indexes with different column orders
-- ============================================

-- Index A: Optimal for our main query (equality + equality + range)
DROP INDEX IF EXISTS idx_service_a_branch_status_date;
CREATE INDEX idx_service_a_branch_status_date
    ON service(branch_id, status, service_date);

-- Index B: Range first (less useful for our query pattern)
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
CREATE INDEX idx_service_b_date_branch_status
    ON service(service_date, branch_id, status);

-- Index C: Different leading column
DROP INDEX IF EXISTS idx_service_c_status_date_branch;
CREATE INDEX idx_service_c_status_date_branch
    ON service(status, service_date, branch_id);

-- PostgreSQL 16 B-tree multicolumn behavior:
-- For index on (a, b, c), constraints on leading columns determine how much
-- of the index scan range can be bounded efficiently.
--
-- - WHERE a = ?           → constrains index range efficiently
-- - WHERE a = ? AND b = ? → constrains further
-- - WHERE a = ? AND b = ? AND c > ? → bounds tight range
-- - WHERE b = ?           → can use index, but may scan large/complete portion;
--                            planner often prefers Seq Scan or another index
-- - WHERE c = ?           → same reasoning as above

-- ============================================
-- QUERY 1: WHERE branch_id = 2
-- ============================================

-- With Index A (branch_id, status, service_date)
-- Constrains index range to branch_id=2; very efficient
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;

-- With Index B (service_date, branch_id, status)
-- branch_id is not the leading column. PostgreSQL 16 can in principle use
-- this index, but without a constraint on service_date it may need to scan
-- a large or complete portion. Planner often prefers Seq Scan here.
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;

-- With Index C (status, service_date, branch_id)
-- Same reasoning: branch_id is not leading, so planner may prefer Seq Scan.
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;

-- Expected: Index A is the most efficient for Query 1.
-- Index B and C may still be used (not literally impossible), but the planner
-- will likely choose Seq Scan because scanning a large portion of the index
-- is more expensive than a sequential table scan.

-- ============================================
-- QUERY 2: WHERE branch_id = 2 AND status = 'FINISHED'
-- ============================================

-- Index A: constrains on branch_id then status; very efficient
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED';

-- Index B: service_date first; without date constraint, planner likely Seq Scan
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED';

-- Index C: constrains on status then service_date then branch_id;
-- can use, but branch_id filtering is less efficient than Index A
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE status = 'FINISHED' AND branch_id = 2;

-- Key distinction:
-- - Reducing scan range: Index narrows WHERE clause from the start
-- - Filtering: Index finds rows but needs extra checks deeper in

-- ============================================
-- QUERY 3: Full query with ORDER BY
-- ============================================

EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- ============================================
-- QUERY 4-6: Individual predicates
-- ============================================

-- Query 4: WHERE status = 'FINISHED'
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'FINISHED';

-- Query 5: WHERE service_date BETWEEN ...
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE service_date BETWEEN '2026-01-01' AND '2026-01-31';

-- Query 6: WHERE branch_id = 2
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE branch_id = 2;

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_a_branch_status_date;
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
DROP INDEX IF EXISTS idx_service_c_status_date_branch;