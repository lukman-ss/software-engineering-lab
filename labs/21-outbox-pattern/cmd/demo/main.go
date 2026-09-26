package main

import (
	"fmt"
	"time"

	"github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox"
)

func main() {
	fmt.Println("=== Lab 21: Transactional Outbox Pattern Demo ===")

	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 50*time.Millisecond)
	service := outbox.NewOrderService(db)
	consumer := outbox.NewConsumer()

	// 1. Demonstrate Dual-Write Flaw
	fmt.Println("\n[Scenario 1: The Dual-Write Problem]")
	broker.SetFailNext(true) // broker failure simulation
	err := service.CreateOrderDualWriteNaive(broker, "order-dual-fail", "cust-1", 150.0)
	if err != nil {
		fmt.Printf("Direct write failed: %v\n", err)
	}
	_, dbFound := db.GetOrder("order-dual-fail")
	fmt.Printf("State Inconsistency: Order in DB = %v, Broker Message Count = %d\n", dbFound, len(broker.GetPublished()))

	// 2. Start Message Relay for Outbox
	fmt.Println("\n[Scenario 2: Transactional Outbox Solution]")
	relay.Start()
	defer relay.Stop()

	fmt.Println("Creating order with transactional outbox...")
	err = service.CreateOrderWithOutbox("order-outbox-success", "cust-2", 300.0)
	if err != nil {
		fmt.Printf("Failed to create order: %v\n", err)
		return
	}
	fmt.Println("Order and Outbox record atomically saved to DB.")

	// Let relay poll and dispatch
	time.Sleep(120 * time.Millisecond)

	published := broker.GetPublished()
	fmt.Printf("Broker received messages: %d\n", len(published))
	for _, m := range published {
		fmt.Printf(" - Event ID: %s, Type: %s, Payload: %s\n", m.ID, m.EventType, m.Payload)
		// Consumer processes message
		accepted := consumer.Handle(m)
		fmt.Printf(" - Consumer processing initial message: accepted=%v\n", accepted)
	}

	// 3. Demonstrate At-Least-Once & Idempotency
	fmt.Println("\n[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]")
	duplicateMsg := published[0]
	fmt.Println("Simulating duplicate delivery to consumer...")
	acceptedAgain := consumer.Handle(duplicateMsg)
	fmt.Printf("Consumer processing duplicate delivery: accepted=%v (Duplicate safely skipped!)\n", acceptedAgain)
	fmt.Printf("Total events processed by consumer: %d\n", consumer.GetReceivedCount())

	fmt.Println("\n=== Demo Complete ===")
}
