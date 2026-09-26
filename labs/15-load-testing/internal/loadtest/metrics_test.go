package loadtest

import (
	"testing"
	"time"
)

func TestCalculateMetrics(t *testing.T) {
	latencies := make([]time.Duration, 100)
	for i := 0; i < 100; i++ {
		latencies[i] = time.Duration(i+1) * time.Millisecond
	}

	res := CalculateMetrics(latencies, 0, 1*time.Second)

	if res.TotalRequests != 100 {
		t.Fatalf("expected 100 requests, got %d", res.TotalRequests)
	}
	if res.MinLatency != 1*time.Millisecond {
		t.Fatalf("expected min 1ms, got %s", res.MinLatency)
	}
	if res.MaxLatency != 100*time.Millisecond {
		t.Fatalf("expected max 100ms, got %s", res.MaxLatency)
	}
	// Average: (1+100)*100/2 / 100 = 50.5ms = 50ms (integer division)
	if res.AvgLatency != 50*time.Millisecond && res.AvgLatency != 50500*time.Microsecond {
		t.Fatalf("unexpected avg %s", res.AvgLatency)
	}
	// P50: index 99 * 0.5 = 49 -> 50ms
	if res.P50Latency != 50*time.Millisecond {
		t.Fatalf("expected P50 50ms, got %s", res.P50Latency)
	}
	// P95: index 99 * 0.95 = 94 -> 95ms
	if res.P95Latency != 95*time.Millisecond {
		t.Fatalf("expected P95 95ms, got %s", res.P95Latency)
	}
	// P99: index 99 * 0.99 = 98 -> 99ms
	if res.P99Latency != 99*time.Millisecond {
		t.Fatalf("expected P99 99ms, got %s", res.P99Latency)
	}
}

func TestCalculateMetrics_Empty(t *testing.T) {
	res := CalculateMetrics(nil, 0, time.Second)
	if res.TotalRequests != 0 {
		t.Fatalf("expected 0 requests, got %d", res.TotalRequests)
	}
	if res.P95Latency != 0 {
		t.Fatalf("expected 0 P95, got %s", res.P95Latency)
	}
}

func TestCalculateMetrics_Invariants(t *testing.T) {
	// Unordered list of latencies with long-tail spikes
	latencies := []time.Duration{
		45 * time.Millisecond, 12 * time.Millisecond, 100 * time.Millisecond,
		5 * time.Millisecond, 20 * time.Millisecond, 500 * time.Millisecond,
		15 * time.Millisecond, 22 * time.Millisecond, 18 * time.Millisecond,
		80 * time.Millisecond, 250 * time.Millisecond, 30 * time.Millisecond,
	}

	res := CalculateMetrics(latencies, 0, time.Second)

	if res.MinLatency > res.P50Latency {
		t.Errorf("invariant violated: Min (%s) > P50 (%s)", res.MinLatency, res.P50Latency)
	}
	if res.P50Latency > res.P90Latency {
		t.Errorf("invariant violated: P50 (%s) > P90 (%s)", res.P50Latency, res.P90Latency)
	}
	if res.P90Latency > res.P95Latency {
		t.Errorf("invariant violated: P90 (%s) > P95 (%s)", res.P90Latency, res.P95Latency)
	}
	if res.P95Latency > res.P99Latency {
		t.Errorf("invariant violated: P95 (%s) > P99 (%s)", res.P95Latency, res.P99Latency)
	}
	if res.P99Latency > res.MaxLatency {
		t.Errorf("invariant violated: P99 (%s) > Max (%s)", res.P99Latency, res.MaxLatency)
	}
}
