## Evidence 1

Claim: Deadlock terjadi ketika dua transaksi saling menahan lock yang dibutuhkan satu sama lain, sehingga keduanya terblokir secara permanen.

Evidence: "A deadlock occurs when two or more tasks permanently block each other by each task having a lock on a resource that the other tasks are trying to lock."

Source: Deadlocks Guide - SQL Server

URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide

Confidence: HIGH

Corroborated By: Source 1 (PostgreSQL Documentation: "two (or more) transactions each hold locks that the other wants")

Notes: Fakta dasar tentang kondisi deadlock.

## Evidence 2

Claim: Database memiliki deadlock monitor untuk mendeteksi siklus dan akan membatalkan (abort) salah satu transaksi sebagai victim.

Evidence: "The Database Engine deadlock monitor periodically checks for tasks that are in a deadlock. If the monitor detects a cyclic dependency, it chooses one of the tasks as a victim and terminates its transaction with an error."

Source: Deadlocks Guide - SQL Server

URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide

Confidence: HIGH

Corroborated By: Source 1 ("PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved")

Notes: Menjelaskan mekanisme internal database dalam menangani siklus deadlock.

## Evidence 3

Claim: Aplikasi sebaiknya menahan transaksi sesingkat mungkin dan memodifikasi objek dengan urutan yang sama untuk mencegah deadlock.

Evidence: "To help minimize deadlocks: Access objects in the same order... Keep transactions short and in one batch."

Source: Deadlocks Guide - SQL Server

URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide

Confidence: HIGH

Corroborated By: Source 1 ("The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order.")

Notes: Praktik terbaik dari tingkat aplikasi.

## Evidence 4

Claim: Aplikasi harus siap merespons error deadlock (seperti error 1205 di SQL Server) dengan melakukan retry.

Evidence: "Implementing an error handler that catches error 1205 allows an application to handle deadlocks and take remedial action (for example, automatically resubmitting the query that was involved in the deadlock)."

Source: Deadlocks Guide - SQL Server

URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide

Confidence: HIGH

Corroborated By: Source 1 ("deadlocks can be handled on-the-fly by retrying transactions that abort due to deadlocks.")

Notes: Konfirmasi praktik retry dari sistem database yang berbeda.

## Evidence 5

Claim: Pemeriksaan deadlock memakan waktu, sehingga database biasanya menunggu (deadlock timeout) sebelum melakukan pengecekan.

Evidence: "The check for deadlock is relatively expensive, so the server doesn't run it every time it waits for a lock. We optimistically assume that deadlocks are not common in production applications and just wait on the lock for a while before checking for a deadlock."

Source: PostgreSQL Lock Management (`deadlock_timeout`)

URL: https://www.postgresql.org/docs/current/runtime-config-locks.html

Confidence: HIGH

Corroborated By: Source 3 ("The default interval is 5 seconds... periodic deadlock detection helps reduce the overhead")

Notes: Menjelaskan mengapa terjadi jeda waktu tunggu sebelum deadlock victim dipilih.
