package engine

import "fmt"

type MutationType string

const (
	RelationalOpReplace MutationType = "RelationalOpReplace"
	BooleanOpFlip       MutationType = "BooleanOpFlip"
	ArithmeticOpReplace MutationType = "ArithmeticOpReplace"
	BoundaryValueMutate MutationType = "BoundaryValueMutate"
	StatementDelete     MutationType = "StatementDelete"
)

type MutantStatus string

const (
	StatusKilled    MutantStatus = "KILLED"
	StatusSurvived  MutantStatus = "SURVIVED"
	StatusEquivalent MutantStatus = "EQUIVALENT"
)

type Mutant struct {
	ID          int
	Type        MutationType
	Description string
	LineNumber  int
	Original    string
	Mutated     string
}

type MutantResult struct {
	Mutant Mutant
	Status MutantStatus
	Output string
}

type Report struct {
	TotalMutants int
	Killed       int
	Survived     int
	Score        float64
	Results      []MutantResult
}

func (r Report) String() string {
	return fmt.Sprintf("Total: %d, Killed: %d, Survived: %d, Score: %.2f%%",
		r.TotalMutants, r.Killed, r.Survived, r.Score)
}
