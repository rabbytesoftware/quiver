package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed config.json
var embeddedConfig []byte

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Logger   LoggerConfig   `json:"logger"`
	Packages PackagesConfig `json:"packages"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// LoggerConfig holds logger-specific configuration
type LoggerConfig struct {
	Show      bool   `json:"show"`       // Whether to show logs on CLI
	Level     string `json:"level"`
	LogDir    string `json:"log_dir"`
	MaxSize   int    `json:"max_size"`
	MaxAge    int    `json:"max_age"`
	Compress  bool   `json:"compress"`
}

// PackagesConfig holds package-specific configuration
type PackagesConfig struct {
	Repositories []string `json:"repositories"`  // List of repositories (URLs or local directories)
	InstallDir   string   `json:"install_dir"`   // Directory to install packages	
	DatabasePath string   `json:"database_path"` // Local package database path
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Logger: LoggerConfig{
			Show:     false,
			Level:    "info",
			LogDir:   "./logs",
			MaxSize:  100, // megabytes
			MaxAge:   7,   // days
			Compress: true,
		},
		Packages: PackagesConfig{
			Repositories: []string{"./pkgs"},
			InstallDir:   "./pkgs",
			DatabasePath: "./pkgs/packages.db",
		},
	}
}

// Load loads configuration from embedded data or environment variables
func Load() (*Config, error) {
	cfg := Default()

	// Load from embedded config first
	if err := loadFromEmbedded(cfg); err != nil {
		return nil, fmt.Errorf("failed to load embedded config: %w", err)
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromEmbedded loads configuration from embedded JSON data
func loadFromEmbedded(cfg *Config) error {
	return json.Unmarshal(embeddedConfig, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	if port := os.Getenv("QUIVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			cfg.Server.Port = p
		}
	}

	if host := os.Getenv("QUIVER_HOST"); host != "" {
		cfg.Server.Host = host
	}

	if level := os.Getenv("QUIVER_LOG_LEVEL"); level != "" {
		cfg.Logger.Level = strings.ToLower(level)
	}

	if logDir := os.Getenv("QUIVER_LOG_DIR"); logDir != "" {
		cfg.Logger.LogDir = logDir
	}
}

// parsePort parses port from string
func parsePort(portStr string) (int, error) {
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid port number: %d", port)
	}
	return port, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
} 