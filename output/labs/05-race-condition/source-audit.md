# Source Audit – Lab 05 – Race Condition

**Resolved lab**: `labs/05-race-condition`

**Series**: Senior Software Engineer Daily #05

**Repo URL**: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/05-race-condition

**Git status**: clean (working tree matches HEAD)

**Files read (human‑authored)**:
- README.md
- datarace.go, datarace_test.go, datarace_unsafe_test.go
- atomic_inventory.go, atomic_update_test.go
- unsafe_inventory.go, mock_inventory.go, lost_update_test.go
- inventory.go, booking_test.go
- postgres_rowlock.go, postgres_booking.go
- schema.sql
- go.mod, go.sum

**Supporting files outside lab**: none (all needed code resides in lab directory).

**Excluded / inaccessible**: generated binaries, Docker volumes, external docs.

**Demonstrates**: business‑level race conditions (lost update) and safe alternatives (atomic conditional update, row‑level lock, UNIQUE constraint).

**Evidence ledger** (selected examples):
- `TestLostUpdate_Deterministic` shows invariant break (initial 1 ≠ 2 + 0) – lines 24‑33, 146‑154 in `lost_update_test.go`.
- `TestAtomicUpdate_StockOne` verifies atomic decrement preserves invariant – lines 42‑63, 100‑115 in `atomic_update_test.go`.
- `Test500_ConcurrentBooking` proves UNIQUE constraint prevents double booking – lines 61‑73, 120‑144 in `booking_test.go`.
- `PostgresRowLockRepository.TrySell` uses `SELECT … FOR UPDATE` to serialize updates – lines 30‑67 in `postgres_rowlock.go`.
- `AtomicInventory.DecrementStock` implements single‑statement conditional update – lines 46‑60 in `atomic_inventory.go`.

**Title candidates**:
1. “Lost Update & Safe Countermeasures in Inventory Systems”
2. “Race Condition pada Penjualan Stok: Dari Bug ke Solusi”
3. “Mengatasi Lost Update dengan Atomic Update & Row Lock”
4. “Invariant Rusak Karena Race Condition – Analisis & Perbaikan”
5. “Concurrency pitfalls in inventory: lost update, booking clash, and fixes”

**Selected title**: #2 – clear, concrete, matches Indonesian language requirement, and reflects actual lab focus.

**English variant**: “Race Condition in Stock Sales: From Bug to Solution”

**Style references**: not accessed.