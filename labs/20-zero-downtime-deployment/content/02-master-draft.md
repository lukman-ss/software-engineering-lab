# Panduan Lengkap Zero-Downtime Deployment

## Problem
Setiap kali aplikasi di-deploy ulang, mematikan instance yang sedang berjalan secara paksa (misalnya dengan `SIGKILL`) dapat memutus koneksi HTTP pengguna yang sedang aktif, menghentikan job yang sedang berjalan di *background*, dan memicu kesalahan sistem jika skema database berubah secara tidak kompatibel. Akibatnya, pengguna mengalami *downtime* atau menerima pesan *error* (seperti 502 Bad Gateway atau koneksi terputus).

## Mental Model
Kunci dari Zero-Downtime Deployment (ZDD) adalah prinsip **"Start new before stopping old"** dan **Koeksistensi**.
Versi aplikasi yang baru (v2) dan yang lama (v1) akan berjalan secara bersamaan selama proses deployment. Ini mengharuskan rute *traffic* diatur secara bertahap dan skema database mendukung pembacaan serta penulisan oleh kedua versi secara simultan tanpa konflik.

## Core Concept
1. **Health Probes**: Memisahkan *Liveness* (status hidup proses) dan *Readiness* (kesiapan menerima lalu lintas jaringan). *Load balancer* hanya mengirim *traffic* ke instance yang berstatus *Ready*.
2. **Graceful Shutdown**: Saat menerima sinyal `SIGTERM`, aplikasi berhenti menerima *request* baru namun membiarkan koneksi HTTP atau proses yang sedang berjalan (*in-flight*) selesai sebelum benar-benar mati (*connection draining*).
3. **Database Parallel Change (Expand & Contract)**: Migrasi database dipisah ke dalam beberapa fase. Pada fase *Expand*, kolom baru ditambahkan tanpa menghapus kolom lama. Aplikasi disesuaikan untuk mendukung baca/tulis di format lama maupun baru. Penghapusan format lama (*Contract*) baru boleh dilakukan setelah seluruh aplikasi versi lama berhenti berjalan.
4. **Cooperative Worker Termination**: Sistem pemrosesan asinkron (*background worker*) menyelesaikan eksekusi pekerjaan (*job*) yang sedang dikerjakan saat ini sebelum keluar, mencegah kerusakan atau *state* parsial.

## Architecture
Lab ini mendemonstrasikan orkestrasi deployment menggunakan komponen berikut:
- `internal/db`: Representasi *in-memory* database yang mengimplementasikan pola Expand/Contract.
- `internal/server`: Server HTTP yang memaparkan *endpoint* untuk *liveness* dan *readiness*, kapabilitas penundaan pemutusan (`preStop`), dan penyelesaian antrean *request*.
- `internal/worker`: Sistem *background queue worker* yang memproses tugas simulasi dan merespons sinyal penghentian.

## How It Works
1. Saat dijalankan, aplikasi melakukan inisialisasi awal. *Readiness probe* mengembalikan HTTP 503.
2. Setelah seluruh komponen siap, *Readiness probe* berubah menjadi HTTP 200. *Load balancer* mulai mengirim *traffic*.
3. Saat *deployment* aplikasi baru dimulai, orkestrator (*orchestrator*) mengirim perintah `SIGTERM` ke versi lama.
4. Aplikasi langsung mengubah *Readiness probe* menjadi tidak siap, sehingga dilepas dari *load balancer*.
5. Aplikasi menjalankan penundaan sementara (`preStop` delay) untuk memberi jeda sampai *routing table* eksternal diperbarui.
6. Masuk ke mode *draining*: server HTTP menunggu seluruh *request* in-flight selesai, dan *worker* merampungkan *job* yang sedang aktif.
7. Aplikasi berhenti dengan aman tanpa menjatuhkan interaksi tunggal.

## What the Tests Prove
- *Liveness* dan *Readiness* beroperasi terpisah; aplikasi menolak *traffic* jika kondisi *readiness* belum valid.
- Modul HTTP sukses menangani penyelesaian operasi pada *request* yang *in-flight* di tengah proses *graceful shutdown*.
- Hook `preStop` terbukti berhasil menunda dimulainya penutupan *listener* sesuai durasi yang dikonfigurasi.
- Komponen *Worker* sukses mendeteksi sinyal *stop*, menyelesaikan 1 buah tugas yang sedang berjalan secara utuh, lalu berhenti beroperasi.
- Mekanisme *fallback* baca/tulis di lapisan database lulus uji kompatibilitas silang untuk rekaman data lama (`name`) dan baru (`first_name`, `last_name`).

## Production Considerations & Common Mistakes

Berikut peringatan teknis yang telah diverifikasi berdasarkan hasil audit research dan engineering:

- **PreStop Hook Implementation**: Implementasi `server.Shutdown` menggunakan `select` pada `time.After(s.preStop)` atau `ctx.Done()`, bukan blok `time.Sleep`. Ini memungkinkan pembatalan jeda saat context deadline habis, mencegah shutdown terkunci. Pada orchestrator seperti Kubernetes, nilai `preStop` (contoh: 1 detik) harus melebihi latency propagation routing table eksternal (kube-proxy, cloud LB) untuk menghindari HTTP 502/504 saat routing masih mengarah ke pod yang sedang di-terminate.
- **Database DDL Lock Contention**: PostgreSQL `ALTER TABLE ... ADD COLUMN` dengan nilai default konstan (misal: `DEFAULT NULL`) adalah metadata-only operation dan aman untuk akses bersamaan. Namun, default volatil (`clock_timestamp()`, `gen_random_uuid()`) memaksa table rewrite lengkap dan `ACCESS EXCLUSIVE` lock. Di produksi, aplikasi harus menggunakan `lock_timeout` pada koneksi database saat mengeksekusi DDL untuk mencegah antrean query yang memanjang.
- **Worker Drain Timeout Behavior**: Worker menyelesaikan job yang sedang aktif selama drain timeout masih tersedia. Ketika timeout (`w.Stop(5 * time.Second)`) tercapai, `w.cancel()` membatalkan konteks, menyebabkan job yang belum selesai dihentikan dan job yang tersisa dalam channel dibuang. Implementasi produksi harus mempertimbangkan job-level context preemption atau hard deadline escalation untuk menjamin tidak ada job yang terlewat.
- **HTTP Worker Signal Handling (PHP-FPM vs Horizon)**: PHP-FPM tidak menangani SIGTERM secara bawaan untuk menyelesaikan in-flight HTTP requests. Implementasi Laravel pada nginx + PHP-FPM memerlukan `process_control_timeout` di `php-fpm.conf` dan konfigurasi nginx `proxy_read_timeout`/`proxy_next_upstream` untuk connection draining. Laravel Horizon, sebaliknya, menyediakan `horizon:terminate` dengan Supervisor `stopwaitsecs` yang memungkinkan worker menyelesaikan job yang sedang berjalan sebelum shutdown.

## Key Takeaways
(Lihat modul Takeaways: `05-key-takeaways.md`)
