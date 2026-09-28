# Audit Verdict

Target Lab: `labs/35-websocket-and-sse`

Audit Date: 2026-09-28

## Summary

Major Claims Reviewed: 11
Sources Reviewed: 8 (7 in source catalogue, 1 additional in report)
Unsupported Claims: 0
Contradictions: 1 (overgeneralized sticky session claim across proxies)
Code Issues: NOT_APPLICABLE (research-only audit per pipeline override)
Test Failures: NOT_APPLICABLE (research-only audit per pipeline override)
Research Gaps: 3 (1 MEDIUM, 2 LOW)

## Quality Gates

Source Integrity:
PASS

Claim Support:
PASS

Internal Consistency:
PASS

Code Correctness:
NOT_APPLICABLE

Tests:
NOT_APPLICABLE

Documentation Accuracy:
PASS

## Blocking Issues

None.

## Non-Blocking Issues

1. **Sticky Session Generalization (`05-report.md: Finding 10`)**: The claim that WebSocket requires sticky sessions on load balancers while SSE does not is an oversimplification. Once established, WebSocket connections are pinned to the TCP socket regardless of sticky cookies; sticky sessions are only relevant if client reconnects must reach node-local state. Similarly, stateful SSE implementations require stickiness unless backed by a pub/sub bus.
2. **Citation Source for 100k Scaling (`03-evidence.md: Evidence 12 & 13`, `05-report.md: Finding 11`)**: Resource scaling numbers and patterns rely on the internal lab specification rather than external empirical benchmark citations. The Research Agent properly disclosed this limitation in `05-report.md` and `06-open-questions.md`.
3. **HTTP/2 Spec Reference & Typo (`05-report.md: Finding 8`)**: Cited URL has a typo (`rfc7541` instead of `rfc7540`). Also RFC 7540 was superseded by RFC 9113.

## Required Revisions

1. During lab implementation / content authoring, qualify sticky session recommendations to distinguish between connection persistence and reconnection routing to stateful nodes.
2. Benchmark runtime memory per connection (Go goroutine vs Node event loop) directly during the implementation phase.
3. Update HTTP/2 citations to RFC 9113 and fix the minor URL typo in Finding 8.

## Final Status

APPROVED_WITH_WARNINGS
