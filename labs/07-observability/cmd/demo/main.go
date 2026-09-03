package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/07-observability"
	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

func resolveServiceName() string {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		return "lab07-observability"
	}
	return serviceName
}

func initTracerProvider(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317" // Default Jaeger OTLP gRPC port
	}

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func resolveScenario(r *http.Request, baseDeps observability.Dependencies) observability.Dependencies {
	scenario := r.URL.Query().Get("scenario")
	var cfg observability.ScenarioConfig
	switch scenario {
	case "slow-pdf":
		cfg = observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          4800 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "slow-database":
		cfg = observability.ScenarioConfig{
			RepoDelay:         3000 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "slow-external":
		cfg = observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 3500 * time.Millisecond,
		}
	case "pdf-error":
		cfg = observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFErr:            observability.ErrPDFGeneration,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "normal":
		fallthrough
	default:
		cfg = observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	}
	deps := observability.NewDependencies(cfg)
	if baseDeps.Notification != nil {
		if httpNotif, isHTTP := baseDeps.Notification.(*observability.HTTPNotificationClient); isHTTP {
			baseURL := httpNotif.BaseURL
			if scenario == "slow-external" {
				baseURL = baseURL + "?scenario=slow-external"
			}
			deps.Notification = observability.NewHTTPNotificationClient(baseURL, httpNotif.HTTPClient)
		} else if baseDeps.Notification != nil {
			deps.Notification = baseDeps.Notification
		}
	}
	return deps
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Prometheus registry
	reg := prometheus.NewRegistry()
	collector := observability.NewPrometheusCollector(reg)

	serviceName := resolveServiceName()

	// OTel TracerProvider
	tp, err := initTracerProvider(context.Background(), serviceName)
	if err != nil {
		logger.Error("failed to initialize tracer provider", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("tracer provider shutdown failed", "error", err)
		}
	}()
	tracer := otel.Tracer(serviceName)

	// Notification client endpoint for simulated downstream service
	notifURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	var notifClient observability.NotificationService
	if notifURL != "" {
		notifClient = observability.NewHTTPNotificationClient(notifURL, nil)
	}

	// Safe Service
	defaultSafeDeps := observability.NewDependencies(observability.ScenarioConfig{})
	if notifClient != nil {
		defaultSafeDeps.Notification = notifClient
	}
	safeService := observability.NewSafeInvoiceService(
		defaultSafeDeps,
		logger, tracer, collector,
	)

	// Unsafe Service
	defaultUnsafeDeps := observability.NewDependencies(observability.ScenarioConfig{})
	if notifClient != nil {
		defaultUnsafeDeps.Notification = notifClient
	}
	unsafeService := &observability.UnsafeInvoiceService{
		Deps:      defaultUnsafeDeps,
		LogWriter: os.Stdout,
	}

	mux := http.NewServeMux()

	// Simulated downstream notification service endpoint (service boundary simulation in same process)
	mux.HandleFunc("POST /notifications", func(w http.ResponseWriter, r *http.Request) {
		extractedCtx := otel.GetTextMapPropagator().Extract(
			r.Context(),
			propagation.HeaderCarrier(r.Header),
		)
		ctx, span := tracer.Start(extractedCtx, "POST /notifications",
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.URLPath(r.URL.Path),
			),
		)
		defer span.End()

		invoiceID := r.URL.Query().Get("invoice_id")
		scenario := r.URL.Query().Get("scenario")
		if scenario == "slow-external" {
			delay := 3500 * time.Millisecond
			timer := time.NewTimer(delay)
			defer timer.Stop()
			select {
			case <-timer.C:
			case <-ctx.Done():
				span.RecordError(ctx.Err())
				span.SetStatus(codes.Error, ctx.Err().Error())
				span.SetAttributes(semconv.HTTPResponseStatusCode(499))
				w.WriteHeader(499)
				return
			}
		}

		span.SetStatus(codes.Ok, "")
		span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusOK))

		requestID := middleware.GetRequestID(ctx)
		logger := observability.ContextLogger(ctx, logger, requestID)
		logger.InfoContext(ctx, "received notification webhook",
			"invoice_id", invoiceID,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "delivered"}); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.ErrorContext(ctx, "failed to write notification response",
				"operation", "encode_response",
				"error", err.Error(),
			)
		}
	})

	// Metrics endpoint
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// Safe Route
	mux.Handle("POST /invoices/{id}/process", safeService.HTTPHandler(resolveScenario))

	// Unsafe Route
	mux.Handle("POST /unsafe/invoices/{id}/process", unsafeService.HTTPHandler(resolveScenario))

	// RequestID middleware wrapped
	handler := middleware.RequestID(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		logger.Info("shutting down server", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown failed", "error", err)
		}
	}()

	fmt.Printf("Lab 07 Observability demo server listening on :%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
	}
}
