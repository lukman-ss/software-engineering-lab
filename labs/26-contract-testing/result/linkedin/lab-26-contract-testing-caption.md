Lab 26 kali ini membahas Consumer-Driven Contract testing, deteksi breaking schema changes otomatis, serta implementasi CI/CD deployment gate di Go.

Masalah teknis yang sering terjadi pada microservices adalah perubahan skema API tanpa koordinasi: penggantian casing enum status, penggantian nama field kunci JSON, atau mutasi tipe data numerik menjadi teks. Pada pengujian unit lokal dengan mock statis, kesalahan ini tidak terdeteksi karena mock cepat usang. Sementara pengujian end-to-end di staging memakan waktu lama dan sering terganggu flakiness jaringan, sehingga breaking change kerap lolos hingga merusak klien seluler di lingkungan produksi.

Lab ini membuktikan empat hal konkret:
• Minimal Subset Contract: Klien seluler mendefinisikan ekspektasi payload minimal yang hanya memuat field yang benar-benar dikonsumsi.
• Build-Time Verification: Engine verifikasi mengeksekusi interaksi HTTP terhadap server ephemeral dan memvalidasi struktur secara rekursif dalam waktu kurang dari satu detik.
• Breaking Change Barrier: Tiga mutasi payload terbukti langsung memicu kegagalan verifikasi dengan pesan kesalahan field yang spesifik, memblokir deployment otomatis di tingkat CI.
• Safe API Evolution: Pola dual DTO paralel mempertahankan kepatuhan kontrak V1 sekaligus mengekspos skema V2 untuk migrasi bertahap tanpa downtime.

Mental model inti: Mocks drift into false confidence; contract tests turn consumer assumptions into executable provider guardrails.

Satu batasan penting: Lab ini menggunakan engine verifikasi Go murni berbasis pustaka standar tanpa binary eksternal Pact, serta menyederhanakan manajemen provider state melalui routing path deterministik.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing

#softwareengineering #golang #microservices #testing #devops #systemdesign
