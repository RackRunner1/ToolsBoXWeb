package main

import (
	"log/slog"
	"os"

	"tbxm/internal/config"
	"tbxm/internal/server"
	"tbxm/internal/web"
)

func main() {
	cfg, err := config.Load()
	if err == nil {
		err = server.Run(cfg, web.NewHandler(cfg.Root))
	}
	if err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
