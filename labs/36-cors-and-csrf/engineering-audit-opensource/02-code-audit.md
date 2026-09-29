# Code Audit

## Finding 1

Location: `internal/cors/middleware.go:52-58`
Claimed Behavior: CORS middleware rejects disallowed origins on preflight with 403; simple cross-origin requests from disallowed origins continue to the handler (no `ACAO` header returned).
Observed Implementation: For disallowed origins on OPTIONS, returns 403. For disallowed origins on non-OPTIONS, calls `next.ServeHTTP(w, r)` without setting any CORS headers. This is spec-compliant: the server processes the request; it is the browser that enforces the disallowed response.
Assessment: PASS
Severity: LOW
Notes: This is intentional and is the core educational point of the lab — CORS does NOT stop server execution.

---

## Finding 2

Location: `internal/cors/middleware.go:62-72`
Claimed Behavior: When `AllowCredentials: true`, `Access-Control-Allow-Origin` must reflect the specific requesting origin rather than `*`.
Observed Implementation: When `AllowCredentials` is true, code always reflects `origin` value and sets `Allow-Credentials: true`. No wildcard `*` emitted. When `AllowCredentials` is false and single allowed origin is `*`, emits `*`. Otherwise reflects `origin`.
Assessment: PASS
Severity: LOW
Notes: Spec-correct per Fetch spec section 3.2.

---

## Finding 3

Location: `internal/csrf/token.go:72-74`
Claimed Behavior: Token parsing rejects malformed tokens.
Observed Implementation: Token format is `base64(sessionID:ts:nonce:signature)`. Split on `:` expects exactly 4 parts. FLAW: If `sessionID` itself contains `:`, the split will produce more than 4 parts and the token will be rejected with `ErrInvalidToken`. However, in the current bank app, session IDs are fixed strings like `"session-victim-secret"` with no colons, so this is not exploitable in the current setup, but it is a latent structural bug if session IDs with colons are used.
Assessment: WARNING
Severity: MEDIUM
Notes: `strings.Split(string(decoded), ":")` used for format that embeds arbitrary user-controlled session IDs. Correct approach is `SplitN(..., 4)` to preserve colons in session ID. Does not affect current test set but breaks correctness contract for IDs containing `:`.

---

## Finding 4

Location: `internal/csrf/token.go:105`
Claimed Behavior: Constant-time comparison to prevent timing attacks.
Observed Implementation: `subtle.ConstantTimeCompare` used correctly.
Assessment: PASS
Severity: LOW
Notes: Correct.

---

## Finding 5

Location: `internal/bank/app.go:149-151`
Claimed Behavior: `HandleTransferProtected` implements CSRF-validated transfer.
Observed Implementation: `HandleTransferProtected` simply delegates to `HandleTransferVulnerable`. CSRF protection is applied as HTTP middleware layer at the mux level in `tests/integration_test.go:42`. This design is correct — the protection is in the middleware, not the handler. However, the function comment implies the handler itself has CSRF verification, which is inaccurate.
Assessment: WARNING
Severity: LOW
Notes: Comment at line 150 says "CSRF middleware handles protection wrapper" — acceptable design. Not a bug.

---

## Finding 6

Location: `internal/bank/app.go:122-138`
Claimed Behavior: Account balance mutations are thread-safe.
Observed Implementation: `sync.RWMutex` used; balance deduction (`sender.Balance -= amount`) and credit (`recipient.Balance += amount`) both happen inside `b.mu.Lock()`. Correct.
Assessment: PASS
Severity: LOW
Notes: No TOCTOU issue. Lock held for full transaction.

---

## Finding 7

Location: `internal/bank/app.go:55-73` (authenticate)
Claimed Behavior: Session authentication reads from session map safely.
Observed Implementation: `b.mu.RLock()` used for `sessions` lookup. Then calls `b.GetAccount()` which re-acquires `b.mu.RLock()`. `sync.RWMutex` in Go allows multiple concurrent readers, so double RLock from same goroutine is safe (no deadlock for RLock).
Assessment: PASS
Severity: LOW
Notes: Correct usage.

---

## Finding 8

Location: `internal/csrf/middleware.go:27`
Claimed Behavior: Safe/idempotent methods (GET, HEAD, OPTIONS, TRACE) bypass CSRF check.
Observed Implementation: Check passes GET, HEAD, OPTIONS, TRACE directly to next handler. TRACE inclusion is pedantically correct per RFC 7231 idempotent method list.
Assessment: PASS
Severity: LOW
Notes: No issues.

---

## Finding 9

Location: `cmd/demo/main.go` (not read yet — inferred from test output)
Claimed Behavior: Demo shows: (1) vulnerable endpoint executes despite CORS block; (2) protected endpoint rejects attack; (3) legit flow with CSRF token succeeds.
Observed Implementation: Execution output confirms all three scenarios matched. Balances computed correctly: victim $1000 → $600 (attack 1: $400 stolen), then $600 → $500 (legit: $100 transferred). Attacker $50 → $450 → $550.
Assessment: PASS
Severity: LOW
Notes: Demo output is real and consistent with implementation behavior.

---

## Finding 10

Location: `internal/csrf/middleware.go:61`
Claimed Behavior: `FetchMetadataMiddleware` rejects cross-site non-safe requests.
Observed Implementation: Only rejects if `Sec-Fetch-Site == "cross-site"`. Requests with absent `Sec-Fetch-Site` header (e.g., from non-browser clients, old browsers) pass through unchecked. This is by design as per the Fetch Metadata spec (old browsers that don't send the header must not be blocked).
Assessment: PASS
Severity: LOW
Notes: Correct behavior per spec. Absence of header = allow.

---

## Finding 11

Location: Concurrency test (`tests/integration_test.go:189-207`)
Claimed Behavior: Concurrent token generation is race-free.
Observed Implementation: 20 goroutines call `/api/csrf-token` simultaneously. `go test -race` passed. `GenerateToken` uses `rand.Read` which is goroutine-safe in Go stdlib (since Go 1.20). `BankServer` read path for sessions uses `RLock`.
Assessment: PASS
Severity: LOW
Notes: Race detector passed.
