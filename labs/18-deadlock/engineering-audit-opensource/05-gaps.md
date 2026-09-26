# Gaps

Target Lab: labs/18-deadlock

## Gap 1

Type: IMPLEMENTATION_OVERCLAIM
Severity: LOW
Location: README.md:12 ("Completely prevents deadlocks")
Evidence: True within 2-account scope; unproven for 3+ locks. Engineering notes already scope to A-to-B.
Action: Qualify wording to 2-account scope. Non-blocking.

## Gap 2

Type: UNHANDLED_ERROR (dead branch)
Severity: LOW
Location: internal/transfer/transfer.go:63 (`err != context.DeadlineExceeded` check)
Evidence: `bank.Lock` maps all ctx expiry to `ErrDeadlock`; `DeadlineExceeded` unreachable. Retry correct regardless.
Action: Optional cleanup. Non-blocking.

## Gap 3

Type: MISSING_TEST
Severity: LOW
Location: tests/transfer_test.go
Evidence: No single-transfer no-contention happy path; no maxRetries-exhaustion path. Core claims covered without them.
Action: Optional additions. Non-blocking.

No HIGH. No CRITICAL. No MISSING_EDGE_CASE beyond out-of-scope 3+ locks. No RACE_CONDITION. No BROKEN_IMPLEMENTATION. No FAKE_DEMO. No FAKE_BENCHMARK. No UNVERIFIED_RESULT.
