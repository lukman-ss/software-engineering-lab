## Snippet 1 — Circuit Breaker State Constants

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Defines the three operational states of the circuit breaker state machine.

```go
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Exported aliases for external compatibility
const (
	StateClosed   = Closed
	StateOpen     = Open
	StateHalfOpen = HalfOpen
)
```

## Snippet 2 — Circuit Breaker Configuration

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Configuration struct, default values, and constructor for the circuit breaker.

```go
type Config struct {
	FailureThreshold int
	OpenTimeout      time.Duration
	HalfOpenMaxCalls int
}

func DefaultConfig() Config {
	return Config{
		FailureThreshold: 3,
		OpenTimeout:      300 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
}

func New(cfg Config) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 3
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 300 * time.Millisecond
	}
	if cfg.HalfOpenMaxCalls <= 0 {
		cfg.HalfOpenMaxCalls = 1
	}

	return &Breaker{cfg: cfg, state: Closed}
}
```

## Snippet 3 — State Transition Check

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Checks if OPEN circuit should transition to HALF_OPEN based on elapsed time.

```go
func (b *Breaker) advanceLocked(now time.Time) {
	if b.state == Open && now.Sub(b.openedAt) >= b.cfg.OpenTimeout {
		b.state = HalfOpen
		b.halfOpenIn = 0
		b.generation++
	}
}
```

## Snippet 4 — Execute Method (Full State Machine)

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Core execution method implementing the three-state machine with fail-fast, probe limiting, and panic handling.

```go
func (b *Breaker) Execute(fn func() error) (err error) {
	b.mu.Lock()
	b.advanceLocked(time.Now())
	gen := b.generation
	switch b.state {
	case Open:
		b.mu.Unlock()
		return ErrCircuitOpen
	case HalfOpen:
		if b.halfOpenIn >= b.cfg.HalfOpenMaxCalls {
			b.mu.Unlock()
			return ErrCircuitOpen
		}
		b.halfOpenIn++
	case Closed:
	}
	b.mu.Unlock()

	panicked := true
	defer func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if panicked {
			b.onFailureLocked(gen, time.Now())
		} else if err == nil {
			b.onSuccessLocked(gen)
		} else {
			b.onFailureLocked(gen, time.Now())
		}
	}()

	err = fn()
	panicked = false
	return err
}
```

## Snippet 5 — Checkout Service with Circuit Breaker

Source File: `internal/checkout/service.go`

Purpose: Application-level consumer that routes payment calls through the circuit breaker.

```go
func (s *Service) Checkout(ctx context.Context) error {
	return s.breaker.Execute(func() error {
		return s.payment.ProcessPayment(ctx)
	})
}
```

## Snippet 6 — Demo Circuit Breaker Configuration

Source File: `cmd/demo/main.go`

Purpose: Shows the configuration used in the demo scenarios.

```go
cb := circuitbreaker.New(circuitbreaker.Config{
	FailureThreshold: 3,
	OpenTimeout:      300 * time.Millisecond,
	HalfOpenMaxCalls: 1,
})
```

## Snippet 7 — Fake Server Mode Switching

Source File: `internal/payment/fake_server.go`

Purpose: Simulates a controllable downstream dependency for testing.

```go
func NewFakeServer(slowDelay time.Duration) *FakeServer {
	fs := &FakeServer{
		slowDelay: slowDelay,
	}
	fs.mode.Store(ModeHealthy)

	mux := http.NewServeMux()
	mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		fs.requestCount.Add(1)
		mode := fs.mode.Load().(ServerMode)

		switch mode {
		case ModeSlow:
			time.Sleep(fs.slowDelay)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok_slow"}`))
		case ModeDown:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal payment server failure"))
		default: // ModeHealthy
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	})

	fs.server = httptest.NewServer(mux)
	return fs
}
```
