package featureflags

import (
	"fmt"
	"sync"
	"time"
)

// Metrics tracks feature flag usage and outcomes
type Metrics struct {
	mu                   sync.Mutex
	TotalRequests        int64
	NewFlowRequests      int64
	LegacyFlowRequests   int64
	BookingSuccess       int64
	BookingFailed        int64
	ResponseTime         []time.Duration
	ResponseTimeDuration time.Duration
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		ResponseTime: make([]time.Duration, 0),
	}
}

// IncrementTotal increments total request counter
func (m *Metrics) IncrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
}

// RecordFlow records which flow was used
func (m *Metrics) RecordFlow(flow string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if flow == "online_booking" {
		m.NewFlowRequests++
	} else {
		m.LegacyFlowRequests++
	}
}

// RecordBooking records booking outcome
func (m *Metrics) RecordBooking(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if success {
		m.BookingSuccess++
	} else {
		m.BookingFailed++
	}
}

// RecordResponseTime adds response time measurement
func (m *Metrics) RecordResponseTime(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseTime = append(m.ResponseTime, d)
	m.ResponseTimeDuration += d
}

// Reset resets metrics counters for next simulation phase
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests = 0
	m.NewFlowRequests = 0
	m.LegacyFlowRequests = 0
	m.BookingSuccess = 0
	m.BookingFailed = 0
	m.ResponseTime = make([]time.Duration, 0)
	m.ResponseTimeDuration = 0
}
// Snapshot returns current metrics as a FeaturesSnapshot
func (m *Metrics) Snapshot() FeaturesSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	var avgResponseTime time.Duration
	if len(m.ResponseTime) > 0 {
		avgResponseTime = m.ResponseTimeDuration / time.Duration(len(m.ResponseTime))
	}

	return FeaturesSnapshot{
		TotalRequests:      m.TotalRequests,
		NewFlowRequests:    m.NewFlowRequests,
		LegacyFlowRequests: m.LegacyFlowRequests,
		BookingSuccess:     m.BookingSuccess,
		BookingFailed:      m.BookingFailed,
		ErrorRate:            float64(m.BookingFailed) / float64(max(1, m.TotalRequests)) * 100,
		AvgResponseTime:    avgResponseTime,
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// FeaturesSnapshot is a point-in-time view of metrics
type FeaturesSnapshot struct {
	TotalRequests      int64
	NewFlowRequests    int64
	LegacyFlowRequests int64
	BookingSuccess     int64
	BookingFailed      int64
	ErrorRate          float64
	AvgResponseTime    time.Duration
}

func (s FeaturesSnapshot) String() string {
	return fmt.Sprintf(`Feature Flag Metrics:
  Total Requests: %d
  New Flow: %d (%.1f%%)
  Legacy Flow: %d (%.1f%%)
  Booking Success: %d
  Booking Failed: %d (%.2f%% error rate)
  Avg Response Time: %v`,
		s.TotalRequests,
		s.NewFlowRequests, float64(s.NewFlowRequests)/float64(max(1, s.TotalRequests))*100,
		s.LegacyFlowRequests, float64(s.LegacyFlowRequests)/float64(max(1, s.TotalRequests))*100,
		s.BookingSuccess,
		s.BookingFailed, s.ErrorRate,
		s.AvgResponseTime)
}