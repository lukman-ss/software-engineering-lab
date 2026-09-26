# Source Map

Sections of `02-master-draft.md` mapped to source files in `labs/26-contract-testing`.

## Introduction & Problem

- Research: `research/05-report.md` — Executive Summary, Limitations, Conclusion.
- Lab README: `README.md` — Overview and problem statement.
- Implementation: `internal/consumer/client.go`, `internal/provider/server.go` — Mobile App / Order Service example.

## Why This Matters

- Research: `research/03-evidence.md` Evidence 5 — Unit tests passing while contract breaks.
- Engineering: `engineering/03-execution-result.md` — Demo output showing breaking change detected.

## Mental Model

- Research: `research/05-report.md` — Consumer-Driven Contract definition, Pact mechanism.
- Engineering: `engineering/01-design.md` — Architecture diagram and expected behavior.

## Core Concept

- Research: `research/05-report.md` Findings 1–4, 7, 8, 9 — CDC definition, contract contents, breaking change detection, expand/contract, Pact Broker.
- Engineering: `engineering/01-design.md` — Verifier engine purpose.
- Code: `internal/contract/verifier.go` — `Verify()` method.

## Failure Scenario

- Research: `research/05-report.md` — Finding 5 (breaking change causes production failure), Finding 5 (additive vs breaking).
- Engineering: `engineering/01-design.md` — Failure scenario section.
- Engineering execution: `engineering/03-execution-result.md` — Stage 3 output with 3 diffs.

## How It Works

- Research: `research/03-evidence.md` Evidence 3 — Pact generation/verification flow.
- Engineering: `engineering/01-design.md` — Architecture and component list.
- Code: `internal/contract/verifier.go` — `Verify()` and `diffValues()`.

## Architecture

- Engineering: `engineering/01-design.md` — Architecture diagram, component list, test strategy.
- Code: `internal/model/order.go`, `internal/provider/server.go`, `internal/consumer/client.go`, `internal/contract/verifier.go`.

## Implementation

- Engineering: `engineering/02-implementation-notes.md` — Files added, design decisions, ponytail notes.
- Code: `internal/contract/verifier.go` — `decoder.UseNumber()` for type preservation.

## Code Walkthrough

- Code snippets source mapping:
  - Snippet 1: `internal/consumer/client.go:82-111` — `GenerateMobileContract()`.
  - Snippet 2: `internal/provider/server.go:12-44` — `ProviderV1`.
  - Snippet 3: `internal/provider/server.go:46-80` — `ProviderBreaking`.
  - Snippet 4: `internal/provider/server.go:82-131` — `ProviderDual`.
  - Snippet 5: `internal/consumer/client.go:51-67` — `FetchOrder` field validation.
  - Snippet 6: `internal/contract/verifier.go:57-115` — `Verify()` method.
  - Snippet 7: `internal/contract/verifier.go:117-176` — `diffValues()` engine.
  - Snippet 8: `internal/model/order.go` — DTO definitions.
  - Snippet 9: `cmd/demo/main.go:14-68` — Demo orchestrator.

## What the Tests Prove

- Tests: `tests/contract_test.go` — Full test suite.
  - `TestConsumerContractGeneration` — Line 13.
  - `TestProviderV1_ContractVerification_Success` — Line 27.
  - `TestProviderBreaking_ContractVerification_Fails` — Line 50.
  - `TestProviderDual_ContractVerification_Success` — Line 74.
  - `TestConcurrentContractVerification` — Line 96.

## Recovery / Rollback

- Research: `research/05-report.md` Finding 7 — expand/contract pattern phases.
- Engineering: `engineering/01-design.md` — Safe API Evolution point 4.
- Code: `internal/provider/server.go:82-131` — ProviderDual dual routing implementation.

## Production Considerations

- Research: `research/05-report.md` — Limitations (collaboration, business logic, test maintenance, E2E complement).
- Research: `research/02-sources.md` Source 10, 11 — Spring Cloud Contract archived (July 2026), AsyncAPI positioning.
- Engineering audit: `engineering-audit-opensource/06-verdict.md` — 6 LOW gaps (array diff recursion, response headers, V2 untested, no timeout, etc.).

## Common Mistakes

- Research: `research/05-report.md` — Limitations, Conclusion.
- Research: `research/03-evidence.md` Evidence 6 — anti-pattern of testing validation rules in contracts (Source 6).
- Engineering: `engineering/01-design.md` — ponytail about minimal subset rule.

## Case Study

- Lab README: `README.md` — Overview, project structure.
- Engineering: `engineering/03-execution-result.md` — Stage 1–4 demo output.

## Checklist

- Research: `research/05-report.md` — Findings 1–10, Limitations.
- Engineering: `engineering/01-design.md` — Success criteria.

## Key Takeaways

- Derived from all sections above.

## Sources

- Research: `research/02-sources.md` — 11 sources with tiers and relevance.
- Research: `research/05-report.md` — Areas of agreement/disagreement, Conclusion.
- Engineering: `engineering-audit/06-verdict.md` — Quality gates (Pass: compilation, tests, race, demo, alignment).
- Engineering: `engineering-audit-opensource/06-verdict.md` — Additional audit findings and gaps.
