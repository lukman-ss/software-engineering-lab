# Open Questions

## Unanswered Questions

1. **Contract testing for GraphQL**: Pact FAQ confirms GraphQL contract testing is possible but provides minimal practical detail. How do GraphQL schema changes differ from REST contract changes in practice?

2. **Schema registry integration**: Should Pact contracts be cross-validated with AsyncAPI or OpenAPI schemas? The Pact FAQ states Pact does not use JSON Schema — but should schemas still exist as documentation alongside pacts?

3. **Performance of contract tests at scale**: With many consumer-provider pairs, what is the verification overhead on provider CI? How do teams manage provider states across dozens of consumers?

4. **Provider state management for complex domains**: Setting up provider states (e.g., "user 123 exists") requires inserting data into the provider DB without the real API. How do teams maintain this reliably in production-like environments?

5. **Contract testing for legacy systems**: How does one introduce contract testing when neither consumer nor provider currently has tests, and the API contract is undocumented?

6. **Message contract evolution in Kafka**: How does one handle schema evolution for Kafka topics (e.g., Avro/Schema Registry compatibility) alongside Pact Message Pact verification?

7. **Consumer-driven contracts in open/partner APIs**: Pact explicitly states it's not ideal for public APIs. What alternatives exist for third-party API consumers?

8. **Spring Cloud Contract sunset migration**: With Spring Cloud Contract archived in July 2026, what migration path exists for existing Spring users? Does Stubborn.sh fully support CDC?

9. **Quantitative impact metrics**: What is the measurable reduction in production incidents attributable to contract test adoption? DORA metrics correlation is mentioned but not quantified in primary sources.

10. **Contract testing in event sourcing/CQRS**: When events are the primary contract (not API responses), how do consumer-driven contract tests handle event schema evolution across aggregated event streams?

## Weak Evidence / Claims Needing Deeper Research

- **AsyncAPI as contract testing tool**: Only visited AsyncAPI docs overview; spec details for automated validation not verified. Claims about contract enforcement remain LOW confidence.
- **Spring Cloud Contract current state**: Repository archived; Stubborn.sh details beyond the blog post link NOT VERIFIED.
- **Pact vs OpenAPI/JSON Schema performance comparisons**: No direct empirical comparison found; Pact FAQ states preference for concrete examples but no benchmark data.

## Next Research Directions

1. Deep-dive into Pact provider state management strategies for stateful APIs.
2. Comparative study of Pact vs Spring Cloud Contract (legacy) vs AsyncAPI in event-driven contexts.
3. Empirical case studies of contract test adoption and measurable impact on deployment frequency/MTTR.
4. Investigation into contract testing for gRPC/protobuf-based microservices.
5. Detailed guide on implementing expand/contract pattern in real microservice migration scenarios.