# Open Questions

## Unanswered Questions

1. **Distributed Token Bucket Consistency**: How to maintain exact token bucket consistency across distributed instances without centralized Redis? What are the trade-offs of eventual consistency vs strong consistency in distributed rate limiting?

2. **Rate Limiting at Different Layers**: What are the specific implementation patterns and trade-offs for rate limiting at:
   - API Gateway level (AWS API Gateway, Kong, Envoy)
   - Application middleware level
   - Database connection pool level
   - OS/kernel level (TC, cgroups)

3. **Circuit Breaker Integration**: How do circuit breakers interact with rate limiting and backpressure? Best practices for cascading failure prevention when combining these patterns?

4. **gRPC Flow Control**: How does gRPC's HTTP/2 flow control (window-based) compare to application-level rate limiting? When to use each?

5. **Kubernetes Backpressure**: How do Kubernetes HPA, KEDA, and queue-proxy interact with application-level backpressure? What metrics should trigger scaling vs backpressure?

6. **Observability Standards**: What are the standard Prometheus metrics for rate limiting and backpressure? (e.g., rate_limit_exceeded_total, queue_depth, queue_oldest_age, retry_rate)

## Weak Evidence Areas

1. **Multi-tenant Fair Queueing Algorithms**: Limited public documentation on specific fair queueing implementations (Deficit Round Robin, Weighted Fair Queueing) in application-level systems.

2. **Database-Level Backpressure**: ScyllaDB's token bucket for IO is one example; how do PostgreSQL, MySQL, MongoDB handle backpressure from application connections?

3. **Webhook Provider Rate Limits**: Industry standards for webhook rate limiting (Stripe, GitHub, Slack) - are they documented publicly?

4. **Real-world Failure Case Studies**: Limited public post-mortems of production outages caused by rate limiting/backpressure failures.

5. **Machine Learning / AI Workloads**: How do these patterns apply to GPU inference workloads with highly variable processing times?

## Claims Needing Deeper Research

1. **Token Bucket vs Sliding Window Log**: Medium article claims sliding window log is alternative for distributed rate limiting. Need quantitative comparison of accuracy, memory, latency.

2. **Hierarchical Token Bucket (HTB)**: Linux HTB used for traffic control. Can/should this be applied at application layer for multi-tenant systems?

3. **Rate Limiting with Cost-based Tokens**: ScyllaDB uses normalized sum of weight + length. How to generalize cost-based tokens for arbitrary API operations?

4. **Dynamic Rate Limit Adjustment**: How to automatically adjust rate limits based on system health metrics (CPU, latency, error rate) without oscillation?

5. **Priority Queue Implementation**: Best practices for priority queues in multi-tenant systems where enterprise tenants get higher priority without starving free tier.

## Possible Next Research Directions

1. **Empirical Benchmarking**: Compare rate limiting algorithm implementations under realistic workloads (burst, sustained, adversarial).

2. **Failure Mode Analysis**: Systematic analysis of failure modes when rate limiting/backpressure is misconfigured or absent.

3. **Cloud Provider Implementations**: Deep dive into AWS API Gateway, Google Cloud Load Balancing, Azure API Management rate limiting implementations.

4. **Service Mesh Integration**: How Istio/Linkerd rate limiting and circuit breaking integrates with application-level patterns.

5. **Formal Verification**: Can rate limiting/backpressure policies be formally verified for correctness? (e.g., using TLA+)

6. **Adaptive Rate Limiting**: Research on ML-based adaptive rate limiting that learns traffic patterns and adjusts limits dynamically.

7. **Cost-Aware Scheduling**: Extending token bucket to support cost-aware scheduling where different operations have different "costs" in terms of resources.

## Research Freshness Note

- **Current research date**: 2026-09-26
- **Most recent source edits**: Wikipedia articles edited Aug-Sep 2026
- **AWS blog**: Updated May 2023
- **IEEE paper**: May 2018 (may need more recent datacenter research)
- **ScyllaDB blog**: Aug 2022

**Recommendation**: Re-verify cloud provider specific implementations and emerging patterns (e.g., gRPC flow control, Kubernetes-native queue management) as these evolve rapidly.