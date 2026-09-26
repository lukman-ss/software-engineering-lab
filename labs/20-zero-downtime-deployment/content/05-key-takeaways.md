# Key Takeaways

1. **Start New Before Stopping Old**: Zero-downtime deployment bergantung pada koeksistensi, di mana versi baru (v2) dan versi lama (v1) harus mampu berjalan berdampingan pada infrastruktur yang sama dan mengakses database yang sama.
2. **Pemisahan Liveness dan Readiness**: Liveness probe hanya mendeteksi proses yang hang atau deadlock, sedangkan Readiness probe menentukan apakah instance diizinkan menerima lalu lintas jaringan dari load balancer.
3. **Pentingnya PreStop Delay**: Pada orchestrator modern seperti Kubernetes, pengiriman SIGTERM dan pembaruan iptables/routing table berlangsung asinkron. Jeda `preStop` mutlak diperlukan untuk mencegah timbulnya galat HTTP 502/504 selama masa transisi.
4. **Connection Draining Menjaga In-Flight Traffic**: Pemutusan server harus kooperatif menggunakan graceful shutdown (`Server.Shutdown`), menolak koneksi baru namun menunggu seluruh koneksi aktif selesai diproses sebelum proses dihentikan.
5. **Pola Expand and Contract untuk Database**: Skema database tidak boleh diubah secara destruktif dalam satu langkah deployment. Tambahkan kolom baru tanpa menghapus kolom lama (Expand), jalankan kode aplikasi baru dengan fallback, lalu bersihkan kolom usang di deployment terpisah (Contract).
6. **Mitigasi DDL Lock di Produksi**: Menjalankan migrasi database di produksi memerlukan kehati-hatian terhadap `ACCESS EXCLUSIVE` lock pada tabel besar. Parameter seperti `lock_timeout` penting untuk mencegah antrean query yang menumpuk.
7. **Penyelesaian Job pada Worker Asinkron**: Background worker harus merespons sinyal terminasi (`SIGTERM`/`SIGQUIT`) secara kooperatif dengan menuntaskan pekerjaan yang sedang dieksekusi agar tidak meninggalkan state data setengah jalan.
8. **Drain Buffer Queue**: Penghentian worker tidak cukup hanya menunggu satu job aktif, melainkan harus memperhitungkan dan mengamankan tugas-tugas yang telah telanjur masuk ke antrean buffer internal aplikasi.
