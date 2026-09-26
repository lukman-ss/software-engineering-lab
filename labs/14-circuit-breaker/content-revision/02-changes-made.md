# Content Revision Changes — Lab 14 Circuit Breaker

**Date:** 2026-09-26  
**Based on:** Audit Verdict APPROVED_WITH_WARNINGS  
**Auditor Issues:** AWS URL 301 redirect; 4xx/5xx error classification; metrics export absence

## Files Revised

### content/02-master-draft.md

**Added:**
- Line 143-144: Note on intentional error classification simplification and observability omission
  - "Pedagogical minimalism: lab `Breaker.Execute` counts any `err != nil` with no `IsFailure(error) bool` predicate, and does not export metrics — intentional simplification to keep state machine readable"
- Lines 145-147: Observability subsection clarifying architectural recommendation vs implementation
  - Lists metrics: `circuit_state` gauge, `circuit_open_count`, `rejected_call_count`, `failure_count`, `dependency_latency` P50/P95/P99
  - Notes: Lab exposes only `State() State`; production requires Prometheus/OpenTelemetry

### content/01-content-brief.md

**Updated:**
- Line 58: AWS URL now points to canonical Builder Center URL with 301 note
  - Changed: `AWS Builders Library, *Timeouts, retries, and backoff with jitter* (301 redirect noted)`
  - To: `AWS Builders Library, *Timeouts, retries, and backoff with jitter* — canonical https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter (301 from https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)`

## Changes Summary

| Issue | Fix Applied | Status |
|-------|-------------|--------|
| AWS URL 301 redirect (audit Gap 1) | Updated to canonical Builder Center URL in `01-content-brief.md` | DONE |
| 4xx/5xx error classifier missing (audit Gap 2) | Added pedagogical minimalism note in `02-master-draft.md` explaining no predicate | DONE |
| Observability metrics not exported (audit Gap 3) | Added Observability subsection in `02-master-draft.md` clarifying architectural spec vs implementation | DONE |
| Snippet 7 FakeServer ModeDown mismatch (content-audit Issue #1) | Fixed ModeDown case in `content/03-code-snippets.md` to match `internal/payment/fake_server.go` lines 42-43: replaced `http.Error(w, "internal payment server failure", http.StatusInternalServerError)` with `w.WriteHeader(http.StatusInternalServerError)` and `_, _ = w.Write([]byte("internal payment server failure"))` | DONE |
| Broken source map reference (content-audit Issue #2) | Fixed `content/06-source-map.md` line 198: changed non-existent `research/04-contradictions.md` to `research-audit/04-contradictions.md` | DONE |

## Verification

- All content changes reflect actual implementation (`internal/circuitbreaker/circuit_breaker.go`, `internal/payment/fake_server.go`)
- Both critical issues from content-audit/09-verdict.md now fixed
- Audit verdict was NEEDS_REVISION — both critical issues addressed
- Next: Content reaudit

---

**Revisor:** Technical Content Reviser  
**Status:** READY_FOR_CONTENT_REAUDIT
