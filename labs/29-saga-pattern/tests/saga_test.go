package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"labs/29-saga-pattern/internal/saga"
	"labs/29-saga-pattern/internal/services"
)

func TestOrchestrator_HappyPath(t *testing.T) {
	ctx := context.Background()
	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"item-1": 10})

	orch := saga.NewOrchestrator()
	orderID := "ord-123"

	orch.AddStep(saga.Step{
		Name: "CreateOrder",
		Execute: func(ctx context.Context) error {
			return orderSvc.CreateOrder(orderID)
		},
		Compensate: func(ctx context.Context) error {
			return orderSvc.CancelOrder(orderID)
		},
	})

	orch.AddStep(saga.Step{
		Name: "ProcessPayment",
		Execute: func(ctx context.Context) error {
			return paymentSvc.ProcessPayment(orderID, 100.0, false)
		},
		Compensate: func(ctx context.Context) error {
			return paymentSvc.RefundPayment(orderID)
		},
	})

	orch.AddStep(saga.Step{
		Name: "ReserveInventory",
		Execute: func(ctx context.Context) error {
			return invSvc.Reserve("item-1", 2)
		},
		Compensate: func(ctx context.Context) error {
			return invSvc.Release("item-1", 2)
		},
	})

	orch.AddStep(saga.Step{
		Name: "ApproveOrder",
		Execute: func(ctx context.Context) error {
			return orderSvc.ApproveOrder(orderID)
		},
	})

	if err := orch.Execute(ctx); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if orderSvc.GetOrder(orderID) != services.OrderApproved {
		t.Fatalf("expected order approved, got %v", orderSvc.GetOrder(orderID))
	}
	if !paymentSvc.HasPayment(orderID) {
		t.Fatal("expected payment to exist")
	}
	if invSvc.GetStock("item-1") != 8 {
		t.Fatalf("expected stock 8, got %d", invSvc.GetStock("item-1"))
	}
}

func TestOrchestrator_FailureCompensatesLIFO(t *testing.T) {
	ctx := context.Background()
	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"item-1": 1})

	orch := saga.NewOrchestrator()
	orderID := "ord-fail"

	orch.AddStep(saga.Step{
		Name: "CreateOrder",
		Execute: func(ctx context.Context) error {
			return orderSvc.CreateOrder(orderID)
		},
		Compensate: func(ctx context.Context) error {
			return orderSvc.CancelOrder(orderID)
		},
	})

	orch.AddStep(saga.Step{
		Name: "ProcessPayment",
		Execute: func(ctx context.Context) error {
			return paymentSvc.ProcessPayment(orderID, 50.0, false)
		},
		Compensate: func(ctx context.Context) error {
			return paymentSvc.RefundPayment(orderID)
		},
	})

	// This step fails (request 5 when stock is 1)
	orch.AddStep(saga.Step{
		Name: "ReserveInventory",
		Execute: func(ctx context.Context) error {
			return invSvc.Reserve("item-1", 5)
		},
		Compensate: func(ctx context.Context) error {
			return invSvc.Release("item-1", 5)
		},
	})

	err := orch.Execute(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify rollbacks
	if orderSvc.GetOrder(orderID) != services.OrderCancelled {
		t.Fatalf("expected order cancelled, got %v", orderSvc.GetOrder(orderID))
	}
	if paymentSvc.HasPayment(orderID) {
		t.Fatal("expected payment refunded")
	}
	if invSvc.GetStock("item-1") != 1 {
		t.Fatalf("expected untouched stock 1, got %d", invSvc.GetStock("item-1"))
	}

	logs := orch.Logs()
	expectedOrder := []saga.StepStatus{
		saga.StatusExecuted,    // CreateOrder
		saga.StatusExecuted,    // ProcessPayment
		saga.StatusFailed,      // ReserveInventory
		saga.StatusCompensated, // ProcessPayment
		saga.StatusCompensated, // CreateOrder
	}

	if len(logs) != len(expectedOrder) {
		t.Fatalf("expected %d logs, got %d", len(expectedOrder), len(logs))
	}
	for i, l := range logs {
		if l.Status != expectedOrder[i] {
			t.Errorf("log[%d] expected status %v, got %v", i, expectedOrder[i], l.Status)
		}
	}
}

func TestPayment_Idempotency(t *testing.T) {
	paymentSvc := services.NewPaymentService()
	id := "tx-1"

	if err := paymentSvc.ProcessPayment(id, 100, false); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if err := paymentSvc.ProcessPayment(id, 100, false); err != nil {
		t.Fatalf("second duplicate call should succeed idempotently: %v", err)
	}
}

func TestSemanticLock(t *testing.T) {
	orderSvc := services.NewOrderService()
	id := "ord-lock"

	if err := orderSvc.CreateOrder(id); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	if err := orderSvc.CreateOrder(id); err == nil {
		t.Fatal("second create should fail due to semantic lock")
	}
}

func TestOrchestrator_Concurrency(t *testing.T) {
	ctx := context.Background()
	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"item-1": 100})

	var wg sync.WaitGroup
	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			orderID := fmt.Sprintf("ord-%d", id)
			orch := saga.NewOrchestrator()
			orch.AddStep(saga.Step{
				Name: "CreateOrder",
				Execute: func(ctx context.Context) error {
					return orderSvc.CreateOrder(orderID)
				},
				Compensate: func(ctx context.Context) error {
					return orderSvc.CancelOrder(orderID)
				},
			})
			orch.AddStep(saga.Step{
				Name: "ProcessPayment",
				Execute: func(ctx context.Context) error {
					return paymentSvc.ProcessPayment(orderID, 10, false)
				},
				Compensate: func(ctx context.Context) error {
					return paymentSvc.RefundPayment(orderID)
				},
			})
			orch.AddStep(saga.Step{
				Name: "ReserveInventory",
				Execute: func(ctx context.Context) error {
					return invSvc.Reserve("item-1", 1)
				},
				Compensate: func(ctx context.Context) error {
					return invSvc.Release("item-1", 1)
				},
			})
			_ = orch.Execute(ctx)
		}(i)
	}

	wg.Wait()
	if invSvc.GetStock("item-1") != 90 {
		t.Fatalf("expected stock 90, got %d", invSvc.GetStock("item-1"))
	}
}

func TestChoreography_Flow(t *testing.T) {
	ctx := context.Background()
	bus := saga.NewEventBus()
	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"item-1": 10})

	orderID := "ch-ord-1"

	bus.Subscribe(saga.OrderCreated, func(c context.Context, e saga.Event) {
		id := e.SagaID
		err := paymentSvc.ProcessPayment(id, 20.0, false)
		if err != nil {
			bus.Publish(c, saga.Event{Type: saga.PaymentFailed, SagaID: id})
		} else {
			bus.Publish(c, saga.Event{Type: saga.PaymentCompleted, SagaID: id})
		}
	})

	bus.Subscribe(saga.PaymentCompleted, func(c context.Context, e saga.Event) {
		id := e.SagaID
		err := invSvc.Reserve("item-1", 1)
		if err != nil {
			bus.Publish(c, saga.Event{Type: saga.InventoryFailed, SagaID: id})
		} else {
			bus.Publish(c, saga.Event{Type: saga.InventoryReserved, SagaID: id})
		}
	})

	bus.Subscribe(saga.InventoryReserved, func(c context.Context, e saga.Event) {
		id := e.SagaID
		_ = orderSvc.ApproveOrder(id)
	})

	bus.Subscribe(saga.PaymentFailed, func(c context.Context, e saga.Event) {
		id := e.SagaID
		_ = orderSvc.CancelOrder(id)
	})

	_ = orderSvc.CreateOrder(orderID)
	bus.Publish(ctx, saga.Event{Type: saga.OrderCreated, SagaID: orderID})

	if orderSvc.GetOrder(orderID) != services.OrderApproved {
		t.Fatalf("expected approved order in choreography, got %v", orderSvc.GetOrder(orderID))
	}
}

func TestChoreography_FailureCompensates(t *testing.T) {
	ctx := context.Background()
	bus := saga.NewEventBus()
	orderSvc := services.NewOrderService()
	paymentSvc := services.NewPaymentService()
	invSvc := services.NewInventoryService(map[string]int{"item-1": 0}) // zero stock to force InventoryFailed

	orderID := "ch-ord-fail"

	bus.Subscribe(saga.OrderCreated, func(c context.Context, e saga.Event) {
		id := e.SagaID
		err := paymentSvc.ProcessPayment(id, 50.0, false)
		if err != nil {
			bus.Publish(c, saga.Event{Type: saga.PaymentFailed, SagaID: id})
		} else {
			bus.Publish(c, saga.Event{Type: saga.PaymentCompleted, SagaID: id})
		}
	})

	bus.Subscribe(saga.PaymentCompleted, func(c context.Context, e saga.Event) {
		id := e.SagaID
		err := invSvc.Reserve("item-1", 1)
		if err != nil {
			bus.Publish(c, saga.Event{Type: saga.InventoryFailed, SagaID: id})
		} else {
			bus.Publish(c, saga.Event{Type: saga.InventoryReserved, SagaID: id})
		}
	})

	bus.Subscribe(saga.InventoryFailed, func(c context.Context, e saga.Event) {
		id := e.SagaID
		_ = paymentSvc.RefundPayment(id)
		_ = orderSvc.CancelOrder(id)
	})

	_ = orderSvc.CreateOrder(orderID)
	bus.Publish(ctx, saga.Event{Type: saga.OrderCreated, SagaID: orderID})

	if orderSvc.GetOrder(orderID) != services.OrderCancelled {
		t.Fatalf("expected cancelled order in choreography failure, got %v", orderSvc.GetOrder(orderID))
	}
	if paymentSvc.HasPayment(orderID) {
		t.Fatal("expected payment refunded in choreography failure")
	}
}

func TestOrchestrator_CompensationErrorPropagated(t *testing.T) {
	ctx := context.Background()
	orch := saga.NewOrchestrator()

	orch.AddStep(saga.Step{
		Name: "Step1",
		Execute: func(ctx context.Context) error {
			return nil
		},
		Compensate: func(ctx context.Context) error {
			return fmt.Errorf("rollback failed intentionally")
		},
	})

	orch.AddStep(saga.Step{
		Name: "Step2",
		Execute: func(ctx context.Context) error {
			return fmt.Errorf("step 2 trigger error")
		},
	})

	err := orch.Execute(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	logs := orch.Logs()
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}
	if logs[0].Status != saga.StatusExecuted || logs[1].Status != saga.StatusFailed || logs[2].Status != saga.StatusCompensateFailed {
		t.Fatalf("unexpected logs: %+v", logs)
	}
}

func TestOrchestrator_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	orch := saga.NewOrchestrator()

	step1Compensated := false
	orch.AddStep(saga.Step{
		Name: "Step1",
		Execute: func(ctx context.Context) error {
			cancel() // cancel context before step 2
			return nil
		},
		Compensate: func(ctx context.Context) error {
			step1Compensated = true
			return nil
		},
	})

	orch.AddStep(saga.Step{
		Name: "Step2",
		Execute: func(ctx context.Context) error {
			t.Fatal("step 2 should not execute when context cancelled")
			return nil
		},
	})

	err := orch.Execute(ctx)
	if err == nil {
		t.Fatal("expected error due to cancellation, got nil")
	}
	if !step1Compensated {
		t.Fatal("expected step 1 to be compensated upon cancellation")
	}
}
