package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppName    = "legiontret"
	Version    = "1.0.0"
	OrgName    = "Death Legion Team"
	DefaultHost = "127.0.0.1"
	DefaultPort = 11434
)

// Config holds the application configuration
type Config struct {
	ModelsDir   string
	BinariesDir string
	Host        string
	Port        int
	Debug       bool
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	baseDir := filepath.Join(homeDir, ".legiontret")

	return &Config{
		ModelsDir:   filepath.Join(baseDir, "models"),
		BinariesDir: filepath.Join(baseDir, "bin"),
		Host:        DefaultHost,
		Port:        DefaultPort,
		Debug:       false,
	}
}

// EnsureDirs creates all necessary directories
func (c *Config) EnsureDirs() error {
	dirs := []string{c.ModelsDir, c.BinariesDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// APIBaseURL returns the base URL for the API server
func (c *Config) APIBaseURL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

// LlamaCppBinaryPath returns the path to the llama.cpp server binary
func (c *Config) LlamaCppBinaryPath() string {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return filepath.Join(c.BinariesDir, "llama-server"+ext)
}

// ModelPath returns the full path for a model file
func (c *Config) ModelPath(modelName string) string {
	return filepath.Join(c.ModelsDir, modelName+".gguf")
}

// LoadConfig loads or creates configuration
func LoadConfig() *Config {
	cfg := DefaultConfig()
	
	configFile := cfg.ConfigPath()
	data, err := os.ReadFile(configFile)
	if err == nil {
		// Parse config file if it exists
		_ = data // TODO: parse JSON/YAML config
	}

	return cfg
}

// ConfigPath returns the path to the config file
func (c *Config) ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".legiontret", "config.json")
}
