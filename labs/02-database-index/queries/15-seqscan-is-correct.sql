-- Lab 02: Database Index Foundation
-- Experiment: When Seq Scan IS Correct
-- Destroys the misconception: "Seq Scan = Bad, Index Scan = Good"

-- PRIMARY LESSON:
-- The goal is NOT to force index usage.
-- The goal is to MINIMIZE TOTAL QUERY COST for the real workload.

================================
-- PART 1: Ensure Index Exists
================================

-- Create the index
DROP INDEX IF EXISTS idx_service_branch_status_date;
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date DESC);

-- Sanity check: total rows in the table
SELECT count(*) FROM service;

================================
-- PART 2: Low Selectivity Query (Sequential Scan should win)
================================

-- Query: Fetch most FINISHED services (highly common, ~75% of rows)
-- Predicate selectivity is ~0.75 -> very low
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'FINISHED';

-- EXPECTATION: PostgreSQL should choose Seq Scan!
-- REASON: Returning 75% of the table means the database would have to read
--         almost every heap page anyway, plus pay the cost of traversing the index.
--         A simple sequential read of the heap is faster.

-- Record: Plan was Seq Scan / Index Scan?
-- Execution Time: _____ ms
-- Rows Examined: _____

================================
-- PART 3: High Selectivity Query (Index Scan should win)
================================

-- Query: Find a single service by rare branch and ID
-- Predicate selectivity is ~0.0001 (very high)
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 9
  AND status = 'PENDING_REFUND';

-- EXPECTATION: PostgreSQL should choose Index Scan or Bitmap Index Scan!
-- REASON: Very few rows match. Reading from the index is much cheaper than
--         scanning the entire table.

-- Record: Plan was Seq Scan / Index Scan / Bitmap?
-- Execution Time: _____ ms
-- Rows Examined: _____

================================
-- PART 4: Highly selective range query
================================

-- Query: Fetch very few rows with exact primary key match
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE id = 12345;

-- EXPECTATION: Index Scan or Index Only Scan on PK
-- REASON: PK index directly points to the single row.

================================
-- PART 5: Forced index comparison
================================

-- What if we FORCE index usage even when it's bad?
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE status = 'FINISHED'
ORDER BY service_date;

-- Compare: Without index hinting (planner chooses)
EXPLAIN (ANALYZE, BUFFERS)
SELECT /*+ NoIndex(service idx_service_branch_status_date) */ *
FROM service
WHERE status = 'FINISHED';

-- Compare the execution times!
-- Did forcing or blocking the index help or hurt performance?

-- Verify PostgreSQL's choice for Query 2
-- If PostgreSQL chose Seq Scan, it's CORRECT.
-- The planner optimizes for the total cost, not for showing off index usage.

================================
-- CLEANUP
================================

DROP INDEX IF EXISTS idx_service_branch_status_date;