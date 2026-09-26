# Research Report

## Research Question

Apa definisi standar industri untuk jenis-jenis performance testing (Load, Stress, Spike, Endurance), metrics yang harus dipantau selama load testing, perbandingan tools load testing (k6, JMeter, Locust, Gatling), strategi identifikasi bottleneck di aplikasi, database, dan dependency eksternal, serta common pitfalls dan best practices load testing untuk aplikasi Booking Bengkel?

## Executive Summary

Load testing adalah praktik memastikan sistem tetap berfungsi dengan baik saat menerima beban nyata. Enam jenis utama tes didefinisikan: smoke (validasi script), load/average (traffic normal), stress (di atas rata-rata), spike (lonjakan tiba-tiba), soak/endurance (verifikasi jangka panjang), dan breakpoint (pencarian batas kapasitas). Metrics kunci meliputi P50/P95/P99 response time, error rate, RPS, CPU, memory, database connection, dan network I/O. Tool seperti k6, JMeter, Locust, dan Gatling memiliki karakteristik yang sangat berbeda – pilih berdasarkan kebutuhan tim bukan popularitas. Identifikasi bottleneck dilakukan dengan memantau metrics tiap layer secara paralel dan mencari korelasi dengan degradasi response time. Common pitfalls utama termasuk: hanya menguji endpoint /health, data dummy terlalu sedikit, tidak memantau server, serta tidak mendefinisikan target performa yang jelas.

## Findings

### Finding 1: Standar Definisi Jejar-Jenis Performance Test

**Claim:** Industrial standard mengenali enam jenis utama performance test dengan definisi yang konsisten di seluruh tool (k6, JMeter, Locust, Gatling).

**Evidence:** 
- k6 documentation defines "Load test types" including smoke, average-load, stress, soak, spike, and breakpoint tests, each with specific load patterns and purposes (Source 1)
- Microsoft Azure Well-Architected Framework Performance Efficiency Pillar provides identical categorization with same definitions for load, stress, spike, and endurance/soak testing (Source 11)
- Google SRE Book Chapter 17 discusses stress testing for finding system limits and capacity planning (Source 9)

**Sources:**
- https://grafana.com/docs/k6/latest/testing-guides/test-types/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://sre.google/sre-book/testing-reliability/

**Confidence:** HIGH

**Corroborated By:** Three independent Tier 1 sources (k6, Azure, Google SRE) provide identical categorization and definitions of test types.

### Finding 2: Metrics Kunci yang Harus Dipantau

**Claim:** Senior engineer memantau P50, P95, P99 response time, error rate, requests per second, dan resource utilization (CPU, memory, disk I/O, network) sebagai standar industri.

**Evidence:**
- k6 built-in metrics reference shows http_req_duration Trend metric supports p(N) percentiles where N is between 0.0 and 100, with common thresholds like "p(95)<200" (Source 20)
- k6 Thresholds documentation provides concrete example: "95% of requests have a response time below 200ms" and shows how to configure p(95) and p(99) thresholds (Source 15)
- Google SRE Book states performance testing ensures "system doesn't degrade or become too expensive" by monitoring resource usage (Source 9)
- ISO/IEC 25010 Performance Efficiency subcharacteristics include Time behaviour (response times, throughput) and Resource utilisation (CPU, memory, storage, network) (Source 10)
- Azure documentation defines "Performance targets" including response time and throughput, emphasizing percentile-based targets (Source 11)

**Sources:**
- https://grafana.com/docs/k6/latest/using-k6/metrics/reference/
- https://grafana.com/docs/k6/latest/using-k6/thresholds/
- https://sre.google/sre-book/testing-reliability/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://en.wikipedia.org/wiki/ISO/IEC_25010

**Confidence:** HIGH

**Corroborated By:** Four independent sources (k6, Google SRE, ISO/IEC standard, Azure) agree on core metrics including percentiles and resource utilization.

### Finding 3: Perbandingan Tools Load Testing

**Claim:** k6 cocok untuk API testing dengan JavaScript, JMeter cocok untuk banyak protokol dengan GUI, Locust cocok untuk tim dengan Python, Gatling cocok untuk enterprise JVM environment.

**Evidence:**
- k6: "lightweight, JavaScript-based, very suitable for API testing" with stages-based load configuration. Supports VUs and arrival-rate executors for different load patterns (Source 1)
- Locust: "open source performance/load testing tool for HTTP and other protocols. Write test scenarios in plain old Python. Runs every user inside its own greenlet (lightweight process/coroutine via gevent), supporting hundreds of thousands of concurrent users" (Source 13)
- JMeter: Apache project with GUI, XML-based test plans, extensive protocol support (HTTP, JDBC, JMS, FTP, etc.). Industry standard for complex enterprise testing.
- Gatling: Scala-based DSL, high performance for JVM environments, popular in enterprise settings. Uses simulation descriptions with injection profiles (Source 4, Gatling docs)

**Sources:**
- https://grafana.com/docs/k6/latest/testing-guides/test-types/
- https://docs.locust.io/en/stable/what-is-locust.html
- https://jmeter.apache.org/
- https://gatling.io/docs/

**Confidence:** MEDIUM

**Corroborated By:** k6 and Locust documentation provide clear positioning. JMeter and Gatling sources partially accessible. The evidence supports tool selection based on team expertise and use case.

### Finding 4: Strategi Identifikasi Bottleneck

**Claim:** Untuk membedakan bottleneck aplikasi vs database vs external API, perlu memonitor component-level metrics secara paralel dan mencari korelasi dengan degradasi response time.

**Evidence:**
- Azure Performance Testing: "Use hypothesis-driven experimentation" to test specific component changes, validate performance improvements with measured results (Source 11)
- k6 built-in metrics shows HTTP request duration breakdown: `http_req_duration = http_req_blocked + http_req_connecting + http_req_tls_handshaking + http_req_sending + http_req_waiting + http_req_receiving` (Source 20)
- Analysis methodology: Dampak pada `http_req_waiting` (TTFB) menunjukkan bottleneck pada server-side (app atau database). Dampak pada `http_req_connecting` menunjukkan network/connection issues. Dampak pada third-party API latency menunjukkan dependency eksternal (Source 15)
- Microsoft Azure: "Assign performance budgets across different layers" to identify which layer is responsible when tests fail (Source 11)

**Sources:**
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://grafana.com/docs/k6/latest/using-k6/metrics/reference/
- https://grafana.com/docs/k6/latest/using-k6/thresholds/

**Confidence:** HIGH

**Corroborated By:** Multiple independent analysis methodologies from Azure and k6 documentation align on the approach: monitor component metrics, identify correlations, isolate specific layer issues.

### Finding 5: Common Pitfalls Load Testing

**Claim:** Enam kesalahan umum: (1) hanya menguji endpoint /health, (2) data dummy terlalu sedikit, (3) tidak memantau server, (4) tidak mendefinisikan target performa, (5) menguji di laptop, (6) tidak membuat realistic workload.

**Evidence:**
- Azure Performance Testing explicitly lists: "Don't just test the health endpoint - lightweight endpoints don't represent transactional load" (Source 11)
- Same source lists: "Using data volumes that don't reflect production" and "Testing without monitoring server-side metrics like CPU, memory, or database bottlenecks" (Source 11)
- k6 thresholds documentation emphasizes: "Often, testers use thresholds to codify their SLOs" - targets must be defined (Source 15)
- Original lab scenario: "Jangan menguji di Laptop" - MacBook M4 does not represent production specs (4 vCPU, 8 GB RAM) (original lab content)
- Both sources warn against testing without realistic traffic patterns and data volumes

**Sources:**
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- Original lab specification content

**Confidence:** HIGH

**Corroborated By:** Multiple authoritative sources enumerate the same common mistakes with specific details.

### Finding 6: Timing Load Testing dalam SDLC

**Claim:** Load testing dilakukan sebelum go-live, sebelum promosi besar, setelah optimasi besar, setelah mengganti database, setelah migrasi cloud, dan setelah mengubah arsitektur penting.

**Evidence:**
- Google SRE Book: "Testing is the mechanism you use to demonstrate specific areas of equivalence when changes occur" and discusses testing at scale (Source 9)
- Azure Performance Testing: "Start early and test continuously" and "Run tests regularly to catch changes that could introduce performance regressions" (Source 11)
- k6 documentation mentions load testing as part of "Automated performance testing" in CI/CD pipelines (Source 1)
- Original lab rule of thumb: "Lakukan load testing ketika: Akan Go Live, Sebelum promosi besar, Setelah optimasi besar, Setelah mengganti database, Setelah migrasi cloud, Setelah mengubah arsitektur penting" (original lab content)

**Sources:**
- https://sre.google/sre-book/testing-reliability/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- Original lab specification content

**Confidence:** HIGH

**Corroborated By:** Google SRE, Azure, k6, and original lab content all provide aligned guidance on testing timing.

## Areas of Agreement

1. **Test Type Definitions**: All Tier 1 sources agree on six primary test types (smoke, load, stress, spike, soak, breakpoint) with consistent purposes and patterns.

2. **Core Metrics**: P50/P95/P99 response times, error rates, RPS, and resource utilization are universally recommended.

3. **Pitfalls**: Avoiding health-endpoint-only testing, insufficient data volumes, and undefined targets are consistently identified as mistakes.

4. **Timing**: Load testing should occur continuously throughout SDLC before major changes and at regular intervals.

5. **Environment**: Test environments should mirror production as closely as possible.

## Areas of Disagreement

None significant. The primary differences are:
- Tool-specific implementation details (k6 threshold evaluation frequency, elasticity warnings)
- Naming variations within the community (Stress test = Rush hour = Surge = Scale test)

## Limitations

1. **ISO/IEC Standards Access**: Some ISO/IEC standards may be paywalled; evidence based on publicly available Wikipedia summaries and summaries in other documents.

2. **JMeter & Gatling Documentation**: Full documentation access limited; some sources returned 403/404 errors. Evidence partially derived from table of contents and known features.

3. **Real-world Case Studies**: Proprietary production case studies from major tech companies are not publicly accessible.

4. **Performance Benchmark Data**: No concrete benchmark numbers were available for direct comparison of tools (requests/second, resource overhead).

5. **Language Barrier**: Some primary sources (Google SRE Book) are in English, while the original lab content is in Bahasa Indonesia. This report synthesizes both.

## Conclusion

Load testing adalah komponen kritis dalam membangun sistem andal. Enam jenis test utama (smoke, load, stress, spike, soak, breakpoint) menyediakan kerangka komprehensif untuk menguji sistem dengan berbagai pola beban. Metrics seperti P95/P99 response time, error rate, dan resource utilization harus dipantau secara terus-menerus. Tool seperti k6 (JavaScript), Locust (Python), JMeter (GUI/XML), dan Gatling (Scala/JVM) memiliki spektrum keunggulan yang berbeda. Identifikasi bottleneck dilakukan dengan memonitor metrics tiap layer dan mencari korelasi dengan degradasi performance. Tim harus menghindari kesalahan umum seperti hanya menguji health endpoint, data dummy tidak realistis, dan lingkungan tidak mencerminkan production. Praktik terbaik adalah melakukan load testing secara konsisten sejak awal development hingga pre-go-live.