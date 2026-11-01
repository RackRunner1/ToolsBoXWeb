package urls

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Clean serves extensionless paths from matching .html files and redirects
// explicit .html URLs to their canonical clean URLs.
func Clean(root string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "/" || strings.HasSuffix(p, "/") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasSuffix(p, ".html") {
			target := strings.TrimSuffix(p, ".html")
			if path.Base(p) == "index.html" {
				target = strings.TrimSuffix(p, "index.html")
			}
			target = path.Clean(target)
			if path.Base(p) == "index.html" && target != "/" {
				target += "/"
			}
			u := *r.URL
			u.Path = target
			u.RawPath = ""
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			return
		}

		fp := filepath.Join(root, filepath.FromSlash(path.Clean(p)))
		if _, err := os.Stat(fp); errors.Is(err, fs.ErrNotExist) {
			if info, err := os.Stat(fp + ".html"); err == nil && info.Mode().IsRegular() {
				r.URL.Path = path.Clean(p) + ".html"
				r.URL.RawPath = ""
			}
		}
		next.ServeHTTP(w, r)
	})
}
