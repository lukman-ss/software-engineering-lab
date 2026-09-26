package outbox

import "sync"

type Consumer struct {
	mu           sync.Mutex
	processedIDs map[string]bool
	received     []OutboxMessage
}

func NewConsumer() *Consumer {
	return &Consumer{
		processedIDs: make(map[string]bool),
		received:     make([]OutboxMessage, 0),
	}
}

// Handle processes the message idempotently. Returns true if processed, false if duplicate.
func (c *Consumer) Handle(msg OutboxMessage) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.processedIDs[msg.ID] {
		// Duplicate detected, ignore
		return false
	}

	c.processedIDs[msg.ID] = true
	c.received = append(c.received, msg)
	return true
}

func (c *Consumer) GetReceivedCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.received)
}

func (c *Consumer) IsProcessed(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.processedIDs[id]
}
