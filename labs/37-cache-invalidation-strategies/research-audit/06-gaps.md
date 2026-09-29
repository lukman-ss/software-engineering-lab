# 06 — Research Gap Analysis

## Gap 1

Type: UNVERIFIED_CLAIM

Severity: HIGH

Location: `research/05-report.md: Finding 6`

Problem: Mathematical proofs and experimental benchmarks for XFetch optimality in Vattani et al. (2015) could not be verified directly from primary PDF (binary stream decompression failed).

Required Revision: Explicitly mark optimality proof as NOT VERIFIED from primary source text; rely on DOI metadata and secondary citations.

Can Be Approved Without Fix: YES (Revision updated `05-report.md` with explicit disclaimer).

---

## Gap 2

Type: FORMULA_ERROR

Severity: CRITICAL

Location: `research/05-report.md: Finding 6`

Problem: Lab prompt contained formula sign error `Δ·β·ln(rand()) > TTL_remaining` which evaluates to negative > positive for $rand \in (0,1)$, rendering early refresh inactive.

Required Revision: Provide correct formula $- \Delta \cdot \beta \cdot \ln(rand()) > \text{TTL\_remaining}$ with clear engineer-facing warning.

Can Be Approved Without Fix: YES (Revision added explicit WARNING block with correct formula).

---

## Gap 3

Type: OUTDATED_SOURCE / UNREACHABLE

Severity: MEDIUM

Location: `research/02-sources.md: Source 10`

Problem: Redis official documentation links returned 404/403 due to URL structure changes on redis.io.

Required Revision: Flag Redis-specific claims as NOT VERIFIED; rely on Microsoft Learn for general Cache-Aside vendor guidance.

Can Be Approved Without Fix: YES (Revision added REVISER UPDATE to Source 10).

---

## Gap 4

Type: OVERGENERALIZATION

Severity: LOW

Location: `research/05-report.md: Finding 4`

Problem: Synthetic lab test parameters (10,000 RPS, 500 goroutines, P99 benchmark) could be mistaken for empirical production benchmark claims.

Required Revision: Clarify that lab parameters are synthetic educational inputs, not empirical benchmarks.

Can Be Approved Without Fix: YES (Research already documented synthetic nature).

---

## Gap 5

Type: INFERENTIAL_CLAIM

Severity: MEDIUM

Location: `research/05-report.md: Finding 8`

Problem: Single-key stampede limitation of jitter was an un-cited inferential deduction.

Required Revision: Label claim explicitly as logical inference rather than cited fact.

Can Be Approved Without Fix: YES (Revision updated Finding 8 & Evidence 16).
