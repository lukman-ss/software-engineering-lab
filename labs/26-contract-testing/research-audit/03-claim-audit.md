# Claim Audit

## Claim 1

Claim:
Contract testing is a technique for testing integration points by checking each application in isolation to ensure messages conform to a shared contract, avoiding brittle and expensive end-to-end integration environments.

Location:
`05-report.md` (Finding 1) & `03-evidence.md` (Evidence 1)

Evidence Provided:
Direct quotes from Fowler (2011) and Pact Documentation Introduction.

Source:
- Pact Foundation Introduction (`https://docs.pact.io/`)
- Martin Fowler Bliki (`https://martinfowler.com/bliki/ContractTest.html`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Definition is uniform across industry and foundational literature.

---

## Claim 2

Claim:
Consumer-Driven Contracts (CDC) derive the provider's obligations from the union of actual consumer expectations, allowing providers to evolve unused functionality freely without breaking consumers.

Location:
`05-report.md` (Finding 2) & `03-evidence.md` (Evidence 2, Evidence 4)

Evidence Provided:
Quotations from Robinson/Fowler (2006) and Pact Foundation Introduction.

Source:
- Ian Robinson / Martin Fowler — Consumer-Driven Contracts (`https://martinfowler.com/articles/consumerDrivenContracts.html`)
- Pact Foundation Introduction (`https://docs.pact.io/`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core architecture pattern accurately captured.

---

## Claim 3

Claim:
Pact workflow runs in two phases: consumer unit test execution (generates pact file via local mock server) and provider verification (replays interactions against real provider service with provider states).

Location:
`05-report.md` (Finding 3) & `03-evidence.md` (Evidence 5)

Evidence Provided:
Step-by-step description and quotations from Pact docs.

Source:
- Pact Foundation — How Pact Works (`https://docs.pact.io/getting_started/how_pact_works`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately details test lifecycle and state setup.

---

## Claim 4

Claim:
Contract tests should focus on message structure/schemas and error response formatting, not provider business logic validation or database side-effects.

Location:
`05-report.md` (Finding 4) & `03-evidence.md` (Evidence 7)

Evidence Provided:
Guidance and examples on username validation rules and testing responsibility table.

Source:
- Pact Foundation — Contract Tests vs Functional Tests (`https://docs.pact.io/consumer/contract_tests_not_functional_tests`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Important distinction to prevent brittle contract test suites.

---

## Claim 5

Claim:
The `can-i-deploy` CLI tool in Pact Broker checks the verification matrix of consumer/provider version pairs and blocks incompatible deployments in CI/CD before release to environments.

Location:
`05-report.md` (Finding 5) & `03-evidence.md` (Evidence 10)

Evidence Provided:
CLI command syntax, exit codes, and matrix explanation.

Source:
- Pact Foundation — Can I Deploy (`https://docs.pact.io/pact_broker/can_i_deploy`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard deployment gate mechanism verified in official docs.

---

## Claim 6

Claim:
Contract testing applies equally to asynchronous messaging architectures (e.g. Kafka, RabbitMQ, SNS/SQS) by abstracting transport protocols and testing message payload schemas.

Location:
`05-report.md` (Finding 6) & `03-evidence.md` (Evidence 6)

Evidence Provided:
Pact non-HTTP testing documentation and message pact model.

Source:
- Pact Foundation — How Pact Works (`https://docs.pact.io/getting_started/how_pact_works`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Supported by Pact v2+ message specification.

---

## Claim 7

Claim:
Additive schema changes (adding new fields) are backward-compatible and safe when consumers practice 'just enough' validation / minimal response matching, whereas field deletions, renames, and type mutations are breaking changes requiring major version increments.

Location:
`05-report.md` (Finding 7) & `03-evidence.md` (Evidence 11, Evidence 14)

Evidence Provided:
Pact minimal-response verification rules, Schematron patterns, and Google AIP-185 major versioning guidelines.

Source:
- Pact Foundation Docs
- Martin Fowler / Ian Robinson (2006)
- Google AIP-185 (`https://google.aip.dev/185`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Robustness principle applied to API contracts is supported by evidence.

---

## Claim 8

Claim:
Contract testing rebalances the test pyramid by replacing most end-to-end integration tests with isolated contract tests, reducing flakiness and execution time.

Location:
`05-report.md` (Finding 8) & `03-evidence.md` (Evidence 9)

Evidence Provided:
Visual pyramid citations from Pactflow vendor blog and industry heuristics.

Source:
- Pactflow Blog (`https://pactflow.io/blog/contract-testing-vs-integration-testing/`)

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
While widely accepted as an industry testing heuristic, quantitative metrics (e.g. specific test execution speedups or exact pyramid ratios) are architectural heuristics rather than universal constants. Properly noted in `06-open-questions.md`.
