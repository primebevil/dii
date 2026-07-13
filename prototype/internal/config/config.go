package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Defaults applied when a field is omitted.
const (
	defaultListen                = ":8080"
	defaultModelServer           = "http://localhost:11434/v1"
	defaultStartupTimeout        = 5 * time.Second
	defaultResponseHeaderTimeout = 30 * time.Second
	defaultMaxBodyBytes          = 1 << 20 // 1 MiB
)

// Config is one node's resolved settings (defaults applied, durations parsed).
// One binary is a node; two nodes is the same binary run twice with different
// config files.
type Config struct {
	Listen        string   // OpenAI endpoint listen addr, e.g. ":8080"
	NodeID        string   // label for this node, e.g. "A"
	ModelServer   string   // OpenAI backend base URL (Ollama's /v1)
	ConsumerToken string   // stub bearer token that selects the consumer door
	Advertise     string   // this node's own URL, published in its manifest
	Peers         []string // peer node endpoints

	StartupTimeout        time.Duration // boot-time model-list and peer-manifest fetches
	ResponseHeaderTimeout time.Duration // max wait for first byte from an upstream (0 = no limit)
	MaxBodyBytes          int64         // request body cap
}

// rawConfig mirrors the YAML on disk before defaulting and duration parsing.
type rawConfig struct {
	Listen        string   `yaml:"listen"`
	NodeID        string   `yaml:"node_id"`
	ModelServer   string   `yaml:"model_server"`
	ConsumerToken string   `yaml:"consumer_token"`
	Advertise     string   `yaml:"advertise"`
	Peers         []string `yaml:"peers"`

	StartupTimeout        string `yaml:"startup_timeout"`         // Go duration, e.g. "5s"
	ResponseHeaderTimeout string `yaml:"response_header_timeout"` // Go duration, e.g. "30s"; "0s" = no limit
	MaxBodyBytes          int64  `yaml:"max_body_bytes"`
}

// Load reads a node's YAML config, applies defaults, and validates it.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if raw.NodeID == "" {
		return nil, fmt.Errorf("config %s: node_id is required", path)
	}

	startupTimeout, err := parseDuration(raw.StartupTimeout, defaultStartupTimeout)
	if err != nil {
		return nil, fmt.Errorf("config %s: startup_timeout: %w", path, err)
	}
	responseHeaderTimeout, err := parseDuration(raw.ResponseHeaderTimeout, defaultResponseHeaderTimeout)
	if err != nil {
		return nil, fmt.Errorf("config %s: response_header_timeout: %w", path, err)
	}
	if raw.MaxBodyBytes < 0 {
		return nil, fmt.Errorf("config %s: max_body_bytes must not be negative", path)
	}

	return &Config{
		Listen:                orDefaultStr(raw.Listen, defaultListen),
		NodeID:                raw.NodeID,
		ModelServer:           orDefaultStr(raw.ModelServer, defaultModelServer),
		ConsumerToken:         raw.ConsumerToken,
		Advertise:             raw.Advertise,
		Peers:                 raw.Peers,
		StartupTimeout:        startupTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
		MaxBodyBytes:          orDefaultInt64(raw.MaxBodyBytes, defaultMaxBodyBytes),
	}, nil
}

func orDefaultStr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func orDefaultInt64(v, def int64) int64 {
	if v == 0 {
		return def
	}
	return v
}

// parseDuration parses a Go duration string; empty -> def. Allows 0 (meaning
// "no limit" for a timeout) but rejects negative values.
func parseDuration(s string, def time.Duration) (time.Duration, error) {
	if s == "" {
		return def, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if d < 0 {
		return 0, fmt.Errorf("must not be negative: %q", s)
	}
	return d, nil
}
