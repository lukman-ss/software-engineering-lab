// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"regexp"
)

const HeaderRequestID = "X-Request-ID"

type contextKey string

const contextKeyRequestID contextKey = "request_id"

// requestIDRegex validates incoming X-Request-ID values.
// Accepts alphanumeric IDs with dashes/underscores, max 128 chars.
var requestIDRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,128}$`)

// RequestID extracts or generates a request ID.
// If the incoming X-Request-ID is valid, it is preserved; otherwise a new one is generated.
// The ID is stored in the request context and returned as X-Request-ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if !isValidRequestID(id) {
			id = generateRequestID()
		}

		w.Header().Set(HeaderRequestID, id)

		ctx := context.WithValue(r.Context(), contextKeyRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID returns the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(contextKeyRequestID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func isValidRequestID(s string) bool {
	if s == "" {
		return false
	}
	return requestIDRegex.MatchString(s)
}

func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
