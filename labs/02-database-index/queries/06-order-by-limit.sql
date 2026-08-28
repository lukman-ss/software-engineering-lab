-- Lab 02: Database Index Foundation
-- Experiment: ORDER BY + LIMIT Optimization
-- Critical production pattern for dashboards and paginated results

-- The pattern: WHERE + ORDER BY + LIMIT
-- Key insight: Proper index can avoid explicit Sort AND stop early

-- Create index that can satisfy ORDER BY
CREATE INDEX idx_service_branch_status_date_desc
    ON service(branch_id, status, service_date DESC);

================================
-- QUERY: Dashboard-style - latest 20 finished services for branch 2
================================

EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- Record your observations:
-- - Index Scan or Seq Scan?
-- - Was Sort present or absent?
-- - Rows Examined: _____
-- - Actual Rows Returned: 20
-- - Execution Time: _____ ms
-- - Buffers: read _____, hit _____

-- Key insight: With proper index, EXPLAIN shows:
-- - Index Scan (not Seq Scan)
-- - Limit pushed down - stops after 20 rows
-- - No explicit Sort node
-- - Much faster than Seq Scan + Sort + LIMIT

================================
-- COMPARE: Without appropriate index
================================

-- Drop the index to simulate missing index
DROP INDEX idx_service_branch_status_date_desc;

-- Run same query
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- Record: How many rows were scanned? _____
-- Records showed: _____
-- Sort was: Present/Absent
-- Execution Time: _____ ms

-- This demonstrates the performance difference:
-- With index: ~hundreds of rows scanned
-- Without index: ~millions of rows scanned

================================
-- COMPARE: With single-column indexes
================================

-- Create three single-column indexes
CREATE INDEX idx_single_branch ON service(branch_id);
CREATE INDEX idx_single_status ON service(status);
CREATE INDEX idx_single_date ON service(service_date DESC);

-- Run same query
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
ORDER BY service_date DESC
LIMIT 20;

-- What plan did PostgreSQL choose?
-- - Bitmap Heap Scan? Seq Scan? Index Scan?
-- - Sort present?

-- Clean up
DROP INDEX IF EXISTS idx_single_branch;
DROP INDEX IF EXISTS idx_single_status;
DROP INDEX IF EXISTS idx_single_date;