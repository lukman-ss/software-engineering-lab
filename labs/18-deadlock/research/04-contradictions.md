# Contradictions

No material contradictions discovered on core deadlock concepts across sources.

All authoritative sources (PostgreSQL, MySQL/InnoDB, Microsoft SQL Server, academic literature)
agree on:

1. **Definition of deadlock**: Circular wait where each transaction holds a resource the other needs.
   - Confirmed by Coffman conditions (Source 1), PostgreSQL docs (Source 2), MySQL docs (Sources 3-5),
     SQL Server guide (Source 6), and academic sources (Sources 8-10).

2. **Necessary conditions**: Mutual exclusion, hold-and-wait, no preemption, circular wait.
   - Universal agreement across DB vendors and academic sources.

3. **Detection mechanism**: Automatic detection via wait-for graph or equivalent.
   - PostgreSQL (Source 2), MySQL (Source 4), SQL Server (Source 6) all describe automatic detection.

4. **Resolution**: Abort/rollback one transaction (deadlock victim) to break the cycle.
   - All sources confirm victim selection and rollback.

5. **Primary cause**: Inconsistent lock ordering between concurrent transactions.
   - Universally cited as main cause and solution (Sources 2, 3-5, 6, 7).

6. **Mitigation strategies**: Consistent lock ordering, short transactions, application retry.
   - All sources recommend these approaches (Sources 2, 3-5, 6, 7).

## Implementation Variations (Not Contradictions)

These are differences in how databases implement deadlock handling, not disagreements about
whether deadlock exists or how it fundamentally works:

| Aspect | PostgreSQL | MySQL/InnoDB | SQL Server |
|--------|------------|--------------|------------|
| **Default lock timeout** | No general lock timeout; only deadlock_timeout (1s) for detection cycle | `innodb_lock_wait_timeout` default 50s (statement-level) | `LOCK_TIMEOUT` configurable (default -1 = wait forever) |
| **Victim selection** | Unpredictable; not to be relied upon | Smallest transaction (fewest rows changed) | Least expensive to roll back; configurable priority |
| **Detection tuning** | `deadlock_timeout` (default 1s) | `innodb_deadlock_detect` ON/OFF; wait-for graph limit 200 | Lock monitor interval (5s → 100ms when frequent) |
| **Error code** | SQLSTATE 40P01 (deadlock_detected) | MySQL 1213 (ER_LOCK_DEADLOCK) | Error 1205 (deadlock victim) |
| **Locking reads** | `SELECT FOR UPDATE/SHARE` can cause deadlock | Same; affects gap locking | Same; affects locking hints |

These variations reflect different engineering trade-offs but do not contradict the fundamental
understanding of deadlock as a circular wait condition resolved by aborting one participant.

For the research topic's core questions — why correct transactions deadlock, how databases handle it,
and how to mitigate — there is complete consensus across sources.