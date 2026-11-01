package cache

import (
	"net/http"
	"path"
	"strings"
)

// Middleware sets cache headers for HTML, assets, and error responses.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHTML := path.Ext(strings.ToLower(r.URL.Path)) == ".html" || strings.HasSuffix(r.URL.Path, "/")
		next.ServeHTTP(&responseWriter{ResponseWriter: w, isHTML: isHTML}, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	isHTML bool
	wrote  bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wrote {
		return
	}
	w.wrote = true
	switch {
	case status >= http.StatusBadRequest:
		w.Header().Set("Cache-Control", "no-store")
	case w.isHTML:
		w.Header().Set("Cache-Control", "no-cache")
	default:
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}
