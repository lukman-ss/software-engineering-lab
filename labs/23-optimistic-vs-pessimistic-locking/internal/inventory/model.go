package inventory

import "errors"

var (
	ErrNotFound          = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOptimisticLock    = errors.New("optimistic lock conflict: record modified by another transaction")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
)

type Product struct {
	ID      int
	Name    string
	Stock   int
	Version int
}
