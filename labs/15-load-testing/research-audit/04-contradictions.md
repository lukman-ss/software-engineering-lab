# Contradictions Audit

## Contradiction 1
Statement A: "Mocking external dependencies makes tests run faster... but it hides real-world performance problems."
Location: Source 15 (Azure Well-Architected Framework)
Statement B: High-volume stress testing real external APIs (like WhatsApp) can lead to rate limiting, account bans, and excessive cost.
Location: Topic Spec / Practical Engineering Experience (Synthesized in Evidence 17 & Finding 7)
Type: SOURCE_CONFLICT (Theory vs Practice)
Impact: MEDIUM
Assessment: The research report resolves this adequately by stating that real external APIs should be used in *controlled sandbox* tests to find end-to-end latency, but high-fidelity stubs/mocks must be used for massive stress/spike tests. 

## Overall Assessment
No material, unresolved contradictions found. The Research Agent properly identified and contextualized discrepancies (e.g., k6 vs Azure managed service, production vs staging fidelity).
