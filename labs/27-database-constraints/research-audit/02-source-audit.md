# Source Audit: Database Constraints

## Source 1
Claimed Title: PostgreSQL 18 - Chapter 5.5: Constraints  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/ddl-constraints.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Section numbers and verbatim citations accurately reflect PostgreSQL documentation on CHECK, NOT NULL, UNIQUE, FOREIGN KEY, and EXCLUSION constraints.

Assessment: PASS

---

## Source 2
Claimed Title: PostgreSQL 18 - Chapter 11.6: Unique Indexes  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/indexes-unique.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Verbatim citations match official documentation regarding B-tree unique indexes and NULL handling (`NULLS NOT DISTINCT`).

Assessment: PASS

---

## Source 3
Claimed Title: PostgreSQL 18 - Chapter 11.8: Partial Indexes  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/indexes-partial.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Verbatim citations match regarding partial indexes, planner implication requirements, and parameterized queries.

Assessment: PASS

---

## Source 4
Claimed Title: PostgreSQL 18 - Appendix A: PostgreSQL Error Codes  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/errcodes-appendix.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. SQLSTATE class 23 codes (`23502`, `23503`, `23505`, `23514`, `23P01`, `23001`) and structured error fields accurately reflected.

Assessment: PASS

---

## Source 5
Claimed Title: PostgreSQL 18 - Chapter 13.4: Data Consistency Checks at the Application Level  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/applevel-consistency.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Correctly cites SERIALIZABLE transaction mechanics and SSI conflict detection for multi-row invariants.

Assessment: PASS

---

## Source 6
Claimed Title: PostgreSQL 18 - Chapter 13.3: Explicit Locking  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/explicit-locking.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- Note: Locking documentation documents table and row-level locks, but exact internal B-tree page lock latching during UNIQUE insertion is engine-internal. Research appropriately noted this as an inference/limitation.

Assessment: PASS

---

## Source 7
Claimed Title: PostgreSQL 18 - Chapter 5.12: Table Partitioning Limitations  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/ddl-partitioning.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Partition key column requirements in partitioned unique constraints match official documentation.

Assessment: PASS

---

## Source 8
Claimed Title: PostgreSQL 18 - ALTER TABLE: NOT VALID and VALIDATE CONSTRAINT  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/sql-altertable.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Lock modes (`SHARE UPDATE EXCLUSIVE` for `VALIDATE CONSTRAINT`) and semantics correctly documented.

Assessment: PASS

---

## Source 9
Claimed Title: SQLite - CREATE TABLE Documentation  
Claimed Publisher: SQLite.org  
URL: https://www.sqlite.org/lang_createtable.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Correctly cites SQLite's independent verification of UNIQUE, NOT NULL, and CHECK constraints.

Assessment: PASS

---

## Source 10
Claimed Title: Stripe API Documentation - Idempotent Requests  
Claimed Publisher: Stripe, Inc.  
URL: https://docs.stripe.com/api/idempotent_requests  

Reachable: YES  
Source Type: PRIMARY (Industry Standard)  
Relevant: YES  
Supports Claimed Topic: PARTIAL  

Problems:
- Stripe's documentation describes application-level idempotency key caching, while database idempotency uses storage-level UNIQUE constraints. Research report explicitly clarifies this distinction (Finding 9 & Limitations).

Assessment: PASS
