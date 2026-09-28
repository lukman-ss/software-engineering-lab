# Audit Plan

## Target Lab
`labs/35-websocket-and-sse`

## Files Reviewed
- `labs/35-websocket-and-sse/research/01-plan.md`
- `labs/35-websocket-and-sse/research/02-sources.md`
- `labs/35-websocket-and-sse/research/03-evidence.md`
- `labs/35-websocket-and-sse/research/04-contradictions.md`
- `labs/35-websocket-and-sse/research/05-report.md`
- `labs/35-websocket-and-sse/research/06-open-questions.md`

## Claims To Verify
1. WebSocket is a bidirectional full-duplex protocol over a single TCP connection initiated via HTTP/1.1 Upgrade (RFC 6455).
2. WebSocket frame format supports text (UTF-8) and binary payloads (RFC 6455).
3. WebSocket specification does not mandate client auto-reconnect or backpressure mechanisms (RFC 6455, MDN).
4. SSE (`EventSource`) is strictly unidirectional (server-to-client) using `text/event-stream` MIME type (WHATWG, MDN).
5. SSE specifies native browser auto-reconnect, reconnection delay, and `Last-Event-ID` tracking (WHATWG HTML Standard Section 9.2).
6. SSE is restricted to UTF-8 text framing only (WHATWG HTML Standard Section 9.2).
7. WebSocket over HTTP/2 cannot use HTTP/1.1 Upgrade headers and requires RFC 8441 Extended CONNECT method with `:protocol = websocket` (RFC 8441, RFC 7540).
8. SSE runs natively over HTTP/2 multiplexed streams without protocol extension (RFC 7540, WHATWG).
9. Browsers enforce a ~6 connection per domain limit over HTTP/1.1 for SSE, mitigated to default 100 streams under HTTP/2 (MDN, RFC 7540).
10. Reverse proxies require specific configuration (disabling buffering via `proxy_buffering off`, tuning read timeouts) for real-time streaming (NGINX docs).
11. 100,000 concurrent connection scaling requirements (OS file descriptors `ulimit -n`, per-connection runtime memory, Redis pub/sub broker).

## Code To Execute
None. PIPELINE OVERRIDE: Audit research only. No code or implementation files exist in this stage.

## Primary Risks
- Untracked or arbitrary claims regarding 100k scaling numbers (file descriptors, memory footprint per connection).
- Lack of independent primary source for internal lab references cited in Evidence 12 and 13.
- Inaccurate RFC references or invalid URLs (e.g. RFC 7541 cited in Finding 8 when RFC 7540 was meant).
- Overgeneralization of browser connection caps (HTTP/1.1 6-connection limit treated as universal standard rather than browser implementation policy).

## Audit Strategy
1. Source verification: Verify reachability, authenticity, publisher, and relevance of all 7 cited sources in `02-sources.md` plus inline sources in `05-report.md`.
2. Claim validation: Cross-reference findings in `05-report.md` and evidence entries in `03-evidence.md` against official normative standards (RFC 6455, RFC 8441, RFC 7540, WHATWG).
3. Contradiction & gap detection: Evaluate internal consistency between research documents and pinpoint unverified assumptions (especially regarding 100k scaling and proxy behavior).
4. Verdict determination: Formulate objective verdict based on evidence reliability and pipeline readiness.
