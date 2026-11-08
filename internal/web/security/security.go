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

// HideDotfiles returns 404 for URL paths containing a dot-prefixed segment,
// unless the segment is listed in allowed.
func HideDotfiles(next http.Handler, allowed []string) http.Handler {
	allow := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		allow[name] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, segment := range strings.Split(r.URL.Path, "/") {
			if strings.HasPrefix(segment, ".") && !allow[segment] {
				w.Header().Set("Cache-Control", "no-store")
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// AllowMethods returns 405 for any method other than GET and HEAD.
func AllowMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			w.Header().Set("Cache-Control", "no-store")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
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
