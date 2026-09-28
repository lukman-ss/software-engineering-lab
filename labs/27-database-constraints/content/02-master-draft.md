# Database Constraints: Jangan Serahkan Integritas Data Hanya ke Application Code

## Problem

Validasi di application layer umumnya ditulis dengan pola check-then-act:

```text
SELECT → cek apakah email sudah ada → jika belum, INSERT
```

Pola ini terlihat benar pada pengujian sekuensial. Di bawah beban konkuren, pola ini gagal. Dua request dapat sama-sama membaca "belum ada" sebelum salah satunya menulis, sehingga keduanya melakukan INSERT dan menghasilkan data duplikat.

Lab ini membuktikan kegagalan tersebut secara konkret. Implementasi `UnsafeStore` melakukan pengecekan di memori aplikasi tanpa constraint database. Pada uji konkuren 20 goroutine dengan email yang sama, `UnsafeStore` menghasilkan lebih dari 1 baris (`TestConcurrentRegistration_Unsafe_SuffersRaceCondition`). Itu bukan anomali teoritis. Itu perilaku yang direproduksi oleh test.

## Why This Matters

Duplikat data bukan sekadar masalah tampilan. Sekali data duplikat masuk, perbaikannya mahal: deduplikasi, rekonsiliasi transaksi, konflik dengan sistem downstream. Pencegahan di titik tulis jauh lebih murah daripada perbaikan setelah fakta.

Database constraint memindahkan pengecekan ke titik yang atomik: saat baris ditulis. Tidak ada jeda antara "cek" dan "tulis" yang bisa diserobot transaksi lain.

Riset disetujui membuktikan bahwa constraints menggunakan B-tree index yang diperbarui atomik dengan `ROW EXCLUSIVE` lock pada saat INSERT/UPDATE [PostgreSQL DDL Constraints 5.5, Explicit Locking 13.3]. PostgreSQL menggambarkan ini: "menambah unique constraint akan secara otomatis membuat unique B-tree index pada kolom atau grup kolom yang digunakan dalam constraint". Ini artinya PostgreSQL tidak membiarkan ada jeda antara pengecekan keunikan dan penulisan baris.

Selain itu, dokumen SQLite mengonfirmasi kedua database memverifikasi NOT NULL dan CHECK hanya pada saat INSERT/UPDATE, bukan pada saat SELECT. Ini membuat constraints menjadi defense-in-depth yang konsisten di lapisan penyimpanan.

## Mental Model

Pikirkan dalam dua lapis:

```text
Application validation → UX lapis pertama (pesan ramah)
Database constraint   → kebenaran lapis terakhir (penolakan absolut)
```

Validasi aplikasi memberi pesan seperti "Email sudah digunakan". Constraint database menjamin "tidak mungkin ada dua email aktif yang sama", bahkan jika semua validasi aplikasi dilewati atau berlomba.

Riset yang disetujui merangkumnya sebagai defense-in-depth: aplikasi untuk pengalaman pengguna, database untuk proteksi absolut.

## Core Concept

Lab ini mengimplementasikan lima mekanisme yang bekerjasama untuk mencegah berbagai jenis ketidakakuratan data:

1. **NOT NULL (`23502`)**: kolom wajib seperti `email`, `username`, `user_id` tidak boleh kosong.
2. **CHECK (`23514`)**: predikat boolean pada baris, contoh lab: `age >= 18`, `status IN (...)`, `total_cents > 0`.
3. **UNIQUE (`23505`)**: satu kemunculan per nilai, mencegah race condition read-then-write.
4. **FOREIGN KEY (`23503`)**: referential integrity, mencegah orphan row.
5. **PARTIAL UNIQUE INDEX**: uniqueness kondisional (`WHERE deleted_at IS NULL`), memungkinkan registrasi ulang email setelah soft delete sambil menjaga uniqueness untuk data aktif.

Dua kelompok saling melengkapi: NOT NULL dan CHECK menjamin integritas kolom per baris, sedangkan UNIQUE, FOREIGN KEY, dan partial index menjamin integritas lintas baris. Pelanggaran masing-masing dipetakan ke SQLSTATE standar PostgreSQL Class 23 yang stabil dan tidak terlokalisasi. Aplikasi memeriksa kode (mis. `23505`) bukan parsing teks error untuk mencegah bug akibat perubahan lokal.

Sebagaimana dokumentasi PostgreSQL 18 menyatakan: "Unique constraint memastikan bahwa tidak ada dua baris dengan nilai indeks yang sama boleh muncul dalam tabel". Ketika kombinasi dengan indeks unique, database memanfaatkan mekanisme lock internal untuk melindungi pengecekan dari gangguan konkuren.

## Failure Scenario

Skenario gagal yang dibuktikan lab adalah registrasi konkuren.

`UnsafeStore.RegisterUser` membaca jumlah user, memindai email satu per satu, lalu `Sleep(1ms)` untuk mensimulasikan jeda scheduling, kemudian insert tanpa constraint:

```text
Goroutine A: scan → tidak ketemu → sleep → INSERT sukses
Goroutine B: scan → tidak ketemu → sleep → INSERT sukses (DUPLIKAT)
```

Test membuktikan `GetUsersCount() > 1` setelah 20 goroutine konkuren memakai email sama.

Skenario aman memakai `SafeStore.RegisterUser` yang mendelegasikan langsung ke engine. Hasil yang dibuktikan test `TestConcurrentRegistration_Safe_EnforcesUniqueness` (20 goroutine) dan demo `cmd/demo` (50 goroutine):

```text
Total Goroutine: 50
Successful: 1
Rejected 23505: 49
Integrity intact: true
```

## How It Works

Engine lab adalah simulator relasional in-memory berbasis Go standard library (`sync.RWMutex`, `atomic.Int64`, map). Ia bukan PostgreSQL live. Ia memodelkan semantik PostgreSQL: pengecekan deklaratif, error SQLSTATE `23xxx`, dan perilaku locking konkuren. Implementasi memakai coarse-grained table lock, bukan B-tree page latch — trade-off yang didokumentasikan demi readability dan zero-dependency.

`Engine.InsertUser` menahan `e.mu.Lock()` untuk seluruh operasi insert. Di dalam satu critical section atomik itu, ia mengevaluasi berurutan:

1. NOT NULL
2. CHECK
3. UNIQUE / PARTIAL UNIQUE
4. Pembangkitan primary key dan commit baris + index

Karena cek dan tulis terjadi di bawah lock yang sama, tidak ada interleaving yang memungkinkan dua insert lolos dengan kunci sama. Ini adalah analogi dari `ROW EXCLUSIVE` lock dan update index B-tree atomik pada PostgreSQL yang dibahas di riset. Implementasi lab memakai coarse-grained table lock, bukan B-tree page latch — trade-off yang didokumentasikan demi readability dan zero-dependency.

`Engine.InsertOrder` mengikuti urutan serupa: NOT NULL pada `user_id`, CHECK pada `total_cents > 0`, lalu FOREIGN KEY dengan memeriksa keberadaan `user_id` di map `users`.

## Architecture

Komponen aktual dalam lab:

```text
cmd/demo/main.go
  → internal/store/store.go (UnsafeStore, SafeStore)
    → internal/engine/engine.go (Engine, constraint checks)
      → internal/dberr/errors.go (SQLSTATE taxonomy, MapToDomainError)
      → internal/model/model.go (User, Order)
```

- `internal/model`: entitas domain `User` dan `Order`.
- `internal/engine`: table store thread-safe dengan validasi skema, index UNIQUE dan partial unique, dan referential integrity.
- `internal/dberr`: taksonomi error Class 23 dan mapper ke domain error.
- `internal/store`: `UnsafeStore` (pola rentan) vs `SafeStore` (delegasi ke constraint).
- `cmd/demo`: demonstrasi eksekusi yang membandingkan keduanya.

## Implementation

Aturan konkret yang diimplementasikan engine:

- `users.email == ""` → `23502` (`users_email_not_null`)
- `users.username == ""` → `23502` (`users_username_not_null`)
- `orders.user_id == 0` → `23502` (`orders_user_id_not_null`)
- `users.age < 18` → `23514` (`users_age_check`)
- `users.status` di luar `active/suspended/pending` → `23514` (`users_status_check`)
- `orders.total_cents <= 0` → `23514` (`orders_total_cents_check`)
- Duplikat `email` pada mode full UNIQUE → `23505` (`users_email_key`)
- Duplikat `email` aktif pada mode partial index → `23505` (`users_active_email_idx`)
- `orders.user_id` tidak ada di `users` → `23503` (`fk_orders_user`)

Nilai seperti `age >= 18` dan daftar status adalah contoh lab, bukan rekomendasi universal. Jangan dibaca sebagai aturan bisnis umum.

Engine memodelkan PostgreSQL `ROW EXCLUSIVE` lock dan atomik B-tree index update dengan coarse-grained mutex lock. PostgreSQL mengacu pada urutan pemeriksaan: NOT NULL sebelum CHECK, CHECK diurutkan secara alfabetis menurut nama constraint. Engine lab mengikuti urutan logis serupa: NOT NULL → CHECK → UNIQUE → FK.

Engine lab tidak mendukung DEFERRABLE constraints, EXCLUDE constraints dengan GiST operator, atau migrasi `NOT VALID + VALIDATE CONSTRAINT`. Topik-topik ini ada di riset sebagai konteks, bukan sebagai perilaku terverifikasi implementasi ini.

Partial unique index dimodelkan dengan dua map: `emailIndex` untuk UNIQUE penuh dan `activeEmails` untuk index kondisional. Baris hanya masuk `activeEmails` jika `DeletedAt == nil`. `SoftDeleteUser` menghapus email dari `activeEmails` dan menandai `DeletedAt`, sehingga email tersebut dapat dipakai ulang oleh baris aktif baru.

Error mapping memisahkan kode storage dari pesan user:

```text
23505 → conflict: resource with this unique attribute already exists
23502 → invalid input: mandatory field is missing
23514 → validation failed: value outside permissible boundary
23503 → reference error: referenced entity does not exist
```

Setiap pesan menyertakan nama rule/constraint agar dapat dipetakan ke respons domain seperti HTTP `409 Conflict`.

## Code Walkthrough

Alur tulis aman (`SafeStore.RegisterUser`): terima `User`, panggil `eng.InsertUser(u, false)`, jika error petakan via `dberr.MapToDomainError`. Tidak ada SELECT pendahulu. Constraint yang memutuskan.

Alur tulis rentan (`UnsafeStore.RegisterUser`): pindai semua user via `GetUser(i)`, beri jeda `Sleep`, lalu panggil `InsertUserUnsafe` yang melewati pengecekan index. Jeda itu mensimulasikan I/O atau scheduling delay yang memperlebar race window.

Alur partial index (`SafeStore.RegisterUserPartial`): panggil `eng.InsertUser(u, true)`. Jika `DeletedAt == nil` dan email sudah ada di `activeEmails`, tolak. Jika baris membawa `DeletedAt`, lewati pengecekan — ini yang memungkinkan multiple soft-deleted duplikat.

Alur order (`SafeStore.CreateOrder`): panggil `eng.InsertOrder(o)`. Engine menolak `UserID == 0`, `TotalCents <= 0`, dan `UserID` yang tidak ada sebelum menyimpan.

Detail kode lengkap ada di `03-code-snippets.md` dengan referensi file sumber.

## What the Tests Prove

Test suite (`internal/store/store_test.go`) lolos penuh termasuk dengan `-race`. Yang dibuktikan, tidak lebih:

- `TestNotNullConstraints`: email kosong, username kosong, dan `user_id` order kosong ditolak.
- `TestCheckConstraints`: umur 16 ditolak, status `banned` ditolak, `total_cents = 0` ditolak setelah user valid dibuat.
- `TestUniqueConstraint`: registrasi pertama sukses, duplikat email ditolak.
- `TestForeignKeyConstraint`: order dengan `UserID=99999` ditolak; order dengan user valid sukses dan mendapat ID.
- `TestPartialUniqueIndex`: duplikat aktif ditolak; setelah soft delete, email sama dapat dipakai ulang; duplikat aktif kedua ditolak lagi; insert langsung dengan `DeletedAt` terisi melewati partial index.
- `TestConcurrentRegistration_Safe_EnforcesUniqueness`: 20 goroutine, tepat 1 sukses, 19 error, total count 1.
- `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`: 20 goroutine tanpa constraint menghasilkan count > 1.
- `TestErrorClassification`: `IsConstraintViolation` benar untuk keempat kode `23502/23505/23514/23503`.

Demo (`go run ./cmd/demo`) memverifikasi perilaku yang sama pada skala 50 goroutine: 1 sukses, 49 ditolak `23505`.

Yang tidak dibuktikan lab: constraint terdistribusi lintas node, parser SQL penuh, EXCLUDE dengan operator GiST, migrasi `NOT VALID + VALIDATE CONSTRAINT` pada tabel produksi, isolasi `SERIALIZABLE`, dan limitasi tabel partisi. Topik-topik itu ada di riset sebagai konteks, bukan sebagai perilaku terverifikasi implementasi ini.

## Recovery / Rollback

Lab ini tidak mengimplementasikan strategi migrasi produksi. Tidak ada kode `NOT VALID`, `VALIDATE CONSTRAINT`, atau rollback skema.

Untuk konteks riset saja: dokumentasi PostgreSQL yang dirujuk menyarankan `ADD CONSTRAINT ... NOT VALID` lalu `VALIDATE CONSTRAINT` untuk menghindari scan awal dan mengurangi locking pada tabel produksi. Pola itu tidak didemonstrasikan oleh engine in-memory ini. Jangan menganggapnya terverifikasi oleh test lab.

## Production Considerations

Beberapa batasan yang harus dipertahankan saat membawa pelajaran lab ke produksi:

- **Batas row-scoped.** CHECK hanya melihat kolom baris saat ini. Aturan lintas baris atau lintas tabel (contoh riset: total transaksi 12 bulan > Rp50 juta) tidak bisa diekspresikan sebagai CHECK. Butuh trigger atau isolasi `SERIALIZABLE` dengan retry pada `40001`.
- **Foreign key tidak otomatis membuat index** pada kolom referencing (temuan riset dari PostgreSQL 5.5.5). Index manual pada kolom referencing disarankan untuk performa DELETE/UPDATE pada tabel parent.
- **NULL handling.** PostgreSQL dan SQLite menganggap NULL distinct pada UNIQUE secara default (multiple NULL diizinkan). Perilaku multi-kolom `(1, NULL)` vs `(1, NULL)` tidak diuji secara empiris di riset dan tidak dicakup test lab. Verifikasi pada database target sebelum mengandalkan.
- **Partial index predicate matching.** Predikat index harus cocok dengan klausa query; planner tidak mengenali ekspresi yang ekuivalen secara matematis tetapi ditulis berbeda, dan klausa parameterized tidak bekerja dengan partial index (temuan riset 11.8). Gunakan `EXPLAIN` untuk verifikasi.
- **MySQL parity belum terverifikasi.** Dokumentasi MySQL tidak dapat diakses saat riset (HTTP 403). Jangan mengasumsikan perilaku identik lintas database.
- **Engine lab adalah simulator.** Lock coarse-grained, bukan B-tree page latch. Cocok untuk pembelajaran deterministik tanpa Docker, bukan untuk benchmark produksi.

## Common Mistakes

1. **Hanya validasi di aplikasi.** Pola `exists() → INSERT` gagal di bawah konkurensi. Test lab membuktikan duplikat.
2. **Menebak dari teks error.** Teks bisa terlokalisasi. Periksa SQLSTATE (`23505`, bukan string "duplicate key").
3. **Mengekspos error mentah sebagai 500.** Petakan `23505` ke `409 Conflict` dengan pesan bisnis.
4. **UNIQUE penuh untuk soft delete.** Tanpa partial index, email yang sudah soft-deleted tidak bisa dipakai ulang. Gunakan `WHERE deleted_at IS NULL`.
5. **CHECK untuk aturan lintas tabel.** CHECK tidak melihat baris lain. Paksa aturan seperti itu via CHECK akan tampak lolos pada test sederhana tetapi gagal di bawah beban konkuren.
6. **Lupa index pada foreign key.** Integritas terjaga, tetapi DELETE/UPDATE parent menjadi lambat karena full scan referencing table.

## Case Study

Dua studi kasus yang benar-benar dijalankan lab:

**Registrasi konkuren.** 50 goroutine mendaftar dengan email sama via `SafeStore`. Tepat 1 menang, 49 menerima domain error dari `23505`. Integritas utuh tanpa retry logic di aplikasi. Lawannya, `UnsafeStore` dengan 20 goroutine menghasilkan duplikat.

**Soft-delete email.** User aktif `alice@company.com` dibuat (ID=1). Duplikat aktif ditolak. Setelah `SoftDeleteUser`, email sama dapat dipakai ulang (ID=2). Insert langsung baris dengan `DeletedAt` terisi juga diizinkan karena predikat partial index tidak mengindeksnya.

## Checklist

- [ ] Kolom wajib ditandai NOT NULL, bukan hanya validasi form
- [ ] Nilai unik yang harus tunggal memakai UNIQUE, bukan `SELECT` pendahulu
- [ ] Pelanggaran dipetakan via SQLSTATE, bukan parsing teks error
- [ ] `23505` dipetakan ke respons idempoten/konflik yang ramah (mis. `409`)
- [ ] Soft-delete uniqueness memakai partial unique index
- [ ] Foreign key kolom referencing di-index untuk performa
- [ ] Aturan lintas baris tidak dipaksa menjadi CHECK
- [ ] Query partial index diverifikasi dengan `EXPLAIN`
- [ ] Perilaku NULL pada UNIQUE multivariat diverifikasi pada DB target
- [ ] Test konkurensi + `-race` lolos sebelum klaim aman-konkuren

## Key Takeaways

Lihat `05-key-takeaways.md` untuk 8 poin ringkas.

## Sources

- `research/02-sources.md`: PostgreSQL 18 docs (constraints 5.5, unique indexes 11.6, partial indexes 11.8, error codes Appendix A, app-level consistency 13.4, explicit locking 13.3, partitioning 5.12, ALTER TABLE), SQLite CREATE TABLE, Stripe idempotency.
- `research/05-report.md`: 11 findings, keterbatasan single-database focus dan locking internals.
- Implementasi: `internal/engine/engine.go`, `internal/dberr/errors.go`, `internal/store/store.go`, `internal/model/model.go`, `cmd/demo/main.go`.
- Test: `internal/store/store_test.go` (8 test, lolos + `-race`).
- Audit: `research-audit/07-verdict.md` (APPROVED), `engineering-audit/06-verdict.md` (APPROVED).

Pemetaan lengkap ada di `06-source-map.md`.