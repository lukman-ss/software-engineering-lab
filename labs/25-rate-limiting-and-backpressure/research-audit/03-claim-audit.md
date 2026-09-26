# Claim Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

---

## Claim 1: HTTP 429 and Retry-After Standardization

Claim: HTTP 429 Too Many Requests and Retry-After header are standardized by IETF for rate limiting.  
Location: `research/02-findings.md` §Claim 1, `research/05-report.md` §2.1  
Evidence Provided: RFC 6585 Section 4, RFC 9110 Section 10.2.3 and 15.5.20.  
Source: Sources 1 & 2 (RFC 6585, RFC 9110)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurate. Notes delay-seconds and HTTP-date formats.

---

## Claim 2: Token Bucket vs Leaky Bucket Mechanics

Claim: Token Bucket allows bursts up to capacity $B$ while replenishing at rate $R$, whereas Leaky Bucket smooths traffic to a strict constant output rate.  
Location: `research/02-findings.md` §Claim 2, `research/05-report.md` §2.2  
Evidence Provided: RFC 2697 Section 2 (srTCM).  
Source: Source 3 (RFC 2697)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurately distinguishes burst tolerance between Token Bucket and Leaky Bucket.

---

## Claim 3: Little's Law vs Deterministic Queue Buildup

Claim: Little's Law is $L = \lambda W$ (average inventory, average throughput, average residence time in steady state). It is mathematically distinct from the deterministic fluid accumulation formula $\Delta Q = (r_{in} - r_{out})\Delta t$.  
Location: `research/02-findings.md` §Claim 3, `research/05-report.md` §3.1 & §3.2  
Evidence Provided: J.D.C. Little (1961) Operations Research 9(3):383-387.  
Source: Source 5 (Little 1961)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: HIGH (Critical distinction rigorously resolved)  
Notes: The research explicitly corrects previous misattributions and provides precise operational boundaries.

---

## Claim 4: Unbounded Queue Operational Risks

Claim: An unbounded queue does not prevent system failure; it shifts failure mode from immediate rejection to memory exhaustion (OOM), age degradation, timeout cascades, and head-of-line blocking.  
Location: `research/02-findings.md` §Claim 4, `research/05-report.md` §3.3  
Evidence Provided: Google SRE Book, Chapter 22 §22.10.  
Source: Source 7 (Google SRE Book)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: MEDIUM  
Notes: Directly supported by SRE cascade failure literature.

---

## Claim 5: Exponential Backoff with Full Jitter

Claim: Full Jitter ($sleep = random(0, \min(cap, base \cdot 2^{attempt}))$) mitigates thundering herd and reduces total client work compared to unjittered exponential backoff.  
Location: `research/02-findings.md` §Claim 5, `research/05-report.md` §4.2  
Evidence Provided: AWS Architecture Blog (Marc Brooker, 2015).  
Source: Source 6 (AWS Blog)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Formula and trade-offs match authoritative reference.

---

## Claim 6: IP-Based Rate Limiting Limitations

Claim: IP-based rate limiting unfairly throttles shared NAT / CGNAT subscribers masked under shared IPv4 ranges (RFC 6598 100.64.0.0/10). Multi-tenant systems should limit by API key/tenant token.  
Location: `research/02-findings.md` §Claim 6, `research/05-report.md` §5  
Evidence Provided: RFC 6598 Section 2.  
Source: Source 4 (RFC 6598)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Correctly contextualizes limitations for public consumer web vs authenticated API traffic.

---

## Claim 7: Queue Age as Early Warning Indicator

Claim: Queue age (`ApproximateAgeOfOldestMessage` in SQS, message age / lag in Kafka) directly measures latency SLA violation risk and is a superior early warning signal compared to queue depth alone, though both are complementary.  
Location: `research/02-findings.md` §Claim 7, `research/05-report.md` §6.1 & §6.2  
Evidence Provided: AWS SQS CloudWatch Metrics documentation, Apache Kafka Monitoring documentation.  
Source: Sources 8 & 9 (AWS SQS, Kafka docs)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: MEDIUM  
Notes: Metric definitions and operational significance are verified.

---

## Claim 8: Fair Queuing and Multi-Tenant Isolation

Claim: Fair queuing (FQ, WFQ, WF2Q+) provides equitable bandwidth distribution across tenants, though specific weight parameters and queue structures are service-specific.  
Location: `research/02-findings.md` §Claim 8, `research/05-report.md` §8  
Evidence Provided: Nagle (1987), Jiang et al. (2005).  
Source: Academic literature referenced in findings.  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: LOW  
Notes: Properly labeled as INTERPRETATION due to implementation-specific weight configurations.

---

## Claim 9: Autoscaling Downstream Saturation Risks

Claim: Autoscaling upstream consumers/workers without scaling downstream dependencies (DB, cache, downstream APIs) leads to downstream saturation and cascade failure.  
Location: `research/02-findings.md` §Claim 9, `research/05-report.md` §7  
Evidence Provided: Google SRE Book, Chapter 22 §22.2.  
Source: Source 7 (Google SRE Book)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: MEDIUM  
Notes: Framed accurately as an architectural failure pattern rather than an invariant physical law.

---

## Claim 10: Retry Amplification Across Multi-Tier Systems

Claim: Multi-tier retries cause exponential attempt amplification (e.g. 3 retries at 3 layers = $4^3 = 64$ attempts). Mitigations include per-process retry budgets, avoiding multi-tier retries, and non-retriable error codes.  
Location: `research/02-findings.md` §Claim 10, `research/05-report.md` §4.3  
Evidence Provided: Google SRE Book §22.8, AWS Architecture Blog.  
Source: Sources 6 & 7 (AWS Blog, Google SRE Book)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: MEDIUM  
Notes: Supported by SRE retry dynamics principles.
