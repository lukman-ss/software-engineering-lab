# Engineering Audit Verdict

Target Lab: labs/35-websocket-and-sse
Audit Date: 2026-09-28

## Summary

Code Files Reviewed: 4 (internal/sse/sse.go, internal/ws/ws.go, internal/server/server.go, cmd/demo/main.go)
Tests Reviewed: 1 file, 4 tests (tests/protocol_test.go: TestSSE_Formatting, TestSSE_LastEventID_Resumption, TestWebSocket_TextAndBinary, TestSSE_Concurrency)
Commands Executed: go build ./..., go test ./..., go test -count=1 -v ./tests/..., go test -race ./tests/..., go vet ./..., go run ./cmd/demo
Failures: 0
Warnings: 6 (5 LOW code robustness + 1 LOW docs + test-coverage thinness MEDIUM)

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (pipeline override: implementation + tests only, research not audited)
Documentation Accuracy: WARNING

## Blocking Issues
None.

## Non-Blocking Issues
1. [LOW, DOC_CODE_MISMATCH] /publish endpoint implemented (server.go:58-66) but undocumented in README.
2. [MEDIUM, MISSING_EDGE_CASE] ws ReadFrame trusts remote 64-bit length with no cap; no control-frame length / fragmentation validation (ws.go:84-127).
3. [LOW] Upgrade header case-sensitive, Sec-WebSocket-Version not validated (ws.go:31-76).
4. [LOW] WriteFrame client-mask uses fixed key, Dial ignores urlStr param and skips Accept verification (ws.go:129-202).
5. [MEDIUM, MISSING_TEST] No failure/negative tests: handshake rejection, ping/pong, close, extended-length, invalid Last-Event-ID, /publish, slow-subscriber drop.
6. [LOW] SSE history unbounded + silent drop to slow subscribers — disclosed in engineering notes, lab-scoped.

## Required Revisions
None blocking. Recommended before publication: document or remove /publish; add max-frame cap; add negative tests for handshake, close/ping, invalid Last-Event-ID.

## Final Status

APPROVED_WITH_WARNINGS
