-- Lab 02: Database Index Foundation
-- Queries for investigating the main performance problem

-- The primary query we need to optimize
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;

-- Count matching rows (for understanding data volume)
SELECT COUNT(*) as matching_rows
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31';

-- Distribution check - how many rows per branch?
SELECT branch_id, COUNT(*) as total_per_branch
FROM service
GROUP BY branch_id
ORDER BY branch_id;