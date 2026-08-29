-- Lab 02: Database Index Foundation
-- Experiment: Write Cost of Indexes (Inserts, Updates, Deletes)
-- Proves that indexes are not free and introduce write amplification

-- KEY INSIGHT: Index design is a read/write trade-off.
-- Every secondary index requires:
-- 1. Additional B-tree maintenance on INSERT (page splits, WAL)
-- 2. Additional disk I/O and buffer cache pressure
-- 3. Write amplification (1 row change = multiple index page changes)

-- ============================================
-- PART 1: Measure INSERT Performance under different index counts
-- ============================================

-- RESET: Clean state for fair comparison
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status_single;
DROP INDEX IF EXISTS idx_service_date_single;
DROP INDEX IF EXISTS idx_service_customer;
TRUNCATE TABLE service RESTART IDENTITY;

-- Scenario 1: No secondary indexes (table + primary key only)
-- Measure INSERT time for 1000 rows without indexes
\timing on
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    (random() * 5)::int + 1,
    (random() * 500)::int + 1,
    (random() * 50)::int + 1,
    'FINISHED',
    '2026-06-01'::date + (random() * 30)::int,
    'TEST-INS-' || gs,
    NOW()
FROM generate_series(1, 1000) AS gs;
\timing off

-- Record: INSERT time without indexes = _____ ms

-- RESET: Clean state
TRUNCATE TABLE service RESTART IDENTITY;
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status_single;
DROP INDEX IF EXISTS idx_service_date_single;
DROP INDEX IF EXISTS idx_service_customer;

-- Scenario 2: With one composite index
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date DESC);

-- Measure INSERT time for 1000 rows with 1 secondary index
\timing on
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    (random() * 5)::int + 1,
    (random() * 500)::int + 1,
    (random() * 50)::int + 1,
    'FINISHED',
    '2026-06-01'::date + (random() * 30)::int,
    'TEST-INS2-' || gs,
    NOW()
FROM generate_series(1, 1000) AS gs;
\timing off

-- Record: INSERT time with 1 index = _____ ms

-- RESET: Clean state
TRUNCATE TABLE service RESTART IDENTITY;

-- Scenario 3: With several secondary indexes (over-indexed anti-pattern)
CREATE INDEX idx_service_status_single ON service(status);
CREATE INDEX idx_service_service_date ON service(service_date);
CREATE INDEX idx_service_customer ON service(customer_id);

-- Measure INSERT time for 1000 rows with 4 secondary indexes
\timing on
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    (random() * 5)::int + 1,
    (random() * 500)::int + 1,
    (random() * 50)::int + 1,
    'FINISHED',
    '2026-06-01'::date + (random() * 30)::int,
    'TEST-INS3-' || gs,
    NOW()
FROM generate_series(1, 1000) AS gs;
\timing off

-- Record: INSERT time with 4 indexes = _____ ms

-- ============================================
-- PART 2: UPDATE Performance: Non-indexed vs Indexed Column
-- ============================================

-- Scenario A: UPDATE a non-indexed column (e.g., mechanic_id)
-- PostgreSQL can use Heap-Only Tuple (HOT) updates if space permits on the same page!
-- HOT updates modify only the heap, not secondary indexes.
\timing on
UPDATE service
SET mechanic_id = 99
WHERE invoice_no = 'TEST-INS3-1';
\timing off

-- Record: HOT UPDATE time = _____ ms

-- Scenario B: UPDATE an indexed column (e.g., status)
-- Forces index modification for EVERY matched row
\timing on
UPDATE service
SET status = 'CANCELLED'
WHERE invoice_no = 'TEST-INS3-2';
\timing off

-- Record: Indexed column UPDATE time = _____ ms

-- HOT (Heap-Only Tuple) Updates concept:
-- When you UPDATE a row in PostgreSQL:
-- 1. If NO indexed columns change AND the new row fits on the same disk page,
--    PostgreSQL creates a new heap tuple and marks the old one dead.
--    Only the heap page is modified (no index updates needed).
-- 2. This avoids write amplification for that row.
-- 3. If an indexed column changes, PostgreSQL MUST update every secondary index
--    containing that column.

-- IMPORTANT NOTE: DELETE does NOT 'remove keys from indexes' directly.
-- DELETE creates a dead heap tuple version. Index entries remain until VACUUM.
-- This contributes to index bloat over time.

-- ============================================
-- PART 3: DELETE Performance
-- ============================================

\timing on
DELETE FROM service WHERE invoice_no LIKE 'TEST-INS%';
\timing off

-- Record: DELETE 1000 rows = _____ ms

-- NOTE: DELETE behavior:
-- - PostgreSQL marks heap tuples as deleted (dead tuples)
-- - Index entries remain in place until VACUUM processes them
-- - This can lead to index bloat over time
-- - VACUUM later removes the index entries
-- - DELETE does NOT cause immediate index page splits (that's an INSERT phenomenon)

-- Cleanup test data and indexes
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status_single;
DROP INDEX IF EXISTS idx_service_service_date;
DROP INDEX IF EXISTS idx_service_customer;