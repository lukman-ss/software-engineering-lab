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
