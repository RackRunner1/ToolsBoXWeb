package compress

import (
	"compress/gzip"
	"net/http"
	"strings"
)

var compressibleTypes = map[string]bool{
	"application/javascript": true,
	"application/json":       true,
	"application/xhtml+xml":  true,
	"application/xml":        true,
	"image/svg+xml":          true,
}

// Gzip compresses text responses for clients that accept gzip. Responses to
// range requests, HEAD requests, and already-encoded responses are left alone.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !acceptsGzip(r) || r.Method == http.MethodHead || r.Header.Get("Range") != "" {
			next.ServeHTTP(w, r)
			return
		}
		gw := &gzipWriter{ResponseWriter: w}
		next.ServeHTTP(gw, r)
		gw.close()
	})
}

func acceptsGzip(r *http.Request) bool {
	for _, part := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		name, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		if strings.EqualFold(strings.TrimSpace(name), "gzip") && !strings.Contains(strings.ReplaceAll(params, " ", ""), "q=0") {
			return true
		}
	}
	return false
}

type gzipWriter struct {
	http.ResponseWriter
	gz      *gzip.Writer
	started bool
}

func (w *gzipWriter) WriteHeader(status int) {
	if w.started {
		return
	}
	w.started = true
	if compressible(w.Header(), status) {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
		w.gz = gzip.NewWriter(w.ResponseWriter)
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipWriter) Write(p []byte) (int, error) {
	if !w.started {
		w.WriteHeader(http.StatusOK)
	}
	if w.gz != nil {
		return w.gz.Write(p)
	}
	return w.ResponseWriter.Write(p)
}

func (w *gzipWriter) close() {
	if w.gz != nil {
		w.gz.Close()
	}
}

func compressible(h http.Header, status int) bool {
	if status < http.StatusOK || status == http.StatusNoContent || status == http.StatusNotModified {
		return false
	}
	if h.Get("Content-Encoding") != "" {
		return false
	}
	mediaType, _, _ := strings.Cut(h.Get("Content-Type"), ";")
	mediaType = strings.ToLower(strings.TrimSpace(mediaType))
	return strings.HasPrefix(mediaType, "text/") || compressibleTypes[mediaType]
}
