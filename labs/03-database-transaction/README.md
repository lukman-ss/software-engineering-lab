# Lab 03 — Distributed Transaction: Kenapa Database Transaction Saja Tidak Cukup?

> **Mental Model**: Database transactions (ACID) menjamin atomicity **hanya pada transactional resource yang berpartisipasi dalam transaction tersebut**. External HTTP APIs, email, WhatsApp, object storage, dan kebanyakan external systems tidak ikut berpartisipasi dalam local transaction database aplikasi.

---

## 1. Local Atomicity

### Flow yang Benar
```
[BEGIN TRANSACTION]
     │
     ├── INSERT payment
     ├── UPDATE orders SET status = 'paid'
     └── INSERT wallet_transactions
     │
     └── [COMMIT]
```

Jika ada error → **ROLLBACK** → Semua pembatalan bersama.

> Use DB transaction to protect **local invariants** (resources yang harus konsisten secara atomik dalam database yang sama).

---

## 2. Where DB Transaction Stops

A normal database transaction guarantees atomicity only for resources participating in that transaction. External HTTP APIs, email, WhatsApp, object storage, dan kebanyakan external systems **tidak ikut berpartisipasi**.

*(XA/2PC dapat koordinasi distributed resource, tapi bukan fokus lab.)*

---

## 3. External Side Effects

```
[BEGIN] → INSERT payment → HTTP WhatsApp → ERP sync → [COMMIT]
```

**Risiko**: Transaction bertahan lama → Connection pool tertahan → Throughput menurun.

**External call di dalam transaksi**:
- BEGIN → UPDATE → HTTP → COMMIT (transaction terblokir latency)
- BEGIN → UPDATE → COMMIT → HTTP (transaction cepat, external tidak terpengaruh)

> **Catatan senior**: HTTP timeout tidak selalu berarti external operation gagal. Remote service bisa berhasil tapi response hilang akibat jaringan. Karena itu external commands seperti payment/refund/reservation juga perlu idempotency key + status reconciliation. Lihat Daily #1 Idempotency untuk detail.

---

## 4. HTTP Inside DB Transaction

```
[BEGIN TRANSACTION]
     ├── UPDATE invoice SET status = 'paid'
     ├── HTTP call ke payment gateway (50ms–5s)
     └── [COMMIT]
```

**Perbandingan:**

| Flow | Transaction Lifetime | Konsekuensi |
|------|---------------------|-------------|
| BEGIN → UPDATE → HTTP → COMMIT | Selama HTTP call | Connection pool tertahan |
| BEGIN → UPDATE → COMMIT → HTTP | Durasi DB saja | Tidak terpengaruh latency |

Test bukti: `TestHTTPInsideTransactionDuration`, `TestCommitThenExternalCall`.

---

## 5. Dual-Write Problem

```
[DB commit]   →   success
                       │
                  X process crash
                       │
              [Publish event]  → TIDAK pernah terjadi
```

**Hasil**: Invoice = Paid, **event tidak pernah sampai ke consumers**.

### Reverse Dual-Write juga bermasalah

```
[Publish event]  →  success
                         │
                    X process crash
                         │
                    [DB commit]  → TIDAK pernah terjadi
```

Mengubah urutan **tidak menyelesaikan** atomicity — hanya memindahkan window kegagalan.

---

## 6. Transactional Outbox

```
[BEGIN TRANSACTION]
     ├── UPDATE invoice SET status = 'paid'
     └── INSERT INTO outbox_events (...)
     │
     └── [COMMIT]
```

Transactional Outbox **atomically records the business change and the intent to publish an event.**

**Garanti:**
- Dispatcher dapat retry pada failure transien
- Duplicate delivery mungkin terjadi (at-least-once)
- Consumer **harus idempotent**

---

## 7. Dispatcher Concurrency — At-Least-Once Overlap

```
Dispatcher A: SELECT pending → [evt-1]
Dispatcher B: SELECT pending → [evt-1]
Dispatcher A: Publish evt-1 ✓
Dispatcher B: Publish evt-1 ✓   ← DUPLICATE
```

**Production alternatives:**

| Strategi | Keterangan |
|----------|-----------|
| `SELECT ... FOR UPDATE SKIP LOCKED` | Pessimistic lock per row (PostgreSQL) |
| Claim column (`claimed_by`, `claimed_at`) | Optimistic claim sebelum publish |
| Lease expiry | Dispatcher crash lepas klaim setelah timeout |
| Partitioning | Setiap dispatcher handle range event_id |
| Single dispatcher | Scale via queue bukan fanout |

---

## 8. At-Least-Once Delivery

**Failure scenario:**
```
worker receives event → business op succeeds → crash before ACK → message redelivered
```

Exactly-once end-to-end **sulit**. Sistem yang praktis gunakan:

```
at-least-once delivery
+
idempotent processing
```

---

## 9. Idempotent Consumer

### Why Per-Consumer Deduplication?

```
InvoicePaid evt-123
        │
        ├── InventoryConsumer
        ├── CommissionConsumer
        └── ERPConsumer
```

Ketiganya harus memproses event yang sama. Jadi `event_id` sendiri tidak cukup sebagai global processed key.

**Schema processed_events:**
```
processed_events
---------------------------------
consumer_name
event_id
processed_at
UNIQUE (consumer_name, event_id)
```

### Consumer Transaction Pattern

```
Message Delivered
      │
      ▼
BEGIN
      │
      ├── Claim (consumer_name, event_id)
      │
      ├── Business DB Mutation
      │
      └── Commit
      │
      ▼
ACK
```

Failure sebelum commit:
- ROLLBACK
- → message may retry

Crash setelah commit sebelum ACK:
- message redelivered
- → dedup prevents duplicate business effect

### At-Least-Once Delivery Semantics

Gunakan pola:
```
at-least-once delivery
+
idempotent consumer
```

**Tidak** mengubah menjadi "exactly-once". Lebih tepat: "effectively-once business effect" atau "idempotent processing under redelivery". Delivery tetap dapat duplicate, tapi business effect akhirnya hanya sekali.

**Test mappings:**
- `TestIdempotentConsumerDeduplication` - same consumer duplicate
- `TestDifferentConsumersSameEvent` - different consumers same event  
- `TestConcurrentDuplicateConsumer` - concurrent same event dedup
- `TestConcurrentDifferentConsumersSameEvent` - different consumers, same event, concurrent
- `TestConcurrentDifferentEvents` - concurrent different events, no lost update
- `TestConcurrentSameEvent` - concurrent same consumer same event
- `TestSequentialRedeliveryIdempotency` - sequential retry (idempotent skip, not crash)
- `TestAtomicConsumerFlow` - dedup + business state atomic commitment
- `TestConsumerCrashAfterCommitBeforeAck` - separation of deliveries vs business rows (consumer restart with persistence-backed dedup)
- `TestMockDBNoLostUpdates` - verify concurrency does not lose updates
- `TestMockDBRollbackIsolation` - verify rollback isolation
- `TestBusinessMutationFailureRollback` - business mutation fails → full rollback
- `TestRedeliveryAfterBusinessFailure` - redelivery succeeds after rollback clears state
- `TestProcessedMarkerFailureRollback` - claim failure, no business mutation
- `TestConsumerAtomicCommitHappyPath` - both claim + business state committed atomically
- `TestMockDBOnConflictSemantics` - ON CONFLICT DO NOTHING returns correct RowsAffected
- `TestConsumerRestartRedelivery` - consumer restart, redelivered event deduplicated
- `TestEventualConsistencyCorrectness` - outbox eventual consistency flow
- `TestREADMESyncIdempotentConsumer` - README idempotent consumer pattern verified

---

## 10. Retry Policy

| Error Type | Action |
|------------|--------|
| Transient (timeout, 503, connection reset) | Retry with backoff |
| Permanent (invalid payload, validation, unsupported event) | No retry, move to DLQ |

---

## 11. DLQ Semantics

DLQ **bukan tempat sampah**. Menyimpan pesan yang melewati retry policy untuk investigation/reprocessing.

**Minimal DLQ record:**
- `event_id`
- `reason`
- `attempts`
- `failed_at`

Production memerlukan: monitoring, alerting, replay tooling.

---

## 12. Saga Business Semantics

```
Create payment → Cash out → Generate journal → Sync ERP fails
```

**Jika payment sudah final**, ERP hanya downstream integration:
- **JANGAN compensate** payment atau cash out
- **Retry ERP** (event-driven integration)

Saga compensation untuk **semantic undo**, bukan technical rollback.

**Saga hanya menjalankan compensation untuk step yang sebelumnya berhasil diselesaikan. Step yang gagal sebelum dianggap complete tidak memiliki completed effect yang perlu dikompensasi, kecuali API step tersebut memiliki business semantics khusus.**

---

## 13. Saga Compensation vs DB Rollback

|  | Rollback | Compensation |
|--|----------|----------------|
| Scope | Transaction scope | Business scope |
| Effect | Undo uncommitted | New action yang reverses |
| Example | `ROLLBACK` | `RefundPayment` (baru, bukan delete) |

Compensation **membuat record baru** (Contoh: RefundPayment), tidak menghapus data historis.

### Compensation Failure Handling

**Compensation itself can fail.** Jika compensation gagal:

```
A → B → C fails
Compensate C → skipped (C never completed)
Compensate B → FAILS (e.g. refund gateway timeout)
Compensate A → STILL ATTEMPTED (don't leave A un-compensated)
```

**Implementation principle:**
- Jangan swallow compensation error
- Continue compensating remaining steps
- Return `SagaExecutionError` that includes both original error and compensation errors

**Test:** `TestSagaCompensationFailureHandling` — Step B compensation fails, Step A compensation still executes.

---

## 14. Choreography vs Orchestration

### Choreography (Event-Driven Flow)
```
InvoicePaid ↓ Inventory listens ↓ InventoryReserved ↓ Accounting listens
```

Keuntungan: loose coupling, scaling per consumer, tidak ada SPOF.  
Kekurangan: flow visibility terhalang, debugging tersebar.

### Orchestration (Explicit Workflow)
```
PaymentSaga
├─ ReserveInventory
├─ CreateJournal
├─ NotifyERP
└─ ...
```

Keuntungan: explicit workflow, easy to modify, compensation eksplisit.  
Kekurangan: orchestrator complexity, tight coupling, SPOF (bisa direplikasi).

**Jangan mengatakan salah satu selalu lebih baik. Pilih berdasarkan domain complexity.**

---

## 15. Eventual Consistency

Setelah DB commit, downstream belum tentu up-to-date. Ini **normal**, bukan bug.

### Kapan Acceptable?
✅ Notification (WhatsApp, Email)  
✅ Analytics / Reporting  
✅ Audit logs  
✅ Search index updates

### Kapan Butuh Strong Consistency?
❌ Payment + ledger + OPL marking → harus satu DB transaction  
❌ Inventory stok fisik  
❌ Balance / akuntansi

---

## 16. Event Naming: Fact vs Command

| Tipe | Naming | Contoh |
|------|--------|--------|
| **Event** | Past tense | `InvoicePaid`, `PaymentRecorded`, `VendorPaymentCompleted` |
| **Command** | Imperative | `PayInvoice`, `SendNotification`, `SyncToERP` |

---

## 17. Event-Driven Trade-offs

### Keuntungan (+)
- **Resilience**: Proses dapat melanjutkan setelah crash
- **Loose Coupling**: Consumers tidak harus tahu satu sama lain
- **Retryable**: Gagal bisa dicoba ulang
- **Scalable Independently**: Consumer bisa scale terpisah

### Risiko / Kekurangan (-)
- **Eventual Consistency**: Data tidak konsisten secara serentak
- **Operational Complexity**: Monitoring harus mencakup multiple services
- **Duplicate Messages**: At-least-once delivery berarti ada duplikasi
- **Harder Debugging**: Flow menyebar di banyak service
- **Ordering Problems**: Event mungkin sampai out-of-order
- **Observability Requirement**: Diperlukan distributed tracing

Jangan menggambarkan event-driven architecture sebagai solusi universal. Pilih bila membutuhkan retry, compensation, atau loose coupling.

---

## 17. MockDB Concurrency Model (Educational)

MockDB implementasi transaction model yang **serialized**:

- `BEGIN` → acquire global transaction mutex
- `INSERT/UPDATE/DELETE` → operate on transaction snapshot
- `COMMIT` → merge snapshot to committed, release mutex
- `ROLLBACK` → discard snapshot, release mutex

### Why Serialized?

Untuk educational purposes, mockdb menggunakan **single write transaction at a time** untuk:
1. Simplicity (tanpa MVCC/row-level locks)
2. Deterministic test results
3. Atomic ON CONFLICT DO NOTHING semantics

### Concurrency Simulation

Goroutines are concurrent at the application level, while MockDB intentionally serializes write transactions to provide deterministic educational transaction semantics. This verifies application-level logical race paths and atomic dedup behavior under the MockDB model.

Ini secara otomatis memungkinkan:
- Concurrent goroutine tests exercise concurrent application paths.
- Go's `-race` flag detects shared-memory data races.
- MockDB itself serializes database write transactions.
- ON CONFLICT behavior verification
- Rollback isolation verification

---

## 18. Test Organization

| Category | Tests | Purpose |
|----------|-------|---------|
| **Local Transaction** | TestUnsafeLocalTransaction, TestSafeLocalTransaction | ACID verification |
| **External Side Effects** | TestDistributedTransactionExternalSideEffectLimitation | local DB transaction boundary / external side-effect rollback limitation |
| **Dual-Write Problems** | TestDualWriteProblemEventLost, TestReverseDualWriteFailure | Atomicity window analysis |
| **Transactional Outbox** | TestTransactionalOutboxPatternAtomicity, TestOutboxDispatcherPublishesPending | Outbox pattern |
| **Idempotent Consumer** | TestIdempotentConsumerDeduplication, TestAtomicConsumerFlow | Deduplication key design |
| **Concurrent Tests** | TestConcurrentDuplicateConsumer, TestConcurrentDifferentEvents, TestConcurrentSameEvent, TestConcurrentDifferentConsumersSameEvent | Race condition verification |
| **Failure Injection** | TestBusinessMutationFailureRollback, TestProcessedMarkerFailureRollback | Error path rollback |
| **Redelivery/Reprocibility** | TestConsumerRestartRedelivery, TestRedeliveryAfterBusinessFailure | At-least-once correctness |
| **Eventual Consistency** | TestEventualConsistencyCorrectness, TestEventualConsistencyDemo | Outbox flow |

---

## 18. Ordering Problem

Queue/event processing juga memiliki ordering problem.

```
InvoicePaid
InvoiceCancelled
```

Consumer **tidak selalu** dapat mengasumsikan arrival order tanpa mechanism tertentu.

**Production strategies:**
- **Partition key**: Semua event untuk satu aggregate ke partition yang sama
- **Aggregate sequence/version**: Consumer validasi urutan secara eksplisit
- **Consumer validation**: Gagal atau abaikan event out-of-order

---

## 19. Event Versioning

Event contract dapat berubah.

```json
{
  "event_type": "InvoicePaid",
  "event_version": 1,
  "invoice_id": 123,
  "occurred_at": "2024-01-01T00:00:00Z"
}
```

**Consumers harus:**
- Check `event_version` dan handle/ignore versi tidak dikenal
- **Backward compatibility**: new consumer can still read events produced by older producers.
- **Forward compatibility**: older consumer can tolerate events produced by newer producers, usually by ignoring unknown optional fields.

Tidak perlu implement schema registry di lab ini.

---

## 20. Poison Message

**Permanent failure** biasanya menghasilkan poison message:

- Payload malformed (JSON invalid)
- Unsupported event version
- Missing required field (`event_id` kosong)

Ini **tidak membaik** hanya karena retry 100 kali.

**Aksi yang benar:**

1. Event masuk DLQ
2. Investigation oleh operator
3. Fix di producer atau consumer
4. Replay secara manual jika perlu

---

## 21. Replay Safety

Event dari DLQ/outbox dapat di-replay, sehingga **consumer tetap harus idempotent**.

```
[DLQ] → [Replay] → [Consumer] → [Idempotent check] → SUCCESS
```

Jika tidak idempotent → double effect (double refund, double deduction, dll).

**DLQ replay di lab:**
```go
dlq.Replay(eventID) // mengambil event dari DLQ dan mengembalikannya ke queue
```

---

## 22. Local vs External Inventory Clarification

### Local inventory (di dalam DB transaction yang sama)

Jika inventory berada di **local inventory table** yang sama dalam transaksi:

```
BEGIN → UPDATE invoice → UPDATE inventory → COMMIT
```

Jika ada error → **ROLLBACK** → pengurangan inventory ikut dibatalkan. State kembali seperti sebelum transaction.

### External inventory (service boundary berbeda)

Jika inventory merupakan **inventory microservice** atau API terpisah:

```
[BEGIN] → UPDATE invoice → [HTTP Inventory Service] → [COMMIT]
```

Jika transaction rollback → **inventory microservice tidak otomatis dikembalikan** (berada di luar DB transaction). Diperlukan **compensation** (cancel reservation).

---

## 23. Queue Publish Failure

```
[OUTBOX]
     ├── INSERT event_123 (pending)
     └── COMMIT (business transaction committed ✓)

[PUBLISHER DOWN]
     └── HTTP 503 — cannot reach broker

[DISPATCHER]
     └── Publish fails → event remains PENDING
```

**Expected behavior:**
- Business transaction **tetap committed**
- Outbox event **tetap pending**
- Dispatcher **bisa retry** ketika broker recovery
- **Jangan rollback** business transaction hanya karena broker down

---

## 24. Outbox Pending Recovery

```
publisher unavailable → event pending → process restart → publisher recovers → dispatcher retries → success
```

**Test: `TestOutboxPendingRecovery`**
- Broker down (failUpTo=100)
- Dispatcher tidak bisa publish
- Event tetap di outbox
- Setelah broker up → event berhasil deliver

Menunjukkan **durable intent**: event disimpan sebelum dipublish, tidak kehilangan.

---

## Cara Menjalankan

```bash
cd labs/03-database-transaction
go test -v -count=1
```

### Test Coverage Mapping

| Test | Konsep | File |
|------|--------|------|
| `TestUnsafeLocalTransaction` | Partial state corruption | local_transaction.go |
| `TestSafeLocalTransaction` | ACID rollback | local_transaction.go |
| `TestExternalSideEffectRollback` | External effect survives DB rollback | external.go |
| `TestHTTPInsideTransactionDuration` | HTTP latency extends transaction | external.go |
| `TestCommitThenExternalCall` | Commit sebelum external | external.go |
| `TestTransactionStaysOpenDuringExternalCall` | IsTxOpen verification | external.go |
| `TestTransactionCommitDoesNotBlockConnection` | Tx state transitions closed→open→closed | external.go |
| `TestDualWriteProblemEventLost` | Invoice paid, event never delivered | external.go |
| `TestDualWriteCrashWindow` | Dual-write crash window demonstrated | external.go |
| `TestReverseDualWriteFailure` | Publish-then-commit also problematic | external.go |
| `TestTransactionalOutboxPatternAtomicity` | Business state + outbox atomic | external.go |
| `TestTransactionalOutboxRollback` | Outbox insert failure → full rollback | external.go |
| `TestOutboxHappyPathAssertions` | Outbox happy path complete assertions | external.go |
| `TestOutboxDispatcherPublishesPending` | Dispatcher publish | external.go |
| `TestOutboxDuplicateDeliveryAtLeastOnce` | Crash after publish = duplicate | external.go |
| `TestOutboxDispatcherConcurrencyOverlap` | Concurrent dispatchers overlap | external.go |
| `TestDispatcherPublishedAtSemantics` | published_at and attempts semantics | external.go |
| `TestOutboxPendingRecovery` | Broker down → event pending → recovery | external.go |
| `TestIdempotentConsumerDeduplication` | Idempotent deduplication | external.go |
| `TestAtomicConsumerFlow` | Dedup + business state atomic | external.go |
| `TestSequentialRedeliveryIdempotency` | Sequential retry (idempotent skip) | external.go |
| `TestConsumerCrashAfterCommitBeforeAck` | Separation of deliveries vs business rows | external.go |
| `TestMockDBNoLostUpdates` | Mock DB no lost updates | external.go |
| `TestMockDBRollbackIsolation` | Rollback isolation | external.go |
| `TestConcurrentDifferentEvents` | Concurrent different events, no lost update | external.go |
| `TestConcurrentSameEvent` | Concurrent same event, one succeeds | external.go |
| `TestDifferentConsumersSameEvent` | Different consumers, same event | external.go |
| `TestTransientFailureSuccessAfterRetry` | Retry: fail-fail-succeed | external.go |
| `TestDeadLetterQueue` | Retry policy exhaust → DLQ | external.go |
| `TestSagaPaymentWithCompensatingAction` | Saga with failure | external.go |
| `TestSagaCompensationOrderFourSteps` | Reverse compensation order | external.go |
| `TestSagaCompensationFailureHandling` | Continue compensation on failure | external.go |
| `TestCompensationIdempotency` | Double compensation = single execution | external.go |
| `TestInvoicePaidPayloadRoundTrip` | JSON serialize/deserialize | external.go |
| `TestEventualConsistencyDemo` | invoice=paid, worker later | external.go |
| `TestDistributedTransactionExternalSideEffectLimitation` | External side effect survives rollback | external.go |
| `TestConcurrentDifferentConsumersSameEvent` | Different consumers, same event concurrently | external.go |
| `TestBusinessMutationFailureRollback` | Business mutation fails → full rollback | external.go |
| `TestRedeliveryAfterBusinessFailure` | Redelivery after rollback succeeds | external.go |
| `TestProcessedMarkerFailureRollback` | Claim insert fails → no business mutation | external.go |
| `TestConsumerAtomicCommitHappyPath` | Both claim+business commit atomically | external.go |
| `TestMockDBOnConflictSemantics` | ON CONFLICT DO NOTHING RowsAffected=1,0 | mockdb.go |
| `TestMockDBDirectConcurrency` | Two TX A+B insert concurrently both persist | mockdb.go |
| `TestMockDBDuplicateClaimDirect` | Two TX duplicate claim only one succeeds | mockdb.go |
| `TestMockDBDifferentUniqueKeys` | Same consumer, different events both succeed | mockdb.go |
| `TestMockDBDifferentConsumersSameEvent` | Different consumers, same event both succeed | mockdb.go |
| `TestConsumerRestartRedelivery` | Worker restart + redelivery | external.go |
| `TestEventualConsistencyCorrectness` | Outbox dispatch ≠ consumer processing | external.go |
| `TestREADMESyncIdempotentConsumer` | README pattern verified | external.go |

---

## Failure Matrix

| Failure Point | Expected Behavior |
|---------------|-------------------|
| Business mutation fails | Full transaction rollback (processed_events=0, commissions=0) |
| Dedup marker fails | No business mutation committed (processed_events=0, commissions=0) |
| Consumer crashes before commit | Redelivery can retry (processed_events=0, commissions=0) |
| Consumer crashes after commit before ACK | Duplicate delivery; idempotent skip (processed_events=1, commissions=1) |
| Two consumers same event | Both may process independently (different business effects) |
| Same consumer concurrent duplicate | Only one business mutation (processed_events=1, commissions=1) |
| Publisher/transient failure | Retry with backoff |
| Retry exhausted (DLQ) | Moved to Dead Letter Queue |
| Crash after publish to broker | Event may be published twice (at-least-once) |
| Saga compensation fails | Continue compensating remaining steps, return error with details |

---

## Rule of Thumb

### Transaction Design Rules

| Rule | When to Use |
|------|-------------|
| Same DB + same business invariant | Local DB transaction |
| Cross-system side effect | Asynchronous event when appropriate |
| DB update + publish event | Transactional Outbox |
| Message redelivery | Idempotent Consumer |
| Transient failure | Retry + exponential backoff |
| Repeated permanent failure | DLQ |
| Cross-service semantic undo | Saga / Compensation |

**Remember**: These are guidelines, not absolute rules. Choose based on business requirements.

---

## Navigasi

- **Previous**: [Lab 02 — Database Index](../02-database-index/)
- **Next**: [Lab 04 — Caching](../04-caching/)