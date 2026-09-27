# Research Plan

## Research Topic
Load Testing untuk Aplikasi Software — Praktik Terbaik, Tools, Metrics, dan Strategi Identifikasi Bottleneck pada Aplikasi Booking Bengkel

## Objective
Mengumpulkan bukti empiris dan best practice mengenai load testing berdasarkan skenario aplikasi Booking Bengkel (Login, Booking, Pilih Cabang, Pembayaran, Generate Invoice, WhatsApp Konfirmasi). Fokus pada:
1. Definisi standar industri: Load Test, Stress Test, Spike Test, Endurance/Soak Test
2. Metrics kunci yang harus dipantau (P50, P95, P99, RPS, Error Rate, Resource Utilization)
3. Tools populer dan karakteristiknya (k6, JMeter, Locust, Gatling) — kekuatan, kelemahan, dan rekomendasi penggunaan
4. Strategi identifikasi bottleneck (aplikasi vs database vs external API) — metode isolation dan correlation metrics
5. Common pitfalls dan best practices dari senior engineer
6. Best practice untuk menentukan dan memvalidasi target SLA (contoh ilustratif: P95 < 500ms, Error Rate < 1%, CPU < 75%, Memory < 80% — angka-angka ini bersifat kontekstual tergantung domain/aplikasi, bukan standar universal)
7. Kapan seharusnya load testing dilakukan dalam SDLC

## Research Questions
1. Apa definisi standar industri untuk Load Test, Stress Test, Spike Test, dan Endurance Test?
2. Metrics apa saja yang direkomendasikan oleh organisasi standar (ISO/IEC, NIST) dan praktisi senior?
3. Perbandingan tools: k6 vs JMeter vs Locust vs Gatling — kapan masing-masing optimal?
4. Bagaimana cara membedakan bottleneck di layer aplikasi, database, dan dependency eksternal?
5. Apa saja kesalahan umum yang dilakukan tim engineering saat load testing?
6. Bagaimana cara menentukan target SLA yang realistis dan cara validasi target tersebut?
7. Kapan seharusnya load testing dilakukan dalam lifecycle development?
8. Apa metodologi yang benar untuk merancang skema load testing berdasarkan profil traffic aplikasi Booking Bengkel?

## Search Strategy
- Dokumentasi resmi tools: k6 (grafana.com/docs/k6), JMeter (jmeter.apache.org), Locust (docs.locust.io), Gatling (gatling.io/docs)
- Standar: ISO/IEC 25010 (Software Quality Model), IEEE 829/1012 (Software V&V), NIST SP 800-53
- Sumber teknis otoritatif: Google SRE Book & Workbook, AWS Well-Architected, Microsoft Azure Performance Testing docs, CNCF Observability Whitepaper
- Akademik: IEEE Xplore / ACM Digital Library — papi penelitian "Performance Testing of Web Applications", "Load Testing Best Practices"
- Engineering blogs: Netflix Tech Blog, Uber Engineering, Shopify Engineering
- Cross-check setiap klaim penting melawan 2-3 sumber independen
- Verifikasi semua URL dan buka sumber asli, bukan hanya snippet

## Expected Primary Sources (Tier 1)
- ISO/IEC 25010:2011 — Software Quality Model (definisi performance efficiency)
- Google SRE Book — Testing Chapter 14 (load/stress testing methodology)
- k6 official documentation (grafana.com/docs/k6/latest)
- Apache JMeter User Manual (jmeter.apache.org/usermanual)
- Locust documentation (docs.locust.io)
- Gatling documentation (docs.gatling.io)
- AWS Well-Architected Framework — Performance Efficiency Pillar
- Microsoft Azure — Performance testing guidance
- IEEE 829 — Software Test Documentation
- NIST SP 800-53 — Performance and Capacity Planning (SC-7, etc.)

## Expected Secondary Sources (Tier 2)
- Martin Fowler — blik tentang load testing vs functional testing
- DZone, InfoQ, The New Stack — artikel performa testing
- Engineering blogs: Netflix, Uber, Shopify, GitLab
- CNCF Observability Whitepaper (performance & scalability)
- Vendor whitepapers: Grafana Labs, Datadog, New Relic

## Risks / Unknowns
- Beberapa standar ISO/IEC berbayar (paywalled) — akses melalui versi publik atau open-access summary
- Benchmark tool sering outdated — perlu verifikasi versi dan best practice terbaru
- Claim "best tool untuk X" sering biased (vendor-sponsored) — butuh multiple independent sources
- Real-world case studies perusahaan besar sering proprietary — fokus pada sumber terbuka dan dokumentasi resmi
- Perbedaan metodologi antar industri (web app vs microservice vs API) bisa menciptakan kebingungan definisi
