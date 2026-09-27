# Contradictions and Disagreements Audit

## Contradiction 1: HTTP Status Code Selection (429 vs 503)

Statement A:
"The 429 status code indicates that the user has sent too many requests in a given amount of time ('rate limiting')."

Location:
`research/02-sources.md`: Source 6 (RFC 6585)

Statement B:
NGINX `limit_req_module` defaults to HTTP 503 Service Unavailable when rate limits are exceeded.

Location:
`research/02-sources.md`: Source 5 (NGINX docs)

Type:
SOURCE_CONFLICT

Impact:
LOW

Assessment:
Not a material contradiction. HTTP 429 is the client-quota rate limit standard (client behavior), whereas HTTP 503 is server overload/capacity shedding (server state). NGINX allows setting `limit_req_status 429;` if desired. The research report accurately captures this distinction in Finding 4.

---

## Contradiction 2: Token Bucket vs Leaky Bucket Operational Model

Statement A:
Token bucket allows instantaneous bursts up to bucket capacity $b$, refilling tokens at rate $r$.

Location:
`research/03-evidence.md`: Evidence 1 (Token bucket Wikipedia)

Statement B:
Leaky bucket (as a queue) smooths traffic to a strict constant output rate without burst output.

Location:
`research/03-evidence.md`: Evidence 2 (Leaky bucket Wikipedia)

Type:
SOURCE_CONFLICT

Impact:
LOW

Assessment:
Resolved. As documented in Source 12 and Research Report Finding 1, "Leaky bucket as a meter" (conformance checking) is mathematically equivalent to token bucket, whereas "Leaky bucket as a queue" (traffic shaping) eliminates bursts. The research report correctly clarifies this nuance.

---

## Contradiction 3: Retrying vs Dropping Failed Requests During Overload

Statement A:
AWS SDK standard retry mode automatically retries failed requests using exponential backoff with full jitter up to 3 or 4 attempts.

Location:
`research/03-evidence.md`: Evidence 7 & 8 (AWS SDK Retry Behavior)

Statement B:
Google SRE recommends: "If a large subset of backend tasks in the datacenter are overloaded, requests should not be retried and errors should bubble up all the way to the caller."

Location:
`research/04-contradictions.md`: Area 8; Google SRE Chapter 21

Type:
SOURCE_CONFLICT

Impact:
MEDIUM

Assessment:
Architectural level difference. AWS SDK manages client-to-cloud API retries with a client token-bucket retry quota (500 tokens). Google SRE manages deep intra-cluster RPC call chains where per-layer retries cause exponential amplification. The research report correctly notes that retries must be budgeted and stopped at upper layers using "overloaded; don't retry" signals.

---

## Contradiction 4: Source Publication / Access Date Anomalies

Statement A:
"Published: 14 September 2026 (last revision)" / "Accessed: 2026-09-26"

Location:
`research/02-sources.md`: Sources 1-4, 12

Statement B:
Today's date in environment context: Sun Sep 27 2026.

Location:
Environment prompt context.

Type:
INTERNAL

Impact:
LOW

Assessment:
The research agent recorded current system dates (September 2026) for Wikipedia revisions. All URLs were fetched live during audit and verified to exist with identical contents. No fabricated citations exist.

---

## Summary Verdict on Contradictions

No material unhandled contradictions found. All nuances between sources (HTTP 429 vs 503, Token vs Leaky Bucket variants, and Client retry vs Server shedding) are accurately analyzed and explained in `research/04-contradictions.md` and `research/05-report.md`.
