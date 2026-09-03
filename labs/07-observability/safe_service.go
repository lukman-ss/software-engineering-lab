package observability

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SafeInvoiceService struct {
	Deps         Dependencies
	Logger       *slog.Logger
	Tracer       trace.Tracer
	Collector    *PrometheusCollector
}

func NewSafeInvoiceService(
	deps Dependencies,
	logger *slog.Logger,
	tracer trace.Tracer,
	collector *PrometheusCollector,
) *SafeInvoiceService {
	if tracer == nil {
		tracer = otel.Tracer("lab07-observability")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &SafeInvoiceService{
		Deps:         deps,
		Logger:       logger,
		Tracer:       tracer,
		Collector:    collector,
	}
}

type dependencyStep struct {
	component string
	operation string
	execute   func(ctx context.Context) error
}

func (s *SafeInvoiceService) Process(ctx context.Context, invoiceID string) error {
	return s.ProcessWithDeps(ctx, invoiceID, s.Deps)
}

func (s *SafeInvoiceService) ProcessWithDeps(ctx context.Context, invoiceID string, deps Dependencies) error {
	ctx, rootSpan := s.Tracer.Start(ctx, "invoice.process",
		trace.WithAttributes(attribute.String("invoice_id", invoiceID)),
	)
	defer rootSpan.End()

	requestID := middleware.GetRequestID(ctx)

	steps := []dependencyStep{
		{
			component: "database",
			operation: "load_invoice",
			execute: func(ctx context.Context) error {
				return deps.Repo.Load(ctx, invoiceID)
			},
		},
		{
			component: "inventory",
			operation: "reserve",
			execute: func(ctx context.Context) error {
				return deps.Inventory.Reserve(ctx, invoiceID)
			},
		},
		{
			component: "commission",
			operation: "calculate",
			execute: func(ctx context.Context) error {
				return deps.Commission.Calculate(ctx, invoiceID)
			},
		},
		{
			component: "pdf",
			operation: "generate",
			execute: func(ctx context.Context) error {
				return deps.PDF.Generate(ctx, invoiceID)
			},
		},
		{
			component: "notification",
			operation: "send",
			execute: func(ctx context.Context) error {
				return deps.Notification.Send(ctx, invoiceID)
			},
		},
	}

	for _, step := range steps {
		stepSpanName := fmt.Sprintf("%s.%s", step.component, step.operation)
		stepCtx, stepSpan := s.Tracer.Start(ctx, stepSpanName,
			trace.WithAttributes(
				attribute.String("component", step.component),
				attribute.String("operation", step.operation),
			),
		)

		start := time.Now()
		err := step.execute(stepCtx)
		duration := time.Since(start)

		traceID := stepSpan.SpanContext().TraceID().String()
		spanID := stepSpan.SpanContext().SpanID().String()

		if err != nil {
			stepSpan.RecordError(err)
			stepSpan.SetStatus(codes.Error, err.Error())
			stepSpan.End()

			rootSpan.RecordError(err)
			rootSpan.SetStatus(codes.Error, err.Error())

			outcome := "error"
			if s.Collector != nil {
				s.Collector.DependencyDurationSeconds.WithLabelValues(step.component, step.operation, outcome).Observe(duration.Seconds())
				s.Collector.DependencyErrorsTotal.WithLabelValues(step.component, step.operation, outcome).Inc()
			}

			s.Logger.Error("dependency completed",
				"request_id", requestID,
				"trace_id", traceID,
				"span_id", spanID,
				"component", step.component,
				"operation", step.operation,
				"duration_ms", duration.Milliseconds(),
				"outcome", outcome,
				"error", err.Error(),
			)

			return fmt.Errorf("%s.%s: %w", step.component, step.operation, err)
		}

		stepSpan.SetStatus(codes.Ok, "")
		stepSpan.End()

		outcome := "success"
		if s.Collector != nil {
			s.Collector.DependencyDurationSeconds.WithLabelValues(step.component, step.operation, outcome).Observe(duration.Seconds())
		}

		s.Logger.Info("dependency completed",
			"request_id", requestID,
			"trace_id", traceID,
			"span_id", spanID,
			"component", step.component,
			"operation", step.operation,
			"duration_ms", duration.Milliseconds(),
			"outcome", outcome,
		)
	}

	rootSpan.SetStatus(codes.Ok, "")
	return nil
}

// HTTPHandler returns an http.Handler that instruments HTTP request telemetry.
func (s *SafeInvoiceService) HTTPHandler(resolveScenario func(r *http.Request) ScenarioConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const route = "/invoices/{id}/process"
		start := time.Now()

		if s.Collector != nil {
			s.Collector.HTTPInFlightRequests.WithLabelValues(route).Inc()
			defer s.Collector.HTTPInFlightRequests.WithLabelValues(route).Dec()
		}

		ctx, rootSpan := s.Tracer.Start(r.Context(), "http.request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.route", route),
			),
		)
		defer rootSpan.End()

		r = r.WithContext(ctx)
		requestID := middleware.GetRequestID(ctx)

		invoiceID := r.PathValue("id")
		if invoiceID == "" {
			invoiceID = r.URL.Query().Get("id")
		}

		if invoiceID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"error":"missing invoice id","request_id":"%s"}`, requestID)))

			duration := time.Since(start)
			s.recordHTTPMetrics(r.Method, route, "4xx", duration, false)
			s.logHTTPRequest(r.Method, route, http.StatusBadRequest, duration, requestID, rootSpan)
			return
		}

		// Resolve request-scoped dependencies
		deps := s.Deps
		if resolveScenario != nil {
			cfg := resolveScenario(r)
			deps = NewDependencies(cfg)
		}

		err := s.ProcessWithDeps(ctx, invoiceID, deps)
		duration := time.Since(start)

		if err != nil {
			statusClass := "5xx"
			s.recordHTTPMetrics(r.Method, route, statusClass, duration, true)
			s.logHTTPRequest(r.Method, route, http.StatusInternalServerError, duration, requestID, rootSpan)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"%s","request_id":"%s"}`, err.Error(), requestID)))
			return
		}

		statusClass := "2xx"
		s.recordHTTPMetrics(r.Method, route, statusClass, duration, false)
		s.logHTTPRequest(r.Method, route, http.StatusOK, duration, requestID, rootSpan)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"success","invoice_id":"%s","request_id":"%s"}`, invoiceID, requestID)))
	})
}

func (s *SafeInvoiceService) recordHTTPMetrics(method, route, statusClass string, duration time.Duration, isError bool) {
	if s.Collector == nil {
		return
	}
	s.Collector.HTTPRequestsTotal.WithLabelValues(method, route, statusClass).Inc()
	s.Collector.HTTPRequestDuration.WithLabelValues(method, route, statusClass).Observe(duration.Seconds())
	if isError {
		s.Collector.HTTPRequestErrorsTotal.WithLabelValues(method, route, statusClass).Inc()
	}
}

func (s *SafeInvoiceService) logHTTPRequest(method, route string, status int, duration time.Duration, requestID string, span trace.Span) {
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 {
		level = slog.LevelWarn
	}

	s.Logger.Log(context.Background(), level, "http request",
		"request_id", requestID,
		"trace_id", traceID,
		"span_id", spanID,
		"method", method,
		"route", route,
		"status", status,
		"duration_ms", duration.Milliseconds(),
	)
}
