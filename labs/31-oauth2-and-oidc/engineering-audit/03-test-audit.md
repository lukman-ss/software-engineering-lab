# Test Audit

Target Lab: labs/31-oauth2-and-oidc

## Test Execution Results

### 1. `go test -v ./...`
```text
=== RUN   TestPKCE_S256_Valid
--- PASS: TestPKCE_S256_Valid (0.00s)
=== RUN   TestPKCE_InvalidMethod
--- PASS: TestPKCE_InvalidMethod (0.00s)
=== RUN   TestPKCE_Mismatch
--- PASS: TestPKCE_Mismatch (0.00s)
=== RUN   TestOIDC_IDToken_Valid
--- PASS: TestOIDC_IDToken_Valid (0.00s)
=== RUN   TestOIDC_IDToken_TamperedSignature
--- PASS: TestOIDC_IDToken_TamperedSignature (0.00s)
=== RUN   TestOIDC_IDToken_Expired
--- PASS: TestOIDC_IDToken_Expired (0.00s)
=== RUN   TestOIDC_IDToken_MismatchClaims
--- PASS: TestOIDC_IDToken_MismatchClaims (0.00s)
=== RUN   TestOAuth2_FullFlowAndPKCEInterception
--- PASS: TestOAuth2_FullFlowAndPKCEInterception (0.00s)
=== RUN   TestOAuth2_RefreshTokenRotation_AndReplayDetection
--- PASS: TestOAuth2_RefreshTokenRotation_AndReplayDetection (0.00s)
=== RUN   TestPKCE_Plain_Method
--- PASS: TestPKCE_Plain_Method (0.00s)
=== RUN   TestOIDC_MalformedJWT
--- PASS: TestOIDC_MalformedJWT (0.00s)
=== RUN   TestOAuth2_NegativePaths
--- PASS: TestOAuth2_NegativePaths (0.00s)
=== RUN   TestOAuth2_ConcurrencyAndRace
--- PASS: TestOAuth2_ConcurrencyAndRace (0.00s)
PASS
ok  	labs/31-oauth2-and-oidc/tests	0.231s
```

### 2. `go test -race ./...`
```text
ok  	labs/31-oauth2-and-oidc/tests	1.174s
```

### 3. `go run ./cmd/demo`
```text
=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===

[Step 1] Initiate Authorization Request with PKCE & OIDC scope...
 Generated PKCE code_challenge (S256): naR7QYtJk8EVqY7GgdcNIvb7gTFJKgIBfKuVjimFJDw
 Generated OIDC nonce: e1648403dc71dcbe7103bb915ec0c92d

[Step 2] Authorization Server issues Authorization Code...
 Issued Code: 014caa33a3c2217dff36544d1f5cab76 (expires: 16:11:04)

[Step 3] Client exchanges Code + code_verifier for Tokens...
 Access Token : at_0090009aabf2f3ce1547171537277f835e1e47cc47daa370
 Refresh Token: rt_eeffbe9fc5a8e2f0514594db3e84a22f88d9c713e9c89920
 ID Token     : eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...

[Step 4] Validate OIDC ID Token claims...
 ID Token Subject  : user_alice_99
 ID Token Issuer   : https://auth.example.com
 ID Token Audience : spa-client-123
 ID Token Nonce Match: SUCCESS (e1648403dc71dcbe7103bb915ec0c92d)

[Step 5] Access Resource Server using Access Token...
 Resource Server verified Access Token. Subject: user_alice_99

[Step 6] Test Interception Attack (Attempt code exchange with wrong verifier)...
 Interception Attack Blocked! Error: pkce verification failed: code_verifier does not match code_challenge

[Step 7] Test Refresh Token Rotation...
 Rotated New Access Token : at_378aecbbb6efd1e8ab0e5e746029b59b991a0ffeb97caac5
 Rotated New Refresh Token: rt_b51bb6806f28dbadc3a2c274971a3717674acc7922054a57

[Step 8] Test Refresh Token Replay Detection (Attacker reuses stolen Refresh Token)...
 Refresh Token Replay Detected! Error: refresh token reuse detected: entire token family revoked
 Re-testing active refresh token after family revocation...
 Active session revoked due to family revocation! Error: refresh token reuse detected: entire token family revoked

=== Demo Completed Successfully ===
```

## Coverage and Assertion Assessment

- **PKCE Protection**: Tests cover valid generation (`TestPKCE_S256_Valid`), invalid methods (`TestPKCE_InvalidMethod`), verifier/challenge mismatches (`TestPKCE_Mismatch`), and plain mode (`TestPKCE_Plain_Method`).
- **OIDC Claims & Signatures**: Tests cover legitimate claims validation (`TestOIDC_IDToken_Valid`), signature tampering (`TestOIDC_IDToken_TamperedSignature`), expired tokens (`TestOIDC_IDToken_Expired`), claim mismatches for issuer, audience, and nonce (`TestOIDC_IDToken_MismatchClaims`), and malformed payloads (`TestOIDC_MalformedJWT`).
- **OAuth Authorization Code Grant & Code Reuse**: Tests confirm interception failure on verifier mismatch and rejection of authorization code reuse (`TestOAuth2_FullFlowAndPKCEInterception`).
- **Refresh Token Rotation & Family Revocation**: Tests verify rotation mints a fresh token and replay of previous token revokes subsequent rotated tokens within the family (`TestOAuth2_RefreshTokenRotation_AndReplayDetection`).
- **Negative Paths**: Exhaustive negative test coverage for bad clients, invalid redirect URIs, missing challenges, invalid codes, and mismatched scopes (`TestOAuth2_NegativePaths`).
- **Concurrency & Race Detection**: 20 concurrent goroutines executing simultaneous authorize, exchange, validate, and refresh flows pass under `-race` with no data races detected (`TestOAuth2_ConcurrencyAndRace`).
