package gatus

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Manager struct {
	configPath string
	dryRun     bool
	mu         sync.Mutex
	logger     *slog.Logger
}

func NewManager(path string, dryRun bool, logger *slog.Logger) *Manager {
	return &Manager{
		configPath: path,
		dryRun:     dryRun,
		logger:     logger,
	}
}

// GetEndpoints returns all endpoints, optionally filtered by a specific group
func (m *Manager) GetEndpoints(groupFilter string) ([]Endpoint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var cfg Config
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Info("Config file does not exist yet, returning empty list", slog.String("path", m.configPath))
			return []Endpoint{}, nil
		}
		m.logger.Error("Failed to read config file", slog.Any("error", err), slog.String("path", m.configPath))
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		m.logger.Error("Failed to parse YAML", slog.Any("error", err))
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	if groupFilter == "" {
		return cfg.Endpoints, nil
	}

	// Filter by group
	var filtered []Endpoint
	for _, ep := range cfg.Endpoints {
		if ep.Group == groupFilter {
			filtered = append(filtered, ep)
		}
	}
	return filtered, nil
}

// AddEndpoint safely appends a new endpoint. Respects DRY_RUN flag.
func (m *Manager) AddEndpoint(newEp Endpoint) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var cfg Config
	data, err := os.ReadFile(m.configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			// We log a warning instead of failing, because we might just want to overwrite a broken file
			m.logger.Warn("Failed to unmarshal existing config file", slog.Any("error", err))
		}
	}

	for _, ep := range cfg.Endpoints {
		if ep.Name == newEp.Name && ep.Group == newEp.Group {
			return false, nil // Already exists
		}
	}

	cfg.Endpoints = append(cfg.Endpoints, newEp)
	updatedData, err := yaml.Marshal(&cfg)
	if err != nil {
		m.logger.Error("Failed to marshal YAML", slog.Any("error", err))
		return false, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	if m.dryRun {
		m.logger.Info("DRY_RUN: Would have written endpoint",
			slog.String("name", newEp.Name),
			slog.String("group", newEp.Group),
		)
		return true, nil
	}

	if err := os.WriteFile(m.configPath, updatedData, 0644); err != nil {
		m.logger.Error("Failed to write to config file", slog.Any("error", err), slog.String("path", m.configPath))
		return false, fmt.Errorf("failed to write file: %w", err)
	}

	m.logger.Info("Successfully wrote to file", slog.String("path", m.configPath))
	return true, nil
}

// DeleteEndpoint safely removes an endpoint. Respects DRY_RUN flag.
func (m *Manager) DeleteEndpoint(name, group string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var cfg Config
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // File doesn't exist, nothing to delete
		}
		return false, fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return false, fmt.Errorf("failed to parse yaml: %w", err)
	}

	initialLen := len(cfg.Endpoints)
	var newEndpoints []Endpoint

	// Keep everything EXCEPT the one we want to delete
	for _, ep := range cfg.Endpoints {
		if ep.Name == name && ep.Group == group {
			continue
		}
		newEndpoints = append(newEndpoints, ep)
	}

	// If the length is the same, we didn't find it
	if len(newEndpoints) == initialLen {
		return false, nil
	}

	cfg.Endpoints = newEndpoints
	updatedData, err := yaml.Marshal(&cfg)
	if err != nil {
		m.logger.Error("Failed to marshal YAML", slog.Any("error", err))
		return false, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	if m.dryRun {
		m.logger.Info("DRY_RUN: Would have deleted endpoint", slog.String("name", name))
		return true, nil
	}

	if err := os.WriteFile(m.configPath, updatedData, 0644); err != nil {
		m.logger.Error("Failed to write config file", slog.Any("error", err))
		return false, fmt.Errorf("failed to write file: %w", err)
	}

	m.logger.Info("Successfully deleted endpoint from file", slog.String("path", m.configPath))
	return true, nil
}
