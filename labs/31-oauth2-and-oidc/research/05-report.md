# Research Report

## Research Question
Mendefinisikan secara preskriptif perbedaan OAuth 2.0 (Authorization) vs OIDC (Authentication), mekanisme Authorization Code Flow + PKCE, struktur dan validasi ID Token JWT, serta jebakan keamanan umum (token storage, rotation, signature/audience validation).

## Executive Summary
OAuth 2.0 adalah protokol *delegated authorization*: menghasilkan access token yang menjawab "resource apa yang boleh diakses", bukan "siapa penggunanya". OIDC menambahkan identity layer di atas OAuth 2.0 dan menghasilkan ID Token (JWT) berisi klaim autentikasi. Memakai access token OAuth untuk login adalah celah keamanan karena access token tidak terikat ke identity RP dan tidak memiliki validation rule `iss`/`aud`/`exp` yang distandarisasi untuk auth. Standards modern mewajibkan Authorization Code Flow + PKCE untuk semua client (RFC 9700, OAuth 2.1), deprecates implicit flow, dan mewajibkan refresh token rotation/sender-constraining untuk public clients.

## Findings

### Finding 1 — OAuth 2.0 bukan protokol autentikasi
Claim: OAuth 2.0 hanya mendefinisikan delegated access (authorization grant -> access token -> protected resource), tanpa standardisasi siapa pengguna.
Evidence: RFC 6749 Abstract/Sec 1 ("limited access to an HTTP service"); Sec 1.4 (access token "usually opaque to the client"). OIDC Core Sec 1 secara eksplisit: "without profiling OAuth 2.0, it is incapable of providing information about the authentication of an End-User."
Sources: RFC 6749; OIDC Core Sec 1
Confidence: HIGH

### Finding 2 — OIDC menambahkan identity layer via ID Token JWT + UserInfo
Claim: OIDC = OAuth 2.0 + identity. Verifikasi identity dilakukan oleh RP terhadap ID Token JWT bertanda tangan, bukan terhadap access token.
Evidence: OIDC Core Abstract/Sec 1.3: "The primary extension... is the ID Token data structure... represented as a JWT." Required claims: iss, sub, aud (must contain client_id), exp, iat. UserInfo Endpoint (REST) berisi profile claims, dilindungi oleh access token.
Sources: OIDC Core Sec 1, 2, 3, 5.3
Confidence: HIGH

### Finding 3 — ID Token wajib divalidasi multi-langkah (signature + iss + aud + exp + nonce)
Claim: Validasi ID Token mencakup 13 langkah: decrypt-if-encrypted, iss exact match, aud contains client_id, JWS signature via issuer JWKS, current time before exp, nonce anti-replay, alg pinning, optional acr/auth_time.
Evidence: OIDC Core Sec 3.1.3.7 (13 verbatim steps), Sec 2 (claim requirements).
Sources: OIDC Core Sec 3.1.3.7, 2
Confidence: HIGH
Corroborated By: RFC 8725 Sec 3.1/3.8/3.9 (algorithm verification, issuer/subject/audience validation)

### Finding 4 — PKCE mencegah authorization-code interception; S256 wajib
Claim: PKCE binds authorization request to token request via code_verifier/code_challenge (S256 = BASE64URL(SHA256(verifier))); S256 adalah MTI; plain deprecated.
Evidence: RFC 7636 Sec 4.1/4.2/4.6 formulas + Appendix B worked example; Sec 7.2 plain SHOULD NOT be used; RFC 9700 Sec 2.1.1 public clients MUST use PKCE.
Sources: RFC 7636; RFC 9700 Sec 2.1.1
Confidence: HIGH

### Finding 5 — Implicit flow deprecated; Authorization Code Flow + PKCE standar modern
Claim: Implicit flow (response_type=token/id_token) mengembalikan token di URL fragment, rentan terhadap XSS, referer, history; RFC 9700: SHOULD NOT digunakan; OAuth 2.1: dihapus.
Evidence: RFC 9700 Sec 2.1.2; OAuth 2.1 summary (implicit omitted); Browser-based-apps draft Sec 7.2 (threat analysis).
Sources: RFC 9700; https://oauth.net/2.1/
Confidence: HIGH

### Finding 6 — Jebakan umum: (a) tidak verifikasi signature + aud, (b) access token di localStorage, (c) refresh token tidak di-rotate
Claim: Tanpa verifikasi signature/iss/aud, siapa saja bisa buat JWT palsu. LocalStorage rentan XSS. Refresh token statis meningkatkan dampak pembajakan.
Evidence: RFC 8725 Sec 2.1 (alg=none / RS256->HS256 confusion); Browser-based-apps draft Sec 5 (malicious JS same privileges, steal from localStorage/IndexedDB); RFC 9700 Sec 4.14 (rotation invalidates old token, detects replay); OIDC Core Sec 2 (aud wajib, else reject).
Sources: RFC 8725; draft-ietf-oauth-browser-based-apps-27; RFC 9700 Sec 4.14; OIDC Core Sec 2
Confidence: HIGH

### Finding 7 — Refresh token: rotation wajib untuk public client
Claim: Rotation = server menerbitkan refresh token baru tiap refresh, lama di-invalidate; jika keduanya dipakai, server deteksi breach dan revoke.
Evidence: RFC 9700 Sec 4.14.2 (verbatim rotation text); Sec 2.2.2 (MUST).
Sources: RFC 9700
Confidence: HIGH

## Areas of Agreement
- OIDC = OAuth 2.0 + identity (ID Token JWT) — sepakat antara OIDC Core dan RFC 9700.
- PKCE mandatory — RFC 9700, OAuth 2.1, RFC 7636 sepakat.
- Implicit flow insecure — sepakat (meskipun OIDC Core masih mendokumentasikannya secara historis).
- Refresh token harus di-rotate atau sender-constrained — RFC 9700, OAuth 2.1 sepakat; OIDC Core contoh output kompatibel tapi tidak mewajibkan secara eksplisit.
- JWT validation harus pin alg, validasi iss/aud/exp — RFC 7519, 8725, OIDC Core sepakat.

## Areas of Disagreement
- OIDC Core masih menyebut implicit/hybrid flow; RFC 9700/OAuth 2.1 mendeprecate/hapus. Status: OIDC Core sudah usang secara praktis untuk browser-based apps; RFC 9700/OAuth 2.1 adalah BCP terkini.
- OIDC Core Sec 12 menyebut refresh_token pada response tanpa menyebut "rotation"; RFC 9700 mewajibkan rotation. Status: OIDC Core tidak kontradiktif, hanya kurang preskriptif — RFC 9700 berlaku sebagai BCP keamanan.

## Limitations
- OIDC Core 1.0 (2014/2023 errata) tidak mencerminkan hardening terbaru (RFC 9700 Jan 2025).
- OAuth 2.1 masih draft (saat penelitian); beberapa terminologi bisa berubah saat final.
- Browser-based-apps draft (Jul 2026, expires Jan 2027) belum final RFC (RFC 10017 belum diterbitkan saat akses).
- Studi kasus numerik (SHA256 example) diambil dari RFC 7636 Appendix B — nilai deterministik, bukan benchmark performa.

## Conclusion
Oleh karena itu, memisahkan Authorization (OAuth 2.0: access token, apa yang boleh diakses) dari Authentication (OIDC: ID Token JWT, siapa penggunanya) adalah fondasi arsitektur yang aman. Authorization Code Flow + PKCE + ID Token signature/iss/aud/exp validation + refresh token rotation + BFF/token storage yang aman adalah satu set prasyarat berbasis standar yang konsisten antar RFC 6749, 7636, 7519, 8725, 9700, dan OIDC Core.

## Lab Implementation Guidance (evidence-based)
1. Buat code_verifier = crypto.random(32 octets) base64url (43 char, >=256-bit entropy).
2. code_challenge = BASE64URL(SHA256(ASCII(verifier))) (S256, mandatory).
3. Auth request: response_type=code + scope=openid... + code_challenge + code_challenge_method=S256 + state + nonce.
4. Token exchange: POST /token grant_type=authorization_code + code + redirect_uri + code_verifier + client_auth.
5. ID Token validation (13 langkah): decrypt -> iss exact -> aud contains client_id -> signature via issuer JWKS -> alg pinned -> exp -> nonce -> optional auth_time/acr.
6. Token storage: BFF (httpOnly Secure SameSite cookie) atau in-memory untuk SPA; hindari localStorage.
7. Refresh: rotasi + invalidasi lama; deteksi reuse -> revoke semua token + force re-auth.
