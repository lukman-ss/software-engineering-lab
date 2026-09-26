# 04 Contradictions Audit

## Overview
Inspection conducted across `01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md`, and `06-open-questions.md`.

---

## Contradiction Analysis 1: Real Calls vs. Mocks for External Dependencies
- **Statement A:** External third-party API dependencies must be tested with real calls to expose real-world latency (Azure Well-Architected / `03-evidence.md` Evidence 17).
- **Statement B:** High-volume stress testing with real calls to external third parties risks rate limits, throttling, cost, and terms of service violations (`04-contradictions.md` Item 2).
- **Type:** SOURCE_CONFLICT / METHODOLOGICAL_TRADEOFF
- **Impact:** LOW
- **Assessment:** Resolved. The research explicitly synthesizes a phased approach: use real sandbox endpoints for baseline average-load tests and latency validation; use high-fidelity stubs/mocks with simulated latency for extreme stress and spike testing.

---

## Contradiction Analysis 2: Staging Testing vs Production Testing
- **Statement A:** Test environments must mirror production as closely as practical (`03-evidence.md` Evidence 11).
- **Statement B:** Testing in staging can never fully replicate production traffic dynamics; production testing exposes problems that only surface under actual usage (`05-report.md` Finding 10, Azure docs).
- **Type:** METHODOLOGICAL_TRADEOFF
- **Impact:** LOW
- **Assessment:** Resolved. The research demonstrates a tiered strategy: staging tests validate baseline capacity safely; synthetic monitoring and controlled progressive canary tests in production validate actual live conditions.

---

## Summary
No material unresolved contradictions found. All apparent conflicts are documented and reconciled with technical context.
