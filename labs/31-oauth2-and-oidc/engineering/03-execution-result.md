# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
Build succeeded without errors.
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
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
=== RUN   TestOAuth2_ConcurrencyAndRace
--- PASS: TestOAuth2_ConcurrencyAndRace (0.00s)
PASS
ok  	labs/31-oauth2-and-oidc/tests	0.338s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
ok  	labs/31-oauth2-and-oidc/tests	1.366s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===

[Step 1] Initiate Authorization Request with PKCE & OIDC scope...
 Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w
 Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de

[Step 2] Authorization Server issues Authorization Code...
 Issued Code: 5fabbdcdb6bbd582c3811481395a5161 (expires: 15:31:47)

[Step 3] Client exchanges Code + code_verifier for Tokens...
 Access Token : at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2
 Refresh Token: rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928
 ID Token     : eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...

[Step 4] Validate OIDC ID Token claims...
 ID Token Subject  : user_alice_99
 ID Token Issuer   : https://auth.example.com
 ID Token Audience : spa-client-123
 ID Token Nonce Match: SUCCESS (2152fa0c9254b322d11b27dba96b99de)

[Step 5] Access Resource Server using Access Token...
 Resource Server verified Access Token. Subject: user_alice_99

[Step 6] Test Interception Attack (Attempt code exchange with wrong verifier)...
 Interception Attack Blocked! Error: pkce verification failed: code_verifier does not match code_challenge

[Step 7] Test Refresh Token Rotation...
 Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd
 Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590

[Step 8] Test Refresh Token Replay Detection (Attacker reuses stolen Refresh Token)...
 Refresh Token Replay Detected! Error: refresh token reuse detected: entire token family revoked
 Re-testing active refresh token after family revocation...
 Active session revoked due to family revocation! Error: refresh token reuse detected: entire token family revoked

=== Demo Completed Successfully ===
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
