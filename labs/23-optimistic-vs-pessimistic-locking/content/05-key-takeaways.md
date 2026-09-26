# Key Takeaways

1. **Lost update adalah silent data corruption**—dua transaksi konkurensi menimpa satu sama lain tanpa error. Ini bukan bug, tapi konsekuensi logis dari isolation tanpa concurrency control.

2. **Transaction alone tidak cukup**—`BEGIN ... COMMIT` tanpa explicit locking atau version guard atau atomic statement *tidak mencegah* lost update di READ COMMITTED.

3. **Pessimistic locking** (`SELECT ... FOR UPDATE`) = *prevent conflict* (blokir writer lain via row-level exclusive lock). Cocok untuk high contention atau correctness-critical (balances, seat reservations, stock).

4. **Optimistic locking** (version/timestamp + affected_rows check) = *detect conflict* (version mismatch = conflict). Cocok untuk low contention atau read-heavy (profile/CMS/CRM). Harus handle retry atau 409.

5. **Atomic single-statement** (`UPDATE ... SET stock = stock - N WHERE stock >= N`) = *eliminate race window* entirely (statement-level atomicity). Cocok untuk simple counter/quota/stock decrement.

6. **Selection criteria:** Atomic first (simple counter), optimistic (low contention/read-heavy), pessimistic (high contention/correctness-critical). Distributed locks hanya untuk cross-database/microservice scenario.

7. **Anti-patterns:** Transaction alone (mengira itu cukup), holding lock across network call (payment gateway), ignoring 0-rows-affected (silent lose update), distributed locks when database can solve it (prefer database-native).

8. **Isolation level:** PostgreSQL/Oracle default READ COMMITTED = lost update mungkin tanpa guard. MySQL default REPEATABLE READ = lost update mungkin tanpa guard. SERIALIZABLE = error harus thrown.

9. **Production monitoring:** Optimistic—retry rate tinggi = need more pessimistic. Pessimistic—lock wait time tinggi = need caching atau sharding.

10. **Demo test verified:** Naive → stock 99 (lost update), Pessimistic → stock 50 (exact), Optimistic (no retry) → 1 success/19 conflicts (state guarded), Optimistic+retry → stock 80 (converged), Atomic → stock 50 (lockless).