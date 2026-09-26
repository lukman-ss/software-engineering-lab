package store

import (
	"context"
	"fmt"

	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/dberr"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/engine"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/model"
)

// UnsafeStore lacks database constraints and relies entirely on application memory checks,
// making it vulnerable to read-then-write race conditions under concurrency.
type UnsafeStore struct {
	eng *engine.Engine
}

func NewUnsafeStore(eng *engine.Engine) *UnsafeStore {
	return &UnsafeStore{eng: eng}
}

// RegisterUser performs application-level check first, then inserts. Vulnerable to race conditions.
func (s *UnsafeStore) RegisterUser(ctx context.Context, u model.User) (model.User, error) {
	// Vulnerable app-level check: no database constraint protecting email uniqueness!
	// If two goroutines reach this concurrently, both see count == 0 and both proceed to insert.
	var duplicateFound bool
	count := s.eng.GetUsersCount()
	for i := 1; i <= count; i++ {
		existing, ok := s.eng.GetUser(int64(i))
		if ok && existing.Email == u.Email {
			duplicateFound = true
			break
		}
	}

	if duplicateFound {
		return model.User{}, fmt.Errorf("app check failed: email %s already taken", u.Email)
	}

	// Engine insert bypassing unique constraints (simulating unconstrained table)
	res, err := s.eng.InsertUser(u, true) // partial unique deactivated for unsafe demo
	return res, err
}

// SafeStore relies on declarative database constraints at the storage engine layer.
type SafeStore struct {
	eng *engine.Engine
}

func NewSafeStore(eng *engine.Engine) *SafeStore {
	return &SafeStore{eng: eng}
}

// RegisterUser delegates uniqueness and validation directly to storage engine constraints.
func (s *SafeStore) RegisterUser(ctx context.Context, u model.User) (model.User, error) {
	// Directly insert; database constraint handles race condition and returns 23505
	res, err := s.eng.InsertUser(u, false)
	if err != nil {
		return model.User{}, dberr.MapToDomainError(err)
	}
	return res, nil
}

// RegisterUserPartial delegates to database partial unique index.
func (s *SafeStore) RegisterUserPartial(ctx context.Context, u model.User) (model.User, error) {
	res, err := s.eng.InsertUser(u, true)
	if err != nil {
		return model.User{}, dberr.MapToDomainError(err)
	}
	return res, nil
}

// CreateOrder enforces foreign key integrity via database engine.
func (s *SafeStore) CreateOrder(ctx context.Context, o model.Order) (model.Order, error) {
	res, err := s.eng.InsertOrder(o)
	if err != nil {
		return model.Order{}, dberr.MapToDomainError(err)
	}
	return res, nil
}
