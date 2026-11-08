package web

import (
	"net/http"

	"tbxm/internal/config"
	"tbxm/internal/web/cache"
	"tbxm/internal/web/files"
	"tbxm/internal/web/security"
	"tbxm/internal/web/urls"
)

// NewHandler builds the static-site handler and its middleware chain.
func NewHandler(cfg config.Config) http.Handler {
	allowedDotfiles := append([]string(nil), cfg.AllowDotfiles...)
	if cfg.AllowWellKnown {
		allowedDotfiles = append(allowedDotfiles, ".well-known")
	}

	handler := files.NewHandler(cfg.Root)
	handler = cache.Middleware(handler)
	handler = urls.Clean(cfg.Root, handler)
	handler = security.HideDotfiles(handler, allowedDotfiles)
	handler = security.HideEnv(handler)
	handler = security.HideGit(handler)
	handler = security.AllowMethods(handler)
	handler = security.Headers(handler)
	return handler
}
