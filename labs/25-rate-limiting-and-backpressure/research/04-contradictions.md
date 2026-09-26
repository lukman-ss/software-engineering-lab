# Contradictions Analysis

## No Material Contradictions Discovered

After reviewing all sources, no significant factual contradictions were found between authoritative sources. The following minor differences in emphasis were observed but do not represent contradictions:

### Areas of Consistent Agreement:

1. **HTTP 429 Status Code**: RFC 6585 and Wikipedia both confirm 429 is the standard rate limiting response with optional Retry-After header.

2. **Token Bucket Algorithm**: All sources (Wikipedia, IEEE paper, ScyllaDB blog, Medium article) describe the same core algorithm: tokens added at fixed rate, bucket capacity limits burst, packet requires n tokens to pass.

3. **Backpressure Definition**: Reactive Streams specification and general systems literature agree: backpressure prevents downstream from being overwhelmed by upstream production rate.

4. **Little's Law**: Universally accepted as L = λW with the same conditions (ergodic, stationary system).

5. **Exponential Backoff + Jitter**: AWS blog and general practice agree jitter is essential to prevent retry storms; multiple jitter variants exist (Full, Equal, Decorrelated).

### Minor Differences in Emphasis:

1. **Leaky Bucket Confusion**: Wikipedia Leaky Bucket article notes there are TWO versions (as meter and as queue) causing confusion in literature. Token bucket article states they are "fundamentally the same" when implemented correctly with same parameters. This is a terminology issue, not a factual contradiction.

2. **Distributed Rate Limiting**: Wikipedia mentions Redis/Aerospike for distributed rate limiting; Medium article discusses sliding window log alternative. These are complementary approaches, not contradictory.

3. **Rate Limiting vs Throttling**: RFC 6585 says servers "not required to use 429; may drop connections during attacks." Wikipedia mentions "should be used along with throttling pattern." Both agree on the mechanism; difference is in terminology and when to apply each.

4. **Retry Behavior**: AWS blog focuses on client-side retry logic. RFC 6585 mentions Retry-After header for server-side guidance. Both are complementary layers (client and server).

### Assessment:

All sources are consistent on core technical facts. Differences are in:
- Level of abstraction (standard vs implementation vs theory)
- Specific use case focus (network vs application vs database)
- Terminology preferences

No source contradicts another on fundamental mechanisms or mathematical relationships.