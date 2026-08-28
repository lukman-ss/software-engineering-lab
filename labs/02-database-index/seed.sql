-- Lab 02: Database Index Foundation
-- Realistic large dataset using generate_series
-- Target: ~500,000 rows

-- Clear any existing data
TRUNCATE TABLE service RESTART IDENTITY CASCADE;

-- Insert ~500,000 rows with realistic distributions
-- Uses generate_series for reproducible, fast generation
-- Include rare status for selectivity experiments

INSERT INTO service (branch_id, customer_id, mechanic_id, status, service_date, invoice_no, created_at)
SELECT
    -- branch_id: skewed distribution
    -- Branch 2 and 5 are busiest (~20% each), others ~10%
    CASE (random() * 10)::int
        WHEN 0, 1 THEN 1  -- 20%
        WHEN 2, 3 THEN 2  -- 20% - our query target
        WHEN 4, 5 THEN 3  -- 20%
        WHEN 6, 7 THEN 4  -- 20%
        ELSE 5             -- 20%
    END,

    -- customer_id: 1-500 (realistic range)
    1 + (random() * 499)::int,

    -- mechanic_id: 1-50 (small team)
    1 + (random() * 49)::int,

    -- status: FINISHED common (~75%), CANCELLED (~24%), rare status (~0.1%)
    CASE (random() * 1000)::int
        WHEN 0 THEN 'PENDING_REFUND'  -- 0.1% - rare status for selectivity experiments
        WHEN 1, 2, 3, 4, 5, 6, 7, 8   -- 0.8% total - rare statuses
        WHEN 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
             21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
             31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
             41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
             51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
             61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
             71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
             81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
             91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
             101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
             111, 112, 113, 114, 115, 116, 117, 118, 119, 120,
             121, 122, 123, 124, 125, 126, 127, 128, 129, 130,
             131, 132, 133, 134, 135, 136, 137, 138, 139, 140 THEN 'FINISHED'  -- 75%
        WHEN 141, 142, 143, 144, 145, 146, 147, 148, 149, 150,
             151, 152, 153, 154, 155, 156, 157, 158, 159, 160,
             161, 162, 163, 164, 165, 166, 167, 168, 169, 170,
             171, 172, 173, 174, 175, 176, 177, 178, 179, 180,
             181, 182, 183, 184, 185, 186, 187, 188, 189, 190,
             191, 192, 193, 194, 195, 196, 197, 198, 199, 200,
             201, 202, 203, 204, 205, 206, 207, 208, 209, 210,
             211, 212, 213, 214, 215, 216, 217, 218, 219, 220,
             221, 222, 223, 224, 225, 226, 227, 228, 229, 230 THEN 'CANCELLED'  -- ~24%
        ELSE 'WAITING'  -- remaining ~1%
    END,

    -- service_date: spread across 2025-2026
    -- Concentrated in 2026 (our query target)
    '2025-01-01'::date + (random() * 730)::int,

    -- invoice_no: unique-ish format
    'INV-' || lpad((100000 + gs)::text, 6, '0'),

    -- created_at: service_date + random hours
    ('2025-01-01'::date + (random() * 730)::int) + (random() * 24 || ' hours')::interval

FROM generate_series(1, 500000) AS gs;

-- Update statistics for accurate query planning
ANALYZE service;

-- Verify data volume
SELECT COUNT(*) AS total_rows FROM service;

-- Distribution checks
SELECT status, COUNT(*) AS count
FROM service
GROUP BY status
ORDER BY count DESC;

SELECT branch_id, COUNT(*) AS count
FROM service
GROUP BY branch_id
ORDER BY branch_id;

-- Statistics inspection
SELECT
    attname,
    n_distinct,
    array_to_string(most_common_vals, ', ') AS most_common_vals,
    array_to_string(most_common_freqs, ', ') AS most_common_freqs
FROM pg_stats
WHERE tablename = 'service'
AND attname = 'status';