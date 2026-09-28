package shrinker

import (
	"math/rand"
	"testing"
)

func TestFindAndShrink(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	// Invariant: all elements must be non-negative
	res, err := FindAndShrink(BuggySortPredicate, rng, 100)
	if err != nil {
		t.Fatalf("failed to find and shrink: %v", err)
	}

	if len(res.InitialFailing) <= len(res.Minimal) {
		t.Fatalf("shrinking did not reduce size: initial=%d, minimal=%d", len(res.InitialFailing), len(res.Minimal))
	}

	// Minimal counterexample for "no negative numbers" is a slice of length 1 containing a negative number
	if len(res.Minimal) != 1 {
		t.Fatalf("expected minimal slice of length 1, got %v", res.Minimal)
	}
	if res.Minimal[0] >= 0 {
		t.Fatalf("expected minimal element to be negative, got %d", res.Minimal[0])
	}
}
