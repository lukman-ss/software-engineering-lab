package observability

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

type SuccessResponse struct {
	Status    string `json:"status"`
	InvoiceID string `json:"invoice_id"`
	RequestID string `json:"request_id"`
}

func writeJSONError(w http.ResponseWriter, status int, msg, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:     msg,
		RequestID: requestID,
	})
}

func writeJSONSuccess(w http.ResponseWriter, invoiceID, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Status:    "success",
		InvoiceID: invoiceID,
		RequestID: requestID,
	})
}

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

	traceID := rootSpan.SpanContext().TraceID().String()

	logger := s.Logger.With(
		"request_id", requestID,
		"trace_id", traceID,
		"invoice_id", invoiceID,
	)

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

			logger.ErrorContext(stepCtx, "dependency completed",
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

		logger.InfoContext(stepCtx, "dependency completed",
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
		const rootSpanName = "POST /invoices/{id}/process"
		start := time.Now()

		if s.Collector != nil {
			s.Collector.HTTPInFlightRequests.WithLabelValues(route).Inc()
			defer s.Collector.HTTPInFlightRequests.WithLabelValues(route).Dec()
		}

		ctx, rootSpan := s.Tracer.Start(r.Context(), rootSpanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.HTTPRoute(route),
				semconv.URLPath(r.URL.Path),
				semconv.URLScheme(r.URL.Scheme),
				semconv.ServerAddress(r.Host),
			),
		)
		defer rootSpan.End()

		r = r.WithContext(ctx)
		requestID := middleware.GetRequestID(ctx)
		traceID := rootSpan.SpanContext().TraceID().String()
		spanID := rootSpan.SpanContext().SpanID().String()
		logger := ContextLogger(ctx, s.Logger, requestID)

		invoiceID := r.PathValue("id")
		if invoiceID == "" {
			invoiceID = r.URL.Query().Get("id")
		}

		if invoiceID == "" {
			duration := time.Since(start)
			rootSpan.SetStatus(codes.Error, "missing invoice id")
			rootSpan.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusBadRequest))
			s.recordHTTPMetrics(r.Method, route, "4xx", duration, false)
			logger.Warn("http request",
				"method", r.Method,
				"route", route,
				"status", http.StatusBadRequest,
				"duration_ms", duration.Milliseconds(),
			)
			writeJSONError(w, http.StatusBadRequest, "missing invoice id", requestID)
			return
		}

		rootSpan.SetAttributes(attribute.String("invoice.id", invoiceID))

		deps := s.Deps
		if resolveScenario != nil {
			cfg := resolveScenario(r)
			deps = NewDependencies(cfg)
		}

		err := s.ProcessWithDeps(ctx, invoiceID, deps)
		duration := time.Since(start)

		rootSpan.SetAttributes(semconv.HTTPResponseStatusCode(httpStatusFromError(err)))

		if err != nil {
			statusCode := httpStatusFromError(err)
			statusClass := statusClassFor(statusCode)
			rootSpan.RecordError(err)
			rootSpan.SetStatus(codes.Error, err.Error())

			s.recordHTTPMetrics(r.Method, route, statusClass, duration, true)
			logger.Error("http request",
				"method", r.Method,
				"route", route,
				"status", statusCode,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)

			writeJSONError(w, statusCode, "invoice processing failed", requestID)
			return
		}

		rootSpan.SetStatus(codes.Ok, "")
		s.recordHTTPMetrics(r.Method, route, "2xx", duration, false)
		logger.Info("http request",
			"method", r.Method,
			"route", route,
			"status", http.StatusOK,
			"duration_ms", duration.Milliseconds(),
		)
		_ = traceID
		_ = spanID

		writeJSONSuccess(w, invoiceID, requestID)
	})
}

func httpStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout
	}
	if errors.Is(err, context.Canceled) {
		return 499
	}
	return http.StatusBadGateway
}

func statusClassFor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	default:
		return "5xx"
	}
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
