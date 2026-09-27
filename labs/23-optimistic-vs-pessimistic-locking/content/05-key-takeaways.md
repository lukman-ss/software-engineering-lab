# Key Takeaways

1. **Lost Update adalah korupsi diam-diam** — dua transaksi yang membaca nilai yang sama, menghitung ulang, dan menulis tanpa penguncian menyebabkan perubahan pertama hilang tanpa error.

2. **Pessimistic locking (SELECT ... FOR UPDATE)** mencegah konflik dengan memblokir baris lain. Setiap goroutine/thread yang berlawanan menunggu sampai lock dilepas. Sesuai untuk sumber daya kritis (balance, stok terbatas, seat reservation).

3. **Optimistic locking (version guard + affected_rows check)** mendeteksi konflik saat commit. Jika `affected_rows == 0`, kembalikan `ErrOptimisticLock`. Aplikasi harus retry atau return 409 Conflict. Sesuai untuk beban baca-banyak dengan konflik jarang (profil, CMS).

4. **Atomic single-statement update** (`UPDATE ... SET stock=stock-N WHERE stock>=N`) menghilangkan jendela read-modify-write tanpa lock manual. Pilihan terbaik untuk counter/sederhana.

5. **Isolation level tidak cukup** — READ COMMITTED (default PostgreSQL/Oracle) tetap memungkinkan lost update tanpa query design yang benar. Transaction semata tidak mencegahnya.

6. **Pilih strategi berdasarkan frekuensi konflik** — Pessimistic (konflik sering), Optimistic (konflik jarang), Atomic (sederhana). Tidak ada satu-solution-fits-all.

7. **Deadlock pada pessimistic** — PostgreSQL/Oracle mendeteksi otomatis dan abort satu transaksi. Solusi: lock order konsisten, transaksi pendek, jangan pegang lock selama HTTP call.

8. **Retry wajib pada optimistic** — Tanpa retry, konflik menghasilkan 409. Dengan exponential backoff (jitter), semua goroutine pada akhirnya konvergen.

9. **Race detector harus lulus** — `go test -race ./...` tidak boleh menghasilkan warning apapun. Lab ini pass tanpa peringatan.

10. **Verifikasi invariant** — FinalStock = InitialStock - TotalDeductions (atau accounting conflicts). Ini adalah bukti utama korektif.

11. **Simulasi dalam lab ≠ RDBMS nyata** — ini adalah pusat memori Go (`sync.Mutex`) yang meniru semantik `SELECT FOR UPDATE`, version guard, dan update atomik. MySQL 8.0 docs tidak diverifikasi langsung (403) — klaim MySQL didasarkan pada SQL standar.

12. **Lab melewatkan beberapa path error** — engineering audit mencatat: `ErrInvalidQuantity`, `ErrNotFound`, kegagalan retry optimis (maxRetries tercapai), dan overdraft konkuren tidak diuji. Perilaku ini ada pada kode, tetapi tidak divisualisasikan lewat tes otomatis.
