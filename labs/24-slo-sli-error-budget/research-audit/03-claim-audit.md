# 03 Claim Audit

## Claim 1
Claim: An SLI is a quantitative measure of service level provided; an SLO is a target value or range; an SLA is an agreement with business consequences; an error budget is the allowable failure rate (100% - SLO).  
Location: `05-report.md` (Finding 1), `03-evidence.md` (Evidence 1, 2, 3)  
Evidence Provided: Direct quotes from Google SRE Book Ch.4 & SRE Workbook Ch.2.  
Source: Google SRE Book / SRE Workbook  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Well-established foundational definitions supported across industry.

## Claim 2
Claim: SLIs should be user-centric (good events / total events) rather than internal infrastructure metrics (e.g. CPU/RAM); latency SLIs should use percentiles (P95/P99) rather than averages.  
Location: `05-report.md` (Finding 2), `03-evidence.md` (Evidence 4, 5)  
Evidence Provided: SRE Book Ch.4 text and distribution explanations.  
Source: Google SRE Book Ch.4, Prometheus Alerting Best Practices  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Averages obscure tail latency outliers.

## Claim 3
Claim: Error Budget Remaining % formula is `100 * (current status - target) / (100 - target)`.  
Location: `05-report.md` (Finding 3), `03-evidence.md` (Evidence 7)  
Evidence Provided: Datadog documentation formula and SRE Workbook ratio definitions.  
Source: Datadog SLO Documentation, Google SRE Workbook Ch.2  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Arithmetic verified.

## Claim 4
Claim: 100% availability is the wrong reliability target because marginal utility approaches zero while costs escalate exponentially, and prevents deployment iteration.  
Location: `05-report.md` (Finding 4), `03-evidence.md` (Evidence 9)  
Evidence Provided: SRE Book Ch.3 "Embracing Risk" analysis.  
Source: Google SRE Book Ch.3  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Core philosophical principle of SRE.

## Claim 5
Claim: Availability "nines" downtime allowances are 99% (~7.2h/mo), 99.9% (~43.2m/mo), 99.99% (~4.32m/mo), 99.999% (~25.9s/mo) under 30-day month convention.  
Location: `05-report.md` (Finding 5), `03-evidence.md` (Evidence 8)  
Evidence Provided: Google SRE Book Appendix A Availability Table.  
Source: Google SRE Book Appendix A  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Noted minor discrepancy if using 30.44-day average month, which is explicitly reconciled in contradictions.

## Claim 6
Claim: Recommended SLO alerting uses multi-window multi-burn-rate parameters (e.g., 2% budget in 1h @ 14.4x burn rate, 5% in 6h @ 6x burn rate).  
Location: `05-report.md` (Finding 6), `03-evidence.md` (Evidence 6)  
Evidence Provided: Google SRE Workbook Ch.5 Table 5-8.  
Source: Google SRE Workbook Ch.5  
Source Actually Supports Claim: YES  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: MEDIUM  
Notes: Properly categorized as recommended starting points requiring contextual tuning.

## Claim 7
Claim: OpenSLO provides an open, vendor-neutral declarative YAML specification for SLOs.  
Location: `05-report.md` (Finding 1 & Open Questions), `03-evidence.md` (Evidence 10)  
Evidence Provided: OpenSLO specification schema repository and site.  
Source: OpenSLO Specification (openslo.github.io)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Ecosystem adoption breadth noted in open questions.

## Claim 8
Claim: Organizations like The Home Depot successfully scaled SLOs to 800+ services using the VALET framework (Volume, Availability, Latency, Errors, Tickets).  
Location: `05-report.md` (Finding 7), `03-evidence.md` (Evidence 11)  
Evidence Provided: Case study in Google SRE Workbook Ch.3.  
Source: Google SRE Workbook Ch.3  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Real-world organizational case study.
