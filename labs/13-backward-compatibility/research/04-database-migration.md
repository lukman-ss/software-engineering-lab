# Database Schema Evolution and Zero-Downtime Migration

## 1. Zero-Downtime Migration Principles
In production, a schema migration must run while the database actively handles concurrent transactions. This requires treating migrations non-destructively. Any change that locks the table heavily or removes state immediately breaks clients.

### Additive Schema Change
Never modify or drop existing structures that clients still depend on.
- **Good**: `ALTER TABLE users ADD COLUMN phone_new VARCHAR(20) DEFAULT NULL;`
- **Bad**: `ALTER TABLE users RENAME COLUMN phone TO phone_number;` (Instant outage for legacy clients).

## 2. Table Locking and Long Transactions
Operations like adding constraints or changing column types typically rewrite tables, resulting in heavy locking. To prevent downtime:
- **Avoid default constraints on new columns on large tables** (in older database versions, this forces full table rewrite). Use nullable columns, backfill, and then set defaults.
- Create indexes concurrently (e.g., `CREATE INDEX CONCURRENTLY` in Postgres).

## 3. Data Backfill Strategies
Large datasets cannot be migrated in a single UPDATE query (it locks rows, exhausts memory, blocks vacuuming, and kills replication).
- **Batch Processing**: Select/update rows in small chunks (`LIMIT X OFFSET Y` or ID range queries).
- **Throttling**: Add deliberate sleep periods between batch executions to reduce database load.
- **Idempotency**: Backfill jobs should track state or be capable of resuming correctly if interrupted.
- **Background Jobs**: Execute backfills out-of-band via background workers.

## 4. Dual Write and Dual Read Transitions

### Dual Write
- **Concept**: The application synchronously writes data to both the old and new database schemas to keep them in sync during migration.
- **Risks**: Increased latency, transaction complexity, and data race conditions (which system is the source of truth?). Data might drift if one write fails.

### Fallback Read (Dual Read)
- **Concept**: The application attempts to read from the new structure. If the data is empty/null (not yet backfilled), it falls back to reading from the old structure.

### Idempotency in Dual Writes
If application N+1 writes to both Table A (old) and Table B (new), the transaction must be robust. If the migration spans databases or services, distributed transaction risks emerge. Usually, a message queue or Outbox Pattern is preferred over synchronous cross-service dual writes.
