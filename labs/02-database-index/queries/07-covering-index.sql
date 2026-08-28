-- Lab 02: Database Index Foundation
-- Experiment: SELECT * and Covering Indexes (INCLUDE)

-- Why is `SELECT *` sometimes bad?
-- Not universally bad! But specifically impacts index usage:
-- 1. Forces fetching from heap (table pages)
-- 2. Prevents Index Only Scans
-- 3. Increases I/O and network transfer

================================
-- CONTROLLED QUERY: Select only needed columns
================================

-- Imagine a dashboard only needs these specific fields
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC
LIMIT 20;

================================
-- EXPERIMENT WITH COVERING INDEX (INCLUDE)
================================

-- PostgreSQL's INCLUDE feature adds non-key columns to leaf nodes
-- Allows index to satisfy query without visiting heap (table)

CREATE INDEX idx_service_dashboard
ON service(branch_id, status, service_date DESC)
INCLUDE (customer_id, mechanic_id, invoice_no);

-- Run the query
EXPLAIN (ANALYZE, BUFFERS)
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC
LIMIT 20;

-- Look for:
-- - "Index Only Scan"
-- - "Heap Fetches: X" (ideally 0 or small)

================================
-- COMPARE: WITH SELECT *
================================

-- Now run the same query but fetch ALL columns
-- The index doesn't have `id` or `created_at`
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC
LIMIT 20;

-- What changed?
-- - Degraded from "Index Only Scan" to "Index Scan"
-- - Buffer "shared read/hit" count likely increased

================================
-- THE VISIBILITY MAP
================================

-- IMPORTANT: INCLUDE does NOT guarantee Index Only Scan!
-- PostgreSQL still needs to check if row is visible (not deleted/updated)
-- It uses the Visibility Map to know if it can skip the heap fetch
-- If Visibility Map says page might have unvacuumed changes,
-- PostgreSQL MUST visit heap anyway (Heap Fetches > 0)

-- Check current visibility map status (approximate)
-- VACUUM helps update visibility map
VACUUM service;

-- Try Index Only query again after VACUUM
-- Heap Fetches should drop
EXPLAIN (ANALYZE, BUFFERS)
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC
LIMIT 20;

================================
-- TRADE-OFFS OF INCLUDE
================================

-- Why not INCLUDE every column?
-- 1. Index becomes huge (storage cost)
-- 2. INSERTs must write more data to index (write penalty)
-- 3. UPDATEs to INCLUDEd columns must update index (HOT update penalty)

-- Clean up
DROP INDEX idx_service_dashboard;