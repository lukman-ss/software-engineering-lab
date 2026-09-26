# 06 — Research Gap Analysis

Target Lab: `labs/20-zero-downtime-deployment`
Research Run: `research/runs/2026-09-26-zero-downtime-deployment/`
Audit Date: 2026-09-26

---

## Gap 1

Type: `WEAK_SOURCE`
Severity: `LOW`
Location: `02-sources.md` (Source 14), `03-evidence.md` (Evidence 14), `06-open-questions.md`
Problem: Redis rolling upgrade doc fetch returned HTTP 404. Official primary documentation for Redis cluster rolling upgrade behavior was not captured.
Required Revision: Future research iterations can update the Redis URL to valid documentation if Redis cluster rolling upgrade is explicitly included in the lab requirements.
Can Be Approved Without Fix: YES (The research agent correctly isolated this claim as `NOT VERIFIED` and did not use it as a foundation for core findings).

---

## Gap 2

Type: `SCOPE_ERROR`
Severity: `LOW`
Location: `05-report.md` (Limitations section)
Problem: NGINX Open Source active health-check limitation means pure OSS implementations require either NGINX reload, OpenResty/Lua script, or an external load balancer (e.g. AWS ALB / HAProxy) for dynamic endpoint removal without reloading NGINX.
Required Revision: The engineering design stage must explicitly select either Kubernetes Service readiness gating or NGINX HUP reload for the demo.
Can Be Approved Without Fix: YES (The limitation is transparently documented in the research report).

---

## Gap 3

Type: `MISSING_CASE`
Severity: `LOW`
Location: `05-report.md` (Finding 5 & Limitations)
Problem: Application-specific timeout values (`stopwaitsecs`, `terminationGracePeriodSeconds`, HTTP request timeout) vary based on workload (e.g., fast API vs long-running background export).
Required Revision: Engineering lab implementation should define explicit, standard timeouts for demonstration purposes (e.g., 5s grace for HTTP, 10s for workers).
Can Be Approved Without Fix: YES (Acknowledged as application-dependent in research findings).

---

## Summary of Gaps

- Total Gaps Identified: 3
- Critical / High Severity Gaps: 0
- Medium / Low Severity Gaps: 3
- Can Be Approved Without Fix: YES (All gaps are transparently recorded with appropriate severity and non-blocking status).
