# Key Takeaways

1. **Rata-rata Menipu**: Metrik rata-rata waktu respons menyembunyikan lonjakan latensi ekor (*tail latency*). Analisis performa wajib menggunakan persentil (P50, P90, P95, P99).
2. **Pengujian Bertahap**: Load testing wajib dimulai dengan Smoke Test (konkurensi rendah, 2–5 VU) untuk memvalidasi integritas pengujian sebelum menjalankan Stress Test skala penuh.
3. **Saturasi Sumber Daya Membentuk Antrean**: Ketika konkurensi melebihi kapasitas sumber daya hilir (seperti connection pool basis data), request terpaksa mengantre sehingga latensi P95 melonjak tajam secara non-linier.
4. **Korelasi Klien-Server**: Mengukur latensi sisi klien saja tidak cukup untuk menemukan akar masalah; metrik engine harus dikorelasikan dengan metrik server (CPU, RAM, utilisasi connection pool).
5. **Generator Beban Bebas Lock**: Runner pengujian beban harus dirancang thread-safe tanpa kontensi kunci internal agar generator itu sendiri tidak menjadi penghambat pengukuran.
6. **Optimasi Algoritma Statistik**: Pengurutan slice in-memory cocok untuk pengujian berskala kecil, namun pengujian berskala jutaan metrik memerlukan histogram streaming untuk menghemat memori.
