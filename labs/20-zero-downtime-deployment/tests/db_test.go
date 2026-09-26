package tests

import (
	"testing"

	"zero-downtime-deployment/internal/db"
)

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
