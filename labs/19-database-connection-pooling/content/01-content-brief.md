# Content Brief

Topic: Database Connection Pooling — Overhead, Exhaustion, Leak dan Pool Sizing

Target Reader: Backend engineer, platform engineer, dan SRE yang mengelola aplikasi dengan beban konkuren terhadap PostgreSQL (atau RDBMS process-per-connection) dan perlu mencegah kegagalan koneksi di production.

Problem: Aplikasi mengalami error koneksi (`FATAL: sorry, too many clients already`) sementara CPU dan memory database terlihat sehat. Respons umum — memperbesar pool atau menaikkan `max_connections` — justru memperburuk latency dan throughput karena perebutan resource, atau menunda kegagalan yang sama pada skala deployment yang lebih besar.

Core Mental Model: Kapasitas koneksi ditentukan oleh hardware database dan limit server, bukan oleh jumlah goroutine/konkurensi aplikasi. Pool harus sekecil mungkin yang masih memenuhi throughput, dan koneksi harus dipegang sesingkat mungkin (jangan menahan koneksi selama I/O eksternal).

Approved Research Status: APPROVED (research-audit/07-verdict.md, 2026-09-26 — Source Integrity PASS, Claim Support PASS, Internal Consistency PASS; 2 non-blocking issues: SSD formula belum terverifikasi empiris, klaim 50x Oracle dari demonstrasi vendor)

Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md, 2026-09-26 — Compilation PASS, Tests 10/10 PASS, Race Detector PASS, Demo PASS, Research Alignment PASS; 1 recommended fix: go 1.26.7 adalah metadata versi non-eksisten)

Main Concepts: Connection handshake overhead dan reuse; `max_connections` sebagai slot keras; rumus sizing `((core_count * 2) + effective_spindle_count)`; performance knee (6 mekanisme degradasi); leak & pool starvation; global deployment sizing (`instance_count × pool_size`); PgBouncer pool modes (session/transaction/statement); pool-locking deadlock formula; monitoring pool-level metrics; HikariCP `leakDetectionThreshold`.

Verified Behaviors: Pooled lebih cepat dari unpooled (5 query, connectDelay 5ms) dan membuat lebih sedikit koneksi; oversized pool ditolak server (server max 10 / pool 20 → ≥10 gagal; demo server 15 / pool 50 → 15 dari 30 ditolak); satu koneksi yang ditahan selama external call memblokir request sehat hingga `context deadline exceeded`; double-close tidak double-decrement; error externalCall diteruskan dan koneksi dibersihkan; pool size 1 memblokir akuisisi koneksi kedua.

Available Case Studies: Demo overhead 54ms vs 6μs (angka ilustratif, bervariasi runtime); demo rejection 15/30; demo starvation dengan 2 unsafe order memblokir pool size 2; studi eksternal Oracle 2048→96 koneksi diklaim 50x lebih cepat (video vendor, confidence MEDIUM, bukan benchmark independen).

Warnings: Lab memakai mock driver berbasis memori (`int32` dengan operasi atomik via `sync/atomic`), bukan Postgres nyata — lock contention internal, context switching CPU, dan OOM per backend tidak disimulasikan; proxy logic PgBouncer tidak didemonstrasikan; formula sizing untuk SSD belum tervalidasi; contoh 3000 users/6000 TPS adalah ilustrasi hedged ("we'd wager"); angka demo adalah penalti handshake yang dipilih sendiri (5–10ms) dan bervariasi antar runtime; fokus riset PostgreSQL, bukan MySQL/SQL Server; go.mod `go 1.26.7` adalah metadata tidak akurat.
