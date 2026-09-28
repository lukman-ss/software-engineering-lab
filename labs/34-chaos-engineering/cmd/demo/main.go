package main

import (
	"context"
	"fmt"
	"time"

	"labs/34-chaos-engineering/internal/circuitbreaker"
	"labs/34-chaos-engineering/internal/experiment"
	"labs/34-chaos-engineering/internal/fault"
	"labs/34-chaos-engineering/internal/monitor"
)

func main() {
	fmt.Println("=== Chaos Engineering & Fault Injection Demo ===")

	injector := fault.NewInjector()
	mon := monitor.NewMonitor(0.25) // Max 25% error rate allowed
	cb := circuitbreaker.NewCircuitBreaker(3, 200*time.Millisecond)

	// Simulated downstream payment service
	callPaymentService := func(ctx context.Context) error {
		if err := injector.Execute(ctx); err != nil {
			return err
		}
		return nil
	}

	// Resilient client with circuit breaker & graceful degradation
	clientRequest := func(ctx context.Context) (string, error) {
		var resp string
		err := cb.Execute(func() error {
			if err := callPaymentService(ctx); err != nil {
				return err
			}
			resp = "Payment Success (Primary)"
			return nil
		}, func() error {
			// Graceful fallback
			resp = "Payment Queued (Fallback)"
			return nil
		})

		if err != nil {
			mon.RecordFailure()
			return "", err
		}
		mon.RecordSuccess()
		return resp, nil
	}

	fmt.Println("\n[1] Baseline Steady-State Traffic (Normal Operation)...")
	for i := 1; i <= 5; i++ {
		resp, _ := clientRequest(context.Background())
		fmt.Printf("Req #%d -> %s | CB State: %s\n", i, resp, cb.State())
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Println("\n[2] Running Chaos Experiment: Injecting Downstream Failures with Resilient Client...")
	expCfg := experiment.Config{
		Name:            "Downstream-Fault-Injection",
		Duration:        300 * time.Millisecond,
		ForceError:      true,
		MonitorInterval: 40 * time.Millisecond,
	}
	exp := experiment.NewExperiment(expCfg, injector, mon)

	go func() {
		_ = exp.Run(context.Background())
	}()

	time.Sleep(20 * time.Millisecond)

	for i := 6; i <= 15; i++ {
		resp, err := clientRequest(context.Background())
		if err != nil {
			fmt.Printf("Req #%d -> ERROR: %v | CB State: %s\n", i, err, cb.State())
		} else {
			fmt.Printf("Req #%d -> %s | CB State: %s\n", i, resp, cb.State())
		}
		time.Sleep(30 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("\nExperiment State: %s\n", exp.State())
	metrics := mon.Metrics()
	fmt.Printf("Steady-State Metrics with Fallback: Total=%d, Success=%d, Failed=%d, ErrorRate=%.2f%%\n",
		metrics.TotalRequests, metrics.SuccessRequests, metrics.FailedRequests, mon.ErrorRate()*100)

	fmt.Println("\n[3] Running Second Chaos Experiment: Unmitigated Downstream Fault (Triggers Auto-Abort)...")
	unmitigatedMon := monitor.NewMonitor(0.20) // 20% limit
	for i := 0; i < 4; i++ {
		unmitigatedMon.RecordSuccess()
	}

	abortExpCfg := experiment.Config{
		Name:            "Unmitigated-Fault-Auto-Abort",
		Duration:        500 * time.Millisecond,
		ForceError:      true,
		MonitorInterval: 20 * time.Millisecond,
	}
	abortExp := experiment.NewExperiment(abortExpCfg, injector, unmitigatedMon)

	go func() {
		_ = abortExp.Run(context.Background())
	}()

	time.Sleep(10 * time.Millisecond)
	// Simulate raw service calls failing directly without graceful fallback
	for i := 1; i <= 3; i++ {
		err := callPaymentService(context.Background())
		if err != nil {
			unmitigatedMon.RecordFailure()
			fmt.Printf("Raw Call #%d Failed: %v\n", i, err)
		}
		time.Sleep(30 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Experiment State: %s\n", abortExp.State())
	if abortExp.AbortReason() != "" {
		fmt.Printf("Abort Reason: %s\n", abortExp.AbortReason())
	}
	fmt.Printf("Fault Injector Active: %v (Neutralized on Abort)\n", injector.IsEnabled())

	fmt.Println("\n[4] Traffic After Experiment Ends / Recovery...")
	time.Sleep(200 * time.Millisecond)
	for i := 16; i <= 20; i++ {
		resp, _ := clientRequest(context.Background())
		fmt.Printf("Req #%d -> %s | CB State: %s\n", i, resp, cb.State())
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Println("\n=== Demo Complete ===")
}
