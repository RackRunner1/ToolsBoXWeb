package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func makeSite(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"index.html":        "home page",
		"page.html":         "clean page",
		"fooindex.html":     "not an index file",
		"assets/site.css":   "body { color: black }",
		"folder/secret.txt": "must not be listed",
		".env":              "secret",
		".env.local":        "secret",
		".git/config":       "git data",
		"docs/index.html":   "docs page",
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

func request(t *testing.T, handler http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	handler.ServeHTTP(rec, req)
	return rec
}
