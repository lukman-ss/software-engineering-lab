-- Lab 02: Database Index Foundation
-- Experiment: Functions on Indexed Columns (Expression Indexes)

-- KEY CONCEPT:
-- Wrapping an indexed column in a function typically prevents index usage.
-- Query structure determines whether PostgreSQL can use an index effectively.

================================
-- PART 1: Query With Expression
================================

-- Create standard index on service_date
DROP INDEX IF EXISTS idx_service_service_date;
CREATE INDEX idx_service_service_date ON service(service_date);

-- Query using DATE() function on indexed column
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE DATE(service_date) = DATE '2026-01-15';

-- Look at the plan: it's likely a Sequential Scan or Bitmap Heap Scan.
-- Why? The index stores raw `service_date` values.
-- The function `DATE(service_date)` changes the value.
-- PostgreSQL must apply DATE() to every row to check the condition.
-- Therefore, the index isn't directly searchable for `DATE '2026-01-15'`.

================================
-- PART 2: Query With Range Predicate
================================

-- Same intent, but rewritten to use a range condition (SARGable)
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE service_date >= '2026-01-15 00:00:00'
  AND service_date < '2026-01-16 00:00:00';

-- Look at the plan: Index Scan or Bitmap Index Scan!
-- Why? The query predicate uses the raw `service_date` column.
-- PostgreSQL can search the index for values between these bounds.

-- SARGable = Search Argument Able
-- A predicate is SARGable if the DBMS engine can take advantage of an index to speed up the execution of the query.
-- Using `DATE(column)` makes it non-SARGable.

================================
-- PART 3: Expression Indexes
================================

-- What if you CANNOT rewrite the query?
-- You can create an Expression Index.

DROP INDEX IF EXISTS idx_service_date_expr;
-- Create an index specifically for the `DATE()` function result
CREATE INDEX idx_service_date_expr ON service(DATE(service_date));

-- Re-run the expression query
EXPLAIN (ANALYZE, BUFFERS)
SELECT count(*)
FROM service
WHERE DATE(service_date) = DATE '2026-01-15';

-- Now the plan should show "Index Scan using idx_service_date_expr".
-- Because the index literally stores the results of DATE(service_date).

================================
-- PART 4: Trade-Off Analysis
================================

-- Should you always use expression indexes? NO!

-- Option A: Rewrite query to use SARGable predicates
-- Pros: Works with existing standard indexes. No extra index maintenance.
-- Cons: Requires changing application code.

-- Option B: Create Expression Index
-- Pros: Speeds up existing non-SARGable queries without app code changes.
-- Cons: Adds another index to maintain (write overhead, storage).
--       Index is highly specific to that exact function expression.

-- The Trade-off:
-- - Change Query vs Add Expression Index
-- - Generally, rewriting to SARGable is preferred.
-- - Use expression indexes only when rewriting is impossible or complex.

-- Cleanup
DROP INDEX IF EXISTS idx_service_service_date;
DROP INDEX IF EXISTS idx_service_date_expr;