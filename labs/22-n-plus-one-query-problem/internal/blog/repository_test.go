package blog

import (
	"reflect"
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

	// Deep equivalence check with N+1 result
	nPlusOneResult := repo.GetAuthorsWithPostsNPlusOne()
	if !reflect.DeepEqual(result, nPlusOneResult) {
		t.Errorf("eager result mismatch with N+1 result:\nEager: %+v\nN+1: %+v", result, nPlusOneResult)
	}
}

func TestEmptyStore(t *testing.T) {
	store := &Store{authors: nil, posts: nil}
	repo := NewRepository(store)

	n1 := repo.GetAuthorsWithPostsNPlusOne()
	eager := repo.GetAuthorsWithPostsEager()

	if len(n1) != 0 || len(eager) != 0 {
		t.Fatalf("expected empty results, got n1=%d eager=%d", len(n1), len(eager))
	}
	if !reflect.DeepEqual(n1, eager) {
		t.Errorf("empty result mismatch: n1=%#v eager=%#v", n1, eager)
	}
}
