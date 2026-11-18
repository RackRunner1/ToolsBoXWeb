package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNotFoundServesCustomPage(t *testing.T) {
	cfg := testConfig(t)
	if err := os.WriteFile(filepath.Join(cfg.Root, "404.html"), []byte("custom not found"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(cfg)

	for _, target := range []string{"/missing", "/missing/deep/", "/folder/"} {
		got := request(t, handler, target)
		if got.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d; want 404", target, got.Code)
		}
		if body := got.Body.String(); body != "custom not found" {
			t.Errorf("GET %s body = %q; want custom page", target, body)
		}
		if ct := got.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
			t.Errorf("GET %s Content-Type = %q; want text/html", target, ct)
		}
		if cc := got.Header().Get("Cache-Control"); cc != "no-store" {
			t.Errorf("GET %s Cache-Control = %q; want no-store", target, cc)
		}
	}

	got := requestMethod(t, handler, http.MethodHead, "/missing")
	if got.Code != http.StatusNotFound || got.Body.Len() != 0 {
		t.Errorf("HEAD /missing = %d with %d body bytes; want 404 without body", got.Code, got.Body.Len())
	}
}

func TestNotFoundFallsBackToDefaultPage(t *testing.T) {
	handler := NewHandler(testConfig(t))

	got := request(t, handler, "/missing")
	if got.Code != http.StatusNotFound || !strings.Contains(got.Body.String(), "Oops! Page not found.") {
		t.Errorf("GET /missing = %d %q; want built-in 404 page", got.Code, got.Body.String())
	}
}
