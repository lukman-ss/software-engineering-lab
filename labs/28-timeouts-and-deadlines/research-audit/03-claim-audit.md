# Claim Audit: Lab 28 (Timeouts & Deadlines)

## Claim 1
Claim: Single slow dependencies cause cascading failures through resource exhaustion (threads, CPU, memory, socket descriptors).
Location: `05-report.md: Finding 1`, `03-evidence.md: Evidence 1`
Evidence Provided: Google SRE Book - Addressing Cascading Failures (Chapter 22)
Source: Source 1 (https://sre.google/sre-book/addressing-cascading-failures/)
Source Actually Supports Claim: YES
Classification: FACT
Severity: CRITICAL (Core architecture principle)
Notes: Fully verified by SRE Chapter 22. Slow dependencies hold worker threads/connections, exhausting upstream pool capacity and causing domino failures.

---

## Claim 2
Claim: Retries at multiple layers multiply exponentially (e.g. 4 layers with 3 retries yields 4^3 = 64 attempts on the bottom storage layer).
Location: `05-report.md: Finding 3`, `03-evidence.md: Evidence 2`
Evidence Provided: Google SRE Book - Addressing Cascading Failures
Source: Source 1 (https://sre.google/sre-book/addressing-cascading-failures/)
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH
Notes: Exact example matches Google SRE Book verbatim.

---

## Claim 3
Claim: Absolute deadlines must be propagated across service calls (subtracting elapsed processing time) to avoid downstream servers performing work after the client has given up.
Location: `05-report.md: Finding 2`, `03-evidence.md: Evidence 3`
Evidence Provided: gRPC Documentation - Deadlines
Source: Source 2 (https://grpc.io/docs/guides/deadlines/)
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH
Notes: gRPC documentation and SRE Book both confirm deadline propagation mechanics and remaining-timeout conversion to prevent clock-skew issues.

---

## Claim 4
Claim: Retrying non-idempotent operations without deduplication keys causes double execution (e.g., double billing); atomic persistence of deduplication markers and business state is required.
Location: `05-report.md: Finding 5`, `03-evidence.md: Evidence 6`
Evidence Provided: Azure Architecture Center - Idempotent Consumer Pattern
Source: Source 5 (https://learn.microsoft.com/en-us/azure/architecture/patterns/idempotent-consumer)
Source Actually Supports Claim: YES
Classification: FACT
Severity: CRITICAL
Notes: Directly supported by Azure patterns and RabbitMQ reliability documentation.

---

## Claim 5
Claim: Database statement timeouts (`statement_timeout`, `lock_timeout`) are necessary guardrails to abort runaway queries and unblock connection pools.
Location: `05-report.md: Finding 6`, `03-evidence.md: Evidence 4 & 13`
Evidence Provided: PostgreSQL Documentation
Source: Source 3 & 13 (https://www.postgresql.org/docs/current/runtime-config-client.html)
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH
Notes: Fully supported by official PostgreSQL documentation.

---

## Claim 6
Claim: Retries must use exponential backoff with jitter to prevent synchronized retry ripples/thundering herds.
Location: `05-report.md: Finding 3`, `03-evidence.md: Evidence 7`
Evidence Provided: AWS Architecture Blog - Exponential Backoff And Jitter
Source: Source 7 (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH
Notes: Verified by AWS Architecture blog and SRE Chapter 22.

---

## Claim 7
Claim: Health check probes must enforce strict individual timeouts on dependency checks and keep liveness probes decoupled from external dependencies.
Location: `05-report.md: Finding 8`, `03-evidence.md: Evidence 10`
Evidence Provided: Azure Architecture Center - Health Endpoint Monitoring Pattern
Source: Source 10 (https://learn.microsoft.com/en-us/azure/architecture/patterns/health-endpoint-monitoring)
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH
Notes: Verified by Azure pattern documentation.

---

## Claim 8
Claim: Dynamic/adaptive timeout tuning using machine learning is standard modern practice.
Location: `02-sources.md: Source 4`, `06-open-questions.md: 4`
Evidence Provided: Azure Circuit Breaker Pattern text mentioning AI/ML
Source: Source 4 (https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker)
Source Actually Supports Claim: PARTIAL
Classification: HYPOTHESIS / OVERGENERALIZED
Severity: MEDIUM
Notes: Azure doc mentions modern adaptive trends in passing, but static P99+margin configuration remains standard engineering practice. Correctly flagged in research `06-open-questions.md`.
