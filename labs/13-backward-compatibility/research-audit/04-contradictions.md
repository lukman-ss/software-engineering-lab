# Contradictions Audit

No material contradictions found.

### Minor Nuance Observation
- **Statement A**: `research/04-database-migration.md` advises avoiding default constraints when adding columns to large tables due to full table rewrite risks in older database engines.
- **Statement B**: `research/07-deployment-and-rollback.md` mentions adding columns as "nullable, or with defaults".
- **Assessment**: Non-conflicting. Modern relational engines (PostgreSQL 11+, MySQL 8.0.12+) perform instant `ADD COLUMN ... DEFAULT` without table rewrites, whereas older systems require care. Both statements reflect valid context-dependent practices.
