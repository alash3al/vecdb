package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		ListenAddr string `yaml:"listen"`
	} `yaml:"server"`

	Store struct {
		Driver string         `yaml:"driver"`
		Args   map[string]any `yaml:"args"`
	} `yaml:"store"`

	Embedder struct {
		Enabled bool           `yaml:"enabled"`
		Driver  string         `yaml:"driver"`
		Args    map[string]any `yaml:"args"`
	} `yaml:"embedder"`
}

func NewFromFile(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content = []byte(os.ExpandEnv(string(content)))

	var cfg Config

	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
