# Research Plan

## Research Topic
Rate Limiting & Backpressure - Sistem Sehat Harus Bisa Bilang "Cukup"

## Objective
Investigate rate limiting and backpressure mechanisms in distributed systems, focusing on practical engineering approaches, Token Bucket and Leaky Bucket algorithms, Exponential Backoff with Jitter, queue management, and their role in system scalability and fault isolation.

## Research Questions
1. What are the primary rate limiting algorithms (Token Bucket, Leaky Bucket) and their practical implementations?
2. How does backpressure differ from rate limiting, and how is it implemented in queue-based systems?
3. What are the best practices for implementing exponential backoff with jitter?
4. How does Little's Law provide the mathematical foundation for queue analysis?
5. What are the multi-layer approaches to rate limiting used in production systems?

## Search Strategy
- Search for official documentation on rate limiting algorithms
- Look for Google SRE Workbook on Handling Overload
- Find industry best practices from cloud providers (AWS, Stripe)
- Search for Redis rate limiter documentation
- Verify HTTP 429 status code specification (RFC 6585)

## Expected Primary Sources
- RFC 6585: HTTP 429 Too Many Requests
- Google SRE Workbook: Handling Overload chapter
- AWS Architecture Blog: Exponential Backoff and Jitter
- Stripe Engineering Blog: Scaling your API with rate limiters
- NGINX documentation: ngx_http_limit_req_module
- Redis documentation: Rate limiter pattern
- AWS SDK Documentation: Retry behavior

## Risks / Unknowns
- Calculation verification: 5,000,000 items / 2,000 items/sec = 2,500 sec ≈ 41m 40s needs verification
- JITT implementation variations (Full Jitter vs Equal Jitter vs Decorrelated Jitter)
- NGINX uses leaky bucket as meter, not as queue