# Source Audit: Optimistic vs Pessimistic Locking

## Source 1

Claimed Title: PostgreSQL 18 Documentation - 13.3. Explicit Locking
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative source for PostgreSQL row-level locks (FOR UPDATE, FOR SHARE), deadlock handling, and advisory locks.

Assessment:
PASS

---

## Source 2

Claimed Title: PostgreSQL 18 Documentation - 13.2. Transaction Isolation
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative source for PostgreSQL MVCC isolation levels, Read Committed re-evaluation semantics, and Snapshot Isolation behavior.

Assessment:
PASS

---

## Source 3

Claimed Title: Concurrency Control
Claimed Publisher: Wikipedia (citing Bernstein et al. 1987; Weikum and Vossen 2001)
URL: https://en.wikipedia.org/wiki/Concurrency_control

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Secondary reference; relies on academic citations for formal definitions of lost update and 2PL.

Assessment:
PASS

---

## Source 4

Claimed Title: Optimistic Offline Lock
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Classic architectural reference (David Rice / Martin Fowler, 2003) focusing on business transactions spanning system transactions.

Assessment:
PASS

---

## Source 5

Claimed Title: Pessimistic Offline Lock
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Focuses on application-level offline locks across multiple web requests rather than raw SQL `SELECT FOR UPDATE` within a single ACID transaction.

Assessment:
PASS

---

## Source 6

Claimed Title: Optimistic Concurrency Control
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Optimistic_concurrency_control

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- High-level theoretical summary (Kung & Robinson 1981).

Assessment:
PASS

---

## Source 7

Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.4 Locking Reads
Claimed Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html

Reachable:
PARTIAL (dev.mysql.com returns 403 to automated bots; verified via official Oracle mirror https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Direct dev.mysql.com URL automated fetch blocked by CDN bot detection; contents mirrored verbatim on Oracle documentation CDN.

Assessment:
PASS

---

## Source 8

Claimed Title: MySQL 8.0 Reference Manual - 15.7.1 InnoDB Locking
Claimed Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html

Reachable:
PARTIAL (dev.mysql.com returns 403; verified via Oracle CDN mirror)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Direct dev.mysql.com URL automated fetch blocked by bot protection; content verified via Oracle CDN mirror.

Assessment:
PASS

---

## Source 9

Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.1 Transaction Isolation Levels
Claimed Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html

Reachable:
PARTIAL (dev.mysql.com returns 403; verified via Oracle CDN mirror)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Same mirror dependency as Sources 7 and 8.

Assessment:
PASS

---

## Source 10

Claimed Title: Oracle Database Concepts 19c - 9 Data Concurrency and Consistency
Claimed Publisher: Oracle
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Contains concrete lost-update example (Table 10-2 Banda salary scenario) and TX row lock semantics.

Assessment:
PASS

---

## Source 11

Claimed Title: Oracle Database Concepts 19c - 10 Transactions
Claimed Publisher: Oracle
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative for SQL statement-level atomicity and transaction lifecycle.

Assessment:
PASS

---

## Source 12

Claimed Title: Optimistic Locking in JPA
Claimed Publisher: Baeldung
URL: https://www.baeldung.com/jpa-optimistic-locking

Reachable:
YES

Source Type:
COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Tutorial site; focuses on Hibernate/JPA `@Version` annotation rather than engine-level database primitives.

Assessment:
PASS

---

## Source 13

Claimed Title: DynamoDB Transaction APIs
Claimed Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Scope is distributed NoSQL transactions (`ConditionCheck`, `TransactWriteItems`), somewhat tangential to relational RDBMS locking lab.

Assessment:
PASS

---

## Source 14

Claimed Title: Hibernate ORM 6.x User Guide - Optimistic Locking
Claimed Publisher: JBoss / Hibernate team
URL: https://docs.jboss.org/hibernate/orm/6.3/userguide/html_single/Hibernate_User_Guide.html#locking-optimistic

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- ORM-level implementation details (`StaleObjectStateException`, versionless optimistic locking).

Assessment:
PASS

---

## Source 15

Claimed Title: PostgreSQL 18 Documentation - 13.4. Data Consistency Checks at the Application Level
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/applevel-consistency.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Directly validates requirement for actual `UPDATE` to protect against concurrent modification.

Assessment:
PASS

---

## Source 16

Claimed Title: MySQL 8.0 Reference Manual - 17.7.2.4 Locking Reads (Oracle CDN mirror)
Claimed Publisher: Oracle / MySQL
URL: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Primary mirror of official MySQL 8.0 locking reads documentation.

Assessment:
PASS
