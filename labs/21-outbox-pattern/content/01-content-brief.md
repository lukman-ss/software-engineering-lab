# Content Brief

Topic: Transactional Outbox Pattern — menyelesaikan dual-write problem antara database dan message broker

Target Reader: Backend engineer dan software architect yang membangun microservices dengan database + message broker (Kafka / RabbitMQ). Diasumsikan paham transaksi database dan event-driven dasar.

Problem: Update database lalu publish message ke broker (atau sebaliknya) tidak atomik. Transaksi database hanya mengontrol operasi database, tidak bisa ROLLBACK Redis, RabbitMQ, Kafka, atau HTTP request. Jika DB commit sukses tetapi publish gagal, sistem inkonsisten: order tersimpan, event hilang.

Core Mental Model: Pindahkan write kedua (publish message) ke dalam boundary transaksi yang sama dengan write pertama. Service menyimpan message ke tabel outbox dalam transaksi yang sama dengan update business entity. Proses terpisah (message relay) membaca tabel outbox secara async dan publish ke broker. Hasil: DB commit = event pasti terkirim (eventual, at-least-once). Relay boleh gagal dan retry karena pesan aman di outbox.

Approved Research Status: APPROVED (research-audit/07-verdict.md, 2026-09-27, 0 unsupported claims)

Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md, 2026-09-27, build + 8 tests + race detector + demo PASS)

Main Concepts:
1. Dual-write problem dan kenapa 2PC tidak viable
2. Outbox table sebagai bagian dari transaksi bisnis
3. Message relay: Polling Publisher (diimplementasikan lab ini) vs Transaction Log Tailing / CDC
4. At-least-once delivery dan idempotent consumer via event ID tracking
5. Outbox table design: kolom id, aggregatetype, aggregateid, type, payload
6. Operasional: cleanup processed events, monitoring unprocessed count dan oldest-event age

Verified Behaviors:
1. CreateOrderWithOutbox menyimpan Order + OutboxMessage atomik dalam satu Tx; commit bersamaan
2. Rollback membuang staged order dan outbox; tidak ada yang terpersist, tidak ada pesan ke broker
3. Relay polling mengambil PENDING, Publish ke broker, tandai PROCESSED; gagal publish = tetap PENDING untuk retry berikutnya
4. Dual-write naive terbukti inkonsisten: order ada di DB, broker kosong saat broker down
5. Consumer.Handle idempotent: delivery pertama accepted=true, duplikat accepted=false, received count tetap 1
6. Concurrent: 10 worker x 10 order = 100 pesan terpublish, 0 pending; race detector bersih
7. Relay retry: gagal pertama karena SetFailNext, poll berikutnya sukses dan tandai PROCESSED
8. Cleanup: PurgeProcessedOutbox menghapus yang PROCESSED, mempertahankan PENDING

Available Case Studies:
1. Demo Scenario 1 (cmd/demo/main.go): dual-write naive dengan broker down — Order in DB = true, Broker Message Count = 0
2. Demo Scenario 2: outbox atomik — 1 pesan OrderCreated, payload JSON Order, consumer accepted=true
3. Demo Scenario 3: duplikat delivery — accepted=false, total processed tetap 1

Warnings:
1. Lab ini simplified: DB adalah in-memory transactional map (BeginTx/Commit/Rollback), bukan Postgres/MySQL; tidak ada partisi throughput tinggi
2. Relay yang diimplementasikan hanya Polling Publisher; Transaction Log Tailing (Debezium / WAL) tidak didemonstrasikan
3. Tidak ada dead-letter queue untuk poison payload non-retryable
4. Nilai SLA contoh (oldest-event age normal ~2 detik, abnormal 47+ menit) bersifat ilustratif dari spesifikasi lab, bukan rekomendasi universal — harus di-tuning ke SLO masing-masing
5. Hindari payload raksasa di outbox; lab memakai JSON Order penuh sebagai contoh, bukan anjuran menyimpan 500 field
