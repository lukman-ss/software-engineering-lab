package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"zero-downtime-deployment/internal/server"
)

func waitForServerReady(addr string) error {
	deadline := time.Now().Add(2 * time.Second)
	url := fmt.Sprintf("http://%s/healthz/live", addr)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fmt.Errorf("server at %s did not start within timeout", addr)
}

func TestServerProbes(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8081", 0)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8081"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}

	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	resp, err := http.Get("http://127.0.0.1:8081/healthz/live")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected liveness 200 OK, got %v, err: %v", resp.StatusCode, err)
	}

	resp, err = http.Get("http://127.0.0.1:8081/healthz/ready")
	if err != nil || resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected readiness 503, got %v, err: %v", resp.StatusCode, err)
	}

	srv.SetReady(true)

	resp, err = http.Get("http://127.0.0.1:8081/healthz/ready")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected readiness 200 OK, got %v, err: %v", resp.StatusCode, err)
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8082", 0)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8082"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}

	done := make(chan bool)
	go func() {
		resp, err := http.Get("http://127.0.0.1:8082/work?d=100ms")
		if err != nil {
			t.Errorf("expected request to succeed, got %v", err)
			done <- false
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "WORK COMPLETED" {
			t.Errorf("unexpected body: %s", string(body))
		}
		done <- true
	}()

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	success := <-done
	if !success {
		t.Fatalf("in-flight request failed during graceful shutdown")
	}
}

func TestServerPreStopHook(t *testing.T) {
	preStopDuration := 100 * time.Millisecond
	srv := server.NewServer("127.0.0.1:8083", preStopDuration)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8083"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < preStopDuration {
		t.Fatalf("expected preStop delay of at least %v, but took %v", preStopDuration, elapsed)
	}
}

func TestServerPreStopContextCancellation(t *testing.T) {
	preStopDuration := 500 * time.Millisecond
	srv := server.NewServer("127.0.0.1:8084", preStopDuration)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8084"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err == nil {
		t.Fatalf("expected error due to context cancellation during preStop, got nil")
	}

	elapsed := time.Since(start)
	if elapsed >= preStopDuration {
		t.Fatalf("shutdown blocked for full preStop (%v) instead of aborting promptly on context cancellation (%v)", preStopDuration, elapsed)
	}
}

func TestServerInvalidDurationFallback(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8086", 0)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8086"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	start := time.Now()
	resp, err := http.Get("http://127.0.0.1:8086/work?d=INVALID")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if elapsed < 50*time.Millisecond {
		t.Fatalf("expected fallback delay >=50ms, got %v", elapsed)
	}
}

func TestServerReadyUnreadyTransition(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8087", 0)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8087"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	srv.SetReady(true)
	resp, err := http.Get("http://127.0.0.1:8087/healthz/ready")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when ready, got %v err=%v", resp.StatusCode, err)
	}

	srv.SetReady(false)
	resp, err = http.Get("http://127.0.0.1:8087/healthz/ready")
	if err != nil || resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 after unready, got %v err=%v", resp.StatusCode, err)
	}
}

func TestServerMultiRequestDrain(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8088", 0)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8088"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}

	const n = 3
	results := make([]string, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resp, err := http.Get("http://127.0.0.1:8088/work?d=100ms")
			if err != nil {
				results[idx] = err.Error()
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			results[idx] = string(body)
		}(i)
	}

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	wg.Wait()

	for i, res := range results {
		if res != "WORK COMPLETED" {
			t.Errorf("request %d did not complete successfully during drain, got %q", i, res)
		}
	}
}

func TestServerWorkRequestCancellation(t *testing.T) {
	srv := server.NewServer("127.0.0.1:8085", 0)
	srv.SetReady(true)

	go func() {
		_ = srv.Start()
	}()
	if err := waitForServerReady("127.0.0.1:8085"); err != nil {
		t.Fatalf("failed waiting for server: %v", err)
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	reqCtx, reqCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer reqCancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "http://127.0.0.1:8085/work?d=200ms", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err == nil {
		t.Fatalf("expected request to fail due to client context cancellation")
	}

	time.Sleep(50 * time.Millisecond)
	if active := srv.ActiveRequests(); active != 0 {
		t.Fatalf("expected active requests to be 0 after cancellation, got %d", active)
	}
}
