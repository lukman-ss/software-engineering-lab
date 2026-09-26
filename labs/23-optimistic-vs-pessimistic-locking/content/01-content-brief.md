# Content Brief

Topic:
Optimistic vs Pessimistic Locking & Atomic Updates — mencegah anomali "lost update" pada read-modify-write konkuren.

Target Reader:
Pengembang perangkat lunak, insinyur basis data, dan arsitek sistem yang sudah akrab dengan konsep transaksi dan concinnity dasar, ingin memahami kapan dan bagaimana menerapkan strategi penguncian pada sumber daya yang diakses oleh banyak klien secara bersamaan.

Problem:
Dua proses (atau goroutine) yang membaca nilai yang sama, menghitung nilai baru secara independen, lalu menuliskan kembali secara berurutan menyebabkan perubahan dari transaksi pertama disamarkan — "lost update". Anomali ini terjadi secara diam-diam tanpa error, bahkan di bawah isolasi default (READ COMMITTED) pada PostgreSQL, Oracle, dan MySQL/InnoDB. Transaksi semata tidak mencegahnya.

Core Mental Model:
Locking = strategi konkurensi. Pessimistic = cegah konflik dengan memblokir (SELECT ... FOR UPDATE). Optimistic = deteksi konflik saat commit (version guard + affected_rows = 0). Atomic = hilangkan jeda read-modify-write lewat satu pernyataan UPDATE bersyarat. Setiap strategi adalah trade-off antara konsistensi, concinnity, dan kompleksitas operasional. Pilih berdasarkan frekuensi konflik dan kritisnya data.

Approved Research Status:
APPROVED

Approved Engineering Status:
APPROVED

Main Concepts:
1. Lost update — definisi, mekanisme, dan reproduksi under READ COMMITTED
2. Pessimistic locking — SELECT ... FOR UPDATE, row-level TX lock, held hingga commit/rollback
3. Optimistic locking — version/timestamp guard di WHERE clause, affected_rows = 0 sebagai sinyal konflik, perlu retry atau 409 Conflict
4. Atomic single-statement — UPDATE ... SET stock = stock - N WHERE stock >= N, statement-level atomicity
5. Isolation level limitation — hanya query design (lock, version guard, atau atomic statement) yang mencegah lost update, bukan isolasi semata
6. Anti-patterns — transaksi semata, lock via jaringan/HTTP, ignore 0-rows-affected, naive read-modify-write
7. Deadlock — auto-detected di PostgreSQL, satu transaksi abort secara tidak pasti; solusi konsisten lock order & transaksi pendek

Verified Behaviors:
- 50 goroutine melakukan DeductNaive(1, 1) pada stok awal 100 → stok akhir 99 (lost update terjadi; bukan 50)
- 50 goroutine melakukan DeductPessimistic(2, 1) pada stok awal 100 → stok akhir 50 (invariant terjaga)
- 20 goroutine melakukan DeductOptimisticDirect(3, 1) pada stok awal 100 → 1 sukses, 19 konflik (affected_rows = 0 setara)
- 20 goroutine dengan retry (max 10, jittered exponential backoff) → semua 20 konvergen, stok akhir 80
- 50 goroutine melakukan DeductAtomic(5, 1) pada stok awal 100 → stok akhir 50 (lockless, statement-atomic)
- `go test -race ./...` lulus tanpa peringatan race
- Stock insufficiency dengan pessimistic → ErrInsufficientStock yang benar (Seed 2, DeductPessimistic(20, 2) sukses, lalu DeductPessimistic(20, 1) gagal)

Available Case Studies:
1. Demo CLI (`go run ./cmd/demo`): lima skenario berdampingan — Naive, Pessimistic, Optimistic Direct, Optimistic With Retry, Atomic — ditampilkan side-by-side dengan stok awal 100
2. Test suite `tests/locking_test.go`: lima tes otomatis (NaiveLostUpdate, PessimisticLocking, PessimisticLockingInsufficientStock, OptimisticLockingConflict, OptimisticLockingWithRetry, AtomicConditionalUpdate)

Warnings:
- Lab ini adalah simulasi in-memory dengan `sync.Mutex` (bukan RDBMS riil). Perilaku row lock dimodelkan, bukan diprasyaratkan oleh database.
- Artificial micro-delay (`time.Sleep`) pada NaiveDeduct (100µs) dan OptimisticDeduct (50µs) disisipkan untuk membuat lost update / conflict dapat direproduksi secara andal.
- Deadlock detection disimulasi tidak mewujudkan karena ordering lock tunggal (single-resource).
- MySQL docs diverifikasi via Oracle CDN mirror — domain dev.mysql.com return 403 pada fetch otomatis.
- Versi optimisitis integer overflow tidak diuji; rekomendasi 64-bit/timestamp tersedia di research-revision.
- Distribusi concurrency di antara goroutine bersifat nondeterministik; nilai stok akhir pada skenario naive/optimistic dapat bervariasi antar-run.
- Nilai stock akhir 99 pada demo naive bersifat illustrative — jumlah presisi tergantung timing scheduler, tapi anomali lost update pada prinsipnya konsisten.
