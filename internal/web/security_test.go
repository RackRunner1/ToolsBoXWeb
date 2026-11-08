package web

import (
	"net/http"
	"strings"
	"testing"
)

func TestRestrictedPathsAndDirectoryListing(t *testing.T) {
	handler := NewHandler(testConfig(t))
	for _, target := range []string{"/.env", "/.env.local", "/.git/config", "/.DS_Store", "/folder/"} {
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

func TestDotfilesAllowList(t *testing.T) {
	cfg := testConfig(t)
	cfg.AllowDotfiles = []string{".nojekyll"}
	handler := NewHandler(cfg)

	if got := request(t, handler, "/.nojekyll"); got.Code != http.StatusOK {
		t.Errorf("GET /.nojekyll = %d; want 200", got.Code)
	}
	if got := request(t, handler, "/.well-known/security.txt"); got.Code != http.StatusOK {
		t.Errorf("GET /.well-known/security.txt = %d; want 200", got.Code)
	}
	if got := request(t, handler, "/.DS_Store"); got.Code != http.StatusNotFound {
		t.Errorf("GET /.DS_Store = %d; want 404", got.Code)
	}
}

func TestWellKnownCanBeDisabled(t *testing.T) {
	cfg := testConfig(t)
	cfg.AllowWellKnown = false
	handler := NewHandler(cfg)

	if got := request(t, handler, "/.well-known/security.txt"); got.Code != http.StatusNotFound {
		t.Errorf("GET /.well-known/security.txt = %d; want 404", got.Code)
	}
}

func TestHardBlockedPathsCannotBeAllowed(t *testing.T) {
	cfg := testConfig(t)
	cfg.AllowDotfiles = []string{".env", ".git"}
	handler := NewHandler(cfg)

	for _, target := range []string{"/.env", "/.git/config"} {
		if got := request(t, handler, target); got.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d; want 404", target, got.Code)
		}
	}
}

func TestOnlyGetAndHeadAreAllowed(t *testing.T) {
	handler := NewHandler(testConfig(t))
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		rec := requestMethod(t, handler, method, "/")
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s / = %d; want 405", method, rec.Code)
		}
		if rec.Header().Get("Allow") != "GET, HEAD" {
			t.Errorf("%s / Allow = %q; want %q", method, rec.Header().Get("Allow"), "GET, HEAD")
		}
	}
	if got := requestMethod(t, handler, http.MethodHead, "/"); got.Code != http.StatusOK {
		t.Errorf("HEAD / = %d; want 200", got.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	got := request(t, NewHandler(testConfig(t)), "/")
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
