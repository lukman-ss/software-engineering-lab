# Research Report

## Research Question
Bagaimana merancang deployment zero-downtime untuk stack Nginx + Laravel + PostgreSQL + Redis Queue pada transisi v1.4 → v1.5 sehingga HTTP tidak downtime, job tidak hilang, v1.4 & v1.5 coexist, DB tetap compatible, dan rollback cepat?

## Executive Summary
Zero-downtime pada stack ini tercapai dengan menggabungkan: (a) **Expand–Contract** untuk skema DB agar v1.4 dan v1.5 coexist, (b) **readiness-gated routing** (NGINX reload / K8s readiness / Laravel `/up` + DiagnosingHealth), (c) **graceful shutdown / connection draining** (NGINX HUP+quit, K8s SIGTERM→grace→SIGKILL, `queue:restart`/`horizon:terminate` + Supervisor `stopwaitsecs` > job terlama), dan (d) **deploy sequence yang atomik** dengan checkpoint rollback. Tidak ada sumber Tier 1 yang bertentangan; perbedaan hanya trade-off biaya (Blue-Green) vs durasi coexist (Rolling) dan keterbatasan fitur NGINX OSS.

## Findings

### Finding 1 — Traffic Orchestration: Readiness Menentukan Kapan Instance Masuk Load Balancer
Claim: Instance baru tidak boleh menerima traffic sebelum readiness PASS (bukan sekadar liveness/container running). Instance lama di-drain sebelum di-terminate.
Evidence: K8s readiness probe menghapus pod dari Service endpoints; liveness hanya restart loop; Pod Lifecycle: delete → endpoint removal → SIGTERM → grace window → SIGKILL. NGINX HUP: master buka config/barusocket, start worker baru, kirim graceful-shutdown ke worker lama yang terus layani koneksi aktif hingga selesai.
Sources: Kubernetes Pod Lifecycle (https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/), Kubernetes Probes (https://kubernetes.io/docs/concepts/workloads/pods/probes/), NGINX Control (https://nginx.org/en/docs/control.html), Laravel Deployment The Health Route (https://laravel.com/docs/11.x/deployment#the-health-route)
Confidence: HIGH

### Finding 2 — Graceful Shutdown / Connection Draining Mencegah Request Terputus
Claim: Terminasi paksa (SIGKILL / `docker stop` tanpa grace) memutus `POST /checkout` yang sedang berjalan → payment terpotong tapi response hilang.
Evidence: K8s mengirim SIGTERM lalu menunggu `terminationGracePeriodSeconds` (default 30s, configurable) sebelum SIGKILL; NGINX QUIT = graceful, TERM/INT = fast. Laravel: di dalam deployment panggil `php artisan horizon:terminate` atau `queue:restart`; Supervisor harus `stopwaitsecs` > durasi job terlama, sonst job dibunuh.
Sources: Kubernetes Pod Lifecycle, NGINX Control, Laravel Horizon Deploying (https://laravel.com/docs/11.x/horizon#deploying-horizon), Laravel Queues (https://laravel.com/docs/11.x/queues)
Confidence: HIGH

### Finding 3 — Database Backward Compatibility via Expand–Contract
Claim: Selama rolling deploy v1 dan v2 hidup bersama dan query tabel yang sama; rename/drop kolom yang breaking membuat deployment gagal walau load balancer sempurna.
Evidence: Fowler Sato: fase Expand (tambah kolom baru nullable, kedua kolom ada), Migrate (deploy app yang paham keduanya, migrasikan data), Contract (DROP kolom lama setelah v1 mati). PostgreSQL: `ADD COLUMN` nullable tanpa volatile default = metadata-only, tanpa rewrite, aman untuk rolling; `ADD COLUMN` dengan `NOT NULL` + volatile default atau `SET DATA TYPE` = rewrite + lock berat; `ADD CONSTRAINT NOT VALID` + `VALIDATE CONSTRAINT` memberi jalan non-blocking untuk constraint baru.
Sources: Martin Fowler Parallel Change (https://martinfowler.com/bliki/ParallelChange.html), PostgreSQL ALTER TABLE (https://www.postgresql.org/docs/current/sql-altertable.html), PostgreSQL DDL (https://www.postgresql.org/docs/current/ddl.html)
Confidence: HIGH

### Finding 4 — Blue-Green Memberi Rollback Tercepat, Rolling Paling Hemat Kapasitas
Claim: Blue-Green: deploy ke GREEN idle, smoke test, switch router; rollback = switch balik (detik). Rolling: ganti instance A→B→C satu-per-satu; rollback = `kubectl rollout undo` atau re-deploy image lama, lebih lambat tetapi tidak butuh 2× kapasitas.
Evidence: Fowler Blue-Green: dua env identik, salah satu live, switch router; rollback cepat; DB tetap butuh pemisahan schema change dulu. K8s Deployment `strategy: RollingUpdate` adalah default rolling; pelajaran NGINX HUP/USR2 menunjukkan mekanisme drain yang sama mendasari keduanya.
Sources: Martin Fowler Blue Green Deployment (https://martinfowler.com/bliki/BlueGreenDeployment.html), Kubernetes Pod Lifecycle, NGINX Control (binary upgrade rollback path), NGINX upstream (drain/backup params)
Confidence: HIGH

### Finding 5 — Queue Worker Deployment: Job Lama Tidak Hilang
Claim: Queue worker bukan HTTP instance; ia long-lived dan harus disuruh berhenti setelah job aktif selesai, serta tidak mengambil job baru selama drain.
Evidence: Laravel Queues & Horizon: gunakan `queue:restart` (cache signal; worker selesaikan job sekarang lalu exit) atau `horizon:terminate` (graceful, job aktif selesai lalu Supervisor restart dengan kode baru). `block_for=0` memblokir indefinite dan mengabaikan SIGTERM sampai job berikutnya; hindari. `after_commit=true` mencegah dispatch job dalam transaksi yang belum commit.
Sources: Laravel Queues, Laravel Horizon Deploying, Laravel Task Scheduling `onOneServer`/`withoutOverlapping` (https://laravel.com/docs/11.x/scheduling)
Confidence: HIGH

### Finding 6 — Deployment Sequence v1.4 → v1.5 yang Memenuhi Lima Kriteria Soal
Sequence (tanpa downtime, tanpa hilang job, coexist, DB compatible, rollback cepat). Dua varian yang valid; pilih sesuai budget infra:

**Varian A — Rolling (hemat, default K8s)**
1. **Migrate: Expand DB (backward-compatible, non-blocking).** `ALTER TABLE ... ADD COLUMN full_name TEXT` (nullable, tanpa default volatile), `ADD CONSTRAINT ... NOT VALID` jika perlu. Tidak ada `DROP name` di tahap ini. Deploy ini sebagai migration terpisah sebelum code v1.5, di-verify di staging yang menjalankan v1.4. Sumber: Fowler + PG ALTER TABLE.
2. **Deploy v1.5 app yang paham kedua skema.** Code membaca `COALESCE(full_name, name)` dan menulis ke keduanya (dual-write) atau menulis `full_name` + backfill trigger. v1.4 tetap jalan tanpa error karena `name` masih ada. Sumber: Parallel Change migrate phase.
3. **Backfill data.** Job/ migration `UPDATE ... SET full_name = name WHERE full_name IS NULL` dalam batch kecil, `NOT VALID` → `VALIDATE CONSTRAINT` setelah backfill.
4. **Rolling update web tier dengan readiness gate.** K8s: `readinessProbe: httpGet: /up` (listener DiagnosingHealth cek DB+Redis+Migration), `startupProbe` jika boot lambat, `terminationGracePeriodSeconds` ≥ `stopwaitsecs` ≥ job terlama + `preStop: sleep 5` untuk LB deregister. NGINX (non-K8s): `nginx -s reload` setelah upstream baru healthy; atau `drain` di Plus. Sumber: K8s probes/lifecycle + NGINX control + Laravel health route.
5. **Queue worker rolling.** Scale worker baru berdampingan atau `horizon:terminate` / `queue:restart` setelah web baru ready. Pastikan `retry_after` > `timeout` terlama, `block_for` > 0, `after_commit=true`, job idempotent (key transaksi/webhook). `onOneServer` untuk scheduler. Sumber: Queues/Horizon/Scheduling.
6. **Verifikasi production health.** Metrics tetap sehat (p95 latency, 5xx, queue depth, DB connections) — bukan hanya CI hijau (sesuai lab).
7. **Contract DB (setelah semua v1.4 mati & observasi).** `ALTER TABLE ... DROP COLUMN name` hanya setelah tidak ada instance v1.4. Ini rilis terpisah. Sumber: Parallel Change contract phase.

Rollback pada varian A: `kubectl rollout undo` + DB tetap aman karena `name` belum di-drop; durasi menit (image pull + readiness).

**Varian B — Blue-Green (rollback detik, butuh 2× kapasitas)**
Langkah 1–3 sama. Deploy v1.5 ke GREEN pool yang belum terima traffic; smoke test GREEN via `/up`. Switch router (NGINX upstream weight / K8s Service selector / ALB target group) dari BLUE→GREEN setelah health PASS. BLUE tetap idle untuk instant rollback (`GREEN→BLUE`).

### Finding 7 — Health Check Tidak Boleh Asal `{status: ok}`
Claim: Endpoint yang hanya return 200 tanpa cek dependensi memberi false-positive; readiness harus cek DB, Redis, storage, status migration.
Evidence: Laravel Deployment: default `/up` hanya 200 jika boot tanpa exception; `DiagnosingHealth` event dilempar untuk tambahan cek; listener dapat throw exception → 500 → readiness FAIL. K8s readiness probe yang salah = traffic masuk ke pod belum siap → 502/503 massal.
Sources: Laravel Deployment, Kubernetes Probes
Confidence: HIGH

## Areas of Agreement
- Semua sumber setuju: **Start new before stopping old** adalah invarian; coexist v1+v2 tidak terhindarkan sehingga DB harus backward-compatible.
- Terminasi harus kooperatif (SIGTERM/grace/drain/`queue:restart`), bukan SIGKILL.
- Rollback harus direncanakan sebelum deploy, bukan setelah gagal; Blue-Green menguji disaster-recovery tiap rilis.

## Areas of Disagreement
Tidak ada pertentangan material antar sumber Tier 1. Perbedaan hanya trade-off: Blue-Green = rollback tercepat & butuh kapasitas 2×; Rolling = kapasitas efisien & rollback lebih lambat. NGINX OSS vs Plus soal `health_check`/`drain` adalah perbedaan fitur, bukan kontradiksi arsitektural. Redis rolling-upgrade tidak ditemukan sumber resmi → dicatat sebagai NOT VERIFIED, bukan sumber perselisihan.

## Limitations
- Redis rolling-upgrade tanpa Tier 1: klaim cluster rolling upgrade = NOT VERIFIED.
- Durasi grace spesifik (K8s 30s default, `stopwaitsecs`, `terminationGracePeriodSeconds`) tergantung workload (checkout, upload, WebSocket, webhook timeout) — harus diukur per aplikasi.
- NGINX `drain`/`health_check` di doc upstream ditandai komersial; implementasi OSS butuh reload atau openresty/Lua — detailnya di luar scope lab minimal.
- Laravel Octane (Swoole/RoadRunner/FrankenPHP) mengubah startup/reload semantics; lab ini diasumsikan PHP-FPM + `queue:work`/`horizon`.

## Conclusion
Zero-downtime untuk Nginx+Laravel+PostgreSQL+Redis Queue dicapai bukan dengan satu flag, melainkan orkestrasi: DB expand dulu (kolom baru nullable, constraint `NOT VALID`), deploy code yang handle kedua skema, backfill, rolling/blue-green dengan readiness yang cek DB/Redis/migration + graceful drain pada web dan queue worker, verifikasi metrics, baru contract DB. Kegagalan di poin mana pun — health check asal, migration breaking, kill tanpa drain, atau rollback tanpa rencana — akan mengembalikan `502/503` atau yang lebih buruk: payment terpotong tanpa response. Senior engineer menjawab ketujuh pertanyaan Rule of Thumb sebelum menekan deploy.
