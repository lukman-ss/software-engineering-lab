package saga

import (
	"context"
	"sync"
)

type EventType string

const (
	OrderCreated      EventType = "OrderCreated"
	OrderCancelled    EventType = "OrderCancelled"
	PaymentCompleted  EventType = "PaymentCompleted"
	PaymentFailed     EventType = "PaymentFailed"
	InventoryReserved EventType = "InventoryReserved"
	InventoryFailed   EventType = "InventoryFailed"
)

type Event struct {
	Type    EventType
	SagaID  string
	Payload any
}

type EventHandler func(ctx context.Context, e Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (b *EventBus) Subscribe(t EventType, h EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

func (b *EventBus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	handlers := append([]EventHandler(nil), b.handlers[e.Type]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, e)
	}
}
