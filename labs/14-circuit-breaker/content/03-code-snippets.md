## Snippet 1 — Circuit Breaker State Constants

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Defines the three operational states of the circuit breaker state machine.

```go
type State string

const (
	StateClosed   State = "CLOSED"
	StateOpen     State = "OPEN"
	StateHalfOpen State = "HALF_OPEN"
)
```

## Snippet 2 — Circuit Breaker Configuration

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Configuration struct and default values for the circuit breaker.

```go
type Config struct {
	FailureThreshold int
	OpenTimeout      time.Duration
	HalfOpenMaxCalls int
}

func New(cfg Config) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 3
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 5 * time.Second
	}
	if cfg.HalfOpenMaxCalls <= 0 {
		cfg.HalfOpenMaxCalls = 1
	}

	return &CircuitBreaker{
		state:           StateClosed,
		config:          cfg,
		lastStateChange: time.Now(),
		now:             time.Now,
	}
}
```

## Snippet 3 — State Transition Check

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Checks if OPEN circuit should transition to HALF_OPEN based on elapsed time.

```go
func (cb *CircuitBreaker) checkStateTransitionLocked() {
	if cb.state == StateOpen && cb.now().Sub(cb.lastStateChange) >= cb.config.OpenTimeout {
		cb.state = StateHalfOpen
		cb.halfOpenCalls = 0
		cb.consecutiveSuccesses = 0
		cb.lastStateChange = cb.now()
	}
}
```

## Snippet 4 — Execute Method (Full State Machine)

Source File: `internal/circuitbreaker/circuit_breaker.go`

Purpose: Core execution method implementing the three-state machine with fail-fast, probe limiting, and panic handling.

```go
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	cb.checkStateTransitionLocked()

	switch cb.state {
	case StateOpen:
		cb.mu.Unlock()
		return ErrCircuitOpen

	case StateHalfOpen:
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		cb.halfOpenCalls++
		cb.mu.Unlock()

		defer func() {
			if r := recover(); r != nil {
				cb.mu.Lock()
				defer cb.mu.Unlock()
				cb.state = StateOpen
				cb.lastStateChange = cb.now()
				cb.failureCount = 0
				cb.halfOpenCalls = 0
				panic(r)
			}
		}()

		err := fn()

		cb.mu.Lock()
		defer cb.mu.Unlock()
		if err != nil {
			cb.state = StateOpen
			cb.lastStateChange = cb.now()
			cb.failureCount = 0
			cb.halfOpenCalls = 0
			return err
		}

		cb.consecutiveSuccesses++
		if cb.consecutiveSuccesses >= cb.config.HalfOpenMaxCalls {
			cb.state = StateClosed
			cb.lastStateChange = cb.now()
			cb.failureCount = 0
			cb.halfOpenCalls = 0
		}
		return nil

	default: // StateClosed
		cb.mu.Unlock()

		defer func() {
			if r := recover(); r != nil {
				cb.mu.Lock()
				defer cb.mu.Unlock()
				cb.failureCount++
				if cb.failureCount >= cb.config.FailureThreshold {
					cb.state = StateOpen
					cb.lastStateChange = cb.now()
				}
				panic(r)
			}
		}()

		err := fn()

		cb.mu.Lock()
		defer cb.mu.Unlock()
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.config.FailureThreshold {
				cb.state = StateOpen
				cb.lastStateChange = cb.now()
			}
			return err
		}

		cb.failureCount = 0
		return nil
	}
}
```

## Snippet 5 — Checkout Service with Circuit Breaker

Source File: `internal/checkout/service.go`

Purpose: Application-level consumer that routes payment calls through the circuit breaker.

```go
func (s *Service) Checkout(ctx context.Context) error {
	if s.cb != nil {
		err := s.cb.Execute(func() error {
			return s.paymentClient.ProcessPayment(ctx)
		})
		if err != nil {
			return fmt.Errorf("checkout payment failed (with CB): %w", err)
		}
		return nil
	}

	// Without Circuit Breaker
	err := s.paymentClient.ProcessPayment(ctx)
	if err != nil {
		return fmt.Errorf("checkout payment failed (no CB): %w", err)
	}
	return nil
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
			http.Error(w, "internal payment server failure", http.StatusInternalServerError)
		default: // ModeHealthy
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	})

	fs.server = httptest.NewServer(mux)
	return fs
}
```
