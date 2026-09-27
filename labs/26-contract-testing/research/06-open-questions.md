# Open Questions

## Unanswered Questions
- Which Pact implementations (JS/Go/Python/Java) are most mature for the lab exercise targeting Mobile + Backend REST?
- What concrete CI pipeline pattern (GitHub Actions?) should the lab recommend for provider verification + can-i-deploy?
- How to handle contract testing for GraphQL or gRPC — Pact supports HTTP/message, but what about schema-first IDL approaches (Buf, Apollo)?
- What is current state of Spring Cloud Contract successor (Stubborn.sh) after archival Jul 2026?

## Weak Evidence
- Claim additive vs breaking changes: inferred from minimal-response behavior + anti-pattern warnings, not explicit Tier 1 definition of taxonomy. Needs stronger source (e.g., Google AIP-180 Backward Compatibility).
- Effectiveness metrics (e.g., reduction in integration bugs, deployment speed) rely on vendor case studies (Pactflow), not independent measurements.
- Test pyramid rebalancing claim is heuristic (Mike Cohn) + vendor blog, not empirical study.

## Claims Needing Deeper Research
- Bi-directional contract testing (Pactflow) vs consumer-driven: when to choose which? Pact docs mention provider contract testing (OpenAPI) as less effective alone — needs comparison.
- Message pact best practices for Kafka schema registry (Avro/Protobuf) vs Pact message pacts — overlap / complement?
- Versioning strategies for REST breaking changes: URL versioning vs header vs consumer-specific contracts — Google AIP-185 vs Stripe-style date versioning.

## Possible Next Research Directions
- Fetch Google AIP-180 (Backward Compatibility) for explicit breaking-change taxonomy.
- Inspect Pact Broker / Pactflow docs for concrete CI setup checklist (pact_nirvana guide).
- Survey alternative tools: Hoverfly, WireMock, Mountebank, Karate, Schemathesis — positioning vs Pact.
- Look at real-world case studies (M1 Finance, Sngular) for quantitative evidence.
- Research contract testing for event-driven architectures with schema registries.
