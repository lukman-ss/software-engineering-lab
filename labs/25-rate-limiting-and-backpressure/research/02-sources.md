# Research Sources

## Source 1

Title: Token bucket - Wikipedia
Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Token_bucket
Published: 14 September 2026 (last revision)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable reference work
Relevance: Primary source on Token Bucket algorithm - one of the key rate limiting algorithms
Evidence: "A conforming flow can thus contain traffic with an average rate up to the rate at which tokens are added to the bucket, and have a burstiness determined by the depth of the bucket."

## Source 2

Title: Exponential backoff - Wikipedia
Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Exponential_backoff
Published: 21 August 2026 (last revision)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable reference work
Relevance: Primary source on Exponential backoff algorithm with jitter, including binary exponential backoff and collision avoidance

## Source 3

Title: Little's law - Wikipedia
Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Little%27s_law
Published: 20 August 2026 (last revision)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable reference work
Relevance: Primary source on Little's Law for queue analysis - L = λW

## Source 4

Title: Rate limiting - Wikipedia
Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Rate_limiting
Published: 2 September 2026 (last revision)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable reference work
Relevance: Primary source on rate limiting concepts, algorithms (Token bucket, Leaky bucket, Fixed window, Sliding window)

## Source 5

Title: Module ngx_http_limit_req_module
Publisher: NGINX (F5)
URL: https://nginx.org/en/docs/http/ngx_http_limit_req_module.html
Published: N/A (Documentation)
Accessed: 2026-09-26
Source Tier: Tier 1 - Official documentation
Relevance: NGINX implements rate limiting using the leaky bucket method via limit_req_module
Evidence: "The limitation is done using the 'leaky bucket' method." Default status code: 503

## Source 6

Title: RFC 6585: Additional HTTP Status Codes
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc6585#section-4
Published: April 2012
Accessed: 2026-09-26
Source Tier: Tier 1 - Standards document
Relevance: Defines HTTP 429 Too Many Requests status code
Evidence: "The 429 status code indicates that the user has sent too many requests in a given amount of time ('rate limiting')." MUST NOT be stored by a cache.

## Source 7

Title: Scaling your API with rate limiters
Publisher: Stripe Engineering Blog (Paul Tarjan)
URL: https://stripe.com/blog/rate-limiters
Published: March 30, 2017
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable technical publication
Relevance: Four types of limiters in production: Request rate limiter, Concurrent requests limiter, Fleet usage load shedder, Worker utilization load shedder

## Source 8

Title: Exponential Backoff And Jitter
Publisher: AWS Architecture Blog (Marc Brooker)
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: March 4, 2015 (Updated May 2023)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable technical publication
Relevance: Three jitter algorithms: Full Jitter, Equal Jitter, Decorrelated Jitter. Full Jitter reduces call count by more than half with 100 contending clients.

## Source 9

Title: Retry behavior
Publisher: AWS SDK Documentation
URL: https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
Published: May 2023 (updated)
Accessed: 2026-09-26
Source Tier: Tier 1 - Official documentation
Relevance: Full Jitter implementation: delay = random(0,1) × min(20000 ms, base_delay × 2^retry). Retry quota uses token bucket. Transient: 50ms base, Throttling: 1000ms base. Max cap: 20 seconds.

## Source 10

Title: Handling Overload
Publisher: Google SRE Workbook (Alejandro Forero Cuervo)
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Published: 2016 (O'Reilly)
Accessed: 2026-09-26
Source Tier: Tier 1 - Authoritative systems engineering resource
Relevance: Rate limiting vs backpressure distinction, client-side throttling, per-customer quotas, criticality levels, utilization signals

## Source 11

Title: Redis rate limiter documentation
Publisher: Redis Documentation
URL: https://redis.io/docs/latest/develop/use-cases/rate-limiter/
Published: N/A (Documentation)
Accessed: 2026-09-26
Source Tier: Tier 1 - Official documentation
Relevance: Redis supports fixed window (INCR/EXPIRE), sliding window (sorted sets), and token bucket (Lua scripts) for rate limiting

## Source 12

Title: Leaky bucket - Wikipedia
Publisher: Wikimedia Foundation
URL: https://en.wikipedia.org/wiki/Leaky_bucket
Published: 2026 (last revision)
Accessed: 2026-09-26
Source Tier: Tier 2 - Reputable reference work
Relevance: Two variants: leaky bucket as a meter (equivalent to token bucket) and leaky bucket as a queue. NGINX uses leaky bucket as a meter.

## Source 13

Title: Consumer Prefetch & Queue Flow Control
Publisher: RabbitMQ (Broadcom)
URL: https://www.rabbitmq.com/docs/consumer-prefetch
Published: N/A (Documentation)
Accessed: 2026-09-26
Source Tier: Tier 1 - Official documentation
Relevance: RabbitMQ consumer prefetch for queue-based backpressure/flow control
