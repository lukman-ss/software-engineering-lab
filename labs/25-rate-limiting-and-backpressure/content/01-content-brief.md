# Content Brief

Topic: Rate Limiting & Backpressure — Sistem Skalabel yang Stabil di Bawah Beban Lebih

Target Reader:
- Engineer backend dan SRE yang sudah menguasai dasar Go
- Audience Indonesia (bacaan teknis, istilah inggris standar)

Problem:
Sistem distribusi runtuh ketika laju permintaan masuk melebihi kapasitas pemrosesan. Tanpa rate limiting, antrian tidak terbatas tumbuh eksis, memicu kehabisan memori dan *timeout cascade*. Tanpa backpressure, worker overload dan terus-menerima beban melebihi kemampuan sistem. Tanpa jitter pada *retry*, klien-seribu klien selihut setelah jeda yang sama menciptakan *retry storm*.

Core Mental Model:
"Apakah pekerjaan masuk lebih cepat daripada kemampuan sistem menyelesaikannya?" Jika ya → masalah kapasitas, butuh intervensi. Semua mekanisme (token bucket, leaky bucket, bounded queue, jittered retry) bertujuan satu hal: **katakan "cukup" secara terukur** sebelum sistem kehabisan sumber daya.

Approved Research Status: APPROVED
Approved Engineering Status: APPROVED

Main Concepts:
- Token Bucket: *burst accommodation* dengan *refill rate* jangka panjang
- Leaky Bucket: *traffic smoothing* dengan kecepatan drain konstan
- HTTP 429 (RFC 6585): standar respon rate limit + `Retry-After`
- Bounded Queue: *load shedding* cepat lewat channel non-blocking
- Exponential Backoff + Jitter: Full/Equal/Decorrelated/No Jitter (AWS Marc Brooker)
- Multi-tenant isolation: key berbasis API key / tenant, bukan IP (hindari CGNAT RFC 6598)
- Little's Law (L = λW): fondasi perencanaan kapasitas antrian

Verified Behaviors:
- Token bucket me-allow *burst* hingga kapasitas B, lalu menolak sampai refill cukup
- Token bucket menghitung `RetryAfterSeconds` yang akurat ketika token habir
- Leaky bucket menolak *burst* saat kapasitas penuh, mengizinkan setelah *leak*
- Registry tenant isolation: tenant A dan B punya kuota mandiri
- Bounded queue (`TrySubmit`: `ErrQueueFull`) menolak cepat saat buffer penuh
- Token bucket + bounded queue bersaing *race-free* di bawah 50 goroutine (race detector bersih)
- Backoff Full Jitter berada dalam batas `[0, min(cap, base * 2^attempt)]`
- HTTP middleware mengembalikan 429 + `Retry-After` pada overload tenant
- Demo CLI menampilkan semua perilaku di atas secara live

Available Case Studies:
- Demo CLI (`cmd/demo/main.go`) menampilkan 4 fase: token bucket, leaky bucket, bounded queue, retry strategies

Warnings:
- Rate limiting state bersifat *in-memory* (replica proses tidak terkoordinasi)
- Empiris benchmark memori/Latency belum disertakan (lihat research-audit/06-gaps.md)
- 5/10 sumber riset termasuk Wikipedia (lihat research-audit poin non-blokir)
- Lab bersifat *single-node*; belum mendemonstrasikan distributed rate limiting
