package observability

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
)

// SafeInvoiceProcessor implements full observability:
// 1. Structured logging with request_id & trace_id
// 2. OpenTelemetry-compatible tracing with child spans
// 3. Prometheus metrics for every dependency (duration and error counts)
// 4. Context propagation across all steps
type SafeInvoiceProcessor struct {
	Repo         InvoiceRepository
	Inventory    InventoryService
	Commission   CommissionService
	PDF          PDFGenerator
	Notification NotificationService
	Logger       *slog.Logger
	Tracer       *Tracer
	Metrics      *Metrics
}

func NewSafeInvoiceProcessor(
	repo InvoiceRepository,
	inv InventoryService,
	comm CommissionService,
	pdf PDFGenerator,
	notif NotificationService,
	logger *slog.Logger,
	tracer *Tracer,
	metrics *Metrics,
) *SafeInvoiceProcessor {
	if tracer == nil {
		tracer = NewTracer()
	}
	if metrics == nil {
		metrics = NewMetrics()
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &SafeInvoiceProcessor{
		Repo:         repo,
		Inventory:    inv,
		Commission:   comm,
		PDF:          pdf,
		Notification: notif,
		Logger:       logger,
		Tracer:       tracer,
		Metrics:      metrics,
	}
}

func (p *SafeInvoiceProcessor) ProcessInvoice(ctx context.Context, invoiceID string) error {
	ctx, rootSpan := p.Tracer.Start(ctx, "invoice.process")
	defer p.Tracer.End(rootSpan)
	rootSpan.SetAttribute("invoice_id", invoiceID)

	requestID := middleware.GetRequestID(ctx)
	traceID := rootSpan.TraceID

	logger := p.Logger.With(
		"request_id", requestID,
		"trace_id", traceID,
		"invoice_id", invoiceID,
	)

	logger.Info("starting invoice processing")

	steps := []struct {
		name string
		op   func(ctx context.Context) error
	}{
		{
			name: "database.load_invoice",
			op: func(ctx context.Context) error {
				return p.Repo.Load(ctx, invoiceID)
			},
		},
		{
			name: "inventory.reserve",
			op: func(ctx context.Context) error {
				return p.Inventory.Reserve(ctx, invoiceID)
			},
		},
		{
			name: "commission.calculate",
			op: func(ctx context.Context) error {
				return p.Commission.Calculate(ctx, invoiceID)
			},
		},
		{
			name: "pdf.generate",
			op: func(ctx context.Context) error {
				return p.PDF.Generate(ctx, invoiceID)
			},
		},
		{
			name: "notification.send",
			op: func(ctx context.Context) error {
				return p.Notification.Send(ctx, invoiceID)
			},
		},
	}

	for _, step := range steps {
		stepCtx, span := p.Tracer.Start(ctx, step.name)
		stepStart := time.Now()

		err := step.op(stepCtx)
		stepDuration := time.Since(stepStart)

		if err != nil {
			span.RecordError(err)
			p.Tracer.End(span)
			rootSpan.RecordError(err)

			p.Metrics.Inc("dependency_errors_total", map[string]string{
				"dependency": step.name,
				"error":      err.Error(),
			})
			p.Metrics.Observe("dependency_duration_seconds", stepDuration, map[string]string{
				"dependency": step.name,
				"status":     "error",
			})

			logger.Error("dependency step failed",
				"step", step.name,
				"duration_ms", stepDuration.Milliseconds(),
				"span_id", span.SpanID,
				"error", err,
			)
			return fmt.Errorf("%s: %w", step.name, err)
		}

		p.Tracer.End(span)
		p.Metrics.Inc("dependency_requests_total", map[string]string{
			"dependency": step.name,
			"status":     "success",
		})
		p.Metrics.Observe("dependency_duration_seconds", stepDuration, map[string]string{
			"dependency": step.name,
			"status":     "success",
		})

		logger.Info("dependency step completed",
			"step", step.name,
			"duration_ms", stepDuration.Milliseconds(),
			"span_id", span.SpanID,
		)
	}

	logger.Info("invoice processing completed successfully",
		"total_duration_ms", time.Since(rootSpan.StartTime).Milliseconds(),
	)
	return nil
}

func (p *SafeInvoiceProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.URL.Query().Get("id")
	if invoiceID == "" {
		http.Error(w, "missing invoice id", http.StatusBadRequest)
		return
	}

	err := p.ProcessInvoice(r.Context(), invoiceID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%s","request_id":"%s"}`, err.Error(), middleware.GetRequestID(r.Context()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
