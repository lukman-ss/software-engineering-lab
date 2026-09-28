package saga

import (
	"context"
	"fmt"
	"sync"
)

type StepStatus string

const (
	StatusPending         StepStatus = "PENDING"
	StatusExecuted        StepStatus = "EXECUTED"
	StatusFailed          StepStatus = "FAILED"
	StatusCompensated     StepStatus = "COMPENSATED"
	StatusCompensateFailed StepStatus = "COMPENSATE_FAILED"
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
		select {
		case <-ctx.Done():
			o.mu.Lock()
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()
			compErr := o.compensate(context.Background(), executed)
			if compErr != nil {
				return fmt.Errorf("saga cancelled: %w; compensation errors: %v", ctx.Err(), compErr)
			}
			return fmt.Errorf("saga cancelled: %w", ctx.Err())
		default:
		}

		err := step.Execute(ctx)
		o.mu.Lock()
		if err != nil {
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()

			// ponytail: LIFO rollback algorithm; compensations run, errors logged/aggregated
			compErr := o.compensate(context.Background(), executed)
			if compErr != nil {
				return fmt.Errorf("step %s failed: %w; compensation errors: %v", step.Name, err, compErr)
			}
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}
		o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusExecuted})
		o.mu.Unlock()
		executed = append(executed, step)
	}

	return nil
}

func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {
	var compErrors []error
	for i := len(executed) - 1; i >= 0; i-- {
		step := executed[i]
		if step.Compensate != nil {
			err := step.Compensate(ctx)
			o.mu.Lock()
			if err != nil {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensateFailed})
				compErrors = append(compErrors, fmt.Errorf("compensation %s: %w", step.Name, err))
			} else {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensated})
			}
			o.mu.Unlock()
		}
	}
	if len(compErrors) > 0 {
		return fmt.Errorf("%v", compErrors)
	}
	return nil
}

func (o *Orchestrator) Logs() []StepLog {
	o.mu.Lock()
	defer o.mu.Unlock()
	res := make([]StepLog, len(o.logs))
	copy(res, o.logs)
	return res
}
