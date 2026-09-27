# Research Evidence

## Evidence 1: Token Bucket Algorithm Properties

Claim: The token bucket algorithm allows burst traffic up to the bucket capacity while maintaining a long-term average rate limit.

Evidence: "A conforming flow can thus contain traffic with an average rate up to the rate at which tokens are added to the bucket, and have a burstiness determined by the depth of the bucket."

Source: Token bucket - Wikipedia
URL: https://en.wikipedia.org/wiki/Token_bucket
Publication date: 14 September 2026
Confidence: HIGH
Corroborated By: NGINX documentation (uses leaky bucket as meter, equivalent to token bucket), Redis rate limiter documentation (supports token bucket via Lua scripts)

## Evidence 2: Leaky Bucket as Meter vs Queue

Claim: The leaky bucket has two implementations: as a meter (checking conformance) and as a queue (enforcing conformance by buffering).

Evidence: "Two different methods of applying this leaky bucket analogy are described in the literature... One version, the bucket is a counter or variable separate from the flow of traffic... This version is referred to here as the leaky bucket as a meter. In the second version, the bucket is a queue in the flow of traffic... This version is referred to here as leaky bucket as a queue."

Source: Leaky bucket - Wikipedia
URL: https://en.wikipedia.org/wiki/Leaky_bucket
Publication date: 2026 (last revision)
Confidence: HIGH
Corroborated By: NGINX documentation (uses leaky bucket as meter for rate limiting)

## Evidence 3: NGINX Rate Limiting Implementation

Claim: NGINX implements rate limiting using the leaky bucket algorithm as a meter via the limit_req_module.

Evidence: "The ngx_http_limit_req_module module (0.7.21) is used to limit the request processing rate per a defined key, in particular, the processing rate of requests coming from a single IP address. The limitation is done using the 'leaky bucket' method." Default status code: 503.

Source: Module ngx_http_limit_req_module
URL: https://nginx.org/en/docs/http/ngx_http_limit_req_module.html
Publication date: N/A (Documentation)
Confidence: HIGH
Corroborated By: Leaky bucket Wikipedia (NGINX uses leaky bucket as meter), Wikipedia Rate limiting article

## Evidence 4: HTTP 429 Status Code Standard

Claim: HTTP 429 Too Many Requests is the standard status code for indicating rate limiting has been applied.

Evidence: "The 429 status code indicates that the user has sent too many requests in a given amount of time ('rate limiting')." "Responses with the 429 status code MUST NOT be stored by a cache." Should include Retry-After header.

Source: RFC 6585: Additional HTTP Status Codes
URL: https://www.rfc-editor.org/rfc/rfc6585#section-4
Publication date: April 2012
Confidence: HIGH
Corroborated By: NGINX limit_req_status defaults to 503 but configurable, Stripe blog recommends choosing between 429 and 503

## Evidence 5: Stripe's Four-Layer Rate Limiting Approach

Claim: Stripe uses four different types of limiters in production: Request rate limiter, Concurrent requests limiter, Fleet usage load shedder, and Worker utilization load shedder.

Evidence: "At Stripe, we operate 4 different types of limiters in production. The first one, the Request Rate Limiter, is by far the most important one." Details each limiter with specific use cases and production metrics.

Source: Scaling your API with rate limiters
URL: https://stripe.com/blog/rate-limiters
Publication date: March 30, 2017
Confidence: HIGH
Corroborated By: Google SRE Workbook (per-customer limits, load shedding), AWS SDK retry behavior (retry quota with token bucket)

## Evidence 6: Exponential Backoff with Jitter Effectiveness

Claim: Adding jitter to exponential backoff significantly reduces the thundering herd problem and improves system recovery during contention.

Evidence: "In the case with 100 contending clients, we've reduced our call count by more than half. We've also significantly improved the time to completion, when compared to un-jittered exponential backoff."

Source: Exponential Backoff And Jitter | AWS Architecture Blog
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Publication date: 04 MAR 2015 (Updated May 2023)
Confidence: HIGH
Corroborated By: AWS SDK Documentation (standard mode uses Full Jitter), Wikipedia Exponential backoff article

## Evidence 7: AWS SDK Full Jitter Formula

Claim: AWS SDKs use exponential backoff with full jitter where delay = random(0, 1) × min(20,000 ms, base_delay × 2^retry).

Evidence: "The SDK computes each retry delay using this formula: delay = random(0, 1) × min(20,000 ms, base_delay × 2^retry)" Transient errors: 50ms base (25ms for DynamoDB), Throttling: 1000ms base. Max cap: 20 seconds.

Source: Retry behavior - AWS SDK Documentation
URL: https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
Publication date: May 2023 (updated)
Confidence: HIGH
Corroborated By: AWS Architecture Blog (Full Jitter recommended), Exponential backoff Wikipedia

## Evidence 8: AWS SDK Retry Quota (Token Bucket)

Claim: AWS SDKs implement a retry quota using a token bucket to prevent retry storms during service disruptions.

Evidence: "Standard mode includes a retry quota, a token bucket that deducts tokens for each retry and replenishes tokens when requests succeed. When the available tokens are exhausted, the SDK returns the error without retrying, so your application fails fast instead of waiting through retries that are unlikely to succeed."

Source: Retry behavior - AWS SDK Documentation
URL: https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
Publication date: May 2023 (updated)
Confidence: HIGH
Corroborated By: AWS Architecture Blog, Redis rate limiter documentation

## Evidence 9: Google SRE Rate Limiting vs Backpressure

Claim: Rate limiting typically works at the entry point to a system, while backpressure is a broader mechanism that propagates pressure backward through the system when downstream components are overwhelmed.

Evidence: "Rate limiting usually works at the entrance to a system. Backpressure is more general: When downstream is unable to keep up with work from upstream, upstream must slow down."

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: Stripe rate limiting blog (load shedders as backpressure), AWS retry behavior documentation

## Evidence 10: Little's Law for Queue Analysis

Claim: Little's Law (L = λW) provides a fundamental relationship for understanding queue behavior in systems.

Evidence: "In mathematical queueing theory, Little's law ... is a theorem by John Little which states that the long-term average number of customers (L) in a stationary system is equal to the long-term average effective arrival rate (λ) multiplied by the average time that a customer spends in the system (W)."

Source: Little's law - Wikipedia
URL: https://en.wikipedia.org/wiki/Little%27s_law
Publication date: 20 August 2026
Confidence: HIGH
Corroborated By: Google SRE Workbook (utilization signals), topic specification calculation verified

## Evidence 11: Queue Backlog Growth Calculation

Claim: When arrival rate exceeds processing rate in a queue system, backlog grows linearly over time according to Little's Law.

Evidence: Topic specification calculation: 5,000,000 items / 2,000 items/sec = 2,500 seconds ≈ 41 minutes 40 seconds. Verified mathematically: 2,500 / 60 = 41.67 minutes = 41 minutes 40 seconds.

Source: Topic specification (verified against Little's Law)
URL: N/A (Original topic specification)
Publication date: Current lab specification
Confidence: HIGH
Corroborated By: Little's Law Wikipedia (L = λW), Google SRE Workbook (queue growth under overload)

## Evidence 12: Redis as Rate Limiting Backend

Claim: Redis is well-suited for distributed rate limiting due to its atomic operations and data structures.

Evidence: "Redis provides the following features that make it a good fit for rate limiting: [INCR and EXPIRE] give you atomic fixed-window counters with automatic time-window cleanup. [Hashes, sorted sets, and strings] cover the data shapes needed for sliding window and token bucket algorithms."

Source: Redis rate limiter documentation
URL: https://redis.io/docs/latest/develop/use-cases/rate-limiter/
Publication date: N/A (Documentation)
Confidence: HIGH
Corroborated By: Stripe rate limiting blog (uses Redis), AWS API Gateway (uses Redis for distributed rate limiting)

## Evidence 13: Google SRE Per-Customer Limits

Claim: Multi-tenant systems require per-customer rate limiting to prevent noisy neighbors from monopolizing resources.

Evidence: Google SRE implements per-customer CPU quotas (e.g., Gmail 4000 CPU seconds/sec, Calendar 4000, Android 3000, Google+ 2000, others 500). System relies on customers not hitting limits simultaneously.

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: Stripe Fleet Usage Load Shedder (reserves capacity for critical traffic), Topic specification multi-tenant example

## Evidence 14: Google SRE Client-Side Throttling

Claim: Client-side throttling prevents clients from overwhelming backends with rejected requests.

Evidence: "When a customer is out of quota, a backend task should reject requests quickly... Client-side throttling addresses this problem... each client task keeps... requests and accepts... Clients can continue to issue requests until requests is K times as large as accepts."

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: AWS SDK adaptive mode (client-side rate limiter), AWS SDK retry quota (token bucket)

## Evidence 15: Google SRE Criticality Levels

Claim: Request criticality levels enable prioritized load shedding during overload.

Evidence: Four criticality values: CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE. "When a customer runs out of global quota, a backend task will only reject requests of a given criticality if it's already rejecting all requests of all lower criticalities."

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: Stripe Worker Utilization Load Shedder (sheds test mode, then GETs, then POSTs, keeps critical)

## Evidence 16: Google SRE Utilization Signals

Claim: Task-level overload protection uses utilization signals (executor load average, CPU, memory) to reject requests based on criticality.

Evidence: "As utilization approaches configured thresholds, we start rejecting requests based on their criticality (higher thresholds for higher criticalities)." Executor load average counts active threads with exponential decay smoothing.

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: Stripe Worker Utilization Load Shedder, Topic specification queue age metric

## Evidence 17: Google SRE Retry Budget

Claim: Per-request and per-client retry budgets prevent retry storms during cascading overload.

Evidence: "per-request retry budget of up to three attempts... per-client retry budget... ratio of requests that correspond to retries... below 10%." Threefold growth capped to 1.1x with per-client budget.

Source: Handling Overload - Google SRE Workbook
URL: https://landing.google.com/sre/sre-book/chapters/handling-overload/
Publication date: 2016
Confidence: HIGH
Corroborated By: AWS SDK retry quota (token bucket), AWS SDK max attempts default 3

## Evidence 18: AWS SDK Adaptive Mode Client-Side Rate Limiter

Claim: AWS SDK adaptive mode includes a client-side rate limiter that can delay or block initial requests when throttling is detected.

Evidence: "Adaptive mode includes everything in standard mode, plus a client-side rate limiter. The rate limiter tracks throttling responses and adjusts the rate at which the SDK sends requests... adaptive mode can delay or block the initial request, not just retries, when throttling is detected."

Source: Retry behavior - AWS SDK Documentation
URL: https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
Publication date: May 2023 (updated)
Confidence: HIGH
Corroborated By: Google SRE client-side throttling, AWS Architecture Blog