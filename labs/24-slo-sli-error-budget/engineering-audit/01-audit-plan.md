# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `go.mod`
- `internal/metrics/tracker.go` — sliding-window time-bucketed event tracker
- `internal/slo/evaluator.go` — SLI ratio, error budget, release-freeze (CanDeploy) evaluator
- `internal/alerting/engine.go` — multi-window multi-burn-rate alert engine
- `cmd/demo/main.go` — executable simulation (4 phases)

Tests:
- `tests/slo_test.go` (6 tests: TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate, TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics)

Executable/Demo:
- `go run ./cmd/demo` — 4-phase simulation (baseline traffic, incident, burn-rate alerting, endpoint criticality comparison)

Approved Research Inputs:
- `research-audit/07-verdict.md` = APPROVED (research-only; code/tests NOT_APPLICABLE)

Main Claims To Verify:
1. SLI = good_events / total_events (rolling window).
2. Error Budget = (1 - SLO) * total; consumed by bad events; CanDeploy = false when remaining <= 0.
3. Burn rate = (bad/total) / (1 - SLO); alert fires when both short and long window burn rates exceed the rule factor.
4. Multi-window multi-burn-rate policy (14.4x fast / 6.0x slow, Google SRE-style).
5. Endpoint criticality: stricter SLO (99.9%) vs non-critical (95%).
6. Concurrency safety of WindowTracker (race-free).
7. README commands (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`) all work and match output.
8. Success criterion: "100% test coverage on core math and sliding window calculations" (design/01-design.md).

Commands To Run:
- `go build ./...`
- `go vet ./...`
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
- `go test -count=1 -coverpkg=./internal/... -coverprofile=/tmp/cov.out ./tests/ && go tool cover -func=/tmp/cov.out`

Primary Risks:
- Coverage overclaim ("100%") vs measured reality.
- Stale recorded demo output in engineering/03-execution-result.md (may omit Phase 4).
- Uncovered edge-case branches (CalculateBurnRate zero-total / SLO>=1; NewWindowTracker invalid args).
- Summary() using write Lock instead of RLock (read op serialized unnecessarily).
- Eviction keyed on event timestamp rather than wall clock (out-of-order old events may fail to evict stale buckets).
