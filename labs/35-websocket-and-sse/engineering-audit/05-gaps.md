# Engineering Audit Gaps

## Gap 1: Static Masking Key in WebSocket WriteFrame

- Type: `IMPLEMENTATION_OVERCLAIM`
- Severity: LOW
- Location: `internal/ws/ws.go:158`
- Description: `WriteFrame` with `mask=true` uses hardcoded mask key `[4]byte{0x12, 0x34, 0x56, 0x78}` instead of cryptographically random entropy per RFC 6455 §5.3.
- Impact: Security limitation only. Protocol parsing works correctly with any valid 4-byte mask key.
- Recommendation: Add `crypto/rand` generation if strict RFC compliance is desired in future revisions.

---

## Gap 2: Untested Medium/Large WebSocket Payload Length Framing (16-bit and 64-bit)

- Type: `MISSING_EDGE_CASE`
- Severity: LOW
- Location: `internal/ws/ws.go:94-106`, `internal/ws/ws.go:145-155`
- Description: Extended 126 (16-bit) and 127 (64-bit) payload length framing paths are implemented in code but not exercised by any test in `tests/protocol_test.go`. All tests send payloads <= 125 bytes.
- Impact: Framing logic is visually correct (`binary.BigEndian`), but unverified by unit/integration tests.
- Recommendation: Add a test sending a 256-byte payload to exercise length=126 framing.

---

## Gap 3: Missing Negative Tests for WebSocket Handshake and Corrupted Frames

- Type: `MISSING_TEST`
- Severity: LOW
- Location: `internal/ws/ws.go:31-56`
- Description: No tests verify server behavior when `Upgrade` header is missing, `Sec-WebSocket-Key` is empty, or connection drops unexpectedly.
- Impact: Core happy path verified; error edge cases unverified by automated tests.
- Recommendation: Add negative handshake test cases asserting `http.StatusBadRequest`.

---

## Gap 4: Unbounded Memory in SSE Hub History

- Type: `MISSING_EDGE_CASE`
- Severity: LOW
- Location: `internal/sse/sse.go:40,63`
- Description: `h.history` slice appends indefinitely on every `Broadcast` without compaction, ring-buffering, or TTL-based eviction.
- Impact: In a long-running production service this causes memory exhaustion. For this lab demo/test harness it has no effect.
- Recommendation: Document as an intentional pedagogical simplification or implement a ring buffer.
