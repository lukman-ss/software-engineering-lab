# Poin Kunci Pembelajaran (Key Takeaways)

1. **Circuit Breaker Bukan Obat Penyembuh Downstream**  
   Circuit Breaker tidak memperbaiki kegagalan pada layanan downstream yang rusak, melainkan menghentikan kehancuran diri pemanggil (*caller self-destruction*) dengan mencegah penumpukan alokasi thread, soket, dan memori.

2. **Efisiensi Mekanisme Fail-Fast**  
   Ketika sirkuit beralih ke status `OPEN`, penolakan permintaan dilakukan secara instan dalam skala nanodetik/mikrodetik tanpa melakukan transmisi I/O jaringan, berbeda jauh dengan latensi pemblokiran timeout yang memakan waktu ratusan milidetik.

3. **Uji Coba Pemulihan Terukur via Canary Probe (HALF-OPEN)**  
   Fase `HALF-OPEN` memungkinkan verifikasi pemulihan downstream secara terukur melalui sejumlah kecil permintaan uji coba (*probe*), mencegah lonjakan trafik seketika yang dapat merusak kembali server yang baru bangkit.

4. **Waktu Tenggang Lab Bersifat Ilustratif**  
   Konfigurasi waktu lab (100ms HTTP timeout dan 300ms cooldown) dirancang khusus agar pengujian otomatis berjalan cepat. Konfigurasi produksi nyata wajib dikalibrasi mengikuti SLA, metrik latensi P99, dan waktu pemulihan downstream yang realistis.

5. **Batasan Penghitungan Kegagalan Berurutan**  
   Implementasi lab mengandalkan penghitungan kegagalan berturut-turut (*consecutive failures*). Pada sistem skala enterprise, pertimbangkan penggunaan metrik berbasis rasio persentase kesalahan dengan jendela geser (*sliding time-window*).

6. **Larangan Silent Fallback pada Mutasi Kritis**  
   Jangan gunakan *silent fallback* (pengalihan semu tanpa error) untuk operasi mutasi state yang berisiko tinggi seperti pendebitan saldo, pembayaran, atau pengurangan stok inventaris. Laporkan error pemutus sirkuit secara transparan ke sistem hulu.

7. **Pemisahan Alur Sinkron dan Asinkron**  
   Layanan kritis (seperti *Payment Gateway*) memerlukan proteksi kegagalan langsung (*fail-fast*), sedangkan layanan non-kritis (seperti notifikasi pesan/email) harus diisolasi ke dalam antrean latar belakang (*message queue*) dengan penanganan coba-ulang berbasis *backoff*.

8. **Integritas Konkurensi**  
   Operasi transisi status dan penghitungan metrik sirkuit harus dilindungi mekanisme konkurensi (seperti `sync.Mutex`) untuk menjamin tidak terjadinya kondisi perlombaan data (*data race*) di bawah beban tinggi multi-goroutine.
