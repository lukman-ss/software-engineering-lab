// Package breaker implements a Circuit Breaker pattern.
// A Circuit Breaker prevents an application from repeatedly trying to execute an operation
// that is likely to fail, allowing it to continue without waiting for the fault to be fixed
// or wasting CPU cycles.
package breaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	// StateClosed means the circuit is closed and requests flow normally.
	StateClosed State = iota
	// StateOpen means the circuit is open and requests are rejected immediately.
	StateOpen
	// StateHalfOpen means the circuit is testing if the downstream service has recovered.
	StateHalfOpen
)

// String returns the string representation of the state.
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

var (
	// ErrCircuitOpen is returned when the circuit breaker is in the Open state.
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyProbes is returned when the HalfOpen probe limit is reached.
	ErrTooManyProbes = errors.New("too many half-open probes")
)

// Config holds configuration for a CircuitBreaker.
type Config struct {
	// FailureThreshold is the number of consecutive failures before opening the circuit.
	FailureThreshold int
	// ResetTimeout is the duration to wait before transitioning from Open to HalfOpen.
	ResetTimeout time.Duration
	// MaxHalfOpenProbes limits concurrent requests when in HalfOpen state.
	MaxHalfOpenProbes int
	// Now allows injecting a custom time function for testing.
	Now func() time.Time
}

// CircuitBreaker is a state machine that monitors for failures and prevents
// cascading failures by failing fast when a downstream service is struggling.
// It is fully thread-safe (Prompt 050).
type CircuitBreaker struct {
	mu sync.RWMutex

	cfg Config

	state         State
	consecFailures int
	openedAt      time.Time

	// activeProbes tracks concurrent requests allowed during Half-Open state (Prompt 051)
	activeProbes  int
}

// New returns a new CircuitBreaker.
func New(cfg Config) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.ResetTimeout <= 0 {
		cfg.ResetTimeout = 60 * time.Second
	}
	if cfg.MaxHalfOpenProbes <= 0 {
		cfg.MaxHalfOpenProbes = 1 // Default to 1 probe allowed
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}

	return &CircuitBreaker{
		cfg:   cfg,
		state: StateClosed,
	}
}

// State returns the current state of the circuit breaker.
// It may perform a state transition from Open to HalfOpen if the reset timeout has elapsed.
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.checkStateTransition()
	return cb.state
}

// checkStateTransition verifies if an Open breaker should transition to HalfOpen.
// It assumes the caller holds at least a read lock, but since it modifies state,
// the caller must actually hold a write lock.
func (cb *CircuitBreaker) checkStateTransition() {
	if cb.state == StateOpen {
		if cb.cfg.Now().Sub(cb.openedAt) >= cb.cfg.ResetTimeout {
			cb.state = StateHalfOpen
			cb.activeProbes = 0
		}
	}
}

// Execute runs the given function if the circuit is Closed or HalfOpen.
// If the circuit is Open, it returns ErrCircuitOpen immediately.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if err := cb.beforeExecution(); err != nil {
		return err
	}

	// Execute the actual function
	err := fn()

	cb.afterExecution(err)
	return err
}

// beforeExecution checks if the request is allowed to proceed.
func (cb *CircuitBreaker) beforeExecution() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.checkStateTransition()

	switch cb.state {
	case StateOpen:
		return ErrCircuitOpen
	case StateHalfOpen:
		// Prompt 051 - Half-Open Probe Limiting
		// Prevent thundering herd when cooling down ends
		if cb.activeProbes >= cb.cfg.MaxHalfOpenProbes {
			return ErrTooManyProbes
		}
		cb.activeProbes++
		return nil
	default:
		return nil
	}
}

// afterExecution updates the circuit breaker state based on the result.
func (cb *CircuitBreaker) afterExecution(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	isSuccess := err == nil

	if cb.state == StateHalfOpen {
		// Probe finished
		cb.activeProbes--

		if isSuccess {
			// One success in HalfOpen closes the circuit
			cb.state = StateClosed
			cb.consecFailures = 0
		} else {
			// One failure in HalfOpen immediately trips back to Open
			cb.state = StateOpen
			cb.openedAt = cb.cfg.Now()
		}
		return
	}

	if isSuccess {
		cb.consecFailures = 0
	} else {
		cb.consecFailures++
		if cb.consecFailures >= cb.cfg.FailureThreshold {
			cb.state = StateOpen
			cb.openedAt = cb.cfg.Now()
		}
	}
}
