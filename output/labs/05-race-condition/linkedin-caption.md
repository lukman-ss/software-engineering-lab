Lab 05 — Race Condition pada penjualan stok: 2 kasir, stok 1 unit, 2 penjualan berhasil, invariant rusak (`1 != 2 + 0`).

Setup: Go module dengan unit test deterministik (channel barrier), stress test 500 goroutine, PostgreSQL integration test (row lock, atomic update, UNIQUE constraint).

Mekanisme yang dibahas:
- Lost update via READ → CHECK → WRITE non-atomic
- Atomic conditional update (`UPDATE ... WHERE stock > 0`)
- `SELECT ... FOR UPDATE` row lock
- UNIQUE constraint sebagai last defense
- Optimistic locking, distributed lock bukan default

Mental model: tanyakan "apa yang terjadi jika request lain mengubah data antara READ dan WRITE?" pada setiap business logic.

URL: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/05-race-condition

#RaceCondition #Concurrency #Database #SoftwareEngineering #Go