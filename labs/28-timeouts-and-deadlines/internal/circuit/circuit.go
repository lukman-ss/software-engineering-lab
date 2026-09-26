package circuit

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
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

type Config struct {
	FailureThreshold int
	SuccessThreshold int
	Cooldown         time.Duration
}

type Breaker struct {
	mu           sync.RWMutex
	cfg          Config
	state        State
	failures     int
	successes    int
	lastStateChg time.Time
}

func NewBreaker(cfg Config) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 3
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = 100 * time.Millisecond
	}
	return &Breaker{
		cfg:          cfg,
		state:        StateClosed,
		lastStateChg: time.Now(),
	}
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.checkCooldown()
	return b.state
}

func (b *Breaker) checkCooldown() {
	if b.state == StateOpen && time.Since(b.lastStateChg) >= b.cfg.Cooldown {
		b.state = StateHalfOpen
		b.successes = 0
		b.failures = 0
		b.lastStateChg = time.Now()
	}
}

func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.checkCooldown()

	if b.state == StateOpen {
		return ErrCircuitOpen
	}
	return nil
}

func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.successes++
		if b.successes >= b.cfg.SuccessThreshold {
			b.state = StateClosed
			b.failures = 0
			b.successes = 0
			b.lastStateChg = time.Now()
		}
	} else if b.state == StateClosed {
		b.failures = 0
	}
}

func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.state = StateOpen
		b.failures = 0
		b.successes = 0
		b.lastStateChg = time.Now()
	} else if b.state == StateClosed {
		b.failures++
		if b.failures >= b.cfg.FailureThreshold {
			b.state = StateOpen
			b.failures = 0
			b.successes = 0
			b.lastStateChg = time.Now()
		}
	}
}

func (b *Breaker) Execute(fn func() error) error {
	if err := b.Allow(); err != nil {
		return err
	}
	err := fn()
	if err != nil {
		b.RecordFailure()
		return err
	}
	b.RecordSuccess()
	return nil
}
