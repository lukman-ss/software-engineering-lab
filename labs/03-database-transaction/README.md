# Lab 03 — Database Transaction & Distributed Transaction Boundary

## 1. Transaction Boundary: Definisi Fundamental

**Definisi:**

> **Transaction boundary** adalah batas resource/state yang dapat dijamin atomik oleh satu transaction.

### Visualisasi Transaction Boundary

**Satu Database — Satu Transaction Boundary:**

```
PostgreSQL connection
┌───────────────────────────┐
│ payment                   │
│ invoice                   │
│ inventory                 │
│ journal                   │
└───────────────────────────┘
        local transaction
```

Semua operasi di dalam boundary dijamin atomic. Jika satu gagal, semua rollback.

---

**Multiple Resources — Multiple Transaction Boundaries:**

```
PostgreSQL
    +
ERP API
    +
WhatsApp API
    +
Kafka
```

Batas tersebut **tidak dapat dijamin oleh satu local DB::transaction()**.

---

### Mental Model

| Kondisi | Rollback Bekerja? |
|---------|-------------------|
| Same transaction boundary | ✅ Ya, rollback dapat bekerja |
| Outside transaction boundary | ❌ Tidak, rollback database tidak cukup |

---

## 2. Local Atomicity: Apa yang Dikelola Database Transaction?

**Definisi Teknis yang Akurat:**

> Local database transaction memberikan atomicity pada satu transactional datastore atau transaction boundary yang sama. Ketika workflow melibatkan resource independen seperti database lain, external API, payment gateway, message broker, email, object storage, atau ERP, local transaction tersebut tidak dapat secara otomatis melakukan rollback terhadap side effect di luar transaction boundary.

**Catatan Penting tentang Distributed Transaction:**

Distributed transaction lintas resource secara teknis dapat dilakukan dengan mekanisme seperti Two-Phase Commit (2PC) atau XA pada sistem tertentu, tetapi membawa coupling, availability, operational complexity, dan scalability trade-off. Karena itu banyak arsitektur modern memilih local transaction + messaging + Saga + compensation + Outbox sesuai kebutuhan.

> **Nota**: Lab ini hanya memperkenalkan Two-Phase Commit (2PC) dan XA sebagai konteks bahwa distributed transaction lintas resource secara teknis memang ada, tetapi tidak membahas implementasinya secara mendalam. Lab 07 tetap fokus pada Transactional Outbox, bukan 2PC/XA.

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

## 3. Business Invariant: Aturan yang Harus Selalu Benar

**Definisi:**

> **Business invariant** adalah aturan bisnis yang harus selalu benar meskipun terjadi failure.

### Kapan Membutuhkan Strong Consistency?

Strong/local atomic consistency diperlukan ketika beberapa perubahan merupakan satu business invariant, berada pada transactional datastore yang sama, dan partial state tidak boleh terlihat.

**Contoh:**
Jika \`payment\` + \`journal\` + \`OPL status\` berada pada database transactional yang sama dan merupakan satu business invariant, maka gunakan **satu local database transaction**.

Tetapi jika:
- \`payment-service DB\` + \`accounting-service DB\`
- atau \`local DB\` + \`external inventory service\`

Maka satu local DB transaction tidak dapat memberikan atomicity lintas boundary.

> **Penting**: Jenis data tidak menentukan transaction strategy. Business invariant + transaction boundary yang menentukan.

### Contoh Business Invariant

**Pembayaran:**

Jika payment tercatat PAID, maka payment record dan financial journal yang wajib menjadi bagian dari transaksi tidak boleh berada pada state setengah jadi.

```
Business invariant:
- Payment PAID → Journal harus balance
- Payment FAILED → Tidak ada journal
- State setengah jadi → TIDAK BOLEH
```

---

## 4. Transaction Size: "Large Enough, Small Enough"

Semakin banyak consistency boundary yang terlibat, semakin tidak realistis mengandalkan satu local database transaction untuk menjamin seluruh workflow.

**Gunakan transaction sebesar yang diperlukan untuk menjaga business invariant, tetapi sekecil mungkin agar lock duration, contention, dan coupling tetap terkendali.**

### Mental Model

```
transaction should be:
├── large enough to protect invariant
└── small enough to avoid unnecessary work
```

---

## 5. External Side Effects: Irreversible Actions

### Batasan Transaksi

Jika `InventoryService::deduct()`, `Payment::create()`, `CommissionService::calculate()` semua mengubah tabel yang berada pada **database yang sama** dan menggunakan **transaction yang sama**:

| Table | Status |
|-------|--------|
| payments | payment rows |
| invoices | invoice rows |
| inventory | inventory rows |
| commissions | commission rows |

Jika terjadi:

```
BEGIN
  INSERT payment
  UPDATE inventory
  INSERT commission
  ERP timeout
ROLLBACK
```

Maka **semua** perubahan di database akan dibatalkan:

| Table | Hasil Setelah Rollback |
|-------|------------------------|
| payments | ← rollback |
| inventory | ← rollback |
| commissions | ← rollback |

**Yang TETAP tidak dapat dibatalkan:**
- WhatsApp API (sudah terkirim)
- Email (sudah dikirim)  
- SMS (sudah dikirim)
- ERP API (sudah diproses)
- Payment Gateway API (sudah charging)

---

## 6. Critical vs Non-Critical Side Effects

**External tidak otomatis berarti non-critical. Criticality ditentukan oleh bisnis.**

### Critical Side Effects (tergantung business invariant)
- Record payment
- Update financial state  
- Reserve inventory
- Generate journal
- Sync to ERP (jika requirement wajib sebelum payment complete)

### Non-Critical / Secondary Side Effects (tergantung requirement)
- Email (receipt notification)
- WhatsApp notification
- Analytics / reporting
- Webhook to third party
- ERP reporting sync

> **Decision**: Tentukan apakah external step menentukan keberhasilan business operation, atau sekadar notification/performance optimization.

---

## 7. Two Inventory Scenarios: Local vs External

### Scenario 1: Inventory dalam Local Database

Jika `InventoryService::deduct()` hanya mengubah tabel inventory di **database yang sama**:

```
Payment | Invoice | Inventory | Commission  →  local DB transaction
```

Jika ERP timeout dan DB rollback:

| System | Status |
|--------|--------|
| Payment | ← rollback |
| Invoice | ← rollback |
| Inventory | ← rollback |
| Commission | ← rollback |
| WhatsApp | ← SUDAH terkirim |

**Kesimpulan:** Inventory **bisa** di-rollback jika berada di database yang sama.

---

### Scenario 2: Inventory sebagai External Service

Jika `InventoryService::deduct()` adalah **HTTP request ke inventory-service** yang memiliki database terpisah:

```
Main DB transaction
  ↓ HTTP
inventory-service database (terpisah)
```

Jika ERP timeout dan database utama rollback:

| System | Status |
|--------|--------|
| Main DB (payments, invoices) | ← rollback |
| Inventory API | ← SUDAH berhasil |
| WhatsApp API | ← SUDAH terkirim |

**Kesimpulan:** Inventory **tidak bisa** di-rollback jika berada di service terpisah.

---

## 8. Anti-Pattern: HTTP Inside DB Transaction

```
[BEGIN TRANSACTION]
     ├── UPDATE invoice SET status = 'paid'
     ├── HTTP call ke payment gateway
     └── [COMMIT]
```

**Dua masalah utama:**

1. **Connection pool tertahan** — Transaction menunggu HTTP, tidak bisa dilepaskan kembali ke pool.
2. **External call tidak dapat di-rollback** — WhatsApp, ERP, payment gateway tidak bisa "undone".

### Alasan Lengkap

- Transaction terbuka lebih lama
- Row/table lock hidup lebih lama  
- Connection database tertahan
- External latency tidak predictable
- Timeout external dapat menahan transaction
- Deadlock probability meningkat
- Throughput turun
- External side effect mungkin sukses walaupun DB rollback

> **Nuansa**: Masalah utamanya bukan karena setiap HTTP call di dalam transaction secara mutlak dilarang, tetapi karena local transaction tidak mampu mengontrol atomicity external side effect dan membuat transaction lifetime bergantung pada resource eksternal.

---

## 9. Timeout != Operation Pasti Gagal (Unknown Outcome)

```
POST /erp/payment
    ↓
ERP berhasil memproses
    ↓
network response hilang
    ↓
client timeout
```

Client melihat **TIMEOUT**, tetapi server sebenarnya: **SUCCESS**

### Unknown Outcome

Karena timeout memberikan **unknown outcome**, integration yang critical harus memiliki mekanisme untuk:

```
mencegah duplicate execution
+
mengetahui status operation
+
menemukan inconsistency yang tertinggal
```

Contoh:

**Idempotency-Key**
→ mencegah duplicate side effect

**Query/status endpoint**
→ mengecek apakah operation sebenarnya berhasil

**Reconciliation**
→ mendeteksi state yang belum sinkron

**Correlation/operation ID**
→ menghubungkan request, retry, log, dan external operation

> **Hubungkan ke Lab 01 Idempotency** untuk detail implementasinya.

---

## 10. Dual-Write Problem

> "Commit DB dulu baru publish ke message broker!"

Ini dapat menyebabkan **Dual-Write Problem**:

```
[UPDATE invoice]
[COMMIT]               → SUCCESS
                          │
                 X server crash
                          │
[Publish event]        → TIDAK pernah terjadi
```

**Hasil**: Database = PAID, Event = hilang.

### Reverse Dual-Write juga bermasalah

```
[Publish event]        → SUCCESS
                          │
                 X server crash
                          │
[COMMIT]               → TIDAK pernah terjadi
```

Mengubah urutan **tidak menyelesaikan** atomicity — hanya memindahkan window kegagalan.

---

## 11. Transactional Outbox

Untuk menyelesaikan dual-write problem, gunakan Transactional Outbox:

```
[Local transaction]
     ↓
[Record business state + event intent atomically]
     ↓
[COMMIT]
     ↓
[Publisher/Dispatcher mengirim event asinkron]
```

### Mental Model Outbox:

**Outbox**
→ menyimpan business change + event intent secara atomik

**Dispatcher**
→ mendukung at-least-once delivery

**Consumer**
→ harus idempotent atau memiliki mekanisme deduplication

> At-least-once berarti sistem berusaha agar message tidak hilang dengan mekanisme retry/redelivery, tetapi duplicate delivery tetap mungkin terjadi.

Contoh:
```
Publish berhasil
↓
dispatcher crash sebelum mark sent
↓
event dikirim ulang
```

Karena itu `duplicate delivery` bukan bug yang mustahil terjadi, melainkan sesuatu yang harus diantisipasi oleh consumer.

Pada delivery model yang memungkinkan retry atau redelivery, seperti at-least-once delivery, consumer harus dirancang idempotent atau memiliki mekanisme deduplication yang ekuivalen.

**Contoh:**
\`InvoicePaid event_id = evt-123\`

- consumer menerima evt-123
- berhasil
- ack hilang
- evt-123 dikirim ulang

Consumer tidak boleh menghasilkan duplicate side effect.

> **Hubungkan ke Lab 01 Idempotency** untuk detail implementasinya.

> **Deep Dive**: Implementasi production-grade Transactional Outbox dibahas secara spesifik pada **Lab 07 — Outbox Pattern**.

---

## 12. Saga Pattern & Compensation

Saga mengelola distributed transactions melalui series of local transactions dengan aksi kompensasi.

```
Step 1: Reserve Hotel    → Action: reserve | Compensate: release
Step 2: Book Flight      → Action: book    | Compensate: cancel
Step 3: Book Car (GAGAL)
```

Jika Step 3 gagal → Saga mengeksekusi compensation untuk Step 2 lalu Step 1.

> **Penting**: Compensation adalah **business operation baru** yang mengoreksi efek sebelumnya, bukan teknis rollback.

---

## 13. Eventual Consistency

Retry merupakan mekanisme umum untuk menangani transient failure pada distributed workflow, tetapi tidak semua failure boleh di-retry.

**Retryable:**
- timeout
- temporary network failure
- HTTP 429
- HTTP 502/503/504
- temporary broker unavailable

**Umumnya tidak retryable (tanpa perubahan input/state):**
- validation error
- malformed payload
- unauthorized / forbidden
- business rule rejection
- resource not found yang memang permanent

> Retry policy harus membedakan transient failure dan permanent failure.

---

Eventual consistency berarti beberapa komponen dapat sementara melihat state yang berbeda, tetapi sistem memiliki mekanisme untuk membawa state tersebut menuju kondisi konsisten yang diharapkan.

### Timeline Contoh

```
t=0:  Payment = PAID (di DB langsung)
t=1:  ERP = PENDING
t=2:  ERP retry
t=3:  ERP = SYNCED
```

Temporary inconsistency dapat menjadi kondisi normal pada eventual consistency jika memang sesuai consistency requirement, masih berada dalam batas waktu/staleness yang dapat diterima, dan sistem memiliki mekanisme untuk membawa state tersebut menuju kondisi konsisten.

Eventual consistency != inconsistent selamanya.

Retry, recovery, reconciliation, monitoring, dan manual intervention adalah mekanisme yang dapat digunakan sesuai criticality dan consistency requirement.

> **Penting**: Eventual consistency ≠ state boleh inconsistent selamanya.

Contoh yang sehat:
- Payment = PAID
- ERP = SYNCED setelah retry

Bukan kondisi yang boleh dianggap normal:
- Payment = PAID
- ERP = PENDING selamanya
- tidak ada retry
- tidak ada reconciliation
- tidak ada alert

---

## 14. Failure Matrix

| Failure | Local DB | External System | Strategy |
|---------|----------|-----------------|----------|
| DB update gagal | rollback | belum dipanggil | local transaction |
| WhatsApp gagal | committed | failed | retry |
| ERP timeout (unknown outcome) | committed | unknown | idempotency + status check |
| Worker crash | committed | pending | queue retry |
| Event publish gagal setelah DB commit | committed | missing | Outbox |
| Duplicate delivery | committed | risk duplicate | idempotent consumer |
| Compensation gagal | committed | inconsistent sementara | retry + reconciliation |

---

## 15. Glossary

| Istilah | Definisi |
|---------|----------|
| **ACID** | Properties transaksi: Atomicity, Consistency, Isolation, Durability |
| **atomicity** | Semua operasi dalam transaksi berhasil atau semua gagal |
| **local transaction** | Transaksi pada satu datastore/txn boundary |
| **transaction boundary** | Batas resource yang dapat dijamin atomic oleh satu transaction |
| **business invariant** | Aturan bisnis yang harus selalu benar meski terjadi failure |
| **external side effect** | Aksi luar dari database (API call, notification) yang tidak ikut transaction |
| **partial failure** | Beberapa langkah berhasil, beberapa gagal — tidak ada atomicity lintas sistem |
| **distributed workflow** | Workflow yang melibatkan beberapa sistem/database |
| **distributed consistency** | Strategi menjaga atau mencapai konsistensi state ketika workflow melibatkan beberapa independent resource atau system boundary |
| **dual-write problem** | Masalah konsistensi saat mencoba menulis ke DB lalu publish event secara terpisah |
| **event-driven** | Pola arsitektur di mana komponen menghasilkan dan/atau bereaksi terhadap event untuk mengomunikasikan perubahan state atau kejadian bisnis |
| **retry** | Mencoba kembali operasi yang gagal |
| **idempotency** | Operasi yang dapat dipanggil berulang kali tanpa efek samping tambahan |
| **at-least-once** | Delivery model tempat event bisa dikirim berulang kali |
| **eventual consistency** | State pada beberapa komponen dapat sementara berbeda, tetapi sistem dirancang agar pada akhirnya menuju kondisi konsisten yang diharapkan |
| **Saga** | Pola untuk distributed transactions dengan kompensasi |
| **compensation** | Aksi pengembalian untuk membatalkan efek operation sebelumnya |
| **Outbox** | Mekanisme untuk menyimpan event secara atomic bersama business state |
| **DLQ** | Dead Letter Queue — menampung message yang gagal setelah retry maksimal |
| **reconciliation** | Proses menemukan dan memperbaiki inconsistency yang lolos normal processing |
| **observability** | Kemampuan memantau dan memdebug sistem melalui logs, metrics, tracing |
| **unknown outcome** | Ketika timeout terjadi, kita tidak tahu apakah operasi berhasil atau tidak |

> Event-driven architecture dapat coexist dengan synchronous HTTP/RPC dalam sistem yang sama. Keduanya dipilih berdasarkan kebutuhan interaction, consistency, latency, dan failure handling.

---

## 16. Failure Scenarios untuk Reviewer

### Scenario 1: Inventory Local DB

```
Main DB:
  payment = created
  inventory = deducted

External:
  WhatsApp = sent

ERP = timeout

Jika inventory LOCAL:
  DB rollback
  payment = absent
  inventory = restored
  WhatsApp tetap terkirim
```

### Scenario 2: Inventory External Service

```
Main DB:
  payment = created
  inventory = deducted

External:
  WhatsApp = sent
  Inventory API = success

ERP = timeout

Jika inventory EXTERNAL:
  Main DB rollback
  Inventory API = STILL deducted
  WhatsApp tetap sent
```

---

## 17. Code References

- `local_transaction.go` — Payment service dengan/tanpa transaction
- `external.go` — Implementasi WhatsApp, HTTP, Outbox, Saga, Idempotent Consumer  
- `transaction_test.go` — Semua contoh diverifikasi melalui unit test

---

## 18. Definition of Done

Lab 03 dianggap selesai bila reader dapat menjawab dengan benar:

1. Apa yang sebenarnya dijamin oleh DB transaction?
2. Apa itu transaction boundary?
3. Kenapa external API tidak ikut rollback?
4. Kapan inventory bisa rollback dan kapan tidak?
5. Mengapa HTTP call di dalam transaction berbahaya?
6. Apa itu partial failure?
7. Apa itu distributed consistency problem?
8. Apa itu dual-write problem?
9. Mengapa commit DB lalu publish event belum aman?
10. Apa fungsi Outbox?
11. Mengapa Outbox tidak berarti exactly-once?
12. Mengapa consumer harus idempotent?
13. Mengapa retry bisa menghasilkan duplicate side effect?
14. Apa arti timeout sebagai unknown outcome?
15. Apa itu Saga?
16. Apa itu compensation?
17. Mengapa compensation bukan rollback?
18. Apa itu eventual consistency?
19. Apa fungsi DLQ?
20. Apa fungsi reconciliation?
21. Kapan synchronous lebih tepat?
22. Kapan asynchronous lebih tepat?
23. Bagaimana menentukan operation yang harus berada dalam satu transaction?
24. Apa hubungan business invariant dengan transaction boundary?
25. Bagaimana mendesain flow pembayaran Vendor OPL jika ERP unavailable?