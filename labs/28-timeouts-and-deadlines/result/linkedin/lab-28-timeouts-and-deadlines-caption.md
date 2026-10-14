Lab 28 kali ini membahas mitigasi cascading failure akibat dependency lambat melalui penerapan context deadline, exponential backoff dengan full jitter, circuit breaker, dan idempotency key di Go.

Masalah teknis yang sering terjadi adalah satu dependency eksternal melambat dari 300ms menjadi 60 detik. Berdasarkan Little's Law (L = lambda * W), jumlah request in-flight meledak dan menguras seluruh worker thread pool, sehingga endpoint independen ikut lumpuh total.

Lab ini membuktikan empat pilar ketahanan konkret:
• Context Deadline Propagation: Child context secara hierarkis mewarisi sisa deadline parent dan melepaskan goroutine melalui buffered channel tanpa memory leak.
• Full Jitter Exponential Backoff: Distribusi acak seragam sleep = random(0, min(max, base * 2^(attempt-1))) membubarkan thundering herd retry storm saat dependency pulih.
• Three-State Circuit Breaker: Transisi CLOSED, OPEN, dan HALF_OPEN memotong rantai retry seketika saat downstream terbukti down, memberi ruang isolasi pemulihan.
• Idempotency Deduplication: Penyimpanan hash response transaksi melindungi operasi finansial dari risiko double charge ketika status timeout bersifat indeterminate.

Mental model inti: Fast failure frees resources immediately; slow dependencies hoard threads until total starvation. Timeouts are resource guardrails, not arbitrary numbers.

Satu batasan penting: Lab ini menggunakan in-memory storage untuk idempotency key dalam satu proses Go, bukan Redis terdistribusi atau cluster database transaksional.

Pelajari implementasi lengkap dan seluruh pembuktian test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/28-timeouts-and-deadlines

#softwareengineering #golang #distributedsystems #systemdesign #resilience #microservices
