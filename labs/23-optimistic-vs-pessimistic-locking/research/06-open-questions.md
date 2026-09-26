# Open Questions: Optimistic vs Pessimistic Locking

Research date: 2026-09-26

---

## Unanswered Questions

### OQ-1: Atomic Decrement Recipe — Which Database Docs State It Explicitly?

Status: MEDIUM evidence, not directly quoted from vendor

The `UPDATE ... SET stock = stock - N WHERE stock >= N` + `affected_rows` check pattern is the recommended "third solution" in the topic specification, and its component guarantees (statement atomicity, conditional WHERE) are documented. However, no single Tier 1 source in this research set quotes this exact recipe verbatim as a recommended pattern.

Next step: Search PostgreSQL and MySQL "best practices" or "patterns" docs for the atomic decrement example; check if MySQL docs have a conditional-update example for inventory management; verify if PostgreSQL application-level consistency docs (13.4 Data Consistency Checks at the Application Level) cover this pattern.

---

### OQ-2: Distributed Locking (Redis) When Is It Justified?

Status: LOW evidence

The topic specification asserts that using distributed locks (Redis) for problems solvable within a single database is an anti-pattern. The research confirmed database-native locking as the first option, but did not find a authoritative Tier 1 source stating explicit criteria for when distributed locking is the right choice.

Next step: Research Redis Redlock algorithm discussions (Antirez's original Redlock post + Martin Kleppmann's critique); document the actual boundary conditions for when a single database cannot solve the concurrency problem (multi-database, cross-region, microservices with separate DBs).

---

### OQ-3: Optimistic Locking Version Counter Overflow

Status: NOT VERIFIED

In long-lived systems, a 32-bit integer version counter will eventually overflow. This is a practical concern for optimistic locking deployments. No Tier 1 source in this research discussed overflow behavior or recommended mitigation (64-bit counter, timestamp-based versions, etc.).

Next step: Search Hibernate documentation for timestamp-based versioning recommendation; search production incident reports related to version overflow.

---

### OQ-4: Snapshot Isolation vs SERIALIZABLE — Performance Penalty Quantification

Status: NOT VERIFIED

No quantitative performance benchmarks comparing READ COMMITTED, REPEATABLE READ (snapshot isolation), and SERIALIZABLE throughput under various contention levels were found. Only qualitative claims ("performance reduction") exist in current sources.

Next step: Search for TPC-C benchmarks or academic papers with empirical measurements; PostgreSQL SSI implementation paper (Ports et al. 2012) may contain relevant data.

---

### OQ-5: NOWAIT and SKIP LOCKED Behavior Under Contention

Status: LOW evidence

MySQL supports NOWAIT and SKIP LOCKED options on SELECT ... FOR UPDATE. PostgreSQL also supports both. Behavior when a lock cannot be acquired (NOWAIT immediately errors vs SKIP LOCKED skips locked rows) is documented but the performance characteristics under high contention (flash-sale scenarios like the lab's voucher exercise) were not sourced.

Next step: Search for PostgreSQL SKIP LOCKED benchmark data; compare to MySQL behavior under concurrent queue access patterns.

---

### OQ-6: Idempotency in Optimistic Locking Retry Logic

Status: LOW evidence

The research confirms optimistic locking conflict detection and retry, but does not address idempotency guarantees during retry. If the retrying client has already partially executed side effects (sent email, triggered webhook), the retry may cause duplicate effects. This is outside database-level locking scope but critical for production systems.

Next step: Research idempotency key patterns, exactly-once processing semantics, and how this interacts with optimistic locking retry.

---

## Claims Needing Deeper Research

| Claim | Current Confidence | Upgrade Path |
|-------|-------------------|--------------|
| `SET stock = stock - N WHERE stock >= N` is the recommended atomic pattern | MEDIUM | Verify against PostgreSQL application-level consistency docs (13.4) and MySQL tutorial examples |
| Distributed locks (Redis) are inappropriate when resource is in one database | MEDIUM | Find Kleppmann/Antirez Redlock debate for explicit criteria |
| Versionless optimistic locking (ALL/DIRTY fields in WHERE) works in practice | LOW | Verify Hibernate @OptimisticLock annotation behavior in practice |

---

## Possible Next Research Directions

1. **Voucher flash-sale design (topic spec exercise):** Combine atomic decrement + unique constraint + idempotency key; formalize proof that 1000 concurrent requests cannot exceed quota. Requires research on PostgreSQL's `SERIALIZABLE` vs `INSERT ... ON CONFLICT` (UPSERT) atomicity guarantees.

2. **Cross-database concurrency patterns:** How to implement consistent optimistic locking across PostgreSQL + MySQL + Oracle in a multi-database system; version column naming conventions across ORMs (Laravel Eloquent vs Hibernate vs SQLAlchemy).

3. **Pessimistic lock timeout behavior:** What happens when `SELECT FOR UPDATE` waits indefinitely vs MySQL's `innodb_lock_wait_timeout` vs PostgreSQL's `lock_timeout` configuration. Practical tuning guidance.

4. **MVCC internals impact on locking:** How PostgreSQL's "tuple visibility" mechanism affects row-lock storage (dead tuple creation, VACUUM implications) vs MySQL InnoDB's "lock on index entry" approach — different performance profiles under write-heavy loads.

5. **Optimistic locking in REST API design:** Conflict detection via `ETag` / `If-Match` headers as HTTP-level optimistic locking; how this maps to database version columns. (Relevant for the SaaS/CRM use case in topic spec.)
