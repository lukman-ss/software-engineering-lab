# Lab 40: Property-Based Testing (PBT)

Demonstration of Property-Based Testing (PBT) vs Example-Based Testing in Go.

## Overview
This lab demonstrates how Example-Based Unit Testing can generate false confidence by passing hand-picked examples while missing boundary and edge cases. In contrast, Property-Based Testing tests universal invariants against hundreds of randomized inputs and automatically shrinks failing inputs to minimal counterexamples.

## Canonical Invariants Demonstrated
1. **Roundtrip Invariant (`Decode(Encode(x)) == x`)**:
   - `internal/currency`: Verifies monetary formatting and parsing. Float64 implementations suffer from precision loss, whereas integer cent implementations pass 1,000 randomized roundtrips across negative, zero, and multi-billion cent bounds.
2. **Idempotence Invariant (`f(f(x)) == f(x)`)**:
   - `internal/interval`: Verifies interval merging. Repeated merges on an already-merged slice yield an identical slice.
3. **Oracle / Equivalence Invariant (`f_naive(x) == f_robust(x)`)**:
   - Highlights failure of naive interval merging on unsorted inputs compared to robust interval merging.
4. **Counterexample Shrinking**:
   - `internal/shrinker`: Reduces complex 10+ element failing arrays to a minimal single-element counterexample (`[-1]`).

## Directory Layout
```text
.
├── cmd/
│   └── demo/
│       └── main.go           # Interactive CLI comparison demo
├── internal/
│   ├── currency/             # Roundtrip invariant on currency parsing/formatting
│   │   ├── currency.go
│   │   └── currency_test.go
│   ├── interval/             # Idempotence and oracle invariants on interval merging
│   │   ├── interval.go
│   │   └── interval_test.go
│   └── shrinker/             # Binary sectioning & element shrinking engine
│       ├── shrinker.go
│       └── shrinker_test.go
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── go.mod
└── README.md
```

## Running the Lab

### 1. Run Unit and Property Tests
```bash
go test -v ./...
```

### 2. Run Race Detector
```bash
go test -race ./...
```

### 3. Run Demo
```bash
go run ./cmd/demo
```
