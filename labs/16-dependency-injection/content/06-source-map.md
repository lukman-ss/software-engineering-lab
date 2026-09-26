# Source Map

## Separation of Configuration from Use

Research:
- `research/05-report.md` (Finding 1 — DI definition and construction/use separation)
- `research/03-evidence.md` (Evidence 1, 4, 11 — DI definition, reduced coupling, testability)
- `research/02-sources.md` (Source 1 — Martin Fowler, Source 7 — Microsoft .NET)

Engine:
- `engineering/01-design.md` (Expected Behavior, Architecture)
- `engineering/02-implementation-notes.md` (Core Design Decisions, Implementation-Specific Choices)

Implementation:
- `cmd/demo/main.go` (composition root — externalized object creation)
- `internal/di/gateway.go` (RealGateway, PaymentGateway interface)
- `internal/di/processor.go` (Processor with constructor injection)

## Core Concept: IoC vs DI

Research:
- `research/05-report.md` (Finding 1 — construction/use separation; Finding 2 — IoC broader than DI)
- `research/03-evidence.md` (Evidence 2 — DI as inversion over implementation; Evidence 12 — Hollywood Principle; Evidence 13 — Java "IoC" terminology)
- `research/04-contradictions.md` (no material contradictions — terminology clarified)
- `research/02-sources.md` (Source 2 — Wikipedia IoC; Source 5 — Spring Framework)

Implementation:
- `internal/di/processor.go` (Processor — DI pattern)
- `internal/di/locator.go` (BadProcessor — Service Locator anti-pattern)

## Constructor Injection

Research:
- `research/05-report.md` (Finding 3 — three DI forms; Finding 4 — prefer constructor injection)
- `research/03-evidence.md` (Evidence 7 — three forms; Evidence 8 — start with constructor, switch to setter when needed)
- `research/02-sources.md` (Source 1 — Fowler 2004)

Code Audit:
- `engineering-audit/02-code-audit.md` (Finding 1 — PASS: explicit constructor injection)

Implementation:
- `internal/di/processor.go` (NewProcessor receives PaymentGateway via constructor)

## Service Locator Anti-Pattern

Research:
- `research/05-report.md` (Finding 5 — DI vs Service Locator)
- `research/03-evidence.md` (Evidence 9, 10 — Service Locator creates locator dependency; acceptable for app-internal code)
- `research/05-report.md` (Finding 9 — PSR-11 discourages passing container into objects: "SHOULD NOT")
- `research/03-evidence.md` (Evidence 17 — PSR-11 ContainerInterface spec)
- `research-revision/02-changes-made.md` (Revision 1 — fixed RFC 2119 "SHOULD NOT" vs "MUST NOT")

Code Audit:
- `engineering-audit/02-code-audit.md` (Finding 2 — PASS: Service Locator couples BadProcessor to Container interface)

Implementation:
- `internal/di/locator.go` (Container interface, BadProcessor struct, NewBadProcessor constructor)

## Testing via Mock Substitution

Research:
- `research/05-report.md` (Finding 6 — DI enables mock/stub testing)
- `research/03-evidence.md` (Evidence 5 — ease of testing is first benefit noticed)
- `engineering/01-design.md` (Test Strategy — MockGateway into Processor)
- `engineering/03-execution-result.md` (Test results — all PASS, race detector PASS)

Test Audit:
- `engineering-audit/03-test-audit.md` (6 test cases covered — all PASS)

Implementation:
- `tests/processor_test.go` (MockGateway with ShouldFail, ChargedMoney verification; 6 test cases)

## Value Objects Bypass DI

Research:
- `research/05-report.md` (Finding 12 — don't inject value objects; specific examples are lab heuristics)
- `research/04-contradictions.md` (Divergence 4 — topic spec list, Fowler principle only)
- `research-revision/02-changes-made.md` (Revision 2 — qualified as lab heuristics)
- `research/03-evidence.md` (Evidence 12 — Value Objects Finding)

Code Audit:
- `engineering-audit/02-code-audit.md` (Finding 3 — PASS: Money instantiated directly)

Implementation:
- `internal/di/gateway.go` (Money struct — no interface, no container)
- `internal/di/processor.go` (Money{...} instantiated inline)
- `internal/di/locator.go` (Money{...} instantiated inline)

## Failure Scenarios (Input Validation + Error Propagation)

Research:
- `research/05-report.md` (Finding 10 — DI costs/trade-offs)

Code Audit:
- `engineering-audit/02-code-audit.md` (Finding 4 — PASS: validation before gateway call)

Test Audit:
- `engineering-audit/03-test-audit.md` (TestProcessor_InvalidAmount, TestBadProcessor_InvalidAmount prove gateway not called on invalid input; TestProcessor_GatewayError, TestBadProcessor_GatewayError prove error propagation)

Implementation:
- `internal/di/processor.go:15-18` (if amount <= 0 → return errors.New("invalid amount"))
- `internal/di/locator.go:19-22` (same validation in Service Locator variant)

## Demo Output

Execution:
- `engineering/03-execution-result.md` (Demo result: "RealGateway charging 100 USD" and "RealGateway charging 200 USD")
- `engineering-audit/03-test-audit.md` (go run ./cmd/demo matches)

Implementation:
- `cmd/demo/main.go` (main function demonstrates both patterns)

## Known Limitations of This Lab

Research:
- `research/05-report.md` (Finding 10 — DI has costs; trade-offs acknowledged)
- `research/06-open-questions.md` (Q1 — no empirical evidence for defect reduction; Q2 — performance overhead not quantified)
- `research/03-claim-audit.md` (Claim 10 — 12-parameter threshold is lab heuristic, not industry standard)

Engineering:
- `engineering/02-implementation-notes.md` (Known Limitations — no reflection-based or code-gen DI container; no lifecycle management singleton/scoped/transient)
