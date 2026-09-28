package saga

import (
	"context"
	"fmt"
	"sync"
)

type StepStatus string

const (
	StatusPending   StepStatus = "PENDING"
	StatusExecuted  StepStatus = "EXECUTED"
	StatusFailed    StepStatus = "FAILED"
	StatusCompensated StepStatus = "COMPENSATED"
)

type Step struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type StepLog struct {
	Name   string
	Status StepStatus
}

type Orchestrator struct {
	mu    sync.Mutex
	steps []Step
	logs  []StepLog
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		steps: make([]Step, 0),
		logs:  make([]StepLog, 0),
	}
}

func (o *Orchestrator) AddStep(step Step) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.steps = append(o.steps, step)
}

func (o *Orchestrator) Execute(ctx context.Context) error {
	o.mu.Lock()
	steps := make([]Step, len(o.steps))
	copy(steps, o.steps)
	o.mu.Unlock()

	executed := make([]Step, 0)

	for _, step := range steps {
		err := step.Execute(ctx)
		o.mu.Lock()
		if err != nil {
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()

			// ponytail: LIFO rollback algorithm ceiling; assumes compensations succeed without retries
			o.compensate(ctx, executed)
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}
		o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusExecuted})
		o.mu.Unlock()
		executed = append(executed, step)
	}

	return nil
}

func (o *Orchestrator) compensate(ctx context.Context, executed []Step) {
	for i := len(executed) - 1; i >= 0; i-- {
		step := executed[i]
		if step.Compensate != nil {
			_ = step.Compensate(ctx)
			o.mu.Lock()
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensated})
			o.mu.Unlock()
		}
	}
}

func (o *Orchestrator) Logs() []StepLog {
	o.mu.Lock()
	defer o.mu.Unlock()
	res := make([]StepLog, len(o.logs))
	copy(res, o.logs)
	return res
}
