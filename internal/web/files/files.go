package files

import (
	_ "embed"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strconv"
)

//go:embed default404.html
var defaultNotFoundPage []byte

// NewHandler serves files from root without exposing directory listings.
// Missing paths return root/404.html, or the built-in page when the site has
// none, with a 404 status.
func NewHandler(root string) http.Handler {
	fsys := noDirectoryListingFS{FileSystem: http.Dir(root)}
	fileServer := http.FileServer(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fsys.Open(path.Clean("/" + r.URL.Path))
		if errors.Is(err, fs.ErrNotExist) {
			serveNotFound(w, r, fsys)
			return
		}
		if err == nil {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}

func serveNotFound(w http.ResponseWriter, r *http.Request, fsys http.FileSystem) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if page, err := fsys.Open("/404.html"); err == nil {
		defer page.Close()
		if info, err := page.Stat(); err == nil && info.Mode().IsRegular() {
			w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
			w.WriteHeader(http.StatusNotFound)
			if r.Method != http.MethodHead {
				io.Copy(w, page)
			}
			return
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(defaultNotFoundPage)))
	w.WriteHeader(http.StatusNotFound)
	if r.Method != http.MethodHead {
		w.Write(defaultNotFoundPage)
	}
}

type noDirectoryListingFS struct {
	http.FileSystem
}

func (fsys noDirectoryListingFS) Open(name string) (http.File, error) {
	f, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if !info.IsDir() {
		return f, nil
	}

	index, err := fsys.FileSystem.Open(path.Join(name, "index.html"))
	if err != nil {
		f.Close()
		return nil, fs.ErrNotExist
	}
	indexInfo, statErr := index.Stat()
	index.Close()
	if statErr != nil || !indexInfo.Mode().IsRegular() {
		f.Close()
		return nil, fs.ErrNotExist
	}
	return f, nil
}
