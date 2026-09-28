package main

import (
	"context"
	"fmt"

	"labs/29-saga-pattern/internal/saga"
	"labs/29-saga-pattern/internal/services"
)

func main() {
	ctx := context.Background()
	fmt.Println("=== Saga Pattern Demonstration ===")

	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"laptop": 1})

	// Scenario 1: Successful Orchestrated Saga
	fmt.Println("\n--- Scenario 1: Happy Path (Orchestrator) ---")
	orch1 := saga.NewOrchestrator()
	orderID1 := "ord-success"

	orch1.AddStep(saga.Step{
		Name: "CreateOrder",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [OrderService] Creating order:", orderID1)
			return orderSvc.CreateOrder(orderID1)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [OrderService] Cancelling order:", orderID1)
			return orderSvc.CancelOrder(orderID1)
		},
	})
	orch1.AddStep(saga.Step{
		Name: "ProcessPayment",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [PaymentService] Processing payment for:", orderID1)
			return paymentSvc.ProcessPayment(orderID1, 1500.0, false)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [PaymentService] Refunding payment for:", orderID1)
			return paymentSvc.RefundPayment(orderID1)
		},
	})
	orch1.AddStep(saga.Step{
		Name: "ReserveInventory",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [InventoryService] Reserving 1 laptop")
			return invSvc.Reserve("laptop", 1)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [InventoryService] Releasing 1 laptop")
			return invSvc.Release("laptop", 1)
		},
	})
	orch1.AddStep(saga.Step{
		Name: "ApproveOrder",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [OrderService] Finalizing order approval")
			return orderSvc.ApproveOrder(orderID1)
		},
	})

	err := orch1.Execute(ctx)
	fmt.Printf("Scenario 1 Result: error=%v, OrderState=%s, Remaining Stock=%d\n",
		err, orderSvc.GetOrder(orderID1), invSvc.GetStock("laptop"))

	// Scenario 2: Failure with Rollback (Inventory Out of Stock)
	fmt.Println("\n--- Scenario 2: Rollback on Failure (Orchestrator) ---")
	orch2 := saga.NewOrchestrator()
	orderID2 := "ord-failed"

	orch2.AddStep(saga.Step{
		Name: "CreateOrder",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [OrderService] Creating order:", orderID2)
			return orderSvc.CreateOrder(orderID2)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [OrderService] Compensating: Cancelling order:", orderID2)
			return orderSvc.CancelOrder(orderID2)
		},
	})
	orch2.AddStep(saga.Step{
		Name: "ProcessPayment",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [PaymentService] Processing payment for:", orderID2)
			return paymentSvc.ProcessPayment(orderID2, 1500.0, false)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [PaymentService] Compensating: Refunding payment for:", orderID2)
			return paymentSvc.RefundPayment(orderID2)
		},
	})
	orch2.AddStep(saga.Step{
		Name: "ReserveInventory",
		Execute: func(ctx context.Context) error {
			fmt.Println(" -> [InventoryService] Reserving 1 laptop (current stock:", invSvc.GetStock("laptop"), ")")
			return invSvc.Reserve("laptop", 1)
		},
		Compensate: func(ctx context.Context) error {
			fmt.Println(" <- [InventoryService] Releasing 1 laptop")
			return invSvc.Release("laptop", 1)
		},
	})

	err2 := orch2.Execute(ctx)
	fmt.Printf("Scenario 2 Result: error=%v, OrderState=%s, HasPayment=%v, Stock=%d\n",
		err2, orderSvc.GetOrder(orderID2), paymentSvc.HasPayment(orderID2), invSvc.GetStock("laptop"))

	fmt.Println("\n=== Demo Completed Successfully ===")
}
