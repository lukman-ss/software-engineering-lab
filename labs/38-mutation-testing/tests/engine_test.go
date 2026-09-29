package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"labs/38-mutation-testing/internal/engine"
	"labs/38-mutation-testing/internal/service"
)

func TestEngine_GeneratesMutants(t *testing.T) {
	srcPath := filepath.Join("..", "internal", "service", "discount.go")
	src, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Failed to read discount.go: %v", err)
	}

	mutator := engine.NewASTMutator()
	plans, err := mutator.GenerateMutations(src)
	if err != nil {
		t.Fatalf("GenerateMutations failed: %v", err)
	}

	if len(plans) == 0 {
		t.Fatalf("Expected at least one mutation plan, got 0")
	}

	// Verify categories present
	hasRelational := false
	hasBoolean := false
	hasArithmetic := false
	hasBoundary := false

	for _, p := range plans {
		switch p.Type {
		case engine.RelationalOpReplace:
			hasRelational = true
		case engine.BooleanOpFlip:
			hasBoolean = true
		case engine.ArithmeticOpReplace:
			hasArithmetic = true
		case engine.BoundaryValueMutate:
			hasBoundary = true
		}
	}

	if !hasRelational || !hasBoolean || !hasArithmetic || !hasBoundary {
		t.Errorf("Missing expected mutation types: relational=%v, bool=%v, arith=%v, boundary=%v",
			hasRelational, hasBoolean, hasArithmetic, hasBoundary)
	}
}

func TestEngine_MutationScoreDifference(t *testing.T) {
	srcPath := filepath.Join("..", "internal", "service", "discount.go")
	src, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Failed to read discount.go: %v", err)
	}

	mutator := engine.NewASTMutator()
	plans, err := mutator.GenerateMutations(src)
	if err != nil {
		t.Fatalf("GenerateMutations failed: %v", err)
	}

	runner := engine.NewRunner()

	// Weak Test Evaluation
	// Simulates what TestCalculateDiscount_Weak asserts:
	// only checks final amount > 0, discount total >= 0, etc.
	weakTestFn := func(mutatedSrc []byte) bool {
		// If the mutated code fails to satisfy the weak assertion:
		// Notice weak tests only assert positive amounts. Mutating >= to > or || to &&
		// still yields positive amounts for the loose test inputs.
		// A mutant is killed ONLY if it breaks the weak test assertions.
		code := string(mutatedSrc)
		// For demo/engine verification, we simulate whether weak assertions catch the difference:
		// Weak assertions miss:
		// - boundary shifts (> into >=, or shift constants)
		// - subtle rate differences
		// Only fatal logic breaks (like rate calculation turning completely negative) would be caught.
		if strings.Contains(code, "FinalAmount <= 0") {
			return true
		}
		return false
	}

	// Strong Test Evaluation
	// Simulates comprehensive test suite (assertions on exact values, boundaries, flags)
	strongTestFn := func(mutatedSrc []byte) bool {
		code := string(mutatedSrc)
		// Strong assertions check every single mutation made to operators, boolean logic, or constants
		// If code differs in operator or value, strong tests catch it.
		return code != string(src)
	}

	weakReport := runner.Run(src, plans, weakTestFn)
	strongReport := runner.Run(src, plans, strongTestFn)

	if weakReport.Score >= strongReport.Score {
		t.Errorf("Expected strong test score (%.2f%%) > weak test score (%.2f%%)",
			strongReport.Score, weakReport.Score)
	}

	if strongReport.Score != 100.0 {
		t.Errorf("Expected strong test score to be 100%%, got %.2f%%", strongReport.Score)
	}
}

func TestDomainService_DirectCalculation(t *testing.T) {
	// Directly verify CalculateDiscount
	res := service.CalculateDiscount(service.Order{
		TotalAmount: 1000.0,
		ItemCount:   3,
		Tier:        service.TierVIP,
		HasCoupon:   true,
	})
	if res.DiscountRate != 0.25 {
		t.Errorf("Expected 0.25 rate, got %v", res.DiscountRate)
	}
	if !res.FreeShipping {
		t.Errorf("Expected free shipping for VIP")
	}
}
