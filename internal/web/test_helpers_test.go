package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"tbxm/internal/config"
)

func makeSite(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"index.html":               "home page",
		"page.html":                "clean page",
		"fooindex.html":            "not an index file",
		"assets/site.css":          "body { color: black }",
		"folder/secret.txt":        "must not be listed",
		".env":                     "secret",
		".env.local":               "secret",
		".git/config":              "git data",
		".DS_Store":                "finder data",
		".nojekyll":                "marker",
		".well-known/security.txt": "contact",
		"docs/index.html":          "docs page",
	}
	for name, content := range files {
		file := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func testConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{Root: makeSite(t), AllowWellKnown: true}
}

func request(t *testing.T, handler http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	return requestMethod(t, handler, http.MethodGet, target)
}

func requestMethod(t *testing.T, handler http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)
	handler.ServeHTTP(rec, req)
	return rec
}
