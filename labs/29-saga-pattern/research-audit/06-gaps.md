# Research Gap Analysis: Saga Pattern Research

Target Lab: `labs/29-saga-pattern`
Date: 2026-09-29

---

## Gap 1

Type: WEAK_SOURCE

Severity: MEDIUM

Location:
`research/05-report.md` → Finding 4; `research/03-evidence.md` → Evidence 7

Problem:
The 3-tier transaction taxonomy (Compensable, Pivot, Retryable) is sourced exclusively from Microsoft Azure Architecture Center. Microservices.io (Richardson) does not enumerate this taxonomy on its public page. Richardson's book (*Microservices Patterns* Chapter 4, Manning, paywalled) contains the authoritative description but was not directly opened as a source. The Garcia-Molina 1987 paper does not introduce this language (these are modern framing terms).

Required Revision:
Either (a) acknowledge single-source limitation with stronger wording, or (b) obtain a cross-reference from a second non-Microsoft authoritative source (e.g., Richardson book page or Temporal documentation). The research does correctly self-report MEDIUM confidence.

Can Be Approved Without Fix: YES
The limitation is transparently acknowledged in both `03-evidence.md` (Evidence 7) and `06-open-questions.md`.

---

## Gap 2

Type: WEAK_SOURCE

Severity: MEDIUM

Location:
`research/05-report.md` → Finding 8; `research/03-evidence.md` → Evidence 12

Problem:
The enumeration of 6 isolation countermeasures (semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency) is sourced exclusively from Microsoft Azure Architecture Center. Microservices.io refers readers to *Microservices Patterns* Chapter 4/Section 4.3 (paywalled). The 6-item list as presented represents Microsoft's formulation, which may reflect Richardson's broader taxonomy but is not independently cross-verified from a second readable source.

Required Revision:
Label the countermeasure list as "Microsoft's taxonomy per Azure Architecture Center" rather than a universally agreed-upon canonical list. As presented, readers may interpret the list as a definitive industry standard.

Can Be Approved Without Fix: YES
Research reports MEDIUM confidence on this item. The underlying existence of countermeasures is corroborated by Richardson's public page.

---

## Gap 3

Type: UNVERIFIED_CLAIM

Severity: LOW

Location:
`research/03-evidence.md` → Evidence 1; `research/06-open-questions.md` → Open Question 3

Problem:
The verbatim definition from Garcia-Molina & Salem (1987) — commonly cited as "A saga is a long-lived transaction that can be written as a sequence of transactions that can be interleaved with other transactions" — was not extracted directly from the PDF due to LZW/scanned-image encoding. Verification relies on citation chains.

Required Revision:
None strictly required. The research transparently documents this limitation in `06-open-questions.md` and acknowledges that verbatim extraction was not performed.

Can Be Approved Without Fix: YES
Attribution is academically well-established. DOI 10.1145/62224.62226 verified.

---

## Gap 4

Type: MISSING_CASE

Severity: LOW

Location:
`research/05-report.md` → Conclusion; `research/06-open-questions.md` → Open Question 1

Problem:
No standard compensating-for-compensation recovery protocol is defined or cited. When a compensating transaction itself fails (e.g., Refund Payment fails during saga rollback), all sources agree recovery depends on retries, dead-letter queues, and human-in-the-loop. No canonical specification covers this scenario.

Required Revision:
This limitation is correctly noted in `06-open-questions.md`. Explicit acknowledgment in the lab README or conclusion section that compensation failure recovery is non-standardized would strengthen the research.

Can Be Approved Without Fix: YES

---

## Gap 5

Type: MISSING_SOURCE

Severity: LOW

Location:
`research/01-plan.md` → Risk / Unknowns; `research/06-open-questions.md` → Open Question 5

Problem:
Camunda, Axon, Seata (TCC pattern), and Spring State Machine were not opened as primary sources. The research plan acknowledged this scope limitation.

Required Revision:
If the lab will cover TCC vs Saga or tooling beyond AWS Step Functions / Temporal, primary sources for those tools should be added.

Can Be Approved Without Fix: YES
The lab focus is patterns and concepts, not an exhaustive tooling survey.

---

## Gap 6

Type: MISSING_CASE

Severity: LOW

Location:
`research/05-report.md` → Limitations

Problem:
No quantitative benchmark comparing 2PC latency vs Saga latency was found or cited. The claim "2PC too slow or not feasible" rests on qualitative engineering consensus across authoritative sources (Richardson, Microsoft, AWS). This is appropriate for a conceptual education lab but should not be presented as quantitatively benchmarked.

Required Revision:
None strictly required. Research correctly reports this in `05-report.md` Limitations: "Tidak ada benchmark kuantitatif latency/throughput saga vs 2PC yang diverifikasi."

Can Be Approved Without Fix: YES

---

## Gap 7

Type: IMPLEMENTATION_GAP (Research-level)

Severity: LOW

Location:
`research/05-report.md` → Conclusion

Problem:
The Microservices.io Idempotent Consumer URL (`/patterns/data/idempotent-consumer.html`) is listed as a source in `02-sources.md` (Source 8) and cited in `03-evidence.md` (Evidence 10) but was not directly opened and fetched as a primary inspection step in this audit session. The core idempotency claim is corroborated by Debezium and Microsoft, but the specific `PROCESSED_MESSAGES` table + `(subscriberId, messageID)` PK mechanism attributed to that URL cannot be independently page-confirmed from the current audit.

Required Revision:
Low priority. Mechanism is consistent with the architectural pattern described in the Debezium blog and Microsoft docs.

Can Be Approved Without Fix: YES

---

## Summary

| Gap | Type | Severity | Blocks Approval |
|-----|------|----------|-----------------|
| 1 — Pivot/Retryable single source | WEAK_SOURCE | MEDIUM | NO |
| 2 — 6 countermeasures single source | WEAK_SOURCE | MEDIUM | NO |
| 3 — 1987 PDF verbatim unextracted | UNVERIFIED_CLAIM | LOW | NO |
| 4 — Compensation failure recovery undefined | MISSING_CASE | LOW | NO |
| 5 — TCC/Camunda/Axon not sourced | MISSING_SOURCE | LOW | NO |
| 6 — No quantitative 2PC benchmark | MISSING_CASE | LOW | NO |
| 7 — Idempotent Consumer URL not directly inspected | IMPLEMENTATION_GAP | LOW | NO |
