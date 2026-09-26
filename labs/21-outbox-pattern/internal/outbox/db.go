package outbox

import (
	"errors"
	"sync"
)

var (
	ErrTxClosed = errors.New("transaction already closed")
)

type DB struct {
	mu     sync.RWMutex
	orders map[string]Order
	outbox map[string]OutboxMessage
}

func NewDB() *DB {
	return &DB{
		orders: make(map[string]Order),
		outbox: make(map[string]OutboxMessage),
	}
}

func (db *DB) BeginTx() *Tx {
	return &Tx{
		db:           db,
		stagedOrders: make(map[string]Order),
		stagedOutbox: make(map[string]OutboxMessage),
	}
}

func (db *DB) GetOrder(id string) (Order, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	order, ok := db.orders[id]
	return order, ok
}

func (db *DB) GetOutbox(id string) (OutboxMessage, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	msg, ok := db.outbox[id]
	return msg, ok
}

func (db *DB) GetPendingOutbox() []OutboxMessage {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var pending []OutboxMessage
	for _, msg := range db.outbox {
		if msg.Status == MessageStatusPending {
			pending = append(pending, msg)
		}
	}
	return pending
}

func (db *DB) MarkOutboxProcessed(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	msg, ok := db.outbox[id]
	if !ok {
		return errors.New("message not found")
	}
	msg.Status = MessageStatusProcessed
	db.outbox[id] = msg
	return nil
}

type Tx struct {
	mu           sync.Mutex
	db           *DB
	stagedOrders map[string]Order
	stagedOutbox map[string]OutboxMessage
	closed       bool
}

func (tx *Tx) SaveOrder(order Order) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.stagedOrders[order.ID] = order
	return nil
}

func (tx *Tx) SaveOutbox(msg OutboxMessage) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.stagedOutbox[msg.ID] = msg
	return nil
}

func (tx *Tx) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true

	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()
	for k, v := range tx.stagedOrders {
		tx.db.orders[k] = v
	}
	for k, v := range tx.stagedOutbox {
		tx.db.outbox[k] = v
	}
	return nil
}

func (tx *Tx) Rollback() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true
	// discard staged mutations
	return nil
}
