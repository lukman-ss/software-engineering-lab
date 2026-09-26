# 03 Claim Audit

Target Lab: `labs/24-slo-sli-error-budget`  
Audit Date: 2026-09-26

---

## Claim 1
Claim: SLI adalah ukuran kuantitatif kualitas layanan; Google SRE merekomendasikan format rasio `good events / total events` agar bernilai 0–100%.  
Location: `05-report.md` §Finding 1; `03-evidence.md` §Evidence 1 & 8  
Evidence Provided: SRE Book Ch.4 §Indicators in Practice; SRE Workbook Ch.2 "What to Measure"  
Source: Sources 1 & 7 (Google SRE Book & Workbook)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurately reflected and cited verbatim from the source.

---

## Claim 2
Claim: SLO adalah target value/range untuk SLI (`SLI <= target` atau lower bound <= SLI <= upper bound); dapat memiliki multiple thresholds bertingkat (misal P90 dan P99).  
Location: `05-report.md` §Finding 2; `03-evidence.md` §Evidence 3  
Evidence Provided: SRE Book Ch.4 §Objectives; SRE Workbook Ch.2  
Source: Sources 1 & 7  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Matches official Google definitions and guidelines.

---

## Claim 3
Claim: SLA berbeda dari SLO; SLA memuat konsekuensi eksplisit (seperti penalti finansial/refund). Jika tanpa konsekuensi eksplisit, praktis adalah SLO.  
Location: `05-report.md` §Finding 3; `03-evidence.md` §Evidence 4  
Evidence Provided: SRE Book Ch.4 §Agreements  
Source: Source 1  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Precise quotation and clear distinction between contract vs internal target.

---

## Claim 4
Claim: Empat sinyal emas (The Four Golden Signals) monitoring user-facing sistem adalah Latency, Traffic, Errors, Saturation.  
Location: `05-report.md` §Finding 4; `03-evidence.md` §Evidence 10  
Evidence Provided: SRE Book Ch.6 "The Four Golden Signals"  
Source: Source 2  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Primary canonical definition by Google SRE. Widely adopted across the industry.

---

## Claim 5
Claim: Gunakan percentile (P50/P90/P99), bukan rata-rata (mean), untuk latency karena rata-rata menyembunyikan tail latency.  
Location: `05-report.md` §Finding 5; `03-evidence.md` §Evidence 9  
Evidence Provided: SRE Book Ch.4 §Aggregation & Ch.6 §Worrying About Your Tail  
Source: Sources 1 & 2  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Fully supported mathematically and empirically in the cited chapters.

---

## Claim 6
Claim: 100% availability adalah target yang salah karena tidak realistis, tidak bernilai tambah bagi user di edge, biaya melonjak eksponensial per nine, dan menghambat deployment/perubahan.  
Location: `05-report.md` §Finding 6; `03-evidence.md` §Evidence 6  
Evidence Provided: SRE Workbook Ch.2 "Reliability Targets"; SRE Book Ch.3 "Embracing Risk"  
Source: Sources 4 & 7  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: LOW  
Notes: Core philosophy of SRE. Fully backed by Chapter 3 & Workbook Chapter 2.

---

## Claim 7
Claim: Downtime availability nines: 99% ≈ 7.2 jam/bulan (30 hari), 99.9% = 43.2 menit/bulan, 99.99% = 4.32 menit/bulan, 99.999% = 5.26 menit/tahun.  
Location: `05-report.md` §Finding 7; `03-evidence.md` §Evidence 5  
Evidence Provided: SRE Book Appendix A Table 1-1  
Source: Source 5  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Math verified against Google Appendix A (30-day month baseline). Minor variation with 30.44-day calendar average noted and explained in contradictions.

---

## Claim 8
Claim: Error Budget = 1 − SLO; mengukur sisa ketidakandalan yang diperbolehkan dalam suatu jendela waktu.  
Location: `05-report.md` §Finding 8; `03-evidence.md` §Evidence 7  
Evidence Provided: SRE Book Ch.3; SRE Workbook Appendix B  
Source: Sources 4 & 8  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Exact mathematical identity used in SRE framework.

---

## Claim 9
Claim: Error budget berfungsi sebagai mekanisme pengambilan keputusan rilis/eksperimen vs stabilisasi/pembayaran utang teknis.  
Location: `05-report.md` §Finding 9; `03-evidence.md` §Evidence 7 & 15  
Evidence Provided: SRE Book Ch.3 §Forming Your Error Budget; SRE Workbook Appendix B  
Source: Sources 4 & 8  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: LOW  
Notes: Core organizational concept of Google SRE error budget policies.

---

## Claim 10
Claim: Alerting berbasis multi-window multi-burn-rate (page: 14.4x/1h+5m, 6x/6h+30m; ticket: 1x/3d+6h) memiliki precision dan recall tinggi dibandingkan alert durasi tunggal (`for: 1h`).  
Location: `05-report.md` §Finding 10; `03-evidence.md` §Evidence 11 & 12  
Evidence Provided: SRE Workbook Ch.5 Table 5-8 & PromQL examples  
Source: Source 6  
Source Actually Supports Claim: YES  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: MEDIUM  
Notes: Supported by SRE Workbook Ch.5. Specific parameters are recommendations tested inside Google; other teams/tools may tweak specific multipliers.

---

## Claim 11
Claim: SLO harus berbeda per endpoint berdasarkan criticality (misal `POST /payment/webhook` > `GET /report`).  
Location: `05-report.md` §Finding 11  
Evidence Provided: SRE Book Ch.4 §Objectives; SRE Workbook Ch.5 Table 5-10  
Source: Sources 1 & 6  
Source Actually Supports Claim: YES  
Classification: EXAMPLE  
Severity: LOW  
Notes: Realistic enterprise architectural guidance supported by Workbook bucketing guidance.

---

## Claim 12
Claim: Metrik infrastruktur (CPU, RAM, connection pool) adalah sinyal diagnostik, bukan SLI/SLO user-facing.  
Location: `05-report.md` §Finding 12  
Evidence Provided: SRE Book Ch.4 §Indicators in Practice; Ch.3  
Source: Sources 1 & 4  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: LOW  
Notes: Strongly aligned with SRE user-centric principles, correctly classified as an interpretation.

---

## Claim 13
Claim: Window rolling 4-minggu (28 hari) adalah interval umum yang direkomendasikan.  
Location: `03-evidence.md` §Evidence 16  
Evidence Provided: SRE Workbook Ch.2  
Source: Source 7  
Source Actually Supports Claim: YES  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: LOW  
Notes: Supported by SRE Workbook Ch.2 as Google's operational experience.

---

## Claim 14
Claim: ~70% gangguan (outages) disebabkan oleh perubahan (changes/deployments).  
Location: `03-evidence.md` §Evidence 17; `05-report.md` §Limitations  
Evidence Provided: SRE Workbook Appendix B Background  
Source: Source 8  
Source Actually Supports Claim: PARTIAL  
Classification: HYPOTHESIS  
Severity: MEDIUM  
Notes: Google mentions this as an internal rule of thumb in Appendix B without formal empirical publication. Research agent appropriately marked it LOW confidence and cautioned against treating it as an absolute universal statistic.
