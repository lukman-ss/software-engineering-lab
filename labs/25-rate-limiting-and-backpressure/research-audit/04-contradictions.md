# Contradiction Audit

## Contradiction 1

Statement A:
"The limitation is done using the 'leaky bucket' method." (NGINX docs for `ngx_http_limit_req_module`)

Location:
`research/02-sources.md:52`, `research/03-evidence.md:59`

Statement B:
Token Bucket algorithm adds tokens at fixed rate and deducts on request, allowing bursts up to burst capacity.

Location:
`research/03-evidence.md:5-7`, `research/04-contradictions.md:9-26`

Type:
SOURCE_CONFLICT / CONCEPTUAL_NUANCE

Impact:
Engineers often confuse Leaky Bucket and Token Bucket behavior or treat them as identical. Leaky bucket as a meter checks average rates with an accumulator, whereas leaky bucket as a queue enforces smooth, unbursty egress. Token bucket permits bursty transmission at wire speed up to capacity.

Assessment:
Properly analyzed in `research/04-contradictions.md` Area 1. The research explicitly clarifies the dual definition of leaky bucket (meter vs queue) and notes why NGINX's meter variant behaves similarly to a bounded burst limiter. No unaddressed internal contradiction.

---

## Contradiction 2

Statement A:
RFC 6585 specifies HTTP 429 Too Many Requests as the standard status code indicating rate limiting has occurred.

Location:
`research/03-evidence.md:31-39`, `research/04-contradictions.md:47-50`

Statement B:
Stripe recommends evaluating whether HTTP 429 or HTTP 503 is more accurate depending on whether the limitation is user quota or server overload.

Location:
`research/04-contradictions.md:51-55`

Type:
SOURCE_CONFLICT / OPERATIONAL_NUANCE

Impact:
Client error (4xx) vs Server error (5xx) impacts automated client retry semantics and alerting.

Assessment:
Accurately categorized in `04-contradictions.md`. HTTP 429 implies client quota exhaustion (client should back off / fix request rate); HTTP 503 implies server capacity overload (load shedding / backpressure). Both are valid in their respective contexts.

---

## Contradiction 3

Statement A:
AWS Architecture Blog analyzes three jitter variants: Full Jitter, Equal Jitter, and Decorrelated Jitter.

Location:
`research/04-contradictions.md:31-35`

Statement B:
AWS SDK standard retry behavior enforces a single standard formula: Full Jitter (`random(0, 1) * min(20000ms, base_delay * 2^retry)`).

Location:
`research/03-evidence.md:135-138`, `research/04-contradictions.md:36-39`

Type:
SOURCE_CONFLICT

Impact:
Low. Explains experimental exploration vs production standard.

Assessment:
Consistent and reconciled. The research correctly explains why SDK standardization adopted Full Jitter.

---

## Contradiction 4

Statement A:
Source 11 Title is labeled "Cloudflare's Rate Limiting Documentation (Redis rate limiter)".

Location:
`research/02-sources.md:105`

Statement B:
Publisher is "Redis Documentation" and URL is `https://redis.io/docs/latest/develop/use-cases/rate-limiter/`.

Location:
`research/02-sources.md:106-107`

Type:
INTERNAL

Impact:
Metadata confusion. Conflates Cloudflare's perimeter rate limiting with Redis's in-memory data store documentation.

Assessment:
Minor internal metadata contradiction in `02-sources.md`. Content referenced in evidence and report pertains strictly to Redis primitives (`INCR`, `EXPIRE`, sorted sets).
