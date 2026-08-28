-- Lab 02: Database Index Foundation
-- Experiment: Write Cost of Indexes (Inserts, Updates, Deletes)
-- Proves that indexes are not free and introduce write amplification

-- KEY INSIGHT: Index design is a read/write trade-off.
-- Every secondary index requires:
-- 1. Additional B-tree maintenance on INSERT
-- 2. Additional page writes and disk I/O
-- 3. Additional WAL (Write-Ahead Log) generation
-- 4. Increased buffer cache pressure
-- 5. Write amplification (1 row change = multiple index page changes)

================================
-- PART 1: Measure INSERT Performance under different index counts
================================

-- Scenario 1: No secondary indexes (table + primary key only)
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status;
DROP INDEX IF EXISTS idx_service_service_date;

-- Measure INSERT time for 1,000 rows without indexes
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

-- Scenario 2: With one composite index
CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date DESC);

-- Measure INSERT time for 1,000 rows with 1 secondary index
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

-- Scenario 3: With several secondary indexes (over-indexed anti-pattern)
CREATE INDEX idx_service_status_single ON service(status);
CREATE INDEX idx_service_date_single ON service(service_date);
CREATE INDEX idx_service_customer ON service(customer_id);

-- Measure INSERT time for 1,000 rows with 4 secondary indexes
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

================================
-- PART 2: UPDATE Performance: Non-indexed vs Indexed Column
================================

-- Scenario A: UPDATE a non-indexed column (e.g., mechanic_id)
-- PostgreSQL can use Heap-Only Tuple (HOT) updates if space permits on the same page!
\timing on
UPDATE service
SET mechanic_id = 99
WHERE invoice_no = 'TEST-INS3-1';
\timing off

-- Scenario B: UPDATE an indexed column (e.g., status)
-- Forces index modification for every matched row
\timing on
UPDATE service
SET status = 'CANCELLED'
WHERE invoice_no = 'TEST-INS3-2';
\timing off

-- HOT (Heap-Only Tuple) Updates concept:
-- When you update a row in PostgreSQL:
-- 1. If NO indexed columns change AND the new row fits on the same disk page,
--    PostgreSQL updates only the table heap page without touching secondary indexes!
-- 2. This avoids write amplification and index bloat.
-- 3. If an indexed column changes, PostgreSQL MUST update every secondary index
--    containing that column.

================================
-- PART 3: DELETE Performance
================================

\timing on
DELETE FROM service WHERE invoice_no LIKE 'TEST-INS%';
\timing off

-- Deletes require removing keys from ALL secondary indexes,
-- which can cause index page splits, fragmentation, and bloat.

-- Cleanup test data and indexes
DELETE FROM service WHERE invoice_no LIKE 'TEST-INS%';
DROP INDEX IF EXISTS idx_service_branch_status_date;
DROP INDEX IF EXISTS idx_service_status_single;
DROP INDEX IF EXISTS idx_service_date_single;
DROP INDEX IF EXISTS idx_service_customer;