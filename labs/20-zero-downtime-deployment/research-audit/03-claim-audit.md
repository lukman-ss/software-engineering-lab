# Claim Audit: Zero-Downtime Deployment

**Target Lab:** `labs/20-zero-downtime-deployment`  
**Research Run:** `2026-09-26-zero-downtime-deployment`

---

## Claim 1

Claim: Rolling update gradually replaces old Pods with new ones, keeping the application available throughout the process, controlled by `maxUnavailable` (default 25%) and `maxSurge` (default 25%).  
Location: `05-report.md: Finding 1` & `03-evidence.md: Evidence 1`  
Evidence Provided: Direct quote and configuration details from Kubernetes documentation.  
Source: Source 1 & Source 11 (Kubernetes Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Fully verified against official Kubernetes docs.  

---

## Claim 2

Claim: Liveness probes determine when to restart a container, while readiness probes determine when a container should receive traffic via EndpointSlices.  
Location: `05-report.md: Finding 2` & `03-evidence.md: Evidence 2`  
Evidence Provided: Exact quotes on liveness vs readiness distinction and EndpointSlice removal.  
Source: Source 3 (Kubernetes Probes Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Standard Kubernetes architectural distinction.  

---

## Claim 3

Claim: Blue-green deployment requires two identical environments (blue and green) with instant router switching for rollback.  
Location: `05-report.md: Finding 3` & `03-evidence.md: Evidence 3`  
Evidence Provided: Quoted text from Martin Fowler's canonical article.  
Source: Source 4 (Martin Fowler - Blue Green Deployment)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Classical definition verified.  

---

## Claim 4

Claim: Renaming or dropping columns breaks backward compatibility unless handled via transition phases (e.g. Expand-Contract pattern or views).  
Location: `05-report.md: Finding 4 & Finding 12` & `03-evidence.md: Evidence 4 & Evidence 12`  
Evidence Provided: Quotes from Martin Fowler & ThoughtWorks Evolutionary Database Design.  
Source: Source 4 & Source 9  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Valid pattern for concurrent multi-version application deployment.  

---

## Claim 5

Claim: A container that is merely "running" does not guarantee readiness; health checks must check DB and cache dependencies.  
Location: `05-report.md: Finding 5` & `03-evidence.md: Evidence 5`  
Evidence Provided: Kubernetes probe docs combined with Laravel health endpoint docs.  
Source: Source 3 & Source 6  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Well supported across orchestrator and application layer docs.  

---

## Claim 6

Claim: Pod termination follows: SIGTERM sent → grace period (default 30s) → SIGKILL; endpoint set to ready=false upon deletion initiation.  
Location: `05-report.md: Finding 6` & `03-evidence.md: Evidence 6`  
Evidence Provided: Kubernetes Pod termination lifecycle quotes.  
Source: Source 2 & Source 14  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurately reflects Kubernetes container termination semantics.  

---

## Claim 7

Claim: NGINX active health checks require NGINX Plus, while passive health checks (`max_fails`, `fail_timeout`) are available in NGINX Open Source.  
Location: `05-report.md: Finding 7` & `03-evidence.md: Evidence 7`  
Evidence Provided: NGINX documentation breakdown.  
Source: Source 5 (NGINX Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Crucial operational distinction between OSS and commercial NGINX.  

---

## Claim 8

Claim: Laravel provides built-in health endpoint `/up` returning HTTP 200/500 and supports custom checks via `DiagnosingHealth` event.  
Location: `05-report.md: Finding 8` & `03-evidence.md: Evidence 8`  
Evidence Provided: Laravel 11.x deployment documentation.  
Source: Source 6 (Laravel Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Standard feature in Laravel 11.  

---

## Claim 9

Claim: `php artisan horizon:terminate` gracefully waits for in-flight queue jobs to complete, with Supervisor `stopwaitsecs` set higher than longest running job.  
Location: `05-report.md: Finding 9` & `03-evidence.md: Evidence 9`  
Evidence Provided: Laravel Horizon official documentation.  
Source: Source 8 (Laravel Horizon Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Supported by documentation.  

---

## Claim 10

Claim: PostgreSQL `ALTER TABLE ... ADD COLUMN` with constant default value is a metadata-only operation and does not rewrite the table.  
Location: `05-report.md: Finding 10` & `03-evidence.md: Evidence 10`  
Evidence Provided: PostgreSQL 11+ documentation quote on constant defaults.  
Source: Source 10 (PostgreSQL Documentation)  
Source Actually Supports Claim: PARTIAL  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: MEDIUM  
Notes: True for constant default values, but FALSE for volatile expressions (e.g., `clock_timestamp()`, `gen_random_uuid()`). Report correctly notes this in evidence notes, but executive summary requires explicit warning regarding volatile defaults.  

---

## Claim 11

Claim: Docker container stop sends SIGTERM and waits for grace period (10s Linux default) before sending SIGKILL.  
Location: `05-report.md: Finding 11` & `03-evidence.md: Evidence 11`  
Evidence Provided: Docker CLI reference quotes.  
Source: Source 12 (Docker Documentation)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Verified against Docker docs.  

---

## Claim 12

Claim: Database migrations during zero-downtime deployments should strictly follow Expand-Deploy-Migrate-Contract.  
Location: `05-report.md: Finding 12` & `03-evidence.md: Evidence 12`  
Evidence Provided: Fowler & Evolutionary Database Design patterns.  
Source: Source 4 & Source 9  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: High industry consensus.  
