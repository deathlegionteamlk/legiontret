package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deathlegionteam/legiontret/internal/config"
	"github.com/deathlegionteam/legiontret/internal/registry"
)

// Manager handles model operations
type Manager struct {
	cfg      *config.Config
	registry *registry.Registry
}

// NewManager creates a new model manager
func NewManager(cfg *config.Config, reg *registry.Registry) *Manager {
	return &Manager{
		cfg:      cfg,
		registry: reg,
	}
}

// LocalModel represents a locally stored model
type LocalModel struct {
	Name       string    `json:"name"`
	DisplayName string   `json:"display_name"`
	Family     string    `json:"family"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
	Parameters string    `json:"parameters"`
}

// ListLocal returns all locally available models
func (m *Manager) ListLocal() ([]LocalModel, error) {
	entries, err := os.ReadDir(m.cfg.ModelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []LocalModel
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gguf") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".gguf")
		displayName := name
		family := "unknown"
		params := ""

		// Look up in registry for richer info
		if regModel, ok := m.registry.GetModel(name); ok {
			displayName = regModel.DisplayName
			family = regModel.Family
			params = regModel.Parameters
		}

		models = append(models, LocalModel{
			Name:        name,
			DisplayName: displayName,
			Family:      family,
			Size:        info.Size(),
			ModifiedAt:  info.ModTime(),
			Parameters:  params,
		})
	}

	return models, nil
}

// IsDownloaded checks if a model is already downloaded
func (m *Manager) IsDownloaded(name string) bool {
	modelPath := m.cfg.ModelPath(name)
	_, err := os.Stat(modelPath)
	return err == nil
}

// GetModelPath returns the path to a model file
func (m *Manager) GetModelPath(name string) string {
	return m.cfg.ModelPath(name)
}

// Delete removes a locally stored model
func (m *Manager) Delete(name string) error {
	modelPath := m.cfg.ModelPath(name)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model %q not found", name)
	}

	if err := os.Remove(modelPath); err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	// Also remove any metadata
	metaPath := modelPath + ".json"
	os.Remove(metaPath)

	return nil
}

// ModelMetadata stores additional metadata about a downloaded model
type ModelMetadata struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	DownloadedAt time.Time `json:"downloaded_at"`
	Size        int64     `json:"size"`
	SHA256      string    `json:"sha256"`
	Family      string    `json:"family"`
}

// SaveMetadata saves model metadata
func (m *Manager) SaveMetadata(name string, meta *ModelMetadata) error {
	metaPath := m.cfg.ModelPath(name) + ".json"
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return os.WriteFile(metaPath, data, 0644)
}

// LoadMetadata loads model metadata
func (m *Manager) LoadMetadata(name string) (*ModelMetadata, error) {
	metaPath := m.cfg.ModelPath(name) + ".json"
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}
	var meta ModelMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &meta, nil
}

// GetModelInfo returns combined info about a model (registry + local)
func (m *Manager) GetModelInfo(name string) (*ModelInfo, error) {
	local := m.IsDownloaded(name)

	regModel, inRegistry := m.registry.GetModel(name)

	var size int64
	if local {
		if info, err := os.Stat(m.cfg.ModelPath(name)); err == nil {
			size = info.Size()
		}
	}

	var meta *ModelMetadata
	if local {
		meta, _ = m.LoadMetadata(name)
	}

	return &ModelInfo{
		Name:        name,
		DisplayName: regModel.DisplayName,
		Family:      regModel.Family,
		Parameters:  regModel.Parameters,
		Size:        size,
		IsDownloaded: local,
		InRegistry:  inRegistry,
		Description: regModel.Description,
		URL:         regModel.URL,
		Tags:        regModel.Tags,
		Metadata:    meta,
	}, nil
}

// ModelInfo provides comprehensive model information
type ModelInfo struct {
	Name         string          `json:"name"`
	DisplayName  string          `json:"display_name"`
	Family       string          `json:"family"`
	Parameters   string          `json:"parameters"`
	Size         int64           `json:"size"`
	IsDownloaded bool            `json:"is_downloaded"`
	InRegistry   bool            `json:"in_registry"`
	Description  string          `json:"description"`
	URL          string          `json:"url"`
	Tags         []string        `json:"tags"`
	Metadata     *ModelMetadata  `json:"metadata,omitempty"`
}

// ResolveModelName resolves a shorthand model name to a full name
func ResolveModelName(name string) string {
	// Handle common aliases
	aliases := map[string]string{
		"llama":     "llama3",
		"gemma":     "gemma3",
		"mistral":   "mistral",
		"qwen":      "qwen2.5",
		"deepseek":  "deepseek-r1",
		"phi":       "phi4",
		"coder":     "qwen2.5-coder:7b",
		"tiny":      "tinyllama",
		"mixtral":   "mixtral",
	}

	// Check exact alias
	if full, ok := aliases[strings.ToLower(name)]; ok {
		return full
	}

	// Handle tag syntax: model:tag -> model
	if idx := strings.Index(name, ":"); idx > 0 {
		return name
	}

	return name
}

// FindModelFile finds a GGUF model file in a directory (for extracted archives)
func FindModelFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gguf") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	// Check subdirectories
	for _, entry := range entries {
		if entry.IsDir() {
			path, err := FindModelFile(filepath.Join(dir, entry.Name()))
			if err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("no GGUF file found in %s", dir)
}
