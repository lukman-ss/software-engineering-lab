# Source Audit

## Source 1

Claimed Title: Token bucket - Wikipedia
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Token_bucket

Reachable:
YES

Source Type:
COMMUNITY (Tier 2 reference work)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Wikipedia is a secondary/tertiary reference work. While useful for conceptual overview, it lacks primary academic/industry authority.

Assessment:
PASS

---

## Source 2

Claimed Title: Exponential backoff - Wikipedia
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Exponential_backoff

Reachable:
YES

Source Type:
COMMUNITY (Tier 2 reference work)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None for basic mathematical definitions of exponential backoff.

Assessment:
PASS

---

## Source 3

Claimed Title: Rate limiting - Wikipedia
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Rate_limiting

Reachable:
YES

Source Type:
COMMUNITY (Tier 2 reference work)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Secondary reference; needs primary industry backing for production implementation rules.

Assessment:
PASS

---

## Source 4

Claimed Title: Little's law - Wikipedia
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Little%27s_law

Reachable:
YES

Source Type:
COMMUNITY (Tier 2 reference work)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Mathematical theorem is well-represented.

Assessment:
PASS

---

## Source 5

Claimed Title: Module ngx_http_limit_req_module
Claimed Publisher: NGINX (F5)
URL: https://nginx.org/en/docs/http/ngx_http_limit_req_module.html

Reachable:
YES

Source Type:
PRIMARY (Official Documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official primary documentation for NGINX rate limiting implementation.

Assessment:
PASS

---

## Source 6

Claimed Title: HTTP/1.1 429 Too Many Requests - RFC 6585
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc6585#section-4

Reachable:
YES

Source Type:
PRIMARY (Standards Document)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative standard for HTTP 429 status code.

Assessment:
PASS

---

## Source 7

Claimed Title: Scaling your API with rate limiters
Claimed Publisher: Stripe Engineering Blog
URL: https://stripe.com/blog/rate-limiters

Reachable:
YES

Source Type:
PRIMARY (Official Engineering Blog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Highly authoritative industry post detailing 4 limiter types in high-scale production.

Assessment:
PASS

---

## Source 8

Claimed Title: Exponential Backoff And Jitter
Claimed Publisher: AWS Architecture Blog
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

Reachable:
YES

Source Type:
PRIMARY (Official Engineering Blog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical article introducing Full Jitter and Equal Jitter algorithms.

Assessment:
PASS

---

## Source 9

Claimed Title: Retry behavior
Claimed Publisher: AWS SDK Documentation
URL: https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html

Reachable:
YES

Source Type:
PRIMARY (Official Documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Primary source for AWS SDK retry quota token bucket and backoff behavior.

Assessment:
PASS

---

## Source 10

Claimed Title: Handling Overload
Claimed Publisher: Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/

Reachable:
YES

Source Type:
PRIMARY (Authoritative Engineering Book)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- URL points to Google SRE Book (Chapter 22), while text mentions SRE Workbook. Content on load shedding and client-side throttling is directly present in Chapter 22 of the SRE Book.

Assessment:
PASS

---

## Source 11

Claimed Title: Cloudflare's Rate Limiting Documentation (Redis rate limiter)
Claimed Publisher: Redis Documentation
URL: https://redis.io/docs/latest/develop/use-cases/rate-limiter/

Reachable:
YES

Source Type:
PRIMARY (Official Documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Title mismatch in source listing: Title claims "Cloudflare's Rate Limiting Documentation (Redis rate limiter)", but Publisher and URL belong to Redis Official Documentation. Title text conflates Cloudflare with Redis.

Assessment:
WARNING

---

## Source 12

Claimed Title: Leaky bucket - Wikipedia
Claimed Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Leaky_bucket

Reachable:
YES

Source Type:
COMMUNITY (Tier 2 reference work)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Secondary reference, but provides crucial distinction between "leaky bucket as meter" and "leaky bucket as queue".

Assessment:
PASS

---

## Source 13

Claimed Title: RabbitMQ Tutorials
Claimed Publisher: RabbitMQ (Broadcom)
URL: https://www.rabbitmq.com/tutorials

Reachable:
YES

Source Type:
PRIMARY (Official Documentation)

Relevant:
PARTIAL

Supports Claimed Topic:
PARTIAL

Problems:
- URL points to generic tutorial landing page without linking to specific backpressure/prefetch documentation (e.g. `channel.basicQos`).

Assessment:
WARNING
