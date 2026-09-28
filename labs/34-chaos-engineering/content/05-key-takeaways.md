# Key Takeaways

1. **Chaos bukan merusak, melainkan eksperimen ilmiah.** Eksperimen chaos harus didasarkan pada hipotesis dan metrik steady-state, bukan serangan acak.
2. **Steady state ukur dari output pengguna.** Fokus pada throughput dan error rate, bukan pada kondisi internal CPU atau memori komponen.
3. **Circuit Breaker mencegah cascading failure.** Dengan membatasi jumlah failure yang dapat diterima dan beralih ke state terbuka saat threshold dilampaui, ia melindungi layanan hulu.
4. **Fallback menyelamatkan pengalaman pengguna.** Graceful degradation memastikan bahwa kegagaan downstream tidak langsung menyebabkan error ke layanan hulu.
5. **Auto-abort dan Clear adalah penting.** Mekanisme pembatalan otomatis dan netralisasi injector (*blast radius control*) melindungi sistem dari eksperimen yang berlari terlalu lama.
6. **Gunakan race detector.** Semua komponen yang bersifat thread-safe harus diuji dengan `go test -race` untuk memastikan tidak ada race condition.
7. **Lab ini disederhanakan.** Metrik kumulatif, injeksi in-memory, dan error rate sederhana dipilih untuk kejelasan. Di production, pertimbangkan sliding window, p99 latency, canary routing, dan tracing terdistribusi.
8. **Tetap amati dan pertahankan.** Setelah eksperimen selesai, sistem harus kembali ke keadaan normal — pemulihan otomatis (`Half-Open → Closed`) wajib diverifikasi.
9. **Jalankan eksperimen secara bertahap.** Mulai dari canary/staging sebelum memperluas ke traffic production penuh, untuk meminimalkan blast radius.
10. **Verifikasi secara empiris.** Ketahanan hanya terbukti jika uji coba otomatis (unit test + demo) melewati semua quality gates, termasuk compile, race detector, dan auto-abort.
