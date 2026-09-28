package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type Middleware struct {
	config Config
}

func New(config Config) *Middleware {
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization", "X-CSRF-Token", "X-Requested-With"}
	}
	return &Middleware{config: config}
}

func (m *Middleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range m.config.AllowedOrigins {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

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
