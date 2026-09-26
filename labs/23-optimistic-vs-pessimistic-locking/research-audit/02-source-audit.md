# Source Audit: Optimistic vs Pessimistic Locking

## Source 1

Claimed Title: PostgreSQL 18 Documentation - 13.3. Explicit Locking  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/explicit-locking.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. URL verified and content directly checked via WebFetch.

Assessment: PASS

---

## Source 2

Claimed Title: PostgreSQL 18 Documentation - 13.2. Transaction Isolation  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/transaction-iso.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Official documentation for PostgreSQL MVCC and isolation levels.

Assessment: PASS

---

## Source 3

Claimed Title: Concurrency Control  
Claimed Publisher: Wikipedia  
URL: https://en.wikipedia.org/wiki/Concurrency_control  

Reachable: YES  
Source Type: SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Canonical secondary overview citing primary academic literature (Bernstein et al. 1987, Weikum & Vossen 2001, Kung & Robinson 1981).

Assessment: PASS

---

## Source 4

Claimed Title: Optimistic Offline Lock  
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)  
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html  

Reachable: YES  
Source Type: SECONDARY (Canonical Industry Reference)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. URL verified via WebFetch. Accurately describes pre-commit conflict detection and low-conflict assumptions.

Assessment: PASS

---

## Source 5

Claimed Title: Pessimistic Offline Lock  
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)  
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html  

Reachable: YES  
Source Type: SECONDARY (Canonical Industry Reference)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Companion pattern to Optimistic Offline Lock in Fowler catalog.

Assessment: PASS

---

## Source 6

Claimed Title: Optimistic Concurrency Control  
Claimed Publisher: Wikipedia  
URL: https://en.wikipedia.org/wiki/Optimistic_concurrency_control  

Reachable: YES  
Source Type: SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Cites Kung & Robinson 1981 paper.

Assessment: PASS

---

## Source 7

Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.4 Locking Reads  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html)  

Reachable: YES (via official Oracle CDN mirror)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Direct access to dev.mysql.com domain may be blocked/rate-limited by CDN bot protection in automated runners, but mirror on docs.oracle.com is an official primary publication by Oracle.

Assessment: PASS

---

## Source 8

Claimed Title: MySQL 8.0 Reference Manual - 15.7.1 InnoDB Locking  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking.html)  

Reachable: YES (via official Oracle CDN mirror)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Same as Source 7.

Assessment: PASS

---

## Source 9

Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.1 Transaction Isolation Levels  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html)  

Reachable: YES (via official Oracle CDN mirror)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Same as Source 7.

Assessment: PASS

---

## Source 10

Claimed Title: Oracle Database Concepts 19c - 9 Data Concurrency and Consistency  
Claimed Publisher: Oracle  
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Primary source for Oracle row locks and lost update example.

Assessment: PASS

---

## Source 11

Claimed Title: Oracle Database Concepts 19c - 10 Transactions  
Claimed Publisher: Oracle  
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Primary source for ACID transaction semantics and statement atomicity.

Assessment: PASS

---

## Source 12

Claimed Title: Optimistic Locking in JPA  
Claimed Publisher: Baeldung  
URL: https://www.baeldung.com/jpa-optimistic-locking  

Reachable: YES  
Source Type: COMMUNITY / SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Publication date is not fixed/verified, but content provides standard Jakarta/JPA @Version behavior.

Assessment: PASS

---

## Source 13

Claimed Title: DynamoDB Transaction APIs  
Claimed Publisher: Amazon Web Services  
URL: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: PARTIAL  
Supports Claimed Topic: PARTIAL  

Problems:
- DynamoDB is a NoSQL distributed database. While useful for contrasting distributed OCC/ConditionCheck with single RDBMS locks, it is peripheral to core RDBMS locking. Properly identified as supporting context in research.

Assessment: PASS

---

## Source 14

Claimed Title: Hibernate ORM 6.x User Guide - Optimistic Locking  
Claimed Publisher: JBoss / Hibernate team  
URL: https://docs.jboss.org/hibernate/orm/6.3/userguide/html_single/Hibernate_User_Guide.html#locking-optimistic  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Noted in research as cached reference. Accurately reflects standard Hibernate optimistic locking behaviors.

Assessment: PASS
