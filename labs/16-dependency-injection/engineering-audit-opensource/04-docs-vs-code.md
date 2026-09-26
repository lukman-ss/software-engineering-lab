# Docs vs Code

## Comparison: README vs Code
- README Finding 1 ("Separation of Configuration from Use ... externalized in main.go"): Verified — `cmd/demo/main.go` constructs `RealGateway` and injects it into both processors.
- README Finding 2 ("Fast, Isolated Unit Testing ... mocked without real network calls"): Verified — `tests/processor_test.go` uses `MockGateway` with no network.
- README Finding 3 ("Constructor Injection: NewProcessor ensures components are fully initialized with explicit dependencies"): Verified — `NewProcessor(g)` sets immutable dependency.
- README Finding 4 ("Service Locator Anti-Pattern: NewBadProcessor injects a Container"): Verified — `NewBadProcessor(c)` stores a Container.
- README Finding 5 ("Value Objects Bypass DI: Money is directly instantiated"): Verified — `Money{...}` created inline.
- README Usage commands (`go test -race ./...`, `go run ./cmd/demo`): Verified — both executed successfully.

## Comparison: Engineering Design vs Code
- Design claims map 1:1 to implementation (PaymentGateway interface, RealGateway, Processor, BadProcessor, Money).
- Design "Failure Scenario" (negative amount rejects without calling gateway; mock gateway error propagates) — both confirmed in code and tests.
- Implementation notes state manual DI over framework (YAGNI) — consistent with code (no framework imports).
- Execution result doc matches audit rerun exactly (same demo output).
- Decision scoping correct: design explicitly defers DI containers (dig/wire) and lifecycle management as out-of-scope; README/disclaimers consistent.

## Discrepancies
- None found. No DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, or RESEARCH_IMPLEMENTATION_MISMATCH detected.

## Concern Worth Noting (not a mismatch)
- Zero-amount input: code treats `0` as invalid (`amount <= 0`), but no test covers `0`. Docs/design do not distinguish `0` from negative; behavior is reasonable but untested. Filed as MISSING_EDGE_CASE / MISSING_TEST with LOW severity in gaps.