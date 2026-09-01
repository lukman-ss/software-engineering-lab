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
     ├── UPDATE invoices SET status = 'paid'
     └── INSERT wallet_transactions
     │
     └── [COMMIT]
```

Jika ada error → **ROLLBACK** → Semua pembatalan bersama.

> Gunakan DB transaction untuk melindungi **local invariants** (resources yang harus konsisten secara atomik dalam database yang sama).

---

## 2. Transaction Boundary

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

Contoh buruk:
```sql
DB::transaction(function () {
    Payment::create(...);
    Http::post($erpUrl, ...);
    WhatsappService::send(...);
});
```

Masalahnya:
- Lock hidup terlalu lama
- External service latency tidak predictable
- Timeout → tidak ada garansi rollback bahkan jika gagal
- Database connection tertahan
- Throughput turun
- Deadlock probability meningkat
- API mungkin sudah melakukan side effect walaupun DB rollback

Pendekatan yang benar:
```
BEGIN
  Update local business state
  Insert event/outbox intent
COMMIT
```
Kemudian external work diproses secara asynchronous jika business requirement memungkinkan.

**Caveat**: Tidak semua external call wajib asynchronous. Yang penting adalah memahami consistency requirement dan transaction boundary. Jangan memakai queue hanya karena ingin memakai queue.

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

Eventual consistency berarti beberapa state pada sistem distributed boleh sementara tidak sinkron, tetapi sistem memiliki mekanisme yang membuat state tersebut pada akhirnya menuju kondisi konsisten yang diharapkan.

Contoh:
```
10:00:00 Payment = PAID
10:00:01 ERP = PENDING
10:00:05 ERP retry
10:00:06 ERP = SYNCED
```

Ini **normal**, bukan bug.

### Kapan Eventual Consistency Acceptable?
✅ Notification (WhatsApp, Email)  
✅ Analytics / Reporting  
✅ Audit logs  
✅ Search index updates

### Kapan Butuh Consistency Langsung?
❌ Payment + ledger + OPL marking → harus satu transaksi
❌ Inventory stok fisik
❌ Balance / akuntansi

---

## 8. Retry & Idempotent Consumer (Pengantar)

Karena sistem terdistribusi bisa gagal sebagian (partial failure), **retry** menjadi mekanisme wajib.

**Producer bisa gagal:**
- Payment gateway timeout
- Network error
- Service unavailable

**Broker bisa tidak tersedia:**
- Connection refused
- Disk full
- Maintenance window

**Consumer bisa crash:**
- OOM killer
- Process restart
- Deployment

**Acknowledgment bisa hilang:**
- TCP drop
- Network partition
- Broker crash

Karena itu, **retry berarti event atau request bisa dikirim/diterima lebih dari sekali (at-least-once delivery)**.

Oleh karena itu, setiap penerima pesan (consumer) **harus idempotent**.

### Observability Requirement

Workflow distributed yang retryable tetapi tidak observable akan sulit dioperasikan di production. Setiap event dan proses harus memiliki:

| Field | Purpose |
|-------|---------|
| `correlation_id` | Melacak request chain yang terkait |
| `event_id` | Unik identifier untuk event |
| `aggregate_id` | Entity yang memicu event |
| `attempt` | Nomor upaya retry saat ini |
| `status` | pending / published / dead_lettered |
| `last_error` | Pesan error terakhir |
| `created_at` | Timestamp pembuatan |
| `processed_at` | Timestamp pemrosesan |

Koneksi queue tidak menjamin delivery. Architecture tetap membutuhkan:

- **Durability** (persist to disk)
- **Retry** (handle transient failure)
- **Idempotency** (handle duplicate delivery)
- **Observability** (debug & monitor)
- **Recovery** (handle crashes)

### Queue vs Delivery Guarantee

Jangan mengira bahwa menggunakan queue = aman. Architecture tetap membutuhkan:

- Producer bisa gagal
- Broker bisa unavailable  
- Consumer bisa crash
- ACK bisa hilang
- Message bisa redelivered

### Konsep `processed_events`:
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

### Sistem Synchronous (Transaction Local)

Jika semua state berada di satu database:

```
Admin klik Bayar
    ↓
Buat pembayaran
    ↓
Cash Out
    ↓
Update status OPL
    ↓
Generate jurnal
    ↓
Kirim WhatsApp Vendor
    ↓
Sync ERP
```

Jika `payment`, `cash_out`, `opl`, `journal` berada pada satu database dan merupakan satu atomic business invariant:

```sql
BEGIN
  Create payment
  Create cash_out
  Update OPL
  Create journal
COMMIT
```

Kemudian:

```
VendorPaymentCompleted event
        ↓
WhatsApp Worker
ERP Sync Worker
```

Jika ERP gagal:
- payment tetap committed
- ERP sync = pending/retry

Kecuali business requirement secara eksplisit menyatakan pembayaran tidak boleh dianggap berhasil sebelum ERP confirmation.

Tekankan bahwa architecture mengikuti business invariant.

### Sistem Asynchronous (Saga)

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

> Catatan: Compensation bukanlah cara untuk mengembalikan sistem seperti semalanya. Email yang sudah terkirim tidak bisa benar-benar di-rollback. Yang mungkin dilakukan: kirim correction email atau corrective business action.

---

## 10. Architecture Selection: Synchronous vs Asynchronous

Jangan mengajarkan bahwa `Controller -> Event -> Queue` selalu lebih baik. Architecture dipilih berdasarkan requirement.

### Synchronous Local Transaction cocok jika:
- User membutuhkan immediate result
- Consistency harus langsung diketahui
- Operation sederhana
- Dependency reliability acceptable

### Asynchronous/Event-Driven cocok jika:
- Pekerjaan tidak harus selesai sebelum response
- Latency tinggi
- Integration dapat gagal sementara
- Retry diperlukan
- Loose coupling penting

---

## 11. Local vs Distributed Workflow

### Kondisi A — Semua state berada pada database yang sama

Jika `payments`, `invoice`, `inventory`, dan `commission` berada dalam *satu database relasional yang sama*, dan perubahan ini dianggap sebagai satu *business invariant* di mana kegagalan parsial tidak boleh terlihat:

```sql
BEGIN
  Insert payment
  Update invoice
  Deduct inventory
  Create commission
COMMIT
```
Pendekatan ini **boleh dan tepat** (Local Transaction).

### Kondisi B — State berada pada boundary berbeda

Jika arsitektur sudah terdistribusi:

```
PostgreSQL (Order Service)
+
ERP API
+
WhatsApp API
+
Kafka (Message Broker)
```

Tidak ada satu `DB::transaction()` yang dapat melakukan `ROLLBACK` pada semuanya jika terjadi kegagalan di tengah jalan.

Gunakan: **Local Transaction + Event + Retry + Idempotency + Compensation** sesuai kebutuhan.

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
| **Eventual Consistency** | `TestEventualConsistencyDemo` | Local vs downstream divergence |

---

## 14. How to Run

```bash
cd labs/03-database-transaction
go test -v -count=1
```

---

## Navigasi

- **Previous**: [Lab 02 — Database Index](../02-database-index/)
- **Next**: [Lab 04 — Caching](../04-caching/)