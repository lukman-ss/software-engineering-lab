-- Lab 02: Database Index Foundations
-- Index definitions for experimentation
-- Run only those needed for the current experiment

-- PostgreSQL 16 compatibility note:
-- This repository targets PostgreSQL 16. PostgreSQL 18 introduced B-tree skip-scan
-- optimization, which is NOT assumed in these explanations.

-- ============================================
-- Single-column indexes (for Experiment 02)
-- ============================================
-- Individual indexes - PostgreSQL may combine with BitmapAnd

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
-- Column order: equality columns first, then range column
-- Note: PostgreSQL 16 B-tree indexes can be scanned in both directions,
-- so explicit DESC is NOT required for ORDER BY service_date DESC.
-- DESC becomes relevant only for mixed-order requirements like
-- ORDER BY x ASC, y DESC.

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
-- Partial index (for specific use case)
-- ============================================
-- Only index FINISHED services for branch 2

-- CREATE INDEX idx_service_branch2_finished
--     ON service(branch_id, service_date DESC)
--     WHERE branch_id = 2 AND status = 'FINISHED';