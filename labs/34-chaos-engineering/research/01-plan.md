## Research Topic
Chaos Engineering --- Menguji Ketahanan Sistem dengan Sengaja Menyuntikkan Kegagalan

## Objective
Meneliti dan mengumpulkan bukti otoritatif mengenai prinsip, metode, dan pola implementasi Chaos Engineering, Fault Injection, Circuit Breakers, serta pengendalian Blast Radius untuk layanan terdistribusi.

## Research Questions
1. Apa definisi resmi dan prinsip inti Chaos Engineering menurut pembuatnya (Netflix/Gremlin)?
2. Bagaimana cara kerja Fault Injection (network latency, service crash, resource exhaustion) dalam menguji Circuit Breakers dan Graceful Degradation?
3. Apa strategi pembatasan Blast Radius dan pengukuran Steady State yang aman untuk production GameDay?
4. Bagaimana mitigasi kegagalan berantai (Cascading Failures) menggunakan Timeouts, Retries, dan Fallbacks?

## Search Strategy
1. Cari panduan dan standar resmi dari Chaos Engineering Principles (principlesofchaos.org).
2. Cari publikasi dan dokumentasi resmi dari Netflix TechBlog, Gremlin, AWS Well-Architected Framework (Reliability Pillar), dan Google SRE Book.
3. Cross-check klaim terkait circuit breaker state machine, timeout propagation, dan retry storms.

## Expected Primary Sources
1. Principles of Chaos Engineering (principlesofchaos.org)
2. Netflix TechBlog (Chaos Engineering articles)
3. AWS Well-Architected Framework - Reliability Pillar (Chaos Engineering / Testing)
4. Google SRE Book - Chapter 24 (Managing Incidents) & Chapter 25 (Tracking Reliability)
5. Gremlin Reliability Management Documentation

## Risks / Unknowns
- Beberapa platform proprietary (seperti Gremlin) mungkin mengubah API atau terminologi fitur fault injection.
- Variasi implementasi circuit breaker antar bahasa/framework (Resilience4j, Polly, Hystrix) dapat menghasilkan perilaku timeout yang berbeda.
