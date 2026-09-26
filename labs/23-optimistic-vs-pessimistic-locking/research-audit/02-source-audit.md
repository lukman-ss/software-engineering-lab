# Source Audit: Optimistic vs Pessimistic Locking

## Source 1
Claimed Title: PostgreSQL 18 Documentation - 13.3. Explicit Locking  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/explicit-locking.html  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Primary vendor documentation for `FOR UPDATE`, row-level locking modes, compatibility matrix, and deadlock handling.  
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
Problems: None. Primary vendor documentation for transaction isolation levels, READ COMMITTED lost update re-evaluation, and snapshot isolation.  
Assessment: PASS  

---

## Source 3
Claimed Title: Concurrency Control  
Claimed Publisher: Wikipedia (citing Bernstein et al. 1987, Weikum & Vossen 2001)  
URL: https://en.wikipedia.org/wiki/Concurrency_control  
Reachable: YES  
Source Type: SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Secondary community synthesis, but cites authoritative academic texts for foundational concurrency definitions.  
Assessment: PASS  

---

## Source 4
Claimed Title: Optimistic Offline Lock  
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture, David Rice)  
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html  
Reachable: YES  
Source Type: SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Canonical reference for application-level optimistic offline lock pattern.  
Assessment: PASS  

---

## Source 5
Claimed Title: Pessimistic Offline Lock  
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture, David Rice)  
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html  
Reachable: YES  
Source Type: SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Canonical reference for application-level pessimistic locking and trade-offs.  
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
Problems: General encyclopedia entry summarizing Kung & Robinson (1981).  
Assessment: PASS  

---

## Source 7
Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.4 Locking Reads  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html)  
Reachable: YES (direct domain may return 403 on automated scrapers; reachable via Oracle CDN mirror)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Direct URL occasionally blocks non-browser agents, but official CDN mirror provides identical manual content.  
Assessment: PASS  

---

## Source 8
Claimed Title: MySQL 8.0 Reference Manual - 15.7.1 InnoDB Locking  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking.html)  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Same mirror consideration as Source 7. Content documents shared/exclusive and gap locks accurately.  
Assessment: PASS  

---

## Source 9
Claimed Title: MySQL 8.0 Reference Manual - 15.7.2.1 Transaction Isolation Levels  
Claimed Publisher: Oracle / MySQL  
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html (Mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html)  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Documents InnoDB isolation defaults, semi-consistent reads, and gap locking.  
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
Problems: None. Directly details lost update scenarios (Table 10-2) and WHERE-clause guards.  
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
Problems: None. Documents statement-level atomicity and transaction commit properties.  
Assessment: PASS  

---

## Source 12
Claimed Title: Optimistic Locking in JPA  
Claimed Publisher: Baeldung  
URL: https://www.baeldung.com/jpa-optimistic-locking  
Reachable: YES  
Source Type: COMMUNITY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Tutorial site; publication date not fixed, but accurately portrays `@Version` and `OptimisticLockException` conventions.  
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
Problems: Relevant for comparing distributed NoSQL transactions, but peripheral to core relational locking strategies. Appropriately bounded as secondary comparison in research.  
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
Problems: The research noted direct automated fetch was not re-verified during session. Explicitly downgraded to MEDIUM confidence in evidence logs.  
Assessment: WARNING  
