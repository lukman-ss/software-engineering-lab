# Open Questions

## Unanswered Questions

1. **Exact implementation details of per-tenant fair queueing algorithms:** While the topic mentions "Fair queueing" as a solution for multi-tenant isolation, specific algorithms (e.g., Weighted Fair Queuing variants) and their implementation trade-offs need deeper investigation.

2. **Dynamic cost-based rate limiting:** The topic mentions rate limits should consider "cost of operation" but how to dynamically compute operation costs in real-time is not well documented in public sources.

3. **Queue aging vs. queue length:** While queue age is mentioned as a better metric than queue length, the exact thresholds and alerting strategies for different system types (API vs batch processing) need empirical research.

4. **Autoscaling interaction with backpressure:** The topic mentions "autoscaling can make downstream collapse faster" but specific patterns for coordinating autoscaling with rate limiting/batching need more investigation.

5. **Dead letter queue (DLQ) strategies:** Retry with bounded limits and DLQ handling is mentioned, but optimal DLQ processing strategies for different failure types (transient vs permanent) are not fully covered.

## Weak Evidence

1. **Topic spec claims about "1 public IP for 200 users":** While true that NAT shares an IP, the actual implementation details of how enterprise networks route API traffic through multiple egress IPs or how to handle this for rate limiting need verification from enterprise architecture sources.

2. **"500 request/second downing an endpoint" scenario:** While plausible, specific failure cascade timing (worker saturation → DB connection exhaustion → CPU overload → health check failure → restart) needs empirical validation.

3. **Burst capacity recommendations:** The topic spec suggests specific token bucket values (capacity=100, refill=10/s) but doesn't cite sources for these specific values.

## Claims Needing Deeper Research

1. **Per-tenant concurrency limits of "5 concurrent jobs":** The specific recommendation of equal per-tenant limits (5) may be oversimplified. Better approaches like weighted fair queuing or adaptive limits based on tenant tier need investigation.

2. **Retry storm threshold:** The claim that retrying at the same rate as the original traffic creates a "retry storm" needs quantitative analysis of specific retry multiplier effects.

3. **Health check failure patterns during overload:** The cascade from "endpoint slow" to "health check fails" to "instance restart" is plausible but the exact mechanisms and prevention strategies need deeper investigation.

## Possible Next Research Directions

1. **Investigate actual implementations of rate limiters in production:**
   - Envoy Proxy rate limit filter implementation
   - Kong API gateway rate limiting plugin
   - NGINX Plus rate limiting features

2. **Research specific backpressure implementations:**
   - Reactive Streams backpressure protocol (Project Reactor, RxJava)
   - TCP congestion control mechanisms
   - HTTP/2 flow control

3. **Analyze real-world incident postmortems involving rate limiting/backpressure failures:**
   - AWS service throttling incidents
   - Google SRE incident reports on overload
   - Industry case studies on cascading failures

4. **Benchmark different rate limiting algorithms under realistic load patterns**

5. **Investigate multi-region rate limiting and global consistency approaches**