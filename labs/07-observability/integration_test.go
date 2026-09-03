package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestIntegration_RequestID_Propagation_And_Preservation(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := NewPrometheusCollector(reg)

	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test-tracer")

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{})
	service := NewSafeInvoiceService(repo, inv, comm, pdf, notif, logger, tracer, collector)

	// Test with explicit custom valid X-Request-ID
	customReqID := "demo-slow-pdf-001"
	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-999/process", nil)
	req.SetPathValue("id", "INV-999")
	req.Header.Set(middleware.HeaderRequestID, customReqID)

	rec := httptest.NewRecorder()
	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	respHeaderID := rec.Header().Get(middleware.HeaderRequestID)
	if respHeaderID != customReqID {
		t.Fatalf("expected X-Request-ID response header %s, got %s", customReqID, respHeaderID)
	}

	// Verify log contains request_id and trace_id
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, customReqID) {
		t.Errorf("log output missing request_id %s", customReqID)
	}
	if !strings.Contains(logOutput, "trace_id") {
		t.Errorf("log output missing trace_id")
	}

	// Test failed request retains request ID
	failRepo, failInv, failComm, failPDF, failNotif := NewFakeDependencies(ScenarioConfig{PDFErr: ErrPDFGeneration})
	failService := NewSafeInvoiceService(failRepo, failInv, failComm, failPDF, failNotif, logger, tracer, collector)

	failReq := httptest.NewRequest(http.MethodPost, "/invoices/INV-FAIL/process", nil)
	failReq.SetPathValue("id", "INV-FAIL")
	failReq.Header.Set(middleware.HeaderRequestID, "req-fail-123")
	failRec := httptest.NewRecorder()

	failHandler := middleware.RequestID(failService.HTTPHandler(nil))
	failHandler.ServeHTTP(failRec, failReq)

	if failRec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", failRec.Code)
	}
	if failRec.Header().Get(middleware.HeaderRequestID) != "req-fail-123" {
		t.Fatalf("expected preserved X-Request-ID in error response header")
	}
	var errBody map[string]string
	json.Unmarshal(failRec.Body.Bytes(), &errBody)
	if errBody["request_id"] != "req-fail-123" {
		t.Fatalf("expected request_id in error body, got %v", errBody)
	}
}

func TestIntegration_StructuredLogging_Schema(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")

	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{})
	service := NewSafeInvoiceService(repo, inv, comm, pdf, notif, logger, tracer, nil)

	ctx := context.Background()
	reqID := "demo-req-001"
	ctx = context.WithValue(ctx, middleware.HeaderRequestID, reqID)

	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-1/process", nil)
	req.SetPathValue("id", "INV-1")
	req.Header.Set(middleware.HeaderRequestID, reqID)
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	// Validate JSON lines
	lines := strings.Split(strings.TrimSpace(logBuf.String()), "\n")
	if len(lines) == 0 {
		t.Fatal("expected log lines")
	}

	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("failed to parse log JSON: %v, raw: %s", err, line)
		}

		if _, ok := entry["level"]; !ok {
			t.Errorf("missing level in log: %s", line)
		}
		if _, ok := entry["msg"]; !ok {
			t.Errorf("missing msg in log: %s", line)
		}
		if entry["request_id"] != reqID {
			t.Errorf("expected request_id=%s, got %v", reqID, entry["request_id"])
		}
		if entry["trace_id"] == "" || entry["trace_id"] == nil {
			t.Errorf("missing trace_id in log: %s", line)
		}
		if entry["span_id"] == "" || entry["span_id"] == nil {
			t.Errorf("missing span_id in log: %s", line)
		}

		// Security check: no forbidden fields
		for _, forbidden := range []string{"authorization", "password", "token", "secret", "phone"} {
			if _, ok := entry[forbidden]; ok {
				t.Errorf("forbidden sensitive field '%s' logged", forbidden)
			}
		}
	}
}

func TestIntegration_Prometheus_GoldenSignals(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := NewPrometheusCollector(reg)

	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")

	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		RepoDelay: 2 * time.Millisecond,
		PDFDelay:  5 * time.Millisecond,
	})
	service := NewSafeInvoiceService(repo, inv, comm, pdf, notif, nil, tracer, collector)

	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-1/process", nil)
	req.SetPathValue("id", "INV-1")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// 1. lab07_http_requests_total
	count := testutil.ToFloat64(collector.HTTPRequestsTotal.WithLabelValues("POST", "/invoices/{id}/process", "2xx"))
	if count != 1 {
		t.Fatalf("expected http_requests_total=1, got %f", count)
	}

	// 2. lab07_http_in_flight_requests returns to 0
	inFlight := testutil.ToFloat64(collector.HTTPInFlightRequests.WithLabelValues("/invoices/{id}/process"))
	if inFlight != 0 {
		t.Fatalf("expected in_flight_requests=0 after completion, got %f", inFlight)
	}

	// 3. lab07_dependency_duration_seconds has counts for pdf generate
	pdfLatencyCount := testutil.CollectAndCount(collector.DependencyDurationSeconds)
	if pdfLatencyCount == 0 {
		t.Fatalf("expected dependency duration metrics to be collected")
	}

	// Error scenario metrics check
	failRepo, failInv, failComm, failPDF, failNotif := NewFakeDependencies(ScenarioConfig{
		PDFErr: ErrPDFGeneration,
	})
	failService := NewSafeInvoiceService(failRepo, failInv, failComm, failPDF, failNotif, nil, tracer, collector)

	failReq := httptest.NewRequest(http.MethodPost, "/invoices/INV-FAIL/process", nil)
	failReq.SetPathValue("id", "INV-FAIL")
	failRec := httptest.NewRecorder()

	failHandler := middleware.RequestID(failService.HTTPHandler(nil))
	failHandler.ServeHTTP(failRec, failReq)

	errCount := testutil.ToFloat64(collector.DependencyErrorsTotal.WithLabelValues("pdf", "generate", "error"))
	if errCount != 1 {
		t.Fatalf("expected dependency_errors_total=1, got %f", errCount)
	}

	httpErrCount := testutil.ToFloat64(collector.HTTPRequestErrorsTotal.WithLabelValues("POST", "/invoices/{id}/process", "5xx"))
	if httpErrCount != 1 {
		t.Fatalf("expected http_request_errors_total=1, got %f", httpErrCount)
	}

	// Ensure in-flight gauge is 0 even after error
	inFlightAfterError := testutil.ToFloat64(collector.HTTPInFlightRequests.WithLabelValues("/invoices/{id}/process"))
	if inFlightAfterError != 0 {
		t.Fatalf("expected in_flight_requests=0 after failure, got %f", inFlightAfterError)
	}
}

func TestIntegration_Trace_Diagnosis_SlowPDF(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test")

	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		PDFDelay: 10 * time.Millisecond,
	})
	service := NewSafeInvoiceService(repo, inv, comm, pdf, notif, nil, tracer, nil)

	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-SLOW/process", nil)
	req.SetPathValue("id", "INV-SLOW")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	spans := sr.Ended()
	// Should have: database.load_invoice, inventory.reserve, commission.calculate, pdf.generate, notification.send, invoice.process, http.request
	if len(spans) < 6 {
		t.Fatalf("expected at least 6 spans, got %d", len(spans))
	}

	var pdfSpan sdktrace.ReadOnlySpan
	for _, s := range spans {
		if s.Name() == "pdf.generate" {
			pdfSpan = s
			break
		}
	}

	if pdfSpan == nil {
		t.Fatal("pdf.generate span was not created")
	}

	pdfDuration := pdfSpan.EndTime().Sub(pdfSpan.StartTime())
	if pdfDuration < 10*time.Millisecond {
		t.Fatalf("expected pdf duration >= 10ms, got %v", pdfDuration)
	}
}

func TestUnsafeInvoiceService_LacksGranularity(t *testing.T) {
	buf := &bytes.Buffer{}
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		PDFDelay: 10 * time.Millisecond,
	})
	unsafeService := &UnsafeInvoiceService{
		Repo:         repo,
		Inventory:    inv,
		Commission:   comm,
		PDF:          pdf,
		Notification: notif,
		LogWriter:    buf,
	}

	req := httptest.NewRequest(http.MethodPost, "/unsafe/invoices/INV-UNSAFE/process", nil)
	req.SetPathValue("id", "INV-UNSAFE")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(unsafeService.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	logs := buf.String()
	if !strings.Contains(logs, "INFO request completed duration_ms=") {
		t.Errorf("expected coarse duration log, got: %s", logs)
	}
	if strings.Contains(logs, "pdf.generate") {
		t.Errorf("unsafe logs should not contain child breakdown")
	}
}
