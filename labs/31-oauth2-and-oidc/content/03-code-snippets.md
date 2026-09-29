## Snippet 1 — PKCE S256 Challenge Generation

**Source File:** `pkg/pkce/pkce.go` (lines 23-45)

**Purpose:** Menghasilkan `code_verifier` acak (32 byte, 43 karakter base64url) dan `code_challenge` menggunakan SHA-256.

```go
func GeneratePKCEPair(method string) (*PKCEPair, error) {
    if method != "S256" && method != "plain" {
        return nil, ErrInvalidMethod
    }

    // 32 random bytes -> 43 base64url chars
    buf := make([]byte, 32)
    if _, err := rand.Read(buf); err != nil {
        return nil, fmt.Errorf("failed to generate random verifier: %w", err)
    }
    verifier := base64.RawURLEncoding.EncodeToString(buf)

    challenge, err := ComputeChallenge(verifier, method)
    if err != nil {
        return nil, err
    }

    return &PKCEPair{
        CodeVerifier:  verifier,
        CodeChallenge: challenge,
        Method:        method,
    }, nil
}

func ComputeChallenge(verifier string, method string) (string, error) {
    if len(verifier) < 43 || len(verifier) > 128 {
        return "", ErrInvalidVerifierLength
    }

    switch method {
    case "S256":
        h := sha256.Sum256([]byte(verifier))
        return base64.RawURLEncoding.EncodeToString(h[:]), nil
    case "plain":
        return verifier, nil
    default:
        return "", ErrInvalidMethod
    }
}
```

**Explanation:** Verifier harus 43-128 karakter (entropy minimal 256 bit). S256 menghitung SHA-256 dari verifier dan meng-encode hasilnya ke base64url. Method `plain` hanya menyalin verifier ke challenge (tidak direkomendasikan untuk production).

---

## Snippet 2 — PKCE Verification

**Source File:** `pkg/pkce/pkce.go` (lines 63-72)

**Purpose:** Memverifikasi bahwa verifier yang dikirim client cocok dengan challenge yang tersimpan di server.

```go
func Verify(verifier string, challenge string, method string) error {
    expected, err := ComputeChallenge(verifier, method)
    if err != nil {
        return err
    }
    if expected != challenge {
        return ErrChallengeMismatch
    }
    return nil
}
```

**Explanation:** Server menghitung ulang challenge dari verifier yang diberikan client, lalu membandingkannya dengan challenge yang disimpan saat authorization request. Jika cocok, client terbukti memiliki verifier asli. (Note: audit menandai penggunaan string equality bukan constant-time compare sebagai LOW severity issue.)

---

## Snippet 3 — ID Token Signing (HMAC-SHA256)

**Source File:** `pkg/oidc/oidc.go` (lines 40-62)

**Purpose:** Membentuk JWT ID Token dengan header dan claims, lalu menandatangani menggunakan HMAC-SHA256.

```go
func SignIDToken(claims IDTokenClaims, secret []byte) (string, error) {
    header := Header{Alg: "HS256", Typ: "JWT"}
    headerJSON, err := json.Marshal(header)
    if err != nil {
        return "", err
    }
    claimsJSON, err := json.Marshal(claims)
    if err != nil {
        return "", err
    }

    headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
    claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

    unsignedToken := fmt.Sprintf("%s.%s", headerB64, claimsB64)

    mac := hmac.New(sha256.New, secret)
    mac.Write([]byte(unsignedToken))
    sig := mac.Sum(nil)
    sigB64 := base64.RawURLEncoding.EncodeToString(sig)

    return fmt.Sprintf("%s.%s", unsignedToken, sigB64), nil
}
```

**Explanation:** JWT terdiri dari tiga bagian base64url: header, claims, signature. Signature dihitung menggunakan HMAC-SHA256 dengan shared secret. Output adalah string `header.claims.signature`.

---

## Snippet 4 — ID Token Parsing and Verification

**Source File:** `pkg/oidc/oidc.go` (lines 64-116)

**Purpose:** Memverifikasi signature, issuer, audience, expiration, issued-at, dan nonce pada ID Token.

```go
func ParseAndVerifyIDToken(rawJWT string, secret []byte, expectedIssuer, expectedAudience, expectedNonce string, now time.Time) (*IDTokenClaims, error) {
    parts := strings.Split(rawJWT, ".")
    if len(parts) != 3 {
        return nil, ErrMalformedJWT
    }

    unsignedToken := fmt.Sprintf("%s.%s", parts[0], parts[1])
    sig, err := base64.RawURLEncoding.DecodeString(parts[2])
    if err != nil {
        return nil, ErrMalformedJWT
    }

    mac := hmac.New(sha256.New, secret)
    mac.Write([]byte(unsignedToken))
    expectedSig := mac.Sum(nil)

    if !hmac.Equal(sig, expectedSig) {
        return nil, ErrSignatureInvalid
    }

    claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, ErrMalformedJWT
    }

    var claims IDTokenClaims
    if err := json.Unmarshal(claimsBytes, &claims); err != nil {
        return nil, ErrMalformedJWT
    }

    if claims.Issuer != expectedIssuer {
        return nil, fmt.Errorf("%w: expected %s, got %s", ErrIssuerMismatch, expectedIssuer, claims.Issuer)
    }

    if claims.Audience != expectedAudience {
        return nil, fmt.Errorf("%w: expected %s, got %s", ErrAudienceMismatch, expectedAudience, claims.Audience)
    }

    unixNow := now.Unix()
    if claims.Expiration <= unixNow {
        return nil, ErrTokenExpired
    }

    if claims.IssuedAt > unixNow+300 {
        return nil, ErrIssuedInFuture
    }

    if expectedNonce != "" && claims.Nonce != expectedNonce {
        return nil, fmt.Errorf("%w: expected %s, got %s", ErrNonceMismatch, expectedNonce, claims.Nonce)
    }

    return &claims, nil
}
```

**Explanation:** Validasi mencakup enam langkah kritis: format JWT, signature HMAC (constant-time), issuer exact match, audience exact match, expiration (tidak expired), issued-at (tidak lebih dari 5 menit ke depan), dan nonce match (anti-replay).

---

## Snippet 5 — Authorization Server: Refresh Token Rotation

**Source File:** `pkg/server/server.go` (lines 219-286)

**Purpose:** Mengimplementasikan refresh token rotation dan deteksi replay berdasarkan RFC 9700 Section 4.14.

```go
func (s *AuthorizationServer) Refresh(refreshToken, clientID string) (*TokenResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    meta, exists := s.refreshMeta[refreshToken]
    if !exists {
        return nil, ErrRefreshTokenNotFound
    }

    // Check if family is revoked (due to prior replay detection)
    if s.revokedFams[meta.FamilyID] {
        return nil, ErrTokenReplayDetected
    }

    // Refresh Token Rotation: If already revoked/used, someone is replaying a stolen token!
    if meta.Revoked {
        // Invalidate the entire token family
        s.revokedFams[meta.FamilyID] = true
        return nil, ErrTokenReplayDetected
    }

    if meta.ClientID != clientID {
        return nil, ErrUnauthorizedClient
    }

    if time.Now().After(meta.ExpiresAt) {
        return nil, ErrInvalidGrant
    }

    // Invalidate the consumed refresh token
    meta.Revoked = true

    // Issue new rotated refresh token in the same family
    newRtBytes := make([]byte, 24)
    if _, err := rand.Read(newRtBytes); err != nil {
        return nil, err
    }
    newRefreshToken := "rt_" + hex.EncodeToString(newRtBytes)

    s.refreshMeta[newRefreshToken] = &RefreshTokenMeta{
        FamilyID:  meta.FamilyID,
        Subject:   meta.Subject,
        ClientID:  meta.ClientID,
        Scope:     meta.Scope,
        Revoked:   false,
        ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
    }

    // Issue new access token
    atBytes := make([]byte, 24)
    if _, err := rand.Read(atBytes); err != nil {
        return nil, err
    }
    newAccessToken := "at_" + hex.EncodeToString(atBytes)
    s.tokens[newAccessToken] = &AccessTokenMeta{
        Subject:   meta.Subject,
        Scope:     meta.Scope,
        ExpiresAt: time.Now().Add(1 * time.Hour),
    }

    return &TokenResponse{
        AccessToken:  newAccessToken,
        TokenType:    "Bearer",
        ExpiresIn:    3600,
        RefreshToken: newRefreshToken,
        Scope:        meta.Scope,
    }, nil
}
```

**Explanation:** Alur rotation: (1) cek apakah family sudah dicabut, (2) cek apakah token sudah revoked (replay detection), (3) mark token lama revoked, (4) buat token baru dengan family ID yang sama. Jika terjadi replay, seluruh family di-revok dan client dipaksa re-auth.

---

## Snippet 6 — Client Authorization Request Builder

**Source File:** `pkg/client/client.go` (lines 36-56)

**Purpose:** Membangun authorization request dengan PKCE S256, state, dan nonce.

```go
func (c *Client) BuildAuthorizationRequest(scope string) (string, error) {
    pair, err := pkce.GeneratePKCEPair("S256")
    if err != nil {
        return "", err
    }
    c.Verifier = pair.CodeVerifier

    stateBytes := make([]byte, 16)
    if _, err := rand.Read(stateBytes); err != nil {
        return "", err
    }
    c.State = hex.EncodeToString(stateBytes)

    nonceBytes := make([]byte, 16)
    if _, err := rand.Read(nonceBytes); err != nil {
        return "", err
    }
    c.Nonce = hex.EncodeToString(nonceBytes)

    return pair.CodeChallenge, nil
}
```

**Explanation:** Client menghasilkan verifier, state (untuk anti-CSRF), dan nonce (untuk anti-replay pada ID Token). Challenge dikembalikan untuk disertakan dalam authorization URL. Verifier disimpan di client untuk tahap token exchange.

---

## Snippet 7 — End-to-End Flow Orchestration

**Source File:** `pkg/client/client.go` (lines 58-76)

**Purpose:** Menukar authorization code menjadi tokens dan memverifikasi ID Token.

```go
func (c *Client) Exchange(code string) (*server.TokenResponse, error) {
    resp, err := c.Server.ExchangeCode(code, c.ClientID, c.RedirectURI, c.Verifier)
    if err != nil {
        return nil, err
    }

    c.AccessToken = resp.AccessToken
    c.RefreshToken = resp.RefreshToken

    if resp.IDToken != "" {
        claims, err := oidc.ParseAndVerifyIDToken(resp.IDToken, c.SigningKey, c.Server.Issuer, c.ClientID, c.Nonce, time.Now())
        if err != nil {
            return nil, fmt.Errorf("id_token validation failed: %w", err)
        }
        c.IDClaims = claims
    }

    return resp, nil
}
```

**Explanation:** Client mengirimkan code + verifier ke server. Server memverifikasi PKCE, memeriksa satu-time code, lalu menerbitkan tokens. Client kemudian memverifikasi ID Token menggunakan signing key, issuer, audience, dan nonce.
