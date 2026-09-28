# Claim Audit: Lab 35 (WebSocket vs SSE Research)

## Claim 1: Full-Duplex Bidirectional Channel for WebSocket
Claim: WebSocket is a two-way, bidirectional communication channel over a single TCP connection that operates outside standard HTTP/1.1 request-response lifecycle.
Location: `05-report.md: Finding 1`, `03-evidence.md: Evidence 1`
Evidence Provided: RFC 6455 Section 1.1 & 1.2
Source: RFC 6455
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by RFC 6455 Section 1.1.

---

## Claim 2: Payloads Supported by WebSocket
Claim: The base WebSocket protocol defines distinct frame types for binary data and text (UTF-8) data.
Location: `05-report.md: Finding 2`, `03-evidence.md: Evidence 3`
Evidence Provided: RFC 6455 Section 1.2 & Section 5.2 (Opcodes for text frame 0x1, binary frame 0x2)
Source: RFC 6455
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully supported.

---

## Claim 3: Automatic Reconnection in WebSocket
Claim: The WebSocket protocol and browser API do not mandate automatic reconnection; the client must implement reconnect logic.
Location: `05-report.md: Finding 3`, `03-evidence.md: Evidence 4`
Evidence Provided: RFC 6455 Section 7.2.3, MDN WebSocket
Source: RFC 6455 / MDN
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate. Third-party wrappers (e.g. reconnecting-websocket) are required for auto-reconnection in vanilla browser environments.

---

## Claim 4: Unidirectional Nature of Server-Sent Events (SSE)
Claim: Server-Sent Events deliver data only from the server to the client; there is no native mechanism for client-to-server messaging over the SSE stream.
Location: `05-report.md: Finding 4`, `03-evidence.md: Evidence 5`
Evidence Provided: WHATWG HTML Standard Section 9.2, MDN EventSource
Source: WHATWG HTML / MDN
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate. Client-to-server communication requires separate HTTP requests (e.g., fetch/POST).

---

## Claim 5: SSE Native Reconnection and Last-Event-ID
Claim: The browser-managed `EventSource` object automatically reconnects on dropped connections, supports a configurable reconnection time, and transmits the last event ID (`Last-Event-ID`) to the server on reconnection.
Location: `05-report.md: Finding 5`, `03-evidence.md: Evidence 6`
Evidence Provided: WHATWG HTML Standard Section 9.2.3, 9.2.4
Source: WHATWG HTML Living Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully specified in the living standard.

---

## Claim 6: SSE Payload Character Encoding Limitation
Claim: Server-Sent Events are restricted to UTF-8 text; the spec provides no mechanism for binary payloads.
Location: `05-report.md: Finding 6`, `03-evidence.md: Evidence 7`
Evidence Provided: WHATWG HTML Standard Section 9.2.1 ("Event streams are always decoded as UTF-8. There is no way to specify another character encoding.")
Source: WHATWG HTML Living Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Binary data requires encoding (such as Base64) at the application layer, adding size overhead.

---

## Claim 7: WebSocket over HTTP/2 via RFC 8441 Extended CONNECT
Claim: Because HTTP/2 forbids connection-wide headers such as `Upgrade` and `Connection`, the WebSocket handshake cannot operate over plain HTTP/2. RFC 8441 defines a new Extended CONNECT method using `:protocol = websocket` on a single HTTP/2 stream.
Location: `05-report.md: Finding 7`, `03-evidence.md: Evidence 8`
Evidence Provided: RFC 8441 Section 1 & Section 5, RFC 7540 Section 8.3
Source: RFC 8441 / RFC 7540
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate technical breakdown of why standard WebSocket upgrade fails on HTTP/2 multiplexing without RFC 8441.

---

## Claim 8: Native SSE Operation over HTTP/2
Claim: Because SSE uses standard HTTP semantics (`text/event-stream`), it runs natively over HTTP/2 with stream multiplexing without requiring protocol extensions like RFC 8441.
Location: `05-report.md: Finding 8`, `03-evidence.md: Evidence 9`
Evidence Provided: WHATWG HTML Standard, RFC 7540
Source: WHATWG HTML / RFC 7540
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate. Minor typo in URL citation in report (points to `rfc7541` instead of `rfc7540`), but factual content is correct.

---

## Claim 9: Browser HTTP/1.1 6-Connection Limit for SSE vs HTTP/2
Claim: Browsers enforce a maximum number of simultaneous HTTP/1.1 connections per origin (~6), limiting concurrent SSE tabs/streams. HTTP/2 removes this bottleneck via negotiated stream limits (default 100).
Location: `05-report.md: Finding 9`, `03-evidence.md: Evidence 10`
Evidence Provided: MDN EventSource, RFC 7540 Section 5.1.2
Source: MDN / RFC 7540
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: 6-connection cap is a browser client implementation policy (Chrome, Firefox, Safari) rather than an RFC mandate, properly identified as such in MDN.

---

## Claim 10: Reverse Proxy Buffering and Sticky Session Behavior
Claim: Reverse proxies buffer responses by default, requiring tuning of `proxy_buffering` and `proxy_read_timeout`. WebSocket requires sticky sessions on load balancers, whereas SSE does not.
Location: `05-report.md: Finding 10`, `03-evidence.md: Evidence 11`
Evidence Provided: NGINX Docs (`ngx_http_proxy_module`), AWS ELB Docs
Source: NGINX / AWS Docs
Source Actually Supports Claim: PARTIAL
Classification: INTERPRETATION / IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes:
- Buffering/timeout claim: Supported by NGINX docs (`proxy_buffering off` is standard for SSE/streaming).
- Sticky session claim: Overgeneralized. Established WebSocket connections are long-lived TCP/stream channels pinned to an instance anyway; sticky sessions on load balancers are only required if reconnects must reach the exact same backend state. Similarly, SSE streams holding state or local event offsets without a shared broker also require sticky routing.

---

## Claim 11: Scaling to 100k Concurrent Connections Constraints
Claim: Operating at ~100,000 concurrent connections is constrained by OS file descriptor limits (`ulimit -n`), runtime memory overhead per connection buffer (e.g. Go goroutine stacks vs Node.js event loop), and requires a pub/sub message broker (e.g. Redis Pub/Sub) for multi-node broadcast.
Location: `05-report.md: Finding 11`, `03-evidence.md: Evidence 12 & 13`
Evidence Provided: Internal lab specification #35
Source: Lab specification (Internal)
Source Actually Supports Claim: PARTIAL
Classification: HYPOTHESIS / ARCHITECTURAL PATTERN
Severity: MEDIUM
Notes:
- The principles (file descriptors, per-connection memory, pub/sub for clustering) are standard distributed systems knowledge.
- However, using the internal lab specification as the sole cited source for empirical scaling behavior lacks external primary benchmark evidence. The report appropriately acknowledges this in "Limitations" and "Open Questions".
