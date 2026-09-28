package shrinker

import (
	"fmt"
	"math/rand"
	"sort"
)

// ShrinkStep records a single step during counterexample shrinking.
type ShrinkStep struct {
	Step     int
	Input    []int
	Failed   bool
	Strategy string
}

// ShrinkResult contains the initial failing input, shrinking steps, and the minimal counterexample.
type ShrinkResult struct {
	InitialFailing []int
	Minimal        []int
	Steps          []ShrinkStep
}

// FindAndShrink searches for an input array that violates `predicate`,
// then performs counterexample shrinking (list binary sectioning & element reduction).
func FindAndShrink(predicate func([]int) bool, rng *rand.Rand, maxAttempts int) (*ShrinkResult, error) {
	var initial []int
	found := false

	// Step 1: Find initial failing counterexample
	for i := 0; i < maxAttempts; i++ {
		n := rng.Intn(20) + 5
		candidate := make([]int, n)
		for j := 0; j < n; j++ {
			candidate[j] = rng.Intn(200) - 50
		}
		if !predicate(candidate) {
			initial = candidate
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no failing counterexample found within %d attempts", maxAttempts)
	}

	res := &ShrinkResult{
		InitialFailing: append([]int(nil), initial...),
		Minimal:        append([]int(nil), initial...),
	}

	stepCount := 0
	current := append([]int(nil), initial...)

	// Step 2: Binary / List section shrinking (try removing halves or single elements)
	changed := true
	for changed {
		changed = false

		// Strategy 1: Try halving / chunk removal
		if len(current) > 1 {
			mid := len(current) / 2
			// Try left half
			left := current[:mid]
			stepCount++
			failedLeft := !predicate(left)
			res.Steps = append(res.Steps, ShrinkStep{
				Step:     stepCount,
				Input:    append([]int(nil), left...),
				Failed:   failedLeft,
				Strategy: "remove-right-half",
			})
			if failedLeft {
				current = append([]int(nil), left...)
				changed = true
				continue
			}

			// Try right half
			right := current[mid:]
			stepCount++
			failedRight := !predicate(right)
			res.Steps = append(res.Steps, ShrinkStep{
				Step:     stepCount,
				Input:    append([]int(nil), right...),
				Failed:   failedRight,
				Strategy: "remove-left-half",
			})
			if failedRight {
				current = append([]int(nil), right...)
				changed = true
				continue
			}
		}

		// Strategy 2: Remove individual elements
		for i := 0; i < len(current); i++ {
			reduced := make([]int, 0, len(current)-1)
			reduced = append(reduced, current[:i]...)
			reduced = append(reduced, current[i+1:]...)

			stepCount++
			failed := !predicate(reduced)
			res.Steps = append(res.Steps, ShrinkStep{
				Step:     stepCount,
				Input:    append([]int(nil), reduced...),
				Failed:   failed,
				Strategy: fmt.Sprintf("remove-elem-idx-%d", i),
			})

			if failed {
				current = reduced
				changed = true
				break
			}
		}

		// Strategy 3: Reduce element magnitude toward zero
		for i := 0; i < len(current); i++ {
			if current[i] == 0 {
				continue
			}
			candidates := []int{0, current[i] / 2}
			for _, cand := range candidates {
				if cand == current[i] {
					continue
				}
				testSlice := append([]int(nil), current...)
				testSlice[i] = cand

				stepCount++
				failed := !predicate(testSlice)
				res.Steps = append(res.Steps, ShrinkStep{
					Step:     stepCount,
					Input:    append([]int(nil), testSlice...),
					Failed:   failed,
					Strategy: fmt.Sprintf("reduce-val-idx-%d-to-%d", i, cand),
				})

				if failed {
					current = testSlice
					changed = true
					break
				}
			}
			if changed {
				break
			}
		}
	}

	res.Minimal = current
	return res, nil
}

// BuggySortContainsNegative returns true if array sorted correctly AND contains no negative elements.
// Intentionally fails when slice has negative numbers.
func BuggySortPredicate(arr []int) bool {
	cp := append([]int(nil), arr...)
	sort.Ints(cp)
	for _, v := range cp {
		if v < 0 {
			return false // Invariant violated: contains negative number!
		}
	}
	return true
}
