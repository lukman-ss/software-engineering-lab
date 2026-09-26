# Contradictions

No material contradictions discovered between Tier 1 authoritative sources.

## Minor Differences in Emphasis

1. **CDC integration**: microservices.io catalog treats Transaction Log Tailing and Polling Publisher as separate sibling patterns. Debezium treats CDC as the primary outbox integration. These are complementary — Debezium's outbox SMT is a CDC-based implementation of the Transaction Log Tailing approach. No contradiction in substance, only in framing (pattern catalog vs. implementation guide).

2. **Duplicate handling location**: microservices.io says "consumer must be idempotent." Debezium blog says idempotency is needed both at consumer (eventId tracking) AND optionally at outbox table level (CDC-based approach inserts+removes rows so table stays empty). These reflect different implementation choices, not contradictions.

## Areas of Consensus

- Outbox guarantees database atomicity; message delivery is at-least-once
- Two relay implementations: polling and log-tailing
- Consumers must be idempotent
- 2PC is explicitly the anti-pattern this addresses
