package web

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestRequestsAreNotLogged(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	handler := NewHandler(testConfig(t))
	for _, target := range []string{"/", "/missing", "/.env", "/assets/site.css"} {
		request(t, handler, target)
	}
	if got := logs.String(); got != "" {
		t.Fatalf("HTTP requests produced log output: %q", got)
	}
}
