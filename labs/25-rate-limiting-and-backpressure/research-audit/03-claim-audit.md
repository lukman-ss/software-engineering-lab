# Claim Audit

## Claim 1: HTTP 429 is the Standard Rate Limiting Response
- **Claim**: HTTP 429 "Too Many Requests" status code is the standard mechanism for rate limiting responses, with optional `Retry-After` header.
- **Location**: `research/03-evidence.md:7`, `research/05-report.md:39`
- **Evidence Provided**: RFC 6585 Section 4.
- **Source**: RFC 6585 (IETF).
- **Source Actually Supports Claim**: YES
- **Classification**: FACT
- **Severity**: LOW
- **Notes**: Correctly notes servers are not obligated to emit 429 during DDoS / resource exhaustion and may drop connections directly.

---

## Claim 2: Token Bucket Algorithm Allows Burst Traffic while Maintaining Long-term Rate Limits
- **Claim**: Token bucket maintains tokens added at rate $r$ with maximum capacity $b$. Allows bursts while guaranteeing long-term average rate $r$.
- **Location**: `research/03-evidence.md:22`, `research/05-report.md:22`
- **Evidence Provided**: Formula $T_{max} = b / (M - r)$ for maximum burst duration at transmission rate $M$.
- **Source**: Wikipedia - Token Bucket / IEEE Datacenter Traffic Control.
- **Source Actually Supports Claim**: YES
- **Classification**: FACT
- **Severity**: LOW
- **Notes**: Standard networking formulation accurately described.

---

## Claim 3: Little's Law ($L = \lambda W$) Enables Queue Capacity Planning
- **Claim**: Long-term average number in system $L$ equals arrival rate $\lambda$ times average wait/sojourn time $W$ in stationary, ergodic systems.
- **Location**: `research/03-evidence.md:53`, `research/05-report.md:57`
- **Evidence Provided**: Mathematical proof references from Little (1961) and queueing theory.
- **Source**: Wikipedia - Little's Law / MIT Lecture Notes.
- **Source Actually Supports Claim**: YES
- **Classification**: FACT
- **Severity**: LOW
- **Notes**: Core theorem applies to queuing systems regardless of arrival distribution.

---

## Claim 4: Exponential Backoff with Jitter Prevents Retry Storms
- **Claim**: Exponential backoff without jitter causes synchronized retries. Adding Full Jitter ($\text{random}(0, \min(\text{cap}, 2^{\text{attempt}}))$) spreads calls evenly and reduces client work by $>50\%$ under 100 contending clients.
- **Location**: `research/03-evidence.md:69`, `research/05-report.md:73`
- **Evidence Provided**: Simulation results published by AWS Architecture Blog (Marc Brooker).
- **Source**: AWS Architecture Blog (2015/2023).
- **Source Actually Supports Claim**: YES
- **Classification**: FACT / EMPIRICAL STUDY
- **Severity**: LOW
- **Notes**: Distinction between Full Jitter, Equal Jitter, and Decorrelated Jitter is preserved accurately.

---

## Claim 5: Reactive Streams Standardizes Non-blocking Backpressure
- **Claim**: Reactive Streams protocol standardizes asynchronous stream processing with bounded buffers and demand-driven non-blocking backpressure (JEP 266 in Java 9).
- **Location**: `research/03-evidence.md:37`, `research/05-report.md:105`
- **Evidence Provided**: Reactive Streams specification history (2013-2015) and adoption by RxJava, Akka, Reactor.
- **Source**: Wikipedia - Reactive Streams / Reactive Streams Spec.
- **Source Actually Supports Claim**: YES
- **Classification**: FACT
- **Severity**: LOW
- **Notes**: Accurate historical and technical summary.

---

## Claim 6: Distributed Rate Limiting Requires Shared State or Coordination
- **Claim**: Rate limiting across distributed instances requires centralized storage (e.g. Redis) or sticky routing/distributed synchronization.
- **Location**: `research/03-evidence.md:211`, `research/05-report.md:139`
- **Evidence Provided**: Architectural trade-offs between centralized Redis vs sliding window log vs local approximations.
- **Source**: Medium (Figma) / Wikipedia.
- **Source Actually Supports Claim**: YES
- **Classification**: INTERPRETATION / IMPLEMENTATION-SPECIFIC
- **Severity**: LOW
- **Notes**: Correctly notes the trade-off between consistency, network overhead, and availability.

---

## Claim 7: Multi-Tenant Systems Require Fair Queueing and Per-Tenant Limits
- **Claim**: Without per-tenant rate limits or fair queueing, a single tenant's heavy workload (e.g. bulk batch jobs) can starve other tenants.
- **Location**: `research/03-evidence.md:227`, `research/05-report.md:89`
- **Evidence Provided**: Datacenter traffic control research and multi-tenant isolation patterns.
- **Source**: IEEE Paper / User Lab Doc.
- **Source Actually Supports Claim**: YES
- **Classification**: FACT / SYSTEM DESIGN PATTERN
- **Severity**: LOW
- **Notes**: Well established in multitenant architecture.
