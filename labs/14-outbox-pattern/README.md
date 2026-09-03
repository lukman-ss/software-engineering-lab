# Lab 14: Outbox Pattern - Dual Write Problem

This lab demonstrates the dual write problem and its solution using the transactional outbox pattern.

## The Dual Write Problem (Prompt 040)

### Scenario
When creating an order, we need to:
1. Create the order in PostgreSQL
2. Publish an "OrderCreated" event to a message queue

Both operations must succeed or both must fail. If only one succeeds, we have inconsistent state.

### Unsafe Implementation

The `unsafe/` package demonstrates the bug:

```
Order Created -> DB Insert (SUCCESS)
             -> Event Publish (FAIL - connection lost)
                    ↓
           Order exists in DB but no event sent!
           CONSUMERS ARE NEVER NOTIFIED!
```

**Files:**
- `unsafe/unsafe.go` - Unsafe service with dual write bug
- `unsafe/order/repository.go` - Unsafe repository

### Root Cause
The operations execute in separate transactions:
1. Order insert commits to DB
2. Event publish to queue fails
3. Order is now "orphaned" - exists in DB but consumers don't know

## Transactional Outbox (Prompt 041)

### Solution
Instead of publishing to the queue, we insert the event into a database table:

```
tx Begin
  ├── Insert Order
  ├── Insert Outbox Event
  └── tx Commit
```

Both operations are in the same transaction - atomic!

### Migration

`migrations/001_outbox_pattern.sql`:
- `orders` table - stores order data
- `outbox_events` table - stores events to be published

**Columns:**
- `id` - unique event ID
- `aggregate_type` - "Order"
- `aggregate_id` - order ID
- `event_type` - "OrderCreated"
- `payload` - JSON serialized event
- `created_at` - when event was created
- `published_at` - NULL until published (nullable)
- `attempts` - number of publish attempts
- `next_attempt_at` - when to retry if failed

### Safe Implementation

`safe/transactional_service.go`:
- `CreateOrder()` uses a single transaction for both inserts
- Either both succeed or both rollback

## Navigasi

- **Previous**: [Lab 13 — Deadlock](../13-deadlock/)

---

## Key Differences

| Aspect | Unsafe | Safe |
|--------|--------|------|
| Consistency | Can be lost | Guaranteed |
| Rollback | None possible | Automatic via transaction |
| Recovery | Manual DB cleanup | Reprocessable via outbox |
| Code path | Sequential, 2 transactions | Single transaction |

## The Outbox Worker (Prompt 042 & 043)

The worker's job is to read unpublished events and send them to the broker.

### Features
- **Batching**: Reads N events at a time to reduce DB roundtrips
- **Retry & Backoff**: Exponential backoff (`base * 2^(attempts-1)`) for failures
- **Graceful Shutdown**: Uses `context.Context` and `sync.WaitGroup` to finish in-flight work
- **Concurrency Control**: Uses PostgreSQL's `FOR UPDATE SKIP LOCKED`

### Concurrent Workers (Prompt 043)
Multiple worker instances can poll the same table without locking each other out.
The query `SELECT ... FOR UPDATE SKIP LOCKED` does two things:
1. Locks the rows it returns so other workers can't process them
2. If it sees rows already locked by another worker, it *skips* them instead of waiting

This allows high-throughput parallel event processing.

## At-Least-Once Delivery Reality (Prompt 044)

The outbox pattern provides **at-least-once** delivery guarantees. This means duplicate messages are a reality, not a bug.

Duplicates happen when:
1. The worker publishes to the broker, but crashes before updating `published_at` in the DB
2. Network timeouts occur between the worker and the broker
3. Broker rebalances or network partitions happen on the consumer side

### The Idempotent Consumer

Because consumers *will* receive duplicates, they must be idempotent.
`safe/consumer.go` demonstrates this using a deduplication table:

```sql
INSERT INTO consumer_processed_events (event_id, processed_at) 
VALUES ($1, NOW())
ON CONFLICT (event_id) DO NOTHING;
```

If the insert succeeds, process the event. If it does nothing (conflict), we've seen this event before and safely ignore it.

## Testing

Run tests against a local PostgreSQL database:

```bash
docker-compose up -d postgres
go test ./labs/14-outbox-pattern/... -v
```