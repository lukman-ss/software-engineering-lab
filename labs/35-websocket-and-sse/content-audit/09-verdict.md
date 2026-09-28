# Content Audit Verdict

Target Lab: labs/35-websocket-and-sse
Audit Date: 2026-09-28
Auditor: Technical Writer Auditor

## Verdict

APPROVED_WITH_WARNINGS

## Warnings Summary

1. **EventSource "exponential backoff" wording inaccurate** — WHATWG defines fixed `retry`-based reconnection time; no source mandates exponential backoff. Fix: replace with "reconnect otomatis dengan jeda sesuai field retry (implementation-defined)".
2. **Demo case study event-count inconsistency** — Case study L124 says replay "ID 2 dan 3" but demo only broadcasts 2 events (IDs 1,2) and output shows only ID 2. Fix: change to "event ID 2 harus di-replay".
3. **SSE no-sticky-session claim oversimplified** — In-memory per-node history means cross-node Last-Event-ID replay requires shared broker regardless of protocol. Fix: add clause about shared state.
4. **WebSocket sticky-session stated as absolute** — Also avoidable with a centralized message broker (Redis Pub/Sub per Finding 11). Minor softening recommended.
5. **WS exempt from browser 6-conn limit unsourced** — No research source asserts exemption; browser limits exist for WS too. Minor softening or removal.

All warnings are non-blocking. Content is technically accurate, well-formatted, and properly aligned with research and engineering implementation. No hallucinations, fabricated test results, or platform-specific biases detected.
