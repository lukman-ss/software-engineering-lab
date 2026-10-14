Lab 34 kali ini membahas Chaos Engineering dan Fault Injection, atau bagaimana sebuah service terdistribusi dapat membuktikan ketahanan secara empiris sebelum insiden produksi menimpa pelanggan.

Masalah teknis nyata terjadi ketika sebuah layanan bergantung pada payment gateway downstream yang mulai melambat atau mengembalikan error 500 beruntun. Koneksi HTTP client terjebak, connection pool habis dalam hitungan detik, dan kerusakan menyebar ke seluruh pohon call tree hulu. Pengujian integrasi statis dan mock unit saja tidak mampu mereproduksi fluktuasi jaringan yang tidak terduga ini.

Lab ini membangun eksperimen chaos deterministik di Go 1.22+ murni dengan standard library. Frameworknya menggabungkan tiga mekanisme kunci:

• Steady-state monitor atomik (sync/atomic) yang menilai error rate transaksional bukan CPU atau RAM.
• Circuit Breaker tiga fase (Closed, Open, Half-Open) yang memutus downstream yang sekarat dan memasukkan antrean fallback asinkron.
• Auto-abort runner berbasis ticker yang menetralkan fault injector dalam hitungan milidetik saat steady-state terlanggar.

Demo interaktif menunjukkan steady-state error rate tetap 0.00% di tengah simulasi pemadaman 100% di downstream, sementunya fase eksperimen tak terproteksi memicu pembatalan otomatis dan pemulihan penuh dalam 200ms cooldown.

Mental model: Build safety rails before injecting turbulence; measure customer-visible output, never component internals.

Batasan penting: Lab ini menggunakan injeksi in-memory dan counter kumulatif sederhana, bukan manipulasi lapisan jaringan (iptables/eBPF) atau sliding time window eksponensial seperti di lingkungan produksi sejati.

Pelajari kode sumber penuh dan lima test suite yang lolos race detector di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/34-chaos-engineering

#softwareengineering #golang #chaosengineering #distributedsystems #resilience #systemdesign
