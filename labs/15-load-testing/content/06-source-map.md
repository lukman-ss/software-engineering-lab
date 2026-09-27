# Source Map

## Load Testing Type Definitions

Research:
`research/05-report.md` (Finding 1), `research/02-sources.md` (Sources 1-6)

Implementation:
`engineering/01-design.md` (Architecture)

Tests:
`tests/loadtest_test.go` (TestLoadTest_SmokeVsStress — smoke vs stress contrast)

---

## Key Metrics (P50, P95, P99, RPS, Error Rate)

Research:
`research/05-report.md` (Finding 2), `research/03-evidence.md` (Evidence 2)

Implementation:
`internal/loadtest/metrics.go` (CalculateMetrics, percentile)

Tests:
`internal/loadtest/metrics_test.go` (TestCalculateMetrics, TestCalculateMetrics_Invariants)

---

## Bottleneck Identification Methodology

Research:
`research/05-report.md` (Finding 4), `research/03-evidence.md` (Evidence 3)

Implementation:
`internal/loadtest/runner.go` (HTTP transport tuning), `internal/server/server.go` (semaphore simulation)

Tests:
`internal/loadtest/metrics_test.go` (TestCalculateMetrics_Invariants — verifies percentile monotonicity)

---

## Tool Comparison (k6, JMeter, Locust, Gatling)

Research:
`research/05-report.md` (Finding 3), `research/02-sources.md` (Sources 1, 13, 16, 21), `research/03-evidence.md` (Evidence 7)

Implementation:
`engineering/02-implementation-notes.md` (Known Limitations — no external tool integration; this is a self-contained simulation)

Tests:
(none — lab implements in-process load test, does not invoke external tools)

---

## Server Semaphore (Connection Pool Simulation)

Research:
`research/05-report.md` (Finding 4 — downstream resource saturation causes non-linear latency growth)

Implementation:
`internal/server/server.go` (New, handleBooking, semaphore)

Tests:
`tests/loadtest_test.go` (TestServer_MaxDBConnectionsBound — verifies peak connections ≤ MaxDBConnections)

---

## Load Runner (Concurrency Orchestrator)

Research:
`research/02-sources.md` (Source 18 — concurrent user calculation), `research/03-evidence.md` (Evidence 4)

Implementation:
`internal/loadtest/runner.go` (Run, NewRunner)

Tests:
`tests/loadtest_test.go` (TestLoadTest_SmokeVsStress, TestLoadTest_ErrorCount, TestLoadTest_DialError, TestLoadTest_SuccessAndErrorInvariant)

---

## Demo Execution (Smoke vs Stress Contrast)

Research:
(none — demo is implementation artifact)

Implementation:
`cmd/demo/main.go`, `engineering/03-execution-result.md`

Tests:
(none — demo is not directly tested; validated by manual execution in audit)

---

## Demo Output / Verified Behaviors

Research:
(none)

Implementation:
`engineering/03-execution-result.md` (sample run metrics)

Tests:
`tests/loadtest_test.go` (TestLoadTest_SmokeVsStress — programmatic assertion that stress P95 > smoke P95)

---

## Common Pitfalls

Research:
`research/05-report.md` (Finding 5), `research/03-evidence.md` (Evidence 5)

Implementation:
`engineering/01-design.md` (Implementation Decisions — custom transport to avoid client-side bottleneck), `engineering/02-implementation-notes.md` (Trade-offs)

Tests:
`engineering-audit/04-docs-vs-code.md`

---

## Load Testing Timing in SDLC

Research:
`research/05-report.md` (Finding 6), `research/03-evidence.md` (Evidence 6)

Implementation:
(none — timing is operational guidance, not code)

Tests:
(none)

---

## Audit Status

Research:
`research-audit/07-verdict.md` → APPROVED

Engineering:
`engineering-audit/06-verdict.md` → APPROVED

Audit Gaps:
`research-audit/06-gaps.md` (Gap 3 — concurrent user formula needs think-time adaptation), `engineering-audit/05-gaps.md` (none critical)