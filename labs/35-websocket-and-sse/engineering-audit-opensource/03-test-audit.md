## Test Findings

- TestSSE_Formatting passes: validates Event.Format.
- TestSSE_LastEventID_Resumption passes: verifies replay skipping ID 1.
- TestWebSocket_TextAndBinary passes: checks echo text and reversed binary.
- TestSSE_Concurrency passes: spawns 10 clients, broadcasts 5 events, ensures no panic.

All tests succeed (`go test ./...` PASS).