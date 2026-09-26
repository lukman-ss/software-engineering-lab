package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"circuitbreaker/internal/checkout"
	"circuitbreaker/internal/circuitbreaker"
	"circuitbreaker/internal/payment"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("  LAB 14: CIRCUIT BREAKER PATTERN DEMONSTRATION   ")
	fmt.Println("==================================================")

	// Setup fake payment server
	// Client timeout: 100ms, Server slow delay: 250ms
	clientTimeout := 100 * time.Millisecond
	serverSlowDelay := 250 * time.Millisecond

	fakeServer := payment.NewFakeServer(serverSlowDelay)
	defer fakeServer.Close()

	paymentClient := payment.NewClient(fakeServer.URL(), clientTimeout)

	// ----------------------------------------------------
	// Scenario 1: Without Circuit Breaker (SLOW dependency)
	// ----------------------------------------------------
	fmt.Println("\n=== SCENARIO 1: WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY) ===")
	fakeServer.SetMode(payment.ModeSlow)
	fakeServer.ResetCount()

	svcNoCB := checkout.NewService(paymentClient, nil)
	for i := 1; i <= 3; i++ {
		start := time.Now()
		err := svcNoCB.Checkout(context.Background())
		dur := time.Since(start)

		res := "success"
		if err != nil {
			res = "timeout/error"
		}
		fmt.Printf("request=%d result=%-13s duration=%v\n", i, res, dur.Round(time.Millisecond))
	}
	fmt.Printf("downstream_calls=%d (all requests blocked and hit downstream)\n", fakeServer.RequestCount())

	// ----------------------------------------------------
	// Scenario 2: With Circuit Breaker (DOWN dependency & Fail-Fast)
	// ----------------------------------------------------
	fmt.Println("\n=== SCENARIO 2: WITH CIRCUIT BREAKER (FAIL-FAST ON DOWN DEPENDENCY) ===")
	fakeServer.SetMode(payment.ModeDown)
	fakeServer.ResetCount()

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,
		OpenTimeout:      300 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})
	svcWithCB := checkout.NewService(paymentClient, cb)

	for i := 1; i <= 6; i++ {
		start := time.Now()
		err := svcWithCB.Checkout(context.Background())
		dur := time.Since(start)

		res := "ok"
		if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			res = "circuit_open (fail-fast)"
		} else if err != nil {
			res = "payment_error"
		}
		fmt.Printf("request=%d result=%-26s state=%-9s duration=%v\n", i, res, cb.State(), dur)
	}
	fmt.Printf("downstream_calls=%d (downstream calls stopped once OPEN)\n", fakeServer.RequestCount())

	// ----------------------------------------------------
	// Scenario 3: Recovery (HALF_OPEN -> CLOSED)
	// ----------------------------------------------------
	fmt.Println("\n=== SCENARIO 3: RECOVERY (HALF_OPEN -> CLOSED) ===")
	fmt.Printf("initial state=%s\n", cb.State())
	fmt.Println("waiting for cooldown (300ms)...")
	time.Sleep(350 * time.Millisecond)

	// Heal the server
	fakeServer.SetMode(payment.ModeHealthy)
	fmt.Printf("payment server recovered to HEALTHY. Current CB state=%s\n", cb.State())

	fmt.Println("sending probe request...")
	err := svcWithCB.Checkout(context.Background())
	fmt.Printf("probe result: err=%v, state after probe=%s\n", err, cb.State())

	fmt.Println("sending subsequent regular request...")
	err = svcWithCB.Checkout(context.Background())
	fmt.Printf("subsequent result: err=%v, state=%s\n", err, cb.State())

	// ----------------------------------------------------
	// Scenario 4: Failed Recovery (HALF_OPEN -> OPEN)
	// ----------------------------------------------------
	fmt.Println("\n=== SCENARIO 4: FAILED RECOVERY (HALF_OPEN -> OPEN AGAIN) ===")
	// Make server DOWN again
	fakeServer.SetMode(payment.ModeDown)

	// Force circuit open with failures
	for i := 1; i <= 3; i++ {
		_ = svcWithCB.Checkout(context.Background())
	}
	fmt.Printf("circuit forced back to: %s\n", cb.State())

	fmt.Println("waiting for cooldown (300ms)...")
	time.Sleep(350 * time.Millisecond)

	fmt.Printf("dependency still DOWN. Current CB state=%s\n", cb.State())
	fmt.Println("sending probe request...")
	err = svcWithCB.Checkout(context.Background())
	fmt.Printf("probe result: err=%v, state after failed probe=%s\n", err != nil, cb.State())

	fmt.Println("sending next request while re-opened...")
	err = svcWithCB.Checkout(context.Background())
	fmt.Printf("next request result: err=%v, state=%s\n", err, cb.State())

	fmt.Println("\n==================================================")
	fmt.Println("  DEMO COMPLETED SUCCESSFULLY                     ")
	fmt.Println("==================================================")
}
