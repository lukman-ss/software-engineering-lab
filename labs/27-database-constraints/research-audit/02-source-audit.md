# Source Audit: Lab 27 — Database Constraints

**Target Lab:** `labs/27-database-constraints`  
**Audit Date:** 2026-09-26  

---

## Source 1: PostgreSQL Documentation — CREATE TABLE Constraints
- **Claimed Title:** CREATE TABLE Constraints
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/sql-createtable.html#SQL-CREATETABLE-CONSTRAINT-UNIQUE`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. URL anchor points to valid SQL command documentation.
- **Assessment:** PASS

---

## Source 2: PostgreSQL Documentation — Chapter 5.5 Constraints
- **Claimed Title:** Chapter 5.5 Constraints (DDL Constraints)
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/ddl-constraints.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Authoritative source for CHECK, NOT NULL, UNIQUE, PRIMARY KEY, FOREIGN KEY, and EXCLUDE constraints.
- **Assessment:** PASS

---

## Source 3: PostgreSQL Documentation — Chapter 11.6 Unique Indexes
- **Claimed Title:** Chapter 11.6 Unique Indexes
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/indexes-unique.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Authoritative source for NULL handling and auto-index creation.
- **Assessment:** PASS

---

## Source 4: PostgreSQL Documentation — Chapter 11.8 Partial Indexes
- **Claimed Title:** Chapter 11.8 Partial Indexes
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/indexes-partial.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Authoritative source for predicate matching, theorem prover limitations, and parameterized query behavior.
- **Assessment:** PASS

---

## Source 5: PostgreSQL Documentation — Chapter 5.12 Table Partitioning
- **Claimed Title:** Chapter 5.12 Table Partitioning
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/ddl-partitioning.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Covers partition key requirements for UNIQUE and EXCLUDE constraints.
- **Assessment:** PASS

---

## Source 6: PostgreSQL Documentation — Chapter 13.3 Explicit Locking
- **Claimed Title:** Chapter 13.3 Explicit Locking
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/explicit-locking.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Covers row-level locks, `FOR UPDATE`, and write conflicts.
- **Assessment:** PASS

---

## Source 7: PostgreSQL Documentation — Chapter 13.4 Application-Level Consistency
- **Claimed Title:** Chapter 13.4 Application-Level Consistency Checks
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/applevel-consistency.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Covers SSI (Serializable Snapshot Isolation), conflict cycles, and retry strategies.
- **Assessment:** PASS

---

## Source 8: PostgreSQL Documentation — Appendix A Error Codes
- **Claimed Title:** Appendix A Error Codes
- **Claimed Publisher:** PostgreSQL Global Development Group
- **URL:** `https://www.postgresql.org/docs/current/errcodes-appendix.html`
- **Reachable:** YES
- **Source Type:** PRIMARY
- **Relevant:** YES
- **Supports Claimed Topic:** YES
- **Problems:** None. Validates SQLSTATE Class 23 error codes.
- **Assessment:** PASS

---

## Source 9: MySQL 8.0 Documentation (Various URLs)
- **Claimed Title:** MySQL 8.0 Reference Manual (Constraints, Create Table)
- **Claimed Publisher:** Oracle Corporation / MySQL
- **URL:** `https://dev.mysql.com/doc/refman/8.0/en/constraints.html`
- **Reachable:** NO (HTTP 403 Forbidden reported by research agent)
- **Source Type:** PRIMARY (Attempted)
- **Relevant:** YES
- **Supports Claimed Topic:** UNKNOWN
- **Problems:** Blocked / Unreachable during research phase. Research agent accurately declared this in `02-sources.md` and `04-contradictions.md` rather than fabricating MySQL behavior.
- **Assessment:** WARNING (Documented limitation; scope reduced to PostgreSQL)

---

## Source 10: Industry Articles (Citus Data, Martin Fowler)
- **Claimed Titles:** Five Ways to Race Condition (Citus Data), DBC Constraints (Martin Fowler)
- **URLs:** 
  - `https://www.citusdata.com/blog/2018/03/28/five-ways-to-race-condition/`
  - `https://martinfowler.com/articles/dbc-constraints.html`
- **Reachable:** NO (HTTP 404 Not Found)
- **Source Type:** COMMUNITY / SECONDARY
- **Relevant:** YES
- **Supports Claimed Topic:** NO (Dead links)
- **Problems:** Broken links honestly recorded in `02-sources.md` and omitted from evidence synthesis.
- **Assessment:** PASS (Accurately reported as inaccessible; not used as phantom citations)
