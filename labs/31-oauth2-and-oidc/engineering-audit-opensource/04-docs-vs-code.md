# Docs vs Code – labs/31-oauth2-and-oidc

## Claim‑to‑Implementation Cross‑check

| Claim (README) | Location in Code | Verified |
|---|---|---|
| PKCE (S256) generation & verification | `pkg/pkce/pkce.go` (GeneratePKCEPair, ComputeChallenge, Verify) | ✅ |
| PKCE plain method support | same file (method "plain") | ✅ |
| OIDC ID Token signing (HS256) & claims validation | `pkg/oidc/oidc.go` (SignIDToken, ParseAndVerifyIDToken) | ✅ |
| ID Token claims: iss, sub, aud, exp, iat, nonce | same file, claim checks | ✅ |
| Authorization Code flow (single‑use) | `pkg/server/server.go` Authorize & ExchangeCode | ✅ |
| Mandatory code_challenge & method | same file (Authorize validation) | ✅ |
| Refresh Token rotation (new token per use) | `Refresh` method new token issuance | ✅ |
| Replay detection revokes token family | `Refresh`‑revoked check + `revokedFams` map | ✅ |
| Access token scope checking | `ValidateAccessToken` + `containsScope` | ✅ |
| Concurrency safety (race‑free) | Mutex protects all mutable maps; race test passes | ✅ |
| Demo steps (PKCE attack, token rotation, replay detection) | `cmd/demo/main.go` walks through each step | ✅ |

## Mismatches / Missing Docs

- README mentions "Refresh Token Rotation (RFC 9700 Section 4.14)" – code implements rotation and revocation but does not expose explicit RFC reference in comments. Not a functional issue.
- README claims "no external heavyweight dependencies" – true; only stdlib used.
- No documentation of `ErrIssuedInFuture` branch (ID token future iat) – not exercised in README examples.
- `Engineering Design` notes mention "Token family revocation on replay" – correctly reflected in code.

## Assessment

No DOC_CODE_MISMATCH. Minor omissions are informational only, not affecting correctness. All core claims backed by implementation.
