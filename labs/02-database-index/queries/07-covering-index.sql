-- Lab 02: Database Index Foundation
-- Experiment: SELECT * and Covering Indexes (INCLUDE)
--
-- Why is `SELECT *` sometimes bad?
-- Not universally! But it specifically impacts index usage:
-- - SELECT * often prevents an Index Only Scan with a narrowly designed
--   covering index because all referenced columns must be available from the index.
--   (A theoretical index containing every column could permit it, but such an
--   index may be impractical and is rarely desirable.)
-- - SELECT * forces fetching from heap (table pages)
-- - Increases I/O and network transfer
--
-- Contrast with column-specific queries that match the index key columns.

-- ============================================
-- CONTROLLED QUERY: Columns included in index key
-- ============================================

-- Imagine a dashboard only needs these specific fields
SELECT branch_id, status, service_date, customer_id, mechanic_id, invoice_no
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC
LIMIT 20;

-- ============================================
-- EXPERIMENT WITH COVERING INDEX (INCLUDE)
-- ============================================

-- PostgreSQL's INCLUDE feature adds non-key columns to leaf nodes
-- Allows index to satisfy query without visiting heap (table)
-- when the visibility map conditions are met.

DROP INDEX IF EXISTS idx_service_dashboard;
CREATE INDEX idx_service_dashboard
ON service(branch_id, status, service_date DESC)
INCLUDE (customer_id, mechanic_id, invoice_no);
ANALYZE service;

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
-- - "Heap Fetches: X" (ideally 0 or small when VM is up to date)

-- ============================================
-- COMPARE: WITH SELECT *
-- ============================================

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
-- - Degraded from "Index Only Scan" (if it was one) to "Index Scan"
-- - Buffer "shared read/hit" count likely increased
-- - Heap fetches now necessary because column not in index

-- ============================================
-- THE VISIBILITY MAP
-- ============================================

-- IMPORTANT: INCLUDE does NOT guarantee Index Only Scan!
-- PostgreSQL still needs to check if row is visible (not deleted/updated)
-- It uses the Visibility Map to know if it can skip the heap fetch
-- If Visibility Map says page might have unvacuumed changes,
-- PostgreSQL MUST visit heap anyway (Heap Fetches > 0)
-- VACUUM updates the visibility map, making Index Only Scan more likely.

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

-- ============================================
-- TRADE-OFFS OF INCLUDE
-- ============================================

-- Why not INCLUDE every column?
-- 1. Index becomes huge (storage cost)
-- 2. INSERTs must write more data to index (write penalty)
-- 3. UPDATEs to INCLUDEd columns must update index (HOT update penalty)

-- ============================================
-- CLEANUP
-- ============================================

DROP INDEX IF EXISTS idx_service_dashboard;