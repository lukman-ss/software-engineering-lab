# Claim Audit

Target Lab: `labs/35-websocket-and-sse`

---

## Claim 1: WebSocket Full-Duplex Bidirectional Communication

Claim:
WebSocket provides a two-way (full-duplex) communication channel over a single TCP connection initiated via HTTP/1.1 Upgrade resulting in `101 Switching Protocols`.

Location:
`05-report.md` (Finding 1), `03-evidence.md` (Evidence 1, Evidence 2)

Evidence Provided:
RFC 6455 Section 1.1, Section 1.2, Section 4.2.2.

Source:
RFC 6455

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects normative specification.

---

## Claim 2: Payload Types Supported

Claim:
WebSocket data frames support UTF-8 text and arbitrary binary payloads, whereas SSE supports UTF-8 text only.

Location:
`05-report.md` (Finding 2, Finding 6), `03-evidence.md` (Evidence 3, Evidence 7)

Evidence Provided:
RFC 6455 Section 1.2; WHATWG HTML Standard Section 9.2.1.

Source:
RFC 6455, WHATWG HTML Living Standard

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Normative references accurately support both protocol framing rules.

---

## Claim 3: Automatic Reconnection Support

Claim:
Browser `EventSource` (SSE) provides built-in standardized auto-reconnection with `Last-Event-ID` tracking, whereas `WebSocket` has no native browser auto-reconnect or backpressure mechanism.

Location:
`05-report.md` (Finding 3, Finding 5), `03-evidence.md` (Evidence 4, Evidence 6)

Evidence Provided:
WHATWG HTML Standard Section 9.2.3 / 9.2.4; MDN EventSource & WebSocket.

Source:
WHATWG HTML Living Standard, MDN

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate distinction between the two web APIs.

---

## Claim 4: HTTP/2 Integration (RFC 8441 vs Native SSE)

Claim:
WebSocket requires RFC 8441 (Extended CONNECT method with `:protocol = websocket`) to run over HTTP/2 because HTTP/2 forbids connection-wide headers (`Upgrade`, `Connection`). In contrast, SSE is standard HTTP semantics and runs natively over HTTP/2 multiplexed streams without extension.

Location:
`05-report.md` (Finding 7, Finding 8), `03-evidence.md` (Evidence 8, Evidence 9)

Evidence Provided:
RFC 8441 Section 1 & Section 5; RFC 7540; WHATWG HTML Standard Section 9.2.

Source:
RFC 8441, RFC 7540, WHATWG

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
RFC 8441 text and HTTP/2 stream multiplexing constraints are accurately cited. (Minor URL typo in Finding 8 notes `rfc7541` instead of `rfc7540`).

---

## Claim 5: Browser Connection Limits on HTTP/1.1

Claim:
Browsers enforce a maximum limit of roughly 6 simultaneous HTTP/1.1 connections per origin, constraining SSE over HTTP/1.1 across multiple tabs, which HTTP/2 resolves via negotiated stream limits (default 100).

Location:
`05-report.md` (Finding 9), `03-evidence.md` (Evidence 10)

Evidence Provided:
MDN EventSource documentation; RFC 7540 Section 5.1.2.

Source:
MDN, RFC 7540

Source Actually Supports Claim:
PARTIAL

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
The 6-connection limit is a de facto browser implementation rule (Chrome, Firefox), not an IETF standard requirement. The research acknowledges this in limitations, which is appropriate.

---

## Claim 6: Reverse Proxy Buffering & Timeout Requirements

Claim:
Reverse proxies (e.g. NGINX) buffer responses by default, requiring `proxy_buffering off` and tuned `proxy_read_timeout` to prevent buffering or dropping long-lived streaming connections (SSE/WebSocket).

Location:
`05-report.md` (Finding 10), `03-evidence.md` (Evidence 11)

Evidence Provided:
NGINX ngx_http_proxy_module documentation.

Source:
NGINX Docs

Source Actually Supports Claim:
PARTIAL

Classification:
FACT / IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
`proxy_buffering off` is specific to NGINX/reverse proxy behavior. The source confirms the directive and mechanism.

---

## Claim 7: 100,000 Concurrent Connections Resource Constraints

Claim:
Scaling WebSocket or SSE to ~100,000 concurrent connections is constrained by OS file descriptor limits (`ulimit -n`), runtime per-connection memory allocation, and requires a centralized broker (Redis Pub/Sub) for multi-node broadcast.

Location:
`05-report.md` (Finding 11), `03-evidence.md` (Evidence 12, Evidence 13)

Evidence Provided:
Internal lab specification citation; C10K/C10M networking principles.

Source:
Internal Lab Specification

Source Actually Supports Claim:
PARTIAL

Classification:
HYPOTHESIS / ARCHITECTURAL PRINCIPLE

Severity:
MEDIUM

Notes:
While technically accurate from systems engineering principles, citing internal lab prompts as the primary source is self-referential. No empirical memory numbers or OS benchmark citations are provided. The research report appropriately classifies this under `Limitations` and `Confidence: MEDIUM`.
