# Failure Scenarios

Catalog of failure modes this lab reproduces and mitigates.

## 1. Duplicate Payment (Idempotency Failure)

**Scenario**: Client retries payment request due to timeout/network error

```
Timeline:
T0: Client POST /payments {amount: 100, idempotency_key: "abc"}
T1: Server processes, calls payment gateway
T2: Gateway succeeds, but response lost (timeout)
T3: Client retries with same idempotency_key
T4: Without protection: Second charge created ❌
```

**Root Cause**: No deduplication on retry

**Fix**: Idempotency key + payload hash + unique constraint

---

## 2. Lost Update (Race Condition)

**Scenario**: Two requests decrement inventory simultaneously

```
Timeline:
T0: Request A reads inventory=10
T1: Request B reads inventory=10
T2: Request A writes inventory=9
T3: Request B writes inventory=9  ❌ Should be 8
```

**Root Cause**: Read-modify-write not atomic

**Fix**: Optimistic locking (version), pessimistic locking (SELECT FOR UPDATE), or atomic decrement

---

## 3. Dirty Read / Non-Repeatable Read / Phantom Read (Transaction Isolation)

**Scenario**: Concurrent transactions see intermediate state

```
Tx1: BEGIN; UPDATE orders SET status='paid' WHERE id=1;
Tx2: BEGIN; SELECT * FROM orders WHERE id=1;  -- sees 'paid' before Tx1 commits ❌
Tx1: COMMIT;
```

**Root Cause**: Default isolation (READ COMMITTED) allows non-repeatable reads

**Fix**: REPEATABLE READ or SERIALIZABLE isolation, or application-level locking

---

## 4. Deadlock

**Scenario**: Circular wait on multiple resources

```
Tx1: Lock order 1 → Lock order 2
Tx2: Lock order 2 → Lock order 1
Result: Both wait forever ❌
```

**Root Cause**: Inconsistent lock ordering

**Fix**: Global lock ordering, lock timeout, deadlock detection + retry

---

## 5. Event Loss (Missing Outbox)

**Scenario**: Order created but event not published

```
Timeline:
T0: API creates order in DB
T1: API publishes event to message broker
T2: Broker unavailable, publish fails
T3: API returns success to client
T4: Event lost forever ❌
```

**Root Cause**: Non-transactional event publishing

**Fix**: Transactional outbox pattern — write event to DB in same transaction

---

## 6. Thundering Herd (Cache Stampede)

**Scenario**: Cache expires, many requests hit DB simultaneously

```
Timeline:
T0: Cache key expires
T1: 1000 requests hit DB for same key ❌
```

**Root Cause**: No cache protection

**Fix**: Single-flight, probabilistic early expiration, distributed lock

---

## 7. Cascading Failure (No Circuit Breaker)

**Scenario**: Downstream service slow, callers exhaust resources

```
Timeline:
T0: Payment gateway slow (10s latency)
T1: API threads block waiting
T2: Thread pool exhausted
T3: All requests fail, including unrelated ones ❌
```

**Root Cause**: No failure isolation

**Fix**: Circuit breaker, timeout, bulkhead

---

## 8. Unbounded Load (No Rate Limiting)

**Scenario**: Single client sends 100k requests/second

```
Timeline:
T0: Malicious/spammy client floods API
T1: DB connection pool exhausted
T2: Legitimate requests fail ❌
```

**Root Cause**: No per-client limits

**Fix**: Token bucket per client, global limit, priority queues

---

## 9. Retry Storm

**Scenario**: Transient error, clients retry immediately

```
Timeline:
T0: Service returns 503 (overloaded)
T1: All clients retry immediately
T2: Load spikes further, service stays down ❌
```

**Root Cause**: No backoff, no jitter

**Fix**: Exponential backoff + jitter, circuit breaker

---

## 10. Split Brain (Distributed Lock Failure)

**Scenario**: Redis master fails, two clients get same lock

```
Timeline:
T0: Redis master fails, failover in progress
T1: Client A gets lock on old master
T2: Client B gets lock on new master
T3: Both think they hold lock ❌
```

**Root Cause**: Lock not replicated synchronously

**Fix**: Redlock algorithm, or accept brief window, use DB advisory locks

---

## Summary Matrix

| Lab | Failure | Detection | Prevention |
|-----|---------|-----------|------------|
| 01 | Duplicate payment | Unique constraint violation | Idempotency key |
| 02 | Data race | `-race` detector, flaky tests | Mutex, atomic, CAS |
| 03 | Inconsistent state | Integration tests | Transaction boundaries |
| 12 | Lost update | Concurrent test | Optimistic locking |
| 05 | Lock contention | Load test | Pessimistic locking |
| 06 | Deadlock | Stress test, pg_locks | Lock ordering, timeout |
| 14 | Event loss | Chaos test (kill broker) | Outbox pattern |
| 08 | Retry amplification | Load test | Backoff + jitter |
| 09 | Cascade failure | Chaos test (slow dependency) | Circuit breaker |
| 10 | Resource exhaustion | Load test | Rate limiting |