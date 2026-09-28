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
