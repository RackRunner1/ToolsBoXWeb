package web

import (
	"net/http"
	"strings"
	"testing"
)

func TestRestrictedPathsAndDirectoryListing(t *testing.T) {
	handler := NewHandler(makeSite(t))
	for _, target := range []string{"/.env", "/.env.local", "/.git/config", "/folder/"} {
		got := request(t, handler, target)
		if got.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d; want 404", target, got.Code)
		}
		if strings.Contains(got.Body.String(), "must not be listed") {
			t.Errorf("GET %s exposed a directory entry", target)
		}
		if target != "/folder/" && got.Header().Get("Cache-Control") != "no-store" {
			t.Errorf("GET %s Cache-Control = %q; want no-store", target, got.Header().Get("Cache-Control"))
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	got := request(t, NewHandler(makeSite(t)), "/")
	for name, want := range map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "SAMEORIGIN",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	} {
		if got.Header().Get(name) != want {
			t.Errorf("%s = %q; want %q", name, got.Header().Get(name), want)
		}
	}
}
