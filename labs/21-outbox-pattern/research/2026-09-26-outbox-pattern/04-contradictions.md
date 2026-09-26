# Contradictions

No material contradictions discovered between Tier 1 sources on core claims.

## Nuance 1: Delivery guarantee terminology

SOURCE A (Microservices.io):
Outbox gives at-least-once delivery; consumer must be idempotent.

SOURCE B (Confluent/Kafka transactions blog, 2017):
Kafka transactions enable exactly-once processing for read-process-write cycles within Kafka (offsets topic + output topics atomically).

ASSESSMENT:
Different scopes. Kafka EOS applies Kafka-to-Kafka. Outbox spans DB-to-broker, no shared transaction, so guarantee remains at-least-once end-to-end. No contradiction, but terminology confusion risk. Confidence: HIGH.

## Nuance 2: Polling vs CDC ordering / duplicates

SOURCE A (Microservices.io polling-publisher):
Tricky to publish events in order.

SOURCE B (Microservices.io transaction-log-tailing):
Guaranteed accurate, but tricky to avoid duplicate publishing.

ASSESSMENT:
Both approaches acknowledge duplicate + ordering challenges, emphasis differs. CDC preferred for near-realtime/low overhead (Debezium 2019). Choice is tradeoff, not contradiction. Confidence: HIGH.

## Nuance 3: Outbox table cleanup

SOURCE A (Debezium 2019 blog):
Uses persist()+remove() in same tx; CDC captures INSERT from WAL, table stays empty, no housekeeping needed.

SOURCE B (General outbox guidance / lab spec):
Mark processed_at, archive/delete after retention.

ASSESSMENT:
Two valid variants: ephemeral-row CDC vs persistent outbox + relay + cleanup. Depends on relay type. Confidence: MEDIUM.
