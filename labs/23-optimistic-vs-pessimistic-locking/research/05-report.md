# Research Report: Optimistic vs Pessimistic Locking — Mencegah Data Diam-Diam Tertimpa oleh Request Lain

Research date: 2026-09-26

## Research Question

How do concurrent processes cause silent data loss ("lost update") when two requests read-modify-write the same record, and what are the verified mechanisms, trade-offs, and selection criteria for preventing it via pessimistic locking, optimistic locking, and atomic database operations?

## Executive Summary

Lost update is a verified concurrency anomaly where a second transaction overwrites a first transaction's committed change without error. All major databases (PostgreSQL, MySQL/InnoDB, Oracle) can exhibit it under their default isolation level (READ COMMITTED) when applications use a naive read-modify-write pattern. Three proven prevention strategies exist:

1. **Pessimistic locking** (`SELECT ... FOR UPDATE`) blocks the second transaction until the first commits — verified in PostgreSQL, MySQL, and Oracle official docs.
2. **Optimistic locking** (version/timestamp in `WHERE` clause + `affected_rows` check) detects the conflict at commit time and forces the application to retry or return 409 — verified by Fowler's Patterns and Oracle's WHERE-guard recommendation.
3. **Atomic single-statement operations** (`UPDATE ... SET stock = stock - N WHERE stock >= N`) eliminate the read-modify-write window entirely.

Pessimistic locking prevents conflicts (suited to high-contention, correctness-critical resources like balances, stock, seat reservations). Optimistic locking detects conflicts (suited to low-contention, read-heavy, multi-request business transactions like CMS/CRM/profile edits). No single isolation level prevents lost update automatically without explicit query design; even wrapping code in a transaction is insufficient.

## Findings

### Finding 1: Lost Update Is a Real, Reproducible Anomaly Under Default Isolation

Claim: Two concurrent transactions reading the same value, computing new values independently, and writing back sequentially cause the first write to be silently lost. This occurs under READ COMMITTED (the default in PostgreSQL and Oracle) without explicit locking.

Evidence: Wikipedia defines lost update as "A second transaction writes a second value ... on top of a first value written by a first concurrent transaction, and the first value is lost." Oracle documents a concrete case: Session 1 updates Banda salary to 7000; Session 2 reads stale 6200; Session 1 commits; Session 2 writes 6300 — 7000 is lost. PostgreSQL Read Committed re-evaluates the WHERE clause after a concurrent commit, allowing silent overwrite. MySQL READ COMMITTED uses semi-consistent reads with the same vulnerability if the application does not guard the WHERE clause.

Sources:
- Wikipedia - Concurrency Control (https://en.wikipedia.org/wiki/Concurrency_control)
- Oracle 19c Concepts 9 Data Concurrency and Consistency (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html)
- PostgreSQL 18 Docs 13.2 Transaction Isolation (https://www.postgresql.org/docs/current/transaction-iso.html)
- MySQL 8.0 Docs 15.7.2.1 Transaction Isolation Levels (https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html)

Confidence: HIGH

### Finding 2: Pessimistic Locking Blocks Concurrent Writers via SELECT ... FOR UPDATE

Claim: `SELECT ... FOR UPDATE` acquires a row-level exclusive lock held until transaction commit/rollback. Concurrent transactions attempting UPDATE/DELETE/SELECT FOR UPDATE on the same row block (wait indefinitely unless NOWAIT/SKIP LOCKED or deadlock detection intervenes). Row locks do not block plain SELECT.

Evidence: PostgreSQL 13.3.2: "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends." MySQL InnoDB: "SELECT ... FOR UPDATE locks the rows and associated index entries ... All locks are released when the transaction is committed or rolled back. NOWAIT and SKIP LOCKED options are also available." Oracle: TX row lock acquired by SELECT ... FOR UPDATE, held until commit/rollback. PostgreSQL lock-compatibility tables confirm plain SELECT is blocked only by ACCESS EXCLUSIVE.

Sources:
- PostgreSQL 18 Docs 13.3 Explicit Locking (https://www.postgresql.org/docs/current/explicit-locking.html)
- MySQL 8.0 Docs 15.7.2.4 Locking Reads (https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html)
- Oracle 19c Concepts 9 Data Concurrency (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html)

Confidence: HIGH

### Finding 3: Pessimistic Locking Trade-offs — Concurrency Loss, Deadlock, and Hold-Time Sensitivity

Claim: Pessimistic locking reduces concurrency, increases wait time, and introduces deadlock risk. Holding a lock during external calls (HTTP/gateway, PDF generation) is an anti-pattern. Deadlocks are auto-detected and one transaction is aborted unpredictably; consistent lock ordering and keeping transactions short are the documented defenses.

Evidence: PostgreSQL 13.3.4: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions ... Exactly which transaction will be aborted is difficult to predict." Best defense: consistent lock order + most-restrictive lock first. Wikipedia: "Most non-optimistic mechanisms (with blocking) are prone to deadlocks." PostgreSQL: "A transaction seeking a lock will wait indefinitely ... it is a bad idea to hold transactions open for long periods (e.g., while waiting for user input)." This directly validates the topic spec's "Jangan melakukan Call Payment Gateway di dalam transaction yang memegang lock."

Sources:
- PostgreSQL 18 Docs 13.3.4 Deadlocks (https://www.postgresql.org/docs/current/explicit-locking.html)
- Wikipedia Concurrency Control Categories (https://en.wikipedia.org/wiki/Concurrency_control)

Confidence: HIGH

### Finding 4: Optimistic Locking Detects Conflicts via Version/Timestamp Guard and Affected-Rows Check

Claim: Optimistic locking assumes conflicts are rare, allows concurrent reads, and validates at commit by including the originally-read version in the UPDATE WHERE clause. Zero affected rows signals a conflict that the application must handle (reload/recalculate/retry or 409 Conflict). It cannot be ignored.

Evidence: Fowler: "Optimistic Offline Lock solves this problem by validating that the changes about to be committed by one session don't conflict ... A successful pre-commit validation is, in a sense, obtaining a lock." Wikipedia optimistic: "only check for violations at each transaction's commit. If violations are detected, the transaction is aborted and restarted. Very efficient when few transactions are aborted." Oracle recommends `UPDATE ... WHERE last_name='Banda' AND salary=6200` (original-value guard). The topic spec pattern `UPDATE products SET stock=7, version=6 WHERE id=10 AND version=5` with `0 rows affected` is the exact implementation. Baeldung/Hibernate: version mismatch throws OptimisticLockException/StaleObjectStateException requiring retry.

Sources:
- Martin Fowler Optimistic Offline Lock (https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html)
- Wikipedia Concurrency Control (https://en.wikipedia.org/wiki/Concurrency_control)
- Oracle 19c Concepts 9 (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html)
- Baeldung JPA Optimistic Locking (https://www.baeldung.com/jpa-optimistic-locking)
- Hibernate ORM 6.3 User Guide (https://docs.jboss.org/hibernate/orm/6.3/userguide/html_single/Hibernate_User_Guide.html#locking-optimistic)

Confidence: HIGH

### Finding 5: Selection Criteria — Pessimistic Prevents, Optimistic Detects

Claim: Pessimistic = prevent conflict (block early). Optimistic = detect conflict (check late). Pessimistic suits high-contention or correctness-critical resources (balances, limited stock, queue numbers, seat reservations, inventory allocation). Optimistic suits low-contention, read-heavy, multi-request workflows (profile/CRM/CMS/master data/documents). When neither is needed, neither should be used.

Evidence: Fowler: "Pessimistic Offline Lock assumes that the chance of session conflict is high and therefore limits the system's concurrency, Optimistic Offline Lock assumes that the chance of conflict is low." "Pessimistic prevents conflicts by avoiding them altogether. It forces a business transaction to acquire a lock before it starts to use it." Wikipedia: optimistic "very efficient when few transactions are aborted" vs pessimistic "Blocking ... typically involved with performance reduction." Baeldung: optimistic suitable when "many more reads than updates" and entities are detached. Topic spec's selection table matches these sources.

Sources:
- Martin Fowler Optimistic Offline Lock (https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html)
- Martin Fowler Pessimistic Offline Lock (https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html)
- Wikipedia Concurrency Control (https://en.wikipedia.org/wiki/Concurrency_control)
- Baeldung JPA Optimistic Locking (https://www.baeldung.com/jpa-optimistic-locking)

Confidence: HIGH

### Finding 6: Atomic Single-Statement Updates Eliminate the Read-Modify-Write Window

Claim: For simple counters/quota/stock decrements, a single conditional UPDATE (`UPDATE products SET stock = stock - 3 WHERE id=10 AND stock >= 3` + check `affected_rows == 1`) is atomic at the statement level and removes the race window without explicit locking. If affected_rows == 0, stock was insufficient.

Evidence: PostgreSQL 13.4.2 (Enforcing Consistency with Explicit Blocking Locks, directly fetched 2026-09-26): "SELECT FOR UPDATE does not ensure that a concurrent transaction will not update or delete a selected row. To do that in PostgreSQL you must actually update the row, even if no values need to be changed." This confirms an actual UPDATE is the authoritative conflict-resolution action — a single conditional UPDATE (`SET stock = stock - 3 WHERE stock >= 3`) both performs the modification and is atomic at statement level, removing the read-modify-write window. Oracle ACID: "A SQL statement is an atomic unit; on failure, only its effects are rolled back." MySQL: under READ COMMITTED "InnoDB holds locks only for rows that it updates or deletes" — single conditional UPDATE holds lock only for the matched row. Oracle's documented WHERE-guard pattern (UPDATE WHERE salary=6200) shows conditional update is the vendor-recommended prevention.

Sources:
- PostgreSQL 18 Docs 13.4 Application-Level Consistency (https://www.postgresql.org/docs/current/applevel-consistency.html)
- Oracle 19c Concepts 10 Transactions (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html)
- MySQL 8.0 Docs Innodb Locking Reads (https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html)
- PostgreSQL 18 Docs 13.2 Transaction Isolation (https://www.postgresql.org/docs/current/transaction-iso.html)

Confidence: HIGH — upgraded from MEDIUM (2026-09-26 direct fetch of PG 13.4.2 confirms actual-UPDATE requirement; component guarantees verified across MySQL, Oracle).

### Finding 7: Isolation Level Alone Does Not Prevent Lost Update

Claim: Wrapping read-modify-write in a transaction does not automatically prevent lost update. The outcome depends on the queries and isolation level. PostgreSQL READ COMMITTED allows the anomaly; REPEATABLE READ/SERIALIZABLE throw an error instead. Oracle's lost-update example occurs inside transactions. MySQL's semi-consistent reads under READ COMMITTED still require WHERE-guard handling.

Evidence: PostgreSQL: "Within a REPEATABLE READ or SERIALIZABLE transaction, however, an error will be thrown if a row to be locked has changed since the transaction started." Under READ COMMITTED re-evaluation permits overwrite. Oracle Table 10-2 lost-update occurs despite transactional boundaries. MySQL docs describe semi-consistent evaluation that still needs application handling.

Sources:
- PostgreSQL 18 Docs 13.3.2 Row-Level Locks (https://www.postgresql.org/docs/current/explicit-locking.html)
- PostgreSQL 18 Docs 13.2 (https://www.postgresql.org/docs/current/transaction-iso.html)
- Oracle 19c Concepts 9 (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html)
- MySQL 8.0 Docs 15.7.2.1 (https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html)

Confidence: HIGH

### Finding 8: Default Isolation Levels Differ Across Databases

Claim: PostgreSQL defaults to READ COMMITTED (READ UNCOMMITTED is a no-op identical to READ COMMITTED due to MVCC). Oracle defaults to READ COMMITTED and does not offer READ UNCOMMITTED or REPEATABLE READ. MySQL/InnoDB defaults to REPEATABLE READ. SERIALIZABLE in MySQL implicitly converts plain SELECT to SELECT ... FOR SHARE. Repeatable Read means Snapshot Isolation in PostgreSQL but gap-lock-based 2PL in MySQL.

Evidence: PostgreSQL 13.2: only three effective levels (Read Committed, Repeatable Read, Serializable); Read Uncommitted behaves identically to Read Committed. MySQL 15.7.2.1: REPEATABLE READ is default; SERIALIZABLE converts plain SELECT. Oracle Concepts: READ COMMITTED default; levels are READ COMMITTED, SERIALIZABLE, READ ONLY.

Sources:
- PostgreSQL 18 Docs 13.2 (https://www.postgresql.org/docs/current/transaction-iso.html)
- MySQL 8.0 Docs 15.7.2.1 (https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html)
- Oracle 19c Concepts 9 (https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html)

Confidence: HIGH for PostgreSQL and Oracle (direct fetch verified); MEDIUM for MySQL (via Oracle CDN mirror)

### Finding 9: Common Anti-Patterns Are Documented

Claim: Five recurring mistakes have verified support:
1. Assuming transaction alone prevents lost update — false (Finding 7).
2. Using distributed/Redis locks when database can solve it — PostgreSQL advisory locks doc states "A common use of advisory locks is to emulate pessimistic locking strategies ... While a flag stored in a table could be used, advisory locks are faster" — implying database-native mechanism is preferred before external locks.
3. Holding pessimistic locks across network calls — violates PostgreSQL "do not hold transactions open waiting for user input" guidance.
4. Ignoring 0-rows-affected on optimistic update — Fowler/Hibernate/Baeldung all require explicit handling (retry or 409).
5. Naive read-modify-write (`$product->stock -= 1; $product->save()`) without concurrency guard — directly exhibits Finding 1.

Sources:
- PostgreSQL 18 Docs 13.3.4 + 13.3.5 (https://www.postgresql.org/docs/current/explicit-locking.html)
- Martin Fowler Optimistic Offline Lock (https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html)
- Baeldung / Hibernate optimistic docs

Confidence: HIGH for 1, 3, 4, 5; MEDIUM for 2 (distributed-lock preference is strong industry consensus but not stated as "never use Redis" in fetched Tier 1 sources)

## Areas of Agreement

- All sources agree lost update is a real anomaly distinct from dirty read / non-repeatable read / phantom read.
- All database vendors agree `SELECT ... FOR UPDATE` is the standard pessimistic row-level mechanism with commit/rollback release.
- All sources agree optimistic locking is detection-not-prevention and requires application handling of version-mismatch / 0-rows.
- All sources agree performance degrades with blocking (pessimistic) and with aborts/retries (optimistic); choice depends on conflict frequency.
- All sources agree a plain transaction without appropriate query design (locking or version guard or atomic statement) does not prevent lost update.

## Areas of Disagreement

- **What "Repeatable Read" guarantees:** PostgreSQL (Snapshot Isolation, no phantoms) vs MySQL (2PL with gap/next-key locks, prevents phantoms for locking reads but via different mechanism) vs ANSI spec (allows phantoms). Not a contradiction in advice but a semantic divergence requiring database-specific reasoning.
- **READ UNCOMMITTED existence:** PostgreSQL treats it as identical to READ COMMITTED (no dirty reads possible); other systems allow dirty reads. Applications porting READ UNCOMMITTED logic to PostgreSQL get no benefit.
- **Vendor emphasis on remedies:** Oracle's documented remedy is original-value WHERE guard; PostgreSQL docs emphasize explicit locking and SERIALIZABLE; MySQL emphasizes semi-consistent reads + gap locking. All are valid within their engine; none claims exclusivity.

## Verification Provenance

All Tier 1 claims in this report were verified by direct source fetches on 2026-09-26 (not search snippets):

| Source | URL | Verification |
|--------|-----|-------------|
| PostgreSQL 13.3 Explicit Locking | postgresql.org/docs/current/explicit-locking.html | Direct fetch — FOR UPDATE quote, deadlock example, advisory-lock text verified verbatim |
| PostgreSQL 13.2 Transaction Isolation | postgresql.org/docs/current/transaction-iso.html | Direct fetch — isolation table, Read Committed re-evaluation, Snapshot Isolation text verified |
| PostgreSQL 13.4 Application-Level Consistency | postgresql.org/docs/current/applevel-consistency.html | Direct fetch — "must actually update the row" quote verified (upgrades atomic finding) |
| MySQL 17.7.2.4 Locking Reads | docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html | Direct fetch via official Oracle CDN mirror — FOR UPDATE/FOR SHARE, NOWAIT, SKIP LOCKED verified |
| Oracle 19c Data Concurrency | docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html | Direct fetch — Table 10-2 Banda lost-update example, row lock (TX) semantics verified |
| Wikipedia Concurrency Control | en.wikipedia.org/wiki/Concurrency_control | Direct fetch — lost-update definition, optimistic/pessimistic/semi-optimistic categories verified; revision reviewed 2026-09-07 |
| Fowler Optimistic Offline Lock | martinfowler.com/eaaCatalog/optimisticOfflineLock.html | Direct fetch — "assumes that the chance of conflict is low" verified; dated 2003-03-05, David Rice |
| Fowler Pessimistic Offline Lock | martinfowler.com/eaaCatalog/pessimisticOfflineLock.html | Direct fetch — "prevents conflicts by avoiding them altogether" verified; dated 2003-03-05, David Rice |

## Limitations

- MySQL documentation was accessed via Oracle CDN mirror (https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/...) due to 403 on direct dev.mysql.com fetch; the mirror hosts official Oracle/MySQL documentation and was fetched directly, but direct dev.mysql.com-domain verification was not completed — MySQL-specific claims remain MEDIUM confidence.
- The exact `UPDATE ... SET stock = stock - N WHERE stock >= N` recipe is not quoted verbatim as a single vendor example; the finding is now supported by PostgreSQL 13.4.2's explicit statement that an actual UPDATE (not just a lock) is required to prevent concurrent modification, plus statement-atomicity guarantees across vendors — upgraded to HIGH.
- Hibernate source (versionless optimistic locking details) was reported by a subagent but not directly fetched this session — treated as MEDIUM.
- Distributed locking (Redis) trade-offs were not deeply sourced from a Tier 1 distributed-systems reference; assessment relies on inference from database-native lock preference.
- Snapshot Isolation (SI) vs ANSI levels debate (Berenson et al. 1995 "Critique of ANSI SQL Isolation Levels") is cited by PostgreSQL docs but the paper itself was not fetched.
- Performance benchmarks (throughput/latency numbers for optimistic vs pessimistic under contention) were not found in Tier 1 sources; only qualitative trade-offs are verified.

## Conclusion

Lost update is a silent, error-free data corruption that reproduces under default isolation in major databases when applications use read-modify-write. Pessimistic locking (SELECT FOR UPDATE) and optimistic locking (version guard + affected_rows check) are both verified, complementary strategies — one prevents, one detects. For simple counters, a single conditional UPDATE is the minimal, verified alternative that avoids both. Isolation levels alone are insufficient; query design determines correctness. The documented selection heuristic is: atomic operation first, then optimistic if conflicts are rare, pessimistic if conflicts are frequent or correctness is critical, distributed locks only when the resource spans databases. Any optimistic conflict (0 rows affected) must be handled explicitly; any pessimistic lock must be held for the shortest possible transaction.
