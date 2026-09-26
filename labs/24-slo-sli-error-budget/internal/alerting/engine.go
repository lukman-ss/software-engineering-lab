package alerting

import (
	"time"

	"labs/24-slo-sli-error-budget/internal/metrics"
)

type Severity string

const (
	SeverityPage   Severity = "PAGE"
	SeverityTicket Severity = "TICKET"
	SeverityNone   Severity = "NONE"
)

type BurnRateRule struct {
	Name             string
	Severity         Severity
	LongWindow       time.Duration
	ShortWindow      time.Duration
	BurnRateFactor   float64
	BudgetConsumedPct float64
}

type AlertEngine struct {
	targetSLO    float64
	shortTracker *metrics.WindowTracker
	longTracker  *metrics.WindowTracker
	rules        []BurnRateRule
}

type AlertResult struct {
	Triggered       bool
	RuleName        string
	Severity        Severity
	LongBurnRate    float64
	ShortBurnRate   float64
	ThresholdRate   float64
}

func NewAlertEngine(targetSLO float64, shortTracker, longTracker *metrics.WindowTracker, rules []BurnRateRule) *AlertEngine {
	return &AlertEngine{
		targetSLO:    targetSLO,
		shortTracker: shortTracker,
		longTracker:  longTracker,
		rules:        rules,
	}
}

func (a *AlertEngine) CalculateBurnRate(total, bad int64) float64 {
	if total == 0 {
		return 0.0
	}
	actualErrorRate := float64(bad) / float64(total)
	allowedErrorRate := 1.0 - a.targetSLO
	if allowedErrorRate <= 0 {
		return 0.0
	}
	return actualErrorRate / allowedErrorRate
}

func (a *AlertEngine) Check(now time.Time) []AlertResult {
	shortTotal, _, shortBad := a.shortTracker.Summary(now)
	longTotal, _, longBad := a.longTracker.Summary(now)

	shortBurn := a.CalculateBurnRate(shortTotal, shortBad)
	longBurn := a.CalculateBurnRate(longTotal, longBad)

	var results []AlertResult
	for _, rule := range a.rules {
		triggered := false
		if shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor {
			triggered = true
		}

		if triggered {
			results = append(results, AlertResult{
				Triggered:     true,
				RuleName:      rule.Name,
				Severity:      rule.Severity,
				LongBurnRate:  longBurn,
				ShortBurnRate: shortBurn,
				ThresholdRate: rule.BurnRateFactor,
			})
		}
	}
	return results
}
