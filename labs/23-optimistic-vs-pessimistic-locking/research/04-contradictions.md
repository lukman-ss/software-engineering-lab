# Contradictions: Optimistic vs Pessimistic Locking

## No material contradictions discovered between Tier 1 sources.

All major claims about the mechanisms of optimistic and pessimistic locking, isolation levels, and atomic operations are consistent across PostgreSQL docs, Oracle docs, Wikipedia references, Martin Fowler's patterns, and Microsoft EF Core documentation.

## Minor Observations

### Observation 1: Lost Update Prevention Under Default Isolation

**SOURCE A (PostgreSQL)**:
Under Read Committed (default), an UPDATE in a long-running transaction CAN cause a lost update. PostgreSQL docs describe: "Because Read Committed mode starts each command with a new snapshot that includes all transactions committed up to that instant, subsequent commands in the same transaction will see the effects of the committed concurrent transaction in any case." This means two UPDATE statements in the same transaction may each see a different snapshot, allowing a lost update.

**SOURCE B (Oracle)**:
Under Read Committed (default), Oracle explicitly shows a lost update scenario in Table 10-2 where Transaction 1 updates a row, Transaction 2 waits, then Transaction 2 overwrites. Oracle documentation clearly states: "Devising a strategy to handle lost updates is an important part of application development."

**ASSESSMENT**:
Both agree that default Read Committed isolation level does NOT prevent lost updates. This is consistent. The "lost update" anomaly is NOT prevented by Read Committed in either PostgreSQL or Oracle. This is an important finding that contradicts a common misconception.

### Observation 2: MySQL Behavior (Not Directly Verified)

**SOURCE C (MySQL)**:
MySQL/InnoDB documentation was inaccessible (403) during this research session. The plan mentioned that MySQL REPEATABLE READ behavior under gap locking and 2PL may differ from PostgreSQL's snapshot isolation at the same named isolation level.

**ASSESSMENT**:
This remains unverified from the primary source. The behavioral difference between PostgreSQL REPEATABLE READ and MySQL/InnoDB REPEATABLE READ is known in the industry but needs verification from MySQL official docs. PostgreSQL's REPEATABLE READ is actually Snapshot Isolation; MySQL/InnoDB REPEATABLE READ uses gap locking which behaves differently. This is an area of uncertainty, not a contradiction.

### Observation 3: Serializable Implementation Differences

**SOURCE A (PostgreSQL)**:
PostgreSQL implements Serializable via Serializable Snapshot Isolation (SSI) with predicate locking. It does not use traditional 2PL for serialization.

**SOURCE B (Oracle)**:
Oracle implements Serializable by providing transaction-level read consistency and raising ORA-08177 when a row is modified by another committed transaction after the serializable transaction began.

**ASSESSMENT**:
Different implementation strategies, but both achieve the same guarantee: concurrent serializable transactions produce results equivalent to some serial order. No contradiction — just different mechanisms.
