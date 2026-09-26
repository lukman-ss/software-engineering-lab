package blog

import (
	"testing"
)

func TestGetAuthorsWithPostsNPlusOne(t *testing.T) {
	store := NewStore()
	repo := NewRepository(store)
	
	store.ResetQueryCount()
	
	result := repo.GetAuthorsWithPostsNPlusOne()
	
	if len(result) != 3 {
		t.Fatalf("expected 3 authors, got %d", len(result))
	}
	
	// Query count should be N+1 (3 authors + 1 initial query) = 4
	queryCount := store.GetQueryCount()
	if queryCount != 4 {
		t.Errorf("expected 4 queries (N+1), got %d", queryCount)
	}
}

func TestGetAuthorsWithPostsEager(t *testing.T) {
	store := NewStore()
	repo := NewRepository(store)
	
	store.ResetQueryCount()
	
	result := repo.GetAuthorsWithPostsEager()
	
	if len(result) != 3 {
		t.Fatalf("expected 3 authors, got %d", len(result))
	}
	
	// Query count should be 2 (1 for authors, 1 for all posts)
	queryCount := store.GetQueryCount()
	if queryCount != 2 {
		t.Errorf("expected 2 queries (eager loading), got %d", queryCount)
	}
}
