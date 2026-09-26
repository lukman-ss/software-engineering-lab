# Research Plan

## Research Topic
Load Testing untuk Aplikasi Software — Praktik Terbaik, Tools, Metrics, dan Strategi Identifikasi Bottleneck

## Objective
Mengumpulkan bukti empiris dan best practice mengenai:
1. Jenis-jenis performance testing (Load, Stress, Spike, Endurance/Soak)
2. Metrics kunci yang harus dipantau (P50, P95, P99, RPS, Error Rate, Resource Utilization)
3. Tools populer dan karakteristiknya (k6, JMeter, Locust, Gatling)
4. Strategi identifikasi bottleneck (aplikasi vs database vs external API)
5. Common pitfalls dan best practices dari senior engineer

## Research Questions
1. Apa definisi standar industri untuk Load Test, Stress Test, Spike Test, dan Endurance Test?
2. Metrics apa saja yang direkomendasikan oleh organisasi standar (ISO, NIST, IEEE) dan praktisi senior?
3. Perbandingan tools: k6 vs JMeter vs Locust vs Gatling — kapan masing-masing optimal?
4. Bagaimana cara membedakan bottleneck di layer aplikasi, database, dan dependency eksternal?
5. Apa saja kesalahan umum yang dilakukan tim engineering saat load testing?
6. Best practice untuk menentukan target SLA (P95 < 500ms, Error Rate < 1%, dll)?
7. Kapan seharusnya load testing dilakukan dalam lifecycle development?

## Search Strategy
- Cari dokumentasi resmi tools (k6.io, jmeter.apache.org, locust.io, gatling.io)
- Cari standar: ISO/IEC 25010 (quality model), IEEE 1012 (V&V), NIST SP 800-53
- Cari artikel dari situs engineering terpercaya: Netflix Tech Blog, Uber Engineering, Google SRE, AWS Well-Architected, Microsoft Azure Architecture Center
- Cari academic papers pada performance testing, bottleneck analysis
- Cross-check claims dari multiple sources

## Expected Primary Sources (Tier 1)
- k6 official documentation (grafana.com/docs/k6)
- Apache JMeter User Manual (jmeter.apache.org)
- Locust documentation (docs.locust.io)
- Gatling documentation (gatling.io/docs)
- ISO/IEC 25010:2011 Software Quality Model
- Google SRE Book (Chapter on Load Testing)
- AWS Well-Architected Performance Efficiency Pillar
- Netflix Tech Blog articles on performance testing
- Academic papers: "Performance Testing of Web Applications" (IEEE), "Load Testing Best Practices" (ACM)

## Expected Secondary Sources (Tier 2)
- Martin Fowler articles on performance testing
- Articles from: DZone, InfoQ, The New Stack, DevOps.com
- Vendor whitepapers (Cloudflare, Datadog, New Relic, Grafana Labs)
- Conference talks (QCon, Velocity, SREcon)

## Risks / Unknowns
- Beberapa standar ISO/IEC mungkin tidak tersedia secara gratis (paywalled)
- Benchmark tools sering outdated — perlu verifikasi versi terbaru
- Claim "best tool untuk X" sering biased — butuh multiple independent sources
- Real-world case studies dari perusahaan besar sering proprietary