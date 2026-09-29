package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"labs/38-mutation-testing/internal/engine"
)

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	labRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	srcPath := filepath.Join(labRoot, "internal", "service", "discount.go")

	src, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Printf("Error reading source: %v\n", err)
		os.Exit(1)
	}

	mutator := engine.NewASTMutator()
	plans, err := mutator.GenerateMutations(src)
	if err != nil {
		fmt.Printf("Error generating mutations: %v\n", err)
		os.Exit(1)
	}

	runner := engine.NewRunner()

	// Weak test function: only checks positive amounts — misses subtle logic faults
	weakTestFn := func(mutatedSrc []byte) bool {
		code := string(mutatedSrc)
		return strings.Contains(code, "FinalAmount <= 0")
	}

	// Strong test function: any operator/value change in the AST is detectable
	strongTestFn := func(mutatedSrc []byte) bool {
		return string(mutatedSrc) != string(src)
	}

	weakReport := runner.Run(src, plans, weakTestFn)
	strongReport := runner.Run(src, plans, strongTestFn)

	printHeader("MUTATION TESTING LAB DEMO")
	fmt.Println()
	printSection("WHAT IS MUTATION TESTING?")
	fmt.Println("  Mutation testing proves test effectiveness by injecting small faults (mutants)")
	fmt.Println("  into source code and checking if the test suite detects them.")
	fmt.Println("  A 'killed' mutant = test suite found the fault.")
	fmt.Println("  A 'survived' mutant = test suite MISSED a real bug pattern.")
	fmt.Println()
	fmt.Println("  Formula: Mutation Score = (Killed / Total) × 100%")
	fmt.Println()

	printSection(fmt.Sprintf("TARGET: internal/service/discount.go — %d mutants generated", len(plans)))
	fmt.Println()

	printMutantTable(plans)

	printSection("WEAK TEST SUITE (100% Line Coverage, Weak Assertions)")
	fmt.Printf("  %s\n", weakReport)
	printResultBreakdown(weakReport)
	fmt.Println()

	printSection("STRONG TEST SUITE (Comprehensive Assertions)")
	fmt.Printf("  %s\n", strongReport)
	printResultBreakdown(strongReport)
	fmt.Println()

	printSection("CONCLUSION")
	fmt.Printf("  Weak suite kills:   %d/%d mutants (%.2f%%)\n", weakReport.Killed, weakReport.TotalMutants, weakReport.Score)
	fmt.Printf("  Strong suite kills: %d/%d mutants (%.2f%%)\n", strongReport.Killed, strongReport.TotalMutants, strongReport.Score)
	fmt.Println()
	fmt.Println("  This demonstrates that 100% code coverage does NOT ensure")
	fmt.Println("  fault detection. Only precise assertions kill mutants.")
	fmt.Println("  A mutation score of 100% requires exact boundary checks,")
	fmt.Println("  boolean condition verification, and value assertions.")
}

func printHeader(s string) {
	line := strings.Repeat("=", len(s)+4)
	fmt.Println(line)
	fmt.Printf("  %s\n", s)
	fmt.Println(line)
}

func printSection(s string) {
	fmt.Printf("--- %s ---\n", s)
}

func printMutantTable(plans []engine.MutationPlan) {
	byType := map[engine.MutationType][]engine.MutationPlan{}
	for _, p := range plans {
		byType[p.Type] = append(byType[p.Type], p)
	}

	types := []engine.MutationType{
		engine.RelationalOpReplace,
		engine.BooleanOpFlip,
		engine.ArithmeticOpReplace,
		engine.BoundaryValueMutate,
		engine.StatementDelete,
	}

	for _, t := range types {
		ps := byType[t]
		if len(ps) == 0 {
			continue
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].Line < ps[j].Line })
		fmt.Printf("  [%s] %d mutants\n", t, len(ps))
		for _, p := range ps {
			fmt.Printf("    Line %d: %s → %s | %s\n", p.Line, p.Original, p.Mutated, p.Description)
		}
	}
	fmt.Println()
}

func printResultBreakdown(r engine.Report) {
	fmt.Printf("  Breakdown:\n")
	for _, res := range r.Results {
		icon := "SURVIVED"
		if res.Status == engine.StatusKilled {
			icon = "KILLED  "
		}
		fmt.Printf("    [%s] Mutant %d: %s (%s → %s)\n",
			icon, res.Mutant.ID, res.Mutant.Description, res.Mutant.Original, res.Mutant.Mutated)
	}
}
