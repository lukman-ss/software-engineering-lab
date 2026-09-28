# Content Brief

Topic: Rate Limiting & Backpressure in Distributed Systems
Target Reader: Backend engineers, SREs, distributed systems architects
Problem: Systems without rate limiting and backpressure risk cascading failures, memory exhaustion, and uncontrolled retry storms when downstream components become overwhelmed.
Core Mental Model: Rate limiting controls inbound traffic at system boundaries (token/leaky bucket), backpressure propagates pressure signals backward when downstream is overloaded (bounded queue rejection), retry jitter prevents thundering herd (exponential backoff with Full Jitter).
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Token Bucket (burst accommodation), Leaky Bucket (traffic smoothing), HTTP 429 RFC 6585, Bounded Queue Backpressure, Little's Law (L = λW), Exponential Backoff with Full Jitter, Multi-tenant Rate Limiting
Verified Behaviors: Token bucket allows burst up to capacity then enforces refill rate; Leaky bucket smooths output to constant drain rate; Bounded queue rejects excess jobs immediately with fast backpressure; HTTP middleware returns 429 + Retry-After header; Full Jitter distributes retry intervals uniformly in [0, min(cap, base·2^attempt)] to prevent synchronization
Available Case Studies: Stripe (4-layer rate limiting), Google SRE (per-customer quotas, 4-level criticality, retry budgets), AWS SDK (Full Jitter formula, token bucket retry quota)
Warnings: In-memory implementation lacks distributed synchronization (Redis/Memcached); Default base delays (50ms transient, 1000ms throttling) are AWS SDK-specific; Queue age more useful than raw depth; AWS constants must be calibrated to downstream SLA; No per-job timeout in TrySubmit (optional enhancement); Zero refillRate panic in RetryAfterSeconds (documented precondition)
