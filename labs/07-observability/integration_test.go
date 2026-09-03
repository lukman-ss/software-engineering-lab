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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestIntegration_SpanHierarchy_And_TraceContext(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test")

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), nil, tracer, nil)

	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-HIERARCHY/process", nil)
	req.SetPathValue("id", "INV-HIERARCHY")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	spans := sr.Ended()
	if len(spans) != 7 {
		t.Fatalf("expected 7 spans (1 HTTP root + 1 process + 5 dependencies), got %d", len(spans))
	}

	var httpRootSpan sdktrace.ReadOnlySpan
	var invoiceProcessSpan sdktrace.ReadOnlySpan
	depSpans := make(map[string]sdktrace.ReadOnlySpan)

	for _, s := range spans {
		switch s.Name() {
		case "POST /invoices/{id}/process":
			httpRootSpan = s
		case "invoice.process":
			invoiceProcessSpan = s
		default:
			depSpans[s.Name()] = s
		}
	}

	if httpRootSpan == nil {
		t.Fatal("missing root span 'POST /invoices/{id}/process'")
	}
	if invoiceProcessSpan == nil {
		t.Fatal("missing child span 'invoice.process'")
	}

	traceID := httpRootSpan.SpanContext().TraceID()

	// All spans must share the same trace ID
	if invoiceProcessSpan.SpanContext().TraceID() != traceID {
		t.Errorf("invoice.process has different trace ID: %s vs %s", invoiceProcessSpan.SpanContext().TraceID(), traceID)
	}

	// invoice.process must have httpRootSpan as parent
	if invoiceProcessSpan.Parent().SpanID() != httpRootSpan.SpanContext().SpanID() {
		t.Errorf("invoice.process parent ID %s != root span ID %s", invoiceProcessSpan.Parent().SpanID(), httpRootSpan.SpanContext().SpanID())
	}

	expectedDeps := []string{
		"database.load_invoice",
		"inventory.reserve",
		"commission.calculate",
		"pdf.generate",
		"notification.send",
	}

	for _, depName := range expectedDeps {
		s, ok := depSpans[depName]
		if !ok {
			t.Errorf("missing expected dependency span: %s", depName)
			continue
		}
		if s.SpanContext().TraceID() != traceID {
			t.Errorf("dep span %s has mismatched trace ID: %s vs %s", depName, s.SpanContext().TraceID(), traceID)
		}
		if s.Parent().SpanID() != invoiceProcessSpan.SpanContext().SpanID() {
			t.Errorf("dep span %s parent ID %s != invoice.process span ID %s", depName, s.Parent().SpanID(), invoiceProcessSpan.SpanContext().SpanID())
		}
	}
}

func TestIntegration_RequestID_Validation_And_Generation(t *testing.T) {
	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), nil, nil, nil)
	handler := middleware.RequestID(service.HTTPHandler(nil))

	// 1. Missing request ID -> newly generated
	req1 := httptest.NewRequest(http.MethodPost, "/invoices/INV-1/process", nil)
	req1.SetPathValue("id", "INV-1")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	id1 := rec1.Header().Get(middleware.HeaderRequestID)
	if id1 == "" {
		t.Fatal("expected generated request ID when missing")
	}

	// 2. Invalid request ID (contains invalid chars) -> replaced with new valid ID
	req2 := httptest.NewRequest(http.MethodPost, "/invoices/INV-1/process", nil)
	req2.SetPathValue("id", "INV-1")
	req2.Header.Set(middleware.HeaderRequestID, "invalid id with spaces &!@#")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	id2 := rec2.Header().Get(middleware.HeaderRequestID)
	if id2 == "" || id2 == "invalid id with spaces &!@#" {
		t.Fatalf("expected invalid request ID to be replaced, got %s", id2)
	}

	// 3. Error response retains request ID in body and headers
	failService := NewSafeInvoiceService(NewDependencies(ScenarioConfig{PDFErr: ErrPDFGeneration}), nil, nil, nil)
	failHandler := middleware.RequestID(failService.HTTPHandler(nil))

	req3 := httptest.NewRequest(http.MethodPost, "/invoices/INV-FAIL/process", nil)
	req3.SetPathValue("id", "INV-FAIL")
	req3.Header.Set(middleware.HeaderRequestID, "req-valid-12345")
	rec3 := httptest.NewRecorder()
	failHandler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusBadGateway {
		t.Fatalf("expected 502 Bad Gateway, got %d", rec3.Code)
	}
	if rec3.Header().Get(middleware.HeaderRequestID) != "req-valid-12345" {
		t.Errorf("expected header request ID req-valid-12345, got %s", rec3.Header().Get(middleware.HeaderRequestID))
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(rec3.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to decode error json: %v", err)
	}
	if errResp.RequestID != "req-valid-12345" {
		t.Errorf("expected json error request_id=req-valid-12345, got %s", errResp.RequestID)
	}
	if errResp.Error != "invoice processing failed" {
		t.Errorf("expected generic error message, got %s", errResp.Error)
	}
}

func TestIntegration_NotificationClient_Cancellation(t *testing.T) {
	reqReceived := make(chan bool, 1)
	downstreamSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqReceived <- true
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer downstreamSrv.Close()

	client := NewHTTPNotificationClient(downstreamSrv.URL, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := client.Send(ctx, "INV-TIMEOUT")
	if err == nil {
		t.Fatal("expected cancellation / timeout error, got nil")
	}
}

func TestIntegration_RequestID_Propagation_And_Preservation(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := NewPrometheusCollector(reg)

	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test-tracer")

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), logger, tracer, collector)

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
	failDeps := NewDependencies(ScenarioConfig{PDFErr: ErrPDFGeneration})
	failService := NewSafeInvoiceService(failDeps, logger, tracer, collector)

	failReq := httptest.NewRequest(http.MethodPost, "/invoices/INV-FAIL/process", nil)
	failReq.SetPathValue("id", "INV-FAIL")
	failReq.Header.Set(middleware.HeaderRequestID, "req-fail-123")
	failRec := httptest.NewRecorder()

	failHandler := middleware.RequestID(failService.HTTPHandler(nil))
	failHandler.ServeHTTP(failRec, failReq)

	if failRec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502 Bad Gateway, got %d", failRec.Code)
	}
	if failRec.Header().Get(middleware.HeaderRequestID) != "req-fail-123" {
		t.Fatalf("expected preserved X-Request-ID in error response header")
	}
	var errBody map[string]string
	json.Unmarshal(failRec.Body.Bytes(), &errBody)
	if errBody["request_id"] != "req-fail-123" {
		t.Fatalf("expected request_id in error body, got %v", errBody)
	}
	if errBody["error"] != "invoice processing failed" {
		t.Fatalf("expected sanitized error message in body, got %s", errBody["error"])
	}
}

func TestIntegration_StructuredLogging_Schema(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), logger, tracer, nil)

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
		msg, ok := entry["msg"].(string)
		if !ok || msg == "" {
			t.Errorf("missing or empty msg in log: %s", line)
		}
		if entry["request_id"] != reqID {
			t.Errorf("expected request_id=%s, got %v", reqID, entry["request_id"])
		}
		if entry["trace_id"] == "" || entry["trace_id"] == nil {
			t.Errorf("missing trace_id in log: %s", line)
		}

		if msg == "dependency completed" {
			if _, ok := entry["component"]; !ok {
				t.Errorf("missing component in dependency log: %s", line)
			}
			if _, ok := entry["operation"]; !ok {
				t.Errorf("missing operation in dependency log: %s", line)
			}
			if _, ok := entry["duration_ms"]; !ok {
				t.Errorf("missing duration_ms in dependency log: %s", line)
			}
			if _, ok := entry["outcome"]; !ok {
				t.Errorf("missing outcome in dependency log: %s", line)
			}
		}

		// Structured logging with request_id, trace_id, component, operation, duration_ms, outcome
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

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{
		RepoDelay: 2 * time.Millisecond,
		PDFDelay:  5 * time.Millisecond,
	}), nil, tracer, collector)

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
	failService := NewSafeInvoiceService(NewDependencies(ScenarioConfig{
		PDFErr: ErrPDFGeneration,
	}), nil, tracer, collector)

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

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{
		PDFDelay: 10 * time.Millisecond,
	}), nil, tracer, nil)

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

func TestIntegration_TracePropagation(t *testing.T) {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	var downstreamTraceparent string
	downstreamSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamTraceparent = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer downstreamSrv.Close()

	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test")

	notifClient := NewHTTPNotificationClient(downstreamSrv.URL, nil)

	deps := NewDependencies(ScenarioConfig{})
	deps.Notification = notifClient

	service := NewSafeInvoiceService(deps, nil, tracer, nil)

	req := httptest.NewRequest(http.MethodPost, "/invoices/INV-PROP/process", nil)
	req.SetPathValue("id", "INV-PROP")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if downstreamTraceparent == "" {
		t.Fatal("expected downstream service to receive traceparent header")
	}

	// Validate trace context
	parts := strings.Split(downstreamTraceparent, "-")
	if len(parts) != 4 || parts[0] != "00" {
		t.Fatalf("invalid traceparent format: %s", downstreamTraceparent)
	}
	traceIDStr := parts[1]

	// Ensure trace ID matches the root span
	var rootSpan sdktrace.ReadOnlySpan
	for _, s := range sr.Ended() {
		if s.Name() == "POST /invoices/{id}/process" {
			rootSpan = s
			break
		}
	}
	if rootSpan == nil {
		t.Fatal("root span not found")
	}

	if rootSpan.SpanContext().TraceID().String() != traceIDStr {
		t.Fatalf("trace ID mismatch: root=%s, downstream=%s", rootSpan.SpanContext().TraceID().String(), traceIDStr)
	}
}

func TestNoHighCardinalityLabelsInPrometheusRegistry(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := NewPrometheusCollector(reg)

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), nil, nil, collector)

	reqID := "req-high-cardinality-uuid-12345"
	invoiceID := "INV-SENSITIVE-9999"

	req := httptest.NewRequest(http.MethodPost, "/invoices/"+invoiceID+"/process", nil)
	req.SetPathValue("id", invoiceID)
	req.Header.Set(middleware.HeaderRequestID, reqID)
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(service.HTTPHandler(nil))
	handler.ServeHTTP(rec, req)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	forbiddenValues := []string{
		reqID,
		invoiceID,
	}

	for _, mf := range mfs {
		for _, m := range mf.Metric {
			for _, pair := range m.Label {
				labelName := pair.GetName()
				labelVal := pair.GetValue()

				for _, forbiddenKey := range []string{"request_id", "trace_id", "span_id", "invoice_id", "customer_id", "user_id", "error", "timestamp", "url"} {
					if strings.EqualFold(labelName, forbiddenKey) {
						t.Errorf("found forbidden metric label name '%s' in metric '%s'", labelName, mf.GetName())
					}
				}

				for _, forbiddenVal := range forbiddenValues {
					if strings.Contains(labelVal, forbiddenVal) {
						t.Errorf("found forbidden high-cardinality value '%s' inside label '%s' for metric '%s'", forbiddenVal, labelName, mf.GetName())
					}
				}
			}
		}
	}
}

func TestConcurrentRequests_NoDataRace(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := NewPrometheusCollector(reg)
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("test")

	service := NewSafeInvoiceService(NewDependencies(ScenarioConfig{}), nil, tracer, collector)

	resolveScenario := func(r *http.Request) ScenarioConfig {
		switch r.URL.Query().Get("scenario") {
		case "error":
			return ScenarioConfig{PDFErr: ErrPDFGeneration}
		case "slow":
			return ScenarioConfig{PDFDelay: 5 * time.Millisecond}
		default:
			return ScenarioConfig{RepoDelay: 2 * time.Millisecond}
		}
	}

	handler := middleware.RequestID(service.HTTPHandler(resolveScenario))

	done := make(chan bool)
	scenarios := []string{"normal", "error", "slow", "normal", "error"}

	for i, sc := range scenarios {
		go func(idx int, scenario string) {
			url := "/invoices/INV-CONC-" + string(rune('A'+idx)) + "/process?scenario=" + scenario
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req.SetPathValue("id", "INV-CONC-"+string(rune('A'+idx)))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			done <- true
		}(i, sc)
	}

	for range scenarios {
		<-done
	}
}

func TestUnsafeInvoiceService_LacksGranularity(t *testing.T) {
	buf := &bytes.Buffer{}
	unsafeService := &UnsafeInvoiceService{
		Deps: NewDependencies(ScenarioConfig{
			PDFDelay: 10 * time.Millisecond,
		}),
		LogWriter: buf,
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

func TestUnsafeConcurrentRequests_NoDataRace(t *testing.T) {
	unsafeService := &UnsafeInvoiceService{
		Deps: NewDependencies(ScenarioConfig{}),
	}

	resolveScenario := func(r *http.Request) ScenarioConfig {
		switch r.URL.Query().Get("scenario") {
		case "error":
			return ScenarioConfig{PDFErr: ErrPDFGeneration}
		case "slow":
			return ScenarioConfig{PDFDelay: 5 * time.Millisecond}
		default:
			return ScenarioConfig{RepoDelay: 2 * time.Millisecond}
		}
	}

	handler := middleware.RequestID(unsafeService.HTTPHandler(resolveScenario))

	done := make(chan bool)
	scenarios := []string{"normal", "error", "slow", "normal", "error"}

	for i, sc := range scenarios {
		go func(idx int, scenario string) {
			url := "/unsafe/invoices/INV-CONC-" + string(rune('A'+idx)) + "/process?scenario=" + scenario
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req.SetPathValue("id", "INV-CONC-"+string(rune('A'+idx)))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			done <- true
		}(i, sc)
	}

	for range scenarios {
		<-done
	}
}
