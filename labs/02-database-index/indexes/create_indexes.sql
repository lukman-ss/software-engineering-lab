-- Lab 02: Database Index Foundations
-- Index definitions for experimentation
-- Run only those needed for the current experiment

-- PostgreSQL 16 compatibility note:
-- This repository targets PostgreSQL 16. PostgreSQL 18 introduced B-tree skip-scan
-- optimization, which is NOT assumed in these explanations.

-- ============================================
-- Single-column indexes (for Experiment 02)
-- ============================================
-- Individual indexes — PostgreSQL may combine them with BitmapAnd

-- CREATE INDEX idx_service_branch_id
--     ON service(branch_id);
--
-- CREATE INDEX idx_service_status
--     ON service(status);
--
-- CREATE INDEX idx_service_service_date
--     ON service(service_date);

-- ============================================
-- Composite index (for Experiment 03)
-- ============================================
-- For this query shape (branch_id + status + date range/order),
-- equality columns followed by the range/order column are a strong candidate.
-- Treat this as workload-driven guidance, not a universal indexing rule.

-- Note: PostgreSQL B-tree indexes can be scanned in both forward and backward
-- directions, so explicit DESC is NOT required solely to satisfy ORDER BY
-- service_date DESC when the preceding equality columns branch_id and status
-- are fixed. DESC becomes relevant only for mixed-order requirements like
-- ORDER BY x ASC, y DESC, or for index key sort order in conflict with queries.

CREATE INDEX idx_service_branch_status_date
    ON service(branch_id, status, service_date);

-- ============================================
-- Alternative: Separate indexes for flexibility
-- ============================================
-- These can be useful if different queries filter
-- different columns

-- CREATE INDEX idx_service_branch_id_status
--     ON service(branch_id, status);
--
-- CREATE INDEX idx_service_branch_date
--     ON service(branch_id, service_date);

-- ============================================
-- Partial index example (for specific use case)
-- ============================================
-- IN_PROGRESS is approximately 5% of the canonical dataset (50,000 rows).
-- A partial index on this subset demonstrates the storage/maintenance
-- benefits of indexing a small, stable predicate.

-- CREATE INDEX idx_service_in_progress_branch_date
--     ON service(branch_id, service_date DESC)
--     WHERE status = 'IN_PROGRESS';