# Research Plan

## Research Topic
Rate Limiting & Backpressure — System Design for Scalability

## Objective
Investigate and document evidence-based practices for rate limiting algorithms, backpressure mechanisms, and queue management in distributed systems. Focus on practical engineering patterns used in production systems.

## Research Questions

### Rate Limiting
1. What are the major rate limiting algorithms (Token Bucket, Leaky Bucket, Fixed Window, Sliding Window) and their trade-offs?
2. How do production systems implement rate limiting at different layers (API Gateway, Application, Database)?
3. What are the best practices for multi-tenant rate limiting and fairness?
4. How to handle distributed rate limiting across multiple instances?

### Backpressure
1. What are the core backpressure patterns (explicit signaling, implicit via queue saturation, reactive streams)?
2. How does Little's Law apply to queue capacity planning?
3. What are the retry strategies (exponential backoff, jitter, circuit breaker) and their effectiveness?
4. How to implement backpressure in message queues (Redis, RabbitMQ, Kafka)?

### Queue Management
1. What are the dangers of unbounded queues and memory pressure?
2. How to implement fair queueing and priority queues?
3. What metrics should be monitored (queue depth, oldest job age, processing rate vs arrival rate)?
4. How to handle dead letter queues and retry limits?

### Practical Applications
1. What are the real-world implementations in major systems (Netflix, Uber, Stripe, etc.)?
2. How to calculate capacity requirements using throughput math?
3. What are common anti-patterns and failures?

## Search Strategy

### Primary Sources (Tier 1)
- Official documentation: Redis, RabbitMQ, Kafka, gRPC, HTTP specs
- Academic papers on queueing theory, Little's Law
- Cloud provider docs: AWS API Gateway, Google Cloud Load Balancing
- Standards: RFC 6585 (429 status), Reactive Streams specification

### Secondary Sources (Tier 2)
- Engineering blogs from Netflix, Uber, Stripe, Shopify, GitHub
- Technical articles from ACM Queue, Communications of ACM
- CNCF / Kubernetes documentation on backpressure

### Tertiary Sources (Tier 3)
- Community discussions (Hacker News, Reddit r/sysadmin, Stack Overflow)
- Personal technical blogs

## Expected Primary Sources

1. **Token Bucket Algorithm** - Original paper or standard references
2. **Little's Law** - John Little's original proof / queueing theory texts
3. **RFC 6585** - HTTP 429 status code specification
4. **Reactive Streams Specification** - Backpressure protocol
5. **Netflix Tech Blog** - Rate limiting at scale
6. **Uber Engineering** - Backpressure in dispatch systems
7. **Stripe Engineering** - API rate limiting
8. **Redis Documentation** - Redis rate limiting patterns (Redis Cell)
9. **gRPC Documentation** - Flow control / backpressure
10. **Kubernetes Documentation** - Horizontal Pod Autoscaler, queue proxy

## Risks / Unknowns

- Some proprietary implementations may not be publicly documented
- Distributed rate limiting consistency models vary (eventual vs strong)
- Exact algorithms used by major cloud providers may be opaque
- Performance benchmarks vary significantly by workload
- Indonesian language resources may be limited; primary sources in English