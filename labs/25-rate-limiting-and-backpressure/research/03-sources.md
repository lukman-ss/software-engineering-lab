# Source Index

## Primary Sources (Verified)

### RFC Documents

| # | Title | Publisher | URL | Status | Relevance |
|---|-------|-----------|-----|--------|-----------|
| 1 | RFC 6585 - Additional HTTP Status Codes (Section 4: 429 Too Many Requests) | IETF | https://www.rfc-editor.org/rfc/rfc6585.html | VERIFIED | HTTP rate limiting standard |
| 2 | RFC 9110 - HTTP Semantics (Section 10.2.3 / 15.5.20) | IETF | https://www.rfc-editor.org/rfc/rfc9110.html | VERIFIED | HTTP header semantics |
| 3 | RFC 2697 - A Single Rate Three Color Marker (srTCM) | IETF | https://www.rfc-editor.org/rfc/rfc2697.html | VERIFIED | Token bucket specification |
| 4 | RFC 6598 - Shared Address Space (100.64.0.0/10 for CGNAT) | IETF | https://www.rfc-editor.org/rfc/rfc6598.html | VERIFIED | IP rate limiting limitation |

### Academic Sources

| # | Title | Publisher | URL | Status | Relevance |
|---|-------|-----------|-----|--------|-----------|
| 5 | A Proof for the Queuing Formula: L = λW | J.D.C. Little, Operations Research (1961) | https://doi.org/10.1287/opre.9.3.383 | VERIFIED | Little's Law formal proof |

### Vendor Engineering References

| # | Title | Publisher | URL | Status | Relevance |
|---|-------|-----------|-----|--------|-----------|
| 6 | Exponential Backoff And Jitter | AWS Architecture Blog (Marc Brooker, 2015) | https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/ | VERIFIED | Retry strategy best practices |
| 7 | Google Site Reliability Engineering / SRE Workbook | Google / O'Reilly | https://sre.google/sre-book/addressing-cascading-failures/ | VERIFIED | Cascading failures, queue management, retries |

### Broker Documentation

| # | Title | Publisher | URL | Status | Relevance |
|---|-------|-----------|-----|--------|-----------|
| 8 | SQS CloudWatch Metrics | AWS Documentation | https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-cloudwatch-metrics.html | VERIFIED | `ApproximateAgeOfOldestMessage`, `ApproximateNumberOfMessagesVisible` |
| 9 | Kafka Monitoring | Apache Kafka Documentation | https://kafka.apache.org/documentation/#monitoring | VERIFIED | Consumer Lag, latency metrics |

---

## Secondary Sources (Referenced)

| # | Title | Publisher | URL | Status | Relevance |
|---|-------|-----------|-----|--------|-----------|
| 10 | Netflix/Cloudflare/Stripe Engineering Articles | Various | N/A | NOT VERIFIED | Illustrative examples only; not primary source |

---

## Source Verification Summary

- **Total Primary Sources**: 9
- **Verified**: 9 (100%)
- **Not Verified**: 0
- **Secondary Sources**: 1 (explicitly noted as illustrative)

---

## Source Mapping to Claims

| Claim | Source # | Status |
|-------|----------|--------|
| HTTP 429 standardization | 1, 2 | VERIFIED |
| Token Bucket definition | 3 | VERIFIED |
| Little's Law formula | 5 | VERIFIED |
| Deterministic queue buildup | N/A (text explanation) | VERIFIED |
| Queue memory exhaustion | 7 | VERIFIED |
| Exponential backoff + jitter | 6 | VERIFIED |
| IP rate limiting limitation | 4 | VERIFIED |
| Fair queuing concepts | N/A (academic theory) | INTERPRETATION |
| Queue age metric | 8, 9 | VERIFIED |
| Autoscaling downstream bottleneck | 7 | VERIFIED |
| Retry storm mitigation | 6, 7 | VERIFIED |

---

## Revision Notes

- **Created**: New source index to replace inline source references in research plan
- **Verified**: All cited sources are reachable and support their claimed relevance
- **Removed**: RFC 8305 (not relevant), RFC 5321 (SMTP not relevant to rate limiting)
- **Added**: AWS SQS CloudWatch metrics (critical for queue age claim)
- **Added**: Kafka monitoring documentation (critical for queue age claim)