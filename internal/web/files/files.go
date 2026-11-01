package files

import (
	"io/fs"
	"net/http"
	"path"
)

// NewHandler serves files from root without exposing directory listings.
func NewHandler(root string) http.Handler {
	return http.FileServer(noDirectoryListingFS{FileSystem: http.Dir(root)})
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
