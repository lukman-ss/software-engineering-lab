# Open Questions

- **Kafka topic compaction tradeoffs for outbox replay**: When using compacted topics for long retention, full state must be encoded in each event. Research needed on practical sizing/performance impacts.
- **Idempotency strategy overhead in high-throughput systems**: How do idempotency keys / message log tables scale? What is the performance cost of dedup checks at very high event rates?
- **Cloud provider managed implementations**: Which major cloud providers offer managed outbox-like functionality (e.g. AWS DynamoDB Streams, GCP Spanner change streams)?
- **Exactly-once via transactional outbox + Kafka**: How does the interaction between outbox + Kafka transactions achieve end-to-end exactly-once semantics in practice?
- **Schema evolution for outbox payloads**: What are the best practices for evolving the schema of outbox events over time without breaking existing consumers?
- **Monitoring and alerting thresholds**: What are industry-standard thresholds for "oldest unprocessed outbox event age" before an alert is triggered?
- **Performance impact of outbox writes**: What is the measurable latency and throughput overhead of adding outbox table writes in high-throughput transaction systems?