# Key Takeaways

1. **Rate limiting mengendalikan batas sistem, backpressure mengirim sinyal tekanan ke belakang.** Rate limiting bekerja di entrance (mis. API gateway); backpressure lebih umum — ketika downstream tidak mampu, upstream harus melambat (Google SRE, Source 10).

2. **Token bucket dan leaky bucket (sebagai meter) setara secara matematis, berbeda sudut pandang.** Token bucket mentolerir burst hingga kedalaman bucket sambil menjaga rata-rata jangka panjang; leaky bucket menghaluskan output ke laju tetap. NGINX `limit_req_module` memakai leaky bucket sebagai meter dengan status default 503.

3. **HTTP 429 adalah kode status standar rate limiting (RFC 6585) dan wajib tidak di-cache.** Respons dengan 429 `MUST NOT be stored by a cache`; sertakan header `Retry-After` bila memungkinkan. 503 lebih tepat untuk overload server dibanding 429 untuk kuota client — keduanya berbeda kasus pemakaian, bukan kontradiksi.

4. **Full Jitter mencegah thundering herd: `delay = random(0,1) × min(cap, base × 2^attempt)`.** AWS melaporkan dengan 100 klien berebut, jitter mengurangi jumlah panggilan lebih dari separuh. Nilai default AWS SDK (50ms transient, 1000ms throttling, cap 20s) adalah pilihan spesifik AWS, bukan mandat universal — kalibrasi ke SLA downstream.

5. **Bounded queue adalah implementasi backpressure: tolak cepat, jangan tumpuk.** `TrySubmit` dengan pola `select-default` mengembalikan `ErrQueueFull` seketika (zero-allocation), mencegah antrean tumbuh tak terbatas dan memori habis. Ukur umur antrean (queue age), bukan hanya panjangnya.

6. **Pembatasan per-tenant, bukan per-IP.** CGNAT (RFC 6598) membuat banyak tenant berbagi satu IP; pembatasan IP berbasis alamat menimbulkan throttling kolateral. Registry per-tenant key (`X-API-Key`) memberi isolasi bucket per tenant.

7. **Little's Law (L = λW) adalah dasar analisis antrean.** Contoh terverifikasi: 5.000.000 item ÷ 2.000 item/detik = 2.500 detik ≈ 41 menit 40 detik. Pertumbuhan backlog = (arrival rate − processing rate) × waktu.

8. **Pendekatan berlapis memberi degradasi elegan.** Stripe memakai 4 limiter (request rate, concurrent requests, fleet usage shedder, worker utilization shedder). Google SRE menambah 4 level criticality (CRITICAL_PLUS → SHEDDABLE) dan sinyal utilitas (executor load average, CPU, memori) — reject criticality tinggi hanya setelah semua criticality lebih rendah ditolak.

9. **Retry budget mencegah retry storm.** Per-request: maksimal 3 attempt; per-client: rasio retry di bawah 10%. AWS SDK memakai token bucket retry quota (kapasitas 500 token) yang menghentikan retry ketika token habis agar aplikasi fail fast.

10. **Implementasi lab ini in-memory dan single-process.** State rate limiter tidak tersinkronisasi lintas replica — Redis (INCR/EXPIRE, Lua, sorted set) adalah langkah lanjut untuk deployment multi-instance. Beberapa gap terbuka tetap tercatat: `RetryAfterSeconds` butuh `refillRate > 0`, race latent `Stop` vs `TrySubmit`, error job dibuang di worker loop, dan sifat statistik decorrelated jitter belum dibuktikan test.
