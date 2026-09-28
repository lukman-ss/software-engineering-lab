# Gaps Analysis — labs/35-websocket-and-sse

Allowed gap types list from PIPELINE OVERRIDE.

## MISSING_TEST
- WS close handshake (OpClose read/write)
- WS ping/pong handling test (server already writes pong on ping)
- WS extended-length (126/127) payload test (not implemented, but test would guard)
- WS fragmented frame (OpContinuation) test
- SSE keep-alive heartbeat test
- SSE invalid Last-Event-ID (non-integer, large int) test
- /publish endpoint test
- Subscriber drop test (broadcast silent-drop when channel full)
- WS client reuse after close (not in spec but could test)
- Negative handshake test (missing Upgrade key/version)

## BROKEN_IMPLEMENTATION
- None: core replay and WS text/binary echo are correct; all warnings are lab-scoped robustness.

## DOC_CODE_MISMATCH
- /publish endpoint (server.go line 58-66) never mentioned in README.md

## RACE_CONDITION
- None: -race detector passes. (Potential theoretical race on broadcast iteration + client close mitigated by mutex ordering; detector clean.)

## UNHANDLED_ERROR
- None: all returned errors either handled (client in demo) or discarded with rationale (server echo write failures acceptable).

## MISSING_EDGE_CASE
- WS server does not enforce max frame length (could OOM on malicious client payload of length 2^63-1 bytes).
- WS server does not reject control frames with payload > 125 bytes (per RFC 6455 §5.5).
- SSE Format with ID=0 and empty Data yields only `\n` (valid but surprising); not exercised by tests.
- ws.Dial ignores `urlStr` param (hardcodes `/ws`) – API confusion but not broken for lab usage.

## IMPLEMENTATION_OVERCLAIM
- README claims "full-duplex text & binary protocol" – true as demo shows bidirectional exchange, though the server only echoes, not a general broadcast. Not an overclaim.

## RESEARCH_MISMATCH
- Not audited per PIPELINE OVERRIDE.

## FAKE_DEMO
- Demo output verified: matches engineering/03-execution-result.md character-for-character modulo dynamic port.

## FAKE_BENCHMARK
- Not applicable: no benchmark claims.

## UNVERIFIED_RESULT
- None: all test outputs are deterministic assertions; demo shows live results.

Severity mapping: each MISSING_TEST is LOW (demo-lab). DOC_CODE_MISMATCH LOW. MISSING_EDGE_CASE MEDIUM (frame-length DoS vector though no external service exposure in lab). Remaining N/A.