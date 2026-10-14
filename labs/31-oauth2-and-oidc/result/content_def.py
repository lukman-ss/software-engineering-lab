#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 31: OAuth 2.0 & OpenID Connect (OIDC).
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
"""
import os
import sys
import math
import subprocess
from datetime import datetime
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN
from pptx.enum.shapes import MSO_AUTO_SHAPE_TYPE, MSO_CONNECTOR

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/31-oauth2-and-oidc/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc"
COMMIT_SHA = "cf0cd90745a547032e8fb3a2cb838f40db1de30a"
LAB_NO = "31"
SLUG = "oauth2-and-oidc"
TITLE_ID = "OAuth 2.0 & OpenID Connect (OIDC)"
SUBTITLE_ID = "Otorisasi, Autentikasi, dan Perlindungan Token Modern"
TITLE_EN = "OAuth 2.0 & OpenID Connect (OIDC)"
SUBTITLE_EN = "Authorization, Authentication, and Modern Token Protection"
AUTHOR = "Lukman / lukman-ss"

# Color constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    """Add a slide to the presentation."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })


# =====================================================================
# SLIDE CONTENT DEFINITIONS
# =====================================================================

# Slide 1: Cover (Dark Navy, visual element)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "CRT terminal retro beige dengan scanlines",
            "ASCII art visualisasi JWT dan PKCE",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Motivation & Problem Definition
add_slide(
    categories=["MOTIVASI", "KONSEP"],
    title="Auth vs Authz: Kebingungan yang Berbahaya",
    subtitle="Access Token bukan Identity — mixing keduanya membuka celah keamanan",
    content_items=[
        {"title": "MOTIVASI", "items": [
            "OAuth 2.0 sering disalahpahami sebagai protokol autentikasi",
            "Access Token menjawab 'apa yang boleh diakses', bukan 'siapa penggunanya'",
            "OIDC menambahkan layer identitas di atas OAuth 2.0 untuk menjawab 'siapa user'",
        ], "type": "bullet"},
        {"title": "RISIKO KONFUSI", "items": [
            "Menyimpan bearer token di localStorage → XSS mengeksploitasi",
            "Menggunakan implicit flow → token tereksposisi di URL fragment",
            "Refresh token statis → attacker mempertahankan akses selamanya",
        ], "type": "bullet"},
        {"title": "STANDAR TERKINI", "items": [
            "RFC 9700: Wajib Authorization Code Flow + PKCE untuk browser client",
            "OIDC Core 1.0: ID Token JWT dengan signature & claims wajib",
            "Implicit flow dihapus dari OAuth 2.1",
        ], "type": "bullet"},
    ],
    mental_model="OAuth 2.0 is for authorization (what you can access), OIDC is for authentication (who you are).",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="2 / 24"
)

# Slide 3: Core Concepts & Architecture
add_slide(
    categories=["KONSEP INTI", "ARKITEKTUR"],
    title="Empat Mekanisme Inti",
    subtitle="PKCE, Authorization Code Grant, ID Token Validation, Refresh Token Rotation",
    content_items=[
        {"title": "PKCE (RFC 7636)", "items": [
            "Client generate code_verifier (32 byte random)",
            "code_challenge = SHA-256(code_verifier)",
            "Mencegah authorization code interception",
        ], "type": "bullet"},
        {"title": "Authorization Code Grant", "items": [
            "Client → AS: response_type=code, code_challenge, nonce",
            "AS → Client: authorization code (one-time, expires 5 min)",
            "Client → AS: exchange code + code_verifier for tokens",
        ], "type": "bullet"},
        {"title": "ID Token (OIDC)", "items": [
            "JWT signed dengan HMAC-SHA256",
            "Claims wajib: iss, sub, aud, exp, iat, nonce",
            "Client verifikasi signature, issuer, audience, expiry, nonce",
        ], "type": "bullet"},
        {"title": "Refresh Token Rotation (RFC 9700)", "items": [
            "Setiap refresh menghasilkan token baru",
            "Token lama di-revoke",
            "Jika token lama digunakan ulang → entire token family revoked",
        ], "type": "bullet"},
    ],
    mental_model="Four pillars: PKCE for code protection, Code Flow for security, ID Token for identity, Rotation for replay defense.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="3 / 24"
)

# Slide 4: Authorization Request Flow
add_slide(
    categories=["ALUR AUTHORIZATION"],
    title="Authorization Request dengan PKCE",
    subtitle="Client mengirim challenge, nonce, state untuk keamanan",
    content_items=[
        {"title": "PARAMETER REQUEST", "items": [
            "response_type=code",
            "client_id=spa-client-123",
            "redirect_uri=https://app.example.com/callback",
            "scope=openid profile email",
        ], "type": "bullet"},
        {"title": "KEAMANAN", "items": [
            "state=<random> → anti-CSRF",
            "code_challenge=<challenge> → PKCE S256",
            "code_challenge_method=S256",
            "nonce=<random> → anti-replay ID Token",
        ], "type": "bullet"},
        {"title": "ALUR", "items": [
            "1. Client generates PKCE pair + nonce + state",
            "2. Client redirects to AS with all parameters",
            "3. AS stores challenge & nonce bound to code",
            "4. AS redirects back with authorization code",
        ], "type": "bullet"},
    ],
    mental_model="State, PKCE challenge, and nonce are bound to the code — interception without verifier fails.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="4 / 24"
)

# Slide 5: Token Exchange & PKCE Verification
add_slide(
    categories=["BUKTI LOGIKA", "TOKEN EXCHANGE"],
    title="Token Exchange: Verifikasi PKCE",
    subtitle="Server hitung SHA-256 verifier, bandingkan dengan stored challenge",
    content_items=[
        {"title": "EXCHANGE REQUEST", "items": [
            "grant_type=authorization_code",
            "code=<authorization_code>",
            "redirect_uri=https://app.example.com/callback",
            "client_id=spa-client-123",
            "code_verifier=dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w",
        ], "type": "bullet"},
        {"title": "VERIFIKASI SERVER", "items": [
            "1. Check code is unused & not expired",
            "2. Compute SHA-256(code_verifier)",
            "3. Compare with stored code_challenge",
            "4. If match → issue tokens, mark code used",
        ], "type": "bullet"},
        {"title": "INTERCEPTION ATTACK", "items": [
            "Attacker intercepts authorization code",
            "Attacker does NOT have code_verifier",
            "Exchange with wrong verifier → verification fails",
            "Error: pkce verification failed: code_verifier does not match code_challenge",
        ], "type": "bullet"},
    ],
    mental_model="Only the legitimate client has the original verifier — attackers can't exchange the intercepted code.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="5 / 24"
)

# Slide 6: ID Token Structure & Claims
add_slide(
    categories=["BUKTI KODE", "ID TOKEN"],
    title="ID Token: JWT Structure",
    subtitle="Header + Claims + Signature HMAC-SHA256",
    content_items=[
        {"title": "JWT FORMAT", "items": [
            "header.claims.signature (3 base64url parts)",
            "Header: {\"alg\":\"HS256\",\"typ\":\"JWT\"}",
            "Claims: {\"iss\":\"...\",\"sub\":\"...\",\"aud\":\"...\",\"exp\":...,\"iat\":...,\"nonce\":\"...\"}",
        ], "type": "bullet"},
        {"title": "SIGNING", "items": [
            "1. Encode header & claims to base64url",
            "2. unsignedToken = header + '.' + claims",
            "3. signature = HMAC-SHA256(unsignedToken, shared_secret)",
            "4. final = unsignedToken + '.' + signature",
        ], "type": "bullet"},
        {"title": "VALIDATION", "items": [
            "1. Decode & verify HMAC-SHA256 signature (constant-time)",
            "2. iss exact match → https://auth.example.com",
            "3. aud match → spa-client-123",
            "4. exp > now, iat < now+300s",
            "5. nonce match → anti-replay",
        ], "type": "bullet"},
    ],
    mental_model="JWT is three base64url segments — verify signature with HMAC-SHA256, then validate all claims.",
    source_path="labs/31-oauth2-and-oidc/content/03-code-snippets.md",
    slide_num="6 / 24"
)

# Slide 7: ID Token Validation Code
add_slide(
    categories=["BUKTI KODE"],
    title="Validation Implementation",
    subtitle="Six-step verification process",
    content_items=[
        {"title": "FORMAT CHECK", "items": [
            "parts = rawJWT.split('.')",
            "if len(parts) != 3: return ErrMalformedJWT",
        ], "type": "code"},
        {"title": "SIGNATURE VERIFY", "items": [
            "unsigned = parts[0] + '.' + parts[1]",
            "expected_sig = HMAC-SHA256(unsigned, secret)",
            "if not hmac.Equal(sig, expected_sig): return ErrSignatureInvalid",
        ], "type": "code"},
        {"title": "CLAIMS VALIDATE", "items": [
            "if claims.Issuer != expectedIssuer: return ErrIssuerMismatch",
            "if claims.Audience != expectedAudience: return ErrAudienceMismatch",
            "if claims.Expiration <= now: return ErrTokenExpired",
            "if claims.IssuedAt > now+300: return ErrIssuedInFuture",
            "if claims.Nonce != expectedNonce: return ErrNonceMismatch",
        ], "type": "code"},
    ],
    mental_model="Validate in order: format, signature (constant-time), then all claims — reject on first failure.",
    source_path="labs/31-oauth2-and-oidc/pkg/oidc/oidc.go",
    slide_num="7 / 24"
)

# Slide 8: Token Response
add_slide(
    categories=["RESPONSE", "TOKENS"],
    title="Token Response",
    subtitle="Access Token, ID Token, Refresh Token",
    content_items=[
        {"title": "ACCESS TOKEN", "items": [
            "Bearer token opaque string (e.g., at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2)",
            "Lifetime: 1 hour",
            "Scope-limited, used to access protected resources",
            "NOT for identity assertion",
        ], "type": "bullet"},
        {"title": "ID TOKEN", "items": [
            "Signed JWT with claims (iss, sub, aud, exp, iat, nonce)",
            "Lifetime: 1 hour",
            "Client verifies signature, issuer, audience, expiry, nonce",
            "Used to assert user identity",
        ], "type": "bullet"},
        {"title": "REFRESH TOKEN", "items": [
            "Long-lived token (e.g., 30 days)",
            "Used to obtain new access tokens",
            "Single-use: each refresh invalidates old token",
            "Family tracking: replay revokes entire family",
        ], "type": "bullet"},
    ],
    mental_model="Access token for resources, ID token for identity, refresh token for rotation — never use access token as identity.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="8 / 24"
)

# Slide 9: Protected Resource Access
add_slide(
    categories=["BUKTI LOGIKA"],
    title="Resource Server Verification",
    subtitle="Resource server accepts Access Token, checks scope",
    content_items=[
        {"title": "REQUEST FORMAT", "items": [
            "GET /api/profile HTTP/1.1",
            "Host: api.example.com",
            "Authorization: Bearer at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2",
        ], "type": "code"},
        {"title": "SERVER PROCESS", "items": [
            "1. Extract bearer token from Authorization header",
            "2. Validate token exists in AS token store",
            "3. Check token scope matches requested resource",
            "4. Return protected resource with subject from token",
        ], "type": "bullet"},
        {"title": "VERIFICATION OUTPUT", "items": [
            "Resource Server verified Access Token. Subject: user_alice_99",
            "Scope check: profile → allowed",
            "Response: {\"id\":99,\"name\":\"Alice\",\"email\":\"alice@example.com\"}",
        ], "type": "bullet"},
    ],
    mental_model="Resource server validates token existence & scope — never trust the token as identity without AS validation.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="9 / 24"
)

# Slide 10: Refresh Token Rotation Flow
add_slide(
    categories=["ALUR ROTASI"],
    title="Refresh Token Rotation",
    subtitle="Single-use tokens with family lineage tracking",
    content_items=[
        {"title": "ROTATION PROCESS", "items": [
            "1. Client sends current refresh_token",
            "2. AS marks token.Revoked = true",
            "3. AS generates new refresh_token with same FamilyID",
            "4. AS generates new access_token",
            "5. AS returns both tokens",
        ], "type": "bullet"},
        {"title": "ROTATION OUTPUT", "items": [
            "Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd",
            "Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590",
            "FamilyID preserved: rt_c624... and rt_a0aa... belong to same family",
        ], "type": "bullet"},
        {"title": "REPLAY DETECTION", "items": [
            "Attacker reuses old refresh_token (already revoked)",
            "AS checks: meta.Revoked = true",
            "AS sets s.revokedFams[FamilyID] = true",
            "AS returns: ErrTokenReplayDetected",
        ], "type": "bullet"},
    ],
    mental_model="Rotate tokens on every use — replay of revoked token revokes entire family.",
    source_path="labs/31-oauth2-and-oidc/content/03-code-snippets.md",
    slide_num="10 / 24"
)

# Slide 11: Refresh Token Rotation Code
add_slide(
    categories=["BUKTI KODE"],
    title="Rotation Implementation",
    subtitle="Revocation + new token generation + family tracking",
    content_items=[
        {"title": "CHECK REVOCATION", "items": [
            "meta, exists := s.refreshMeta[refreshToken]",
            "if !exists: return ErrRefreshTokenNotFound",
            "if s.revokedFams[meta.FamilyID]: return ErrTokenReplayDetected",
            "if meta.Revoked: s.revokedFams[meta.FamilyID] = true; return ErrTokenReplayDetected",
        ], "type": "code"},
        {"title": "MARK OLD REVOKED", "items": [
            "meta.Revoked = true",
            "meta.UsedAt = time.Now()",
        ], "type": "code"},
        {"title": "GENERATE NEW", "items": [
            "newRtBytes := make([]byte, 24)",
            "rand.Read(newRtBytes)",
            "newRefreshToken := 'rt_' + hex.EncodeToString(newRtBytes)",
            "s.refreshMeta[newRefreshToken] = &RefreshTokenMeta{FamilyID: meta.FamilyID, Revoked: false}",
        ], "type": "code"},
    ],
    mental_model="Check revocation twice: before (family already compromised) and after (token reused).",
    source_path="labs/31-oauth2-and-oidc/pkg/server/server.go",
    slide_num="11 / 24"
)

# Slide 12: Failure Scenario 1: Code Interception
add_slide(
    categories=["KEGAGALAN", "PKCE INTERCEPTION"],
    title="Scenario 1: Authorization Code Interception",
    subtitle="Attacker intercepts code, lacks verifier",
    content_items=[
        {"title": "INTERCEPTION", "items": [
            "Attacker intercepts authorization code from redirect URL",
            "Attacker attempts token exchange with intercepted code",
        ], "type": "bullet"},
        {"title": "EXCHANGE ATTEMPT", "items": [
            "Client code + attacker's random verifier",
            "AS computes SHA-256(attacker_verifier)",
            "Stored challenge: SHA-256(original_verifier)",
            "Mismatch: challenge does not match",
        ], "type": "bullet"},
        {"title": "ERROR RESPONSE", "items": [
            "Error: pkce verification failed: code_verifier does not match code_challenge",
            "Token exchange rejected",
            "Attacker cannot obtain tokens",
            "Legitimate client with correct verifier succeeds",
        ], "type": "bullet"},
    ],
    mental_model="Interception without verifier = failure — PKCE binds verifier to code.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Failure Scenario 2: Forged ID Token
add_slide(
    categories=["KEGAGALAN", "TOKEN TAMPERING"],
    title="Scenario 2: Forged / Tampered ID Token",
    subtitle="Attacker modifies payload or signature",
    content_items=[
        {"title": "MODIFICATION ATTEMPT", "items": [
            "Attacker intercepts ID Token",
            "Attacker modifies payload (e.g., change sub from user_alice to admin)",
            "Attacker generates new signature with known key",
        ], "type": "bullet"},
        {"title": "VERIFICATION FAILURE", "items": [
            "Client verifies HMAC-SHA256 signature",
            "Computed signature ≠ embedded signature",
            "Return: ErrSignatureInvalid",
            "Token rejected before claims validation",
        ], "type": "bullet"},
        {"title": "ALTERNATIVE TAMPERING", "items": [
            "Attacker forges entirely new JWT without signing key",
            "HMAC-SHA256 verification fails immediately",
            "Return: ErrSignatureInvalid",
            "No claims validation needed — signature invalid",
        ], "type": "bullet"},
    ],
    mental_model="Forged tokens fail signature verification — no key, no valid token.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="13 / 24"
)

# Slide 14: Failure Scenario 3: Replay Attack
add_slide(
    categories=["KEGAGALAN", "TOKEN REPLAY"],
    title="Scenario 3: Refresh Token Replay",
    subtitle="Attacker reuses already-consumed refresh token",
    content_items=[
        {"title": "REPLAY ATTEMPT", "items": [
            "Attacker steals refresh token (rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928)",
            "Attacker uses stolen token after legitimate refresh",
            "Client refreshes token → new refresh token issued, old revoked",
        ], "type": "bullet"},
        {"title": "REPLAY DETECTION", "items": [
            "Attacker sends revoked refresh token to AS",
            "AS checks: meta.Revoked = true",
            "AS sets s.revokedFams[meta.FamilyID] = true",
            "AS returns: ErrTokenReplayDetected",
        ], "type": "bullet"},
        {"title": "FAMILY REVOCATION", "items": [
            "All active tokens in family immediately revoked",
            "Active session: Error: refresh token reuse detected: entire token family revoked",
            "Client must re-authenticate (force login)",
            "Attack mitigated — attacker loses access",
        ], "type": "bullet"},
    ],
    mental_model="Single-use refresh tokens — replay of revoked token revokes entire session family.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="14 / 24"
)

# Slide 15: Failure Scenario 4: Token Storage Vulnerability
add_slide(
    categories=["KEGAGALAN", "STORAGE"],
    title="Scenario 4: Insecure Token Storage",
    subtitle="localStorage XSS vulnerability",
    content_items=[
        {"title": "VULNERABLE PRACTICE", "items": [
            "Store access token in localStorage",
            "Application vulnerable to XSS attack",
            "Malicious script reads localStorage",
            "Attacker steals access token",
        ], "type": "bullet"},
        {"title": "EXPLOITATION", "items": [
            "XSS payload: fetch('/api/steal', {method:'POST', body: localStorage.getItem('access_token')})",
            "Access token sent to attacker server",
            "Attacker uses token to access protected resources",
            "No PKCE or refresh token protection — token is bearer",
        ], "type": "bullet"},
        {"title": "SECURE ALTERNATIVES", "items": [
            "httpOnly Secure SameSite cookie (BFF pattern)",
            "In-memory storage (react state, service variable)",
            "DPoP or mTLS sender-constrained tokens",
            "Never localStorage for bearer tokens",
        ], "type": "bullet"},
    ],
    mental_model="Bearer tokens in localStorage = XSS exposes token — use httpOnly cookies or in-memory.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="15 / 24"
)

# Slide 16: Architecture Components
add_slide(
    categories=["ARKITEKTUR", "KOMPONEN"],
    title="Component Structure",
    subtitle="Five packages + demo + tests",
    content_items=[
        {"title": "CORE PACKAGES", "items": [
            "pkg/pkce/pkce.go: PKCE code verifier + S256 challenge generation/verification",
            "pkg/oidc/oidc.go: ID Token JWT signing + verification + claims validation",
            "pkg/server/server.go: In-memory Authorization Server + rotation + replay detection",
            "pkg/client/client.go: Client helper for orchestration",
        ], "type": "bullet"},
        {"title": "IMPLEMENTATION", "items": [
            "cmd/demo/main.go: End-to-end walkthrough demo",
            "tests/oauth_test.go: 17 unit tests + race detector",
            "go.mod: Go 1.22+ (standard library only)",
            "No external dependencies",
        ], "type": "bullet"},
        {"title": "STANDARD LIBRARY", "items": [
            "crypto/rand: 256-bit entropy for verifiers",
            "crypto/hmac, crypto/sha256: HMAC-SHA256 JWT signature",
            "encoding/json: JWT claims marshaling",
            "encoding/base64: base64url encoding",
        ], "type": "bullet"},
    ],
    mental_model="Standard library only — crypto/rand, hmac/sha256, json, base64 for zero-dependency JWT.",
    source_path="labs/31-oauth2-and-oidc/engineering/01-design.md",
    slide_num="16 / 24"
)

# Slide 17: Test Coverage
add_slide(
    categories=["VERIFIKASI TEST"],
    title="Test Suite Coverage",
    subtitle="17 test cases with race detector",
    content_items=[
        {"title": "PKCE TESTS", "items": [
            "TestPKCE_S256_Valid: verifier 43-128 chars, challenge valid",
            "TestPKCE_InvalidMethod: reject non-S256/plain",
            "TestPKCE_Mismatch: wrong verifier rejected",
            "TestPKCE_VerifierLength_Bounds: reject <43 or >128",
        ], "type": "bullet"},
        {"title": "ID TOKEN TESTS", "items": [
            "TestOIDC_IDToken_Valid: signature valid, claims match",
            "TestOIDC_IDToken_TamperedSignature: modified signature rejected",
            "TestOIDC_IDToken_Expired: exp in past rejected",
            "TestOIDC_IDToken_MismatchClaims: iss/aud/nonce mismatch rejected",
            "TestOIDC_IDToken_IssuedInFuture: iat > now+300 rejected",
            "TestOIDC_MalformedJWT: wrong structure rejected",
        ], "type": "bullet"},
        {"title": "OAUTH2 TESTS", "items": [
            "TestOAuth2_FullFlowAndPKCEInterception: valid flow, interception blocked",
            "TestOAuth2_RefreshTokenRotation_AndReplayDetection: rotation works, replay revoked",
            "TestOAuth2_NegativePaths: unregistered client, bad redirect, invalid code",
            "TestOAuth2_ExpiredAuthCode_And_ExpiredTokens: expired resources rejected",
            "TestOAuth2_ConcurrentRefreshReplay: 10 goroutines, 9+ fail",
            "TestOAuth2_ConcurrencyAndRace: 20 workers, no race detected",
        ], "type": "bullet"},
    ],
    mental_model="17 tests cover PKCE, ID Token validation, full flow, and concurrency — race detector passed.",
    source_path="labs/31-oauth2-and-oidc/tests/oauth_test.go",
    slide_num="17 / 24"
)

# Slide 18: Test Results
add_slide(
    categories=["HASIL PENGUJIAN"],
    title="Test Execution Results",
    subtitle="PASS: 17/17, Race: PASS",
    content_items=[
        {"title": "TEST OUTPUT", "items": [
            "=== RUN   TestPKCE_S256_Valid",
            "--- PASS: TestPKCE_S256_Valid (0.00s)",
            "=== RUN   TestOIDC_IDToken_Valid",
            "--- PASS: TestOIDC_IDToken_Valid (0.00s)",
            "=== RUN   TestOAuth2_FullFlowAndPKCEInterception",
            "--- PASS: TestOAuth2_FullFlowAndPKCEInterception (0.00s)",
            "=== RUN   TestOAuth2_RefreshTokenRotation_AndReplayDetection",
            "--- PASS: TestOAuth2_RefreshTokenRotation_AndReplayDetection (0.00s)",
            "PASS — 17/17 tests",
        ], "type": "code"},
        {"title": "RACE DETECTOR", "items": [
            "Command: go test -race ./...",
            "Result: ok  	labs/31-oauth2-and-oidc/tests	1.366s",
            "No data races detected",
            "Thread-safe operations verified",
        ], "type": "bullet"},
        {"title": "DEMO OUTPUT", "items": [
            "=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===",
            "[Step 6] Interception Attack Blocked! Error: pkce verification failed",
            "[Step 8] Refresh Token Replay Detected! Error: entire token family revoked",
            "=== Demo Completed Successfully ===",
        ], "type": "code"},
    ],
    mental_model="17/17 tests pass, race detector clean, demo confirms attack scenarios fail as expected.",
    source_path="labs/31-oauth2-and-oidc/engineering/03-execution-result.md",
    slide_num="18 / 24"
)

# Slide 19: Checklist
add_slide(
    categories=["PANDUAN OPERASIONAL"],
    title="Security Checklist",
    subtitle="9 checkpoints before production",
    content_items=[
        {"title": "CHECKLIST", "items": [
            "[ ] PKCE S256 wajib untuk semua browser client",
            "[ ] Implicit flow & resource owner password flow dihindari",
            "[ ] ID Token diverifikasi: signature, iss, aud, exp, iat, nonce",
            "[ ] Authorization code one-time use, expires 5 minutes",
            "[ ] Refresh token rotation setiap penggunaan",
            "[ ] Replay detection → family revocation",
            "[ ] Access token tidak digunakan untuk authentication",
            "[ ] Token storage: httpOnly cookie atau in-memory",
            "[ ] Concurrency thread-safe (mutex, race detector)",
        ], "type": "bullet"},
        {"title": "STANDARD", "items": [
            "RFC 6749: OAuth 2.0 Authorization Framework",
            "RFC 7636: Proof Key for Code Exchange",
            "RFC 7519: JSON Web Token",
            "RFC 8725: OAuth 2.0 Security Best Current Practice",
            "RFC 9700: OAuth 2.0 Security Best Current Practice (2025)",
            "OIDC Core 1.0: Authentication layer",
        ], "type": "bullet"},
        {"title": "RECOMMENDED", "items": [
            "Use JWKS endpoint for asymmetric signing (RS256/ES256)",
            "Implement constant-time comparison for PKCE verification",
            "Use DPoP or mTLS for sender-constrained tokens",
            "Graceful eviction for expired codes/tokens",
            "Validate state parameter for anti-CSRF",
        ], "type": "bullet"},
    ],
    mental_model="Nine checks cover PKCE, ID Token validation, rotation, and storage — align with RFC 9700.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="19 / 24"
)

# Slide 20: Common Mistakes
add_slide(
    categories=["KESALAHAN UMUM"],
    title="6 Common Mistakes",
    subtitle="Penghindaran praktik berbahaya",
    content_items=[
        {"title": "MISTAKE 1", "items": [
            "Using access token for authentication",
            "Access token is NOT identity — it's authorization grant",
            "Solution: Use ID Token for user identity",
        ], "type": "bullet"},
        {"title": "MISTAKE 2", "items": [
            "Not validating ID Token signature",
            "Anyone can forge JWT without signature verification",
            "Solution: Verify HMAC-SHA256 signature before claims",
        ], "type": "bullet"},
        {"title": "MISTAKE 3", "items": [
            "Ignoring issuer/audience checks",
            "ID Token could be intended for another client",
            "Solution: Exact match issuer, audience exact or array match",
        ], "type": "bullet"},
        {"title": "MISTAKE 4", "items": [
            "Storing bearer token in localStorage",
            "XSS attacks can read localStorage and steal token",
            "Solution: httpOnly Secure SameSite cookie or in-memory",
        ], "type": "bullet"},
        {"title": "MISTAKE 5", "items": [
            "Using implicit flow",
            "Token exposed in URL fragment — browser history, logs, referrer",
            "Solution: Authorization Code Flow + PKCE (RFC 9700)",
        ], "type": "bullet"},
        {"title": "MISTAKE 6", "items": [
            "Static refresh tokens",
            "Same token reused indefinitely allows permanent access",
            "Solution: Single-use rotation with family tracking",
        ], "type": "bullet"},
    ],
    mental_model="Six mistakes: token misuse, no signature verification, no claim validation, insecure storage, implicit flow, static refresh.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="20 / 24"
)

# Slide 21: Case Study: Attack Scenarios
add_slide(
    categories=["STUDI KASUS"],
    title="Case Study: Attack Flow",
    subtitle="From code interception to family revocation",
    content_items=[
        {"title": "STEP 1-3: LEGITIMATE FLOW", "items": [
            "Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w",
            "Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de",
            "Access Token: at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2",
        ], "type": "bullet"},
        {"title": "STEP 4-5: ID TOKEN VALIDATION", "items": [
            "ID Token Subject: user_alice_99",
            "ID Token Issuer: https://auth.example.com",
            "ID Token Audience: spa-client-123",
            "ID Token Nonce Match: SUCCESS",
        ], "type": "bullet"},
        {"title": "STEP 6-8: ATTACK BLOCKED", "items": [
            "Interception Attack Blocked! Error: pkce verification failed",
            "Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590",
            "Refresh Token Replay Detected! Error: entire token family revoked",
            "Active session revoked due to family revocation! Error: refresh token reuse detected",
        ], "type": "bullet"},
    ],
    mental_model="Legitimate flow succeeds, interception blocked by PKCE, replay detected by rotation — defense in depth.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="21 / 24"
)

# Slide 22: Lab Boundaries
add_slide(
    categories=["BATAS LAB"],
    title="Lab Implementation Boundaries",
    subtitle="Simulasi konseptual — bukan produksi database asli",
    content_items=[
        {"title": "BATAS 1: IN-MEMORY MOCK", "items": [
            "Store menggunakan in-memory map (not real database)",
            "No network latency or real HTTP transport",
            "Tokens represented as strings, not cryptographically secure random",
        ], "type": "bullet"},
        {"title": "BATAS 2: SYMMETRIC SIGNING", "items": [
            "HMAC-SHA256 symmetric signing for demo simplicity",
            "Production should use asymmetric (RS256/ES256) with JWKS",
            "No key rotation mechanism demonstrated",
        ], "type": "bullet"},
        {"title": "BATAS 3: STATE VALIDATION", "items": [
            "Demo does not validate state parameter (client-side only)",
            "Production must validate state for anti-CSRF",
            "No PKCE constant-time comparison (audit: LOW severity)",
        ], "type": "bullet"},
    ],
    mental_model="This lab proves concepts — production needs asymmetric signing, JWKS, state validation, constant-time compare.",
    source_path="labs/31-oauth2-and-oidc/content/01-content-brief.md",
    slide_num="22 / 24"
)

# Slide 23: Core Design Principles
add_slide(
    categories=["PRINSIP INTI"],
    title="Three Core Principles",
    subtitle="Pillars for secure OAuth 2.0 & OIDC",
    content_items=[
        {"title": "PRINCIPLE 1: SEPARATE AUTHZ & AUTHN", "items": [
            "OAuth 2.0 = Authorization (access token, scopes)",
            "OIDC = Authentication (ID Token JWT, identity claims)",
            "Never use access token as identity assertion",
        ], "type": "bullet"},
        {"title": "PRINCIPLE 2: VERIFY EVERY TOKEN", "items": [
            "Access Token: Validate existence + scope in AS token store",
            "ID Token: Verify HMAC-SHA256 signature + all claims (iss, aud, exp, iat, nonce)",
            "Refresh Token: Track usage, rotate, detect replay",
        ], "type": "bullet"},
        {"title": "PRINCIPLE 3: DEFENSE IN DEPTH", "items": [
            "PKCE prevents code interception",
            "ID Token nonce prevents replay",
            "Refresh token rotation prevents long-term token reuse",
            "Family revocation limits blast radius",
        ], "type": "bullet"},
    ],
    mental_model="Separate authz/autn, verify every token, defense in depth — PKCE, nonce, rotation, revocation.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md",
    slide_num="23 / 24"
)

# Slide 24: Closing (Dark Navy, rich visual)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 31",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 17/17 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "KEY RESULTS", "items": [
            "PKCE S256: code_verifier binds client to authorization code",
            "ID Token: JWT with HMAC-SHA256 signature + claims validation",
            "Refresh Token Rotation: single-use, family revocation on replay",
            "Attack Scenarios: All blocked (interception, forgery, replay)",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL", "items": [
            "OAuth 2.0 is authorization (what you can access), OIDC is authentication (who you are).",
            "Never trust a JWT without cryptographic signature verification.",
            "Refresh tokens must be single-use — replay revokes entire session family.",
        ], "type": "bullet"},
        {"title": "RESOURCES", "items": [
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Run demo: go run ./cmd/demo",
            "Run tests: go test -v -race ./...",
        ], "type": "bullet"},
    ],
    mental_model="Separate authz/autn, verify JWT signature, rotate refresh tokens — defense in depth.",
    source_path="labs/31-oauth2-and-oidc/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
