# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
  - cmd/demo/main.go
  - internal/server/server.go
  - internal/loadtest/runner.go
  - internal/loadtest/metrics.go
Tests:
  - internal/loadtest/metrics_test.go
  - tests/loadtest_test.go
Executable/Demo:
  - cmd/demo (go run ./cmd/demo)
Approved Research Inputs:
  - research/01-plan.md
  - research/03-evidence.md
  - research/05-report.md
  - research-audit/03-claim-audit.md
Main Claims To Verify:
  1. Smoke (low VUs) vs Stress (high VUs) load differentiates baseline vs. resource-exhaustion behavior. (research Claim 1, Finding 1)
  2. Metrics reported are P50, P95, P99, error rate, RPS (percentile-over-average). (research Claim 2, Finding 2)
  3. Downstream saturation (DB connection pool / semaphore) causes non-linear tail-latency degradation for P95/P99. (research Finding 4 bottleneck isolation)
  4. All tests pass with zero race conditions. (engineering 01-design success criterion)
  5. Demo output is real and reproducible, not fabricated.
  6. README matches code.
Commands To Run:
  - go build ./...
  - go vet ./...
  - gofmt -l .
  - go test -v ./...
  - go test -race ./...
  - go run ./cmd/demo
Primary Risks:
  - Timing-dependent assertions in TestLoadTest_SmokeVsStress (P95 ordering) could be flaky on a loaded CI host.
  - Latency recorded only for successful requests; error request durations are not surfaced in P95/P99 (semantic gap vs. research Finding 2).
  - engineering/03-execution-result.md test list is stale vs. actual test set (DOC_CODE_MISMATCH).
  - runner.go not gofmt-clean (cosmetic trailing whitespace).
