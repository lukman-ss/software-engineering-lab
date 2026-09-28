# Research Gap Analysis

## Gap 1: Arbitrary Sticky Routing Window

**Type**: WEAK_SOURCE  
**Severity**: MEDIUM  
**Location**: `05-report.md` (Conclusion), `06-open-questions.md` (Question 1)  
**Problem**: The recommendation of a "5-second sticky routing window" after writes is based on conventional tutorials and the lab specification rather than empirical benchmark publications or vendor guidance.  
**Required Revision**: In the lab implementation and documentation, treat the 5-second window explicitly as a configurable default threshold rather than a distributed systems guarantee.  
**Can Be Approved Without Fix**: YES (The research report itself already discloses this limitation in `04-contradictions.md` and `06-open-questions.md`).

---

## Gap 2: Operational Overhead of LSN Polling

**Type**: MISSING_CASE  
**Severity**: LOW  
**Location**: `06-open-questions.md` (Question 3)  
**Problem**: Polling `pg_last_wal_receive_lsn` or querying replica status before read queries introduces connection round-trip latency that may negate read-replica scaling benefits if not implemented via cached heartbeats or sticky connection pools.  
**Required Revision**: The subsequent implementation design should note the trade-off between strict LSN verification and latency overhead.  
**Can Be Approved Without Fix**: YES (Adequately flagged in open questions).

---

## Gap 3: Multi-Region Replication Lag Metrics

**Type**: WEAK_SOURCE  
**Severity**: LOW  
**Location**: `06-open-questions.md` (Question 2)  
**Problem**: Quantitative p99 latency distributions for cross-region replication under peak load are not available in public vendor documentation.  
**Required Revision**: None for baseline research; acknowledged as out of scope for a local educational lab.  
**Can Be Approved Without Fix**: YES.
