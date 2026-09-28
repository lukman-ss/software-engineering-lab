# Docs vs Code

## Sources Compared
- README.md (lab root)
- engineering/01-design.md (design + expected behavior + test strategy)
- engineering/02-implementation-notes.md (design decisions + what is/isn't demonstrated)
- engineering/03-execution-result.md (recorded build/test/race/demo)
- Source code (internal/*, cmd/demo)
- tests/slo_test.go
- Demo runtime output (re-executed)

> Research/content files excluded from this docs-vs-code pass per pipeline override (audit implementation + tests only).

## 1. README vs Code
| README Claim | Code Reality | Status |
|---|---|---|
| internal/metrics = sliding-window time-bucketed tracker for good vs total events | tracker.go implements exactly this | PASS |
| internal/slo = SLI ratios, remaining Error Budget, release freeze policy | evaluator.go computes SLI (good/total), budget, CanDeploy | PASS |
| internal/alerting = multi-window multi-burn-rate alerting (fast + slow) | engine.go computes short/long burn vs factor | PASS |
| cmd/demo = baseline tracking, budget depletion, alert triggering | demo has Phase 1 baseline, Phase 2 incident+depletion, Phase 3 alerts | PASS |
| tests/ = unit + concurrency tests for thread-safety + math correctness | tests present | PASS |
| `go test ./...` | Verified PASS | PASS |
| `go test -race ./...` | Verified PASS | PASS |
| `go run ./cmd/demo` | Verified PASS | PASS |

README is accurate. No DOC_CODE_MISMATCH.

## 2. Execution Record vs Re-executed Output
engineering/03-execution-result.md records demo output. Re-running `go run ./cmd/demo` reproduced **identical** output:
- Phase 1: Total 1000, Good 1000, Bad 0, Budget 1.00, deploy true.
- Phase 2: Total 1100, Good 1090, Bad 10, Budget -8.90, deploy false.
- Phase 3: TICKET Slow Burn 9.09x (6.0x threshold).
- Phase 4: Reports SLI 90.0%, Budget -5.00.

Verdict: recorded execution is authentic. No FAKE_DEMO / FAKE_BENCHMARK / UNVERIFIED_RESULT.

## 3. Design Doc Mismatches (engineering/01-design.md)

### RESEARCH_IMPLEMENTATION_MISMATCH (component attribution)
Design doc `01-design.md:34` states:
> SLOEvaluator: Calculates SLI, remaining Error Budget, **and current Burn Rate**

Code reality (`internal/slo/evaluator.go`): Evaluator computes **only** SLI + error budget + CanDeploy. Burn rate is calculated in `internal/alerting/engine.go` (`CalculateBurnRate`). No burn-rate field exists in slo.Status.
Severity: MEDIUM — design misattributes a component's responsibility. README itself is correct (does not make this claim); only the design notes are wrong.

### RESEARCH_IMPL_MISMATCH (demo recovery not demonstrated)
Design doc `01-design.md:29` states cmd/demo demonstrates:
> "normal traffic, incident budget burn, alerting, and recovery"

Code reality (`cmd/demo/main.go`): Phases 1-4 show baseline, incident/burn, alerting, and endpoint-criticality comparison. **No recovery phase** (no decay of burn rate back to normal after incident, no budget "healing").
Severity: LOW — the demo content is a subset of the design claim.

### TEST_CLAIM_MISMATCH (100% errors not tested)
Design/Test Strategy `01-design.md:38`:
> "Unit tests for ... edge cases (zero traffic, **100% errors**, rolling window expiry)."

Code reality (tests/slo_test.go): No test exercises a 100%-failure scenario (every bad, no good). Highest bad ratio tested is 10% (Phase 2) and 10% transient spike.
Severity: LOW-MEDIUM — test strategy over-states coverage.

## 4. Implementation-Notes Mismatches (engineering/02-implementation-notes.md)

### DOC_CODE_MISMATCH (ring buffer terminology)
Implementation notes `02-implementation-notes.md:21`:
> "In-memory time-bucketed **ring buffer** (WindowTracker) used to aggregate..."

Code reality (`internal/metrics/tracker.go`): buckets stored in a `[]Bucket` slice with append + slice-reslice eviction — **not** a ring buffer (no fixed-capacity overwrite; eviction discards old heads).
Severity: LOW — terminology only.

### DOC_CODE_MISMATCH (latency histogram terminology)
Design doc `01-design.md:26`:
> "internal/metrics: Sliding-window event recorder (histogram latency buckets & success counts)."

Implementation notes `02-implementation-notes.md:16-17`:
> "Error Budget Calculation: Budget defined as (1.0 - target_slo) * total_events. Consumed on bad events or **latency threshold breaches**."

Code reality: `Bucket` holds only `TotalCount/GoodCount/BadCount` (counts) and a single `isGood` predicate that returns bool (e.g., status code < 500 AND duration <= threshold). There is **no latency histogram**; latency is treated as a binary good/bad gate via the predicate, and the `Duration` field is recorded in `Event` but not bucketed by magnitude. The budget-consumed-on-latency-breach claim is accurate at the predicate level (a slow request counts as bad), but the "histogram" terminology is wrong.
Severity: LOW — terminology/imprecision.

## 5. In-Code Comment Mismatches (cmd/demo/main.go)

### DOC mismatch: variable naming
`cmd/demo/main.go:24`:
```go
window30d := 30 * time.Minute
```
Name implies 30 days; value is 30 minutes. Used consistently as a 30-minute window for sloTracker. No logic impact.
Severity: LOW.

### DOC mismatch: rule descriptions
`cmd/demo/main.go:38-49` comment strings:
- `"Critical Fast Burn (14.4x - 2% in 1h)"` — for a 99.9% SLO, 2% error rate ⇒ burn = 20x, not 14.4x.
- `"Slow Burn Alert (6.0x - 5% in 6h)"` — 5% error rate ⇒ burn = 50x, not 6.0x.

The rule *factors* (14.4, 6.0) match the design doc `01-design.md:9` ("14.4x for fast burn, 6x for slow burn"). Only the inline comment describing the error-rate equivalence is wrong. Behavior is correct (10% error ⇒ 100x burn exceeds both thresholds).
Severity: LOW.

## Summary
- README vs code: PASS (no mismatches).
- Execution record vs re-execution: PASS (authentic, no fabricated output).
- Design/implementation notes: 3 mismatches (1 MEDIUM — burn-rate attribution; 1 LOW — missing recovery; 1 LOW — 100%-errors test gap).
- In-code comments: 2 mismatches (both LOW — naming + rule descriptions).
- All mismatches are documentation imprecision EXCEPT the burn-rate attribution, which misrepresents component responsibility. None alter runtime behavior or correctness of what IS implemented.
