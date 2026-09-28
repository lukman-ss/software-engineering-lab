# Research Plan

## Research Topic

WebSocket vs Server-Sent Events (SSE) — choosing the right real-time protocol.

Lab: 35-websocket-and-sse
Research date: 2026-09-28
Language of lab: Bahasa Indonesia
Audience: Senior Software Engineer

## Objective

Collect authoritative, cross-checked evidence on protocol differences, operational complexity, HTTP/2 behavior, reconnection, proxy/load-balancer constraints, and scaling of concurrent connections. Output is research-only; no implementation.

## Research Questions

1. What are the protocol-level differences between WebSocket and SSE (handshake, framing, MIME type, transport)?
2. Is communication direction the correct primary decision criterion (full-duplex vs unidirectional server→client)?
3. Which payload types does each protocol support (binary vs UTF-8 text)?
4. Does the browser EventSource API provide automatic reconnection, Last-Event-ID, and event types that WebSocket does not?
5. How does each protocol interact with HTTP/2 (RFC 8441 vs native HTTP streaming)?
6. What proxy, firewall, caching, and load-balancing constraints apply (sticky sessions, buffering, timeout)?
7. What is the operational cost of WebSocket (upgrade, heartbeat, reconnection) versus SSE for one-way server push?
8. What OS and runtime limits matter at ~100k concurrent connections (file descriptors, memory per connection, goroutine vs event loop)?
9. How is multi-node broadcast typically solved (Redis Pub/Sub or equivalent)?
10. Which use cases are documented as SSE-appropriate (notifications, LLM token streaming, dashboards) versus WebSocket-appropriate (chat, games, collaborative editing)?
11. Are there important caveats that contradict the lab table (browser connection limits on HTTP/1.1, SSE binary limitation, HTTP/2 WebSocket support, IE/legacy browsers)?

## Search Strategy

1. Primary specifications: RFC 6455, RFC 8441, WHATWG HTML EventSource, MDN.
2. HTTP working-group / IETF notes on WebSocket over HTTP/2 and HTTP/3.
3. Official browser and runtime docs (Chrome, MDN, WHATWG; Node, Go net/http).
4. Cloud / proxy vendor docs (nginx, AWS ALB, Cloudflare, Envoy) for buffering, timeouts, sticky sessions.
5. LLM streaming docs (OpenAI SSE, Anthropic SSE) as documented production use of SSE.
6. Scaling literature: OS ulimit, C10K/C10M context, Redis Pub/Sub docs — treat blog numbers as LOW unless corroborated.

## Expected Primary Sources

- RFC 6455 — The WebSocket Protocol
- RFC 8441 — Bootstrapping WebSockets with HTTP/2
- WHATWG HTML Living Standard — Server-sent events / EventSource
- MDN: WebSocket, EventSource, Using server-sent events
- OpenAI / Anthropic API streaming documentation (SSE)
- nginx / Envoy / AWS docs on WebSocket and proxy buffering
- Redis Pub/Sub documentation
- Go/Node runtime documentation on connection and goroutine/event-loop cost (if available)

## Risks / Unknowns

- Memory-per-connection numbers vary by language, buffer size, TLS, and OS; likely not transferable.
- HTTP/2 browser connection limits for SSE may differ from HTTP/1.1 (6-connection-per-host) — must verify.
- WebSocket over HTTP/2 (RFC 8441) browser and proxy support may be incomplete.
- SSE over HTTP/2 vs HTTP/1.1 behavior in reverse proxies may be implementation-specific.
- "100k concurrent connections" is a scaling scenario, not a measured benchmark in this lab.
- Some popular comparison tables mix spec facts with folklore (proxy "cannot cache WebSocket", "SSE is always simpler").
