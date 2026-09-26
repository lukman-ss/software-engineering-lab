# Research Topic

Contract Testing — API Bisa Sama-Sama "Lulus Test", Tapi Integrasi Tetap Rusak (Contract Testing — APIs Can Both "Pass Tests" But Integration Still Broken)

## Objective

Investigate and document the theory, practice, tools, and patterns of contract testing in distributed systems, with emphasis on:
1. Consumer-Driven Contracts (CDC) as defined by Ian Robinson and Martin Fowler
2. Pact as the primary implementation tool
3. Breaking vs. non-breaking change detection
4. Contract testing for REST APIs and event-driven systems (AsyncAPI, Kafka, etc.)
5. CI/CD integration patterns (expand/contract pattern, can-i-deploy)
6. Common mistakes and anti-patterns

## Research Questions

1. What is the theoretical foundation of Consumer-Driven Contract Testing?
2. How does Pact implement contract testing for HTTP and message-based integrations?
3. What constitutes a "contract" beyond just JSON schema (status codes, error behaviors, semantics)?
4. What are the criteria for breaking vs. non-breaking (additive) changes?
5. How does contract testing compare to E2E and integration testing?
6. How does the expand/contract pattern enable safe breaking changes?
7. What is the role of Pact Broker in CI/CD pipelines?
8. How does contract testing apply to event-driven architectures (Kafka, RabbitMQ, webhooks)?
9. What are the common anti-patterns and misconceptions?
10. What are the limitations of contract testing?

## Search Strategy

- Primary sources: Pact official documentation (docs.pact.io), Martin Fowler's CDC article (2006), Martin Fowler's ContractTest bliki (2011), ParallelChange pattern (2014)
- Secondary: AsyncAPI specification for event contracts, Spring Cloud Contract (archived but relevant)
- Focus on: Official docs, pattern originators, tool maintainers
- Cross-reference: Pact FAQ, CDC theory, expand/contract pattern

## Expected Primary Sources

1. Pact.io Documentation (docs.pact.io) - Tier 1, official tool docs
2. Martin Fowler, "Consumer-Driven Contracts: A Service Evolution Pattern" (2006) - Tier 1, pattern originator
3. Martin Fowler, "Contract Test" bliki (2011) - Tier 1, pattern definition
4. Martin Fowler, "Parallel Change" (2014) - Tier 1, breaking change pattern
5. AsyncAPI Initiative - Tier 1, event-driven API specification
6. Pact FAQ & Best Practices (contract_tests_not_functional_tests) - Tier 1, tool-specific guidance

## Risks / Unknowns

- Spring Cloud Contract is archived (July 2026) - may not reflect current practices
- AsyncAPI documentation is more of a spec reference than practical contract testing guide
- Need to verify claims about event contract testing (Pact Message Pact) from Pact docs
- Limited primary sources on "when NOT to use contract testing" beyond Pact's own FAQ
- Lab exercise specific analysis: determining breaking changes for three specific field changes