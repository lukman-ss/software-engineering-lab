package tests

import (
	"sync"
	"testing"
	"time"

	"labs/24-slo-sli-error-budget/internal/alerting"
	"labs/24-slo-sli-error-budget/internal/metrics"
	"labs/24-slo-sli-error-budget/internal/slo"
)

func TestMetricsWindowTracker(t *testing.T) {
	window := 10 * time.Second
	bucketSize := 1 * time.Second
	isGood := func(e metrics.Event) bool {
		return e.StatusCode < 500 && e.Duration <= 100*time.Millisecond
	}

	tracker := metrics.NewWindowTracker(window, bucketSize, isGood)

	now := time.Now()
	// Record 10 good, 2 bad
	for i := 0; i < 10; i++ {
		tracker.Record(metrics.Event{
			Timestamp:  now.Add(time.Duration(i) * 50 * time.Millisecond),
			Duration:   50 * time.Millisecond,
			StatusCode: 200,
		})
	}
	tracker.Record(metrics.Event{
		Timestamp:  now.Add(600 * time.Millisecond),
		Duration:   200 * time.Millisecond, // slow
		StatusCode: 200,
	})
	tracker.Record(metrics.Event{
		Timestamp:  now.Add(700 * time.Millisecond),
		Duration:   50 * time.Millisecond,
		StatusCode: 500, // error
	})

	total, good, bad := tracker.Summary(now.Add(1 * time.Second))
	if total != 12 {
		t.Fatalf("expected total 12, got %d", total)
	}
	if good != 10 {
		t.Fatalf("expected good 10, got %d", good)
	}
	if bad != 2 {
		t.Fatalf("expected bad 2, got %d", bad)
	}

	// Test eviction
	future := now.Add(20 * time.Second)
	total, good, bad = tracker.Summary(future)
	if total != 0 || good != 0 || bad != 0 {
		t.Fatalf("expected 0 events after eviction, got total=%d, good=%d, bad=%d", total, good, bad)
	}
}

func TestSLOEvaluator(t *testing.T) {
	tracker := metrics.NewWindowTracker(1*time.Minute, 1*time.Second, func(e metrics.Event) bool {
		return e.StatusCode == 200
	})
	evaluator := slo.NewEvaluator(slo.Config{
		Name:         "API SLO",
		TargetUptime: 0.99, // 99% allowed failure = 1%
	}, tracker)

	now := time.Now()
	// 99 good, 1 bad -> exactly 99% SLI
	for i := 0; i < 99; i++ {
		tracker.Record(metrics.Event{Timestamp: now, StatusCode: 200})
	}
	tracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})

	status := evaluator.Evaluate(now)
	if status.CurrentSLI < 0.99 {
		t.Fatalf("expected SLI 0.99, got %f", status.CurrentSLI)
	}
	if !status.CanDeploy {
		t.Fatalf("expected CanDeploy=true when remaining budget >= 0, got %v", status.CanDeploy)
	}

	// Another bad event pushes SLI below 99%
	tracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})
	status = evaluator.Evaluate(now)
	if status.CanDeploy {
		t.Fatalf("expected CanDeploy=false when budget exhausted, got %v", status.CanDeploy)
	}
}

func TestAlertEngineBurnRate(t *testing.T) {
	isGood := func(e metrics.Event) bool { return e.StatusCode < 500 }
	shortTracker := metrics.NewWindowTracker(5*time.Minute, time.Second, isGood)
	longTracker := metrics.NewWindowTracker(60*time.Minute, time.Second, isGood)

	rules := []alerting.BurnRateRule{
		{
			Name:           "Page Alert",
			Severity:       alerting.SeverityPage,
			BurnRateFactor: 14.4,
		},
	}

	// Target 99.9% SLO (0.001 allowed error rate)
	engine := alerting.NewAlertEngine(0.999, shortTracker, longTracker, rules)
	now := time.Now()

	// 100 requests with 2 errors = 2% error rate = 20x burn rate (> 14.4x)
	for i := 0; i < 98; i++ {
		ev := metrics.Event{Timestamp: now, StatusCode: 200}
		shortTracker.Record(ev)
		longTracker.Record(ev)
	}
	for i := 0; i < 2; i++ {
		ev := metrics.Event{Timestamp: now, StatusCode: 500}
		shortTracker.Record(ev)
		longTracker.Record(ev)
	}

	alerts := engine.Check(now)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert triggered, got %d", len(alerts))
	}
	if alerts[0].Severity != alerting.SeverityPage {
		t.Fatalf("expected Page alert severity, got %s", alerts[0].Severity)
	}

	// Negative Test: Transient spike in short window only, long window error rate is low.
	// Short window: 10% errors (100x burn rate).
	// Long window: 10,000 requests, 1 error = 0.01% error rate (0.1x burn rate < 14.4x threshold).
	shortOnlyTracker := metrics.NewWindowTracker(5*time.Minute, time.Second, isGood)
	longCleanTracker := metrics.NewWindowTracker(60*time.Minute, time.Second, isGood)
	engineTransient := alerting.NewAlertEngine(0.999, shortOnlyTracker, longCleanTracker, rules)

	for i := 0; i < 90; i++ {
		shortOnlyTracker.Record(metrics.Event{Timestamp: now, StatusCode: 200})
	}
	for i := 0; i < 10; i++ {
		shortOnlyTracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})
	}
	for i := 0; i < 9999; i++ {
		longCleanTracker.Record(metrics.Event{Timestamp: now, StatusCode: 200})
	}
	longCleanTracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})

	transientAlerts := engineTransient.Check(now)
	if len(transientAlerts) != 0 {
		t.Fatalf("expected 0 alerts for transient spike (long window below threshold), got %d", len(transientAlerts))
	}
}

func TestOutOfOrderTimestamps(t *testing.T) {
	tracker := metrics.NewWindowTracker(10*time.Second, 1*time.Second, func(e metrics.Event) bool {
		return e.StatusCode == 200
	})
	now := time.Now()

	// Record later event first
	tracker.Record(metrics.Event{Timestamp: now.Add(5 * time.Second), StatusCode: 200})
	// Record earlier event second
	tracker.Record(metrics.Event{Timestamp: now.Add(2 * time.Second), StatusCode: 200})
	// Record another event in the same earlier bucket
	tracker.Record(metrics.Event{Timestamp: now.Add(2 * time.Second + 100*time.Millisecond), StatusCode: 500})

	total, good, bad := tracker.Summary(now.Add(6 * time.Second))
	if total != 3 || good != 2 || bad != 1 {
		t.Fatalf("expected total 3, good 2, bad 1; got total=%d, good=%d, bad=%d", total, good, bad)
	}

	// Advance time past the earlier bucket to verify sorted slice eviction
	total, good, bad = tracker.Summary(now.Add(13 * time.Second)) // cutoff is now + 3s, earlier bucket at +2s is evicted, +5s bucket remains
	if total != 1 || good != 1 || bad != 0 {
		t.Fatalf("expected total 1, good 1, bad 0 after partial eviction; got total=%d, good=%d, bad=%d", total, good, bad)
	}
}

func TestEvaluatorZeroTraffic(t *testing.T) {
	tracker := metrics.NewWindowTracker(1*time.Minute, 1*time.Second, func(e metrics.Event) bool {
		return e.StatusCode == 200
	})
	evaluator := slo.NewEvaluator(slo.Config{
		Name:         "API SLO",
		TargetUptime: 0.99,
	}, tracker)

	status := evaluator.Evaluate(time.Now())
	if status.TotalEvents != 0 {
		t.Fatalf("expected 0 total events, got %d", status.TotalEvents)
	}
	if status.CurrentSLI != 1.0 {
		t.Fatalf("expected SLI 1.0 on zero traffic, got %f", status.CurrentSLI)
	}
	if !status.CanDeploy {
		t.Fatalf("expected CanDeploy=true on zero traffic, got %v", status.CanDeploy)
	}
}

func TestConcurrencyMetrics(t *testing.T) {
	tracker := metrics.NewWindowTracker(10*time.Second, 100*time.Millisecond, func(e metrics.Event) bool {
		return e.StatusCode == 200
	})

	now := time.Now()
	var wg sync.WaitGroup
	numGoroutines := 20
	requestsPerGoroutine := 100

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < requestsPerGoroutine; i++ {
				status := 200
				if i%10 == 0 {
					status = 500
				}
				tracker.Record(metrics.Event{
					Timestamp:  now.Add(time.Duration(i) * time.Millisecond),
					StatusCode: status,
				})
			}
		}(g)
	}
	wg.Wait()

	total, good, bad := tracker.Summary(now.Add(time.Second))
	expectedTotal := int64(numGoroutines * requestsPerGoroutine)
	if total != expectedTotal {
		t.Fatalf("expected total %d, got %d", expectedTotal, total)
	}
	if good+bad != total {
		t.Fatalf("expected good (%d) + bad (%d) == total (%d)", good, bad, total)
	}
}
