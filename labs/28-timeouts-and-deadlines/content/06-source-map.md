# Source Map

Pemetaan setiap bagian master draft dan content snippets ke sumber riset, kode implementasi, dan test yang memverifikasinya.

---

## Deadline Propagation & ExecuteWithBudget

Research:
- `research/01-plan.md` — RQ6: Context & Deadline Propagation
- `research/03-evidence.md` (Evidence 6) — gRPC deadlines, remaining-duration vs absolute timestamp
- `research/05-report.md` — Finding 6: Context and Deadline Propagation

Implementation:
- `internal/deadline/deadline.go` — `ExecuteWithBudget` (lines 13-27)
- `internal/deadline/deadline.go` — `WorkerFunc`, `ErrDeadlineExceeded`

Tests:
- `internal/deadline/deadline_test.go` — `TestExecuteWithBudget_Success`, `TestExecuteWithBudget_Timeout`, `TestExecuteWithBudget_ParentTimeoutInherited`

---

## Exponential Backoff & Full Jitter

Research:
- `research/01-plan.md` — RQ4: Retry Storms, Exponential Backoff, and Jitter
- `research/03-evidence.md` (Evidence 4) — retry storm math, Full Jitter formula
- `research/04-contradictions.md` — Contradiction 3: Retrying on Timeout (Retry Budgets)
- `research/05-report.md` — Finding 4: Retry Storms and Jittered Exponential Backoff

Implementation:
- `internal/retry/retry.go` — `Retrier`, `Config`, `CalculateBackoff` (lines 22-48), `Do` (lines 50-77)

Tests:
- `internal/retry/retry_test.go` — `TestRetrier_SuccessOnFirstTry`, `TestRetrier_RetryUntilSuccess`, `TestRetrier_ExceedMaxAttempts`, `TestRetrier_ContextCanceled`, `TestRetrier_JitterBoundsAndZeroConfig`

---

## Circuit Breaker (3-State FSM)

Research:
- `research/01-plan.md` — (mentioned) circuit breaking
- `research/03-evidence.md` (Evidence 1) — cascading failure, worker pool exhaustion
- `research/06-open-questions.md` — (retry/circuit as mitigation)

Implementation:
- `internal/circuit/circuit.go` — `Breaker`, `State`, `Config`, `RecordFailure`/`RecordSuccess`/`Allow`/`Execute` (lines 80-126)

Tests:
- `internal/circuit/circuit_test.go` — `TestCircuitBreaker_StateTransitions`, `TestCircuitBreaker_HalfOpenFailureTripsOpen`, `TestCircuitBreaker_DefaultZeroConfig`

---

## Idempotency Deduplication Store

Research:
- `research/01-plan.md` — RQ5: Non-Idempotent Operations & Timeout Ambiguity
- `research/03-evidence.md` (Evidence 5) — Stripe documentation, indeterminate state, Idempotency-Key
- `research/04-contradictions.md` — Contradiction 3 (Practice B: idempotency key guard)
- `research/05-report.md` — Finding 5: Timeout Ambiguity & Idempotency Key Reconciliation

Implementation:
- `internal/idempotency/idempotency.go` — `Store`, `Record`, `NewStore`, `Get`, `Set` (lines 19-52)

Tests:
- `internal/idempotency/idempotency_test.go` — `TestStore_GetSet`, `TestStore_ConcurrentAccess`, `TestStore_LazyEvictionOnGet`

---

## Integration: Retry + Circuit Breaker + Idempotency

Research:
- `research/01-plan.md` — (overall concept)
- `research/04-contradictions.md` — Contradiction 3: retry budget, idempotency guard
- `research/05-report.md` — Findings 2-5

Implementation:
- `cmd/demo/main.go` — Demo 1 (deadline), Demo 2 (retry), Demo 3 (circuit), Demo 4 (idempotency)

Integration Tests:
- `tests/integration_test.go` — `TestIntegration_RetryWithCircuitBreaker`, `TestIntegration_IdempotentRetry`

Verified Execution:
- `engineering/03-execution-result.md` — `go test`, `go test -race`, `go run ./cmd/demo` output

---

## Database & Queue Guardrails (Research Context — Not Implemented in Lab)

Research:
- `research/01-plan.md` — RQ7: Database & Worker Timeout Controls
- `research/03-evidence.md` (Evidence 7 & 8) — PostgreSQL timeouts, queue worker DLQ
- `research/05-report.md` — Finding 7: Database and Queue Worker Guardrails

Note: Dapatkan dalam konteks riset. Tidak ada implementasi di kode Go; hanya disebutkan sebagai guardrail yang direkomendasikan.

---

## Audit Verdicts (Approval Gate)

Research:
- `research-audit/07-verdict.md` — **APPROVED** (2026-09-28)
- `research-audit/02-source-audit.md` — 5 sources reviewed, all PASS
- `research-audit/06-gaps.md` — 3 gaps, semua LOW severity

Engineering:
- `engineering-audit/06-verdict.md` — **APPROVED**
- `engineering-audit/02-code-audit.md` — 4 code findings, all PASS
- `engineering-audit/05-gaps.md` — no gaps
- `engineering-revision/03-revision-result.md` — 0 issues, re-validation PASS

Content:
- `content-audit/09-verdict.md` — REJECTED (no content to audit sebelum draft ini dibuat)
