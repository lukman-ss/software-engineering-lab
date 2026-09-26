# Content Audit: Issues & Recommendations

## Consolidated Issue List

---

## Issue 1: IMPRECISE TRAFFIC REJECTION CLAIM

**Severity**: WARNING  
**Location**: `content/01-content-brief.md:16`, `content/02-master-draft.md:32`  
**Content**:

- Brief: "Server rejects traffic when unready (HTTP 503) and accepts when ready (HTTP 200)."
- Draft: "aplikasi menolak traffic jika kondisi readiness belum valid."

**Problem**: The `/healthz/ready` endpoint returns 503 when not ready, but the server's `/work` endpoint does **not** enforce readiness — it accepts and processes all requests regardless of readiness state. Traffic rejection is the responsibility of an external load balancer that polls the readiness probe.

**Evidence**: `internal/server/server.go:43-63` — the `/work` handler increments counters and processes the job without checking `s.ready.Load()`.

**Recommendation**: Revise to: "Server signals not-ready via `/healthz/ready` (HTTP 503); external load balancer stops routing traffic to the pod based on probe result."

---

## Issue 2: OVER-SIMPLIFIED WORKER DRAIN DESCRIPTION

**Severity**: WARNING  
**Location**: `content/02-master-draft.md:35`  
**Content**: "Worker sukses mendeteksi sinyal stop, menyelesaikan 1 buah tugas yang sedang berjalan secara utuh, lalu berhenti beroperasi."

**Problem**: States the worker completes exactly "1 job" during graceful shutdown. The worker actually drains all jobs from the channel that can complete within the drain timeout, not just one.

**Evidence**:
- `tests/worker_test.go:28-46` (`TestWorkerGracefulShutdown`): 2 jobs enqueued, both complete.
- `tests/worker_test.go:12-26` (`TestWorkerConcurrency`): 6 jobs completed with concurrency=3.
- `internal/worker/worker.go:47-69`: Worker loop continues processing from `jobChan` until channel is drained or context is canceled.

**Recommendation**: Revise to reflect that the worker drains all buffered jobs within the drain timeout, not just a single job.

---

## Issue 3: INCOMPLETE WORKER DRAIN WARNING IN BRIEF

**Severity**: LOW  
**Location**: `content/01-content-brief.md:20`  
**Content**: "Background worker completes the currently active job upon receiving a stop signal; remaining buffered jobs are abandoned after drain timeout."

**Problem**: Does not clarify that the **active** (in-flight) job is also aborted if it exceeds the drain timeout. The engineering audit explicitly notes this in Finding 4.

**Recommendation**: Add qualifier: "The currently active job completes only if it finishes within the drain timeout; otherwise, it is aborted via context cancellation."

---

## Issue 4: MISSING TEST REFERENCE IN SOURCE MAP

**Severity**: INFO  
**Location**: `content/06-source-map.md` (Health Probes section, line 31)  
**Content**: Source Map references `tests/server_test.go: TestServerProbes, TestServerReadyUnreadyTransition` for health probes.

**Problem**: `TestServerReadyUnreadyTransition` is correctly referenced, but `TestServerInvalidDurationFallback` (which tests the `/work` endpoint fallback duration) is not referenced anywhere in the content.

**Evidence**: `tests/server_test.go:158-185` — `TestServerInvalidDurationFallback` tests the 50ms default fallback when the `d` query parameter is invalid.

**Recommendation**: Add `TestServerInvalidDurationFallback` to the graceful shutdown section of the Source Map, since it tests connection-draining behavior edge cases.

---

## Issue 5: DRAFT LANGUAGE IS INDONESIAN — ENGLISH AUDIENCE CONSIDERATION

**Severity**: INFO  
**Location**: `content/02-master-draft.md`, `content/03-code-snippets.md`  
**Content**: Master Draft and Code Snippets explanations are fully in Indonesian.

**Problem**: The target reader (content brief line 4) includes "Software Engineers, Backend Developers, DevOps Engineers" — an international audience. The content brief and source map are in English; master draft and explanations are in Indonesian.

**Evidence**: `content/01-content-brief.md:4` specifies "Target Reader: Software Engineers, Backend Developers, DevOps Engineers" with no language specification.

**Recommendation**: Clarify the intended target language. If English is intended, translate the master draft and explanation text. If Indonesian is intentional, document it in the content brief.

---

## Summary Table

| # | Issue | Severity | File(s) |
|---|---|---|---|
| 1 | Imprecise "rejects traffic" claim (server doesn't enforce readiness on /work) | WARNING | brief:16, draft:32 |
| 2 | "1 job" over-simplification (worker drains all jobs within timeout) | WARNING | draft:35 |
| 3 | Brief doesn't note active job also abortable on timeout | LOW | brief:20 |
| 4 | `TestServerInvalidDurationFallback` not referenced in Source Map | INFO | source-map |
| 5 | Language inconsistency (English brief + Indonesian draft) | INFO | brief, draft, snippets |

**No errors (ERROR severity)**: No hallucinated facts, no broken references, no incorrect code snippets, no platform-specific bias detected.
