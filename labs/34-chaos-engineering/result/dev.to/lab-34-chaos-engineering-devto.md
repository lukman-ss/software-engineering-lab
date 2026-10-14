---
title: "Chaos Engineering: Empirically Proving Distributed System Resilience"
published: false
description: "Explore Chaos Engineering in Go: steady-state monitoring, in-memory fault injection, circuit breakers, graceful fallback, and blast radius auto-abort."
tags: go, architecture, testing, distributed
---

# Chaos Engineering: Empirically Proving Distributed System Resilience

Distributed systems often assume external dependencies like payment gateways, databases, and microservices will always respond reliably. Under real production turbulence, unmitigated downstream failures cascade into thread exhaustion and widespread outages. Lab 34 implements an in-memory chaos engineering experiment in Go 1.22+ to prove how steady-state monitoring, circuit breaking, graceful fallbacks, and automated blast radius controls contain failures deterministically.

> Resilience without empirical verification is merely wishful thinking.

## Ketahanan Sistem Tanpa Bukti Empiris Adalah Asumsi

Engineering teams frequently assert that their distributed services are resilient because they wrote unit tests with static mocks or deployed load balancers. However, static assertions fail to capture runtime turbulence when a downstream provider begins stalling or returning sporadic errors. In production, cascading latency saturates connection pools, depletes goroutines, and degrades entire upstream call trees within seconds. Chaos Engineering converts resilience from an unproven assumption into verifiable empirical evidence by testing failure paths under controlled turbulence.

## Empat Langkah Baku Eksperimen Ilmiah Perangkat Lunak

Chaos Engineering follows an empirical scientific method consisting of four distinct phases. First, engineers define a steady-state metric representing user-perceived health, such as end-to-end transaction success rate. Second, they formulate a falsifiable hypothesis stating that steady-state health will persist despite injected turbulence. Third, they inject targeted faults into downstream dependencies under strict monitoring intervals. Fourth, they enforce automated blast radius controls that immediately abort the experiment and neutralize injectors if metrics breach acceptable boundaries.

## Steady State: Mengukur Dampak Pengguna, Bukan Node

Monitoring internal infrastructure attributes like CPU utilization, memory consumption, or disk I/O provides misleading signals during chaos experiments. A server with low CPU consumption may actually be deadlocked or dropping incoming client requests entirely. Steady state must be defined through customer-facing output metrics, specifically the ratio of failed transactions to total requests. In this lab, the monitor requires a warmup sample of at least five requests before evaluating health thresholds, avoiding premature false-positive aborts during cold starts.

## Topologi Komponen dan Siklus Closed-Loop Kontrol

The system topology separates fault injection from resilient client mitigation and steady-state monitoring. Requests arrive at a resilient client wrapped in a Circuit Breaker, which proxies calls to a downstream service subject to an in-memory fault injector. Successful and failed transactions report directly to the steady-state monitor. Concurrently, an experiment runner queries the monitor on a recurring ticker. If the error rate breaches the configured threshold, the runner halts the run and neutralizes the injector, creating a closed-loop feedback mechanism.

```text
[Client Request] ──> [Resilient Client (Circuit Breaker)]
                            │
                            ├─► [Downstream Service] (Fault Injector)
                            └─► [Graceful Fallback]  (On Failure / Open)
                                    │
                                    ▼
                          [Steady-State Monitor]
                                    │
                         (Ticker Evaluates Health)
                                    ▼
                         [Experiment Auto-Abort]
```

## Injeksi Fault: Simulasi Latensi dan Error Terprogram

The fault injector provides thread-safe latency and error synthesis without external dependencies. Using `sync.RWMutex`, goroutines read the active injection configuration concurrently while administrative methods update settings under write locks. Latency simulation uses `time.After` inside a `select` block that monitors `ctx.Done()`. This guarantees that if a calling client enforces a request deadline, the injected delay aborts immediately, preventing orphan goroutines from leaking in memory.

```go
func (fi *Injector) Execute(ctx context.Context) error {
	fi.mu.RLock()
	enabled := fi.enabled
	latency := fi.latency
	forceErr := fi.forceError
	fi.mu.RUnlock()

	if !enabled {
		return nil
	}
	if latency > 0 {
		select {
		case <-time.After(latency):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if forceErr {
		return ErrInjectedFault
	}
	return nil
}
```

## State Machine Circuit Breaker: Tiga Fasa Proteksi

The circuit breaker operates across three distinct states: Closed, Open, and Half-Open. In the Closed state, all traffic flows directly to the downstream dependency while a failure counter tracks consecutive errors. When failures reach the designated threshold, the breaker trips to Open, instantly rejecting requests with `ErrCircuitOpen` or routing them to a fallback. After a configured cooldown window elapses, the breaker enters Half-Open, admitting a trial request to probe downstream health before deciding whether to reset or reopen.

## Eksekusi Terproteksi dan Pengalihan Fallback

The `Execute` method on the circuit breaker orchestrates execution while minimizing lock contention. It acquires an exclusive lock to check whether cooldown expiration warrants a transition to Half-Open, then releases the lock before invoking the user-provided function. This prevents slow downstream calls from holding the mutex and blocking concurrent readers. If the downstream call fails and a fallback handler is provided, the breaker intercepts the error and executes the fallback transparently.

```go
func (cb *CircuitBreaker) Execute(fn func() error, fallback func() error) error {
	cb.mu.Lock()
	cb.checkStateLocked()

	if cb.state == StateOpen {
		cb.mu.Unlock()
		if fallback != nil {
			return fallback()
		}
		return ErrCircuitOpen
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		if cb.failures >= cb.threshold || cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.lastStateChg = time.Now()
		}
		if fallback != nil {
			return fallback()
		}
		return err
	}

	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failures = 0
		cb.lastStateChg = time.Now()
	} else if cb.state == StateClosed {
		cb.failures = 0
	}
	return nil
}
```

## Graceful Degradation: Menyelamatkan Metrik Bisnis

When downstream failures occur, a naive system returns HTTP 500 errors to consumers, destroying conversion rates and degrading availability metrics. Graceful degradation routes failing requests to secondary mechanisms, such as placing payment intents into an asynchronous reconciliation queue. From the customer perspective, the transaction is accepted ("Payment Queued"), allowing the monitor to record a success. By decoupling immediate external fulfillment from customer confirmation, business availability remains preserved at 100%.

## Pelacakan Rasio Kegagalan Menggunakan Operasi Atomik

The steady-state monitor leverages the `sync/atomic` package to maintain request tallies without lock overhead. Methods `RecordSuccess` and `RecordFailure` invoke `atomic.AddUint64`, allowing concurrent goroutines to update metrics without mutual exclusion stalls. The `ErrorRate` method calculates the current failure ratio on demand, while `IsHealthy` compares the ratio against a maximum permitted threshold once the initial sample threshold of five requests is satisfied.

```go
func (m *Monitor) RecordSuccess() {
	atomic.AddUint64(&m.totalRequests, 1)
}

func (m *Monitor) RecordFailure() {
	atomic.AddUint64(&m.totalRequests, 1)
	atomic.AddUint64(&m.failedRequests, 1)
}

func (m *Monitor) IsHealthy() bool {
	total := atomic.LoadUint64(&m.totalRequests)
	if total < 5 {
		return true
	}
	return m.ErrorRate() <= m.maxErrorRate
}
```

## Orkestrator Eksperimen dan Mekanisme Auto-Abort

The experiment runner manages the chaos lifecycle through a background loop that couples fault activation with automated monitoring. Upon launch, the runner activates the injector and initializes a recurring ticker. On every tick, the runner queries `monitor.IsHealthy()`. If the error rate exceeds the safety boundary, the runner aborts the experiment, invokes `terminate()`, records the breach percentage in `abortReason`, and immediately neutralizes the fault injector.

```go
func (e *Experiment) Run(ctx context.Context) error {
	e.mu.Lock()
	e.state = StateRunning
	e.mu.Unlock()

	e.injector.SetFault(e.config.Latency, e.config.ForceError)

	ticker := time.NewTicker(e.config.MonitorInterval)
	defer ticker.Stop()

	timeout := time.After(e.config.Duration)

	for {
		select {
		case <-ctx.Done():
			e.terminate(StateAborted, "context canceled")
			return ctx.Err()
		case <-timeout:
			e.terminate(StateCompleted, "")
			return nil
		case <-ticker.C:
			if !e.monitor.IsHealthy() {
				reason := fmt.Sprintf("Steady state breached: error rate %.2f%% exceeded threshold", e.monitor.ErrorRate()*100)
				e.terminate(StateAborted, reason)
				return fmt.Errorf("experiment aborted: %s", reason)
			}
		}
	}
}
```

## Pengendalian Blast Radius: Mencegah Malapetaka

Blast radius containment ensures that simulated failures never escape their sandbox or cause lasting damage. The experiment runner enforces containment through three interdependent safeguards: temporal bounds that cap maximum run duration, periodic health evaluations that catch metric deviations, and synchronous injector neutralization via `terminate()`. By clearing injected latency and error flags synchronously upon abort, the system returns to baseline operation without lingering residual effects.

## Demo Kasus Pembayaran: Baseline dan Mitigasi Sukses

The interactive demonstration simulates a payment processing pipeline across four operational phases. In Phase 1, baseline steady-state traffic dispatches five sequential requests under healthy downstream conditions, yielding 100% primary success with the breaker remaining Closed. In Phase 2, a chaos experiment injects forced downstream errors for 300 ms while routing traffic through a resilient client. The first two failures are routed to the fallback while the breaker is Closed; the third failure trips the breaker to Open, causing all subsequent calls to fast-fail into the fallback queue. The aggregate error rate remains 0.00%.

```text
=== Chaos Engineering & Fault Injection Demo ===

[1] Baseline Steady-State Traffic (Normal Operation)...
Req #1 -> Payment Success (Primary) | CB State: CLOSED
Req #2 -> Payment Success (Primary) | CB State: CLOSED
Req #3 -> Payment Success (Primary) | CB State: CLOSED
Req #4 -> Payment Success (Primary) | CB State: CLOSED
Req #5 -> Payment Success (Primary) | CB State: CLOSED

[2] Running Chaos Experiment: Injecting Downstream Failures with Resilient Client...
Req #6 -> Payment Queued (Fallback) | CB State: CLOSED
Req #7 -> Payment Queued (Fallback) | CB State: CLOSED
Req #8 -> Payment Queued (Fallback) | CB State: OPEN
Req #9 -> Payment Queued (Fallback) | CB State: OPEN
...
Req #15 -> Payment Queued (Fallback) | CB State: OPEN

Experiment State: COMPLETED
Steady-State Metrics with Fallback: Total=15, Success=15, Failed=0, ErrorRate=0.00%
```

## Demo Fase 3 & 4: Auto-Abort dan Pemulihan Mandiri

In Phase 3, the demo executes an experiment against an unmitigated client without fallback handling to verify blast radius auto-abort. Two successive downstream failures push the error rate to 33.33%, breaching the 20% tolerance ceiling. The runner detects the breach, aborts the experiment, and neutralizes the injector immediately. In Phase 4, subsequent traffic dispatches after the cooldown period. The circuit breaker probes downstream in Half-Open, succeeds, resets to Closed, and restores normal primary execution.

```text
[3] Running Second Chaos Experiment: Unmitigated Downstream Fault (Triggers Auto-Abort)...
Raw Call #1 Failed: chaos: injected fault failure
Raw Call #2 Failed: chaos: injected fault failure
Experiment State: ABORTED
Abort Reason: Steady state breached: error rate 33.33% exceeded threshold
Fault Injector Active: false (Neutralized on Abort)

[4] Traffic After Experiment Ends / Recovery...
Req #16 -> Payment Success (Primary) | CB State: CLOSED
Req #17 -> Payment Success (Primary) | CB State: CLOSED
Req #18 -> Payment Success (Primary) | CB State: CLOSED
Req #19 -> Payment Success (Primary) | CB State: CLOSED
Req #20 -> Payment Success (Primary) | CB State: CLOSED
```

## Verifikasi Pengujian: Lima Test Case Lolos 100%

The test harness in `tests/chaos_test.go` verifies every architectural invariant deterministically. Unit tests cover fault injection parameters, circuit breaker state machine transitions, graceful fallback delegation, and experiment auto-abort upon steady-state breach. In addition, `TestConcurrencyAndRace` executes 50 concurrent goroutines continuously reading and mutating states under Go's race detector. All tests pass with zero data races.

```sh
cd labs/34-chaos-engineering
go test -v ./...
go test -race ./...
```

The test execution output proves exact state compliance:

```text
=== RUN   TestFaultInjector
--- PASS: TestFaultInjector (0.01s)
=== RUN   TestCircuitBreakerStateTransitions
--- PASS: TestCircuitBreakerStateTransitions (0.06s)
=== RUN   TestCircuitBreakerGracefulDegradation
--- PASS: TestCircuitBreakerGracefulDegradation (0.00s)
=== RUN   TestExperimentAutoAbortOnSteadyStateViolation
--- PASS: TestExperimentAutoAbortOnSteadyStateViolation (0.01s)
=== RUN   TestConcurrencyAndRace
--- PASS: TestConcurrencyAndRace (0.00s)
PASS
ok  	labs/34-chaos-engineering/tests	1.481s
```

## Desain Konkurensi: Harmonisasi RWMutex dan Mutex

Concurrency control employs tailored primitives matched to workload access patterns. The fault injector uses `sync.RWMutex` to facilitate high-throughput concurrent reads across simulated execution paths, acquiring exclusive locks only when toggling faults. Conversely, the circuit breaker utilizes standard `sync.Mutex` to serialize state mutations and cooldown evaluations cleanly. Meanwhile, the steady-state monitor avoids locks entirely by executing lock-free 64-bit atomic operations, eliminating contention hotspots along critical request paths.

## Perbandingan: Sistem Naif vs Sistem Berketahanan

Comparing unmitigated architectures against resilient designs illustrates how chaos engineering validates defensive patterns:

| Evaluated Dimension | Naive Unmitigated System | Resilient Client Architecture (Lab 34) |
| --- | --- | --- |
| Downstream Latency | Threads hang, socket pools exhaust | Fast-fail rejection via Circuit Breaker (<1ms) |
| Downstream Outage | HTTP 500 errors returned to customers | Graceful degradation via asynchronous fallback |
| User-Facing Error Rate | Escalates directly toward 100% | Maintained at 0.00% through fallback mitigation |
| Resilience Verification | Passive assumption / production surprises | Proactive, empirical chaos testing |
| Recovery Mechanism | Manual service restarts and human intervention | Automated self-healing via Half-Open probing |

## Trade-off dan Batasan Implementasi Laboratorium

Engineering honesty requires acknowledging the scope and deliberate simplifications of this lab. The monitor tracks lifetime cumulative counters rather than a sliding time window; in production, historical failures must age out to reflect current conditions. In addition, faults are injected in-memory rather than via network-level packet drops (e.g., iptables or eBPF). Finally, the `errorRate` field in `internal/fault/injector.go` remains unused in favor of deterministic boolean forcing (`forceError`), avoiding non-deterministic statistical flapping during test runs.

## Jalur Evolusi Menuju Skala Produksi Enterprise

Evolving this model into an enterprise chaos platform involves replacing cumulative counters with sliding ring-buffers that evaluate rolling 60-second windows. Fault injection should migrate from application code into the service mesh layer using Envoy sidecars or Chaos Mesh, allowing chaos injection without code modification. Furthermore, traffic routing should leverage canary deployments, directing turbulence toward a constrained 1% user slice tagged via distributed tracing headers.

## Jebakan Umum dalam Mengadopsi Chaos Engineering

Adopting chaos engineering without disciplined guardrails creates severe operational risks. The most dangerous anti-pattern is executing experiments without an automated kill switch, forcing engineers into manual remediation during unexpected outages. Another frequent error is tracking internal hardware metrics rather than business transactions, mistaking low CPU load for system health. Finally, implementing circuit breakers without fallback mechanisms merely accelerates failure delivery to end users rather than preserving service availability.

## Panduan Keputusan: Kapan Menggunakan Chaos Engineering

Chaos Engineering delivers immense value in microservice architectures with critical third-party dependencies, high transaction volumes, and strict availability SLAs. However, teams should avoid chaos experiments on monolithic systems lacking baseline observability, or on codebases where fundamental error handling and request timeouts have not yet been implemented. Chaos testing is designed to reveal subtle systemic vulnerabilities in resilient architectures, not to confirm obvious defects in unhardened code.

## Checklist Kesiapan Pengujian Chaos di Lingkungan Nyata

Before executing chaos experiments in production environments, verify the following operational gates:

- [x] Define steady-state health using customer-facing business metrics rather than infrastructure stats.
- [x] Establish a precise, falsifiable hypothesis and define acceptable error rate thresholds.
- [x] Configure explicit temporal deadlines and automated kill-switch triggers.
- [x] Verify that resilient clients intercept downstream errors and route to graceful fallbacks.
- [x] Confirm that experiment aborts immediately neutralize all injected faults synchronously.
- [x] Validate thread safety and state transitions with race-detection tooling.

## Rangkuman Esensial: Ketahanan Terverifikasi Secara Empiris

Resilience is not an architectural trait achieved by declaration; it is an empirical outcome validated through disciplined experimentation. Circuit breakers prevent downstream failures from cascading upstream, while graceful fallbacks protect user transactions from immediate degradation. When coupled with automated steady-state monitoring and strict blast radius auto-aborts, chaos engineering empowers software engineering teams to uncover systemic weaknesses safely before production outages occur.

## Source

[Software Engineering Lab 34 - Chaos Engineering & Fault Injection](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/34-chaos-engineering), snapshot commit `cfba3b76`. Source: README, implementation (`internal/fault`, `internal/circuitbreaker`, `internal/monitor`, `internal/experiment`), tests, and execution records. Code excerpts belong in their original package context.
