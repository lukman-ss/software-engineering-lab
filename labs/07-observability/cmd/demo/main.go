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
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

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
	tp := sdktrace.NewTracerProvider()
	defer tp.Shutdown(context.Background())
	tracer := tp.Tracer("lab07-observability-demo")

	// Safe Service
	safeRepo, safeInv, safeComm, safePDF, safeNotif := observability.NewFakeDependencies(observability.ScenarioConfig{})
	safeService := observability.NewSafeInvoiceService(
		safeRepo, safeInv, safeComm, safePDF, safeNotif,
		logger, tracer, collector,
	)

	// Unsafe Service
	unsafeRepo, unsafeInv, unsafeComm, unsafePDF, unsafeNotif := observability.NewFakeDependencies(observability.ScenarioConfig{})
	unsafeService := &observability.UnsafeInvoiceService{
		Repo:         unsafeRepo,
		Inventory:    unsafeInv,
		Commission:   unsafeComm,
		PDF:          unsafePDF,
		Notification: unsafeNotif,
		LogWriter:    os.Stdout,
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
