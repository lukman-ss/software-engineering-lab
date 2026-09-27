# Gaps & Discrepancies

1. RACE_CONDITION — internal/backpressure/queue.go:66-81 (TrySubmit) + queue.go:87-94 (Stop). TOCTOU send-on-closed-channel panic under concurrent Submit+Stop (~50% panic rate). Confirmed via throwaway repro. Not caught by -race and not covered by any test. Severity: HIGH.
2. TEST_CLAIM_MISMATCH — engineering/03-execution-result.md lists 9 tests; 2 real tests omitted from the transcript. Severity: MEDIUM.
3. UNVERIFIED_RESULT — 03-execution-result.md §3 backpressure stats (Accepted=3,Rejected=3) do not match a fresh real run (Accepted=4,Rejected=2); jitter values differ by design. Output presented as canonical despite nondeterminism. Severity: MEDIUM.
4. IMPLEMENTATION_OVERCLAIM — design concurrency-safety guarantee extends to an untested/shutdown path that panics. Severity: HIGH (root cause overlaps Finding 1).
5. MISSING_TEST — Registry concurrent same-key Get (double-checked locking) unproven. Severity: LOW.
6. MISSING_TEST — LeakyBucket concurrency unproven. Severity: LOW.
7. MISSING_TEST — Middleware Retry-After value, body JSON shape, anonymous fallback, per-tenant isolation at HTTP level. Severity: LOW.
8. MISSING_EDGE_CASE — TokenBucket.AllowN>1, capacity=0, refill=0, negative allow; Decorrelated prev handling; unknown BackoffStrategy default branch. Severity: LOW.
9. DOC_CODE_MISMATCH — README omits engineering-revision/ and audit artifacts from structure listing. Severity: LOW.
10. gofmt drift on 3 files (cosmetic). Severity: LOW.

## Severity Aggregation

- CRITICAL: 0
- HIGH: 2 (1, 4)
- MEDIUM: 2 (2, 3)
- LOW: 6 (5,6,7,8,9,10)
- Total: 10 findings.
