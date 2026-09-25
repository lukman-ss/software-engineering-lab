package main

import (
	"fmt"
	"lab16/internal/di"
)

type SimpleContainer struct {
	gateway di.PaymentGateway
}

func (c *SimpleContainer) GetPaymentGateway() di.PaymentGateway {
	return c.gateway
}

func main() {
	// Finding 1: Separation of Configuration from Use
	realGateway := &di.RealGateway{}

	// Finding 3: Constructor Injection
	processor := di.NewProcessor(realGateway)
	fmt.Println("--- Running Constructor Injection ---")
	err := processor.ProcessPayment(100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Finding 4: Service Locator (Anti-Pattern)
	container := &SimpleContainer{gateway: realGateway}
	badProcessor := di.NewBadProcessor(container)
	fmt.Println("--- Running Service Locator ---")
	err = badProcessor.ProcessPayment(200)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
