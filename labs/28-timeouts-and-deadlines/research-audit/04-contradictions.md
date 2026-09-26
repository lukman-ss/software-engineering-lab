# Contradiction Audit: Lab 28 (Timeouts & Deadlines)

## Contradiction 1
Statement A: Immediate retry on transient errors (Azure Retry Pattern).
Statement B: Do not retry immediately during overload; return error immediately or backoff (Google SRE Handling Overload).
Type: SOURCE_CONFLICT / SCOPE_DIFFERENCE
Impact: Misconfigured retry logic can crash an already struggling backend.
Assessment: RESOLVED in research. The research correctly distinguishes packet corruption/single-connection transient blips from overload/saturation errors where immediate retry causes cascading meltdowns.

---

## Contradiction 2
Statement A: Lower layers should fail fast and only top layer retries (Azure Retry Pattern).
Statement B: Retries should only happen at the immediate parent layer above the rejecting dependency, propagating "don't retry" upward (Google SRE).
Type: INTERNAL / PATTERN_VARIATION
Impact: Minor implementation variance.
Assessment: RESOLVED in research. Both avoid multi-layer multiplicative retry explosions.

---

## Contradiction 3
Statement A: gRPC example code shows up to 10 retry attempts.
Statement B: Google SRE and Resilience4j mandate max 3 retry attempts and <10% client retry ratio.
Type: INTERNAL
Impact: Naive code can cause retry storm.
Assessment: RESOLVED in research. The 10-attempt code in SRE Book was explicitly presented as an anti-pattern demonstration of naive retry loops.

---

## Summary
No material unresolved contradictions remain in the research documentation. All key tensions between sources have been properly contextualized and reconciled.
