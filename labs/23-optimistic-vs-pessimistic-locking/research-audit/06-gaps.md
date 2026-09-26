# Research Gap Analysis: Optimistic vs Pessimistic Locking

## Gap 1

Type: MISSING_SOURCE  
Severity: LOW  
Location: `research/06-open-questions.md:9-16`, `research/05-report.md:89-100`  
Problem: Exact `UPDATE ... SET stock = stock - N WHERE stock >= N` decrement-with-guard recipe is synthesized from statement atomicity guarantees rather than quoted from an explicit database vendor tutorial/docs snippet.  
Required Revision: Optional future enhancement to locate an explicit PostgreSQL/MySQL docs example or SQL pattern documentation.  
Can Be Approved Without Fix: YES  

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
