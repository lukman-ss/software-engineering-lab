package safe

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metrics holds Prometheus collectors for idempotency operations.
// Labels are deliberately coarse (operation, status) to avoid high-cardinality
// from business identifiers (user_id, order_id) which would blow up the TSDB.
var metrics = struct {
	requestsTotal  *prometheus.CounterVec
	replaysTotal   *prometheus.CounterVec
	conflictsTotal *prometheus.CounterVec
	processingTotal *prometheus.GaugeVec
	failuresTotal  *prometheus.CounterVec
}{
	requestsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "idempotency_requests_total",
			Help: "Total number of idempotency requests by operation and status.",
		},
		[]string{"operation", "status"}, // e.g. operation="payment", status="completed"
	),
	replaysTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "idempotency_replays_total",
			Help: "Total number of idempotent replays (request with existing key+hash).",
		},
		[]string{"operation"}, // e.g. operation="payment"
	),
	conflictsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "idempotency_conflicts_total",
			Help: "Total number of idempotency conflicts (key exists with different payload).",
		},
		[]string{"operation", "reason"}, // e.g. reason="payload_mismatch"
	),
	processingTotal: promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "idempotency_processing_total",
			Help: "Current number of in-flight processing idempotency records.",
		},
		[]string{"operation"},
	),
	failuresTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "idempotency_failures_total",
			Help: "Total number of idempotency failures by operation and failure type.",
		},
		[]string{"operation", "type"}, // e.g. type="gateway_error"
	),
}

// Labels used for Prometheus. Never include user_id, order_id, etc.
const (
	LabelOperationPayment = "payment"
	LabelOperationGeneric = "generic"

	LabelStatusCompleted = "completed"
	LabelStatusConflict = "conflict"
	LabelStatusFailed   = "failed"

	LabelReasonPayloadMismatch = "payload_mismatch"
	LabelReasonDuplicateOrder  = "duplicate_order"

	LabelTypeGatewayError = "gateway_error"
	LabelTypeStorageError = "storage_error"
)
