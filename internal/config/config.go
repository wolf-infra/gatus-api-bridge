package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Port       string
	ConfigPath string
	DryRun     bool
	APIKey     string
	LogLevel   string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	path := os.Getenv("GATUS_CONFIG_PATH")
	if path == "" {
		path = "/data/config.yaml"
	}

	dryRunEnv := strings.ToLower(os.Getenv("DRY_RUN"))
	dryRun := dryRunEnv == "true" || dryRunEnv == "1"

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "info"
	}

	apiKey := os.Getenv("API_KEY")

	cfg := &Config{
		Port:       port,
		ConfigPath: path,
		DryRun:     dryRun,
		APIKey:     apiKey,
		LogLevel:   logLevel,
	}

	slog.Debug("Bridge Configuration Loaded",
		"port", cfg.Port,
		"config_path", cfg.ConfigPath,
		"dry_run", cfg.DryRun,
		"log_level", cfg.LogLevel,
		"api_key_length", len(cfg.APIKey),
	)

	return cfg
}
