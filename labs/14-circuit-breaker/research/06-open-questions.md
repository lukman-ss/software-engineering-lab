# Open Questions

1. **Distributed vs Local State**: Should circuit state be synchronized via Redis/Consul across a microservice fleet, or kept isolated in-process per replica to eliminate single points of failure?
2. **Error Classification Granularity**: In production HTTP APIs, which status codes count as circuit-breaker trip candidates (e.g. 500, 502, 503, 504 vs 429 Too Many Requests vs 400 Bad Request)? Typically 4xx client errors should never trip a circuit breaker.
3. **Adaptive Thresholds**: How do dynamic thresholding systems (e.g. Istio outlier detection) balance latency anomalies versus hard connection drops?
