# Research Report: Transactional Outbox Pattern

Research date: 2026-09-26

## Research Question

How to ensure database changes and event/queue publishing do not leave system half-succeeded? Does Transactional Outbox solve dual-write? How to implement, monitor, avoid pitfalls?

## Executive Summary

Dual-write (DB + broker) cannot be atomic. No shared transaction. Outbox solves by writing business row + outbox row in one DB transaction, then separate relay publishes to broker. Guarantee: at-least-once, ordered per service. Consumer must be idempotent. Two relay types: polling publisher, transaction log tailing (CDC). Debezium provides production implementation. Key ops: small payloads, cleanup, DLQ, monitor oldest unprocessed age.

## Findings

### Finding 1: Dual-write problem is real, 2PC not viable

Claim: Service updating DB + sending broker message cannot do both atomically. DB commits but publish fails, or publish succeeds but DB rolls back.

Evidence: Microservices.io states 2PC spanning DB + broker often unsupported / undesirable; send mid-tx unreliable, send after commit risks crash before send. Debezium 2019 confirms Kafka cannot enlist in XA, network failure causes DB✔/Kafka✘ or reverse.

Sources:
- https://microservices.io/patterns/data/transactional-outbox.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Confidence: HIGH

### Finding 2: DB transaction does not cover external systems

Claim: `DB::transaction(()=>{create + Queue::push})` still leaves half-success. DB rollback cannot rollback Redis/Kafka/RabbitMQ/HTTP.

Evidence: Same as Finding 1. Transaction boundary = DB only.

Sources:
- https://microservices.io/patterns/data/transactional-outbox.html

Confidence: HIGH

### Finding 3: Outbox core mechanism

Claim: INSERT business + INSERT outbox in one DB tx. Separate message relay publishes to broker, then marks processed. Result: both commit or both rollback.

Evidence: Microservices.io solution statement verbatim. Participants: Sender, Database, Message outbox (table for RDBMS, property for NoSQL), Message relay.

Sources:
- https://microservices.io/patterns/data/transactional-outbox.html

Confidence: HIGH

### Finding 4: Debezium canonical outbox schema

Claim: Columns: id uuid (dedup key, Kafka header), aggregatetype (topic routing), aggregateid (Kafka key / partition order), type (event name), payload jsonb. This is Debezium's canonical schema, not a universal industry standard — other implementations may use different column names.

Evidence: Debezium blog + Debezium Outbox Event Router docs define same columns, defaults: `table.field.event.id=id`, `route.by.field=aggregatetype`, `table.field.event.key=aggregateid`, `table.field.event.payload=payload`.

Sources:
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html

Confidence: HIGH

### Finding 5: Two relay implementations

Claim: (a) Polling publisher: poll outbox table. Works any SQL DB, tricky ordering. (b) Transaction log tailing: tail MySQL binlog / Postgres WAL / DynamoDB Streams. Accurate, near-realtime, low overhead, DB-specific, still tricky duplicates.

Evidence: Microservices.io defines both as alternatives. Debezium 2019 demonstrates CDC variant with Postgres WAL + Kafka Connect, sub-second latency.

Sources:
- https://microservices.io/patterns/data/polling-publisher.html
- https://microservices.io/patterns/data/transaction-log-tailing.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Confidence: HIGH

### Finding 6: At-least-once, consumer must be idempotent

Claim: Relay can crash after publish before marking processed → duplicate. No magical exactly-once DB→broker. Consumer tracks processed_event_id / ConsumedMessage table, ignores repeats. Side-effects must not blindly re-apply (e.g. +100 points twice).

Evidence: Microservices.io result context. Debezium example `OrderEventHandler.onOrderEvent()` checks `log.alreadyProcessed(eventId)` first, persists `ConsumedMessage(eventId)` in same tx.

Sources:
- https://microservices.io/patterns/data/transactional-outbox.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Confidence: HIGH

### Finding 7: Kafka transactions ≠ outbox exactly-once

Claim: Kafka EOS (2017) gives atomic read-process-write Kafka→Kafka via transaction coordinator + offsets topic + zombie fencing. Does not solve DB→Kafka dual-write.

Evidence: Confluent blog details atomic multi-partition writes, transactional.id fencing. Scope is Kafka topics only.

Sources:
- https://www.confluent.io/blog/transactions-apache-kafka/

Confidence: HIGH

### Finding 8: Payload design, cleanup, event API stability

Claim: Keep payload minimal (IDs + needed fields), not full 500-field object. Outbox table needs retention: archive/delete processed rows, except ephemeral CDC variant using INSERT+DELETE in same tx leaving table empty. Outbox event schema is public API, evolve compatibly.

Evidence: Debezium blog warns event structure is part of emitting service API, consumers lenient. Cleanup claim MEDIUM, supported by lab spec + Debezium ephemeral-row technique.

Sources:
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- Lab spec (Tier 3, practice guidance)

Confidence: MEDIUM for sizing/cleanup thresholds, HIGH for API stability

### Finding 9: Monitoring + retry/DLQ required

Claim: Monitor unprocessed count, oldest unprocessed age, failure rate, retry count, throughput. Example thresholds "2s normal, 47m serious" are lab-spec values for educational purposes, not derived from external production benchmarks. Need bounded retry + dead-letter, housekeeping message log.

Evidence: Debezium blog notes bounded retry + DLQ + housekeeping needed. Specific metrics and threshold values from lab spec, NOT VERIFIED against independent production source.

Sources:
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- Lab spec

Confidence: MEDIUM

### Finding 10: When to use / not use

Claim: Use when pattern is Update DB + Publish Queue/Event/Webhook and loss/duplication matters. Avoid for simple CRUD no cross-system side effect – adds complexity.

Evidence: Microservices.io context: saga participants, domain event publishers. Martin Fowler event-driven taxonomy (notification vs state transfer vs sourcing) gives placement context.

Sources:
- https://microservices.io/patterns/data/transactional-outbox.html
- https://martinfowler.com/articles/201701-event-driven.html

Confidence: HIGH for use case, MEDIUM for "avoid CRUD" (lab guidance, sensible, no contradicting source)

## Areas of Agreement

- Dual-write unsolvable with DB tx alone. All Tier 1 agree.
- Outbox = single DB tx + async relay. Agree.
- At-least-once + idempotent consumer mandatory. Agree.
- Ordering per aggregate via aggregateid as Kafka key. Agree.
- CDC preferred for low latency, polling most portable. Agree.

## Areas of Disagreement

- None material. Nuances only (see 04-contradictions.md): EOS terminology scope, polling vs CDC emphasis, cleanup variant.

## Limitations

- Monitoring thresholds, payload size limits, retry counts from lab spec lack independent Tier 1 corroboration.
- No primary benchmark for outbox write overhead collected.
- Cloud managed variants (DynamoDB Streams, Spanner change streams) not deeply researched.
- Sources opened via fetch; Kafka transactions page truncated, full content archived in tool output file.

## Conclusion

Outbox eliminates half-success by collapsing dual-write into single DB transaction plus reliable async publish. It trades strong cross-system atomicity for eventual consistency with at-least-once delivery. Correctness depends on idempotent consumers, ordered partitioned publishing, small stable payloads, bounded retry + DLQ, cleanup, and oldest-unprocessed monitoring.
