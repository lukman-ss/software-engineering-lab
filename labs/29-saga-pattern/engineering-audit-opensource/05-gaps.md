# Gaps

## Identified Gaps

1. MISSING_TEST: PaymentService failure path (shouldFail=true) is not tested and not demonstrated. Payment failure handling would be generic compensation but should be verified.
2. MISSING_TEST: Duplicate compensation invocation safety (compensate called twice) not tested. Implementation would re-run Compensate; idempotency depends on underlying service not being idempotent for Cancel/Refund.
3. MISSING_TEST: Orchestrator re-execution or shared orchestrator reuse not tested; not claimed supported.
4. WARNING: Semantic lock on OrderService is never released on process crash or saga abandonment — resource leak / stale lock. Documented limitation.
5. WARNING: Shared Orchestrator instance is not safe for concurrent Execute calls due to steps slice mutation; tests avoid by using per-worker instances.

## No Fake Results / Benchmarks
- Engineering execution result (engineering/03-execution-result.md) was independently re-run and matches exactly:
  - go test -v ./... → all 9 tests PASS.
  - go test -race ./... → ok.
  - go run ./cmd/demo → output matches documented demo output character-for-character.
- No benchmarks claimed, none present.

## No Fabricated Output
- Demo output verified by live execution.
- Test output verified by live execution.
