-- Lab 02: Database Index Foundation
-- Cleanup script to reset the lab environment

-- Drop the table (this cascades to dependent objects)
DROP TABLE IF EXISTS service CASCADE;

-- The schema and seed scripts can be re-run with:
--   psql -f schema.sql
--   psql -f seed.sql