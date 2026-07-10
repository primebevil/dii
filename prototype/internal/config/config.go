package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is one node's settings. One binary is a node; two nodes is the same
// binary run twice with different config files (config-a.yaml, config-b.yaml).
type Config struct {
	Listen        string   `yaml:"listen"`         // OpenAI endpoint listen addr, e.g. ":8080"
	NodeID        string   `yaml:"node_id"`        // label for this node, e.g. "A"
	ModelServer   string   `yaml:"model_server"`   // OpenAI backend base URL (Ollama). Unused in M1 (mock backend).
	ConsumerToken string   `yaml:"consumer_token"` // stub bearer token that selects the consumer door
	Peers         []string `yaml:"peers"`          // peer node endpoints, e.g. http://host:8080
}

// Load reads and validates a node's YAML config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if c.Listen == "" {
		c.Listen = ":8080"
	}
	if c.NodeID == "" {
		return nil, fmt.Errorf("config %s: node_id is required", path)
	}
	return &c, nil
}
