# Code Snippets

## Snippet 1 — Fault Injector Primitive

Source File: `labs/34-chaos-engineering/internal/fault/injector.go`
Purpose: Menyisipkan latensi dan error paksa secara thread-safe pada pemanggilan layanan downstream.

```go
package fault

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrInjectedFault = errors.New("chaos: injected fault failure")

type Injector struct {
	mu          sync.RWMutex
	enabled     bool
	latency     time.Duration
	errorRate   float64 // 0.0 to 1.0
	forceError  bool
}

func NewInjector() *Injector {
	return &Injector{}
}

func (fi *Injector) SetFault(latency time.Duration, forceError bool) {
	fi.mu.Lock()
	defer fi.mu.Unlock()
	fi.enabled = true
	fi.latency = latency
	fi.forceError = forceError
}

func (fi *Injector) Clear() {
	fi.mu.Lock()
	defer fi.mu.Unlock()
	fi.enabled = false
	fi.latency = 0
	fi.forceError = false
}

func (fi *Injector) IsEnabled() bool {
	fi.mu.RLock()
	defer fi.mu.RUnlock()
	return fi.enabled
}

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

Explanation:
Injector menggunakan `sync.RWMutex` untuk sinkronisasi thread-safe. Jika aktif (`enabled = true`), injector menahan eksekusi selama durasi `latency` (dengan dukungan pembatalan konteks `ctx.Done()`) dan mengembalikan `ErrInjectedFault` jika `forceError` bernilai true.

---

## Snippet 2 — Circuit Breaker State Machine

Source File: `labs/34-chaos-engineering/internal/circuitbreaker/circuitbreaker.go`
Purpose: Mengelola transisi state Circuit Breaker (`Closed`, `Open`, `Half-Open`) dan mengeksekusi fallback saat terjadi kegagalan.

```go
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	lastStateChg time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        StateClosed,
		threshold:    threshold,
		cooldown:     cooldown,
		lastStateChg: time.Now(),
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.checkStateLocked()
	return cb.state
}

func (cb *CircuitBreaker) checkStateLocked() {
	if cb.state == StateOpen && time.Since(cb.lastStateChg) > cb.cooldown {
		cb.state = StateHalfOpen
		cb.lastStateChg = time.Now()
	}
}

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

	// Success
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

Explanation:
Circuit Breaker melacak kegagalan berturut-turut. Jika mencapai `threshold`, state berubah ke `Open`. Setelah waktu `cooldown` terlampaui, state bergeser ke `Half-Open` untuk menguji pemulihan. Jika ada fungsi `fallback`, error ditelan dan fallback dieksekusi.

---

## Snippet 3 — Steady-State Monitor

Source File: `labs/34-chaos-engineering/internal/monitor/monitor.go`
Purpose: Mengumpulkan metrik steady-state secara atomik dan mengevaluasi kesehatan sistem terhadap ambang batas error rate.

```go
package monitor

import (
	"sync"
	"sync/atomic"
)

type SteadyStateMetrics struct {
	TotalRequests   uint64
	FailedRequests  uint64
	SuccessRequests uint64
}

type Monitor struct {
	mu            sync.RWMutex
	totalRequests uint64
	failedRequests uint64
	maxErrorRate  float64 // Threshold e.g. 0.20 (20%)
}

func NewMonitor(maxErrorRate float64) *Monitor {
	return &Monitor{
		maxErrorRate: maxErrorRate,
	}
}

func (m *Monitor) RecordSuccess() {
	atomic.AddUint64(&m.totalRequests, 1)
}

func (m *Monitor) RecordFailure() {
	atomic.AddUint64(&m.totalRequests, 1)
	atomic.AddUint64(&m.failedRequests, 1)
}

func (m *Monitor) Metrics() SteadyStateMetrics {
	total := atomic.LoadUint64(&m.totalRequests)
	failed := atomic.LoadUint64(&m.failedRequests)
	return SteadyStateMetrics{
		TotalRequests:   total,
		FailedRequests:  failed,
		SuccessRequests: total - failed,
	}
}

func (m *Monitor) ErrorRate() float64 {
	total := atomic.LoadUint64(&m.totalRequests)
	if total == 0 {
		return 0.0
	}
	failed := atomic.LoadUint64(&m.failedRequests)
	return float64(failed) / float64(total)
}

func (m *Monitor) IsHealthy() bool {
	total := atomic.LoadUint64(&m.totalRequests)
	if total < 5 { // Minimum sample before evaluating breach
		return true
	}
	return m.ErrorRate() <= m.maxErrorRate
}
```

Explanation:
Menggunakan `sync/atomic` untuk performa tinggi dan bebas race condition saat mencatat sukses dan gagal. Minimum 5 sampel diperlukan sebelum evaluasi kesehatan dipicu.

---

## Snippet 4 — Experiment Runner & Auto-Abort

Source File: `labs/34-chaos-engineering/internal/experiment/runner.go`
Purpose: Mengelola siklus eksperimen chaos, pemantauan berkala, dan pembatalan otomatis (*auto-abort*) beserta netralisasi fault injector.

```go
package experiment

import (
	"context"
	"fmt"
	"sync"
	"time"

	"labs/34-chaos-engineering/internal/fault"
	"labs/34-chaos-engineering/internal/monitor"
)

type ExperimentState string

const (
	StatePending   ExperimentState = "PENDING"
	StateRunning   ExperimentState = "RUNNING"
	StateAborted   ExperimentState = "ABORTED"
	StateCompleted ExperimentState = "COMPLETED"
)

type Config struct {
	Name             string
	Duration         time.Duration
	Latency          time.Duration
	ForceError       bool
	MonitorInterval  time.Duration
}

type Experiment struct {
	mu          sync.Mutex
	config      Config
	injector    *fault.Injector
	monitor     *monitor.Monitor
	state       ExperimentState
	abortReason string
}

func NewExperiment(cfg Config, inj *fault.Injector, mon *monitor.Monitor) *Experiment {
	return &Experiment{
		config:   cfg,
		injector: inj,
		monitor:  mon,
		state:    StatePending,
	}
}

func (e *Experiment) State() ExperimentState {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state
}

func (e *Experiment) AbortReason() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.abortReason
}

func (e *Experiment) Run(ctx context.Context) error {
	e.mu.Lock()
	e.state = StateRunning
	e.mu.Unlock()

	// Inject fault
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

func (e *Experiment) terminate(targetState ExperimentState, reason string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.injector.Clear()
	e.state = targetState
	e.abortReason = reason
}
```

Explanation:
Runner menjalankan ticker pemantauan kesehatan. Jika `monitor.IsHealthy()` mengembalikan `false`, eksperimen segera di-abort dan `injector.Clear()` dipanggil secara sinkron untuk menghentikan gangguan seketika.
