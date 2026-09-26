package main

import (
	"circuitbreaker/internal/circuitbreaker"
	"circuitbreaker/internal/checkout"
	"circuitbreaker/internal/payment"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== SCENARIO 1: WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY) ===")
	scenario1()

	fmt.Println("\n=== SCENARIO 2: WITH CIRCUIT BREAKER (FAIL-FAST ON DOWN DEPENDENCY) ===")
	scenario2()

	fmt.Println("\n=== SCENARIO 3: RECOVERY (HALF_OPEN -> CLOSED) ===")
	scenario3()

	fmt.Println("\n=== SCENARIO 4: FAILED RECOVERY (HALF_OPEN -> OPEN AGAIN) ===")
	scenario4()
}

func scenario1() {
	server := payment.NewFakeServer(200 * time.Millisecond)
	defer server.Close()
	server.SetMode(payment.ModeSlow)

	client := payment.NewClient(server.URL(), 100*time.Millisecond)
	service := checkout.New(client, nil)

	downstream := server.RequestCount()
	for i := 1; i <= 3; i++ {
		result := service.CheckoutWithoutBreaker()
		fmt.Printf("request=%d result=%s\n", i, result)
	}
	newDownstream := server.RequestCount()
	fmt.Printf("\ndownstream_calls=%d\n", newDownstream-downstream)
}

func scenario2() {
	server := payment.NewFakeServer(100 * time.Millisecond)
	defer server.Close()
	server.SetMode(payment.ModeDown)

	breaker := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,
		OpenTimeout:      300 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})
	client := payment.NewClient(server.URL(), 100*time.Millisecond)
	service := checkout.New(client, breaker)

	downstream := server.RequestCount()
	for i := 1; i <= 6; i++ {
		result := service.CheckoutWithBreaker()
		fmt.Printf("request=%d result=%s\n", i, result)
	}
	newDownstream := server.RequestCount()
	fmt.Printf("\ndownstream_calls=%d\n", newDownstream-downstream)
}

func scenario3() {
	server := payment.NewFakeServer(100 * time.Millisecond)
	defer server.Close()
	server.SetMode(payment.ModeDown)

	breaker := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 2,
		OpenTimeout:      300 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})
	client := payment.NewClient(server.URL(), 100*time.Millisecond)
	service := checkout.New(client, breaker)

	for i := 0; i < 2; i++ {
		_ = service.CheckoutWithBreaker()
	}
	fmt.Printf("initial state=%s\n", breaker.State())

	fmt.Println("waiting for cooldown (300ms)...")
	time.Sleep(320 * time.Millisecond)
	server.SetMode(payment.ModeHealthy)
	fmt.Printf("payment server recovered to HEALTHY. Current CB state=%s\n", breaker.State())

	fmt.Println("sending probe request...")
	result := service.CheckoutWithBreaker()
	fmt.Printf("probe result: err=%v, state after probe=%s\n", result.Err, result.State)

	fmt.Println("sending subsequent regular request...")
	result = service.CheckoutWithBreaker()
	fmt.Printf("subsequent result: err=%v, state=%s\n", result.Err, result.State)
}

func scenario4() {
	server := payment.NewFakeServer(100 * time.Millisecond)
	defer server.Close()
	server.SetMode(payment.ModeDown)

	breaker := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 2,
		OpenTimeout:      300 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})
	client := payment.NewClient(server.URL(), 100*time.Millisecond)
	service := checkout.New(client, breaker)

	for i := 0; i < 2; i++ {
		_ = service.CheckoutWithBreaker()
	}
	fmt.Printf("circuit forced back to: %s\n", breaker.State())

	fmt.Println("waiting for cooldown (300ms)...")
	time.Sleep(320 * time.Millisecond)
	server.SetMode(payment.ModeDown)
	fmt.Printf("dependency still DOWN. Current CB state=%s\n", breaker.State())

	fmt.Println("sending probe request...")
	result := service.CheckoutWithBreaker()
	fmt.Printf("probe result: err=%v, state after failed probe=%s\n", result.Err, result.State)

	fmt.Println("sending next request while re-opened...")
	result = service.CheckoutWithBreaker()
	fmt.Printf("next request result: err=%v, state=%s\n", result.Err, result.State)
}
