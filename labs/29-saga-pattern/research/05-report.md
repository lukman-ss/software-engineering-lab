# Research Report

## Research Question

Bagaimana menjaga konsistensi data di beberapa service/database berbeda ketika 2PC (Two-Phase Commit) terlalu lambat atau tidak memungkinkan, dan bagaimana Saga Pattern (Choreography vs Orchestration) dengan Compensating Transactions, Dual-Write mitigation, dan idempotency menjawab masalah tersebut?

## Executive Summary

Saga Pattern memecah transaksi terdistribusi menjadi rangkaian transaksi lokal. Setiap langkah commit di database-nya sendiri. Jika langkah gagal, transaksi kompensasi membatalkan langkah sebelumnya secara semantik (bukan rollback ACID). Dua gaya koordinasi: Choreography (event-driven, tanpa koordinator) dan Orchestration (koordinator terpusat). Pola ini **bukan** pengganti ACID isolation; saga menyerahkan "I" (isolation) dan menerima eventual consistency plus anomali data yang harus dimitigasi. Dual-write (update DB + publish event secara atomik) adalah prasyarat yang diselesaikan oleh Transactional Outbox. Idempotency wajib karena delivery at-least-once. Lab spec akurat pada definisi, dua pendekatan, kompensasi, dan kapan memakai/tidak memakai saga.

## Findings

### Finding 1: Definisi dan Asal-Usul Saga

**Claim:**
Saga diperkenalkan Garcia-Molina & Salem (ACM SIGMOD 1987) untuk transaksi database berumur panjang; diadaptasi ke microservices sebagai sequence of local transactions + compensating transactions.

**Evidence:**
Paper asli di-host Cornell (`sagas.pdf`). Definisi modern: "A saga is a sequence of local transactions. Each local transaction updates the database and publishes a message or event to trigger the next." (Richardson). Microsoft: setiap local transaction "Completes its work atomically within a single service."

**Sources:**
- Garcia-Molina & Salem 1987 — https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf
- Microservices.io Pattern: Saga — https://microservices.io/patterns/data/saga.html
- Microsoft Azure Architecture Center — https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

### Finding 2: Choreography vs Orchestration

**Claim:**
Choreography: setiap service publish domain event yang memicu service berikutnya. Orchestration: satu orchestrator mengirim command ke peserta dan menangani kompensasi.

**Evidence:**
Keduanya didokumentasikan Microsoft (tabel benefit/drawback), Richardson (diagram Create Order Saga choreography vs orchestration), Temporal (analogi ant colony vs air-traffic control). Choreography: loose coupling, no SPOF, tapi spaghetti event + cyclic dependency + sulit test. Orchestration: flow jelas, no cyclic dep, tapi coordinator SPOF + complexity.

**Sources:**
- Microsoft Azure Architecture Center — https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Microservices.io — https://microservices.io/patterns/data/saga.html
- Temporal Blog 2023-07-13 — https://temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question

**Confidence:** HIGH

---

### Finding 3: Compensating Transactions adalah Semantic Undo, Bukan Rollback

**Claim:**
Setiap aksi positif yang reversible membutuhkan aksi kompensasi setara (Reserve Stock → Release Stock; Charge Card → Refund; Create Pending Order → Cancel Order). Kompensasi dirancang developer, bukan otomatis dari database.

**Evidence:**
Richardson: "Lack of automatic rollback — a developer must design compensating transactions that explicitly undo changes." Microsoft: "Compensable transactions can be undone or compensated for by other transactions with the opposite effect." Mapping lab spec cocok definisi.

**Sources:**
- Microservices.io Pattern: Saga
- Microsoft Azure Architecture Center
- AWS Prescriptive Guidance Saga Pattern — https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

**Confidence:** HIGH

---

### Finding 4: Pivot dan Retryable Transactions (Extended Concept)

**Claim:**
Microsoft membagi langkah saga menjadi compensable (bisa di-undo), pivot (point of no return), dan retryable (harus selesai, idempotent, setelah pivot). Taksonomi tiga jenis ini didokumentasikan secara rinci di Microsoft Azure Architecture Center. Richardson dalam *Microservices Patterns* (Chapter 4) merujuk countermeasures tanpa menyebut taksonomi eksplisit di halaman publik.

**Evidence:**
Hanya Microsoft Architecture Center yang merinci taksonomi tiga jenis ini secara publik. Richardson merujuk "countermeasures" di buku tanpa enumerasi di halaman publik.

**Sources:**
- Microsoft Azure Architecture Center — https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** MEDIUM (satu sumber otoritatif rinci)

---

### Finding 5: 2PC vs Saga — Tradeoff Availability vs Isolation

**Claim:**
2PC menahan lock lintas service, menurunkan availability, coupling tinggi. Saga melepaskan isolation (ACID I), memakai eventual consistency, tidak memblokir peserta lain.

**Evidence:**
Richardson Forces: "2PC is not an option." AWS: "long-lived transactions and you don't want other microservices to be blocked." Microsoft: "there's no built-in isolation across services." Lab spec ("2PC terlalu lambat atau tidak memungkinkan") terverifikasi.

**Sources:**
- Microservices.io Pattern: Saga
- AWS Prescriptive Guidance
- Microsoft Azure Architecture Center

**Confidence:** HIGH

---

### Finding 6: Dual-Write Problem dan Transactional Outbox

**Claim:**
Tidak bisa atomik update database DAN publish ke message broker (Kafka tidak XA). Solusi: tulis business row + outbox row dalam satu local ACID transaction; CDC (Debezium) atau polling merelay ke broker. Memberi read-your-own-writes di service sumber + eventual consistency ke konsumen.

**Evidence:**
Morling 2019: outbox table schema (id UUID, aggregatetype, aggregateid, type, payload JSONB); persist+delete trick agar table kosong tapi WAL berisi INSERT. Richardson Related Patterns: Event Sourcing, Transactional Outbox.

**Sources:**
- Debezium Blog 2019-02-19 — https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- Microservices.io Pattern: Saga (Related patterns)

**Confidence:** HIGH

---

### Finding 7: Idempotency Wajib karena At-Least-Once Delivery

**Claim:**
Konsumen pesan dan peserta saga harus idempotent. Mekanisme: tabel PROCESSED_MESSAGES keyed (subscriberId, messageID); atau event UUID di header Kafka; atau idempotency key (clientId) di API pihak ketiga (Stripe-style).

**Evidence:**
Richardson Idempotent Consumer: INSERT gagal jika PK duplikat → rollback, ignore. Debezium: eventId header. Temporal: "you, the programmer, need to make sure each Temporal Activity is idempotent."

**Sources:**
- Microservices.io Idempotent Consumer — https://microservices.io/patterns/data/idempotent-consumer.html
- Debezium Blog 2019-02-19
- Temporal Blog 2023-05-24 — https://temporal.io/blog/saga-pattern-made-easy

**Confidence:** HIGH

---

### Finding 8: Anomali Data karena Tidak Ada Isolation

**Claim:**
Tanpa isolation lintas service, saga concurrent bisa menghasilkan lost updates, dirty reads, fuzzy/nonrepeatable reads. Countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.

**Evidence:**
Microsoft enumerasi 3 anomali + 6 countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency. Daftar 6 item spesifik ini berasal dari dokumentasi Microsoft. Richardson: "saga developer must typically use countermeasures... careful analysis is needed" — halaman publik Microservices.io mengonfirmasi keberadaan konsep countermeasures tanpa meng-enumerasi daftar lengkapnya.

**Sources:**
- Microsoft Azure Architecture Center
- Microservices.io Pattern: Saga (chapter 4/section 4.3 reference)

**Confidence:** HIGH untuk adanya anomali dan konsep countermeasures; MEDIUM untuk enumerasi 6 countermeasures spesifik (satu sumber otoritatif: Microsoft Azure Architecture Center)

---

### Finding 9: Kompensasi Bisa Gagal — Tidak Ada Rollback Otomatis Tingkat Kedua

**Claim:**
Compensating transaction tidak dijamin sukses. Sistem bisa inkonsisten. Monitoring, retry, dan intervensi operator diperlukan.

**Evidence:**
Microsoft: "Compensating transactions might not always succeed, which can leave the system in an inconsistent state." AWS: "The saga pattern is difficult to debug and its complexity increases with the number of microservices." Temporal sample: log error pada compensation failure, lanjutkan sisa kompensasi. Tidak ada protokol recovery standar untuk kompensasi yang gagal; operasi produksi harus menambahkan retry dengan backoff, dead-letter queue, dan logging untuk intervensi manual.

**Sources:**
- Microsoft Azure Architecture Center
- AWS Prescriptive Guidance
- Temporal compensating sample (linked from 2023-05-24 post)

**Confidence:** HIGH

---

### Finding 10: Implementasi Nyata — Step Functions, Temporal, Eventuate

**Claim:**
AWS Step Functions: state machine orchestration, Catch/Retry, Standard = exactly-once workflow, Express = at-least-once. Temporal: deterministic replay, activity retry, compensation via Saga helper class; menghindari orchestrator SPOF. Eventuate Tram: framework Richardson untuk orchestration-based sagas.

**Evidence:**
AWS docs: Standard 2000 exec/s, exactly-once, up to 1 year; Express 100k exec/s, at-least-once, 5 min. Temporal: "by running your code with Temporal, you automatically get your state saved and retries on failure." Richardson: Eventuate Tram Sagas examples on GitHub.

**Sources:**
- AWS Step Functions Welcome — https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html
- Temporal Blog 2023-05-24
- Microservices.io example code links (eventuate-tram-sagas)

**Confidence:** HIGH

---

### Finding 11: Lab Exercise Mapping — Orchestration Checkout Failure di Inventory

**Claim:**
Skenario lab (Order Pending → Payment Debit → Inventory Reserve gagal → Refund Payment → Cancel Order) adalah orchestration saga klasik dengan kompensasi LIFO.

**Evidence:**
Richardson orchestration example: Create Order PENDING → Reserve Credit → approve/reject. Microsoft: compensable then compensate on failure. Temporal Java: `saga.addCompensation` sebelum setiap step; pada catch, `saga.compensate()` LIFO. Mapping lab:
1. Create Order (Pending) — kompensasi: Cancel Order
2. Debit Payment — kompensasi: Refund Payment
3. Reserve Inventory — gagal (stok habis)
4. Orchestrator jalankan: Refund Payment, lalu Cancel Order

**Sources:**
- Microservices.io orchestration example
- Temporal Java Saga class example
- Microsoft compensable/pivot/retryable model

**Confidence:** HIGH

---

## Areas of Agreement

- Saga = sequence of local transactions + compensating transactions. (Semua sumber)
- Dua gaya: Choreography dan Orchestration. (Semua sumber)
- 2PC tidak cocok untuk microservices long-running / heterogeneous. (Richardson, Microsoft, AWS)
- Eventual consistency, bukan ACID isolation. (Semua sumber)
- Idempotency wajib. (Richardson, Debezium, Temporal)
- Dual-write diselesaikan Outbox/Event Sourcing. (Richardson, Debezium)
- Kompensasi bisa gagal; monitoring wajib. (Microsoft, AWS)
- Lab spec "Kapan memakai / jangan memakai" sesuai guidance AWS/Microsoft.

## Areas of Disagreement

Tidak ada kontradiksi material. Perbedaan: Temporal mendorong orchestration bahkan untuk workflow sederhana jika pakai durable execution engine; Microsoft/lab menempatkan choreography sebagai opsi valid untuk workflow kecil. Dual-write adalah masalah modern (bukan definisi 1987). Taksonomi pivot/retryable hanya rinci di Microsoft.

## Limitations

- Paper 1987 adalah PDF scan; kutipan verbatim halaman tidak diekstrak (encoding LZW). Klaim historis diverifikasi via citation chain (DOI, Temporal footnote, ACM).
- Camunda, Axon, Seata, Spring State Machine tidak dibuka sebagai sumber primer dalam sesi ini.
- Tidak ada benchmark kuantitatif latency/throughput saga vs 2PC yang diverifikasi dari dataset primer.
- Recovery protocol untuk compensation-of-compensation tidak distandardisasi di sumber yang dibuka.

## Conclusion

Lab specification akurat. Saga Pattern adalah solusi konsistensi terdistribusi tanpa 2PC: pecah jadi transaksi lokal, kompensasi semantik jika gagal. Orchestration cocok untuk latihan checkout e-commerce 3 langkah karena alur kompensasi eksplisit dan mudah dilacak. Prasyarat produksi yang lab harus sebutkan: (1) Transactional Outbox atau Event Sourcing untuk dual-write, (2) idempotency key pada Payment debit/refund dan Inventory reserve/release, (3) handling jika kompensasi sendiri gagal (retry + alert), (4) isolation anomalies jika saga concurrent. Jangan klaim saga memberikan ACID; klaim yang benar: eventual consistency dengan kompensasi terarah.
