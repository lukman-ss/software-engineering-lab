# Engineering Design

Target Lab: labs/17-architecture-decision-record
Research Status: APPROVED_WITH_WARNINGS

## Concept To Prove
Architectural Decision Records (ADRs) provide an immutable, version-controlled mechanism for documenting architecturally significant decisions. A structural validation process ensures that ADRs maintain proper state transitions (e.g., Proposed -> Accepted -> Superseded), follow monotonic numbering, and preserve a valid decision lineage (DAG) without broken references.

## Expected Behavior
- **Parsing**: Can extract Title, Status, and Superseded references from a Markdown ADR.
- **Linting Valid**: An ADR sequence with monotonic numbering, correct statuses (Proposed, Accepted, Superseded, Deprecated, Rejected), and valid supersession references passes validation.
- **Linting Invalid**: Validation fails on missing files, broken supersession references, invalid statuses, or structural defects.
- **Supersession DAG**: Validates that if ADR-B supersedes ADR-A, ADR-A exists and its status correctly reflects it is superseded by ADR-B.

## Failure Scenario
- A superseded ADR points to a non-existent replacement.
- An accepted ADR is silently mutated instead of being superseded by a new ADR.
- An ADR uses an invalid lifecycle state (e.g., "Draft" instead of "Proposed").

## Success Criteria
1. ADR parser correctly extracts structured data from Markdown files.
2. Linter correctly accepts a valid sequence of ADRs (including the Modular Monolith to Microservices scenario from research).
3. Linter correctly rejects invalid structural states (broken links, unknown status).
4. Concurrency safety verified with Go race detector on the linter processing multiple files.

## Architecture
- `adr.Parser`: Reads markdown text and extracts metadata.
- `adr.Linter`: Validates a collection of parsed ADRs against lifecycle rules.
- `Fake File System / In-Memory Repo`: To simulate the `docs/adr` directory for fast testing.

## Components
- `internal/adr`: Core parser and linter implementation.
- `cmd/demo`: Executable demonstration script running validation on sample ADRs.
- `tests`: Unit tests and integration tests.

## Test Strategy
- Unit tests for regex-based markdown parsing.
- Unit tests for structural validation rules.
- Happy path integration test modeling the SaaS ERP scenario (Modular Monolith -> Microservices).

## Execution Plan
1. Implement `adr.Parser`.
2. Implement `adr.Linter`.
3. Implement `cmd/demo/main.go` demonstrating the research scenario.
4. Execute tests and race detector.

## Implementation Decisions
- ADR structure will be simple Markdown with specific headers (Title, Status).
- Status field will be parsed via simple regex matching to simulate a structural linter without requiring complex markdown AST parsing.
- Implementation acts as a programmatic tool (linter) rather than just static files, proving the structural integrity rules defined in the research.
