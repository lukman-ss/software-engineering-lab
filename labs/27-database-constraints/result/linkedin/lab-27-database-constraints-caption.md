Lab 27 kali ini membahas pencegahan race condition dan kerusakan integritas data melalui pemindahan validasi ke level database constraints di Go.

Masalah teknis yang sering terjadi adalah pola check-then-act di application layer (SELECT lalu INSERT). Di bawah konkurensi, dua goroutine dapat membaca bahwa email belum ada secara bersamaan, sehingga keduanya mengeksekusi INSERT dan menghasilkan record duplikat yang merusak database.

Lab ini membuktikan lima pilar integritas data konkret:
• NOT NULL Enforcement: Kolom wajib seperti email, username, dan user_id ditolak seketika pada level engine dengan SQLSTATE 23502.
• CHECK Boolean Predicate: Validasi batas nilai (age >= 18, enum status, total_cents > 0) dievaluasi atomik dengan SQLSTATE 23514.
• UNIQUE B-Tree Protection: Penulisan atomik di bawah lock menolak duplikasi konkuren, menghasilkan tepat 1 pemenang dan N-1 error SQLSTATE 23505.
• Referential Integrity: FOREIGN KEY mencegah orphan rows pada tabel orders dengan SQLSTATE 23503.
• Partial Unique Index: Predikat WHERE deleted_at IS NULL mengizinkan registrasi ulang email setelah soft delete tanpa mengorbankan keunikan aktif.

Mental model inti: Application validation is for user experience; database constraints are for absolute correctness. Never rely on application checks alone.

Satu batasan penting: Lab ini menggunakan simulator in-memory engine Go dengan sync.RWMutex untuk membuktikan konsep relasional tanpa dependensi Docker live.

Pelajari implementasi lengkap dan seluruh pembuktian test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints

#softwareengineering #golang #databases #systemdesign #dataintegrity #postgresql
