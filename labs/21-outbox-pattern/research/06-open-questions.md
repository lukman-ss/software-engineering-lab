# Open Questions

## Unanswered Questions
1. What are recommended retry policies and dead-letter queue configurations for the message relay when the message broker is temporarily unavailable?
2. How should the outbox table be partitioned or indexed for high-throughput scenarios (10K+ events/second)?
3. What is the optimal balance between polling interval (for polling publisher) and log-tailing latency requirements?
4. How should schema evolution be handled for outbox events when the payload structure changes over time?

## Weak Evidence / Implementation Details
1. Exact numeric thresholds for monitoring alerts (e.g., "oldest event > 5 seconds = warning") — varies by SLA
2. Specific dead-letter queue implementation patterns for failed message relay attempts
3. Recommended retention periods for archived outbox events (regulatory vs operational requirements)
4. Performance impact of JSON vs Avro payload serialization at scale

## Claims Needing Deeper Research
1. Comparison of exactly-once delivery mechanisms (idempotent consumers + deduplication) vs true exactly-once semantics in modern streaming platforms
2. Performance comparison of outbox pattern vs event sourcing for audit trail requirements
3. Security considerations for outbox tables containing sensitive payload data (encryption at rest, access controls)

## Possible Next Research Directions
1. Implement a reference outbox pattern service in Java/Spring Boot with Debezium and compare performance characteristics
2. Research transactional outbox implementations in NoSQL databases (MongoDB, DynamoDB) beyond the relational focus
3. Investigate how service mesh technologies (Istio, Linkerd) interact with or complement the outbox pattern for service-to-service communication
4. Study the relationship between outbox pattern and modern event streaming platforms (Pulsar, Kinesis) that have built-in transactional capabilities