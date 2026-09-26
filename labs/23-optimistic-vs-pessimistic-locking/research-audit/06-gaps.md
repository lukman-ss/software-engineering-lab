# Research Gap Analysis: Optimistic vs Pessimistic Locking

## Gap 1

Type: RESOLVED (downgraded from MISSING_SOURCE to VERIFIED)
Severity: CLOSED
Location: `research/06-open-questions.md:9-16`, `research/05-report.md:89-100`
Resolution: 2026-09-26 direct fetch of PostgreSQL 18 Documentation - 13.4 Data Consistency Checks at the Application Level now provides Tier 1 vendor evidence that an actual UPDATE (not merely a lock) is required to prevent concurrent modification: "SELECT FOR UPDATE does not ensure that a concurrent transaction will not update or delete a selected row. To do that in PostgreSQL you must actually update the row, even if no values need to be changed." Combined with ACID statement atomicity (Oracle), conditional WHERE semantics (Oracle WHERE-guard), and single-statement lock scope (MySQL/innodb-locking-reads counter example), the atomic decrement recipe is now HIGH confidence.
Was: MISSING_SOURCE / LOW
Now: VERIFIED / CLOSED  

---

## Gap 2

Type: UNVERIFIED_CLAIM  
Severity: LOW  
Location: `research/06-open-questions.md:19-25`  
Problem: Concrete decision boundaries for when distributed locks (Redis Redlock) are warranted vs single RDBMS locks rely on engineering consensus rather than a formal comparative study.  
Required Revision: Add reference to Kleppmann/Antirez Redlock debate in future iterations.  
Can Be Approved Without Fix: YES  

---

## Gap 3

Type: UNVERIFIED_CLAIM  
Severity: LOW  
Location: `research/06-open-questions.md:29-36`  
Problem: Integer version counter overflow in high-throughput long-running optimistic locking systems is noted as unverified.  
Required Revision: Document BIGINT or timestamp-based version mitigation in production guidelines.  
Can Be Approved Without Fix: YES  

---

## Gap 4

Type: WEAK_SOURCE  
Severity: LOW  
Location: `research/06-open-questions.md:39-46`  
Problem: Performance benchmarks comparing throughput under heavy contention are qualitative ("reduced performance") rather than quantitative (TPC-C / latency graphs).  
Required Revision: Empirical benchmark numbers can be added in performance tuning deep-dives.  
Can Be Approved Without Fix: YES  
