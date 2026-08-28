package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("worker starting")

	// Wait for shutdown signal
	<-quit
	logger.Info("shutting down worker")

	// Give ongoing work time to complete
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	<-ctx.Done()

	logger.Info("worker stopped")
}
