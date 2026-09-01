# Lab 03 — Database Transaction & Distributed Transaction Boundary

Kenapa Database Transaction Saja Tidak Cukup?

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

> Gunakan DB transaction untuk melindungi **local invariants** (resources yang harus konsisten secara atomik dalam database yang sama).

---

## 2. Where DB Transaction Stops (Transaction Boundary)

Satu local database transaction hanya menjamin atomicity untuk resource yang berpartisipasi. External HTTP APIs, email, WhatsApp, message broker, cloud storage, payment gateway request, ERP request, dan external systems lainnya **tidak ikut berpartisipasi**.

*(XA/2PC dapat koordinasi distributed resource, tapi bukan fokus lab ini.)*

---

## 3. External Side Effects

Apa yang terjadi ketika kita mencampur HTTP request ke dalam DB transaction?

```
[BEGIN] → INSERT payment → HTTP WhatsApp → ERP sync → [COMMIT]
```

**Risiko Utama**: Transaction bertahan lama → Connection pool tertahan → Throughput menurun drastis.

**External call di dalam transaksi**:
- BEGIN → UPDATE → HTTP → COMMIT (transaction terblokir latency)
- BEGIN → UPDATE → COMMIT → HTTP (transaction cepat, external service tidak terpengaruh latency)

> **Catatan senior**: HTTP timeout tidak selalu berarti external operation gagal. Remote service bisa berhasil tapi response hilang akibat jaringan. Karena itu external commands seperti payment/refund/reservation juga perlu idempotency key + status reconciliation. Lihat Daily #1 Idempotency untuk detail.

---

## 4. HTTP Inside DB Transaction Limitation

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

Jika terjadi `ROLLBACK` di DB, **HTTP request, WhatsApp, atau email yang sudah terkirim tidak dapat dibatalkan oleh database.** External side effects bersifat independen dari transaction state.

Test bukti: `TestHTTPInsideTransactionDuration`, `TestCommitThenExternalCall`, `TestExternalSideEffectRollback`.

---

## 5. Dual-Write Problem

Ketika kita mencoba menyelesaikan permasalahan integrasi dengan message broker, kita dihadapkan pada **Dual-Write Problem**:

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

## 6. Transactional Outbox (Pengantar)

Untuk menyelesaikan dual-write problem, kita gunakan **Transactional Outbox**.

```
[BEGIN TRANSACTION]
     ├── UPDATE invoice SET status = 'paid'
     └── INSERT INTO outbox_events (...)
     │
     └── [COMMIT]
```

Transactional Outbox **secara atomik mencatat business change dan *intent* (niat) untuk mengirimkan event ke broker.**

Flow lengkapnya:
1. Business transaction + Insert outbox event
2. **Outbox Worker** berjalan asinkron
3. **Publish event** ke broker
4. **Mark as processed** di tabel outbox

> 💡 **Deep Dive**: Implementasi production-grade Transactional Outbox (dispatcher internals, FOR UPDATE SKIP LOCKED, concurrency, retry dengan backoff, crash-after-publish recovery, DLQ) dibahas secara spesifik pada **Lab 07 — Outbox Pattern**. Lab ini hanya memperkenalkan konsep atomicity-nya.

---

## 7. Eventual Consistency (Pengantar)

Setelah DB commit dan event dicatat, system hilir (downstream) belum tentu *up-to-date*. Ini **normal**, bukan bug.

### Kapan Eventual Consistency Acceptable?
✅ Notification (WhatsApp, Email)  
✅ Analytics / Reporting  
✅ Audit logs  
✅ Search index updates

---

## 8. Retry & Idempotent Consumer (Pengantar)

Karena sistem terdistribusi bisa gagal sebagian (partial failure), **retry** menjadi mekanisme wajib.
Namun, retry berarti event atau request bisa dikirim/diterima lebih dari sekali (at-least-once delivery).

Oleh karena itu, setiap penerima pesan (consumer) **harus idempotent**.

Konsep `processed_events`:
```
processed_events
---------------------------------
consumer_name
event_id
processed_at
UNIQUE (consumer_name, event_id)
```

Gunakan pola:
```
At-least-once delivery
+
Idempotent consumer
=
Effectively-once business effect
```

> 💡 **Deep Dive**: Retry policy, eksponensial backoff, dan Dead Letter Queue (DLQ) dibahas lebih dalam di **Lab 08 — Retry** dan **Lab 07 — Outbox Pattern**.

---

## 9. Saga / Compensation (Pengantar)

Database rollback (`ROLLBACK`) hanya bisa membatalkan transaksi yang belum di-commit.
Lalu bagaimana jika kita perlu membatalkan *distributed workflow* yang sudah sebagian berhasil?

```
Create payment → Cash out → Generate journal → Sync ERP fails
```

Jika payment sudah *final*, kita membutuhkan **Saga Compensation** untuk **semantic undo**, bukan technical rollback.

| | Rollback | Compensation |
|--|----------|----------------|
| Scope | Transaction scope | Business scope |
| Effect | Undo uncommitted | New action yang membalikkan efek (reverses) |
| Example | `ROLLBACK` | `RefundPayment` (membuat record refund baru, bukan DELETE payment) |

Saga compensation membatalkan langkah (step) yang *sebelumnya* telah berhasil dengan menjalankan *compensating action*. Compensation itu sendiri juga harus *idempotent*.

---

## 10. Rule of Thumb: Menentukan Transaction Boundary

### CMMS Case Study (Bengkel / Bengkel Management System)

Bayangkan flow **Invoice Paid**. Apa yang terjadi setelah invoice dibayar?
- Mengubah Payment state
- Memotong Inventory stok
- Menghitung Commission mekanik
- Mengubah Membership / Poin pelanggan
- Mengirim WhatsApp notifikasi
- Sinkronisasi ke ERP

**PENTING**: Tidak semua operasi tersebut otomatis dapat (atau harus) berada dalam *satu* transaksi database.

#### Kondisi A — Semua state berada pada database yang sama

Jika `payments`, `invoice`, `inventory`, dan `commission` berada dalam *satu database relasional yang sama*, dan perubahan ini dianggap sebagai satu *business invariant* di mana kegagalan parsial (partial failure) sama sekali tidak boleh terjadi/terlihat:

```sql
BEGIN;
  Insert payment;
  Update invoice;
  Deduct inventory;
  Create commission;
COMMIT;
```
Pendekatan ini **boleh dan tepat** (Local Transaction).

#### Kondisi B — State berada pada boundary berbeda

Jika arsitektur sudah terdistribusi:
- `PostgreSQL` (Order Service)
- + `ERP API`
- + `WhatsApp API`
- + `Kafka` (Message Broker)

Tidak ada satu `DB::transaction()` yang dapat melakukan `ROLLBACK` pada semuanya jika terjadi kegagalan di tengah jalan.

Gunakan: **Local Transaction + Event + Retry + Idempotency + Compensation** sesuai kebutuhan.

---

## 11. Transaction Design Rules & Consistency

Jika beberapa perubahan merupakan satu *business invariant*, berada pada transactional datastore yang sama, dan kegagalan parsial tidak boleh terlihat, perubahan tersebut sebaiknya dilakukan secara atomik dalam satu local database transaction.

Namun, jika state berada pada service, database, atau external system yang berbeda, satu local database transaction tidak dapat memberikan atomicity lintas boundary.

Contoh:
- `payments` + `journal_entries` di dalam satu DB yang sama → **Local Transaction**
- `payment-service DB` + `accounting-service DB` yang terpisah → **Distributed Workflow**

---

## 12. MockDB Concurrency Model (Educational)

MockDB mengimplementasikan transaction model yang **serialized**:

- `BEGIN` → acquire global transaction mutex
- `INSERT/UPDATE/DELETE` → operate on transaction snapshot
- `COMMIT` → merge snapshot to committed, release mutex
- `ROLLBACK` → discard snapshot, release mutex

Goroutines are concurrent at the application level, while MockDB intentionally serializes write transactions to provide deterministic educational transaction semantics. This verifies application-level logical race paths and atomic dedup behavior under the MockDB model.

---

## 13. Test Organization

| Category | Tests | Purpose |
|----------|-------|---------|
| **Local Transaction** | `TestUnsafeLocalTransaction`, `TestSafeLocalTransaction` | ACID verification |
| **External Side Effects** | `TestDistributedTransactionExternalSideEffectLimitation`, `TestHTTPInsideTransactionDuration`, `TestCommitThenExternalCall` | Boundary & external side-effect rollback limitation |
| **Dual-Write Problems** | `TestDualWriteProblemEventLost`, `TestReverseDualWriteFailure` | Atomicity window analysis |
| **Transactional Outbox** | `TestTransactionalOutboxPatternAtomicity`, `TestOutboxDispatcherPublishesPending` | Outbox concept atomicity |
| **Idempotent Consumer** | `TestIdempotentConsumerDeduplication`, `TestAtomicConsumerFlow` | Deduplication key design |
| **Compensation/Saga** | `TestSagaPaymentWithCompensatingAction`, `TestSagaCompensationFailureHandling` | Semantic undo conceptual demo |
| **Eventual Consistency**| `TestEventualConsistencyDemo` | Local vs downstream divergence |

---

## Cara Menjalankan

```bash
cd labs/03-database-transaction
go test -v -count=1
```

---

## Navigasi

- **Previous**: [Lab 02 — Database Index](../02-database-index/)
- **Next**: [Lab 04 — Caching](../04-caching/)