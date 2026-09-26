# Research Gap Analysis: Optimistic vs Pessimistic Locking

## Gap 1
Type: WEAK_SOURCE  
Severity: LOW  
Location: `research/06-open-questions.md:9-16` (OQ-1); `research/05-report.md:89-101` (Finding 6)  
Problem: The conditional atomic decrement pattern (`UPDATE ... SET stock = stock - N WHERE stock >= N`) is synthesized from underlying statement-atomicity guarantees rather than directly cited from a single vendor best-practice guide.  
Required Revision: Keep confidence at MEDIUM or verify against PostgreSQL documentation chapter 13.4 ("Data Consistency Checks at the Application Level") during implementation.  
Can Be Approved Without Fix: YES  

---

## Gap 2
Type: SCOPE_ERROR / WEAK_SOURCE  
Severity: LOW  
Location: `research/06-open-questions.md:19-26` (OQ-2)  
Problem: Distributed lock boundaries (e.g., Redis Redlock vs single database ACID locks) are discussed at a high level without citing Martin Kleppmann's formal analysis or primary distributed systems references.  
Required Revision: Bound Redis locking claims as out-of-scope for single-database relational locking lab or add Redlock analysis if distributed locking is introduced.  
Can Be Approved Without Fix: YES  

---

## Gap 3
Type: UNVERIFIED_CLAIM  
Severity: LOW  
Location: `research/06-open-questions.md:29-36` (OQ-3)  
Problem: Integer version column overflow in optimistic locking is flagged as an open question without mitigation guidelines documented in the primary evidence.  
Required Revision: Note standard schema practices (e.g., BIGINT or timestamp versions) during lab design.  
Can Be Approved Without Fix: YES  

---

## Gap 4
Type: UNVERIFIED_CLAIM  
Severity: LOW  
Location: `research/06-open-questions.md:39-46` (OQ-4); `research/05-report.md:166-167`  
Problem: Quantitative benchmark metrics (latency and throughput exact numbers under varying concurrency levels) are absent; only qualitative trade-offs are supported.  
Required Revision: Ensure lab educational text does not present arbitrary benchmark numbers without running live benchmarks during lab execution.  
Can Be Approved Without Fix: YES  

---

## Gap 5
Type: WEAK_SOURCE  
Severity: LOW  
Location: `research/02-sources.md:163-173` (Source 14)  
Problem: Hibernate User Guide for versionless optimistic locking was logged as an unverified direct fetch by the researcher.  
Required Revision: Maintain the author's downgrade to MEDIUM confidence.  
Can Be Approved Without Fix: YES  
