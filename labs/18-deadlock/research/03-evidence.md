# Evidence

## Evidence 1 — Coffman Conditions (definisi deadlock)
**Claim:** Deadlock membutuhkan keempat syarat Coffman secara bersamaan: mutual exclusion,
hold-and-wait, no preemption, circular wait.

**Evidence:** "A deadlock situation on a resource can arise only if all of the following
conditions occur simultaneously in a system: 1. Mutual exclusion — only one process may
use a resource at a time. 2. Hold and wait — a process holds at least one resource while
requesting others. 3. No preemption — resources are released only voluntarily.
4. Circular wait — a cycle exists where P1 waits for P2, …, PN waits for P1."

**Source:** Coffman et al. 1971 / Silberschatz OS Concepts (Wikipedia summary)

**URL:** https://en.wikipedia.org/wiki/Deadlock_(computer_science)
**Published:** 1971 (original paper); 2006 (Silberschatz textbook)
**Confidence:** HIGH — multiple authoritative sources agree
**Corroborated By:** PostgreSQL documentation implicitly confirms all four in deadlock example
**Notes:** Syarat-syarat ini adalah kondisi perlu (bukan cukup) untuk deadlock.

## Evidence 2 — PostgreSQL: deadlock terdeteksi otomatis, satu transaksi di-abort
**Claim:** PostgreSQL secara otomatis mendeteksi deadlock pada tingkat baris dan mem-batalkan
salah satu transaksi yang terlibat, memungkinkan transaksi lain untuk melanjutkan.

**Evidence:** "Consider the case in which two concurrent transactions modify a table. The
first transaction executes: UPDATE accounts SET balance = balance + 100.00 WHERE acctnum =
11111; — This acquires a row-level lock on the row. Then, the second transaction executes:
UPDATE accounts SET balance = balance + 100.00 WHERE acctnum = 22222; UPDATE accounts SET
balance = balance - 100.00 WHERE acctnum = 11111; — The second UPDATE finds the row already
locked, so it waits. Transaction two is now waiting on transaction one. Now, transaction one
executes: UPDATE accounts SET balance = balance - 100.00 WHERE acctnum = 22222; — Transaction
one attempts to acquire a row-level lock, but it cannot: transaction two already holds such a
lock. So it waits. Thus, transaction one is blocked on transaction two, and transaction two is
blocked on transaction one: a deadlock condition. PostgreSQL will detect this situation and
abort one of the transactions."

**Source:** PostgreSQL 18 Documentation, §13.3.4 Deadlocks

**URL:** https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
**Published:** 2026-09-24
**Confidence:** HIGH — dokumentasi resmi PostgreSQL
**Corroborated By:** MySQL 8.4 InnoDB deadlock example (Source 4)
**Notes:** "Exactly which transaction will be aborted is difficult to predict and should not be
relied upon."

## Evidence 3 — MySQL/InnoDB: deadlock detection via wait-for graph, rollback transaksi terkecil
**Claim:** InnoDB mendeteksi deadlock secara otomatis menggunakan wait-for graph dan melakukan
rollback pada transaksi yang "termalest" (termasuk, rollback transaksi dengan jumlah baris
terkecil yang dimasukkan/diupdate/dihapus).

**Evidence:** "When deadlock detection is enabled (the default), InnoDB automatically detects
transaction deadlocks and rolls back a transaction or transactions to break the deadlock.
InnoDB tries to pick small transactions to roll back, where the size of a transaction is
determined by the number of rows inserted, updated, or deleted."

**Source:** MySQL 8.4 Reference Manual, §17.7.5.2 Deadlock Detection

**URL:** https://web.archive.org/web/20250126061946/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-detection.html
**Published:** © 2025 Oracle (archived 2025-01-26)
**Confidence:** HIGH — dokumentasi resmi MySQL
**Corroborated By:** PostgreSQL §13.3.4 (both rollback one transaction)
**Notes:** Batas 200 transaksi pada wait-for list; kelebihan diperlakukan sebagai deadlock.
innodb_deadlock_detect dapat dimatikan untuk high-concurrency systems.

## Evidence 4 — MySQL/InnoDB: contoh dua transaksi deadlock (lock order berbeda)
**Claim:** InnoDB menunjukkan deadlock dengan dua klien yang mengunci resource dalam urutan
terbalik; klien B blokir menunggu kunci A, klien A blokir menunggu kunci B.

**Evidence:** "Client A locks Animals row (SELECT ... FOR SHARE), then tries to update Birds.
Client B locks Birds row (SELECT ... FOR SHARE), then tries to update Animals. Client B blocks
waiting for Client A's lock on Animals. Client A attempts UPDATE Birds SET value=40 WHERE
name='Buzzard' → ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting
transaction." SHOW ENGINE INNODB STATUS menampilkan Transaction 43260 (HOLDS S lock pada Birds,
WAITING X lock pada Animals) dan Transaction 43261 (HOLDS S lock pada Animals, WAITING X lock
pada Birds). "InnoDB rolls back transaction (2)."

**Source:** MySQL 8.4 Reference Manual, §17.7.5.1 An InnoDB Deadlock Example

**URL:** https://web.archive.org/web/20241214192732/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-example.html
**Published:** © 2024 Oracle (archived 2024-12-14)
**Confidence:** HIGH — contoh nyata dari dokumentasi resmi
**Corroborated By:** PostgreSQL §13.3.4 contoh yang setara
**Notes:** Error 1213 adalah error deadlock khusus MySQL.

## Evidence 5 — Lock ordering konsisten mencegah deadlock
**Claim:** Mengakuisisi lock pada benda-benda yang sama dalam urutan yang konsisten adalah
defense paling efektif melawan deadlock.

**Evidence:** PostgreSQL: "The best defense against deadlocks is generally to avoid them by
being certain that all applications using a database acquire locks on multiple objects in a
consistent order. In the example above, if both transactions had updated the rows in the same
order, no deadlock would have occurred." Microsoft SQL Server: "Access objects in the same
order. If all concurrent transactions access objects in the same order, deadlocks are less
likely to occur." MySQL: "When modifying multiple tables within a transaction, or different
sets of rows in the same table, do those operations in a consistent order each time. Then
transactions form well-defined queues and do not deadlock."

**Source:** PostgreSQL §13.3.4; Microsoft Learn Deadlocks Guide; MySQL §17.7.5.3

**URLs:**
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html

**Published:** 2026-09-24 (PG); 2026-09-21 (MS); ©2024 Oracle (MySQL)
**Confidence:** HIGH — tiga sumber independen (DB vendor berbeda) sepakat
**Corroborated By:** Percona Blog (sumber 7) — "serialize access... change them to access data in the same order"
**Notes:** Prinsip ini berlaku di semua DBMS berbasis lock.

## Evidence 6 — Transaction duration panjang meningkatkan risiko deadlock
**Claim:** Menahan transaksi terlalu lama meningkatkan peluang deadlock karena mengikat lock
lebih lama, memperbesar surface area konflik.

**Evidence:** PostgreSQL: "it is a bad idea for applications to hold transactions open for
long periods of time (e.g., while waiting for user input)." Microsoft: "A deadlock typically
occurs when several long-running transactions execute concurrently. The longer the transaction,
the longer the exclusive or update locks are held, blocking other activity and leading to
possible deadlock situations." MySQL: "Keep transactions small and short in duration to make
them less prone to collision."

**Source:** PostgreSQL §13.3.4; Microsoft Learn Deadlocks Guide; MySQL §17.7.5.3

**URLs:**
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html

**Published:** 2026-09-24 (PG); 2026-09-21 (MS); ©2024 Oracle (MySQL)
**Confidence:** HIGH — tiga DB vendor sependapat
**Corroborated By:** Percona Blog — "splitting a long transaction into smaller ones, so locks are released sooner"
**Notes:** PostgreSQL deadlock_timeout default 1s (Evidence 9) harus idealnya melebihi typcal
transaction time.

## Evidence 7 — Deadlock retry di application-level (dengan retry, bukan retry naif)
**Claim:** Aplikasi harus selalu siap mengulangi (retry) transaksi yang dibatalkan karena
deadlock, karena deadlock bukan berarti bug fatal jika operasi idempotent.

**Evidence:** MySQL: "Normally, you must write your applications so that they are always
prepared to re-issue a transaction if it gets rolled back because of a deadlock." Microsoft:
"applications should have an error handler that can handle error 1205 ... Implementing an
error handler that catches error 1205 allows an application to handle deadlocks and take
remedial action (for example, automatically resubmitting the query)." Percona: "Before and
above all diagnosis, it is always an important practice to have the applications catch
deadlock error (MySQL error no. 1213) and handles it by retrying the transaction." Microsoft
menambahkan: "application should pause briefly before resubmitting... Randomizing the duration
of the pause minimizes likelihood of deadlock reoccurring."

**Source:** PostgreSQL §13.3.4; MySQL §17.7.5.3; Microsoft Learn; Percona Blog

**URLs:**
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
- https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/

**Published:** 2026-09-24 (PG); 2026-09-21 (MS); ©2024 Oracle (MySQL); 2014-10-28 (Percona)
**Confidence:** HIGH — semua sumber ini secara konsisten menyatakan aplikasi harus retry
**Corroborated By:** AWS retry-backoff guide

## Evidence 8 — Exponential backoff + jitter untuk retry yang aman
**Claim:** Retry harus menggunakan exponential backoff dengan jitter (randomisasi) untuk
meminimalkan kemungkinan deadlock terulang, dan operasi harus idempotent.

**Evidence:** AWS: "Exponential backoff is a technique where operations are retried by increasing
wait times... Operations should be idempotent when you use the retry with backoff pattern."
Microsoft: "Randomizing the duration of the pause minimizes likelihood of deadlock reoccurring
when the resubmitted query requests its locks. For example, the error handler might be coded to
pause for a random duration between one and three seconds."

**Source:** AWS Prescriptive Guidance; Microsoft Learn Deadlocks Guide

**URLs:**
- https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/retry-backoff.html
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17

**Published:** Undated (AWS); 2026-09-21 (MS)
**Confidence:** HIGH — konsisten antar sumber
**Corroborated By:** PostgreSQL §13.3.4 ("retry transactions that abort due to deadlocks")
**Notes:** Idempotency adalah prasyarat penting agar retry aman. Ini terkait langsung dengan
prerequisite lab: "Idempotency".

## Evidence 9 — PostgreSQL deadlock_timeout default 1 detik
**Claim:** PostgreSQL default deadlock_timeout = 1s; ini adalah jeda sebelum mengecek deadlock.
Semakin tinggi beban, semakin disarankan meningkatkan nilai ini; idealnya melebihi typcal
transaction time.

**Evidence:** "This is the amount of time to wait on a lock before checking to see if there is
a deadlock condition. The check for deadlock is relatively expensive, so the server doesn't run
it every time it waits for a lock. We optimistically assume that deadlocks are not common in
production applications and just wait on the lock for a while before checking for a deadlock.
Increasing this value reduces the amount of time wasted in needless deadlock checks, but slows
down reporting of real deadlock errors. ... The default is one second (1s) ... On a heavily
loaded server you might want to raise it. Ideally the setting should exceed your typical
transaction time."

**Source:** PostgreSQL 18 Documentation, §19.12 Lock Management

**URL:** https://www.postgresql.org/docs/current/runtime-config-locks.html#GUC-DEADLOCK-TIMEOUT
**Published:** 2026-09-24
**Confidence:** HIGH — dokumentasi resmi PostgreSQL
**Corroborated By:** (tidak ada sumber lain yang menyebut nilai ini secara eksplisit;
relevan untuk tuning)
**Notes:** log_lock_waits juga terpengaruh oleh parameter ini. Hanya superuser yang bisa
mengubah.

## Evidence 10 — Lock contention reduksi via isolation level dan index
**Claim:** Mengurangi lock contention dapat dilakukan dengan isolasi level yang lebih rendah
(READ COMMITTED), row versioning (READ_COMMITTED_SNAPSHOT), dan index yang tepat.

**Evidence:** Microsoft: "Enable the READ_COMMITTED_SNAPSHOT database option to use row
versioning... Snapshot isolation also uses row versioning, which doesn't use shared locks during
read operations." "Avoid higher isolation levels such as REPEATABLE READ and SERIALIZABLE when
not required." MySQL: "try using a lower isolation level such as READ COMMITTED" untuk locking
reads. "Add well-chosen indexes to your tables so that your queries scan fewer index records
and set fewer locks." Percona: "adding indexes to minimize the rows scanned and locked" dan
"change the transaction isolation level to read committed... to avoid [gap locking]."

**Source:** Microsoft Learn; MySQL §17.7.5.3; Percona Blog

**URLs:**
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html
- https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/

**Published:** 2026-09-21 (MS); ©2024 Oracle (MySQL); 2014-10-28 (Percona)
**Confidence:** HIGH — tiga sumber independen setuju
**Corroborated By:** PostgreSQL §13.2.3 — SSI predicate locks "do not cause any blocking and
therefore can not play any part in causing a deadlock" (sumber berbeda: serializable snapshot
isolation tidak menyebabkan deadlock)
**Notes:** Gap lock adalah contri besar lock contention, terutama pada REPEATABLE READ di
InnoDB.

## Evidence 11 — Lock contention: mengunci terlalu banyak baris meningkatkan konflik
**Claim:** Query yang mengunci banyak baris (UPDATE skop luas) meningkatkan kemungkinan konflik
antar transaksi.

**Evidence:** Percona contoh: "TABLE LOCK table mydb.t1 ... lock mode AUTO-INC waiting" vs
"RECORD LOCKS ... S locks rec but not gap" pada tabel induk t2 — akar penyebab S lock adalah
foreign key constraint. Implikasi: query yang tersebar luas (INSERT massal, UPDATE global)
mengunci banyak resource. Microsoft: "If you use locking reads (SELECT ... FOR UPDATE atau FOR
SHARE), try using a lower isolation level."

**Source:** Percona Blog, §17.7.5.3 MySQL

**URLs:**
- https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html

**Published:** 2014-10-28 (Percona); ©2024 Oracle (MySQL)
**Confidence:** MEDIUM — bukti tidak langsung; berasal dari analisis kasus, bukan klaim
ekspresif
**Corroborated By:** MySQL §17.7.5.3 rekomendasi "Add well-chosen indexes"
**Notes:** Prinsip YAGNI — jangan UPDATE semua row jika tidak perlu; gunakan WHERE yang
spesifik, gunakan ORDER BY untuk urutan lock.

## Evidence 12 — Lock timeout ≠ deadlock (perbedaan mekanisme)
**Claim:** Lock timeout dan deadlock adalah dua kondisi berbeda. Deadlock = siklus tunggu
saling; database mendeteksi dan mem-batalkan satu. Lock timeout = satu transaksi menunggu
lama melebihi batas waktu; tidak ada siklus, hanya ekspirasi waktu tunggu.

**Evidence:** Microsoft: "Deadlocks are sometimes called a deadly embrace... Deadlocking is
often confused with normal blocking... In a deadlock... one transaction is chosen as a victim
and terminated with an error... Lock timeout... When a transaction requests a lock on a resource
locked by another transaction, the requesting transaction waits until the lock is released.
By default, transactions in the Database Engine don't time out, unless LOCK_TIMEOUT is set."
PostgreSQL: "So long as no deadlock situation is detected, a transaction seeking either a
table-level or row-level lock will wait indefinitely." MySQL: innodb_lock_wait_timeout
default 50s untuk lock wait, bukan untuk deadlock detection.

**Source:** Microsoft Learn; PostgreSQL §13.3.4; MySQL §17.7.5.2 + §17.7.5.3

**URLs:**
- https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://web.archive.org/web/20250126061946/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-detection.html

**Published:** 2026-09-21 (MS); 2026-09-24 (PG); ©2025 Oracle (MySQL)
**Confidence:** HIGH — tiga sumber independen jelas membedakan kedua konsep
**Corroborated By:** (semua tiga sumber)
**Notes:** PostgreSQL tidak memiliki lock timeout bawaan — hanya deadlock detection. MySQL
menggunakan innodb_lock_wait_timeout (default 50s) sebagai fallback ketika deadlock detection
dinonaktifkan.

## Evidence 13 — Two-Phase Locking (2PL) dan deadlock
**Claim:** 2PL menjamin conflict-serializability tetapi rentan deadlock; Conservative 2PL (C2PL)
menghilangkan deadlock dengan memaksa acquire semua lock di awal, tapi jarang dipakai.

**Evidence:** Wikipedia (mengutip Bernstein 1987): "By the 2PL protocol, locks are applied and
removed in two phases: 1. Expanding phase: locks are acquired and no locks are released.
2. Shrinking phase: locks are released and no locks are acquired." "Using locks that block
processes, 2PL, S2PL, and SS2PL may be subject to deadlocks that result from the mutual
blocking of two or more transactions." "Conservative 2PL (C2PL) eliminates deadlocks by
requiring all locks upfront."

**Source:** Wikipedia (mengutip Bernstein, Hadzilacos, Goodman 1987); Silberschatz OS Concepts

**URL:** https://en.wikipedia.org/wiki/Two-phase_locking
**Published:** 1987 (Buku asli); continuously updated (Wikipedia)
**Confidence:** MEDIUM — Wikipedia adalah Tier 3, namun mengutip buku teks akademik
**Corroborated By:** (tidak ada verifikasi langsung ke buku asli)
**Notes:** Sumber ini lebih bersifat teoretis. Perlu verifikasi langsung ke buku Bernstein
jika tersedia.

## Evidence 14 — Monitoring: SHOW ENGINE INNODB STATUS dan log_lock_waits
**Claim:** Log database dan metrik dapat mendeteksi deadlock. MySQL menyediakan SHOW ENGINE
INNODB STATUS (terakhir) dan innodb_print_all_deadlocks (semua). PostgreSQL menyediakan
log_lock_waits dan pg_locks.

**Evidence:** MySQL: "SHOW ENGINE INNODB STATUS command... Only the latest deadlock can be
reviewed... With MySQL 5.6, you can enable innodb_print_all_deadlocks to have all deadlocks
in InnoDB recorded in mysqld error log." PostgreSQL: "pg_locks system view" untuk melihat lock
yang sedang berlangsung; log_lock_waits menentukan pesan log tentang lock waits.

**Source:** MySQL §17.7.5.x; PostgreSQL Chapter 27 Monitoring; Percona Blog

**URLs:**
- https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html
- https://www.postgresql.org/docs/current/monitoring.html
- https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/

**Published:** ©2024 Oracle (MySQL); 2026-09-24 (PG); 2014-10-28 (Percona)
**Confidence:** HIGH — terdokumentasi jelas di masing-masing DB
**Corroborated By:** Semua sumber DB di atas
**Notes:** Untuk production, penting meng-aga-aktifkan logging deadlock dan menggunakannya
untuk tuning aplikasi, bukan hanya untuk diagnosis satu kali.
