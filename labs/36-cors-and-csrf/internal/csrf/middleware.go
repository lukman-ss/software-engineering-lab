package csrf

import (
	"net/http"
	"strings"
)

type Middleware struct {
	tokenManager *TokenManager
	cookieName   string
	headerName   string
	formFieldName string
}

func NewMiddleware(tm *TokenManager) *Middleware {
	return &Middleware{
		tokenManager:  tm,
		cookieName:    "session_id",
		headerName:    "X-CSRF-Token",
		formFieldName: "csrf_token",
	}
}

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
