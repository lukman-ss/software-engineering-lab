# Docs vs Code

## Comparison Targets
- README.md (Expected Behavior section, line ~105-140)
- `go run ./cmd/demo` actual output
- implementation code
- tests

## Findings

### DOC_CODE_MISMATCH — README "Expected Behavior" error string differs from implementation.

README line 114-116 (Expected Behavior section):
```
request=1 result=err=payment failed: status 500 duration=366.75µs state=CLOSED
request=2 result=err=payment failed: status 500 duration=96.08µs state=CLOSED
request=3 result=err=payment failed: status 500 duration=76.38µs state=OPEN
```

Actual demo output:
```
request=1 result=err=payment failed: status 500 body internal payment server failure
 duration=464.708µs state=CLOSED
request=2 result=err=payment failed: status 500 body internal payment server failure
 duration=164µs state=CLOSED
request=3 result=err=payment failed: status 500 body internal payment server failure
 duration=128.959µs state=OPEN
```

The implementation emits `... status 500 body <body>` (client.go:42), but README's expected output omits the ` body <body>` portion. The README was written before the `body` field was added (or is stale). This is a DOC_CODE_MISMATCH: the printed example does not match the actual error format. Severity: MEDIUM (documentation misrepresents the runtime output format; a reader copying the example would be confused).

### DOC_CODE_MISMATCH — README "Expected Behavior" timing values are illustrative.

README lists exact microsecond timings (366.75µs, etc.) but README's own note (line 63) says timeouts are "illustrative". Actual timings are non-deterministic. This is acceptable and documented, not a mismatch; flagged only because exact numbers are shown. Severity: LOW (informational).

### DOC_OK — State machine prose matches code.
- README CLOSED / OPEN / HALF_OPEN prose matches `advanceLocked`, `onSuccessLocked`, `onFailureLocked` exactly. PASS.

### DOC_OK — Fail-fast, recovery, throttle semantics match demo.
- Scenarios 1-4 outputs match README Expected Behavior for the structural parts (state transitions, fail-fast messages, recovery probe). PASS.

### DOC_OK — Architecture diagram matches implementation layout.
- Checkout Service → Circuit Breaker Proxy → Payment Client → Fake Server matches `checkout.Service` composing `circuitbreaker.Breaker` and `payment.Client`. PASS.

### DOC_OK — Run instructions accurate.
- `go run ./cmd/demo` and `go test -v ./...` / `go test -race ./...` all work as documented. PASS.

### DOC_OK — Observability note honest.
- README correctly states metrics are "omitted in this minimal lab implementation." Code does not instrument metrics. No mismatch. PASS.

## Summary

| Check | Result |
|---|---|
| State machine prose ↔ code | PASS |
| Scenario structure ↔ demo | PASS |
| Run / test instructions ↔ actual commands | PASS |
| Error message format in README "Expected Behavior" ↔ code | WARNING (omitted `body` field — stale example) |
| Timing values in README ↔ code | LOW (illustrative, as documented) |
