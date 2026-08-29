-- Lab 02: Database Index Foundation
-- Experiment: PostgreSQL Planner Statistics and ANALYZE

-- KEY CONCEPT:
-- Indexes do NOT make decisions.
-- The PostgreSQL planner chooses the execution plan based on its statistical model.
-- `ANALYZE` updates these statistics.

-- ================================
-- PART 1: Planner Statistics (`pg_stats`)
-- ================================

-- Important statistics columns in `pg_stats`:
-- n_distinct         : Number (or fraction) of distinct values in a column.
-- most_common_vals   : Common values in the column.
-- most_common_freqs  : Frequency of those common values.
-- histogram_bounds   : Boundary values dividing the data into equal groups (buckets).
-- null_frac          : Fraction of rows where this column is NULL.

-- Inspect statistics for `status` (low-cardinality column)
SELECT
    attname,
    n_distinct,
    most_common_vals,
    most_common_freqs,
    histogram_bounds,
    null_frac
FROM pg_stats
WHERE tablename = 'service'
  AND attname = 'status';

-- Inspect statistics for `service_date` (high-cardinality column)
SELECT
    attname,
    n_distinct,
    most_common_vals,
    most_common_freqs,
    histogram_bounds,
    null_frac
FROM pg_stats
WHERE tablename = 'service'
  AND attname = 'service_date';

-- ================================
-- PART 2: Insert Substantial Data Without ANALYZE
-- ================================

-- Drop existing indexes to observe planner choices based on stale statistics
DROP INDEX IF EXISTS idx_service_branch_status_date;

-- Insert 50,000 new rows for a new branch
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    9 AS branch_id,                          -- Brand new branch
    1 + (random() * 499)::int,
    1 + (random() * 49)::int,
    'FINISHED',
    '2026-08-01'::date + (random() * 30)::int,
    'STATS-' || gs,
    NOW()
FROM generate_series(1, 50000) AS gs;

-- Examine the plan BEFORE ANALYZE (planner uses stale statistics)
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 9;

-- Compare the `rows` (estimate) with the actual row count.
-- Notice if the estimate is wildly incorrect because ANALYZE hasn't run since the insert.

-- ================================
-- PART 3: Run ANALYZE and Observe Plan Change
-- ================================

ANALYZE service;

-- Examine the plan AFTER ANALYZE (planner now has accurate stats)
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service WHERE branch_id = 9;

-- Verify the statistics were updated
SELECT
    attname,
    n_distinct,
    most_common_vals,
    most_common_freqs
FROM pg_stats
WHERE tablename = 'service'
  AND attname = 'branch_id';

-- ================================
-- PART 4: Estimated vs Actual Rows
-- ================================

-- Compare the planner's estimates with reality.
-- Look for large discrepancies (indicating stale stats or skewed data).

EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM service
WHERE branch_id = 5
  AND status = 'PENDING_REFUND'
  AND service_date >= '2026-01-01';

-- Compare the rows (estimate) in the EXPLAIN output with the `actual rows` shown.

-- Key insight:
-- If `rows` (estimate) and `actual rows` diverge significantly, the planner might be making
-- decisions based on bad assumptions. Running `ANALYZE` frequently helps keep this accurate.
-- BUT: Running ANALYZE does NOT guarantee a plan change!

-- ================================
-- CLEANUP
-- ================================

DELETE FROM service WHERE invoice_no LIKE 'STATS-%';
ANALYZE service;