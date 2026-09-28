package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"labs/40-property-based-testing/internal/currency"
	"labs/40-property-based-testing/internal/interval"
	"labs/40-property-based-testing/internal/shrinker"
)

const (
	passed = "PASS"
	failed = "FAIL"
)

func sectionHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 72))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 72))
}

func subHeader(title string) {
	fmt.Printf("\n-- %s --\n", title)
}

// Run example-based currency tests and show which pass despite float bugs
func demoCurrency() {
	sectionHeader("DEMO 1: Currency Formatting – Roundtrip Invariant")
	naive := currency.NaiveCurrency{}

	subHeader("Example-Based Tests (Naive Float Currency)")
	happyCases := []float64{1.00, 5.25, 99.99, 10.50, 0.01}
	examplePassed := 0
	for _, dollars := range happyCases {
		formatted := naive.Format(dollars)
		parsed, _ := naive.Parse(formatted)
		if parsed == dollars {
			fmt.Printf("  EXAMPLE PASS: $%.2f -> %s -> $%.2f\n", dollars, formatted, parsed)
			examplePassed++
		} else {
			fmt.Printf("  EXAMPLE FAIL: $%.6f -> %s -> $%.6f\n", dollars, formatted, parsed)
		}
	}
	fmt.Printf("  Result: %d/%d examples pass (false confidence!)\n", examplePassed, len(happyCases))

	subHeader("Property-Based Test: Detect Float Precision Loss")
	rng := rand.New(rand.NewSource(12345))
	failedInputs := []float64{}
	for i := 0; i < 1000; i++ {
		// Generate dollars with arbitrary sub-cent precision
		raw := (rng.Float64() - 0.5) * 10000
		// Keep full float precision, NOT rounded to 2 decimal places
		formatted := naive.Format(raw)
		parsed, err := naive.Parse(formatted)
		if err != nil || parsed != raw {
			failedInputs = append(failedInputs, raw)
		}
	}
	if len(failedInputs) > 0 {
		fmt.Printf("  Property FAIL: %d/1000 inputs failed roundtrip\n", len(failedInputs))
		fmt.Printf("  Counterexample: %.10f -> %s -> %.10f\n",
			failedInputs[0], naive.Format(failedInputs[0]), func() float64 {
				p, _ := naive.Parse(naive.Format(failedInputs[0]))
				return p
			}())
	} else {
		fmt.Println("  Property PASS: no failures found")
	}

	subHeader("Property-Based Test: Robust Integer Cents Roundtrip (1000 iterations)")
	propertyFailed := 0
	for i := 0; i < 1000; i++ {
		cents := rng.Int63n(2000000000) - 1000000000
		a := currency.NewRobustAmount(cents)
		formatted := a.Format()
		parsed, err := currency.ParseRobust(formatted)
		if err != nil || parsed.Cents != cents {
			propertyFailed++
			fmt.Printf("  FAIL: cents=%d, formatted=%s, parsed_cents=%d\n", cents, formatted, parsed.Cents)
		}
	}
	if propertyFailed == 0 {
		fmt.Printf("  All 1000 iterations PASS: ParseRobust(amount.Format()) == amount\n")
	} else {
		fmt.Printf("  %d/1000 iterations FAIL\n", propertyFailed)
	}
}

func demoInterval() {
	sectionHeader("DEMO 2: Interval Merging – Idempotence Invariant")

	subHeader("Example-Based Test (Sorted Input – Passes Naive)")
	sortedInput := []interval.Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	naive := interval.NaiveMerge(sortedInput)
	fmt.Printf("  Input:  %+v\n", sortedInput)
	fmt.Printf("  Merged: %+v\n", naive)
	idempotenceOnSorted := fmt.Sprintf("%+v", interval.NaiveMerge(naive)) == fmt.Sprintf("%+v", naive)
	fmt.Printf("  Idempotent: %v (passes on sorted example – false confidence!)\n", idempotenceOnSorted)

	subHeader("Property-Based Test: NaiveMerge fails oracle check compared to RobustMerge")
	rng := rand.New(rand.NewSource(99))
	naiveFailures := 0
	var firstNaiveFail []interval.Interval
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		ivs := make([]interval.Interval, n)
		for j := 0; j < n; j++ {
			s := rng.Intn(40) - 10
			ivs[j] = interval.Interval{Start: s, End: s + rng.Intn(20)}
		}
		naiveRes := interval.NaiveMerge(ivs)
		robustRes := interval.RobustMerge(ivs)
		if fmt.Sprintf("%+v", naiveRes) != fmt.Sprintf("%+v", robustRes) {
			naiveFailures++
			if firstNaiveFail == nil {
				firstNaiveFail = ivs
			}
		}
	}
	fmt.Printf("  NaiveMerge vs RobustMerge oracle discrepancies: %d/100\n", naiveFailures)
	if firstNaiveFail != nil {
		naiveRes := interval.NaiveMerge(firstNaiveFail)
		robustRes := interval.RobustMerge(firstNaiveFail)
		fmt.Printf("  Counterexample input: %+v\n", firstNaiveFail)
		fmt.Printf("  NaiveMerge(input):    %+v\n", naiveRes)
		fmt.Printf("  RobustMerge(input):   %+v\n", robustRes)
	}

	subHeader("Property-Based Test: RobustMerge passes idempotence (1000 iterations)")
	robustFail := 0
	for i := 0; i < 1000; i++ {
		n := rng.Intn(10)
		ivs := make([]interval.Interval, n)
		for j := 0; j < n; j++ {
			s := rng.Intn(200) - 50
			ivs[j] = interval.Interval{Start: s, End: s + rng.Intn(50)}
		}
		once := interval.RobustMerge(ivs)
		twice := interval.RobustMerge(once)
		if fmt.Sprintf("%+v", once) != fmt.Sprintf("%+v", twice) {
			robustFail++
		}
	}
	if robustFail == 0 {
		fmt.Printf("  All 1000 iterations PASS: Merge(Merge(x)) == Merge(x)\n")
	} else {
		fmt.Printf("  FAIL: %d/1000 iterations violated idempotence\n", robustFail)
	}

	subHeader("Property-Based Test: Non-overlapping output invariant (1000 iterations)")
	overlapFail := 0
	for i := 0; i < 1000; i++ {
		n := rng.Intn(12)
		ivs := make([]interval.Interval, n)
		for j := 0; j < n; j++ {
			s := rng.Intn(200) - 50
			ivs[j] = interval.Interval{Start: s, End: s + rng.Intn(50)}
		}
		merged := interval.RobustMerge(ivs)
		for k := 0; k < len(merged)-1; k++ {
			if merged[k].End >= merged[k+1].Start {
				overlapFail++
				break
			}
		}
	}
	if overlapFail == 0 {
		fmt.Printf("  All 1000 iterations PASS: output intervals are non-overlapping\n")
	} else {
		fmt.Printf("  FAIL: %d/1000 iterations had overlapping output\n", overlapFail)
	}
}

func demoShrinking() {
	sectionHeader("DEMO 3: Counterexample Shrinking")

	subHeader("Invariant: all elements must be non-negative (buggy predicate for demo)")
	rng := rand.New(rand.NewSource(42))

	res, err := shrinker.FindAndShrink(shrinker.BuggySortPredicate, rng, 200)
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
		return
	}

	fmt.Printf("  Initial failing input (%d elements): %v\n", len(res.InitialFailing), res.InitialFailing)
	fmt.Printf("  Shrinking steps taken: %d\n", len(res.Steps))

	// Print only steps that actually reduced the failure (kept failing=true)
	fmt.Println("\n  Shrinking trace (failing steps only):")
	printed := 0
	for _, step := range res.Steps {
		if step.Failed {
			status := "FAIL (used)"
			fmt.Printf("    Step %3d %-30s [%d elements] %v -> %s\n",
				step.Step, step.Strategy, len(step.Input), step.Input, status)
			printed++
			if printed >= 10 {
				fmt.Printf("    ... (%d more shrinking steps) ...\n", len(res.Steps)-printed)
				break
			}
		}
	}

	fmt.Printf("\n  Minimal counterexample (%d elements): %v\n", len(res.Minimal), res.Minimal)
	negFound := false
	for _, v := range res.Minimal {
		if v < 0 {
			negFound = true
		}
	}
	if negFound && len(res.Minimal) == 1 {
		fmt.Printf("  Shrinking SUCCESS: reduced %d elements to 1 element showing exact cause\n",
			len(res.InitialFailing))
	}
}

func demoSummary() {
	sectionHeader("SUMMARY: Example-Based vs Property-Based Testing")
	fmt.Println()
	fmt.Printf("  %-40s %-20s %-20s\n", "Check", "Example-Based", "Property-Based")
	fmt.Println("  " + strings.Repeat("-", 80))
	rows := [][]string{
		{"Float precision bugs caught", "NO", "YES"},
		{"Unsorted interval edge cases caught", "NO", "YES"},
		{"Negative value handling verified", "NO", "YES"},
		{"Empty / minimal cases tested", "RARELY", "ALWAYS"},
		{"Automatic bug isolation (shrinking)", "NO", "YES"},
		{"Confidence on happy path examples", "HIGH (false)", "MEDIUM (calibrated)"},
	}
	for _, row := range rows {
		fmt.Printf("  %-40s %-20s %-20s\n", row[0], row[1], row[2])
	}
	fmt.Println()
}

var _ = math.Pi
var _ = sort.Ints

func main() {
	demoCurrency()
	demoInterval()
	demoShrinking()
	demoSummary()
}
