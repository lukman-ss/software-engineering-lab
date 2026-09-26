package tests

import (
	"errors"
	"testing"

	"zero-downtime-deployment/internal/db"
)

func TestDBNotFound(t *testing.T) {
	store := db.NewUserStore()
	_, err := store.GetUser("nonexistent")
	if !errors.Is(err, db.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDBSingleNameLegacy(t *testing.T) {
	store := db.NewUserStore()
	store.InsertLegacy("1", "Madonna")
	user, err := store.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.FirstName != "Madonna" || user.LastName != "" {
		t.Fatalf("expected FirstName=Madonna LastName='', got %s / %s", user.FirstName, user.LastName)
	}
}

func TestDBSaveExpandEmptyFields(t *testing.T) {
	store := db.NewUserStore()

	store.SaveExpand("1", "", "Smith")
	u1, err := store.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u1.Name != "Smith" || u1.FirstName != "" || u1.LastName != "Smith" {
		t.Fatalf("empty firstName: unexpected %+v", u1)
	}

	store.SaveExpand("2", "Jane", "")
	u2, err := store.GetUser("2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u2.Name != "Jane" || u2.FirstName != "Jane" || u2.LastName != "" {
		t.Fatalf("empty lastName: unexpected %+v", u2)
	}
}

func TestExpandContractDatabase(t *testing.T) {
	store := db.NewUserStore()

	// 1. Legacy write
	store.InsertLegacy("1", "John Doe")

	// Read using new API (backward compatible)
	user, err := store.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error reading legacy user: %v", err)
	}
	if user.FirstName != "John" || user.LastName != "Doe" {
		t.Fatalf("expected John Doe, got %s / %s", user.FirstName, user.LastName)
	}

	// 2. Expand phase write (dual write or modern write)
	store.SaveExpand("2", "Jane", "Smith")

	user2, err := store.GetUser("2")
	if err != nil {
		t.Fatalf("unexpected error reading modern user: %v", err)
	}
	if user2.Name != "Jane Smith" || user2.FirstName != "Jane" || user2.LastName != "Smith" {
		t.Fatalf("expected Jane Smith fields populated, got %+v", user2)
	}
}
