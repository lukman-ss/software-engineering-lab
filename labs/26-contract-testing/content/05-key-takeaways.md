# Key Takeaways

1. **Contract testing closes the gap** between passing unit/integration tests and integration failures in production caused by broken service contracts (field rename, enum casing change, primitive type mutation).

2. **Consumer-Driven Contracts** flip responsibility: consumers define minimal expectations, not providers; only fields *actually used* by consumers are verified, freeing providers to evolve safely.

3. **Contracts are semantic, not syntactic**: They include HTTP semantics (method, path, status, headers), field names/types, casing, and error behavior — not just JSON schema.

4. **Breaking change detection happens in CI before deployment**: Provider verification failure blocks deployment with precise error lists; this prevents silent integration failures in production.

5. **Safe API evolution via expand/contract pattern** (three phases: add new field/endpoint, migrate consumers, remove old field/endpoint) allows independent consumer/provider deployment.

6. **Over-specification is an anti-pattern**: Testing validation rules (e.g., "username max 20 chars") in contracts blocks provider evolution; contracts should verify *format* and *presence*, not *why* validation failed.

7. **Contract testing complements — does not replace — other tests**: Unit and functional tests validate internal behavior and business logic; contract tests validate integration boundaries; E2E tests validate critical user journeys.

8. **Subset verification** means provider-side extra fields (e.g., `notes`, `created_at`) do not break consumer contracts — only consumer-required fields are checked. Verifier validates status code, response headers, and body fields (`internal/contract/verifier.go:90-98`).

9. **Type precision matters**: Using `json.Number` in contracts and `UseNumber()` decoder allows verifier to distinguish integer `150000` from string `"150000"`, catching semantic type changes that JSON Schema alone would miss.

10. **The pattern applies to event-driven systems** via Message Pact for Kafka, RabbitMQ, SNS, webhooks — same principles: message format and semantics verified between producer and consumer.
