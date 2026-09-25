# Failure Modes and Mitigations

## 1. Failure Mode: Immediate Destructive Schema Alteration
- **Incident**: Engineer runs `ALTER TABLE customers DROP COLUMN phone;` or changes column type directly.
- **Impact**: All running application instances crash instantly with SQL syntax/missing column errors.
- **Mitigation**: CI/CD checks that lint database migrations to prevent `DROP COLUMN`, `RENAME COLUMN`, or non-nullable column additions without default values.

## 2. Failure Mode: Dual-Write Inconsistency (Split-Brain / Drift)
- **Incident**: An application writes to both Old and New tables. If the second write fails (or network times out), data in the two tables drifts out of sync.
- **Mitigation**:
  - Enclose dual writes in the same database transaction.
  - Run continuous data reconciliation scripts / integrity audits to detect discrepancies between old and new tables.

## 3. Failure Mode: Missing Backfill Records
- **Incident**: Reads are switched to the new table before backfill completes, returning 404/Empty for legacy rows.
- **Mitigation**:
  - Use fallback reading: if a record is not found in the new table, fetch from the old table and lazily write to the new table.
  - Implement checkpoint validation ensuring historical row count in old matches new before turning the read switch.

## 4. Failure Mode: Permanent Expand State (Contract Phase Abandoned)
- **Incident**: The engineering team migrates reads and writes, but forgets to clean up the legacy tables and dual-write logic. Technical debt piles up, wasting disk space and CPU cycles.
- **Mitigation**:
  - Set explicit deprecation deadlines.
  - Track metrics on legacy code usage; trigger automated alerts when zero legacy hits are logged for 30 consecutive days to prompt code removal.
