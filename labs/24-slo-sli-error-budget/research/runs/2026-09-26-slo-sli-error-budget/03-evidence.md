# Evidence

## Evidence 1
Claim: SLI adalah ukuran kuantitatif aspek level layanan yang diberikan.
Evidence: "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."
Source: Google SRE Book Ch.4 Service Level Objectives
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2 (SLI as ratio good/total events)
Notes: Definisi kanonis. Contoh umum: request latency, error rate, throughput, availability.

## Evidence 2
Claim: SLI umum mencakup latency, error rate, throughput, availability; agregasi via rate/average/percentile.
Evidence: "Most services consider request latency... Other common SLIs include the error rate... and system throughput... raw data is collected over a measurement window and then turned into a rate, average, or percentile."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2 Table 2-1 (availability, latency, quality, freshness, correctness, coverage, durability)
Notes: Workbook memperluas taksonomi per tipe komponen (request-driven, pipeline, storage).

## Evidence 3
Claim: SLO adalah target value/range untuk SLI; struktur SLI ≤ target.
Evidence: "An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI. A natural structure for SLOs is thus SLI ≤ target, or lower bound ≤ SLI ≤ upper bound."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2
Notes: Contoh: 99% Get RPC < 100ms.

## Evidence 4
Claim: SLA berbeda dari SLO: SLA adalah kontrak dengan konsekuensi; jika tidak ada konsekuensi eksplisit maka itu SLO.
Evidence: "SLAs are service level agreements: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs... if there is no explicit consequence, then you are almost certainly looking at an SLO."
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2, Appendix B Error Budget Policy
Notes: "Most people really mean SLO when they say SLA."

## Evidence 5
Claim: Tabel downtime: 99% = 3.65 hari/tahun, 7.2 jam/bulan; 99.9% = 8.76 jam/tahun, 43.2 menit/bulan; 99.99% = 52.6 menit/tahun, 4.32 menit/bulan.
Evidence: Table 1-1 Availability table rows for 99%, 99.9%, 99.99% per year/month/week/day.
Source: Google SRE Book Appendix A
URL: https://sre.google/sre-book/availability-table/
Confidence: HIGH
Corroborated By: SRE Book Ch.3 (99.99% = up to 52.56 min/year downtime)
Notes: Asumsi no planned downtime. Topik lab menulis 99% → ~7 jam 18 menit/bulan (vs 7.2 jam = 7 jam 12 menit di tabel) — selisih kecil, lihat contradictions.

## Evidence 6
Claim: 100% reliability adalah target yang salah.
Evidence: "Our experience has shown that 100% reliability is the wrong target" + 4 alasan: probabilitas kegagalan nonzero, rantai ke user tidak 100%, marginal utility menurun sementara cost naik, change (#1 sumber outage) menjadi mustahil.
Source: Google SRE Workbook Ch.2 Implementing SLOs
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: SRE Book Ch.3 Embracing Risk (cost tidak linear, 100x per increment; user di 99% smartphone tidak bedakan 99.99% vs 99.999%)
Notes: Inti argumen cost/benefit.

## Evidence 7
Claim: Error Budget = 1 − SLO; selisih antara SLO dan uptime aktual dalam periode.
Evidence: "An error budget is 1 minus the SLO of the service. A 99.9% SLO service has a 0.1% error budget." + "The difference between these two numbers is the budget of how much unreliability is remaining for the quarter."
Source: SRE Workbook Appendix B + SRE Book Ch.3
URL: https://sre.google/workbook/error-budget-policy/
URL: https://sre.google/sre-book/embracing-risk/
Confidence: HIGH
Corroborated By: Kedua sumber independen dalam korpus Google SRE setuju.
Notes: Contoh: 1M request/4 minggu pada SLO 99.9% → budget 1.000 errors.

## Evidence 8
Claim: Rekomendasi bentuk SLI sebagai rasio good events / total events (0–100%), memudahkan error budget dan tooling.
Evidence: "we generally recommend treating the SLI as the ratio of two numbers: the number of good events divided by the total number of events."
Source: SRE Workbook Ch.2
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.5 (error budget/error rate berlaku untuk semua SLI rasio ini)
Notes: Contoh: successful HTTP / total HTTP.

## Evidence 9
Claim: Gunakan percentile (P50/P90/P99/P99.9), bukan mean/average, untuk latency.
Evidence: "Most metrics are better thought of as distributions rather than averages" + Figure 4-1 menunjukkan typical 50ms tapi 5% request 20x lebih lambat; average menyembunyikan tail.
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Book Ch.6 (mean 100ms @1000rps bisa sembunyikan 1% request 5 detik; kumpulkan histogram bucketed)
Notes: User lebih suka sistem sedikit lambat tapi variansi rendah.

## Evidence 10
Claim: Four Golden Signals: Latency, Traffic, Errors, Saturation.
Evidence: "The four golden signals of monitoring are latency, traffic, errors, and saturation. If you can only measure four metrics of your user-facing system, focus on these four."
Source: Google SRE Book Ch.6
URL: https://sre.google/sre-book/monitoring-distributed-systems/
Confidence: HIGH
Corroborated By: NOT VERIFIED untuk sumber kedua independen (tidak dicari di luar korpus Google dalam riset ini)
Notes: Definisi tiap sinyal didokumentasikan rinci di sumber yang sama.

## Evidence 11
Claim: Burn rate = kecepatan konsumsi error budget relatif terhadap SLO; burn rate 1 menghabiskan tepat habis di akhir window.
Evidence: "Burn rate is how fast, relative to the SLO, the service consumes the error budget... With an SLO of 99.9% over 30 days, a constant 0.1% error rate uses exactly all of the error budget: a burn rate of 1."
Source: SRE Workbook Ch.5 Alerting on SLOs
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: NOT VERIFIED di luar korpus (konsep berasal dari sumber ini)
Notes: Tabel: burn 2 → 15 hari, burn 10 → 3 hari, burn 1000 (100% outage) → 43 menit untuk SLO 99.9%/30 hari.

## Evidence 12
Claim: Rekomendasi alerting multiwindow multi-burn-rate: page 14.4x/1h+5m (2%), page 6x/6h+30m (5%), ticket 1x/3d+6h (10%).
Evidence: Table 5-8 Recommended parameters + contoh PromQL expr dengan long+short window AND.
Source: SRE Workbook Ch.5
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: SRE Book Ch.10 (alert error ratio > threshold + minimum duration for 2m untuk hindari flap)
Notes: Short window = 1/12 long window. Jangan pakai `for: 1h` duration saja (poor recall).

## Evidence 13
Claim: Prinsip pemilihan SLO: jangan berdasar performa saat ini, keep it simple, hindari absolut, sesedikit mungkin, perfection can wait.
Evidence: Lima bullet "Choosing Targets" di Ch.4.
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2 (start dengan yang mudah diukur, iterasi; aspirational SLO)
Notes: -

## Evidence 14
Claim: Gunakan safety margin (internal SLO lebih ketat dari eksternal) dan jangan overachieve (user bergantung pada realita, bukan janji; contoh Chubby planned outage).
Evidence: "Using a tighter internal SLO than the SLO advertised to users..." + "If your service's actual performance is much better than its stated SLO, users will come to rely on its current performance." + Chubby case.
Source: Google SRE Book Ch.4
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: MEDIUM (satu sumber otoritatif, studi kasus internal Google)
Corroborated By: NOT VERIFIED independen
Notes: -

## Evidence 15
Claim: Error budget policy konkret: jika budget 4-mingguan habis → halt changes kecuali P0/security; single incident >20% budget → postmortem wajib + P0 action item.
Evidence: "If the service has exceeded its error budget for the preceding four-week window, we will halt all changes..." + "If a single incident consumes more than 20% of error budget over four weeks, then the team must conduct a postmortem."
Source: SRE Workbook Appendix B
URL: https://sre.google/workbook/error-budget-policy/
Confidence: HIGH
Corroborated By: SRE Workbook Ch.2 (Establishing an Error Budget Policy; freeze/deprioritize eksternal)
Notes: Template, bukan standar universal; perlu persetujuan PM/dev/SRE.

## Evidence 16
Claim: Rolling 4-minggu adalah interval general-purpose yang baik; lengkapi ringkasan mingguan + laporan kuartalan.
Evidence: "We have found a four-week rolling window to be a good general-purpose interval. We complement this with weekly summaries... and quarterly summarized reports..."
Source: SRE Workbook Ch.2
URL: https://sre.google/workbook/implementing-slos/
Confidence: MEDIUM (rekomendasi praktik Google, bukan hasil eksperimen terkontrol)
Corroborated By: NOT VERIFIED independen
Notes: Rolling selaras user experience; calendar selaras business planning.

## Evidence 17
Claim: Perubahan adalah sumber outage #1 (~70%).
Evidence: "Changes are a major source of instability, representing roughly 70% of our outages"
Source: SRE Workbook Appendix B Background
URL: https://sre.google/workbook/error-budget-policy/
Confidence: LOW (klaim internal Google tanpa metodologi publik di halaman tersebut)
Corroborated By: NOT VERIFIED
Notes: Jangan kutip sebagai statistik universal.
