package observability

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
)

func TestUnsafeImplementation_LacksObservability(t *testing.T) {
	buf := &bytes.Buffer{}
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		PDFDelay: 5 * time.Millisecond,
	})

	processor := &UnsafeInvoiceProcessor{
		Repo:         repo,
		Inventory:    inv,
		Commission:   comm,
		PDF:          pdf,
		Notification: notif,
		LogWriter:    buf,
	}

	err := processor.ProcessInvoice(context.Background(), "INV-1001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logOutput := buf.String()
	// Unsafe log only gives request completed duration_ms without child breakdown
	if !strings.Contains(logOutput, "INFO request completed duration_ms=") {
		t.Errorf("expected summary log, got: %s", logOutput)
	}
	if strings.Contains(logOutput, "pdf.generate") {
		t.Errorf("unsafe logs should not contain dependency breakdown")
	}
}

func TestSafeImplementation_NormalScenario(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{})

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	tracer := NewTracer()
	metrics := NewMetrics()

	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, logger, tracer, metrics)

	req := httptest.NewRequest(http.MethodPost, "/invoices?id=INV-NORMAL", nil)
	req.Header.Set(middleware.HeaderRequestID, "req-test-123")
	rec := httptest.NewRecorder()

	handler := middleware.RequestID(http.HandlerFunc(processor.ServeHTTP))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// 1. Trace Assertions
	spans := tracer.Spans()
	// 1 root span + 5 child spans = 6 spans
	if len(spans) != 6 {
		t.Fatalf("expected 6 spans, got %d", len(spans))
	}

	spanNames := make(map[string]bool)
	for _, s := range spans {
		spanNames[s.Name] = true
		if s.TraceID == "" {
			t.Errorf("span %s missing TraceID", s.Name)
		}
	}

	expectedSpans := []string{
		"invoice.process",
		"database.load_invoice",
		"inventory.reserve",
		"commission.calculate",
		"pdf.generate",
		"notification.send",
	}
	for _, exp := range expectedSpans {
		if !spanNames[exp] {
			t.Errorf("missing expected span: %s", exp)
		}
	}

	// 2. Metrics Assertions
	pdfCount := metrics.GetCount("dependency_requests_total", map[string]string{
		"dependency": "pdf.generate",
		"status":     "success",
	})
	if pdfCount != 1 {
		t.Errorf("expected 1 pdf.generate success count, got %d", pdfCount)
	}

	// 3. Log Assertions
	logs := logBuf.String()
	if !strings.Contains(logs, "req-test-123") {
		t.Errorf("logs should include request_id")
	}
	if !strings.Contains(logs, "trace_id") {
		t.Errorf("logs should include trace_id")
	}
}

func TestSafeImplementation_SlowPDFScenario(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		PDFDelay: 10 * time.Millisecond,
	})

	tracer := NewTracer()
	metrics := NewMetrics()
	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, nil, tracer, metrics)

	err := processor.ProcessInvoice(context.Background(), "INV-SLOW-PDF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that the tracer accurately isolates pdf.generate as the bottleneck
	var pdfSpan *Span
	var otherMaxDuration time.Duration

	for _, s := range tracer.Spans() {
		if s.Name == "pdf.generate" {
			pdfSpan = s
		} else if s.Name != "invoice.process" {
			if s.Duration > otherMaxDuration {
				otherMaxDuration = s.Duration
			}
		}
	}

	if pdfSpan == nil {
		t.Fatal("pdf.generate span not found")
	}

	if pdfSpan.Duration < 10*time.Millisecond {
		t.Errorf("pdf duration should be at least 10ms, got %v", pdfSpan.Duration)
	}
	if pdfSpan.Duration <= otherMaxDuration {
		t.Errorf("pdf duration (%v) should exceed other dependencies (%v)", pdfSpan.Duration, otherMaxDuration)
	}
}

func TestSafeImplementation_SlowDatabaseScenario(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		RepoDelay: 10 * time.Millisecond,
	})

	tracer := NewTracer()
	metrics := NewMetrics()
	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, nil, tracer, metrics)

	err := processor.ProcessInvoice(context.Background(), "INV-SLOW-DB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var dbSpan *Span
	for _, s := range tracer.Spans() {
		if s.Name == "database.load_invoice" {
			dbSpan = s
			break
		}
	}

	if dbSpan == nil {
		t.Fatal("database span not found")
	}
	if dbSpan.Duration < 10*time.Millisecond {
		t.Errorf("database span should be at least 10ms, got %v", dbSpan.Duration)
	}
}

func TestSafeImplementation_SlowExternalScenario(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		NotificationDelay: 10 * time.Millisecond,
	})

	tracer := NewTracer()
	metrics := NewMetrics()
	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, nil, tracer, metrics)

	err := processor.ProcessInvoice(context.Background(), "INV-SLOW-EXT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var notifSpan *Span
	for _, s := range tracer.Spans() {
		if s.Name == "notification.send" {
			notifSpan = s
			break
		}
	}

	if notifSpan == nil {
		t.Fatal("notification span not found")
	}
	if notifSpan.Duration < 10*time.Millisecond {
		t.Errorf("notification span should be at least 10ms, got %v", notifSpan.Duration)
	}
}

func TestSafeImplementation_PDFErrorScenario(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		PDFErr: ErrPDFGeneration,
	})

	tracer := NewTracer()
	metrics := NewMetrics()
	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, nil, tracer, metrics)

	err := processor.ProcessInvoice(context.Background(), "INV-PDF-ERR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrPDFGeneration) {
		t.Errorf("expected ErrPDFGeneration, got %v", err)
	}

	// Verify error metric recorded
	errCount := metrics.GetCount("dependency_errors_total", map[string]string{
		"dependency": "pdf.generate",
		"error":      ErrPDFGeneration.Error(),
	})
	if errCount != 1 {
		t.Errorf("expected 1 recorded error for pdf.generate, got %d", errCount)
	}

	// Verify trace span recorded error
	var pdfSpan *Span
	for _, s := range tracer.Spans() {
		if s.Name == "pdf.generate" {
			pdfSpan = s
			break
		}
	}
	if pdfSpan == nil || pdfSpan.Err == nil {
		t.Errorf("expected pdf.generate span with error attribute recorded")
	}
}

func TestSafeImplementation_ContextCancellation(t *testing.T) {
	repo, inv, comm, pdf, notif := NewFakeDependencies(ScenarioConfig{
		RepoDelay: 50 * time.Millisecond,
	})

	processor := NewSafeInvoiceProcessor(repo, inv, comm, pdf, notif, nil, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	err := processor.ProcessInvoice(ctx, "INV-TIMEOUT")
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
