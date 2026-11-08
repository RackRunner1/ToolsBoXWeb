package web

import (
	"net/http"
	"strings"
	"testing"
)

func TestCleanURLsServeAndRedirect(t *testing.T) {
	handler := NewHandler(testConfig(t))

	for _, tc := range []struct {
		name string
		path string
		body string
	}{
		{name: "root index", path: "/", body: "home page"},
		{name: "clean URL", path: "/page", body: "clean page"},
		{name: "directory index", path: "/docs/", body: "docs page"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := request(t, handler, tc.path)
			if got.Code != http.StatusOK || !strings.Contains(got.Body.String(), tc.body) {
				t.Fatalf("GET %s = %d %q; want body containing %q", tc.path, got.Code, got.Body.String(), tc.body)
			}
		})
	}

	for _, tc := range []struct {
		path     string
		location string
	}{
		{path: "/page.html?lang=fr", location: "/page?lang=fr"},
		{path: "/index.html", location: "/"},
		{path: "/docs/index.html", location: "/docs/"},
		{path: "/fooindex.html", location: "/fooindex"},
	} {
		got := request(t, handler, tc.path)
		if got.Code != http.StatusMovedPermanently || got.Header().Get("Location") != tc.location {
			t.Errorf("GET %s = %d Location %q; want %d Location %q", tc.path, got.Code, got.Header().Get("Location"), http.StatusMovedPermanently, tc.location)
		}
	}
}
