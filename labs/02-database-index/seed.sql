-- Lab 02: Database Index Foundation
-- Exactly 500,000 rows using generate_series
-- Deterministic modulo-based distribution (exact percentages, no random())

-- Clear any existing data
TRUNCATE TABLE service RESTART IDENTITY CASCADE;

-- Insert exactly 500,000 rows
-- Deterministic distribution using two separate deterministic modulo-based permutations:
--
-- Branch distribution (per 10,000-row cycle, gs % 10000):
--   mod 0-1499   = branch 1   15.00%
--   mod 1500-3999 = branch 2  25.00%
--   mod 4000-5999 = branch 3  20.00%
--   mod 6000-7499 = branch 4  15.00%
--   mod 7500-8999 = branch 5  15.00%
--   mod 9000-9999 = branch 6  10.00%
--
-- Status distribution (per 10,000-row cycle, (gs * 7) % 10000):
--   mod 0-6999    = FINISHED       70.00%  → 350,000 total
--   mod 7000-8999 = CANCELLED      20.00%  → 100,000 total
--   mod 9000-9499 = IN_PROGRESS     5.00%   →  25,000 total
--   mod 9500-9998 = WAITING         4.90%   →  24,500 total
--   mod 9999      = PENDING_REFUND  0.10%   →     500 total
--
-- The multiplier 7 is coprime with 10000 (gcd(7, 10000) = 1), ensuring
-- that (gs * 7) % 10000 produces a permutation of 0-9999 in each
-- 10,000-row cycle. This preserves exact marginal status percentages
-- while decoupling status from branch_id (which uses gs % 10000 directly).
--
-- created_at is derived from service_date: same base date + gs-derived hours
-- service_date spans 2025-01-01 to 2026-12-31 (730 days)

INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    CASE
        WHEN (gs % 10000) < 1500  THEN 1
        WHEN (gs % 10000) < 4000  THEN 2
        WHEN (gs % 10000) < 6000  THEN 3
        WHEN (gs % 10000) < 7500  THEN 4
        WHEN (gs % 10000) < 9000  THEN 5
        ELSE 6
    END AS branch_id,
    (gs % 500) + 1 AS customer_id,
    (gs % 50) + 1 AS mechanic_id,
    CASE
        WHEN ((gs * 7) % 10000) < 7000  THEN 'FINISHED'
        WHEN ((gs * 7) % 10000) < 9000  THEN 'CANCELLED'
        WHEN ((gs * 7) % 10000) < 9500  THEN 'IN_PROGRESS'
        WHEN ((gs * 7) % 10000) < 9999  THEN 'WAITING'
        ELSE 'PENDING_REFUND'
    END AS status,
    '2025-01-01'::date + (gs % 730) AS service_date,
    'INV-' || lpad((gs + 100000)::text, 6, '0') AS invoice_no,
    ('2025-01-01'::date + (gs % 730))::timestamp + ((gs % 24) * interval '1 hour') AS created_at
FROM generate_series(1, 500000) AS gs;

-- Update statistics for accurate query planning
ANALYZE service;

-- ============================================================
-- VALIDATION QUERIES
-- ============================================================

SELECT '=== Row Count ===' AS section;
SELECT COUNT(*) AS total_rows FROM service;

SELECT '=== Status Distribution ===' AS section;
SELECT status, COUNT(*) AS cnt,
       ROUND(COUNT(*) * 100.0 / 500000, 2) AS percentage
FROM service GROUP BY status ORDER BY cnt DESC;

SELECT '=== Branch Distribution ===' AS section;
SELECT branch_id, COUNT(*) AS cnt,
       ROUND(COUNT(*) * 100.0 / 500000, 2) AS percentage
FROM service GROUP BY branch_id ORDER BY branch_id;

SELECT '=== Service Date Range ===' AS section;
SELECT MIN(service_date) AS min_service_date,
       MAX(service_date) AS max_service_date
FROM service;

SELECT '=== Created At Range ===' AS section;
SELECT MIN(created_at) AS min_created_at,
       MAX(created_at) AS max_created_at
FROM service;

-- ============================================================
-- JOINT DISTRIBUTION VALIDATION
-- ============================================================
-- Marginal distributions are deterministic.
-- Joint distributions are measured rather than assumed independent.

SELECT
    branch_id,
    status,
    count(*) AS cnt,
    round(
        count(*) * 100.0 /
        sum(count(*)) OVER (PARTITION BY branch_id),
        2
    ) AS pct_within_branch
FROM service
GROUP BY branch_id, status
ORDER BY branch_id, status;

-- Verify: branch 2 contains multiple statuses, and FINISHED is not 100% of branch 2.
-- This ensures the primary composite-index experiment has meaningful data.
DO $$
DECLARE
    distinct_statuses INTEGER;
    finished_pct      NUMERIC;
BEGIN
    SELECT COUNT(DISTINCT status)
    INTO distinct_statuses
    FROM service
    WHERE branch_id = 2;

    ASSERT distinct_statuses > 1,
        'Branch 2 should contain multiple statuses (got ' || distinct_statuses || ')';

    SELECT COUNT(*) FILTER (WHERE status = 'FINISHED') * 100.0 / COUNT(*)
    INTO finished_pct
    FROM service
    WHERE branch_id = 2;

    ASSERT finished_pct < 100,
        'branch_id=2 must not imply status=FINISHED (got ' || ROUND(finished_pct, 2) || '% FINISHED)';
END $$;

-- ============================================================
-- ASSERTIONS (fail if marginal distributions deviate from spec)
-- ============================================================

DO $$
BEGIN
    ASSERT (SELECT COUNT(*) FROM service) = 500000,
        'Row count mismatch: expected 500000';
END $$;

DO $$
DECLARE
    finished_pct      NUMERIC;
    pr_pct            NUMERIC;
    in_progress_pct   NUMERIC;
    branch2_pct       NUMERIC;
BEGIN
    SELECT COUNT(*) FILTER (WHERE status = 'FINISHED') * 100.0 / 500000
    INTO finished_pct FROM service;
    ASSERT finished_pct = 70.0,
        'FINISHED expected 70.0%, got ' || ROUND(finished_pct, 2);

    SELECT COUNT(*) FILTER (WHERE status = 'PENDING_REFUND') * 100.0 / 500000
    INTO pr_pct FROM service;
    ASSERT pr_pct = 0.1,
        'PENDING_REFUND expected 0.1%, got ' || pr_pct;

    SELECT COUNT(*) FILTER (WHERE status = 'IN_PROGRESS') * 100.0 / 500000
    INTO in_progress_pct FROM service;
    ASSERT in_progress_pct = 5.0,
        'IN_PROGRESS expected 5.0%, got ' || ROUND(in_progress_pct, 2);

    SELECT COUNT(*) FILTER (WHERE branch_id = 2) * 100.0 / 500000
    INTO branch2_pct FROM service;
    ASSERT branch2_pct = 25.0,
        'Branch 2 expected 25.0%, got ' || ROUND(branch2_pct, 2);
END $$;

SELECT '=== All assertions passed ===' AS result;