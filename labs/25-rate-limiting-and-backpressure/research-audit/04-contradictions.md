# Contradictions Audit: Rate Limiting & Backpressure

Target Lab: `labs/25-rate-limiting-and-backpressure`  
Audit Scope: Internal, Source, and Conceptual Consistency  
Audit Date: 2026-09-26  

---

## Analysis of Investigated Nuances

The research artifact `research/04-contradictions.md` reviewed five distinct technical areas where literature shows variations in terminology, algorithm choices, and parameter configurations:

1. **Token Bucket vs Leaky Bucket**:
   - Analyzed as mathematically equivalent duals with differing queuing vs meter semantics. No factual contradiction.
2. **Jitter Algorithms (Full vs Equal vs Decorrelated)**:
   - AWS Architecture blog compared multiple formulations; AWS SDK standardized on Full Jitter. No conflict.
3. **HTTP 429 vs 503**:
   - 429 applies to client-side rate limit quotas (RFC 6585); 503 applies to server-side resource saturation/overload (Stripe / RFC 9110). Well-delineated.
4. **Exponential Backoff Base Delays**:
   - Examined protocol-specific base parameters (50ms AWS API vs 500ms SIP vs 51.2μs Ethernet CSMA/CD). Context-specific calibration recognized.
5. **Queue Metrics (Queue Age vs Queue Depth vs Utilization)**:
   - Clarified that queue age measures latency/head-of-line blocking, while utilization measures resource saturation. Complementary rather than conflicting.

---

## Contradiction Checks

- **Research Plan vs Evidence**: Consistent.
- **Evidence vs Report**: Consistent. All findings in `05-report.md` match evidence entries in `03-evidence.md`.
- **Sources vs Claims**: Consistent. Source scopes align with cited claims.
- **Revision Artifacts vs Research Files**: Consistent. All revisions recorded in `research-revision/` are reflected in the target `research/` files.

---

## Conclusion

No material contradictions found.
