package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukman/software-engineer-lab/internal/inventory"
	"github.com/lukman/software-engineer-lab/internal/notification"
	"github.com/lukman/software-engineer-lab/internal/order"
	"github.com/lukman/software-engineer-lab/internal/payment"
	"github.com/lukman/software-engineer-lab/internal/wallet"
	apperrors "github.com/lukman/software-engineer-lab/pkg/errors"
	"github.com/lukman/software-engineer-lab/pkg/database"
	"github.com/lukman/software-engineer-lab/pkg/middleware"
	"github.com/lukman/software-engineer-lab/pkg/validation"
)

type app struct {
	orderSvc      order.Service
	paymentSvc    payment.Service
	inventorySvc  inventory.Service
	walletSvc     wallet.Service
	notifySvc     notification.Service
	db            *sql.DB
	logger        *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	// Database
	cfg := database.FromEnv()
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connected")

	// Initialize repositories
	orderRepo := order.NewPGRepo(db)
	inventoryRepo := inventory.NewPGRepo(db)
	paymentRepo := payment.NewPGRepo(db)
	walletRepo := wallet.NewPGRepo(db)
	notifyRepo := notification.NewPGRepo(db)

	// Initialize services
	inventorySvc := inventory.NewService(inventoryRepo)
	paymentSvc := payment.NewService(paymentRepo)
	walletSvc := wallet.NewService(walletRepo)
	notifySvc := notification.NewService(notifyRepo)
	orderSvc := order.NewService(orderRepo, inventorySvc, paymentSvc, notifySvc)

	// HTTP server
	mux := http.NewServeMux()
	handler := &app{
		orderSvc:      orderSvc,
		paymentSvc:    paymentSvc,
		inventorySvc:  inventorySvc,
		walletSvc:     walletSvc,
		notifySvc:     notifySvc,
		db:            db,
		logger:        logger,
	}

	mux.HandleFunc("GET /health", handler.health)
	mux.HandleFunc("GET /ready", handler.ready)
	mux.HandleFunc("POST /orders", handler.createOrder)
	mux.HandleFunc("GET /orders/{id}", handler.getOrder)
	mux.HandleFunc("POST /orders/{id}/pay", handler.payOrder)

	// Apply middleware chain: RequestID -> Logger
	wrapped := middleware.RequestID(middleware.Logger(logger)(mux))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           wrapped,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("shutdown signal received, stopping server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown failed", "error", err)
		}
	}()

	logger.Info("server starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (a *app) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.db.PingContext(ctx); err != nil {
		a.logger.Error("readiness check failed", "error", err)
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func (a *app) createOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Items  []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
			UnitPrice int64  `json:"unit_price"`
		} `json:"items"`
	}

	if err := validation.ReadJSON(r, &req); err != nil {
		a.respondError(w, r, err)
		return
	}

	if req.UserID == "" {
		a.respondError(w, r, validation.MissingField("user_id"))
		return
	}
	if len(req.Items) == 0 {
		a.respondError(w, r, validation.InvalidValue("items", "must have at least one item"))
		return
	}

	items := make([]order.OrderItem, len(req.Items))
	for i, item := range req.Items {
		if item.Quantity <= 0 || item.UnitPrice <= 0 {
			a.respondError(w, r, validation.InvalidValue("items", "quantity and unit_price must be positive"))
			return
		}
		items[i] = order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ord, err := a.orderSvc.CreateOrder(ctx, req.UserID, items)
	if err != nil {
		a.respondError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ord)
}

func (a *app) getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		a.respondError(w, r, validation.MissingField("id"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ord, err := a.orderSvc.GetByID(ctx, id)
	if err != nil {
		a.respondError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ord)
}

func (a *app) payOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		a.respondError(w, r, validation.MissingField("id"))
		return
	}

	var req struct {
		IdempotencyKey string `json:"idempotency_key"`
		Method         string `json:"method"`
	}
	if err := validation.ReadJSON(r, &req); err != nil {
		a.respondError(w, r, err)
		return
	}

	if req.IdempotencyKey == "" {
		a.respondError(w, r, validation.MissingField("idempotency_key"))
		return
	}
	if req.Method == "" {
		req.Method = "card"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ord, err := a.orderSvc.GetByID(ctx, id)
	if err != nil {
		a.respondError(w, r, err)
		return
	}

	if ord.Status != order.StatusPending {
		a.respondError(w, r, apperrors.ErrInvalidStatus)
		return
	}

	payCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	p, err := a.paymentSvc.ProcessPayment(payCtx, req.IdempotencyKey, id, ord.TotalAmount, req.Method)
	if err != nil {
		a.respondError(w, r, err)
		return
	}

	if err := a.orderSvc.MarkAsPaid(ctx, id); err != nil {
		a.respondError(w, r, err)
		return
	}

	_ = a.notifySvc.NotifyPaymentResult(ctx, ord.UserID, id, true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// respondError writes a structured error response and logs server-side details.
func (a *app) respondError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := middleware.GetRequestID(r.Context())

	appErr := apperrors.FromError(err)

	// Log server-side error with full context (never expose to client)
	a.logger.Error("request error",
		"request_id", requestID,
		"path", r.URL.Path,
		"method", r.Method,
		"error_code", appErr.Code,
		"error_category", appErr.Category,
		"error", appErr.Error(),
	)

	// Return structured error to client (no internal details)
	status := apperrors.HTTPStatus(appErr.Category)
	resp := map[string]any{
		"error": map[string]any{
			"code":       appErr.Code,
			"message":    appErr.Message,
			"request_id": requestID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
