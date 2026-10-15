Lab 25 kali ini membahas strategi pengendalian beban terdistribusi melalui implementasi Token Bucket, Leaky Bucket, Bounded Queue backpressure, dan AWS Full Jitter retry di Go.

Masalah mendasar pada sistem tanpa pembatas adalah ketimpangan laju: ketika trafik masuk melompat melampaui kemampuan downstream, antrean bertumpuk tanpa batas di RAM. Berdasarkan Little's Law (L = lambda * W), 5.000.000 request dengan processing rate 2.000/detik membutuhkan waktu 41 menit 40 detik untuk menguras buffer, memicu lonjakan latensi masif dan crash akibat Out of Memory (OOM).

Lab ini mengimplementasikan empat pilar pertahanan konkret murni menggunakan Go standard library:
• Token Bucket & Leaky Bucket: Pengisian waktu berjalan (lazy monotonic refill) mengizinkan burst terukur atau memaksakan aliran output konstan tanpa background ticker yang boros CPU.
• Isolasi Kunci Per-Tenant: Registri bucket berbasis API key (X-API-Key) mencegah fenomena collateral throttling akibat ribuan pengguna berbagi IP publik di balik Carrier-Grade NAT (RFC 6598).
• Bounded Queue Backpressure: Buffer berkapasitas tetap dengan pola select-default menolak kelebihan pekerjaan seketika (ErrQueueFull) tanpa alokasi memori tambahan.
• AWS Full Jitter Retry: Randomisasi jeda tidur sleep = random(0, min(cap, base * 2^attempt)) membubarkan sinkronisasi thundering herd saat downstream pulih.
• HTTP 429 RFC 6585 Middleware: Menyertakan header Retry-After presisi untuk memandu penjadwalan ulang request klien.

Mental model inti: Rate limiting mengontrol pintu gerbang masuk; backpressure mempropagasi sinyal penolakan ke belakang; Full Jitter mendesinkronisasi gelombang percobaan ulang klien.

Batasan penting implementasi: Lab ini menggunakan state in-memory single process. Pada skala multi-pod produksi terdistribusi, koordinasi token bucket membutuhkan shared state eksternal seperti Redis cluster dengan skrip atomic Lua.

Pelajari arsitektur lengkap, kode implementasi, dan pengujian race-clean di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/25-rate-limiting-and-backpressure

#softwareengineering #golang #distributedsystems #systemdesign #resilience #performance #backend
