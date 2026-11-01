package web

import (
	"net/http"

	"tbxm/internal/web/cache"
	"tbxm/internal/web/files"
	"tbxm/internal/web/security"
	"tbxm/internal/web/urls"
)

// NewHandler builds the static-site handler and its middleware chain.
func NewHandler(root string) http.Handler {
	handler := files.NewHandler(root)
	handler = cache.Middleware(handler)
	handler = urls.Clean(root, handler)
	handler = security.HideGit(handler)
	handler = security.HideEnv(handler)
	handler = security.Headers(handler)
	return handler
}
