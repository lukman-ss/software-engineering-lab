package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zero-downtime-deployment/internal/server"
	"zero-downtime-deployment/internal/worker"
)

func main() {
	log.Println("Starting Zero-Downtime Deployment Demo")

	// Start worker
	w := worker.NewWorker(100)
	w.Start(2)
	log.Println("Background worker started")

	w.Enqueue(worker.Job{ID: "DemoJob-1", Duration: 2 * time.Second})

	// Start server with a 1-second preStop hook
	srv := server.NewServer("127.0.0.1:8080", 1*time.Second)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Simulate initialization time (DB connection, cache warming)
	time.Sleep(1 * time.Second)
	srv.SetReady(true)
	log.Println("Application is ready to receive traffic (Readiness check passes)")

	// Simulate incoming HTTP traffic that takes time to process
	go func() {
		resp, err := http.Get("http://127.0.0.1:8080/work?d=2s")
		if err != nil {
			log.Printf("Client request error: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Client request completed with status: %v", resp.StatusCode)
	}()

	// Simulate an orchestrator sending SIGTERM during a deployment
	time.Sleep(500 * time.Millisecond)
	log.Println("Simulating SIGTERM from orchestrator (e.g. Kubernetes)")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// In a real app we wait for <-sig, here we mock sending it to ourselves
	go func() {
		sig <- syscall.SIGTERM
	}()

	<-sig
	log.Println("SIGTERM received, initiating graceful shutdown procedures")

	// Shutdown server gracefully (includes preStop hook and connection draining)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Stop background workers gracefully
	w.Stop()

	log.Println("Demo finished cleanly. Zero downtime achieved.")
}
