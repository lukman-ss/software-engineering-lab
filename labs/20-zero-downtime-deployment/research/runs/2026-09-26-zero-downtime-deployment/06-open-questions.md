# 06 — Open Questions

## Unanswered / Weak-Evidence Claims
- **Redis cluster/HA rolling upgrade**: Official Redis upgrade doc not found (404 at redis.io). Need verified guidance for `v1.4 → v1.5` queue worker compatibility across a Redis upgrade. (NOT VERIFIED)
- **NGINX OSS active health checks**: Open-source upstream has only passive checks (`max_fails`/`fail_timeout`). Exact mechanism for readiness-gated NGINX upstream swap without Plus or Lua is implementation detail.
- **Laravel Octane graceful drain**: SIGTERM behavior for Swoole/RoadRunner/FrankenPHP workers under rolling deploy not quantified from docs (only `octane:reload`/`octane:stop` described).
- **Postgres `ALTER TABLE ... SET DATA TYPE` without rewrite conditions**: docs mention binary-coercible exceptions; threshold for "large table" needing `USING` rewrite needs empirical sizing for this stack.

## Deeper Research Needed
- Spesifik `terminationGracePeriodSeconds` dan durasi aktual `checkout`/`upload`/`websocket` pada aplikasi lab ini (batas atas job terlama) — harus di-measure, bukan diasumsusi.
- Apakah queue ada job long-running (menit) yang membuat 30s default grace tidak cukup dan harus dipertahankan via `stopwaitsecs` Supervisor.
- Strategi backfill data pada tabel besar: chunk size, lock behavior, kemungkinanan `VALIDATE CONSTRAINT` memakan waktu berjam — perlukah `CONCURRENTLY`-style app-level backfill job.

## Next Research Directions
1. Verifikasi Redis rolling upgrade dari sumber tier 1 (Redis Cluster spec, Sentinel docs) — cari URL yang valid.
2. Tes kecil (jika lab mengizinkan) pada environment production-like: NGINX HUP reload vs container `stop` dan verifikasi `connection draining` via tcpdump / metrics request count.
3. Benchmark `ALTER TABLE ADD COLUMN` vs `ADD COLUMN ... DEFAULT` pada tabel sample ukuran nyata (periksa lock wait vs metadata-only di PG 16+ di environment lab).
4. Validasi `/up` listener dependency-check latency vs health probe timeouts (avoid readiness 500 during transient slow-db).

## Engineering Notes (Pre-Handoff to Engineer Agent)
- Deploy tooling belum ditentukan (K8s vs Docker Compose vs Envoyer vs Forge). Sequence di atas bersifat portable; Engineer Agent dapat memilih orchestrator dan melengkapi `deployment.yaml`/`supervisor.conf`/`nginx.conf` sesuai.
- Database migration script harus dipisahkan menjadi `expand` dan `contract` sebagai file migrasi terpisah (contoh: `2025_01_02_000000_expand_add_fullname_users.php` dan `2025_01_04_000000_contract_drop_name_users.php`), bukan sekali `ALTER`.
- Health listener `DiagnosingHealthEvent` perlu implementasi konkrit — didefinisikan di `app/Listeners/` lab.
