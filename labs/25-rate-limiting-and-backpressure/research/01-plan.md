# Research Plan

## Research Topic
Rate Limiting & Backpressure — konsep, algoritma, dan praktik engineering untuk membatasi tekanan traffic pada sistem terdistribusi (lab: labs/25-rate-limiting-and-backpressure).

## Objective
Mengumpulkan bukti dari sumber otoritatif (dokumentasi resmi, paper akademik, standar) untuk menopang klaim-klaim teknis dalam spesifikasi lab: token bucket, HTTP 429, Little's Law, queue sebagai penyerap burst sementara, exponential backoff + jitter, retry storm, fairness multi-tenant, serta metrik monitoring (queue depth vs queue age). Memisahkan fakta terverifikasi, klaim sumber, interpretasi, dan hal yang belum terverifikasi.

## Research Questions

1. Bagaimana mekanisme rate limiting yang didokumentasikan resmi (HTTP 429 + `Retry-After`, pola API gateway/CDN)?
2. Apa definisi dan karakteristik algoritma Token Bucket (dan perbandingan dengan Leaky Bucket, Fixed/Sliding Window)? Apa klaim soal "burst wajar"?
3. Apa bunyi Little's Law dan bagaimana penerapannya pada queue/throughput? Apakah contoh "backlog = (arrival − processing) × waktu" sejalan dengan Little's Law atau lebih pada konsep capacity/head-of-line?
4. Apakah benar "queue bukan kapasitas tak terbatas" — apa bukti otoritatif (dokumentasi broker queue, artikel teknis)?
5. Apa rekomendasi resmi soal exponential backoff + jitter (AWS, GCP, Azure, RFC)?
6. Apa itu retry storm / retry amplification dan bagaimana mitigasi resmi (budget, cap, jitter)?
7. Apa bukti otoritatif bahwa rate limit berbasis IP bermasalah (shared NAT/CGNAT) dan praktik kunci limit alternatif (user/API key/tenant)?
8. Bagaimana praktik fairness / per-tenant isolation pada sistem antrian (fair queuing, weighted scheduling) menurut sumber otoritatif?
9. Metric apa yang direkomendasikan untuk memantau backpressure (queue age/lag, depth, processing rate) menurut dokumentasi broker/observability?
10. Apakah klaim "autoscaling tidak menyelesaikan bottleneck downstream" dan "unlimited queue memindahkan titik kegagalan" didukung sumber?

## Search Strategy
- Cari dokumentasi resmi: RFC (RFC 6585 429, RFC 9110, RFC 8305?), AWS/GCP/Azure retry & rate limit docs, Cloudflare/Nginx/Twilio docs, Kafka/SQS/RabbitMQ queue depth & lag docs.
- Cari paper akademik: token bucket (generalized cell rate algorithm), Little's Law (Little 1961, Jewett/Lehoczky), fair queueing (Nagle 1987, WF2Q), Google SRE book (publikasi resmi Google).
- Inspect halaman sumber langsung (bukan snippet) untuk kutipan.
- Cross-check klaim penting dengan ≥2 sumber independen.

## Expected Primary Sources
- RFC 6585 (429 Too Many Requests), RFC 9110 (HTTP semantics)
- RFC 8305? (mungkin tidak relevan) — ganti: RFC 5321 (SMTP) tidak relevan; fokus: RFC 6585, RFC 9110.
- AWS Architecture Blog / AWS SDK retry guidance, Google Cloud retry best practices, Azure exponential backoff guidance
- Google SRE Workbook (O'Reilly publik) — chapters on prioritizing work / load shedding / Little's Law?
- Original paper: J.D.C. Little, "A Proof for the Queuing Formula L = λW" (1961)
- Token bucket: Turner/Jain? atau RFC 2697/2698 (srTCM/tcTBM token bucket specs!) — RFC 2697 = Single Rate Three Color Marker, definisi token bucket resmi IETF.
- Kafka docs (consumer lag), AWS SQS (queue depth/age), RabbitMQ (alarms/backpressure)
- RFC 5737? tidak relevan. NAT/CGNAT: RFC 6598 (shared address space 100.64/10) sebagai bukti IP bersama.
- Netflix/Cloudflare/Stripe engineering articles (Tier 2) untuk retry storm & rate limiting praktis.

## Risks / Unknowns
- Beberapa klaim lab adalah ilustrasi numerik internal (mis. 500 req/detik, 200 user per IP) — tidak perlu diverifikasi sebagai statistik dunia nyata; tandai sebagai skenario hipotetis.
- Klaim "queue age lebih berguna dari queue depth" mungkin bersifat praktik (interpretasi) — cari sumber eksplisit jika ada.
- Little's Law vs formula backlog kumulatif: risiko salah atribusi — bedakan dengan jelas.
- Konten yang berubah (blog vendor) — catat tanggal akses 2026-09-26.
