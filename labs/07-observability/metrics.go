package observability

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// Metrics records Prometheus-style dimensional counters and histograms.
type Metrics struct {
	mu           sync.RWMutex
	requestCount map[string]int64
	durations    map[string][]time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{
		requestCount: make(map[string]int64),
		durations:    make(map[string][]time.Duration),
	}
}

func (m *Metrics) Inc(metricName string, labels map[string]string) {
	key := formatMetricKey(metricName, labels)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount[key]++
}

func (m *Metrics) Observe(metricName string, duration time.Duration, labels map[string]string) {
	key := formatMetricKey(metricName, labels)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.durations[key] = append(m.durations[key], duration)
}

func (m *Metrics) GetCount(metricName string, labels map[string]string) int64 {
	key := formatMetricKey(metricName, labels)
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.requestCount[key]
}

func (m *Metrics) GetDurations(metricName string, labels map[string]string) []time.Duration {
	key := formatMetricKey(metricName, labels)
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]time.Duration, len(m.durations[key]))
	copy(res, m.durations[key])
	return res
}

func (m *Metrics) WritePrometheus(w io.Writer) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Write counters
	for key, count := range m.requestCount {
		fmt.Fprintf(w, "%s %d\n", key, count)
	}

	// Write summary / counts for durations
	for key, durations := range m.durations {
		var sum time.Duration
		for _, d := range durations {
			sum += d
		}
		fmt.Fprintf(w, "%s_count %d\n", key, len(durations))
		fmt.Fprintf(w, "%s_sum_ms %d\n", key, sum.Milliseconds())
	}
}

func formatMetricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	res := name + "{"
	for i, k := range keys {
		if i > 0 {
			res += ","
		}
		res += fmt.Sprintf("%s=\"%s\"", k, labels[k])
	}
	res += "}"
	return res
}
