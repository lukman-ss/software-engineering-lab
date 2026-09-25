# Implementation Notes

## Files Added
- `internal/adr/models.go`: Defines the core ADR structs and valid lifecycle statuses (`Proposed`, `Accepted`, `Superseded`, `Deprecated`, `Rejected`).
- `internal/adr/parser.go`: Implements regex-based markdown parsing to extract ID, Title, Status, and supersession relationships.
- `internal/adr/linter.go`: Validates monotonic numbering and structural integrity of the decision graph (no broken or one-sided supersession links).
- `tests/*_test.go`: Unit tests for parser and linter, proving constraints.
- `cmd/demo/main.go`: Executable simulation using the Modular Monolith vs. Microservices scenario from the research.

## Core Design Decisions
- **Regex over AST**: Used line-by-line regex parsing rather than a full Markdown AST library to keep dependencies minimal (standard library only) while maintaining functional sufficiency for the structural rules.
- **In-Memory Validation Graph**: ADR relationships (supersession) are validated using a unified in-memory map across all parsed records, executed concurrently using Goroutines to prove concurrency safety.

## Implementation-Specific Choices
- Rather than a shell script and physical markdown files, the lab implements a programmable linter in Go that reads strings. This proves the logic of the structural rules without requiring complex file IO mocking in unit tests.
- Status strings are strict subsets defined by the AWS Prescriptive Guidance (including `Rejected`).

## Known Limitations
- The parser expects a very specific Markdown heading format (`# 1. Title`) and metadata format (`Status: Accepted`). It is not a robust generalized Markdown parser.
- Does not automatically perform Git file manipulation (e.g. creating/updating physical files).

## Trade-offs
- Traded broad markdown syntax support for strict zero-dependency validation. 

## What Is Demonstrated
- Co-location of decisions with code (conceptually via tooling).
- Strict monotonic numbering enforcement.
- Immutable history through `Superseded` and `Supersedes` bidirectional validation.
- Preservation of rejected choices.

## What Is Not Demonstrated
- Integration with Git pre-commit hooks.
- Generating new boilerplate ADR markdown files.
