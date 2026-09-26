# 03 — Claim Audit

Target Lab: `labs/20-zero-downtime-deployment`
Research Run: `research/runs/2026-09-26-zero-downtime-deployment/`
Audit Date: 2026-09-26

---

## Claim 1

Claim: NGINX supports graceful configuration reload using `HUP` signal and graceful shutdown using `QUIT` signal.
Location: `03-evidence.md` (Evidence 1 & 2), `05-report.md` (Finding 1 & 2)
Evidence Provided: NGINX control doc describes `HUP` signal master behavior (opening new config, starting new workers, sending graceful shutdown to old workers) and `QUIT` signal behavior.
Source: Source 10 (NGINX Control)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully accurate and supported by official NGINX docs.

---

## Claim 2

Claim: Kubernetes pod termination sends SIGTERM, waits terminationGracePeriodSeconds (default 30s), then SIGKILL; Readiness probes determine pod inclusion in Service endpoints.
Location: `03-evidence.md` (Evidence 3 & 4), `05-report.md` (Finding 1 & 2)
Evidence Provided: Kubernetes Pod Lifecycle and Probes documentation.
Source: Source 8 (Pod Lifecycle) & Source 9 (Probes)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly notes 30s as default configurable threshold.

---

## Claim 3

Claim: Laravel `queue:work` handles SIGTERM for graceful exit; `queue:restart` instructs workers to exit after current job; Horizon `horizon:terminate` requires Supervisor `stopwaitsecs` > longest job duration.
Location: `03-evidence.md` (Evidence 5, 6, 13), `05-report.md` (Finding 2 & 5)
Evidence Provided: Laravel Queues & Horizon documentation.
Source: Source 2 (Queues) & Source 4 (Horizon)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Crucial distinction between immediate SIGKILL and signal-based graceful shutdown correctly captured.

---

## Claim 4

Claim: Blue-Green deployment and Expand-Migrate-Contract (Parallel Change) enable zero-downtime transitions and schema coexistence.
Location: `03-evidence.md` (Evidence 7 & 8), `05-report.md` (Finding 3 & 4)
Evidence Provided: Martin Fowler bliki articles on Blue-Green and Parallel Change.
Source: Source 6 & Source 7
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Theoretical foundations for zero-downtime database and traffic orchestration are accurate.

---

## Claim 5

Claim: PostgreSQL `ADD COLUMN` without volatile default or `NOT NULL` is instant (metadata-only); `ADD CONSTRAINT ... NOT VALID` followed by `VALIDATE CONSTRAINT` prevents write-blocking during migration.
Location: `03-evidence.md` (Evidence 9 & 10), `05-report.md` (Finding 3)
Evidence Provided: PostgreSQL ALTER TABLE documentation.
Source: Source 13 (ALTER TABLE)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Lock mechanics and non-blocking DDl paths accurately reflected for modern PostgreSQL versions.

---

## Claim 6

Claim: Laravel `/up` route returns 200 on boot; `DiagnosingHealth` event must be listened to for verifying DB/Redis readiness.
Location: `03-evidence.md` (Evidence 11), `05-report.md` (Finding 7)
Evidence Provided: Laravel Deployment documentation.
Source: Source 1 (Deployment)
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: LOW
Notes: Sound recommendation distinguishing basic HTTP 200 boot check from functional readiness.

---

## Claim 7

Claim: NGINX `upstream` module `drain` mode and active `health_check` require commercial NGINX Plus (or experimental Lua/custom module in OSS).
Location: `03-evidence.md` (Evidence 12), `04-contradictions.md`
Evidence Provided: NGINX Upstream module documentation.
Source: Source 11 (Upstream Module)
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: Research agent correctly highlights limitation of open-source NGINX rather than falsely claiming OSS supports active probes natively out of the box.

---

## Claim 8

Claim: Redis supports zero-downtime rolling upgrades across version changes.
Location: `03-evidence.md` (Evidence 14), `06-open-questions.md`
Evidence Provided: None (Attempted fetch returned 404).
Source: Source 14 (Failed)
Source Actually Supports Claim: NO
Classification: HYPOTHESIS
Severity: MEDIUM
Notes: Explicitly marked as NOT VERIFIED by research agent. No unsubstantiated claim was asserted as fact.

---

## Summary of Claim Audit
- Major Claims Audited: 8
- Fully Supported: 7
- Explicitly Marked Unverified: 1
- Unsupported Claims Presented as Fact: 0
- Overgeneralized Claims: 0
