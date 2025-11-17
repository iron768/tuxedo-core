package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the server configuration
type Config struct {
	Server  ServerConfig  `json:"server"`
	Project ProjectConfig `json:"project"`
	Logging LoggingConfig `json:"logging"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port         string `json:"port"`
	Host         string `json:"host"`
	AllowOrigins string `json:"allowOrigins"`
}

// ProjectConfig holds project path settings
type ProjectConfig struct {
	YukonPath   string `json:"yukonPath"`
	ScenesPath  string `json:"scenesPath"`
	AssetsPath  string `json:"assetsPath"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
	Format  string `json:"format"`
}

var defaultConfig = Config{
	Server: ServerConfig{
		Port:         "3000",
		Host:         "0.0.0.0",
		AllowOrigins: "*",
	},
	Project: ProjectConfig{
		YukonPath:  "../yukon",
		ScenesPath: "src/scenes",
		AssetsPath: "assets",
	},
	Logging: LoggingConfig{
		Enabled: true,
		Level:   "info",
		Format:  "json",
	},
}

// Load loads configuration from file or returns default
func Load(configPath string) (*Config, error) {
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &defaultConfig, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save(configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// GetAssetsPath returns the full path to assets directory
func (c *Config) GetAssetsPath() string {
	return filepath.Join(c.Project.YukonPath, c.Project.AssetsPath)
}

// GetScenesPath returns the full path to scenes directory
func (c *Config) GetScenesPath() string {
	return filepath.Join(c.Project.YukonPath, c.Project.ScenesPath)
}
