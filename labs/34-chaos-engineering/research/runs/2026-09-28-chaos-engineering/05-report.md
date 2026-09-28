# Research Report

## Research Question
Bagaimana prinsip Chaos Engineering dan metodologi Fault Injection diterapkan untuk memverifikasi ketahanan sistem (seperti Circuit Breakers dan Graceful Degradation) sebelum terjadi kegagalan di production?

## Executive Summary
Chaos Engineering adalah disiplin eksperimentasi untuk menguji ketahanan sistem terdistribusi terhadap kondisi turbulen. Dengan mendefinisikan *Steady State*, membuat hipotesis, dan menyuntikkan gangguan terkontrol (*Fault Injection*), tim dapat membuktikan atau membantah ketahanan arsitektur (seperti failover, timeout, dan circuit breakers) sebelum insiden nyata merusak bisnis.

## Findings

### Finding 1: Definisi dan Prinsip Chaos Engineering
Claim: Chaos Engineering memerlukan eksperimen terstruktur berbasis hipotesis pada metrik *Steady State* daripada pengujian acak.
Evidence: Principles of Chaos Engineering menetapkan bahwa eksperimen harus berfokus pada output perilaku sistem dan membandingkannya saat gangguan disuntikkan.
Sources: Principles of Chaos Engineering (https://principlesofchaos.org/)
Confidence: HIGH

### Finding 2: Pengendalian Blast Radius
Claim: Eksperimen chaos wajib membatasi area dampak melalui subset trafik production atau lingkungan staging dengan tombol *abort* otomatis.
Evidence: AWS Well-Architected Reliability Pillar menekankan perlunya teknik *canary analysis* dan prosedur penghentian darurat.
Sources: AWS Well-Architected Reliability Pillar (https://docs.aws.amazon.com/)
Confidence: HIGH

### Finding 3: Verifikasi Circuit Breakers dan Graceful Degradation
Claim: Suntikan latensi atau kegagalan layanan pihak ketiga (seperti Payment Gateway) harus memicu Circuit Breaker untuk mencegah *Cascading Failures*.
Evidence: Google SRE & praktik industri menunjukkan bahwa kegagalan downstream yang tidak dibatasi oleh timeout dan circuit breaker dapat menghabiskan thread pool downstream dan menyebabkan kegagalan menyeluruh.
Sources: Google SRE Book, Netflix TechBlog
Confidence: HIGH

## Areas of Agreement
- Semua sumber sepakat bahwa chaos engineering bukan tentang merusak sistem secara sembarangan, melainkan pengujian ilmiah berhipotesis.
- Pengukuran *steady state* berbasis metrik eksternal/pengguna adalah mutlak.

## Areas of Disagreement
No material contradictions discovered.

## Limitations
- Penelitian dibatasi pada literatur teoretis dan arsitektur umum; implementasi spesifik bahasa pemrograman (Java/Node.js/Go) memiliki pustaka pendukung (Resilience4j, Polly) yang bervariasi.

## Conclusion
Chaos Engineering adalah metode paling efektif untuk membuktikan klaim ketahanan sistem (high availability, fault tolerance) melalui pengujian empiris yang terkontrol.
