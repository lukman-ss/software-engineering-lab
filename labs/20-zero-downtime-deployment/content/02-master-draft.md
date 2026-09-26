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

Berdasarkan hasil audit implementasi dan riset, ada beberapa peringatan teknis yang wajib diperhatikan di tingkat produksi:

- **Kubernetes Asynchronous Routing**: Di lingkungan Kubernetes, pengiriman sinyal `SIGTERM` ke container dan pembaruan aturan jaringan (melalui *EndpointSlice* ke *kube-proxy*) berjalan secara **asinkron**. Tanpa jeda `preStop` (contoh: `sleep 5`), aplikasi mungkin keburu menutup saluran koneksi sementara *load balancer* masih mengirimkan *traffic* kepadanya, mengakibatkan *error* HTTP 502/504.
- **Database DDL Lock Contention**: Walaupun pola *Expand & Contract* menyelesaikan masalah aplikasi, operasi DDL (seperti `ALTER TABLE`) di *relational database* seperti PostgreSQL membutuhkan `ACCESS EXCLUSIVE lock`. *Lock* ini menghalangi seluruh proses baca dan tulis lain. Di tahap produksi, ini wajib disertai mitigasi seperti *timeout* pada level koneksi database (misal: `SET lock_timeout`).
- **Worker Buffered Jobs Abandonment**: Pada implementasi *lab* ini, *worker* hanya dirancang untuk menyelesaikan *job* aktif (yang *sedang* dieksekusi). Sisa *job* yang sudah di-kuota (*buffered*) di *channel* dalam antrean memori dibatalkan (*abandoned*) akibat konfigurasinya menggunakan `context.cancel()` instan. Implementasi nyata harus memikirkan cara membuang atau mengembalikan (*re-queue*) seluruh *buffer* yang tersisa.
- **POSIX vs Framework Signals**: Terminasi *worker* secara universal berstandar POSIX (`SIGTERM` atau `SIGQUIT`), meski sejumlah platform (misal ekosistem Laravel) memanipulasi utilitas penyimpanan sementara seperti Redis (*cache*) untuk mengirimkan sinyal perintah `restart`.

## Key Takeaways
(Lihat modul Takeaways: `05-key-takeaways.md`)
