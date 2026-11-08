package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config contains the server's environment-based configuration.
type Config struct {
	Host string
	Port string
	Root string

	// AllowWellKnown serves the /.well-known directory when true.
	AllowWellKnown bool
	// AllowDotfiles lists additional dot-prefixed names to serve, such as ".nojekyll".
	AllowDotfiles []string
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

	cfg.AllowWellKnown, err = boolOrDefault("TBXM_ALLOW_WELL_KNOWN", true)
	if err != nil {
		return Config{}, err
	}
	cfg.AllowDotfiles, err = dotfilesFromEnv("TBXM_ALLOW_DOTFILES")
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func boolOrDefault(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", key)
	}
	return parsed, nil
}

func dotfilesFromEnv(key string) ([]string, error) {
	var names []string
	for _, name := range strings.Split(os.Getenv(key), ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if !strings.HasPrefix(name, ".") || strings.ContainsAny(name, "/\\") || name == "." || name == ".." {
			return nil, fmt.Errorf("%s entries must be dot-prefixed names without slashes, got %q", key, name)
		}
		names = append(names, name)
	}
	return names, nil
}
