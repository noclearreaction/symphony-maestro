package config

import (
	"fmt"
	"os"
)

// Config holds runtime configuration for Rubato.
type Config struct {
	// UpstreamURL is the base URL of the upstream service
	UpstreamURL string
	// UpstreamAPIKey is an optional bearer token used for upstream requests.
	UpstreamAPIKey string
	// ListenAddr is the address to listen on (host:port)
	ListenAddr string
}

// Load reads configuration from environment variables.
// Defaults to sensible values if not provided.
func Load() *Config {
	upstream := os.Getenv("RUBATO_UPSTREAM_URL")
	if upstream == "" {
		upstream = "http://localhost:8000"
	}

	listenAddr := os.Getenv("RUBATO_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	upstreamAPIKey := os.Getenv("OPENROUTER_API_KEY")

	return &Config{
		UpstreamURL:    upstream,
		UpstreamAPIKey: upstreamAPIKey,
		ListenAddr:     listenAddr,
	}
}

// Validate ensures the configuration is valid.
func (c *Config) Validate() error {
	if c.UpstreamURL == "" {
		return fmt.Errorf("upstream URL is required")
	}
	if c.ListenAddr == "" {
		return fmt.Errorf("listen address is required")
	}
	return nil
}

// String returns a string representation of the config.
func (c *Config) String() string {
	keyStatus := "unset"
	if c.UpstreamAPIKey != "" {
		keyStatus = "set"
	}
	return fmt.Sprintf("Config{UpstreamURL: %s, UpstreamAPIKey: %s, ListenAddr: %s}", c.UpstreamURL, keyStatus, c.ListenAddr)
}
