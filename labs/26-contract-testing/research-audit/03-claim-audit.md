# Claim Audit

Target Lab: labs/26-contract-testing
Scope: Research Stage Claims

---

## Claim 1

Claim: Contract testing checks that inter-application messages conform to a shared understanding in isolation without deploying both applications together.

Location: `research/03-evidence.md` (Evidence 1), `research/05-report.md` (Finding 1)

Evidence Provided: Direct quotes from Pact Foundation Introduction and Martin Fowler Bliki.

Source: Pact Foundation (`https://docs.pact.io/`), Martin Fowler (`https://martinfowler.com/bliki/ContractTest.html`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Well-supported by primary authoritative sources.

---

## Claim 2

Claim: Consumer-Driven Contracts use consumer test execution to generate contracts containing minimal concrete request/response pairs used by consumers, allowing providers to evolve unused fields freely.

Location: `research/03-evidence.md` (Evidence 2, 4), `research/05-report.md` (Finding 2)

Evidence Provided: Pact Introduction docs and Martin Fowler CDC article (2006).

Source: Pact Foundation (`https://docs.pact.io/`), Martin Fowler (`https://martinfowler.com/articles/consumerDrivenContracts.html`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Accurate representation of consumer-driven contract mechanics.

---

## Claim 3

Claim: Pact workflow operates in two isolated phases: consumer tests with mock provider (generating pact file) and provider verification replaying requests against real provider.

Location: `research/03-evidence.md` (Evidence 5), `research/05-report.md` (Finding 3)

Evidence Provided: Pact How Pact Works guide.

Source: Pact Foundation (`https://docs.pact.io/getting_started/how_pact_works`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Matches Pact architecture specifications.

---

## Claim 4

Claim: Contract tests must focus on message structure and generic error responses rather than provider functional behavior or business logic.

Location: `research/03-evidence.md` (Evidence 7), `research/05-report.md` (Finding 4)

Evidence Provided: Pact Contract Tests vs Functional Tests guide.

Source: Pact Foundation (`https://docs.pact.io/consumer/contract_tests_not_functional_tests`)

Source Actually Supports Claim: YES

Classification: FACT / BEST_PRACTICE

Severity: LOW

Notes: Core guidance on avoiding brittle contracts.

---

## Claim 5

Claim: Pact Broker `can-i-deploy` CLI tool queries the Pact Matrix of consumer/provider versions to gate CI/CD deployment pipelines before code hits production.

Location: `research/03-evidence.md` (Evidence 10), `research/05-report.md` (Finding 5)

Evidence Provided: Pact Broker Can I Deploy documentation.

Source: Pact Foundation (`https://docs.pact.io/pact_broker/can_i_deploy`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Accurately describes Pact Broker matrix release strategy.

---

## Claim 6

Claim: Pact supports asynchronous message contract testing (Message Pacts) for Kafka, RabbitMQ, SNS/SQS, and Kinesis by abstracting transport protocols.

Location: `research/03-evidence.md` (Evidence 6), `research/05-report.md` (Finding 6)

Evidence Provided: Pact How Pact Works non-HTTP testing section.

Source: Pact Foundation (`https://docs.pact.io/getting_started/how_pact_works`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Fully supported.

---

## Claim 7

Claim: Additive API changes (adding optional fields) are safe under minimal expected response matching, while breaking changes require major versioning.

Location: `research/03-evidence.md` (Evidence 11), `research/05-report.md` (Finding 7)

Evidence Provided: Pact docs minimal response matching rules + Google AIP-185 API Versioning.

Source: Pact Foundation (`https://docs.pact.io/`), Google AIP-185 (`https://google.aip.dev/185`)

Source Actually Supports Claim: YES

Classification: FACT / IMPLEMENTATION-SPECIFIC

Severity: LOW

Notes: Research correctly notes that "additive vs breaking" terminology is synthesized from matching logic and versioning guidelines.

---

## Claim 8

Claim: Contract testing rebalances the test pyramid by sitting between unit and E2E tests, reducing reliance on slow and flaky integrated test environments.

Location: `research/03-evidence.md` (Evidence 8, 9), `research/05-report.md` (Finding 8)

Evidence Provided: Pactflow Contract Testing vs Integration Testing blog post.

Source: Pactflow (`https://pactflow.io/blog/contract-testing-vs-integration-testing/`)

Source Actually Supports Claim: YES

Classification: INTERPRETATION / HEURISTIC

Severity: LOW

Notes: Properly identified as Tier 2 heuristic / vendor model in research notes.
