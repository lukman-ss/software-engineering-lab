## Revision 1

Audit Issue: Timeout Claim Mismatch
Severity: LOW
Files Changed: `labs/15-load-testing/engineering/01-design.md`
Action: Changed "timeouts for the 95th percentile" to "severe tail degradation for the 95th and 99th percentiles" in Failure Scenario to match test runtime realities.
Verification: Verified documentation strictly aligns with execution logs where `Errors: 0` occurs while retaining `P95` and `P99` severe degradation metrics.
Status: RESOLVED

## Revision 2

Audit Issue: Unhandled Error Context for NewRequest
Severity: LOW
Files Changed: `labs/15-load-testing/internal/loadtest/runner.go`
Action: Wrapped `errs++` block under `http.NewRequestWithContext` with `if ctx.Err() == nil` to avoid counting test shutdown context cancellations as standard load test errors.
Verification: Run `go test -v ./...` and `go run ./cmd/demo` to confirm tests remain green and behavior is unaffected.
Status: RESOLVED
