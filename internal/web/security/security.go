package security

import (
	"net/http"
	"strings"
)

// HideGit returns 404 for URL paths containing a .git segment.
func HideGit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, segment := range strings.Split(r.URL.Path, "/") {
			if segment == ".git" {
				w.Header().Set("Cache-Control", "no-store")
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// HideEnv returns 404 for URL paths that could refer to environment files.
func HideEnv(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, segment := range strings.Split(r.URL.Path, "/") {
			name := strings.ToLower(strings.TrimLeft(segment, "."))
			if name == "env" || strings.HasPrefix(name, "env.") || strings.HasPrefix(name, "env-") || strings.HasPrefix(name, "env_") {
				w.Header().Set("Cache-Control", "no-store")
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Headers adds baseline browser security headers to every response.
func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
