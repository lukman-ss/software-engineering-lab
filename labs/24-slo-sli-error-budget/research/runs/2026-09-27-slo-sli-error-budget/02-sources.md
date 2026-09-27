## Source 1

Title: Chapter 4 - Service Level Objectives
Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/service-level-objectives/
Published: 2016 (O'Reilly SRE Book), continuously updated
Accessed: 2026-09-27
Source Tier: Tier 1 (Google internal primary source / official documentation)
Relevance: Definisi isomorph SLI/SLO/SLA, contoh implementasi, best practice pemilihan target.

## Source 2

Title: Chapter 3 - Embracing Risk
Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/embracing-risk/
Published: 2016
Accessed: 2026-09-27
Source Tier: Tier 1
Relevance: Error budget, availability table, risk tolerance layanan, cost non-linear kaunter reliability.

## Source 3

Title: Appendix A - Availability Table
Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/availability-table/
Published: 2016
Accessed: 2026-09-27
Source Tier: Tier 1
Relevance: Tabel downtime terstandarisasi per tahun/kuartal/bulan/minggu/hari/jam untuk banyak availability level (99%, 99.9%, 99.99%, 99.999%).

## Source 4

Title: Chapter 2 - Implementing SLOs
Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/implementing-slos/
Published: 2018
Accessed: 2026-09-27
Source Tier: Tier 1
Relevance: Step-by-step implementasi SLO, kalkulasi error budget, kebijakan release vs reliability, multi-window SLO, contoh kalkulasi nyata.

## Source 5

Title: Chapter 5 - Alerting on SLOs
Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/alerting-on-slos/
Published: 2018
Accessed: 2026-09-27
Source Tier: Tier 1
Relevance: Burn rate alerting, multi-window multi-burn-rate, precision/recall/detection time, rekomendasi parameter (14.4x/1h, 6x/6h, 1x/3d).

## Source 6

Title: Appendix B - Example Error Budget Policy
Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/error-budget-policy/
Published: 2018-02-19
Accessed: 2026-09-27
Source Tier: Tier 1
Relevance: Template kebijakan budget error, prosedur postmortem bila budget terpaksa >20%, escalation path, statistik "70% outage dari change".

## Source 7

Title: Concepts in Service Monitoring
Publisher: Google Cloud Documentation (Google Inc.)
URL: https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring
Published: 2026-09-25
Accessed: 2026-09-27
Source Tier: Tier 1 (official cloud provider doc)
Relevance: SLI sebagai rasio good/total events, error budget formula, burn rate konsep, jenis SLO (metric-based, monitor-based, time-slice).

## Source 8

Title: Alerting | Prometheus
Publisher: Prometheus Authors / The Linux Foundation
URL: https://prometheus.io/docs/practices/alerting/
Published: 2014-2026
Accessed: 2026-09-27
Source Tier: Tier 1 (official CNCF project documentation)
Relevance: Prinsip alerting symptom bukan cause, golden signals (latency, traffic, errors, saturation), perkiraan for duration vs burn rate.

## Source 9

Title: Service Level Objectives
Publisher: Datadog
URL: https://docs.datadoghq.com/service_level_objectives/
Published: 2026-09-27
Accessed: 2026-09-27
Source Tier: Tier 2 (reputable technical platform documentation)
Relevance: Definisi SLI/SLO/SLA/error budget yang identik dengan Google, error budget remaining formula, burn rate indicator (elevated 1-6, critical >6), SLO status corrections, tipe SLO (metric/monitor/time-slice).