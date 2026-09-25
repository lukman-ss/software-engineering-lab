## Claim 1

Claim: Deadlock terjadi ketika dua transaksi saling menahan lock yang dibutuhkan satu sama lain, membentuk siklus tunggu permanen (circular wait).

Location: `03-evidence.md` (Evidence 1), `05-report.md` (Finding 1)

Evidence Provided: "A deadlock occurs when two or more tasks permanently block each other by each task having a lock on a resource that the other tasks are trying to lock."

Source: Deadlocks Guide - SQL Server; PostgreSQL 18 Documentation: 13.3. Explicit Locking

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Definisi standar deadlock dalam RDBMS yang didukung penuh oleh dokumentasi PostgreSQL dan SQL Server.

## Claim 2

Claim: Database memiliki deadlock monitor yang membatalkan (abort/rollback) salah satu transaksi sebagai deadlock victim secara periodik/otomatis.

Location: `03-evidence.md` (Evidence 2), `05-report.md` (Finding 2)

Evidence Provided: "The Database Engine deadlock monitor periodically checks for tasks that are in a deadlock. If the monitor detects a cyclic dependency, it chooses one of the tasks as a victim and terminates its transaction with an error."

Source: Deadlocks Guide - SQL Server; PostgreSQL 18 Documentation: 13.3. Explicit Locking

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Didukung dokumentasi resmi kedua DBMS. Keduanya menerapkan pemilihan victim untuk memecahkan siklus.

## Claim 3

Claim: Urutan akses yang konsisten (Lock Ordering) mencegah mayoritas deadlock.

Location: `03-evidence.md` (Evidence 3), `05-report.md` (Finding 3)

Evidence Provided: "The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order."

Source: Deadlocks Guide - SQL Server; PostgreSQL 18 Documentation: 13.3. Explicit Locking

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Kedua vendor menyebut konsistensi lock ordering sebagai mitigasi primer.

## Claim 4

Claim: Durasi transaksi yang panjang (termasuk interaksi eksternal) memperbesar probabilitas deadlock.

Location: `03-evidence.md` (Evidence 3), `05-report.md` (Finding 4)

Evidence Provided: "Keep transactions short and in one batch... The longer the transaction, the longer the exclusive or update locks are held, blocking other activity and leading to possible deadlock situations."

Source: Deadlocks Guide - SQL Server

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Didukung dokumentasi SQL Server dan PostgreSQL.

## Claim 5

Claim: Aplikasi sebaiknya menerapkan mekanisme Retry otomatis saat menangkap error deadlock (misal error 1205).

Location: `03-evidence.md` (Evidence 4), `05-report.md` (Finding 5)

Evidence Provided: "Implementing an error handler that catches error 1205 allows an application to handle deadlocks and take remedial action (for example, automatically resubmitting the query that was involved in the deadlock)."

Source: Deadlocks Guide - SQL Server; PostgreSQL 18 Documentation: 13.3. Explicit Locking

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Didukung kedua referensi vendor.

## Claim 6

Claim: Pemeriksaan deadlock memakan waktu (overhead tinggi), sehingga database biasanya menunggu (deadlock timeout) sebelum melakukan pengecekan.

Location: `03-evidence.md` (Evidence 5)

Evidence Provided: "The check for deadlock is relatively expensive, so the server doesn't run it every time it waits for a lock."

Source: PostgreSQL Lock Management (`deadlock_timeout`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Didukung dokumentasi resmi PostgreSQL.
