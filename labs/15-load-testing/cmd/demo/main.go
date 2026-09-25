package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"text/tabwriter"
	"time"

	"github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest"
	"github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server"
)

func main() {
	// 1. Setup server with 5 max DB connections
	cfg := server.Config{
		MaxDBConnections: 5,
		DBQueryDuration:  20 * time.Millisecond,
	}
	srv := server.New(cfg)
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	fmt.Println("Starting Load Test Demo (Booking Bengkel)")
	fmt.Printf("Server simulated DB connections: %d\n", cfg.MaxDBConnections)
	fmt.Printf("Simulated DB query duration: %s\n\n", cfg.DBQueryDuration)

	testDuration := 2 * time.Second

	// 2. Smoke Test (low VUs, no queuing)
	smokeCfg := loadtest.Config{
		URL:         ts.URL + "/booking",
		Method:      http.MethodPost,
		Body:        []byte(`{"vehicle_id":"V123","service":"oil_change"}`),
		ContentType: "application/json",
		VUs:         2, // Below DB max connections
		Duration:    testDuration,
	}
	fmt.Println("--- Running Smoke Test (2 VUs) ---")
	smokeRunner := loadtest.NewRunner(smokeCfg)
	smokeRes := smokeRunner.Run(context.Background())
	printResults(smokeRes)

	// 3. Stress Test (high VUs, queueing expected)
	stressCfg := loadtest.Config{
		URL:         ts.URL + "/booking",
		Method:      http.MethodPost,
		Body:        []byte(`{"vehicle_id":"V123","service":"oil_change"}`),
		ContentType: "application/json",
		VUs:         50, // 10x DB max connections, will cause queuing
		Duration:    testDuration,
	}
	fmt.Println("\n--- Running Stress Test (50 VUs) ---")
	stressRunner := loadtest.NewRunner(stressCfg)
	stressRes := stressRunner.Run(context.Background())
	printResults(stressRes)
}

func printResults(res loadtest.Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Total Requests:\t%d\n", res.TotalRequests)
	fmt.Fprintf(w, "Success:\t%d\n", res.SuccessCount)
	fmt.Fprintf(w, "Errors:\t%d\n", res.ErrorCount)
	fmt.Fprintf(w, "RPS:\t%.2f\n", res.RPS)
	fmt.Fprintf(w, "Average:\t%s\n", res.AvgLatency)
	fmt.Fprintf(w, "P50:\t%s\n", res.P50Latency)
	fmt.Fprintf(w, "P95:\t%s\n", res.P95Latency)
	fmt.Fprintf(w, "P99:\t%s\n", res.P99Latency)
	w.Flush()
}
