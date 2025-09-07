package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config represents application configuration loaded from a YAML file.
type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	DSN string `yaml:"dsn"`
}

// ServerConfig holds HTTP server related settings.
type ServerConfig struct {
	Host           string   `yaml:"host"`
	Port           string   `yaml:"port"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// Load reads the configuration from the given path. If path is empty,
// it returns Config with default values.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{Host: "0.0.0.0", Port: "8080", AllowedOrigins: []string{"*"}},
		DB:     DBConfig{DSN: "file:mosscat.db?mode=memory&cache=shared"},
	}
	if path == "" {
		return cfg, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if len(cfg.Server.AllowedOrigins) == 0 {
		cfg.Server.AllowedOrigins = []string{"*"}
	}
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = "file:mosscat.db?mode=memory&cache=shared"
	}
	return cfg, nil
}
