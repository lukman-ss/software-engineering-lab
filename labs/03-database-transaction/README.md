# Lab 03 — Distributed Transaction: Kenapa Database Transaction Saja Tidak Cukup?

> **Mental Model**: Database transactions (ACID) menjamin atomicity **hanya pada transactional resource yang berpartisipasi dalam transaction tersebut**. External HTTP APIs, email, WhatsApp, object storage, dan kebanyakan external systems tidak ikut berpartisipasi dalam local transaction database aplikasi.

---

## 1. Local Database Transaction

A normal database transaction guarantees atomicity for resources participating in that transaction. External HTTP APIs, email, WhatsApp, object storage, and most external systems do not participate. (XA/2PC can coordinate distributed resources, but that is not the focus.)

### Flow yang Benar
```
[Local Transaction]
     │
     ├── Payment
     ├── Cash Ledger
     └── Outbox Event
```

**Aturan utama:**
> Use DB transaction to protect **local invariants** (resources that must be atomically consistent within the same database).  
> Use **asynchronous messaging** for cross-boundary side effects (WhatsApp, ERP, email, analytics).

---

## 2. Transaction Boundary

```
[Unsafe Distributed Flow]
     │
     ├── [BEGIN]
     ├── UPDATE invoice
     ├── [HTTP WhatsApp] ────► External World
     ├── ERP sync
     └── [COMMIT/Rollback]

[Result]
ROLLBACK DB          ← invoice dibatalkan
WhatsApp Tetap Terkirim  ← TIDAK dapat di-rollback!
```

---

## 3. Anti-Pattern: HTTP Call Di Dalam Transaction

```
[BEGIN TRANSACTION]
     │
     ├── UPDATE invoice SET status = 'paid'
     ├── HTTP call ke payment gateway (50ms–5s)
     └── [COMMIT]
```

**Risiko**: Transaction bertahan lama → Connection pool tertahan → Throughput menurun.

| Flow | Transaction Lifetime | Risiko |
|------|---------------------|--------|
| BEGIN → UPDATE → HTTP → COMMIT | Selama HTTP call | Connection pool tertahan |
| BEGIN → UPDATE → COMMIT → HTTP | Durasi DB saja | Tidak terpengaruh latency |

---

## 4. Dual-Write Problem

```
[DB commit]   →   success
                       │
                  X process crash
                       │
              [Publish event]  → TIDAK pernah terjadi
```

**Hasil**: Invoice = Paid, event tidak sampai ke consumers.

### Reverse Dual-Write juga bermasalah

```
[Publish event]  →  success
                         │
                    X process crash
                         │
                    [DB commit]  → TIDAK pernah terjadi
```

**Hasil**: Consumer melihat event, tapi business state belum committed di DB.

> Mengubah urutan (publish-then-commit) **tidak menyelesaikan** atomicity — hanya memindahkan window kegagalan.

---

## 5. Transactional Outbox Pattern

```
[BEGIN TRANSACTION]
     ├── UPDATE invoice SET status = 'paid'
     └── INSERT INTO outbox_events (id, event_type, payload, status='pending')
     │
     └── [COMMIT]

[Outbox Dispatcher — terpisah dari main flow]
     ├── SELECT FROM outbox_events WHERE status='pending'
     ├── Publish ke broker
     └── UPDATE status = 'published'
```

Transactional Outbox atomically records the business change and the intent to publish an event.

**Garanti (finite):**
- dispatcher can retry on transient failures
- duplicate delivery is possible (at-least-once)
- consumer **must be idempotent**
- **no exactly-once guarantee**

---

## 6. Exactly-Once Is Impossible

Exactly-once end-to-end is difficult. A practical system commonly uses:

```
at-least-once delivery
+
idempotent processing
```

**Failure scenario:**
```
worker receives event → business op succeeds → crash before ACK → message redelivered
```

The Idempotent Consumer pattern (consumer_name + event_id unique key) ensures duplicates are harmless.

---

## 7. DLQ Semantics

DLQ is NOT a trash bin. It stores messages that exceeded retry policy for investigation/reprocessing.

**Minimal DLQ record:**
- `event_id`
- `reason`
- `attempts`
- `failed_at`

Production systems need:
- monitoring (DLQ growth alerts)
- alerting
- manual/replay tooling

---

## 8. Retry Policy

| Error Type | Action |
|------------|--------|
| Transient (timeout, 503, connection reset) | Retry with backoff |
| Permanent (invalid payload, validation, unsupported event) | No retry, move to DLQ |

Retrying permanent errors every time wastes resources.

---

## 9. Retry Backoff

Tight retry loops are bad in production.

Use **exponential backoff + jitter**:
- retry after 100ms, 200ms, 400ms, 800ms...
- add random jitter 0–10% to avoid thundering herd

Tests stay fast with deterministic delays (no sleep loops).

---

## 10. Saga Business Semantics

Not every failure requires compensating everything.

**Example scenario:**
```
Create payment → Cash out → Generate journal → Sync ERP fails
```

If payment is already final and ERP is just downstream integration:
- **do NOT compensate** payment or cash out
- **retry ERP** instead (event-driven integration)

Saga compensation is for when the **business semantics requires an undo action**, not just technical rollback.

---

## 11. Saga Compensation vs DB Rollback

|  | Rollback | Compensation |
|--|----------|----------------|
| Scope | Transaction scope only | Business scope |
| Effect | Undo uncommitted changes | New business action that reverses |
| Example | `ROLLBACK` | `RefundPayment` (creates refund row) |

**Compensation creates new records** (e.g., RefundPayment), it does NOT delete historical data. This is crucial for audit trails.

---

## 12. Dispatcher Concurrency — At-Least-Once Overlap

Tanpa claim/locking, dua dispatcher yang berjalan bersamaan dapat mengambil event yang sama:

```
Dispatcher A: SELECT pending → [evt-1]
Dispatcher B: SELECT pending → [evt-1]
Dispatcher A: Publish evt-1 ✓
Dispatcher B: Publish evt-1 ✓   ← DUPLICATE
```

**At-least-once memang mengizinkan duplikat ini.** Konsekuensinya: consumer harus idempotent.

**Production alternatives untuk mengurangi duplikat:**

| Strategi | Keterangan |
|----------|-----------|
| `SELECT ... FOR UPDATE SKIP LOCKED` | Pessimistic lock per row (PostgreSQL) |
| Claim column (`claimed_by`, `claimed_at`) | Optimistic claim sebelum publish |
| Lease expiry | Dispatcher yang crash lepas klaim setelah timeout |
| Partitioning | Setiap dispatcher handle range event_id tertentu |
| Single dispatcher | Proses tunggal, scale via queue bukan fanout |

Pada lab ini yang menggunakan in-memory fake DB, tidak perlu implementasi `SELECT ... FOR UPDATE`. Yang penting: **consumer harus idempotent terlepas dari strategi dispatcher**.

---

## 13. Event Naming: Fact vs Command

| Tipe | Naming | Contoh |
|------|--------|--------|
| **Event** (fact already happened) | Past tense | `InvoicePaid`, `PaymentRecorded`, `VendorPaymentCompleted` |
| **Command** (request to perform action) | Imperative | `PayInvoice`, `SendNotification`, `SyncToERP` |

✅ Gunakan past-tense untuk events:
- `InvoicePaid`
- `InventoryUpdated`
- `CommissionGenerated`

❌ Jangan gunakan command-style untuk events:
- `SendWhatsapp`
- `GenerateCommission`
- `DoInventory`

**Alasan**: Event mencerminkan **state yang sudah ada**. Consumer yang menerima event harus dapat memproses dua kali (duplicate delivery) dan hasilnya tetap konsisten.

---

## 14. Event Payload: Explicit Struct, Not fmt.Sprintf

**Jangan:**
```go
payload := fmt.Sprintf(`{"invoice_id": %d, "timestamp": "%s"}`, id, t)
```

**Gunakan explicit struct:**
```go
type InvoicePaidPayload struct {
    EventID    string `json:"event_id"`
    InvoiceID  int    `json:"invoice_id"`
    OccurredAt string `json:"occurred_at"` // RFC3339
}
```

**Required fields untuk InvoicePaid:**
- `event_id` — untuk idempotent consumer (deduplikasi)
- `invoice_id` — aggregate identity
- `occurred_at` — RFC3339, untuk ordering dan audit

---

## 15. Separation of Business Fact vs Downstream Effect

`InvoicePaid` berarti **invoice sudah paid**, bukan berarti semua downstream sudah selesai.

**Bedanya:**

| Jenis | Contoh | Konsistensi |
|-------|--------|-------------|
| Business fact | Invoice paid, payment recorded | Strong consistency — satu DB transaction |
| Downstream projection | ERP sync, analytics | Eventual consistency OK |
| Integration side effect | WhatsApp, email | Fire-and-forget + retry |

---

## 16. Eventual Consistency

Setelah DB commit, downstream belum tentu up-to-date. Ini **normal**, bukan bug.

### Kapan Eventual Consistency Acceptable?
✅ Notification (WhatsApp, Email)  
✅ Analytics / Reporting  
✅ Audit logs  
✅ Search index updates

### Kapan Butuh Strong Consistency?
❌ Payment created + cash ledger updated + OPL marked paid → **satu DB transaction**  
❌ Inventory stok fisik  
❌ Balance / akuntansi  
❌ SLA-critical business invariants

---

## 17. At-Least-Once Delivery

Worker receives event → business op succeeds → crash before ACK → message redelivered.

Consumer **harus idempotent**. Gunakan `(consumer_name, event_id)` sebagai unique key untuk deduplikasi.

---

## 18. Saga & Compensation

```
Step A: Reserve Budget   → success
Step B: Process Payment  → FAIL
         │
[Compensate A: Release Budget]  ← hanya step yang sudah berhasil dikompensasi
```

Compensation dijalankan **terbalik** dari urutan eksekusi, hanya untuk step yang **sudah berhasil**.

**Compensation bukan Rollback:**
- ROLLBACK: undo uncommitted DB changes (tidak ada jejak)
- COMPENSATION: action baru yang membalikkan secara semantik (RefundPayment bukan delete payment row)

---

## 19. Cara Menjalankan

```bash
cd labs/03-database-transaction
go test -v -count=1
```

---

## Test Coverage

| Test | Konsep |
|------|--------|
| `TestUnsafeLocalTransaction` | Partial state corruption tanpa transaksi |
| `TestSafeLocalTransaction` | ACID rollback: all-or-nothing |
| `TestDistributedTransactionExternalSideEffectLimitation` | DB rollback ≠ WhatsApp rollback |
| `TestExternalSideEffectRollback` | payment+invoice rollback, notification tetap terkirim |
| `TestHTTPInsideTransactionDuration` | HTTP latency memperpanjang transaction lifetime |
| `TestCommitThenExternalCall` | COMMIT sebelum HTTP → transaction tidak terblokir latency |
| `TestTransactionStaysOpenDuringExternalCall` | `IsTxOpen()=true` selama external call (deterministic) |
| `TestTransactionCommitDoesNotBlockConnection` | `IsTxOpen()` transisi: closed→open→closed |
| `TestDualWriteProblemEventLost` | Invoice paid, event tidak terkirim |
| `TestDualWriteCrashWindow` | DB commit + crash = broker delivery = 0 |
| `TestReverseDualWriteFailure` | Publish-then-commit: consumer lihat event, DB belum committed |
| `TestTransactionalOutboxPatternAtomicity` | Business state + outbox atomic dalam satu tx |
| `TestOutboxDispatcherPublishesPending` | Dispatcher publish pending event |
| `TestOutboxDuplicateDeliveryAtLeastOnce` | Duplicate delivery: crash before mark published |
| `TestOutboxDispatcherConcurrencyOverlap` | Dua dispatcher concurrent publish event yang sama |
| `TestIdempotentConsumerDeduplication` | Deduplikasi per (consumer_name, event_id) |
| `TestTransientFailureSuccessAfterRetry` | Retry: fail-fail-succeed |
| `TestDeadLetterQueue` | Event masuk DLQ setelah max attempts |
| `TestSagaPaymentWithCompensatingAction` | Compensasi hanya untuk step yang berhasil |
| `TestSagaCompensationOrderFourSteps` | Reverse order: C, B, A (D tidak dikompensasi) |
| `TestSagaCompensationFailureHandling` | Saga lanjut kompensasi meski satu compensation gagal |
| `TestTransactionalOutboxRollback` | Outbox insert gagal → full rollback |
| `TestOutboxHappyPathAssertions` | Outbox event field: type, aggregate_id, status, attempts |
| `TestDispatcherPublishedAtSemantics` | published_at=null saat gagal, non-null setelah success |
| `TestInvoicePaidPayloadRoundTrip` | Serialize→deserialize, required fields, no forbidden fields |
| `TestEventualConsistencyDemo` | invoice=paid immediately, outbox processed later |
| `TestCompensationIdempotency` | Double RefundPayment = single refund (idempotent) |

---

## Navigasi

- **Previous**: [Lab 02 — Database Index](../02-database-index/)
- **Next**: [Lab 04 — Caching](../04-caching/)
