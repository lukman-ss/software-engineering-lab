# Test Audit

Target Lab: labs/31-oauth2-and-oidc

## Commands Executed
- `go test -v ./...` → PASS (13/13, tests package)
- `go test -race ./...` → PASS (no races)
- `go run ./cmd/demo` → PASS (8 steps, completed successfully)

## Coverage
- Happy path: PASS (full auth-code + PKCE + ID token + access + refresh)
- Failure path: PASS (PKCE mismatch, tampered sig, expired, iss/aud/nonce mismatch, code reuse, replay)
- Edge cases: PARTIAL (missing expired auth-code, expired access/refresh, iat-future, verifier length)
- Transitions: PASS (code unused→used, refresh active→revoked→family-revoked)
- Recovery/Rollback: PASS (family revocation blocks active token)
- Concurrency: PASS (20 workers independent flows, race clean; no concurrent same-token replay test)
- Negative cases: PASS (bad client, bad redirect, empty challenge, bad method, bad code, scope check, bad RT)

## Strengths
- Attack simulations real: interception + replay + tamper all asserted
- Scope validation positive + negative asserted
- Error strings asserted for replay (`reuse detected`) and revocation (`revoked`)

## Weaknesses
- No test for `ErrIssuedInFuture` (oidc.go:107)
- No test for expired auth-code (`ExchangeCode` expired branch), expired access token, expired refresh token
- No test for `ErrInvalidVerifierLength`
- Concurrency test covers parallel independent flows only, not parallel replay of same RT

## Verdict
Suite passes and proves core claims, but edge expiry paths unproven.
