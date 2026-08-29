-- Lab 02: Database Index Foundation
-- Realistic ~500,000 row dataset using generate_series
-- Uses VALID PostgreSQL 16 syntax with simple modulo-based distribution

-- Clear any existing data
TRUNCATE TABLE service RESTART IDENTITY CASCADE;

-- Insert ~500,000 rows with realistic distributions
-- Using modulo arithmetic for clean percentage-based distribution
INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    -- branch_id distribution (approximate percentages as documented)
    -- Branch 1: ~15%, Branch 2: ~25%, Branch 3: ~20%, Branch 4: ~15%, Branch 5: ~15%, Branch 6: ~10%
    CASE
        WHEN (gs % 20) < 3   THEN 1  -- ~15%
        WHEN (gs % 20) < 8   THEN 2  -- ~25% (3+5)
        WHEN (gs % 20) < 12  THEN 3  -- ~20% (8+4)
        WHEN (gs % 20) < 15  THEN 4  -- ~15% (12+3)
        WHEN (gs % 20) < 18  THEN 5  -- ~15% (15+3)
        ELSE                        6  -- ~10% (18+2)
    END AS branch_id,

    -- customer_id: 1-500 (realistic range)
    (gs % 500) + 1 AS customer_id,

    -- mechanic_id: 1-50 (small team)
    (gs % 50) + 1 AS mechanic_id,

    -- status distribution matching documented values:
    -- FINISHED: ~75%, CANCELLED: ~24%, WAITING: ~0.9%, PENDING_REFUND: ~0.1%
    CASE
        WHEN (gs % 1000) = 0                         THEN 'PENDING_REFUND'  -- 0.1%
        WHEN (gs % 1000) BETWEEN 1 AND 9              THEN 'WAITING'         -- 0.9%
        WHEN (gs % 1000) BETWEEN 10 AND 84           THEN 'FINISHED'        -- 75% (10+75=85)
        ELSE 'CANCELLED'                              -- 15% (85+15=100)
    END AS status,

    -- service_date: spread across 2025-2026 (730 days)
    '2025-01-01'::date + (gs % 730) AS service_date,

    -- invoice_no: unique format
    'INV-' || lpad((gs + 100000)::text, 6, '0') AS invoice_no,

    -- created_at: logically derived from service_date + random hours (0-23)
    ('2025-01-01'::date + (gs % 730)) + ((gs % 24) * interval '1 hour') AS created_at

FROM generate_series(1, 500000) AS gs;

-- Update statistics for accurate query planning
ANALYZE service;

-- Validation output
SELECT COUNT(*) AS total_rows FROM service;

SELECT status, COUNT(*) AS cnt,
       ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM service), 2) AS percentage
FROM service GROUP BY status ORDER BY cnt DESC;

SELECT branch_id, COUNT(*) AS cnt,
       ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM service), 2) AS percentage
FROM service GROUP BY branch_id ORDER BY branch_id;