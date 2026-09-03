package main

import (
	"context"
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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317" // Default Jaeger OTLP gRPC port
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "lab07-observability"
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

func resolveScenario(r *http.Request) observability.ScenarioConfig {
	scenario := r.URL.Query().Get("scenario")
	switch scenario {
	case "slow-pdf":
		return observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          4800 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "slow-database":
		return observability.ScenarioConfig{
			RepoDelay:         3000 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "slow-external":
		return observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 3500 * time.Millisecond,
		}
	case "pdf-error":
		return observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFErr:            observability.ErrPDFGeneration,
			NotificationDelay: 20 * time.Millisecond,
		}
	case "normal":
		fallthrough
	default:
		return observability.ScenarioConfig{
			RepoDelay:         20 * time.Millisecond,
			InventoryDelay:    15 * time.Millisecond,
			CommissionDelay:   10 * time.Millisecond,
			PDFDelay:          30 * time.Millisecond,
			NotificationDelay: 20 * time.Millisecond,
		}
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Prometheus registry
	reg := prometheus.NewRegistry()
	collector := observability.NewPrometheusCollector(reg)

	// OTel TracerProvider
	tp, err := initTracerProvider(context.Background())
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
	tracer := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))

	// Safe Service
	safeService := observability.NewSafeInvoiceService(
		observability.NewDependencies(observability.ScenarioConfig{}),
		logger, tracer, collector,
	)

	// Unsafe Service
	unsafeService := &observability.UnsafeInvoiceService{
		Deps:      observability.NewDependencies(observability.ScenarioConfig{}),
		LogWriter: os.Stdout,
	}

	mux := http.NewServeMux()

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
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Printf("Lab 07 Observability demo server listening on :%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
	}
}
