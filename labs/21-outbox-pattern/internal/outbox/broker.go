package outbox

import "sync"

type Broker interface {
	Publish(msg OutboxMessage) error
	GetPublished() []OutboxMessage
}

type MockBroker struct {
	mu        sync.RWMutex
	published []OutboxMessage
	failNext  bool
}

func NewMockBroker() *MockBroker {
	return &MockBroker{
		published: make([]OutboxMessage, 0),
	}
}

func (b *MockBroker) SetFailNext(fail bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failNext = fail
}

func (b *MockBroker) Publish(msg OutboxMessage) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.failNext {
		b.failNext = false
		return &BrokerError{Message: "broker unavailable"}
	}
	b.published = append(b.published, msg)
	return nil
}

func (b *MockBroker) GetPublished() []OutboxMessage {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]OutboxMessage, len(b.published))
	copy(out, b.published)
	return out
}

type BrokerError struct {
	Message string
}

func (e *BrokerError) Error() string {
	return e.Message
}
