# Audit Plan

Target Lab: labs/35-websocket-and-sse

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`
- `research-revision/01-revision-plan.md`
- `research-revision/02-changes-made.md`
- `research-revision/03-revision-result.md`

## Claims To Verify
1. WebSocket is a full-duplex, non-HTTP protocol operating over a single TCP connection after an HTTP/1.1 `Upgrade` handshake (RFC 6455).
2. WebSocket supports binary and UTF-8 text frames natively.
3. WebSocket does not mandate browser auto-reconnection.
4. SSE (`EventSource`) is strictly unidirectional (server -> client) using MIME type `text/event-stream`.
5. SSE natively supports browser auto-reconnection and `Last-Event-ID` header.
6. SSE supports UTF-8 text only (no binary payload support).
7. WebSocket over HTTP/2 requires RFC 8441 (Extended CONNECT with `:protocol = websocket`).
8. SSE runs natively over HTTP/2 without protocol extensions.
9. Browser per-domain connection limit (6 over HTTP/1.1) constrains SSE; HTTP/2 removes this via stream multiplexing.
10. Reverse proxy configuration requires tuning `proxy_buffering` and `proxy_read_timeout`.
11. Scaling to 100k concurrent connections is governed by OS file descriptor limits (`ulimit -n`, `fs.file-max`), `epoll(7)`, per-connection runtime memory, and Redis Pub/Sub for inter-node routing.

## Code To Execute
- PIPELINE OVERRIDE: Code and implementation auditing is excluded in this stage. No code commands executed.

## Primary Risks
- Over-reliance on secondary/community docs (MDN, man pages, NGINX docs) without primary specification verification.
- Claiming 100k connection scaling and Redis Pub/Sub patterns as established facts without empirical benchmarking or primary performance source backing.
- Citing obsoleted standards (RFC 7540 instead of RFC 9113) or unverified vendor URL behaviors.

## Audit Strategy
1. Audit all 8 listed sources in `research/02-sources.md` for URL validity, tier classification, authority, and accurate scope.
2. Evaluate all major claims in `research/05-report.md` and `research/03-evidence.md` against cited sources.
3. Verify internal consistency across research files and revision logs.
4. Record research gaps and issue final audit verdict.
