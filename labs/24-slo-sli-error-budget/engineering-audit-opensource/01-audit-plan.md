# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Go module: labs/24-slo-sli-error-budget (Go 1.22+)
Workspace root: /Users/tthi/Documents/LUKMAN/software-engineering-lab
Go toolchain verified: go1.26.7 darwin/arm64

## Implementation Files (verified present)
- go.mod
- internal/metrics/tracker.go   -> Event, Bucket, WindowTracker (Record / Summary / evictStaleLocked)
- internal/slo/evaluator.go       -> Config, Evaluator, Evaluate -> Status
- internal/alerting/engine.go     -> BurnRateRule, AlertEngine (CalculateBurnRate / Check -> []AlertResult)
- cmd/demo/main.go                -> runnable demonstration (PHases 1-3)

## Tests (verified present)
- tests/slo_test.go
  - TestMetricsWindowTracker     (sliding window + eviction)
  - TestSLOEvaluator             (SLO/budget math + CanDeploy policy)
  - TestAlertEngineBurnRate      (burn-rate alert triggering)
  - TestConcurrencyMetrics       (concurrent Record under race detector)

## Executable/Demo
- cmd/demo (go run ./cmd/demo)

## Approved Research Inputs
- Per PIPELINE OVERRIDE, research/content is NOT audited in this stage.
- Implementation claims are validated against README.md, engineering design notes, and the code/tests directly.

## Main Claims To Verify
1. README accurately describes package responsibilities and how to build/test/run the demo.
2. Build compiles cleanly (go build ./...).
3. go vet is clean.
4. go test ./... passes.
5. go test -race ./... passes (concurrency safety).
6. Demo runs and output is genuine (reproduced), math internally consistent.
7. SLO/SLI/error-budget math matches the documented formulas.
8. No fabricated benchmark/result.

## Commands To Run
- go build ./...
- go vet ./...
- gofmt -l .
- go test -count=1 -v ./...
- go test -race -count=1 ./...
- go test -coverpkg=./internal/... -coverprofile=/tmp/cov24.out ./tests/ && go tool cover -func=/tmp/cov24.out
- go run ./cmd/demo

## Primary Risks
- Test coverage reported as "100%" in engineering design but actual coverage is < 100% on CalculateBurnRate and NewWindowTracker (edge paths untested).
- gofmt formatting drift in three files (style only, not a correctness gate).
- Demo output depends only on deterministic arithmetic, so it is reproducible and not time-of-day dependent.
- Burn-rate alert fires only when BOTH short and long windows exceed threshold (a deviation from canonical Google SRE which can use OR semantics for fast-burn paging); verified against implementation, not a defect if intentional.
