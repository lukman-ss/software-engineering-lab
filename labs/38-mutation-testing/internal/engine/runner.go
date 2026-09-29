package engine

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sync"
)

// TestFunc is a function that takes an AST-modified source and runs assertions against it.
// It returns true if the test detects the mutation (i.e., kills it).
type TestFunc func(src []byte) bool

// Runner executes a test suite against all mutation plans and collects a Report.
type Runner struct {
	mu sync.Mutex
}

func NewRunner() *Runner {
	return &Runner{}
}

// Run applies each mutation plan, runs testFn against the mutated code, and collects results.
// It uses a fresh AST parse per mutant to avoid interference between mutations.
func (r *Runner) Run(sourceCode []byte, plans []MutationPlan, testFn TestFunc) Report {
	results := make([]MutantResult, len(plans))

	var wg sync.WaitGroup
	for i, plan := range plans {
		wg.Add(1)
		go func(idx int, p MutationPlan) {
			defer wg.Done()

			fset := token.NewFileSet()
			parsedFile, err := parser.ParseFile(fset, "source.go", sourceCode, 0)
			if err != nil {
				results[idx] = MutantResult{
					Mutant: Mutant{
						ID:          idx + 1,
						Type:        p.Type,
						Description: p.Description,
						LineNumber:  p.Line,
						Original:    p.Original,
						Mutated:     p.Mutated,
					},
					Status: StatusEquivalent,
					Output: "parse error",
				}
				return
			}

			undo := p.Apply(parsedFile)

			mutatedSrc, err := renderSource(fset, parsedFile)
			if err != nil {
				undo()
				results[idx] = MutantResult{
					Mutant: Mutant{
						ID:          idx + 1,
						Type:        p.Type,
						Description: p.Description,
						LineNumber:  p.Line,
						Original:    p.Original,
						Mutated:     p.Mutated,
					},
					Status: StatusEquivalent,
					Output: "render error",
				}
				return
			}

			killed := testFn(mutatedSrc)
			undo()

			status := StatusSurvived
			if killed {
				status = StatusKilled
			}

			results[idx] = MutantResult{
				Mutant: Mutant{
					ID:          idx + 1,
					Type:        p.Type,
					Description: p.Description,
					LineNumber:  p.Line,
					Original:    p.Original,
					Mutated:     p.Mutated,
				},
				Status: status,
			}
		}(i, plan)
	}
	wg.Wait()

	total := len(plans)
	killed := 0
	for _, res := range results {
		if res.Status == StatusKilled {
			killed++
		}
	}

	score := 0.0
	if total > 0 {
		score = float64(killed) / float64(total) * 100.0
	}

	return Report{
		TotalMutants: total,
		Killed:       killed,
		Survived:     total - killed,
		Score:        score,
		Results:      results,
	}
}

// renderSource is an independent helper used by Runner (does not rely on the mutator state).
func renderSource(fset *token.FileSet, file *ast.File) ([]byte, error) {
	m := &ASTMutator{Fset: fset}
	return m.RenderSource(file)
}
