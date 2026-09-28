# Docs vs Code Audit

## Comparisons

### 1. README vs Code
- README outlines packages: `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, `tests`.
- Observed: All 5 directories and modules exist and match README structure exactly.
- Test commands in README (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) execute exactly as described.
- Assessment: PASS

### 2. Engineering Notes vs Code
- Notes describe zero third-party dependencies, standard Go library usage, HMAC-SHA256 signed double-submit token structure (`sessionID:timestamp:nonce:hmac`), and Fetch Metadata checks.
- Observed: Implemented precisely as specified.
- Assessment: PASS

### 3. Demo Output vs Test Assertions
- Both demo and integration tests assert:
  1. Transfer execution occurs on vulnerable endpoint even when CORS rejects origin.
  2. Anti-CSRF token middleware returns 403 on missing or invalid token.
  3. Legitimate client with valid token executes transfer.
- Assessment: PASS
