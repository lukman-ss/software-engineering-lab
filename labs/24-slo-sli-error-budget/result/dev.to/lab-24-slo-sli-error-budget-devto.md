---
title: "SLO, SLI & Error Budget: Turning Reliability Into Deployable Numbers"
published: false
description: "A Go implementation of Google SRE principles: sliding-window event tracking, SLI/SLO evaluation, dual-window burn rate alerting, and automated deploy gates — all in pure standard library."
tags: go, sre, reliability, testing, softwareengineering
---

# SLO, SLI & Error Budget: Turning Reliability Into Deployable Numbers

A payment service degrades. CPU sits at 30%, memory is nominal, and the infrastructure dashboard stays green — yet 5% of user requests return 5xx errors. Lab 24 replaces that false confidence with a Google SRE–based SLO/SLI/Error Budget framework written in pure Go.

> Infrastructure metrics track machines. SLO metrics track users. Always measure what the customer experiences.

The lab does not connect to production systems, Prometheus, or PagerDuty. It runs five tests, a four-phase demo, and proves correctness with the race detector — zero external dependencies required.

## The Core Problem: Healthy Servers, Broken User Experience

Reliability teams historically monitor CPU, RAM, and disk — machine-level signals that do not correlate with application-layer errors. An always-green infrastructure dashboard can mask a 5% 5xx error rate on the primary endpoint. Without an objective SLO framework, every release becomes a risk gamble and reliability devolves into opinion rather than data-driven decision.

Dwell time compounds the problem. Without error budget tracking, teams discover incidents only after mass user complaints. A single regression can consume a year's uptime budget in minutes.

## Why Static Thresholds Fail

A single alert threshold cannot distinguish a transient spike from a systemic incident. Ten errors out of 100 requests looks critical (10%), but ten errors out of one million requests is negligible (0.001%). A static threshold produces either excessive false positives or dangerous false negatives depending on traffic volume.

Infrastructure metrics compound this blind spot. CPU spikes do not equal error rates; memory leaks occur without CPU pressure; latency P99 can degrade while error rate appears stable. Single-threshold alerting on machine metrics is inherently reactive.

## Mental Model: Reliability Is Not 100%

SLO sets a realistic target — 99.9% for a critical service, not 100%. The gap between target and actual is the error budget. Each bad event spends from that budget:

```
Error Budget = 1 - SLO
SLI = good_events / total_events
Burn Rate = actual_error_rate / allowed_error_rate
```

A 14.4× burn rate means the budget will be exhausted in 2% of its allotted time. A 6.0× burn rate indicates a slow burn over 5% of the window. The budget is spendable risk capital — when it runs out, innovation pauses until reliability recovers.

## Pipeline Architecture: Event → Tracker → Evaluator → Alert

Three modular components compose the pipeline:

```
Event{Timestamp, Duration, StatusCode, Endpoint}
  ↓ WindowTracker.Record()        (RWMutex + bucket eviction)
  ↓ Summary(now)                  → (total, good, bad)
  ↓ Evaluator.Evaluate()          → SLI, Budget Remaining, CanDeploy
  ↓ AlertEngine.Check()           → Multi-window Burn Rate Alerts
```

- `metrics.WindowTracker` — sliding window aggregation, O(bucket_count) query
- `slo.Evaluator` — SLI calculation, budget math, deploy gate logic
- `alerting.AlertEngine` — dual-window burn rate threshold check

Time windows: SLO tracker uses 30 days (compressed to 30 minutes in demo); alert short tracker uses 5 minutes; alert long tracker uses 60 minutes.

## WindowTracker: Sliding Window with Out-of-Order Support

`WindowTracker` uses `sync.RWMutex` for thread safety when hundreds of goroutines record events concurrently. The `Record()` method acquires a write lock, evicts stale buckets past `windowSize`, then searches for or creates a new bucket for the truncated timestamp.

```go
// internal/metrics/tracker.go:46
func (w *WindowTracker) Record(e Event) {
    w.mu.Lock()
    defer w.mu.Unlock()

    w.evictStaleLocked(e.Timestamp)

    bucketStart := e.Timestamp.Truncate(w.bucketSize)
    // ... linear search insert or append ...
}

func (w *WindowTracker) evictStaleLocked(now time.Time) {
    cutoff := now.Add(-w.windowSize)
    idx := 0
    for idx < len(w.buckets) && w.buckets[idx].StartTime.Before(cutoff) {
        idx++
    }
    if idx > 0 {
        w.buckets = w.buckets[idx:]
    }
}
```

If timestamps arrive out-of-order, the bucket is inserted at the correct chronological position via linear search, ensuring the slice stays sorted for efficient eviction. This handles clock skew and network reordering common in distributed systems.

## Failure Path: Latency Regression Exhausts Budget

A deployment introduces a latency regression: P95 moves from 50ms to 500ms. The `isGood` predicate requires `StatusCode < 500 && Duration <= 200ms`. Every new request now counts as bad.

With SLO 99.9% (allowed error rate = 0.1%) and a 10% error rate from the regression, burn rate hits 100×. A 30-day budget exhausts in approximately 43 minutes. The alert engine compares short and long window burn rates; when both exceed 6.0×, a TICKET alert triggers and `CanDeploy` flips false — halting all non-P0 deployments automatically.

## Mechanism vs Symptom: What Gets Measured

| Metric | Type | User Correlation | Alert Responsiveness |
|--------|------|-----------------|---------------------|
| CPU Usage | Internal system | Weak (indirect) | Very slow |
| Memory Leak | Internal system | Weak (slow correlation) | Slow (hours) |
| Error Rate (SLI) | User-perceived | Direct, real-time | Minutes |
| Latency P99 (SLI) | User-perceived | Direct, real-time | Minutes |
| Burn Rate | Composite metric | Predictive (forecasting) | Hours |

Infrastructure metrics diagnose machines; SLO metrics protect users. SLI signals provide immediate user-perspective data, while burn rate offers predictive warning before problems become systemic.

## Three Critical Decisions in Error Budgeting

1. **SLI Definition**: Choose metrics that reflect user experience, not server state. Define the `isGood` closure — typically a combination of status code and latency threshold.
2. **SLO Target**: Set realistic targets (99.9% for critical services). Avoid 100% — network and device errors are outside your control. Different endpoints may have different SLOs based on criticality.
3. **CanDeploy Policy**: `CanDeploy = false` when `budgetRemaining <= 0` (and `total > 0`). Halt all non-P0 deployments until budget recovers. Require postmortem if a single incident consumes > 20% of budget.

## Execution Evidence: Demo Phases 1–4

The terminal output below is direct execution from `cmd/demo/main.go`:

```
[PHASE 1] Simulating Baseline Traffic (1,000 requests, 100% success)...
Total: 1000 | Good: 1000 | Bad: 0
Target SLO: 99.900% | Current SLI: 100.0000% | Budget Remaining: 1.00
Deployment Allowed: true

[PHASE 2] Simulating Severe Incident (100 total, 10 errors)...
Total: 1100 | Good: 1090 | Bad: 10
Target SLO: 99.900% | Current SLI: 99.0909% | Budget Remaining: -8.90
Deployment Allowed: false (Budget exhausted)

[PHASE 3] Checking Multi-Window Burn Rate Alerts...
>>> ALERT TRIGGERED: [TICKET] Slow Burn Alert (6.0x - 5% in 6h)
    | ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)
```

Phase 1 shows baseline — budget remaining at 1.00 (one error allowed for 99.9% SLO). Phase 2 simulates an incident with 10 errors from 100 additional requests: budget drops to -8.90 and `CanDeploy` flips false. Phase 3 proves the alert engine triggers TICKET severity because both windows exceed the 6.0× threshold.

## Boundary Transparency: What Is Not Proven

Engineering integrity requires explicit boundary statements:

- **No metrics persistence**: All data is in-memory; lost on process restart.
- **No latency histogram**: Only a boolean good/bad predicate — no P50/P95/P99 distribution.
- **Time-compressed windows**: Demo uses 30-minute windows to represent 30-day periods.
- **Burn rate thresholds are recommendations**: 14.4×/6.0× follow Google SRE guidance, not universal standards.
- **Dead field**: `LatencyThreshold` exists in `Evaluator.Config` but is not read by `Evaluate()`.
- **No recovery simulation**: Budget does not refill after incident resolution in the demo.
- **Limited endpoint scope**: Demo uses 2 endpoints; hundreds of concurrent endpoints are untested.

## Architecture Trade-off: In-Memory vs TSDB

The demo prioritizes correctness and zero-dependency over production durability:

**SIM (Demo):**
- + Zero-dependency: Go standard library only
- + Thread-safe: `sync.RWMutex` proven with race detector
- + Out-of-order safe: bucket insertion maintains chronological order
- - No persistence: data loss on restart
- - No horizontal scaling: single-process bottleneck

**Production (Recommendation):**
- + Prometheus/Grafana for time-series persistence
- + Horizontal scaling via distributed TSDB cluster
- + Native SLO query language (Prometheus recording rules)
- - Additional infrastructure complexity
- - Integration testing overhead

The dual-window burn rate pattern can be replicated using Prometheus recording rules and alerting expressions.

## Concurrency: Goroutine-Safe Recording Pipeline

The concurrency test executes 20 goroutines, each recording 100 requests simultaneously to the same `WindowTracker`:

```go
// tests/slo_test.go:200
func TestConcurrencyMetrics(t *testing.T) {
    var wg sync.WaitGroup
    numGoroutines := 20
    requestsPerGoroutine := 100
    for g := 0; g < numGoroutines; g++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for i := 0; i < requestsPerGoroutine; i++ {
                tracker.Record(metrics.Event{...}) // RWMutex protects
            }
        }(g)
    }
    wg.Wait()
}
```

Result: 2,000 events recorded correctly, `good + bad == total` verified, and `go test -race` reports zero data races. Lock granularity is sufficient for high-throughput workloads.

## Evaluator: SLI and Error Budget Mathematics

The `Evaluator` computes three values: SLI ratio, total error budget, and remaining budget:

```go
// internal/slo/evaluator.go:41
func (e *Evaluator) Evaluate(now time.Time) Status {
    total, good, bad := e.tracker.Summary(now)

    var sli float64 = 1.0
    if total > 0 {
        sli = float64(good) / float64(total)
    }

    totalErrorBudget := (1.0 - e.config.TargetUptime) * float64(total)
    budgetRemaining := totalErrorBudget - float64(bad)
    canDeploy := total == 0 || budgetRemaining > 0
    // ...
}
```

Zero traffic returns SLI = 1.0 and `CanDeploy = true` to prevent false freezes on new services. The formula is straightforward: `SLI = good/total`, `budget = (1 - SLO) × N`, and the deploy gate triggers when `budgetRemaining ≤ 0` with active traffic.

## AlertEngine: Dual-Window Burn Rate

The alert engine calculates burn rate for two windows in parallel:

```go
// internal/alerting/engine.go:51
func (a *AlertEngine) CalculateBurnRate(total, bad int64) float64 {
    if total == 0 { return 0.0 }
    actualErrorRate := float64(bad) / float64(total)
    allowedErrorRate := 1.0 - a.targetSLO
    return actualErrorRate / allowedErrorRate
}
```

Alert triggers only when **both** windows exceed the threshold — an AND condition that filters transient spikes:

```go
// internal/alerting/engine.go:73
triggered := shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor
```

**Transient spike proof**: Short window shows 100× burn rate (10% errors in 5 min), long window shows 0.1× (0.01% in 60 min). Result: no alert — the spike is temporary, not systemic.

## Test Evidence: Budget Exhaustion and Zero Traffic

Two test cases cover critical boundary conditions:

```go
// TestBudgetExhaustion: 99 good + 1 bad → CanDeploy=true (budgetRemaining=0)
// Add 1 bad → 99 good + 2 bad → CanDeploy=false (budgetRemaining=-1)

// TestEvaluatorZeroTraffic: total=0 → SLI=1.0, CanDeploy=true
```

The first verifies `CanDeploy` flips false exactly when budget reaches zero or negative. The second prevents false freeze on fresh services with no traffic — a common operational gotcha.

## Test Evidence: Sustained vs Transient Burn

```go
// Sustained burn: 98 OK + 2 error in BOTH short and long windows
// → 20× burn rate, both windows >= 14.4x → PAGE alert triggered

// Transient spike: 90 OK + 10 error in short (100x), 
//                  9999 OK + 1 error in long (0.1x)
// → short >= 14.4x ✓, long >= 14.4x ✗
// → 0 alerts — spike filtered correctly
```

If only one window screams, it might be noise. Both must agree for action.

## Out-of-Order and Eviction Verification

```go
// Record t+5s event FIRST, then t+2s event SECOND
// → Bucket at t+2s inserted before t+5s bucket
// → Summary(t+6s): total=3, good=2, bad=1 ✓

// Advance to t+13s → cutoff = t+3s
// → Bucket at t+2s evicted, t+5s retained
// → Summary(t+13s): total=1, good=1, bad=0 ✓
```

Sorted bucket insertion with time-based eviction handles distributed clock skew gracefully without external coordination.

## Architecture Trade-offs Revisited: Pedagogical vs Production

The lab chose in-memory for pedagogical clarity and zero-dependency proof. The `AlertEngine` dual-tracker pattern and burn rate mathematics are portable to any time-series database. Production systems would replace the in-memory slices with Prometheus metrics, export them at short intervals, and evaluate burn rate via recording rules.

The core invariant — `SLI = good/total`, `budget = (1-SLO)×N`, deploy gate on exhaustion — remains identical regardless of storage backend.

## Before vs After: Observability Transformation

**Before (No SLO):**
- Dashboard: CPU 30%, Memory 60%, Disk I/O normal
- Alert: "CPU high" during load test — user experience unmeasured
- Decision: Deploy based on intuition, no budget data
- Response: Reactive, dwell time measured in hours

**After (With SLO):**
- Dashboard: SLI 99.09%, Budget Remaining -8.90, CanDeploy=false
- Alert: "Slow burn 9.09x" — both windows detected
- Decision: Auto-halt deployment, focus on postmortem
- Response: Proactive, budget consumption rate measurable

SLO dashboards show user reality; infrastructure dashboards show machine state.

## Implementation Checklist

**Fundamental:**
- [ ] Define SLI measuring user experience (good/total ratio)
- [ ] Choose realistic SLO ≤ 99.99% — avoid 100% perfectionism
- [ ] Define `isGood` closure: status code + latency threshold
- [ ] Implement `WindowTracker` with sliding window + bucket eviction

**Alerting:**
- [ ] Setup dual tracker: short (5–15m) and long (1–6h) windows
- [ ] Choose burn rate thresholds: 14.4× (PAGE) and 6.0× (TICKET)
- [ ] Implement AND logic: alert only when both windows exceed threshold
- [ ] Integrate with notification system (PagerDuty, Slack, email)

**Policy:**
- [ ] CanDeploy gate: auto-halt when `budgetRemaining <= 0`
- [ ] Exception policy: P0/security hotfixes deploy with approval
- [ ] Postmortem requirement: mandatory if single incident consumes > 20% budget
- [ ] SLO dashboard: real-time SLI, budget gauge, burn rate trend

Measurement without policy is vanity. Policy without measurement is opinion.

## Anti-Patterns and How to Reproduce Them

1. **100% SLO**: Setting `TargetUptime=0.999999` makes budget nearly zero. One error triggers `CanDeploy=false`. *Reproduce:* Set extreme SLO, record one bad event.
2. **Single-Window Alert**: Alert based only on short window. *Reproduce:* 10 errors in short window, 0 in long → false positive alert.
3. **Infrastructure SLI**: Using CPU/memory as primary SLI. *Reproduce:* CPU 30% with 10% error rate → no alert fires while user complaints rise.
4. **No Deploy Policy**: Budget exhausted but no auto-halt mechanism. *Reproduce:* Set `budgetRemaining=-5`, skip `CanDeploy` check → deployment continues.

The best SLO is the highest one your users will accept — not the highest number you can write.

## Running the Lab

```bash
# Navigate to lab directory
cd labs/24-slo-sli-error-budget

# Run all tests with race detector
go test -race -v ./...

# Run the demo
go run ./cmd/demo
```

Five test cases verify: WindowTracker aggregation, SLO evaluator budget math, AlertEngine dual-window logic, out-of-order timestamp handling, and concurrency stress testing. The four-phase demo shows baseline → incident → alert → cross-SLO comparison in real time.

## Long-term Architectural Principles

1. **Reliability Is a Feature**: SLO is not a meta-metric — it is a user-facing quality contract. Every percentage point represents real money and user trust.
2. **Automate the Halt**: `CanDeploy=false` must block the CI/CD pipeline automatically. Manual override only for P0/security with audit trail.
3. **Measure User Experience**: CPU 30% means nothing if error rate is 5%. CPU 90% is fine if SLI is 99.99%. Dashboards should answer: "Are our users happy?"

Automate the guardrails. Measure the user. Spend the budget wisely.

## Closing

Lab 24 proves that an SLO/Error Budget framework can be implemented in 128 lines of pure Go with zero external dependencies. Three modular components — `WindowTracker` for event aggregation, `Evaluator` for SLI and budget calculation, `AlertEngine` for multi-window burn rate — form a user-experience-direct observability system. Five test cases cover all boundary conditions and concurrency.

The final message: error budget must be spent, not saved. Unused budget means you are being too conservative and missing innovation opportunities.

An unspent error budget is wasted opportunity. Spend it on controlled risk, not hoarded perfection.

**Source:** [Software Engineering Lab 24](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/24-slo-sli-error-budget), `internal/metrics/tracker.go`, `internal/slo/evaluator.go`, `internal/alerting/engine.go`, `tests/slo_test.go`, `cmd/demo/main.go`, and `go.mod`; snapshot commit `13224b48`.
