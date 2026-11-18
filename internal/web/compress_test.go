package web

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func getWithHeaders(t *testing.T, handler http.Handler, target string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func gunzip(t *testing.T, body []byte) string {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gunzip: %v", err)
	}
	return string(out)
}

func TestGzipCompressesTextResponses(t *testing.T) {
	cfg := testConfig(t)
	cfg.Compression = true
	handler := NewHandler(cfg)

	got := getWithHeaders(t, handler, "/assets/site.css", map[string]string{"Accept-Encoding": "gzip, br"})
	if got.Code != http.StatusOK {
		t.Fatalf("GET /assets/site.css = %d; want 200", got.Code)
	}
	if enc := got.Header().Get("Content-Encoding"); enc != "gzip" {
		t.Errorf("Content-Encoding = %q; want gzip", enc)
	}
	if vary := got.Header().Get("Vary"); vary != "Accept-Encoding" {
		t.Errorf("Vary = %q; want Accept-Encoding", vary)
	}
	if body := gunzip(t, got.Body.Bytes()); body != "body { color: black }" {
		t.Errorf("decompressed body = %q; want original css", body)
	}
}

func TestGzipSkipsWhenNotAccepted(t *testing.T) {
	cfg := testConfig(t)
	cfg.Compression = true
	handler := NewHandler(cfg)

	got := getWithHeaders(t, handler, "/assets/site.css", nil)
	if enc := got.Header().Get("Content-Encoding"); enc != "" {
		t.Errorf("Content-Encoding = %q; want none without Accept-Encoding", enc)
	}
	if got.Body.String() != "body { color: black }" {
		t.Errorf("body = %q; want plain css", got.Body.String())
	}
}

func TestGzipSkipsNonTextAndRangeRequests(t *testing.T) {
	cfg := testConfig(t)
	cfg.Compression = true
	if err := os.WriteFile(filepath.Join(cfg.Root, "image.png"), []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(cfg)

	got := getWithHeaders(t, handler, "/image.png", map[string]string{"Accept-Encoding": "gzip"})
	if enc := got.Header().Get("Content-Encoding"); enc != "" {
		t.Errorf("image Content-Encoding = %q; want none", enc)
	}

	got = getWithHeaders(t, handler, "/assets/site.css", map[string]string{"Accept-Encoding": "gzip", "Range": "bytes=0-3"})
	if enc := got.Header().Get("Content-Encoding"); enc != "" {
		t.Errorf("range Content-Encoding = %q; want none", enc)
	}
	if got.Code != http.StatusPartialContent || got.Body.String() != "body" {
		t.Errorf("range GET = %d %q; want 206 %q", got.Code, got.Body.String(), "body")
	}
}

func TestGzipDisabled(t *testing.T) {
	cfg := testConfig(t)
	cfg.Compression = false
	handler := NewHandler(cfg)

	got := getWithHeaders(t, handler, "/assets/site.css", map[string]string{"Accept-Encoding": "gzip"})
	if enc := got.Header().Get("Content-Encoding"); enc != "" {
		t.Errorf("Content-Encoding = %q; want none when compression is disabled", enc)
	}
}
