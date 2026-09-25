package loadtest

import (
	"sort"
	"time"
)

type Result struct {
	TotalRequests int
	SuccessCount  int
	ErrorCount    int
	Duration      time.Duration
	RPS           float64
	MinLatency    time.Duration
	MaxLatency    time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P90Latency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
}

func CalculateMetrics(latencies []time.Duration, errors int, totalDuration time.Duration) Result {
	total := len(latencies) + errors
	if total == 0 {
		return Result{}
	}

	res := Result{
		TotalRequests: total,
		SuccessCount:  len(latencies),
		ErrorCount:    errors,
		Duration:      totalDuration,
	}

	if totalDuration > 0 {
		res.RPS = float64(total) / totalDuration.Seconds()
	}

	if len(latencies) == 0 {
		return res
	}

	// ponytail: slice sorting fine for <1M samples; upgrade to streaming histogram if memory constrained
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	var sum time.Duration
	for _, lat := range sorted {
		sum += lat
	}

	res.MinLatency = sorted[0]
	res.MaxLatency = sorted[len(sorted)-1]
	res.AvgLatency = sum / time.Duration(len(sorted))
	res.P50Latency = percentile(sorted, 50)
	res.P90Latency = percentile(sorted, 90)
	res.P95Latency = percentile(sorted, 95)
	res.P99Latency = percentile(sorted, 99)

	return res
}

func percentile(sorted []time.Duration, pct float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * (pct / 100.0))
	return sorted[idx]
}
