package monitor

import (
	"sync"
	"sync/atomic"
)

type SteadyStateMetrics struct {
	TotalRequests   uint64
	FailedRequests  uint64
	SuccessRequests uint64
}

type Monitor struct {
	mu            sync.RWMutex
	totalRequests uint64
	failedRequests uint64
	maxErrorRate  float64 // Threshold e.g. 0.20 (20%)
}

func NewMonitor(maxErrorRate float64) *Monitor {
	return &Monitor{
		maxErrorRate: maxErrorRate,
	}
}

func (m *Monitor) RecordSuccess() {
	atomic.AddUint64(&m.totalRequests, 1)
}

func (m *Monitor) RecordFailure() {
	atomic.AddUint64(&m.totalRequests, 1)
	atomic.AddUint64(&m.failedRequests, 1)
}

func (m *Monitor) Metrics() SteadyStateMetrics {
	total := atomic.LoadUint64(&m.totalRequests)
	failed := atomic.LoadUint64(&m.failedRequests)
	return SteadyStateMetrics{
		TotalRequests:   total,
		FailedRequests:  failed,
		SuccessRequests: total - failed,
	}
}

func (m *Monitor) ErrorRate() float64 {
	total := atomic.LoadUint64(&m.totalRequests)
	if total == 0 {
		return 0.0
	}
	failed := atomic.LoadUint64(&m.failedRequests)
	return float64(failed) / float64(total)
}

func (m *Monitor) IsHealthy() bool {
	total := atomic.LoadUint64(&m.totalRequests)
	if total < 5 { // Minimum sample before evaluating breach
		return true
	}
	return m.ErrorRate() <= m.maxErrorRate
}
