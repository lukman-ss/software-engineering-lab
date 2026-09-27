# Audit Plan: Research Stage

Target Lab: `labs/26-contract-testing`
Audit Scope: Research deliverables only (`labs/26-contract-testing/research/`)
Audit Date: 2026-09-27

## Files Reviewed
- `labs/26-contract-testing/research/01-plan.md`
- `labs/26-contract-testing/research/02-sources.md`
- `labs/26-contract-testing/research/03-evidence.md`
- `labs/26-contract-testing/research/04-contradictions.md`
- `labs/26-contract-testing/research/05-report.md`
- `labs/26-contract-testing/research/06-open-questions.md`

## Claims To Verify
1. Definition of contract testing: isolating integrations and checking messages conform to shared understanding without deploying both applications.
2. Consumer-Driven Contracts (CDC) pattern: consumer expectations drive provider contracts; contracts are code-first and derived from the union of consumer expectations.
3. Two-phase Pact workflow: Consumer mock testing (generating pact JSON) and Provider verification (replaying against real service).
4. Scope boundary: Contract tests verify message schemas and payloads, not provider internal business logic/side-effects.
5. Deployment gate: Pact Matrix and `can-i-deploy` CLI command prevent breaking changes from reaching environments.
6. Multi-protocol support: Contract testing applies to HTTP REST and asynchronous message brokers (Kafka, RabbitMQ, SNS/SQS) by abstracting payload.
7. Additive vs breaking evolution: Field additions are safe under minimal-response matching; type changes/renames/removals require major versioning (e.g., Google AIP-185).
8. Rebalanced test pyramid: Unit tests base -> Contract tests middle layer -> Minimal E2E tests top layer.

## Source Verification Strategy
- Direct live inspection via HTTP/WebFetch for primary sources:
  - Pact Documentation (`https://docs.pact.io/`)
  - Martin Fowler Bliki (`https://martinfowler.com/bliki/ContractTest.html`)
  - Martin Fowler Article (`https://martinfowler.com/articles/consumerDrivenContracts.html`)
  - Google AIP-185 (`https://google.aip.dev/185`)
  - Pactflow blog references and GitHub repositories.
- Validate claims against exact quotes, author credentials, publication context, and tier categorization.

## Primary Risks
- Over-reliance on vendor-published literature (Pact Foundation / SmartBear Pactflow) for empirical claims (e.g. speed/cost advantages).
- Conflation of provider schema validation (OpenAPI/JSON schema) with consumer-driven contract testing.
- Over-generalizing additive backward compatibility without highlighting consumer parser strictness issues.
