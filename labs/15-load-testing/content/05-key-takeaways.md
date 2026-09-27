## Key Takeaways

1. Load testing menemukan batas sistem sebelum pengguna menembukannya. Tanpa data nyata, deployment adalah asumsi.
2. Gunakan persentil (P95, P99), bukan rata-rata. Rata-rata 21ms tidak menutupi P95 1,5 detik dibawah saturation.
3. Jenis tes menentukan tujuan. Smoke untuk validasi skrip; spike untuk lonjakan trafik; soak untuk memory leak selama jam-hari.
4. Bottleneck terdeteksi dengan memonitor semua lapisan secara paralel. Tunggu `http_req_waiting` tinggi berarti app/DB, bukan jaringan.
5. Connection pool yang terebut menghasilkan antrian + degradasi tail. Ini non-linear: dari 21ms (2 VUs) menjadi 1,35s (50 VUs).
6. Tool pilihan kontekstual, bukan vendor-driven. k6 untuk JS/API, Locust untuk Python, JMeter untuk GUI/XML, Gatling untuk JVM.
7. Persentil (P50/P90/P95/P99) hanya mencakup request sukses (HTTP 2xx); request gagal dihitung sebagai error tanpa latensi tercatat.
8. Thresholds harus ditentukan upfront berdasarkan SLA bisnis. Contoh P95 < 200ms adalah ilustratif, bukan standar universal.
9. Jangan uji endpoint /health atau data fixture kecil. Workload harus merepresentasikan trafik produksi nyata.
10. Load test di lingkungan mirip production. Laptop dev tidak bisa merepresentasikan spesifikasi server.
11. Demo mencetak subset P50/P95/P99, walau `Result` menghitung P90 juga.
12. Tidak ada keajaiban satu benchmark. Semua angka kontekstual tergantung arsitektur, workload, dan resource yang disponibilitas.