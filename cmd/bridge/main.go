package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/wolf-infra/gatus-api-bridge/internal/api"
	"github.com/wolf-infra/gatus-api-bridge/internal/config"
	"github.com/wolf-infra/gatus-api-bridge/internal/gatus"
)

var Version = "dev" // Defaults to "dev" if not built with a tag

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger) // Good practice to set it globally just in case

	// Load Configuration
	cfg := config.Load()

	// Initialize Manager and Server (Injecting the logger)
	manager := gatus.NewManager(cfg.ConfigPath, cfg.DryRun, logger)
	server := api.NewServer(manager, logger, Version)
	mux := server.Mount()

	// Structured Startup Logs
	logger.Info("Wolf-Infra Gatus Bridge starting",
		slog.String("port", cfg.Port),
		slog.String("config_path", cfg.ConfigPath),
		slog.Bool("dry_run", cfg.DryRun),
	)

	if cfg.DryRun {
		logger.Warn("DRY_RUN is enabled. Config file will not be modified.")
	}

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
