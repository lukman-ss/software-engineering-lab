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
