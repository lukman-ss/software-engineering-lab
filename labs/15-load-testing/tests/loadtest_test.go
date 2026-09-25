package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest"
	"github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server"
)

func TestLoadTest_SmokeVsStress(t *testing.T) {
	srv := server.New(server.Config{
		MaxDBConnections: 2,
		DBQueryDuration:  10 * time.Millisecond,
	})
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	// Smoke test: 1 VU, zero queueing
	smokeCfg := loadtest.Config{
		URL:      ts.URL + "/booking",
		Method:   http.MethodPost,
		Body:     []byte(`{}`),
		VUs:      1,
		Duration: 500 * time.Millisecond,
	}
	smokeRes := loadtest.NewRunner(smokeCfg).Run(context.Background())

	if smokeRes.TotalRequests == 0 {
		t.Fatal("smoke test recorded zero requests")
	}
	if smokeRes.ErrorCount != 0 {
		t.Fatalf("expected 0 errors in smoke test, got %d", smokeRes.ErrorCount)
	}

	// Stress test: 10 VUs against 2 DB slots -> heavy queuing
	stressCfg := loadtest.Config{
		URL:      ts.URL + "/booking",
		Method:   http.MethodPost,
		Body:     []byte(`{}`),
		VUs:      10,
		Duration: 500 * time.Millisecond,
	}
	stressRes := loadtest.NewRunner(stressCfg).Run(context.Background())

	if stressRes.TotalRequests == 0 {
		t.Fatal("stress test recorded zero requests")
	}
	if stressRes.P95Latency <= smokeRes.P95Latency {
		t.Fatalf("expected stress P95 (%s) to exceed smoke P95 (%s)", stressRes.P95Latency, smokeRes.P95Latency)
	}
}

func TestServer_MethodNotAllowed(t *testing.T) {
	srv := server.New(server.Config{})
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/booking")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}
