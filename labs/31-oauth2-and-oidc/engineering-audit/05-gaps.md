# Gap Analysis

## Gaps Identified

### Gap 1: Non-Constant-Time String Comparison in PKCE Verification
- **Gap Type**: `MISSING_EDGE_CASE`
- **Severity**: LOW
- **Description**: `pkg/pkce/pkce.go:68` uses string equality `expected != challenge` instead of `crypto/subtle.ConstantTimeCompare`.
- **Impact**: In practice, PKCE code challenge is held server-side, so side-channel timing attacks are impractical, but constant-time comparison is recommended security hygiene.
- **Remediation**: Use `subtle.ConstantTimeCompare([]byte(expected), []byte(challenge)) == 1`.

### Gap 2: Unbounded In-Memory Token Storage
- **Gap Type**: `MISSING_EDGE_CASE`
- **Severity**: LOW
- **Description**: `AuthorizationServer` maps (`authCodes`, `tokens`, `refreshMeta`, `revokedFams`) grow indefinitely without TTL eviction or memory caps.
- **Impact**: Lab memory usage increases over time under long runtime.
- **Remediation**: Add a background cleanup routine or purge expired entries during operations.

### Gap 3: Client State Parameter Unused
- **Gap Type**: `DOC_CODE_MISMATCH`
- **Severity**: LOW
- **Description**: `Client.BuildAuthorizationRequest` populates `c.State`, but `c.State` is never sent to `Authorize` or validated during exchange.
- **Impact**: CSRF state parameter is unused. Does not affect PKCE, OIDC, or Refresh Rotation scope.
- **Remediation**: Either include `State` validation in Client flow or add comment clarifying in-memory demo omission.

---

## Allowed Gap Types Check

| Allowed Gap Type | Found? | Details |
|---|---|---|
| MISSING_TEST | No | Test coverage is thorough (17 unit/concurrency tests) |
| BROKEN_IMPLEMENTATION | No | All flows execute correctly as designed |
| DOC_CODE_MISMATCH | Yes (LOW) | Unused client `State` field |
| RACE_CONDITION | No | `go test -race ./...` passed cleanly |
| UNHANDLED_ERROR | No | Error propagation is complete |
| MISSING_EDGE_CASE | Yes (LOW) | Non-constant-time PKCE compare; unbounded map growth |
| IMPLEMENTATION_OVERCLAIM | No | All claims backed by code and tests |
| RESEARCH_MISMATCH | No | Implementation aligns with research inputs |
| FAKE_DEMO | No | Demo runs live and produces real output |
| FAKE_BENCHMARK | No | No benchmarks claimed |
| UNVERIFIED_RESULT | No | All execution output verified live |
