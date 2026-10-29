package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
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
		Handler:           http.FileServer(http.Dir(root)),
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
