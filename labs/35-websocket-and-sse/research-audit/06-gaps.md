# Research Gap Analysis: Lab 35 (WebSocket vs SSE Research)

## Gap 1

Type:
WEAK_SOURCE

Severity:
MEDIUM

Location:
`03-evidence.md: Evidence 12 & 13`, `05-report.md: Finding 11`

Problem:
The claims concerning 100k concurrent connections resource scaling (file descriptor limits, runtime memory overhead, Redis Pub/Sub cluster routing) cite the internal lab specification as the sole source rather than external primary technical sources, benchmarks, or system engineering literature (e.g. C10K/C10M papers or OS kernel socket benchmarks).

Required Revision:
When implementing the lab and writing educational materials, substantiate the scaling parameters with concrete runtime measurements (e.g., Goroutine memory allocation baseline in Go runtime docs or Linux socket memory buffers `tcp_rmem`/`tcp_wmem`).

Can Be Approved Without Fix:
YES (Research agent noted this constraint under "Limitations" and "Open Questions").

---

## Gap 2

Type:
OVERGENERALIZATION

Severity:
LOW

Location:
`05-report.md: Finding 10 & Executive Summary`

Problem:
Asserts that WebSocket requires sticky sessions on load balancers whereas SSE does not. In reality, established WebSocket connections stay on their target backend once upgraded; sticky sessions on load balancers are only required if reconnects must reach node-local state. Similarly, SSE connections reconnecting with `Last-Event-ID` require sticky routing if the event backlog is held in local instance memory rather than an external cache or broker.

Required Revision:
Clarify that sticky session necessity depends on node statefulness across reconnects rather than the underlying framing protocol itself.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
OUTDATED_SOURCE

Severity:
LOW

Location:
`02-sources.md: Source 6`, `05-report.md: Finding 7 & 8`

Problem:
RFC 7540 is cited for HTTP/2. RFC 7540 was formally obsoleted by RFC 9113 (HTTP/2) in June 2022. Additionally, in Finding 8's URL citation, the link typo points to `rfc7541` (HPACK) instead of `rfc7540`/`rfc9113`.

Required Revision:
Reference RFC 9113 alongside RFC 7540 and correct the URL typo in Finding 8.

Can Be Approved Without Fix:
YES
