package interval

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// Example-based test: passes when hand-picked examples are sorted.
func TestExampleBasedIntervalMerge(t *testing.T) {
	input := []Interval{
		{Start: 1, End: 3},
		{Start: 2, End: 6},
		{Start: 8, End: 10},
		{Start: 15, End: 18},
	}
	expected := []Interval{
		{Start: 1, End: 6},
		{Start: 8, End: 10},
		{Start: 15, End: 18},
	}

	res := NaiveMerge(input)
	if !intervalsEqual(res, expected) {
		t.Fatalf("expected %+v, got %+v", expected, res)
	}
}

type IntervalSlice []Interval

func (IntervalSlice) Generate(r *rand.Rand, size int) reflect.Value {
	n := r.Intn(8) + 2
	slice := make([]Interval, n)
	for i := 0; i < n; i++ {
		s := r.Intn(100) - 20
		e := s + r.Intn(30)
		slice[i] = Interval{Start: s, End: e}
	}
	return reflect.ValueOf(IntervalSlice(slice))
}

// Property-based test: NaiveMerge fails oracle property compared to RobustMerge on arbitrary input.
func TestPropertyNaiveMergeFails(t *testing.T) {
	prop := func(ivs IntervalSlice) bool {
		naive := NaiveMerge([]Interval(ivs))
		robust := RobustMerge([]Interval(ivs))
		return intervalsEqual(naive, robust)
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	err := quick.Check(prop, cfg)
	if err == nil {
		t.Fatal("expected NaiveMerge oracle property to fail, but it passed")
	}
}

// Property-based test 1: Idempotence Invariant
// Merge(Merge(intervals)) == Merge(intervals)
func TestPropertyRobustMergeIdempotence(t *testing.T) {
	prop := func(ivs IntervalSlice) bool {
		once := RobustMerge([]Interval(ivs))
		twice := RobustMerge(once)
		return intervalsEqual(once, twice)
	}

	cfg := &quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(prop, cfg); err != nil {
		t.Fatalf("property violation in RobustMerge idempotence: %v", err)
	}
}

// Property-based test 2: Non-overlapping invariant
// For all consecutive intervals i and i+1 in output, out[i].End < out[i+1].Start
func TestPropertyRobustMergeNonOverlapping(t *testing.T) {
	prop := func(ivs IntervalSlice) bool {
		merged := RobustMerge([]Interval(ivs))
		for i := 0; i < len(merged)-1; i++ {
			if merged[i].End >= merged[i+1].Start {
				return false
			}
		}
		return true
	}

	cfg := &quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(prop, cfg); err != nil {
		t.Fatalf("property violation in RobustMerge non-overlapping output: %v", err)
	}
}
