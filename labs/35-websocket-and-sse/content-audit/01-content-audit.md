# Content Audit Report

Target Lab: labs/35-websocket-and-sse
Audit Date: 2026-09-28
Auditor: Technical Writer Auditor
Scope: Content accuracy against engineering implementation and research evidence
Directive: Audit content only — do not audit research/code, do not modify files

---

## 1. Summary

Content files reviewed:
- 01-content-brief.md
- 02-master-draft.md
- 03-code-snippets.md
- 04-diagrams.md
- 05-key-takeaways.md
- 06-source-map.md

Research audit status: APPROVED (11 claims, 0 unsupported, 0 contradictions)
Engineering audit status: APPROVED (4 tests PASS, race PASS, demo PASS)
Content audit result: APPROVED_WITH_WARNINGS

---

## 2. Alignment Assessment

### 2.1 Research Alignment
- Directionality (WS full-duplex vs SSE unidirectional): matches Finding 1/4. PASS.
- Binary vs UTF-8-only: matches Finding 2/6. PASS.
- Auto-reconnect + Last-Event-ID: matches Finding 5. PASS, except backoff wording (§3.1).
- RFC 8441 / HTTP-2 native SSE: matches Finding 7/8. PASS.
- 6-connection HTTP/1.1 limit + 100-stream HTTP/2: matches Finding 9. PASS.
- Proxy buffering/timeout: matches Finding 10 (MEDIUM confidence, correctly hedged in draft L106). PASS.
- 100k scaling as architectural principle, not benchmark: draft L106 explicitly hedges. PASS — matches research Limitations and research-audit Gap 1.
- RFC 9113 (obsoletes RFC 7540) cited in Sources. PASS — resolves research-audit Gap 3.

### 2.2 Code Alignment
- All 10 snippets in 03-code-snippets.md verified byte-for-byte against `internal/sse/sse.go`, `internal/ws/ws.go`, `internal/server/server.go`, `tests/protocol_test.go`. PASS.
- Snippet 6 correctly discloses static mask key `{0x12,0x34,0x56,0x78}` vs RFC 6455 §5.3 randomness. Matches engineering-audit Gap 1. PASS.
- Snippet 7 line range `server.go:26-56` correct for /ws handler. PASS.
- Test assertions quoted exactly: `TestSSE_Formatting` expected string, `TestWebSocket_TextAndBinary` payloads `0x01..0x04` → reversed. PASS.
- Draft correctly lists 3 endpoints (/sse, /ws, /publish) where README lists only 2. Draft is more accurate than README. PASS.
- Headers (`text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`), 15s keep-alive comment, capacity-64 non-blocking broadcast, `ID > lastID` replay: all match code. PASS.
- Content-brief Warnings (static mask, unbounded history, untested extended framing) match engineering-audit Gaps 1/2/4 exactly. PASS.
- Demo vs test payloads correctly separated: L74 cites test bytes `01 02 03 04`, L137 cites demo bytes `DEADBEEF`. No confusion. PASS.

### 2.3 Diagrams / Source Map
- Diagram 2 (SSE replay ID:2 after Last-Event-ID:1) matches demo 2-event scenario. PASS.
- Diagram 3 (FIN/opcode/length 7-16-64 + XOR formula) matches ws.go Read/WriteFrame. PASS.
- Diagram 4 (history + clients map cap 64) matches Hub struct. PASS.
- Diagram 5 includes /publish. PASS.
- 06-source-map.md traces every claim to research finding + code + test + audit gap, and explicitly notes "HTTP/2 multiplexing not demonstrated" and "100k not empirically benchmarked". PASS.

---

## 3. Issues Found

### 3.1 Low-Medium — "Exponential backoff" for EventSource unsupported
- Location: content/02-master-draft.md:94 ("secara otomatis reconnect dengan exponential backoff").
- Reality: WHATWG defines `retry` field reconnection time; research Evidence 6 notes "reconnection time is implementation-defined", Open Questions #6 says "exact backoff algorithm ... implementation-defined". No source claims exponential backoff.
- Impact: Readers may expect spec-mandated exponential backoff that browsers do not guarantee.
- Recommendation: Replace with "reconnect otomatis dengan jeda sesuai field `retry` (implementation-defined, biasanya beberapa detik)".

### 3.2 Low-Medium — Case study event-count inconsistency
- Location: content/02-master-draft.md:124 ("broadcast 2 event ... event ID 2 dan 3 harus di-replay").
- Reality: Demo (`cmd/demo/main.go:28-29`) broadcasts 2 events (ID 1, 2); replay is only ID 2 (demo output block below it correctly shows only id:2). The "ID 2 dan 3" numbering belongs to the 3-event test (`TestSSE_LastEventID_Resumption`), not the demo. Draft conflates the two.
- Impact: Minor numeric inconsistency; output block itself is correct.
- Recommendation: Change to "event ID 2 harus di-replay" in demo case-study paragraph.

### 3.3 Low — SSE "no sticky session" oversimplified
- Location: content/02-master-draft.md:32 ("proxy tidak memerlukan sticky session ... bisa ke node manapun").
- Reality: True for connection routing, but `Last-Event-ID` replay against in-memory `history` (sse.go:40,63, unbounded, per-node) fails cross-node without shared broker. Research Finding 11 requires Redis Pub/Sub (or equivalent) for multi-node broadcast — applies to SSE history as well as WS.
- Impact: Readers may assume SSE multi-node replay works with zero shared state.
- Recommendation: Add one clause: "untuk replay lintas node tetap perlu shared history/broker seperti Redis Pub/Sub".

### 3.4 Low — WS sticky-session stated as absolute
- Location: content/02-master-draft.md:30 ("harus dikonfigurasi untuk mendukung WebSocket sticky sessions").
- Reality: Required only for in-memory state without broker; with Redis Pub/Sub (Finding 11) sticky is avoidable. Research-audit Claim 10 marks proxy/sticky claims PARTIAL / implementation-specific.
- Impact: Overgeneralization; scenario is framed as "Bayangkan" so clearly illustrative.
- Recommendation: Soften to "umumnya membutuhkan sticky session kecuali ada broker pesan terpusat".

### 3.5 Low — WS exempt from browser connection limits unsourced
- Location: content/02-master-draft.md:104 ("WebSocket tidak terkena batasan ini karena menggunakan koneksi TCP terpisah").
- Reality: Research Finding 9 documents only the SSE 6-connection cap; no research source asserts WS exemption. Browsers do impose separate (higher) WS limits.
- Impact: Minor unsupported extrapolation.
- Recommendation: Soften to "WebSocket tidak dihitung dalam batas 6 koneksi HTTP/1.1 tersebut (tetapi browser tetap membatasi jumlah koneksi WS secara terpisah)" or drop sentence.

---

## 4. Verification Evidence

| Check | Result |
|-------|--------|
| Snippets match implementation byte-for-byte (10/10) | PASS |
| Test assertions quoted exactly (formatting string, payloads, Last-Event-ID) | PASS |
| Headers, keep-alive interval, channel capacity, replay predicate | PASS |
| RFC 6455 / RFC 8441 / RFC 9113 / WHATWG citations | PASS |
| 100k + HTTP/2 scope correctly hedged as non-benchmarked | PASS |
| Engineering warnings (mask key, unbounded history, untested framing) disclosed | PASS |
| Research audit verdict | APPROVED |
| Engineering audit verdict | APPROVED |
| Docs vs code consistency | PASS |

---

## 5. Hallucination Check

| Claim | Source | Verdict |
|-------|--------|---------|
| WS full-duplex, Upgrade → 101, opcodes 0x1/0x2 | RFC 6455, ws.go | PASS |
| SSE unidirectional, UTF-8 only, EventSource + Last-Event-ID | WHATWG, MDN | PASS |
| RFC 8441 Extended CONNECT `:protocol=websocket` | RFC 8441 | PASS |
| SSE native HTTP/2, ~6 conn HTTP/1.1 limit | MDN, RFC 9113 | PASS |
| proxy_buffering / X-Accel-Buffering / read_timeout | NGINX docs, sse.go | PASS |
| EventSource exponential backoff (spec-mandated) | None — impl-defined | WARN (§3.1) |
| Demo replays "ID 2 dan 3" | Demo has only ID 1,2 | WARN (§3.2) |
| SSE never needs sticky / shared state | Contradicted by in-memory history | WARN (§3.3) |
| LLM streaming, notifications, dashboards = SSE; chat/games/collab = WS | Research conclusion | PASS |
| No platform-specific bias | Neutral Go stdlib framing | PASS |

No fabricated test results. No invented APIs. No blocking hallucinations.

---

## 6. Recommendation Summary

- No blocking issues.
- Content accuracy: HIGH.
- Clarity: HIGH.
- Formatting: Consistent, Indonesian, appropriate for backend/architect audience.

Actions (all non-blocking):
1. Fix "exponential backoff" wording (§3.1).
2. Fix demo replay numbering "2 dan 3" → "2" (§3.2).
3. Optional: qualify SSE sticky-session (§3.3), WS sticky absolute (§3.4), WS limit exemption (§3.5).

---

## 7. Final Verdict

APPROVED_WITH_WARNINGS
