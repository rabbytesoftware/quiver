package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
	Level     string `json:"level"`
	LogDir    string `json:"log_dir"`
	MaxSize   int    `json:"max_size"`
	MaxAge    int    `json:"max_age"`
	Compress  bool   `json:"compress"`
}

// PackagesConfig holds package-specific configuration
type PackagesConfig struct {
	Repository    string   `json:"repository"`
	TemplateDir  string   `json:"template_dir"`
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
			Level:    "info",
			LogDir:   "./logs",
			MaxSize:  100, // megabytes
			MaxAge:   7,   // days
			Compress: true,
		},
		Packages: PackagesConfig{
			Repository:    "./pkgs",
			TemplateDir:  "./template",
		},
	}
}

// Load loads configuration from file or environment variables
func Load() (*Config, error) {
	cfg := Default()

	// Try to load from config file
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		if err := loadFromFile(cfg, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	if path := os.Getenv("QUIVER_CONFIG"); path != "" {
		return path
	}
	return "./config.json"
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
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

	if pkgDir := os.Getenv("QUIVER_PACKAGES_DIR"); pkgDir != "" {
		cfg.Packages.Repository = pkgDir
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