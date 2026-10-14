---
title: "SLO, SLI, and Error Budgets in Go: Multi-Window Multi-Burn-Rate Alerting"
published: false
description: "Implement time-bucketed sliding windows, mathematical error budget depletion, and dual-window burn rate alert engines in Go."
tags: "go, devops, sre, architecture"
---

Operating distributed services without user-centric reliability targets produces operational confusion. Infrastructure metrics such as CPU saturation, host memory pressure, or disk I/O do not reflect user pain. A fleet running at 15% CPU utilization can return HTTP 500 errors to every caller, while saturated machines might successfully process every critical transaction within latency targets.

In **Software Engineering Lab 24**, we implement Service Level Indicators (SLI), Service Level Objectives (SLO), mathematical Error Budget tracking, and Google SRE multi-window multi-burn-rate alerting in Go.

```
Incoming Request ---> [ Caller / Router ]
                            |
                   (isGood Predicate)
                            |
                            v
               [ WindowTracker (Buckets) ]
               TotalCount, GoodCount, BadCount
                      |               |
             (Short Window)     (Long Window)
                      \               /
                       v             v
                    [ AlertEngine (Check) ]
                Dual Burn Rate Threshold Verification
```

---

## The Core Mental Model: Reliability Is Not 100%

Reliability targets must not be set to 100%. Aiming for 100% availability is economically unfeasible, creates misaligned engineering incentives, and exceeds the reliability ceiling of users' consumer hardware and cellular networks.

An SLO defines a quantitative target:

$$\text{SLO} = 99.9\% \quad (\text{allowed failure rate } = 0.1\%)$$

The Error Budget is the exact mathematical inverse of the target:

$$\text{Error Budget} = 1.0 - \text{SLO} = 0.001 \quad (0.1\%)$$

The error budget represents an explicit budget for innovation, deployments, and infrastructure risk. As long as the error budget remains positive, engineering teams ship features rapidly. When the budget is exhausted, release policies pause non-emergency deployments and mandate architectural reliability investments.

---

## 1. Defining SLI: Ratios Over Raw Metrics

A Service Level Indicator (SLI) is a quantifiable ratio of good events over total valid events:

$$\text{SLI} = \frac{\text{Good Events}}{\text{Total Events}}$$

In `internal/metrics/tracker.go`, events pass through an evaluation closure provided by the caller:

```go
type Event struct {
	Timestamp  time.Time
	Duration   time.Duration
	StatusCode int
	Endpoint   string
}

// Caller evaluates qualification before classification:
isGood := func(e metrics.Event) bool {
	return e.StatusCode < 500 && e.Duration <= latencyThreshold
}
```

By separating qualification logic into caller predicates, the metrics tracker remains clean and focused on aggregation.

---

## 2. Sliding-Window Aggregation and In-Memory Bucketing

Tracking rolling windows requires granular time bucketing. Rather than storing millions of individual request structs in RAM, `WindowTracker` groups events into fixed intervals (e.g., 10-second or 1-minute buckets):

```go
type Bucket struct {
	StartTime  time.Time
	TotalCount int64
	GoodCount  int64
	BadCount   int64
}

type WindowTracker struct {
	mu         sync.RWMutex
	windowSize time.Duration
	bucketSize time.Duration
	buckets    []*Bucket
}
```

### Thread-Safe Event Recording and Eviction

When recording events, `WindowTracker` enforces atomic mutex locking, evicts stale buckets beyond the rolling horizon, and increments the appropriate counters:

```go
func (w *WindowTracker) Record(e Event, isGood bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evict(e.Timestamp)

	bucketStart := e.Timestamp.Truncate(w.bucketSize)
	b := w.getOrCreateBucket(bucketStart)

	b.TotalCount++
	if isGood {
		b.GoodCount++
	} else {
		b.BadCount++
	}
}
```

When out-of-order timestamps arrive within the valid window, the tracker locates or inserts buckets at the correct index while pruning records that fall outside the active duration.

---

## 3. The SLO Evaluator and Release Freeze Enforcement

The evaluator computes current SLI status, remaining error budget, and deployment eligibility:

```go
type Config struct {
	Name             string
	TargetUptime     float64       // e.g. 0.999 for 99.9%
	LatencyThreshold time.Duration
	WindowDuration   time.Duration
}

type Status struct {
	TargetSLO       float64
	CurrentSLI      float64
	TotalRequests   int64
	GoodRequests    int64
	BadRequests     int64
	BudgetRemaining float64
	CanDeploy       bool
}
```

### Mathematical Budget Formulation

In `internal/slo/evaluator.go`, budget calculation uses absolute event counts:

```go
func (e *Evaluator) Evaluate(now time.Time) Status {
	total, good, bad := e.tracker.Summary(now)

	if total == 0 {
		return Status{
			TargetSLO:       e.config.TargetUptime,
			CurrentSLI:      1.0,
			BudgetRemaining: 0,
			CanDeploy:       true,
		}
	}

	sli := float64(good) / float64(total)
	totalErrorBudget := (1.0 - e.config.TargetUptime) * float64(total)
	budgetRemaining := totalErrorBudget - float64(bad)
	canDeploy := budgetRemaining > 0

	return Status{
		TargetSLO:       e.config.TargetUptime,
		CurrentSLI:      sli,
		TotalRequests:   total,
		GoodRequests:    good,
		BadRequests:     bad,
		BudgetRemaining: budgetRemaining,
		CanDeploy:       canDeploy,
	}
}
```

When zero requests exist, the evaluator returns `CurrentSLI: 1.0` and `CanDeploy: true` to prevent false deployment freezes during service bootstrap.

---

## 4. Multi-Window Multi-Burn-Rate Alerting

Naïve alerting based on single-window thresholds creates severe operational problems:
- Short windows (e.g., 5 minutes) trigger false alerts on transient network glitches.
- Long windows (e.g., 24 hours) delay notification during catastrophic outages, causing high mean time to detect (MTTD).

The Google SRE framework solves this using **Multi-Window Multi-Burn-Rate Alerting**.

### Burn Rate Mechanics

Burn rate represents the consumption speed of an error budget relative to nominal allowance:

$$\text{Burn Rate} = \frac{\text{Actual Error Rate}}{\text{Allowed Error Rate}} = \frac{\text{Bad Events} / \text{Total Events}}{1.0 - \text{Target SLO}}$$

A burn rate of 1.0 consumes exactly 100% of the budget across the entire measurement period. A burn rate of 14.4x consumes 2% of a 30-day budget within 1 hour.

```
Burn Rate Table:
========================================================================
Alert Tier    Severity   Burn Rate   Short Window   Long Window   Budget
========================================================================
Page Fast     PAGE       14.4x       5 min          1 hour        2%
Ticket Slow   TICKET      6.0x       30 min         6 hours       5%
========================================================================
```

### Dual-Window Alert Verification

In `internal/alerting/engine.go`, an alert fires only if **both** windows exceed the burn rate threshold:

```go
func (e *Engine) Check(now time.Time) []Alert {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var alerts []Alert
	sTotal, _, sBad := e.shortTracker.Summary(now)
	lTotal, _, lBad := e.longTracker.Summary(now)

	if sTotal == 0 || lTotal == 0 {
		return alerts
	}

	sErrorRate := float64(sBad) / float64(sTotal)
	lErrorRate := float64(lBad) / float64(lTotal)

	allowedRate := 1.0 - e.sloConfig.TargetUptime
	sBurn := sErrorRate / allowedRate
	lBurn := lErrorRate / allowedRate

	for _, rule := range e.rules {
		if sBurn >= rule.BurnRateThreshold && lBurn >= rule.BurnRateThreshold {
			alerts = append(alerts, Alert{
				RuleName:      rule.Name,
				Severity:      rule.Severity,
				ShortBurnRate: sBurn,
				LongBurnRate:  lBurn,
				FiredAt:       now,
			})
		}
	}
	return alerts
}
```

If a burst of errors spikes the short window to 100x but the long window remains at 0.1x, no alert fires. The system suppresses transient noise.

---

## 5. Empirical Verification and Concurrency Guarantees

The test suite in `tests/slo_test.go` verifies behavioral correctness across five deterministic test cases:

```bash
$ go test -race -v ./...
=== RUN   TestMetricsWindowTracker
--- PASS: TestMetricsWindowTracker (0.00s)
=== RUN   TestSLOEvaluator
--- PASS: TestSLOEvaluator (0.00s)
=== RUN   TestAlertEngineBurnRate
--- PASS: TestAlertEngineBurnRate (0.01s)
=== RUN   TestOutOfOrderTimestamps
--- PASS: TestOutOfOrderTimestamps (0.00s)
=== RUN   TestEvaluatorZeroTraffic
--- PASS: TestEvaluatorZeroTraffic (0.00s)
=== RUN   TestConcurrencyMetrics
--- PASS: TestConcurrencyMetrics (0.00s)
PASS
ok  	labs/24-slo-sli-error-budget/tests	2.904s
```

Key verified characteristics:
- **Noise Suppression**: `TestAlertEngineBurnRate` verifies that a 100x spike restricted to the short window produces zero alerts when the long window remains clean.
- **Race Safety**: `TestConcurrencyMetrics` spawns 20 goroutines executing 100 requests simultaneously (2,000 requests total) under `-race` without data races or missing records.
- **Zero Traffic Handling**: `TestEvaluatorZeroTraffic` validates `CanDeploy: true` when request volume is zero.

---

## 6. Real-World Execution: 4-Phase Lifecycle

Running the demonstration in `cmd/demo/main.go` illustrates the complete operational lifecycle:

```bash
$ go run ./cmd/demo
============================================================
  LAB 24: SLO, SLI, AND ERROR BUDGET TRACKING DEMO
============================================================
Target SLO: 99.900% | Allowed Error Rate: 0.100%

--- PHASE 1: Baseline Traffic (Healthy) ---
Total: 1000 | Good: 1000 | Bad: 0
Current SLI: 100.0000%
Budget Remaining: +1.0000 errors
Deploy Status: ALLOWED (CanDeploy = true)

--- PHASE 2: Incident Injection (10% Error Burst) ---
Injected 100 requests with 10 failures...
Total: 1100 | Good: 1090 | Bad: 10
Current SLI: 99.0909%
Budget Remaining: -8.9000 errors
Deploy Status: BLOCKED (CanDeploy = false)

--- PHASE 3: Alert Engine Evaluation ---
Active Alerts: 1
>>> ALERT TRIGGERED: [TICKET] Slow Burn Alert (6.0x - 5% in 6h)
    | ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)

--- PHASE 4: Criticality-Based SLO Comparison ---
[CRITICAL]  Payment Service: Target=99.900% | SLI=99.091% | Budget=-8.90 | Deploy=BLOCKED
[STANDARD]  Reports Service: Target=95.000% | SLI=90.000% | Budget=-5.00 | Deploy=BLOCKED
```

### Tiered Criticality in Practice

In Phase 4, the demo contrasts critical financial endpoints against standard reporting pipelines:
- **Payment Service**: 99.9% SLO allows only 0.1% errors. An incident generating 10 errors on 1,100 requests exhausts the budget by -8.90 errors.
- **Reports Service**: 95.0% SLO allows a 5.0% failure rate. While a severe 10% outage still triggers a freeze, smaller transient failures that breach payment thresholds are safely absorbed by reports.

---

## Implementation Limitations & Production Boundaries

This laboratory focuses on core mathematical patterns and thread-safe data structures. Translating this directly to production requires acknowledging distinct boundaries:

1. **In-Memory Volatility**: The implementation stores bucket counters in RAM. A pod restart discards metric state. Production architectures record counters into Prometheus, VictoriaMetrics, or Datadog.
2. **Time Window Compression**: The demo compresses time scales (evaluating 5m/60m windows over seconds) to facilitate deterministic testing rather than waiting 30 days.
3. **Latency Histograms**: The metrics tracker uses boolean qualification closures rather than full quantile latency histograms (P95/P99).
4. **Static Engine Windows**: `AlertEngine` evaluates two global tracker windows constructed at initialization, rather than instantiating dynamic windows per rule struct.

---

## Repository and Code Reference

The complete Go implementation, multi-window alert engine, and test harness are available on GitHub:

- Repository: `https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/24-slo-sli-error-budget`
- Verified Commit SHA: `ffb1d0c42b78d71a36d1e3b9769d01bd8f458f67`
- Test Verification: `go test -race -v ./...` (All tests passing)
