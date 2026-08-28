package inventory

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrConcurrentUpdate  = errors.New("inventory updated by another request")
)

type InventoryItem struct {
	ID        string
	ProductID string
	Quantity  int
	Version   int // for optimistic locking
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	GetByProductID(ctx context.Context, productID string) (*InventoryItem, error)
	// Reserve decreases quantity. It must handle concurrency (e.g. optimistic/pessimistic locking)
	Reserve(ctx context.Context, productID string, quantity int) error
	// Restock increases quantity.
	Restock(ctx context.Context, productID string, quantity int) error
}

type Service interface {
	CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error)
	ReserveItems(ctx context.Context, items map[string]int) error
	ReleaseItems(ctx context.Context, items map[string]int) error
}
