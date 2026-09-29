# OAuth 2.0 & OpenID Connect (OIDC): Authentication vs Authorization dan Mekanisme Keamanan

## Problem

Sistem web modern sering membingungkan konsep **authorization** (otorisasi) dengan **authentication** (autentikasi). OAuth 2.0 adalah protokol untuk delegated access: ia memberikan access token yang menjawab pertanyaan "sumber daya apa yang boleh diakses oleh client ini?" OIDC menambahkan identity layer di atas OAuth 2.0 sehingga kita juga bisa menjawab "siapa pengguna yang terautentikasi?". Kegagalan memisahkan kedua konsep ini menghasilkan celah keamanan: akses token yang tidak diverifikasi ditukar sebagai bukti identitas, flow implicit yang rentan terhadap XSS masih dipakai, atau refresh token yang statis menyebabkan serangan replay.

Lab ini mendemonstrasikan bagaimana Authorization Code Flow dikombinasikan dengan PKCE mencegah code interception, bagaimana ID Token (JWT) divalidasi secara kriptografis, dan bagaimana refresh token rotation dengan token family revocation mendeteksi token yang disadap.

## Why This Matters

Standar keamanan terkini (RFC 9700, OAuth 2.1) mewajibkan Authorization Code Flow + PKCE untuk semua client berbasis browser dan menghapus implicit flow. Tanpa pemahaman yang tepat, implementasi dapat melakukan kesalahan kritis seperti menyimpan bearer token di localStorage, mengabaikan validasi `iss`/`aud`/`exp`, atau menggunakan refresh token statis yang memungkinkan attacker mempertahankan sesi setelah token dicuri.

## Mental Model

Bayangkan dua gerbang berbeda:

1. **Authorization Gate (OAuth 2.0)**: Memberikan kartu akses (`access_token`) yang menunjukkan area mana yang boleh dimasuki oleh client tertentu. Kartu ini berbatas waktu dan memiliki scope.
2. **Authentication Gate (OIDC)**: Memberikan kartu identitas (`id_token`, berupa JWT bertanda tangan) yang menyatakan siapa pengguna sebenarnya. Kartu ini diverifikasi tanda tangannya oleh RP.

PKCE berfungsi seperti segel pengaman pada authorization code: hanya pihak yang memiliki `code_verifier` asli yang bisa menukarnya menjadi token. Refresh token rotation memastikan bahwa setiap penggunaan token segar menggantikan token lama; jika token lama muncul kembali (replay), seluruh keluarga token dianggap rusak dan dicabut.

## Core Concept

Lab ini mengimplementasikan empat mekanisme inti:

- **PKCE (Proof Key for Code Exchange)**: `code_verifier` acak di-generate oleh client, `code_challenge` (SHA-256 hash) dikirim dalam authorization request. Server menyimpan challenge, client mengirim verifier saat token exchange.
- **Authorization Code Grant**: Client bermitra dengan Authorization Server (AS) untuk mendapatkan authorization code, lalu menukarnya dengan access token, refresh token, dan ID token.
- **ID Token Validation**: ID Token adalah JWT dengan header `alg: HS256` dan claims wajib (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`). Validasi meliputi dekripsi (jika ada), kecocokan issuer, audience, signature HMAC, expiration, dan nonce.
- **Refresh Token Rotation & Replay Detection**: Setiap refresh menghasilkan token baru; token lama di-invalidate. Jika token lama digunakan ulang, seluruh `FamilyID` ditandai compromised dan semua token aktif dicabut.

## Failure Scenario

- **Code Interception**: Attacker mencuri authorization code dari URL callback. Tanpa `code_verifier` yang sesuai, exchange akan gagal.
- **Forged ID Token**: Modifikasi payload atau signature pada JWT akan ditolak pada tahap verifikasi HMAC dan claims.
- **Refresh Token Replay**: Attacker yang menangkap refresh token lama akan memicu deteksi reuse, yang kemudian mencabut seluruh session terkait.

## How It Works

### 1. Authorization Request

Client membuat PKCE pair (verifier + challenge S256) dan nonce acak. Client mengarahkan authorization request ke AS dengan parameter:

```
response_type=code
client_id=spa-client-123
redirect_uri=https://app.example.com/callback
scope=openid profile email
state=<random>
code_challenge=<challenge>
code_challenge_method=S256
nonce=<random>
```

### 2. Authorization Code Issuance

AS memvalidasi client ID dan redirect URI, menyimpan code_challenge serta nonce pada authorization code yang baru dibuat (umur 5 menit). AS mengembalikan authorization code kepada client melalui redirect.

### 3. Token Exchange

Client mengirimkan POST ke endpoint `/token` dengan:
- `grant_type=authorization_code`
- `code=<authorization_code>`
- `redirect_uri`
- `client_id`
- `code_verifier`

AS memverifikasi PKCE: menghitung SHA-256 dari verifier dan membandingkannya dengan stored challenge. Jika cocok, AS menerbitkan access token (opaque bearer), refresh token, dan ID Token (JWT). Authorization code ditandai `Used=true`.

### 4. ID Token Validation

Client memverifikasi ID Token menggunakan kunci signing yang sama dengan AS. Langkah validasi:
1. Decode header dan claims.
2. Verifikasi signature HMAC-SHA256.
3. Validasi `iss` exact match.
4. Validasi `aud` mengandung client ID.
5. Cek `exp` (masih berlaku).
6. Cek `iat` (tidak melebihi sekarang + 5 menit).
7. Cocokkan `nonce` dengan yang dikirim pada request awal.
*Catatan: Implementasi ini menggunakan 7 langkah validasi dengan HMAC-SHA256 symmetric signing untuk tujuan pendidikan; lengkap OIDC Core 3.1.3.7 meliputi ekstra checks seperti alg pinning, JWKS-based signature verification, at_hash/c_hash, dan decryption untuk JWE.*

### 5. Protected Resource Access

Resource server menerima Bearer Access Token, mencocokkannya dengan metadata di AS, dan memverifikasi scope yang diminta.

### 6. Refresh Token Rotation

Ketika access token hampir expired, client mengirimkan refresh token ke AS. AS:
1. Menandai refresh token lama sebagai `Revoked=true`.
2. Menerbitkan access token baru dan refresh token baru (dengan `FamilyID` yang sama).
3. Jika refresh token yang sudah di-revoke muncul kembali, AS menandai `FamilyID` sebagai compromised dan mengembalikan error replay detection.

## Architecture

```text
+-------------+
|   Client    |----+ (1) Auth Request (code_challenge, nonce)
| Application |    |
+-------------+    v
       |         +-----------------------+
       |         | Authorization Server  | (Handles AS + OIDC Provider)
       |         +-----------------------+
       |            | (2) Auth Code
       |<-----------+
       |
       | (3) Token Exchange (code + code_verifier)
       |------------------->+-----------------------+
       |<-------------------| AS validates PKCE     |
       | (4) Tokens         | Issues:               |
       |  - Access Token    | - Access Token        |
       |  - ID Token (JWT)  | - ID Token (signed)   |
       |  - Refresh Token   | - Refresh Token       |
       |                    +-----------------------+
       v
+-------------+
|  Resource   | (5) Request with Bearer Access Token
|   Server    |------------------------------------> [Returns Protected Resource]
+-------------+
```

## Implementation

Lab ini menggunakan Go standar library tanpa dependensi eksternal. Empat paket utama:

- `pkg/pkce/pkce.go`: Generator verifier/challenge dan verifikasi PKCE.
- `pkg/oidc/oidc.go`: Pembentukan dan validasi ID Token JWT (HMAC-SHA256).
- `pkg/server/server.go`: Authorization Server in-memory dengan rotasi refresh token dan deteksi replay.
- `pkg/client/client.go`: Helper client untuk meng orchestrate flow.
- `cmd/demo/main.go`: Demo end-to-end.
- `tests/oauth_test.go`: 17 unit test dengan race detector.

Keputusan implementasi:
- Simetri HMAC-SHA256 dipilih agar lab tetap mandiri (self-contained) tanpa perlu JWKS endpoint atau kunci asimetris.
- Refresh token dilacak menggunakan `FamilyID` untuk memungkinkan revokasi seluruh keluarga saat replays terdeteksi.
- Semua operasi stateful dilindungi mutex agar aman dari race condition.

## Code Walkthrough

Berikut adalah potongan kode penting yang merepresentasikan alur inti. Setiap snippet menyertakan sumber file asli.

### PKCE Generation and Verification

**File:** `pkg/pkce/pkce.go`

```go
// Generate 32 random bytes -> 43 base64url chars
buf := make([]byte, 32)
rand.Read(buf)
verifier := base64.RawURLEncoding.EncodeToString(buf)

// S256 challenge
h := sha256.Sum256([]byte(verifier))
challenge := base64.RawURLEncoding.EncodeToString(h[:])
```

Penjelasan: `code_verifier` memiliki entropy minimal 256 bit. `code_challenge` adalah BASE64URL(SHA256(verifier)). Validasi dilakukan dengan menghitung ulang challenge dari verifier dan membandingkannya dengan stored challenge.

### ID Token Signing and Verification

**File:** `pkg/oidc/oidc.go`

```go
header := Header{Alg: "HS256", Typ: "JWT"}
headerJSON, _ := json.Marshal(header)
claimsJSON, _ := json.Marshal(claims)

headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

unsignedToken := fmt.Sprintf("%s.%s", headerB64, claimsB64)
mac := hmac.New(sha256.New, secret)
mac.Write([]byte(unsignedToken))
sig := mac.Sum(nil)
sigB64 := base64.RawURLEncoding.EncodeToString(sig)

return fmt.Sprintf("%s.%s", unsignedToken, sigB64), nil
```

Penjelasan: ID Token dibentuk sebagai JWT tiga bagian (header.claims.signature). Signature dihitung menggunakan HMAC-SHA256 dengan shared secret. Pada verifikasi, signature dihitung ulang dan dibandingkan menggunakan `hmac.Equal` (constant-time comparison).

### Refresh Token Rotation

**File:** `pkg/server/server.go`

```go
// Mark old refresh token as revoked
meta.Revoked = true

// Issue new rotated refresh token in the same family
newRtBytes := make([]byte, 24)
rand.Read(newRtBytes)
newRefreshToken := "rt_" + hex.EncodeToString(newRtBytes)

s.refreshMeta[newRefreshToken] = &RefreshTokenMeta{
    FamilyID: meta.FamilyID,
    Subject:  meta.Subject,
    Revoked:  false,
}
```

Penjelasan: Setelah refresh berhasil, token lama di-set `Revoked=true`. Token baru dibuat dengan `FamilyID` yang sama. Jika token lama (yang sudah revoked) digunakan kembali, server menandai `FamilyID` compromised dan menolak semua request berikutnya dengan error replay detection.

### Full Flow Orchestration

**File:** `pkg/client/client.go`

```go
func (c *Client) BuildAuthorizationRequest(scope string) (string, error) {
    pair, _ := pkce.GeneratePKCEPair("S256")
    c.Verifier = pair.CodeVerifier
    // ... generate state, nonce ...
    return pair.CodeChallenge, nil
}

func (c *Client) Exchange(code string) (*server.TokenResponse, error) {
    resp, _ := c.Server.ExchangeCode(code, c.ClientID, c.RedirectURI, c.Verifier)
    c.AccessToken = resp.AccessToken
    c.RefreshToken = resp.RefreshToken
    if resp.IDToken != "" {
        claims, _ := oidc.ParseAndVerifyIDToken(resp.IDToken, c.SigningKey, ...)
        c.IDClaims = claims
    }
    return resp, nil
}
```

Penjelasan: Client menyimpan `Verifier` secara lokal, membangun authorization request dengan PKCE dan nonce, kemudian menukar authorization code dengan tokens. ID Token divalidasi sebelum disimpan.

## What the Tests Prove

Test suite (`tests/oauth_test.go`) mencakup 17 kasus:

1. **PKCE S256 Valid**: Verifier panjang 43-128 karakter, challenge valid, verifikasi berhasil.
2. **PKCE Invalid Method**: Metode selain S256/plain ditolak.
3. **PKCE Mismatch**: Verifier dari pasangan berbeda ditolak.
4. **ID Token Valid**: Signature valid, claims cocok.
5. **ID Token Tampered Signature**: Modifikasi signature ditolak.
6. **ID Token Expired**: Token dengan `exp` di masa lalu ditolak.
7. **ID Token Mismatch Claims**: Issuer, audience, nonce yang salah ditolak.
8. **Full Flow and PKCE Interception**: Attacker dengan verifier salah tidak bisa exchange; code one-time use.
9. **Refresh Token Rotation and Replay Detection**: Rotasi berhasil; replay memicu family revocation.
10. **Plain Method**: Fungsi plain tetap bekerja (meski deprecated).
11. **Malformed JWT**: Struktur JWT yang salah ditolak.
12. **Negative Paths**: Client tidak terdaftar, redirect URI salah, code invalid, scope mismatch.
13. **ID Token Issued In Future**: `iat` lebih dari 5 menit ke depan ditolak.
14. **Verifier Length Bounds**: Verifier < 43 atau > 128 karakter ditolak.
15. **Expired Auth Code**: Authorization code yang kadaluarsa ditolak.
16. **Concurrent Refresh Replay**: 10 goroutine simultan mencoba refresh token curian; minimal 9 gagal.
17. **Concurrency and Race**: 20 worker concurrent; race detector tidak mendeteksi race condition.

Demo berjalan tanpa error (`go run ./cmd/demo`) mengonfirmasi bahwa setiap skenario attack gagal seperti yang diharapkan.

## Recovery / Rollback

Jika refresh token replay terdeteksi:
1. Server menandai `FamilyID` compromised.
2. Semua token aktif dalam keluarga tersebut dibatalkan.
3. Client menerima error `ErrTokenReplayDetected`.
4. Client harus meminta user untuk re-authenticate ulang (force login).

Rollback tidak otomatis; recovery memerlukan intervention dari user untuk memulai flow baru dari awal dengan PKCE dan nonce baru.

## Production Considerations

Implementasi lab ini adalah referensi pendidikan. Untuk produksi, perhatikan:

- Gunakan asimetris signing (RS256/ES256) dengan JWKS endpoint untuk key rotation.
- Gunakan constant-time comparison untuk PKCE verification (`crypto/subtle.ConstantTimeCompare`).
- Implementasikan token storage yang aman (httpOnly Secure SameSite cookie untuk BFF pattern, atau DPoP/mTLS untuk sender-constrained tokens).
- Tambahkan mekanisme graceful eviction untuk authorization codes dan tokens yang expired.
- Validasi state parameter untuk mencegah CSRF.
- Implementasikan refresh token binding (sender-constraining) sesuai RFC 9700 Section 4.14.

## Common Mistakes

1. **Menggunakan access token untuk authentication**: Access token OAuth bukan identitas; gunakan ID Token untuk membuktikan siapa pengguna.
2. **Tidak memvalidasi signature ID Token**: Siapapun bisa membuat JWT palsu tanpa validasi HMAC.
3. **Mengabaikan audience/issuer check**: ID Token bisa ditujukan untuk client lain.
4. **Menyimpan bearer token di localStorage**: Rentan terhadap XSS; gunakan httpOnly cookie atau in-memory storage.
5. **Menggunakan implicit flow**: Flow ini deprecated karena mengekspos token di URL fragment.
6. **Refresh token statis**: Token yang sama digunakan berulang memungkinkan attacker mempertahankan akses selamanya.

## Case Study

Berikut adalah contoh alur legimitas dan serangan dari demo (`cmd/demo/main.go`):

**Step 1**: Client membuat PKCE challenge S256 dan nonce.
```
Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w
Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de
```

**Step 2**: AS menerbitkan authorization code.
```
Issued Code: 5fabbdcdb6bbd582c3811481395a5161 (expires: 15:31:47)
```

**Step 3**: Client menukar code + verifier untuk tokens.
```
Access Token : at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2
Refresh Token: rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928
ID Token     : eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...
```

**Step 4**: ID Token divalidasi.
```
ID Token Subject  : user_alice_99
ID Token Issuer   : https://auth.example.com
ID Token Audience : spa-client-123
ID Token Nonce Match: SUCCESS (2152fa0c9254b322d11b27dba96b99de)
```

**Step 5**: Resource server memverifikasi access token.
```
Resource Server verified Access Token. Subject: user_alice_99
```

**Step 6**: Serangan interception diblokir.
```
Interception Attack Blocked! Error: pkce verification failed: code_verifier does not match code_challenge
```

**Step 7**: Refresh token rotation berhasil.
```
Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd
Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590
```

**Step 8**: Replay detection memicu family revocation.
```
Refresh Token Replay Detected! Error: refresh token reuse detected: entire token family revoked
Active session revoked due to family revocation! Error: refresh token reuse detected: entire token family revoked
```

## Checklist

- [ ] PKCE S256 digunakan, verifier disimpan aman.
- [ ] ID Token diverifikasi signature, issuer, audience, expiration, nonce.
- [ ] Authorization code bersifat one-time.
- [ ] Refresh token di-rotate setiap penggunaan.
- [ ] Replay detection memicu family revocation.
- [ ] Tidak ada implicit flow atau resource owner password flow.
- [ ] State parameter divalidasi (untuk anti-CSRF) — *catatan: demo ini tidak memvalidasi state server-side (arsitektur lab Demo)*.
- [ ] Token storage menghindari localStorage untuk bearer tokens.
- [ ] Concurrency aman (mutex, race detector passed).

## Key Takeaways

1. OAuth 2.0 dan OIDC memiliki peran berbeda: OAuth untuk authorization (access token), OIDC untuk authentication (ID token JWT).
2. PKCE S256 wajib digunakan untuk mencegah authorization code interception.
3. ID Token harus diverifikasi secara kriptografis (signature + iss + aud + exp + nonce).
4. Refresh token rotation memastikan single-use dan replay detection melindungi dari token replay.
5. Family revocation mencabut seluruh session ketika satu token dalam keluarga terdeteksi direplay.
6. Implisit flow dan penyimpanan token di localStorage adalah praktik berbahaya yang harus dihindari.
7. Implementasi standarisasi (RFC 7636, RFC 9700, OIDC Core) menyediakan landasan yang konsisten untuk keamanan.

## Sources

- RFC 6749 (OAuth 2.0 Authorization Framework)
- RFC 7636 (Proof Key for Code Exchange by OAuth Public Clients)
- RFC 7519 (JSON Web Token - JWT)
- RFC 8725 (OAuth 2.0 Security Best Current Practice)
- RFC 9700 (OAuth 2.0 Security Best Current Practice - 2025)
- OpenID Connect Core 1.0
- `labs/31-oauth2-and-oidc/research/` (research plan, evidence, report, contradictions)
- `labs/31-oauth2-and-oidc/engineering/` (design, implementation notes, execution result)
- Source code: `pkg/pkce/`, `pkg/oidc/`, `pkg/server/`, `pkg/client/`, `cmd/demo/`
- Tests: `tests/oauth_test.go`
- Demos: `go run ./cmd/demo`
