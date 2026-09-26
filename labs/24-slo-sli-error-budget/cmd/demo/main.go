package main

import (
	"fmt"
	"time"

	"labs/24-slo-sli-error-budget/internal/alerting"
	"labs/24-slo-sli-error-budget/internal/metrics"
	"labs/24-slo-sli-error-budget/internal/slo"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("  SLI / SLO / ERROR BUDGET & BURN RATE ALERTING DEMO")
	fmt.Println("================================================================")

	targetSLO := 0.999 // 99.9% availability
	latencyThreshold := 200 * time.Millisecond

	isGood := func(e metrics.Event) bool {
		return e.StatusCode < 500 && e.Duration <= latencyThreshold
	}

	window30d := 30 * time.Minute
	shortWin := 5 * time.Minute
	longWin := 60 * time.Minute

	sloTracker := metrics.NewWindowTracker(window30d, 10*time.Second, isGood)
	shortTracker := metrics.NewWindowTracker(shortWin, 5*time.Second, isGood)
	longTracker := metrics.NewWindowTracker(longWin, 10*time.Second, isGood)

	evaluator := slo.NewEvaluator(slo.Config{
		Name:             "Payment Service Availability",
		TargetUptime:     targetSLO,
		LatencyThreshold: latencyThreshold,
	}, sloTracker)

	alertRules := []alerting.BurnRateRule{
		{
			Name:           "Critical Fast Burn (14.4x - 2% in 1h)",
			Severity:       alerting.SeverityPage,
			BurnRateFactor: 14.4,
		},
		{
			Name:           "Slow Burn Alert (6.0x - 5% in 6h)",
			Severity:       alerting.SeverityTicket,
			BurnRateFactor: 6.0,
		},
	}

	alertEngine := alerting.NewAlertEngine(targetSLO, shortTracker, longTracker, alertRules)

	now := time.Now()

	fmt.Println("\n[PHASE 1] Simulating Baseline Traffic (1,000 requests, 100% success)...")
	for i := 0; i < 1000; i++ {
		ev := metrics.Event{
			Timestamp:  now.Add(time.Duration(i) * 50 * time.Millisecond),
			Duration:   50 * time.Millisecond,
			StatusCode: 200,
			Endpoint:   "/api/v1/pay",
		}
		sloTracker.Record(ev)
		shortTracker.Record(ev)
		longTracker.Record(ev)
	}

	simTime := now.Add(50 * time.Second)
	status := evaluator.Evaluate(simTime)
	fmt.Printf("Total: %d | Good: %d | Bad: %d\n", status.TotalEvents, status.GoodEvents, status.BadEvents)
	fmt.Printf("Target SLO: %.3f%% | Current SLI: %.4f%% | Budget Remaining: %.2f\n", status.TargetSLO*100, status.CurrentSLI*100, status.BudgetRemaining)
	fmt.Printf("Deployment Allowed: %v\n", status.CanDeploy)

	fmt.Println("\n[PHASE 2] Simulating Severe Incident (100 total requests, 10 errors = 10% error rate)...")
	// Target SLO = 99.9% (allowed error rate = 0.1%). With 10% error rate, burn rate = 0.10 / 0.001 = 100x
	incStart := simTime.Add(time.Second)
	for i := 0; i < 90; i++ {
		ev := metrics.Event{
			Timestamp:  incStart.Add(time.Duration(i) * 10 * time.Millisecond),
			Duration:   50 * time.Millisecond,
			StatusCode: 200,
			Endpoint:   "/api/v1/pay",
		}
		sloTracker.Record(ev)
		shortTracker.Record(ev)
		longTracker.Record(ev)
	}
	for i := 0; i < 10; i++ {
		ev := metrics.Event{
			Timestamp:  incStart.Add(time.Duration(90+i) * 10 * time.Millisecond),
			Duration:   500 * time.Millisecond,
			StatusCode: 500,
			Endpoint:   "/api/v1/pay",
		}
		sloTracker.Record(ev)
		shortTracker.Record(ev)
		longTracker.Record(ev)
	}

	evalTime := incStart.Add(2 * time.Second)
	status = evaluator.Evaluate(evalTime)
	fmt.Printf("Total: %d | Good: %d | Bad: %d\n", status.TotalEvents, status.GoodEvents, status.BadEvents)
	fmt.Printf("Target SLO: %.3f%% | Current SLI: %.4f%% | Budget Remaining: %.2f\n", status.TargetSLO*100, status.CurrentSLI*100, status.BudgetRemaining)
	fmt.Printf("Deployment Allowed: %v (Budget exhausted)\n", status.CanDeploy)

	fmt.Println("\n[PHASE 3] Checking Multi-Window Burn Rate Alerts...")
	alerts := alertEngine.Check(evalTime)
	if len(alerts) == 0 {
		fmt.Println("No alerts triggered.")
	} else {
		for _, a := range alerts {
			fmt.Printf(">>> ALERT TRIGGERED: [%s] %s | ShortBurn: %.2fx | LongBurn: %.2fx (Threshold: %.2fx)\n",
				a.Severity, a.RuleName, a.ShortBurnRate, a.LongBurnRate, a.ThresholdRate)
		}
	}

	fmt.Println("\n[PHASE 4] Endpoint Criticality Comparison (Payment 99.9% vs Reports 95.0%)...")
	reportsTracker := metrics.NewWindowTracker(window30d, 10*time.Second, isGood)
	reportsEvaluator := slo.NewEvaluator(slo.Config{
		Name:             "Reports Service Availability",
		TargetUptime:     0.95, // 95% non-critical SLO
		LatencyThreshold: latencyThreshold,
	}, reportsTracker)

	// Ingest identical 10 errors out of 100 requests (10% error rate) to Reports
	for i := 0; i < 90; i++ {
		reportsTracker.Record(metrics.Event{
			Timestamp:  incStart.Add(time.Duration(i) * 10 * time.Millisecond),
			Duration:   50 * time.Millisecond,
			StatusCode: 200,
			Endpoint:   "/api/v1/reports",
		})
	}
	for i := 0; i < 10; i++ {
		reportsTracker.Record(metrics.Event{
			Timestamp:  incStart.Add(time.Duration(90+i) * 10 * time.Millisecond),
			Duration:   500 * time.Millisecond,
			StatusCode: 500,
			Endpoint:   "/api/v1/reports",
		})
	}
	reportStatus := reportsEvaluator.Evaluate(evalTime)
	fmt.Printf("Reports Target SLO: %.1f%% | Current SLI: %.1f%% | Budget Remaining: %.2f\n",
		reportStatus.TargetSLO*100, reportStatus.CurrentSLI*100, reportStatus.BudgetRemaining)
	fmt.Printf("Payment CanDeploy: %v | Reports CanDeploy: %v (Reports has wider 5%% error tolerance)\n",
		status.CanDeploy, reportStatus.CanDeploy)

	fmt.Println("\n================================================================")
	fmt.Println("  DEMO COMPLETE")
	fmt.Println("================================================================")
}
