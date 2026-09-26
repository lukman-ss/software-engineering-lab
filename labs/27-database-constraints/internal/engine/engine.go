package engine

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/dberr"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/model"
)

// Engine simulates an ACID relational table store with declarative constraint checking.
type Engine struct {
	mu           sync.RWMutex
	userSeq      atomic.Int64
	orderSeq     atomic.Int64
	users        map[int64]model.User
	orders       map[int64]model.Order
	emailIndex   map[string]int64 // UNIQUE (email)
	activeEmails map[string]int64 // PARTIAL UNIQUE (email) WHERE deleted_at IS NULL
}

func NewEngine() *Engine {
	return &Engine{
		users:        make(map[int64]model.User),
		orders:       make(map[int64]model.Order),
		emailIndex:   make(map[string]int64),
		activeEmails: make(map[string]int64),
	}
}

// InsertUser evaluates NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints atomically.
func (e *Engine) InsertUser(u model.User, usePartialUniqueIndex bool) (model.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. NOT NULL checks
	if u.Email == "" {
		return model.User{}, dberr.NewNotNullViolation("users", "email", "users_email_not_null")
	}
	if u.Username == "" {
		return model.User{}, dberr.NewNotNullViolation("users", "username", "users_username_not_null")
	}

	// 2. CHECK constraints
	// CHECK (age >= 18)
	if u.Age < 18 {
		return model.User{}, dberr.NewCheckViolation("users", "users_age_check", fmt.Sprintf("value %d for column 'age' violates check constraint (age >= 18)", u.Age))
	}
	// CHECK (status IN ('active', 'suspended', 'pending'))
	switch u.Status {
	case "active", "suspended", "pending":
	default:
		return model.User{}, dberr.NewCheckViolation("users", "users_status_check", fmt.Sprintf("invalid status %q violates check constraint", u.Status))
	}

	// 3. UNIQUE / PARTIAL UNIQUE constraint
	if usePartialUniqueIndex {
		// Partial Unique Index: CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
		if u.DeletedAt == nil {
			if existingID, exists := e.activeEmails[u.Email]; exists {
				return model.User{}, dberr.NewUniqueViolation("users", "users_active_email_idx", fmt.Sprintf("duplicate key value violates unique constraint 'users_active_email_idx' (email=%s, existing_id=%d)", u.Email, existingID))
			}
		}
	} else {
		// Standard full UNIQUE: CONSTRAINT users_email_key UNIQUE (email)
		if existingID, exists := e.emailIndex[u.Email]; exists {
			return model.User{}, dberr.NewUniqueViolation("users", "users_email_key", fmt.Sprintf("duplicate key value violates unique constraint 'users_email_key' (email=%s, existing_id=%d)", u.Email, existingID))
		}
	}

	// Primary Key generation
	if u.ID == 0 {
		u.ID = e.userSeq.Add(1)
	}

	// Commit row and indexes
	e.users[u.ID] = u
	if !usePartialUniqueIndex {
		e.emailIndex[u.Email] = u.ID
	} else if u.DeletedAt == nil {
		e.activeEmails[u.Email] = u.ID
	}

	return u, nil
}

// SoftDeleteUser updates DeletedAt and updates partial indexes.
func (e *Engine) SoftDeleteUser(id int64, deletedAt model.User) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	u, exists := e.users[id]
	if !exists {
		return fmt.Errorf("user %d not found", id)
	}

	if deletedAt.DeletedAt == nil {
		return fmt.Errorf("deleted_at cannot be nil for soft delete")
	}

	// Remove from partial active index
	delete(e.activeEmails, u.Email)
	u.DeletedAt = deletedAt.DeletedAt
	e.users[id] = u
	return nil
}

// InsertOrder evaluates NOT NULL, CHECK, and FOREIGN KEY constraints.
func (e *Engine) InsertOrder(o model.Order) (model.Order, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. NOT NULL constraint on user_id
	if o.UserID == 0 {
		return model.Order{}, dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null")
	}

	// 2. CHECK constraint (total_cents > 0)
	if o.TotalCents <= 0 {
		return model.Order{}, dberr.NewCheckViolation("orders", "orders_total_cents_check", fmt.Sprintf("total_cents %d must be greater than zero", o.TotalCents))
	}

	// 3. FOREIGN KEY constraint: CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id)
	if _, exists := e.users[o.UserID]; !exists {
		return model.Order{}, dberr.NewForeignKeyViolation("orders", "fk_orders_user", fmt.Sprintf("key (user_id)=(%d) is not present in table \"users\"", o.UserID))
	}

	if o.ID == 0 {
		o.ID = e.orderSeq.Add(1)
	}

	e.orders[o.ID] = o
	return o, nil
}

func (e *Engine) GetUser(id int64) (model.User, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	u, ok := e.users[id]
	return u, ok
}

func (e *Engine) GetUsersCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.users)
}
