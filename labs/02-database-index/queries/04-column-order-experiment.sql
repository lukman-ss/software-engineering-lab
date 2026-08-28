-- Lab 02: Database Index Foundation
-- Experiment: Composite Index Column Order
-- Demonstrates why column order matters in multicolumn B-tree indexes

================================
-- SETUP: Create three indexes with different column orders
================================

-- Index A: Optimal for our main query (equality + equality + range)
DROP INDEX IF EXISTS idx_service_a_branch_status_date;
CREATE INDEX idx_service_a_branch_status_date
    ON service(branch_id, status, service_date);

-- Index B: Wrong order for equality predicates (range first - bad)
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
CREATE INDEX idx_service_b_date_branch_status
    ON service(service_date, branch_id, status);

-- Index C: Wrong order (equality but different prefix)
DROP INDEX IF EXISTS idx_service_c_status_date_branch;
CREATE INDEX idx_service_c_status_date_branch
    ON service(status, service_date, branch_id);

-- The key insight: PostgreSQL can use a multicolumn B-tree index
-- when the query has conditions on an ORDERED PREFIX of indexed columns.
--
-- For index on (a, b, c):
-- - WHERE a = ?           ✓ Uses index
-- - WHERE a = ? AND b = ? ✓ Uses index
-- - WHERE a = ? AND b = ? AND c > ? ✓ Uses index
-- - WHERE b = ?           ✗ CANNOT use index (no prefix match)
-- - WHERE c = ?           ✗ CANNOT use index
-- - WHERE b = ? AND c = ? ✗ CANNOT use index

================================
-- QUERY 1: WHERE branch_id = 2
-- This query needs to understand how tight each index can make the scan
================================

-- With Index A (branch_id, status, service_date)
-- PostgreSQL knows branch_id = 2, can narrow scan to ~1/6 of index
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;

-- With Index B (service_date, branch_id, status)
-- PostgreSQL CANNOT use this index (branch_id is not leftmost)
-- Wait, let me check - actually it CAN skip to branch_id but...
-- No, wait! This is a common misconception. Let's verify:
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;
-- Actually, PostgreSQL can ONLY use index when prefix matches
-- So Index B cannot be used for WHERE branch_id = 2

-- With Index C (status, service_date, branch_id)
-- Same issue - branch_id is not leftmost
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2;

-- Record: Index A is the ONLY one usable for Query 1
-- Index B and C must do Sequential Scan
-- Why? Because B-tree can only be traversed when leftmost keys are constrained

================================
-- QUERY 2: WHERE branch_id = 2 AND status = 'FINISHED'
-- Understanding index prefix matching
================================

-- Index A: branch_id first, then status
-- PostgreSQL can first narrow to branch_id=2, then within that subset narrow to status='FINISHED'
-- This is called "reducing the index scan range" at each level
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED';

-- Index B: service_date first - cannot use at all for this query
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 2 AND status = 'FINISHED';

-- Index C: status first, then service_date, then branch_id
-- CAN use for this query! PostgreSQL can use status to find 'FINISHED'
-- Then sort within each status value, then filter branch_id
-- Wait - but the order matters for efficiency
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE status = 'FINISHED' AND branch_id = 2;

-- Key distinction:
-- - Reducing scan range: Index narrows WHERE clause
-- - Filtering: Index finds rows but needs extra checks

================================
-- QUERY 3: Full query with ORDER BY
-- The complete picture
================================

EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Which index works best and why?

================================
-- QUERY 4-6: Individual predicates
-- Understanding selectivity at each level
================================

-- Query 4: WHERE status = 'FINISHED'
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE status = 'FINISHED';

-- Query 5: WHERE service_date BETWEEN ...
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE service_date BETWEEN '2026-01-01' AND '2026-01-31';

-- Query 6: WHERE branch_id = 2
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*) FROM service WHERE branch_id = 2;

================================
-- CLEANUP
-- Drop test indexes to move to next experiment
================================

DROP INDEX IF EXISTS idx_service_a_branch_status_date;
DROP INDEX IF EXISTS idx_service_b_date_branch_status;
DROP INDEX IF EXISTS idx_service_c_status_date_branch;