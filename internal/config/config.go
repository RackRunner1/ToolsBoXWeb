package config

import (
	"errors"
	"os"
)

// Config contains the server's environment-based configuration.
type Config struct {
	Host string
	Port string
	Root string
}

// Load reads configuration from the environment and validates the site root.
func Load() (Config, error) {
	cfg := Config{
		Host: envOrDefault("TBXM_HOST", "0.0.0.0"),
		Port: envOrDefault("TBXM_PORT", "8080"),
		Root: envOrDefault("TBXM_ROOT", "/var/www/tbxm"),
	}

	info, err := os.Stat(cfg.Root)
	if err != nil || !info.IsDir() {
		return Config{}, errors.New("TBXM_ROOT must point to an existing directory")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
