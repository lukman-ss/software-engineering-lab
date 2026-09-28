# Accuracy vs Research (05-report.md, 11 findings)

PASS — all core claims trace to research:

- F1 full-duplex non-HTTP, Upgrade + 101 → draft L24, L44 correct.
- F2 text 0x1 / binary 0x2 → draft L24, tests description correct.
- F3 no WS auto-reconnect, RFC 6455 §7.2.3 implement concern, ReconnectingWebSocket gap-filler → draft L94-96 correct, hedged.
- F4 SSE unidirectional → draft L26, brief, takeaways correct.
- F5 EventSource auto-reconnect, retry, Last-Event-ID → draft L26, L94, L104 correct.
- F6 SSE UTF-8 only, no binary → draft L7, L155, takeaways #5 correct.
- F7 WS over H2 needs RFC 8441 Extended CONNECT :protocol=websocket (Upgrade/Connection forbidden) → draft L102 correct.
- F8 SSE native H2, no extension → draft L102, takeaways #2 correct.
- F9 H1.1 ~6 conn/origin, H2 default 100 streams → draft L104 correct, attributes as de-facto browser policy.
- F10 proxy_buffering off, proxy_read_timeout, X-Accel-Buffering: no; AWS ELB sticky for WS → draft L100 correct, matches MEDIUM-confidence source.
- F11 100k = FD limits + per-conn memory + broker (Redis Pub/Sub), epoll, C10k → draft L106, L156 correctly labeled architectural principle, not lab benchmark.

No research contradiction. No unsupported protocol claim.
