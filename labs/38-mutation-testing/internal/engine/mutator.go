package engine

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
)

type ASTMutator struct {
	Fset *token.FileSet
}

func NewASTMutator() *ASTMutator {
	return &ASTMutator{
		Fset: token.NewFileSet(),
	}
}

type MutationPlan struct {
	Type        MutationType
	Description string
	Line        int
	Original    string
	Mutated     string
	Apply       func(file *ast.File) func()
}

// GenerateMutations inspects Go source code and returns executable mutation plans.
func (m *ASTMutator) GenerateMutations(sourceCode []byte) ([]MutationPlan, error) {
	parsedFile, err := parser.ParseFile(m.Fset, "source.go", sourceCode, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var plans []MutationPlan

	ast.Inspect(parsedFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch expr := n.(type) {
		case *ast.BinaryExpr:
			pos := m.Fset.Position(expr.Pos())

			// 1. Relational operators
			switch expr.Op {
			case token.GEQ:
				plans = append(plans, MutationPlan{
					Type:        RelationalOpReplace,
					Description: "Replace '>=' with '>'",
					Line:        pos.Line,
					Original:    ">=",
					Mutated:     ">",
					Apply: func(f *ast.File) func() {
						expr.Op = token.GTR
						return func() { expr.Op = token.GEQ }
					},
				})
			case token.GTR:
				plans = append(plans, MutationPlan{
					Type:        RelationalOpReplace,
					Description: "Replace '>' with '>='",
					Line:        pos.Line,
					Original:    ">",
					Mutated:     ">=",
					Apply: func(f *ast.File) func() {
						expr.Op = token.GEQ
						return func() { expr.Op = token.GTR }
					},
				})
			case token.EQL:
				plans = append(plans, MutationPlan{
					Type:        RelationalOpReplace,
					Description: "Replace '==' with '!='",
					Line:        pos.Line,
					Original:    "==",
					Mutated:     "!=",
					Apply: func(f *ast.File) func() {
						expr.Op = token.NEQ
						return func() { expr.Op = token.EQL }
					},
				})

			// 2. Boolean operators
			case token.LOR:
				plans = append(plans, MutationPlan{
					Type:        BooleanOpFlip,
					Description: "Replace '||' with '&&'",
					Line:        pos.Line,
					Original:    "||",
					Mutated:     "&&",
					Apply: func(f *ast.File) func() {
						expr.Op = token.LAND
						return func() { expr.Op = token.LOR }
					},
				})
			case token.LAND:
				plans = append(plans, MutationPlan{
					Type:        BooleanOpFlip,
					Description: "Replace '&&' with '||'",
					Line:        pos.Line,
					Original:    "&&",
					Mutated:     "||",
					Apply: func(f *ast.File) func() {
						expr.Op = token.LOR
						return func() { expr.Op = token.LAND }
					},
				})

			// 3. Arithmetic operators
			case token.SUB:
				plans = append(plans, MutationPlan{
					Type:        ArithmeticOpReplace,
					Description: "Replace '-' with '+'",
					Line:        pos.Line,
					Original:    "-",
					Mutated:     "+",
					Apply: func(f *ast.File) func() {
						expr.Op = token.ADD
						return func() { expr.Op = token.SUB }
					},
				})
			case token.MUL:
				plans = append(plans, MutationPlan{
					Type:        ArithmeticOpReplace,
					Description: "Replace '*' with '/'",
					Line:        pos.Line,
					Original:    "*",
					Mutated:     "/",
					Apply: func(f *ast.File) func() {
						expr.Op = token.QUO
						return func() { expr.Op = token.MUL }
					},
				})
			}

		case *ast.BasicLit:
			// 4. Boundary / Value Mutation for integers and floats
			pos := m.Fset.Position(expr.Pos())
			if expr.Kind == token.INT {
				val, err := strconv.Atoi(expr.Value)
				if err == nil && val > 0 {
					orig := expr.Value
					plans = append(plans, MutationPlan{
						Type:        BoundaryValueMutate,
						Description: "Shift integer constant +1",
						Line:        pos.Line,
						Original:    orig,
						Mutated:     strconv.Itoa(val + 1),
						Apply: func(f *ast.File) func() {
							expr.Value = strconv.Itoa(val + 1)
							return func() { expr.Value = orig }
						},
					})
				}
			}
		}

		return true
	})

	return plans, nil
}

// RenderSource outputs the modified AST back into source code.
func (m *ASTMutator) RenderSource(file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, m.Fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
