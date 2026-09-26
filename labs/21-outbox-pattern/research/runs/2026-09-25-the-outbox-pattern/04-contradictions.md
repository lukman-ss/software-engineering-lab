# 04-contradictions.md

No material contradictions discovered.

All examined sources (Microservices.io, AWS Prescriptive Guidance, and Debezium Blog) consistently agree on:
1. The definition and hazards of the dual-write problem across disparate data stores / brokers without 2PC.
2. The mechanism of writing events into an outbox table in the same local database transaction.
3. The existence and trade-offs of the two main relay patterns (Polling Publisher vs. Transaction Log Tailing/CDC).
4. The guarantee of at-least-once delivery semantics and the subsequent requirement for idempotent consumers.
