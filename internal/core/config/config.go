package config

import (
	_ "embed"

	"os"

	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed default.yaml
	defaultConfigByte []byte
	config *Config
)

type Netbridge struct {
	Enabled      bool   `yaml:"enabled"`
	AllowedPorts string `yaml:"allowed_ports"`
}

type Arrows struct {
	Repositories []string `yaml:"repositories"`
	InstallDir   string   `yaml:"install_dir"`
}

type API struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

type Database struct {
	Path string `yaml:"path"`
}

type Watcher struct {
	Enabled  bool   `yaml:"enabled"`
	Level    string `yaml:"level"`
	Folder   string `yaml:"folder"`
	MaxSize  int    `yaml:"max_size"`
	MaxAge   int    `yaml:"max_age"`
	Compress bool   `yaml:"compress"`
}

type ConfigData struct {
	Netbridge Netbridge `yaml:"netbridge"`
	Arrows    Arrows    `yaml:"arrows"`
	API       API       `yaml:"api"`
	Database  Database  `yaml:"database"`
	Watcher   Watcher   `yaml:"watcher"`
}

type Config struct {
	Config ConfigData `yaml:"config"`
}

func Get() *Config {
	if config != nil {
		return config
	}

	configPath := metadata.GetDefaultConfigPath()
	
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		config = getDefaultConfig()
		return config
	}

	config = &Config{}
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		config = getDefaultConfig()
		return config
	}

	return config
}

func GetNetbridge() Netbridge {
	return Get().Config.Netbridge
}

func GetArrows() Arrows {
	return Get().Config.Arrows
}

func GetAPI() API {
	return Get().Config.API
}

func GetDatabase() Database {
	return Get().Config.Database
}

func GetWatcher() Watcher {
	return Get().Config.Watcher
}

func GetConfigPath() string {
	return metadata.GetDefaultConfigPath()
}

func ConfigExists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

func getDefaultConfig() *Config {
	config = &Config{}
	err := yaml.Unmarshal(defaultConfigByte, config)
	if err == nil {
		return config
	}
	
	return &Config{
		Config: ConfigData{
			Netbridge: Netbridge{
				Enabled:      true,
				AllowedPorts: "40128-40256",
			},
			Arrows: Arrows{
				Repositories: []string{
					"./pkgs",
				},
				InstallDir: "./arrows",
			},
			API: API{
				Host:    "0.0.0.0",
				Port:    40257,
			},
			Database: Database{
				Path: "./.db",
			},
			Watcher: Watcher{
				Enabled:  true,
				Level:    "info",
				Folder:   "./logs",
				MaxSize:  100,
				MaxAge:   7,
				Compress: true,
			},
		},
	}
}