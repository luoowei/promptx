package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Provider represents an AI provider configuration
type Provider struct {
	Name    string `json:"name"`
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
	Model   string `json:"model"`
}

// Config holds all configuration for PromptX
type Config struct {
	DefaultProvider string              `json:"default_provider"`
	Providers       map[string]Provider `json:"providers"`
	Theme           string              `json:"theme"`
	MaxTokens       int                 `json:"max_tokens"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultProvider: "openai",
		Theme:           "dracula",
		MaxTokens:       4096,
		Providers: map[string]Provider{
			"openai": {
				Name:    "openai",
				BaseURL: "https://api.openai.com/v1",
				Model:   "gpt-4o",
			},
			"anthropic": {
				Name:    "anthropic",
				BaseURL: "https://api.anthropic.com/v1",
				Model:   "claude-sonnet-4-6",
			},
			"ollama": {
				Name:    "ollama",
				BaseURL: "http://localhost:11434/v1",
				Model:   "llama3.2",
			},
		},
	}
}

// ConfigDir returns the config directory path
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".promptx")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config directory: %w", err)
	}
	return dir, nil
}

// ConfigPath returns the config file path
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads configuration from disk or returns defaults
func Load() (*Config, error) {
	cfg := DefaultConfig()
	path, err := ConfigPath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // return defaults
		}
		return cfg, fmt.Errorf("cannot read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("cannot parse config: %w", err)
	}

	// Merge with defaults for any missing providers
	defaults := DefaultConfig()
	for k, v := range defaults.Providers {
		if _, ok := cfg.Providers[k]; !ok {
			cfg.Providers[k] = v
		}
	}

	return cfg, nil
}

// Save writes configuration to disk
func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}

	return nil
}

// GetAPIKey retrieves API key from environment or config
func GetAPIKey(provider string) string {
	// Check environment variable first
	envMap := map[string]string{
		"openai":    "OPENAI_API_KEY",
		"anthropic": "ANTHROPIC_API_KEY",
		"ollama":    "OLLAMA_API_KEY",
	}

	if envVar, ok := envMap[provider]; ok {
		if key := os.Getenv(envVar); key != "" {
			return key
		}
	}

	// Fall back to config
	cfg, err := Load()
	if err != nil {
		return ""
	}
	if p, ok := cfg.Providers[provider]; ok {
		return p.APIKey
	}
	return ""
}
