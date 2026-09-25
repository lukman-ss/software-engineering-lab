package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State string

const (
	StateClosed   State = "CLOSED"
	StateOpen     State = "OPEN"
	StateHalfOpen State = "HALF_OPEN"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type Config struct {
	FailureThreshold int
	OpenTimeout      time.Duration
	HalfOpenMaxCalls int
}

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	consecutiveSuccesses int
	halfOpenCalls    int
	lastStateChange  time.Time
	config           Config
	now              func() time.Time
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

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.checkStateTransitionLocked()
	return cb.state
}

func (cb *CircuitBreaker) checkStateTransitionLocked() {
	if cb.state == StateOpen && cb.now().Sub(cb.lastStateChange) >= cb.config.OpenTimeout {
		cb.state = StateHalfOpen
		cb.halfOpenCalls = 0
		cb.consecutiveSuccesses = 0
		cb.lastStateChange = cb.now()
	}
}

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
