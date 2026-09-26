# Content Brief

**Topic:** Optimistic vs Pessimistic Locking — Mencegah Lost Update pada Aplikasi Concurrent

**Target Reader:** Backend engineers, database engineers, software architects yang membangun sistem dengan konkurensi tinggi pada resource terbatas (stok, saldo, reservasi).

**Problem:** Dua transaksi konkurensi membaca nilai yang sama, memodifikasinya secara independen, dan menulis kembali secara berurutan menyebabkan penimpaan data diam-diam (lost update) tanpa error — bahkan di isolation level default (READ COMMITTED) di PostgreSQL, MySQL, Oracle.

**Core Mental Model:**
- Lost update ≠ dirty read ≠ phantom read. Ini adalah penimpaan *committed* write oleh write kedua tanpa deteksi.
- Pessimistic locking = **prevent conflict early** (blokir writer lain via `SELECT ... FOR UPDATE`).
- Optimistic locking = **detect conflict late** (version guard di `WHERE` clause, cek `affected_rows == 0`).
- Atomic single-statement = **eliminate window entirely** (`UPDATE ... SET stock = stock - N WHERE stock >= N`).

**Approved Research Status:** APPROVED (research-audit/07-verdict.md)

**Approved Engineering Status:** APPROVED (engineering-audit/06-verdict.md)

**Main Concepts:**
1. Lost Update Anomaly — reproduktibel di default isolation tanpa explicit locking
2. Pessimistic Locking (`SELECT ... FOR UPDATE`) — row-level exclusive lock hingga commit/rollback
3. Optimistic Locking (version/timestamp + affected_rows check) — deteksi konflik saat commit
4. Atomic Conditional Update — statement-level atomicity tanpa lock manual
5. Selection Criteria — atomic first, optimistic untuk low-contention/read-heavy, pessimistic untuk high-contention/correctness-critical
6. Anti-patterns — transaction alone ≠ protection; holding lock across network call; ignoring 0-rows-affected

**Verified Behaviors (dari tests & demo):**
- Naive read-modify-write: 50 deduct calls → final stock 99 (expected 50) — **lost update terbukti**
- Pessimistic locking: 50 goroutines → final stock 50 (exact invariant) — **fully synchronized**
- Optimistic locking (no retry): 20 goroutines → 1 success, 19 rejected conflicts, stock 99 — **zero corruption, state guarded**
- Optimistic locking (with exponential backoff retry): 20 goroutines → 20 success, ~61 retries, stock 80 — **converged successfully**
- Atomic conditional update: 50 goroutines → final stock 50 — **lockless single statement works**
- Race detector: **zero race conditions** di semua test
- All tests pass, build passes, demo runs successfully

**Available Case Studies:**
- Lab demo in-memory simulation (5 skenario side-by-side)
- Automated concurrency tests (6 test cases dengan invariant assertions)
- Real database behavior documented via PostgreSQL, MySQL, Oracle official docs

**Warnings:**
- Lab menggunakan in-memory store (simulasi), bukan database nyata — perilaku deadlock, lock escalation, network partition tidak ditampilkan
- Atomic decrement recipe (`SET stock = stock - N WHERE stock >= N`) disintesis dari jaminan atomicitas statement ACID, bukan kutipan vendor verbatim (confidence MEDIUM di research)
- MySQL docs diakses via Oracle CDN mirror — claims MySQL-specific downgrade ke MEDIUM
- No performance benchmarks (throughput/latency) — hanya qualitative trade-offs
- Version overflow & Redis/distributed lock boundary conditions = open questions