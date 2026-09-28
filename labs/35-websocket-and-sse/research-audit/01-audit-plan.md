# Audit Plan: Lab 35 (WebSocket vs SSE Research)

## Target Lab
`labs/35-websocket-and-sse`

## Scope of Audit
Research Agent outputs under `labs/35-websocket-and-sse/research/`:
- `01-plan.md`
- `02-sources.md`
- `03-evidence.md`
- `04-contradictions.md`
- `05-report.md`
- `06-open-questions.md`

*(Note: Per pipeline instruction, audit scope is research-only; implementation and code execution are skipped in this stage).*

## Files Reviewed
1. `labs/35-websocket-and-sse/research/01-plan.md`
2. `labs/35-websocket-and-sse/research/02-sources.md`
3. `labs/35-websocket-and-sse/research/03-evidence.md`
4. `labs/35-websocket-and-sse/research/04-contradictions.md`
5. `labs/35-websocket-and-sse/research/05-report.md`
6. `labs/35-websocket-and-sse/research/06-open-questions.md`

## Claims To Verify
1. WebSocket is a full-duplex protocol over a single TCP connection initiated via HTTP/1.1 Upgrade handshake (`101 Switching Protocols`) (RFC 6455).
2. WebSocket supports both UTF-8 text and binary payloads (RFC 6455).
3. WebSocket does not have native/mandated auto-reconnection in the browser API (RFC 6455 / MDN).
4. Server-Sent Events (SSE) via `EventSource` is strictly unidirectional (server → client) over standard HTTP using `text/event-stream` (WHATWG HTML / MDN).
5. `EventSource` provides automatic reconnection with `Last-Event-ID` tracking (WHATWG HTML).
6. SSE supports UTF-8 text payloads only (WHATWG HTML).
7. WebSocket over HTTP/2 requires RFC 8441 Extended CONNECT (`:protocol = websocket`) because HTTP/2 bans connection-specific headers like `Upgrade` (RFC 8441 / RFC 7540).
8. SSE runs natively over HTTP/2 multiplexing without protocol extensions (RFC 7540 / WHATWG).
9. HTTP/1.1 browsers enforce a ~6 connections-per-origin limit which restricts SSE concurrent tabs, whereas HTTP/2 allows multiplexing up to negotiated stream limits (default 100) (MDN / RFC 7540).
10. Reverse proxies (e.g. NGINX) buffer responses by default, requiring tuning of `proxy_buffering` and `proxy_read_timeout` for real-time streams (NGINX docs).
11. 100k concurrent connections resource scaling depends on OS file descriptors (`ulimit -n`), runtime memory per connection/goroutine/event-loop, and multi-node broker architecture (Redis Pub/Sub) (Architectural patterns / internal lab spec).

## Primary Risks
- Use of internal lab specification as a primary source for scaling/architecture claims without external peer-reviewed benchmark sources.
- RFC 7541 cited instead of RFC 7540 in Report Finding 8 URL (`rfc7541` is HPACK, whereas text references RFC 7540 HTTP/2).
- Verification of external URL availability and correctness of cited normative text.

## Audit Strategy
1. Cross-reference all 7 sources in `02-sources.md` against authoritative IETF RFCs, WHATWG standards, and MDN specifications.
2. Audit all 13 evidence items in `03-evidence.md` and 11 findings in `05-report.md` for factual correctness, classification, and source fidelity.
3. Check for internal contradictions, scope overgeneralizations, and undocumented assumptions across report files.
4. Evaluate research limitations and open questions in `06-open-questions.md`.
5. Produce final verdict based on quality gates.
