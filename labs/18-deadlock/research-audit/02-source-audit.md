# Source Audit

## Source 1

Claimed Title: PostgreSQL 18 Documentation — Chapter 13: Concurrency Control, §13.3.4 Deadlocks
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference on deadlock mechanisms, detection, and mitigation in PostgreSQL.

Assessment:
PASS

---

## Source 2

Claimed Title: PostgreSQL 18 Documentation — Chapter 19: Server Configuration, §19.12 Lock Management (deadlock_timeout)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-locks.html#GUC-DEADLOCK-TIMEOUT

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Confirms default 1s parameter and performance trade-offs.

Assessment:
PASS

---

## Source 3

Claimed Title: MySQL 8.4 Reference Manual — §17.7.5.1 An InnoDB Deadlock Example
Claimed Publisher: Oracle (via Wayback Machine)
URL: https://web.archive.org/web/20241214192732/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-example.html

Reachable:
YES

Source Type:
PRIMARY (Archived)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Live dev.mysql.com endpoint returns HTTP 403 to automated clients; properly documented as Wayback Machine snapshot.

Assessment:
PASS

---

## Source 4

Claimed Title: MySQL 8.4 Reference Manual — §17.7.5.2 Deadlock Detection
Claimed Publisher: Oracle (via Wayback Machine)
URL: https://web.archive.org/web/20250126061946/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-detection.html

Reachable:
YES

Source Type:
PRIMARY (Archived)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Archived snapshot used due to upstream access protection; accurately reflects InnoDB wait-for graph and victim selection rules.

Assessment:
PASS

---

## Source 5

Claimed Title: MySQL 8.4 Reference Manual — §17.7.5.3 How to Minimize and Handle Deadlocks
Claimed Publisher: Oracle (via Wayback Machine)
URL: https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html

Reachable:
YES

Source Type:
PRIMARY (Archived)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Confirms ordering, transaction length, commit frequency, and application retry policies.

Assessment:
PASS

---

## Source 6

Claimed Title: Microsoft Learn — Deadlocks Guide (SQL Server / Azure SQL Database)
Claimed Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Extremely comprehensive, verified live.

Assessment:
PASS

---

## Source 7

Claimed Title: How to deal with MySQL deadlocks
Claimed Publisher: Percona Blog (Peiran Song)
URL: https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Article dates from 2014, but foundational InnoDB deadlock analysis principles remain relevant.

Assessment:
PASS

---

## Source 8

Claimed Title: Coffman, E.G.; Elphick, M.J.; Shoshani, A. — "System Deadlocks"
Claimed Publisher: ACM Computing Surveys
URL: https://doi.org/10.1145/356586.356588

Reachable:
YES

Source Type:
PRIMARY (Academic)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Research acknowledged accessing this primarily via secondary citations rather than reading the original ACM publication text directly.

Assessment:
WARNING

---

## Source 9

Claimed Title: Wikipedia — Deadlock (computer science)
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Deadlock_(computer_science)

Reachable:
YES

Source Type:
COMMUNITY / TERTIARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Tertiary crowdsourced reference. While accurate regarding Coffman conditions, it should not be treated as a primary technical source for an authoritative engineering lab.

Assessment:
WARNING

---

## Source 10

Claimed Title: Wikipedia — Two-phase locking
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Two-phase_locking

Reachable:
YES

Source Type:
COMMUNITY / TERTIARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Tertiary source summarizing Bernstein et al. and Weikum & Vossen without direct primary paper verification.

Assessment:
WARNING

---

## Source 11

Claimed Title: AWS Prescriptive Guidance — Retry and Backoff Pattern
Claimed Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/retry-backoff.html

Reachable:
YES

Source Type:
SECONDARY (Vendor Architecture Guidance)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Generic cloud pattern, not database-deadlock specific, but directly supports the idempotent retry and exponential backoff claim.

Assessment:
PASS

---

## Source 12

Claimed Title: Oracle Database Concepts 19c — Data Concurrency and Consistency
Claimed Publisher: Oracle
URL: NOT VERIFIED

Reachable:
NO

Source Type:
UNKNOWN

Relevant:
YES

Supports Claimed Topic:
NO

Problems:
- Research team honestly excluded this source after link failure / confusion with MySQL docs.

Assessment:
FAIL
