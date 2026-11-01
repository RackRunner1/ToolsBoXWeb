package web

import "testing"

func TestCachePolicy(t *testing.T) {
	handler := NewHandler(makeSite(t))
	for _, tc := range []struct {
		path  string
		cache string
	}{
		{path: "/", cache: "no-cache"},
		{path: "/page", cache: "no-cache"},
		{path: "/assets/site.css", cache: "public, max-age=3600"},
		{path: "/missing", cache: "no-store"},
	} {
		got := request(t, handler, tc.path)
		if value := got.Header().Get("Cache-Control"); value != tc.cache {
			t.Errorf("GET %s Cache-Control = %q; want %q", tc.path, value, tc.cache)
		}
	}
}
