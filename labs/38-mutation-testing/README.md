# Lab 38: Mutation Testing

Demonstrates how AST-based mutation testing exposes false confidence in unit test suites achieving 100% line coverage.

## Overview

Mutation testing evaluates test suite quality by injecting small syntax-level faults (mutants) into source code and verifying whether existing tests detect (kill) them.

Formula:

$$\text{Mutation Score} = \left(\frac{\text{Killed Mutants}}{\text{Total Mutants}}\right) \times 100\%$$

## Structure

```text
labs/38-mutation-testing/
├── cmd/demo/main.go            # CLI demonstration of mutation testing score contrast
├── internal/
│   ├── engine/                 # AST mutation generator & concurrent runner
│   │   ├── mutator.go          # Relational, boolean, arithmetic, & boundary AST mutator
│   │   ├── runner.go           # Parallel mutant test execution engine
│   │   └── types.go            # Data structures for mutants & reports
│   └── service/                # Domain logic (discount calculation engine)
│       ├── discount.go         # Core pricing and discount rules
│       ├── discount_weak_test.go   # 100% line coverage weak assertions
│       └── discount_strong_test.go # Comprehensive boundary assertions
├── tests/
│   └── engine_test.go          # Unit tests verifying mutation engine behavior
├── engineering/                # Engineering docs (01-design, 02-notes, 03-execution)
└── go.mod
```

## Running the Lab

### Run Unit Tests & Race Detector

```bash
go test -v ./...
go test -race ./...
```

### Run Demo

```bash
go run ./cmd/demo
```

## Implemented Mutation Operators

1. **Relational Operator Replacement**: `>` $\leftrightarrow$ `>=`, `==` $\leftrightarrow$ `!=`
2. **Boolean Flip**: `&&` $\leftrightarrow$ `||`
3. **Arithmetic Operator Replacement**: `*` $\leftrightarrow$ `/`, `-` $\leftrightarrow$ `+`
4. **Boundary Value Shift**: Integer constants shifted $+1$
