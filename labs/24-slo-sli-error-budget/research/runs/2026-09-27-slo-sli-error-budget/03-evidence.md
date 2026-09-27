## Evidence 1 — SLI Definition

Claim: SLI adalah ukuran kuantitatif aspek level layanan yang diberikan.
Evidence: "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Datadog docs, Prometheus docs (golden signals)
Notes: Definisi kanonis. Contoh: request latency, error rate, throughput, availability.

## Evidence 2 — Common SLI Types

Claim: SLI umum mencakup latency, error rate, throughput, availability; agregasi via rate/percentile.
Evidence: "Most services consider request latency... Other common SLIs include the error rate... and system throughput... raw data is collected over a measurement window and then turned into a rate, average, or percentile."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Datadog (metric-based SLI), Prometheus (golden signals)
Notes: Sumber tambahan: Datadog mendefinisikan SLI sebagai "quantitative measurement... a metric or aggregation of one or more monitors."

## Evidence 3 — SLI Bentuk Rasio Good/Total

Claim: SLI ideal berupa rasio good events / total events (0–100%).
Evidence: "we generally recommend treating the SLI as the ratio of two numbers: the number of good events divided by the total number of events."
Source: Google SRE Workbook Ch.2
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: Datadog "SLI is the sum of good events divided by the sum of total events" (metric-based SLO)
Notes: Memudahkan kalkulasi error budget dan tooling.

## Evidence 4 — SLA vs SLO

Claim: SLA mengandung konsekuensi (mis. denda); SLO tidak punya konsekuensi eksplisit.
Evidence: "SLAs are service level agreements: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs they contain... if there is no explicit consequence, then you are almost certainly looking at an SLO."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Datadog "SLA: explicit or implicit agreement... consequences for not meeting them"
Notes: "Most people really mean SLO when they say SLA."

## Evidence 5 — Availability Tabel (Time-based)

Claim: 99% = 7.2 jam/bulan (432 menit), 99.9% = 43.2 menit/bulan, 99.99% = 4.32 menit/bulan (259 detik), 99.999% = 25.9 detik/bulan.
Evidence: Google SRE Book Appendix A Table 1-1.
Source: Google SRE Book Appendix A
URL: https://sre.google/sre-book/availability-table/
Confidence: HIGH
Corroborated By: Google SRE Book Ch.3 (aggregate availability formula)
Notes: Topic specification menyebut "7 jam 18 menit" untuk 99% → selisih 6 menit karena basis perhitungan (30.44 hari vs 30 hari); bukan kontradiksi teknis, dijelaskan di contradictions.

## Evidence 6 — 100% Reliability Bukan Target

Claim: 100% reliability adalah target yang salah karena: (a) probabilitas kegagalan nonzero, (b) user tidak merasakan beda antara 99.99% vs 99.999% (device/network lebih lemah), (c) cost naik non-linear (~100x per nines), (d) change = sumber outage #1.
Evidence: "100% reliability is the wrong target... four reasons" + cost discussion.
Source: Google SRE Workbook Ch.2
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: Google SRE Book Ch.3 "cost does not increase linearly as reliability increments"
Notes: Konsisten lintas sumber.

## Evidence 7 — Error Budget = 1 − SLO

Claim: Error budget = 100% - SLO%; selisih antara target dan actual.
Evidence: "An error budget is 1 minus the SLO of the service. A 99.9% SLO service has a 0.1% error budget."
Source: Google SRE Book Ch.3
URL: https://sre.google/sre-book/embracing-risk/
Confidence: HIGH
Corroborated By: Datadog "error budget = 100% - SLO target percentage"; SRE Workbook Appendix B
Notes: Contoh valid: 1M request SLO 99.9% → budget 1000 errors.

## Evidence 8 — Error Budget sebagai Decision Tool

Claim: Budget tinggi → release lebih cepat; budget hampir habis → stop deployment berisiko, fokus reliability.
Evidence: "As long as there is error budget remaining, new releases can be pushed."
Source: Google SRE Book Ch.3
URL: https://sre.google/sre-book/embracing-risk/
Confidence: HIGH
Corroborated By: Datadog, SRE Workbook Appendix B (SLO Miss Policy)
Notes: Ini inti value proposition: mengubah reliability jadi alat keputusan, bukan sekadar dashboard.

## Evidence 9 — Burn Rate Alerting

Claim: Burn rate = laju konsumsi budget relatif SLO. Multi-window recommended: page 14.4x (2%/1h+5m), page 6x (5%/6h+30m), ticket 1x (10%/3d+6h). Short window = 1/12 long window.
Evidence: Table 5-8 + contoh PromQL expr.
Source: Google SRE Workbook Ch.5
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: Datadog "Burn rate indicators: red icon >6, yellow icon 1-6 in past 2 hours"
Notes: Datadog pakai 2 jam window vs Google rekomendasi 1h/6h/3d — bukan kontradiksi, hanya implementasi berbeda.

## Evidence 10 — Percentile bukan Average

Claim: Mean menyembunyikan tail latency; gunakan P50/P95/P99.
Evidence: "Most metrics are better thought of as distributions rather than averages" + Figure 4-1: typical 50ms tapi 5% request 20x lambat.
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Datadog "Percentiles show the shape of the distribution", Prometheus "focus on high percentile values"
Notes: User lebih suka latency stabil meski sedikit lebih tinggi daripada variansi tinggi.

## Evidence 11 — CPU/RAM Bukan SLO User

Claim: SLO baik ukur pengalaman user (request berhasil? cepat?); CPU/RAM = diagnostic signal.
Evidence: "Ideally, the SLI directly measures a service level of interest" + contoh AdSense vs AdWords latency tradeoff.
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Prometheus "alert on symptoms, not causes"; Datadog "SLI should be user-facing"
Notes: Prinsip ala Google: ukur output user, bukan input infra.

## Evidence 12 — SLO Berbeda per Endpoint

Claim: Criticality berbeda → SLO berbeda (mis. POST /payment/webhook > GET /report).
Evidence: "Not all requests are considered equal... bucketing untuk menambah label ke SLI dan menerapkan SLO berbeda" + Tabel 5-10 (CRITICAL/HIGH_FAST/HIGH_SLOW/LOW/NO_SLO).
Source: Google SRE Workbook Ch.5
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: Datadog "Different SLOs for different request classes"
Notes: Topic specification menyebut contoh ini → didukung.

## Evidence 13 — Four Golden Signals

Claim: Latency, Traffic, Errors, Saturation adalah empat sinyal utama monitoring.
Evidence: "The four golden signals of monitoring are latency, traffic, errors, and saturation. If you can only measure four metrics of your user-facing system, focus on these four."
Source: Google SRE Book Ch.6
URL: https://sre.google/sre-book/monitoring-distributed-systems/
Confidence: HIGH
Corroborated By: Datadog "Golden signals are latency, traffic, errors, and saturation"; Prometheus golden signals
Notes: Definisi konsisten di ekosistem.

## Evidence 14 — Error Budget Remaining Formula (Datadog)

Claim: error budget remaining = 100 * (current_status - target) / (100 - target)
Evidence: Datadog Documentation explicitly states this formula.
Source: Datadog SLO docs
URL: https://docs.datadoghq.com/service_level_objectives/
Confidence: MEDIUM (Datadog-specific implementation, mathematically sound)
Notes: Umum digunakan, tetapi sumber Google tidak menyebut formula eksplisit.

## Evidence 15 — SLO Miss Policy (Google Template)

Claim: Budget habis → halt changes kecuali P0/security; single incident >20% budget → postmortem + P0 action item.
Evidence: SRE Workbook Appendix B full policy.
Source: Google SRE Workbook Appendix B
URL: https://sre.google/workbook/error-budget-policy/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2 "Establishing an Error Budget Policy"
Notes: Template; perlu persetujuan PM/dev/SRE di organisasi masing-masing.

## Evidence 16 — Statistik "70% Outage from Change"

Claim: "Changes are a major source of instability, representing roughly 70% of our outages"
Evidence: SRE Workbook Appendix B Background section.
Source: Google SRE Workbook Appendix B
URL: https://sre.google/workbook/error-budget-policy/
Confidence: LOW (klaim internal Google tanpa metodologi publik di halaman tersebut)
Corroborated By: NOT VERIFIED (tidak ditemukan sumber independen dalam riset ini)
Notes: Gunakan dengan hati-hati; bukan fakta universal.

## Evidence 17 — Status Corrections

Claim: Status corrections allow exclude time periods (maintenance, non-business hours, deployments) dari SLO calculation.
Evidence: Datadog "Status corrections allow you to exclude specific time periods... Prevent expected downtime... Ignore non-business hours... Ensure that temporary issues caused by deployments do not negatively impact your SLOs."
Source: Datadog SLO docs
URL: https://docs.datadoghq.com/service_level_objectives/
Confidence: HIGH
Corroborated By: Google SRE Book Ch.4 "SLOs Set Expectations" (planned outages contoh Chubby)
Notes: Datadog formalisasi konsep yang disebut sekilas di Google SRE.

## Evidence 18 — Window SLO Recommendation

Claim: Rolling 4 minggu (28 hari) adalah interval general-purpose yang baik; lengkapi weekly summary + quarterly report.
Evidence: "We have found a four-week rolling window to be a good general-purpose interval."
Source: Google SRE Workbook Ch.2
URL: https://sre.google/workbook/implementing-slos/
Confidence: MEDIUM (rekomendasi praktik Google, bukan hasil eksperimen terkontrol)
Corroborated By: Datadog (7d, 30d, 90d options)
Notes: Topic specification pakai 30 hari; bukan kontradiksi, trade-off: 28 hari = weekend konsisten, 30 hari = lebih intuitive untuk bisnis bulanan.

## Evidence 19 — Metric vs Monitor vs Time-Slice SLO

Claim: Datadog mendefinisikan 3 tipe SLO: metric-based (good/total count), monitor-based (uptime dari monitor), time-slice (custom uptime definition).
Evidence: Datadog "SLO types: Metric-based, Monitor-based, Time Slice SLOs" + definition table.
Source: Datadog SLO docs
URL: https://docs.datadoghq.com/service_level_objectives/
Confidence: HIGH
Corroborated By: NOT VERIFIED di Google SRE (konsep serupa tapi tidak dibedakan dalam tipe eksplisit)
Notes: Implementasi vendor; prinsip SLI rasio tetap konsisten.

## Evidence 20 — Kalkulasi Latihan Lab (Validasi)

Claim: Payment Webhook SLO 99.99% dari 200.000 request → budget = 20 request; 15 gagal = 75% budget habis.
Evidence: Manual calculation: 0.01% × 200.000 = 20; 15/20 = 0.75 = 75%.
Source: Topic specification exercise (verified against SRE Book aggregate availability formula)
Confidence: HIGH
Corroborated By: Google SRE Book Ch.3 aggregate availability formula
Notes: Contoh valid; sesuai dengan konsep.

## Evidence 21 — Alert on Symptoms, Not Causes

Claim: Alerts should fire on user-visible symptoms (latency, errors, traffic), bukan infrastructure (CPU, memory).
Evidence: Prometheus "aim to have as few alerts as possible, by alerting on symptoms that are associated with end-user pain rather than trying to catch every possible way that pain could be caused."
Source: Prometheus Alerting docs
URL: https://prometheus.io/docs/practices/alerting/
Confidence: HIGH
Corroborated By: Google SRE Book Ch.6 "monitoring symptoms is easier the further 'up' your stack you monitor"
Notes: Ini prinsip ala Google; konsisten di ekosistem.

## Evidence 22 — Multi-dimensional SLOs

Claim: SLO bisa multiple thresholds (90% < 100ms, 99% < 400ms) untuk menangkap distribusi user experience.
Evidence: Google SRE Book Ch.4 contoh "90% Get RPC < 1ms, 99% Get RPC < 10ms, 99.9% Get RPC < 100ms".
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Datadog "Multiple target values per SLO"
Notes: Bukan kontradiksi dengan topic specification yang menyebut single SLO; ini extension best practice.
