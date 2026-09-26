package slo

import (
	"math"
	"time"

	"labs/24-slo-sli-error-budget/internal/metrics"
)

type Config struct {
	Name            string
	TargetUptime    float64
	LatencyThreshold time.Duration
}

type Evaluator struct {
	config  Config
	tracker *metrics.WindowTracker
}

type Status struct {
	SLOName          string
	TargetSLO        float64
	CurrentSLI       float64
	TotalEvents      int64
	GoodEvents       int64
	BadEvents        int64
	TotalErrorBudget float64
	BudgetConsumed   float64
	BudgetRemaining  float64
	CanDeploy        bool
}

func NewEvaluator(cfg Config, tracker *metrics.WindowTracker) *Evaluator {
	return &Evaluator{
		config:  cfg,
		tracker: tracker,
	}
}

func (e *Evaluator) Evaluate(now time.Time) Status {
	total, good, bad := e.tracker.Summary(now)

	var sli float64 = 1.0
	if total > 0 {
		sli = float64(good) / float64(total)
	}

	allowedFailureRate := 1.0 - e.config.TargetUptime
	totalErrorBudget := allowedFailureRate * float64(total)
	budgetConsumed := float64(bad)
	budgetRemaining := totalErrorBudget - budgetConsumed

	canDeploy := true
	if total > 0 && budgetRemaining <= 0 {
		canDeploy = false
	}

	return Status{
		SLOName:          e.config.Name,
		TargetSLO:        e.config.TargetUptime,
		CurrentSLI:       math.Round(sli*10000) / 10000,
		TotalEvents:      total,
		GoodEvents:       good,
		BadEvents:        bad,
		TotalErrorBudget: math.Round(totalErrorBudget*100) / 100,
		BudgetConsumed:   budgetConsumed,
		BudgetRemaining:  math.Round(budgetRemaining*100) / 100,
		CanDeploy:        canDeploy,
	}
}
