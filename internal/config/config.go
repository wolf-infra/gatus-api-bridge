package config

import (
	"os"
	"strings"
)

type Config struct {
	Port       string
	ConfigPath string
	DryRun     bool
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

	// DRY_RUN is true if the env var is set to "true" or "1"
	dryRunEnv := strings.ToLower(os.Getenv("DRY_RUN"))
	dryRun := dryRunEnv == "true" || dryRunEnv == "1"

	return &Config{
		Port:       port,
		ConfigPath: path,
		DryRun:     dryRun,
	}
}
