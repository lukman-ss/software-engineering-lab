# Gaps Analysis

Target Lab: labs/35-websocket-and-sse

## Identified Gaps
1. MISSING_TEST (MEDIUM): No test for error paths – missing Sec-WebSocket-Key, invalid Upgrade header, malformed frames, large payloads, back-pressure drops.
2. UNHANDLED_ERROR (LOW): Broadcast default-drop silently loses events if client channel buffer full.
3. DOC_CODE_MISMATCH (LOW): `/publish` endpoint undocumented in README.
4. IMPLEMENTATION_OVERCLAIM (MEDIUM): Design states "WS detects connection drop and reconnects" – not present.
5. IMPLEMENTATION_OVERCLAIM (MEDIUM): Design states "WS broadcast" – code only echoes.
6. UNVERIFIED_RESULT (LOW): Demo output recorded but not programmatically validated in CI (no assertions on demo output).
7. RESEARCH_MISMATCH (MEDIUM): Design mentions "native standard library implementation" – code uses manual hijack, which is standard but deviates from typical net/http usage; acceptable but noted.

## Severity Summary
- HIGH: none
- MEDIUM: items 1,4,5,7
- LOW: items 2,3,6

## Recommendations
- Add negative-path tests for WS handshake and frame errors.
- Document `/publish` in README.
- Clarify design to match actual capabilities (no broadcast, no auto-reconnect).
- Consider back-pressure handling or log when events dropped.
- Add demo output verification (e.g., capture stdout, assert strings).
