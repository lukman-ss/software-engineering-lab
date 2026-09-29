# Code Snippets

Dokumen ini mendokumentasikan cuplikan kode terverifikasi dari implementasi `labs/36-cors-and-csrf`.

---

## Snippet 1 — CORS Spec Compliance & Preflight Handling

Source File: `internal/cors/middleware.go:44-93`  
Purpose: Menangani negosiasi origin, melarang wildcard `*` ketika kredensial diaktifkan, dan merespons preflight request `OPTIONS` secara spec-compliant.

```go
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !m.isOriginAllowed(origin) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Set CORS headers
		if m.config.AllowCredentials {
			// Spec: If credentials are included, wildcard "*" is illegal. Must reflect specific origin.
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			if len(m.config.AllowedOrigins) == 1 && m.config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		w.Header().Add("Vary", "Origin")

		if len(m.config.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.config.ExposedHeaders, ", "))
		}

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))
			if m.config.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(m.config.MaxAge))
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

Explanation:
- Jika origin tidak cocok dan method adalah `OPTIONS` (preflight), middleware mengembalikan `403 Forbidden`.
- Jika origin tidak cocok tetapi method adalah simple request (misal `POST`), middleware membiarkan request dieksekusi oleh handler aplikasi tanpa menyetel header CORS. Inilah mengapa CORS tidak melindungi backend dari form submission cross-origin.
- Menegakkan larangan spesifikasi W3C Fetch terhadap wildcard `*` ketika `AllowCredentials: true`.

---

## Snippet 2 — HMAC-SHA256 Session-Bound Token Generation

Source File: `internal/csrf/token.go:44-59`  
Purpose: Menghasilkan token Anti-CSRF stateless bertanda tangan HMAC yang terikat pada ID sesi dan timestamp kedaluwarsa.

```go
// GenerateToken creates a signed token: base64(sessionID:timestamp:nonce:hmac)
func (tm *TokenManager) GenerateToken(sessionID string) string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	nonceStr := base64.RawURLEncoding.EncodeToString(nonce)
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	payload := sessionID + ":" + ts + ":" + nonceStr
	h := hmac.New(sha256.New, tm.secret)
	h.Write([]byte(payload))
	signature := h.Sum(nil)

	tokenRaw := payload + ":" + base64.RawURLEncoding.EncodeToString(signature)
	return base64.RawURLEncoding.EncodeToString([]byte(tokenRaw))
}
```

Explanation:
- Menggunakan CSPRNG (`crypto/rand`) untuk menghasilkan nilai nonce 16-byte acak.
- Menggabungkan `sessionID`, `timestamp`, dan `nonce` ke dalam payload.
- Menandatangani payload dengan kunci rahasia server menggunakan algoritma HMAC-SHA256.

---

## Snippet 3 — Constant-Time Token & Session Validation

Source File: `internal/csrf/token.go:62-110`  
Purpose: Memvalidasi integritas HMAC, mencocokkan session ID aktif pengguna, mengecek TTL token, dan mencegah serangan timing attacks.

```go
// ValidateToken verifies HMAC signature, sessionID match, and expiration.
func (tm *TokenManager) ValidateToken(tokenStr, expectedSessionID string) error {
	if tokenStr == "" {
		return ErrTokenMissing
	}

	decoded, err := base64.RawURLEncoding.DecodeString(tokenStr)
	if err != nil {
		return ErrInvalidToken
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 4 {
		return ErrInvalidToken
	}

	sessionID := parts[0]
	tsStr := parts[1]
	nonceStr := parts[2]
	sigB64 := parts[3]

	if sessionID != expectedSessionID {
		return ErrInvalidToken
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return ErrInvalidToken
	}

	if time.Since(time.Unix(ts, 0)) > tm.ttl {
		return ErrExpiredToken
	}

	payload := sessionID + ":" + tsStr + ":" + nonceStr
	h := hmac.New(sha256.New, tm.secret)
	h.Write([]byte(payload))
	expectedSig := h.Sum(nil)

	actualSig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return ErrInvalidToken
	}

	if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {
		return ErrInvalidToken
	}

	return nil
}
```

Explanation:
- Mengecek kesesuaian `sessionID` dengan sesi aktif pengguna pemanggil untuk mencegah token reuse lintas sesi.
- Memeriksa batas TTL token untuk mencegah replay attack jangka panjang.
- Menggunakan `crypto/subtle.ConstantTimeCompare` saat membandingkan tanda tangan digital untuk mencegah kebocoran informasi melalui *timing side-channel attacks*.

---

## Snippet 4 — Anti-CSRF Token Enforcement Middleware

Source File: `internal/csrf/middleware.go:24-53`  
Purpose: Mengekstrak token dari HTTP Header atau Form Body pada method state-changing dan menolak request tidak terverifikasi.

```go
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Safe methods are idempotent / read-only under HTTP specs
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve session ID from cookie
		cookie, err := r.Cookie(m.cookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized: missing session cookie", http.StatusUnauthorized)
			return
		}
		sessionID := cookie.Value

		// Check token from header or form body
		token := r.Header.Get(m.headerName)
		if token == "" {
			token = r.PostFormValue(m.formFieldName)
		}

		if err := m.tokenManager.ValidateToken(token, sessionID); err != nil {
			http.Error(w, "Forbidden: CSRF token validation failed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

Explanation:
- Mengabaikan method aman (`GET`, `HEAD`, `OPTIONS`, `TRACE`) sesuai spesifikasi RFC HTTP.
- Memeriksa cookie sesi pengguna; jika tidak ada, menolak dengan `401 Unauthorized`.
- Membaca token dari header kustom `X-CSRF-Token` atau field form `csrf_token`.
- Jika validasi token gagal, menghentikan eksekusi dengan `403 Forbidden` sebelum handler bisnis disentuh.

---

## Snippet 5 — Fetch Metadata (`Sec-Fetch-Site`) & Custom Header Middleware

Source File: `internal/csrf/middleware.go:56-86`  
Purpose: Pertahanan modern berbasis header browser dan penegakan custom header API untuk memblokir penyerangan form HTML langsung.

```go
// FetchMetadataMiddleware checks Sec-Fetch-Site modern browser header.
func FetchMetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchSite := r.Header.Get("Sec-Fetch-Site")

		// If header is present and request is cross-site on state-changing methods, reject.
		if fetchSite == "cross-site" {
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				http.Error(w, "Forbidden: Cross-site request rejected by Sec-Fetch-Site", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireCustomHeaderMiddleware enforces custom headers for API endpoints (defeats simple form requests).
func RequireCustomHeaderMiddleware(headerKey, expectedValue string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				val := r.Header.Get(headerKey)
				if val == "" || (expectedValue != "" && !strings.EqualFold(val, expectedValue)) {
					http.Error(w, "Forbidden: Missing or invalid required custom header", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

Explanation:
- `FetchMetadataMiddleware` memanfaatkan header `Sec-Fetch-Site` yang disetel otomatis oleh browser. Menolak request mutasi state yang berasal dari konteks `cross-site`.
- `RequireCustomHeaderMiddleware` mewajibkan kehadiran header kustom (seperti `X-Requested-With: XMLHttpRequest`), yang mustahil dikirim oleh formulir HTML standar tanpa memicu preflight CORS.
