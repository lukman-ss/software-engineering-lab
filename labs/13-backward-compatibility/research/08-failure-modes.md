# Failure Modes

1. **Incomplete Backfill**: Data gagal termigrasi penuh sehingga read yang bergantung ke skema baru mengembalikan null.
2. **Dual-write Consistency**: App crash di antara write lama dan write baru, membuat data out of sync. Solusi: CDC (Change Data Capture) atau Outbox Pattern.
3. **Database Lock**: Alter table yang mem-block read/write table. Solusi: Gunakan operator non-blocking (e.g. `CREATE INDEX CONCURRENTLY` di Postgres).
4. **Premature Contract Phase**: Menghapus field lama saat mobile client lawas masih menggunakannya. App client force crash.
