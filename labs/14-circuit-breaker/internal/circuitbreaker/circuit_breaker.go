package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Exported constants for integration test compatibility
const (
	StateClosed = Closed
	StateOpen   = Open
	StateHalfOpen = HalfOpen
)

func (s State) String() string {
	switch s {
	case Closed:
		return "CLOSED"
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

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

type Breaker struct {
	mu         sync.Mutex
	cfg        Config
	state      State
	failures   int
	openedAt   time.Time
	halfOpenIn int
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

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.advanceLocked(time.Now())
	return b.state
}

func (b *Breaker) advanceLocked(now time.Time) {
	if b.state == Open && now.Sub(b.openedAt) >= b.cfg.OpenTimeout {
		b.state = HalfOpen
		b.halfOpenIn = 0
	}
}

func (b *Breaker) Execute(fn func() error) error {
	b.mu.Lock()
	b.advanceLocked(time.Now())
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

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()
	if err == nil {
		b.onSuccessLocked()
		return nil
	}
	b.onFailureLocked(time.Now())
	return err
}

func (b *Breaker) onSuccessLocked() {
	if b.state == HalfOpen {
		b.state = Closed
		b.failures = 0
		b.halfOpenIn = 0
		return
	}
	b.failures = 0
}

func (b *Breaker) onFailureLocked(now time.Time) {
	if b.state == HalfOpen {
		b.state = Open
		b.openedAt = now
		b.halfOpenIn = 0
		return
	}
	b.failures++
	if b.failures >= b.cfg.FailureThreshold {
		b.state = Open
		b.openedAt = now
	}
}
