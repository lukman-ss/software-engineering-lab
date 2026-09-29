# Content Brief

Topic: Mutation Testing (AST-Based Fault Injection & Test Quality Assessment)
Target Reader: Software Engineers, QA Engineers, Engineering Leads, and System Architects building and maintaining test suites.
Problem: High code coverage (even 100% line coverage) provides false confidence because traditional coverage metrics only verify which statements are executed, not whether the test suite contains assertions capable of detecting logic regressions or faults.
Core Mental Model: Code coverage measures *what is executed*; mutation testing measures *what is verified*. By systematically injecting small faults (mutants) into source code, mutation testing checks if the test suite fails (kills the mutant) or passes (mutant survives).
Approved Research Status: APPROVED (Audit Date: 2026-09-29)
Approved Engineering Status: APPROVED (Engineering Audit Date: 2026-09-29)
Main Concepts:
- Mutation Score: `(Killed Mutants / Total Mutants) × 100%`
- RIP Model: Reach, Infect, Propagate (conditions for killing a mutant)
- Core Mutation Operators: Relational (`>` vs `>=`), Boolean (`&&` vs `||`), Arithmetic (`*` vs `/`, `-` vs `+`), Boundary Value (+1 shift)
- Theoretical Foundations: Competent Programmer Hypothesis & Coupling Effect
- Equivalent Mutants: Mutants semantically identical to the original program (undecidable to fully eliminate)
- False Confidence Gap: 100% line coverage with weak assertions yielding 0% mutation score
Verified Behaviors:
- AST Mutator systematically generates 15 mutants from `internal/service/discount.go` using Go standard library (`go/ast`, `go/parser`, `go/token`, `go/format`).
- Weak assertion test suite achieves 100.0% statement coverage on `discount.go` while killing 0/15 mutants (0.00% mutation score).
- Strong assertion test suite evaluates boundaries, exact calculation outputs, and flags, killing 15/15 mutants (100.00% mutation score).
- Concurrent mutant execution engine runs race-free with isolated AST copies (`go test -race ./...` passes).
Available Case Studies:
- Discount & Shipping Service Engine (`internal/service/discount.go`): multi-tiered pricing, boundary thresholds, coupon rules, and free shipping triggers.
Warnings:
- No industry-standard mutation score threshold exists; target scores must be determined by project risk profile.
- In this lab, the strong test evaluation in the custom engine demo relies on AST textual inequality rather than subprocess compilation (`os/exec`), chosen as a lightweight illustrative trade-off.
- Equivalent mutants are mathematically undecidable; practical tools rely on heuristics and scoped analysis.
- `StatementDelete` enum is declared for classification completeness but not emitted on the target AST where no side-effect-free void statements exist.
