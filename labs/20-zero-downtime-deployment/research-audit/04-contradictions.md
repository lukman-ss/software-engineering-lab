# Contradiction Audit: Zero-Downtime Deployment

**Target Lab:** `labs/20-zero-downtime-deployment`  
**Research Run:** `2026-09-26-zero-downtime-deployment`

---

## Contradiction 1: Rolling Update vs Blue-Green Infrastructural Requirements

Statement A:
Blue-Green deployment requires maintaining two identical environments (blue and green) with instant traffic switching via router.  
Location: `research/.../05-report.md: Finding 3`, Martin Fowler  

Statement B:
Rolling deployment operates incrementally in-place within the same environment using `maxSurge` / `maxUnavailable`.  
Location: `research/.../05-report.md: Finding 1`, Kubernetes Docs  

Type: INTERNAL / PATTERN SELECTION  
Impact: LOW. Both achieve zero downtime through different resource/operational tradeoffs.  
Assessment: Properly analyzed and reconciled in `04-contradictions.md`.  

---

## Contradiction 2: Endpoint Deregistration vs Cloud Load Balancer Draining

Statement A:
Kubernetes EndpointSlice immediately sets `ready=false` on pod deletion, preventing new traffic.  
Location: `research/.../03-evidence.md: Evidence 6`, Kubernetes Pod Lifecycle  

Statement B:
External cloud load balancers (AWS ALB/NLB, NGINX) exhibit connection draining latency (up to hundreds of seconds) before completely removing targets.  
Location: `research/.../04-contradictions.md: Contradiction 3` & `06-open-questions.md: #6`  

Type: CODE_DOC_MISMATCH / REALITY GAP  
Impact: MEDIUM. If `terminationGracePeriodSeconds` is shorter than the load balancer drain timeout, client connections drop with 502/504 errors.  
Assessment: Properly documented in contradictions and open questions. Mitigated by recommending preStop hooks and coordinated timeouts.  

---

## Contradiction 3: Fast Metadata Column Addition vs Lock Contention

Statement A:
Adding a column with constant default is metadata-only and fast.  
Location: `research/.../05-report.md: Finding 10`  

Statement B:
PostgreSQL `ALTER TABLE` still requires an `ACCESS EXCLUSIVE` lock briefly, which can queue behind long-running queries and cause connection starvation.  
Location: PostgreSQL lock semantics (noted in `06-open-questions.md: #4`).  

Type: SCOPE_LIMITATION  
Impact: MEDIUM. A "metadata-only" DDL can still cause an outage on high-concurrency tables if lock acquisition blocks queries.  
Assessment: Acknowledged in open questions as an area needing explicit lock timeout safeguards.  

---

## Summary

No fatal or unaddressed contradictions exist between sources or report sections. The research appropriately identified nuanced discrepancies between theoretical zero-downtime and real-world failure modes.
