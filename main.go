package main

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func hideGit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, segment := range strings.Split(r.URL.Path, "/") {
			if segment == ".git" {
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func hideEnv(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, segment := range strings.Split(r.URL.Path, "/") {
			name := strings.ToLower(strings.TrimLeft(segment, "."))
			if name == "env" || strings.HasPrefix(name, "env.") || strings.HasPrefix(name, "env-") || strings.HasPrefix(name, "env_") {
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func cleanURLs(root string, next http.Handler) http.Handler {
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
			if strings.HasSuffix(p, "index.html") && target != "/" {
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

func main() {
	host := getenv("TBXM_HOST", "0.0.0.0")
	port := getenv("TBXM_PORT", "8080")
	root := getenv("TBXM_ROOT", "/var/www/tbxm")

	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           hideEnv(hideGit(cleanURLs(root, http.FileServer(http.Dir(root))))),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
