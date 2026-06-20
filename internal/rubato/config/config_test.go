package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("RUBATO_UPSTREAM_URL")
	os.Unsetenv("RUBATO_LISTEN_ADDR")

	cfg := Load()

	if cfg.UpstreamURL != "http://localhost:8000" {
		t.Errorf("expected default upstream URL, got %s", cfg.UpstreamURL)
	}

	if cfg.ListenAddr != ":8080" {
		t.Errorf("expected default listen addr, got %s", cfg.ListenAddr)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	os.Setenv("RUBATO_UPSTREAM_URL", "http://api.example.com")
	os.Setenv("RUBATO_LISTEN_ADDR", ":9000")
	defer os.Unsetenv("RUBATO_UPSTREAM_URL")
	defer os.Unsetenv("RUBATO_LISTEN_ADDR")

	cfg := Load()

	if cfg.UpstreamURL != "http://api.example.com" {
		t.Errorf("expected upstream URL from env, got %s", cfg.UpstreamURL)
	}

	if cfg.ListenAddr != ":9000" {
		t.Errorf("expected listen addr from env, got %s", cfg.ListenAddr)
	}
}

func TestValidate(t *testing.T) {
	// Valid config
	validCfg := &Config{
		UpstreamURL: "http://localhost:8000",
		ListenAddr:  ":8080",
	}

	if err := validCfg.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}

	// Invalid - no upstream
	invalidCfg := &Config{
		UpstreamURL: "",
		ListenAddr:  ":8080",
	}

	if err := invalidCfg.Validate(); err == nil {
		t.Errorf("expected validation error for missing upstream")
	}

	// Invalid - no listen addr
	invalidCfg2 := &Config{
		UpstreamURL: "http://localhost:8000",
		ListenAddr:  "",
	}

	if err := invalidCfg2.Validate(); err == nil {
		t.Errorf("expected validation error for missing listen addr")
	}
}

func TestString(t *testing.T) {
	cfg := &Config{
		UpstreamURL: "http://localhost:8000",
		ListenAddr:  ":8080",
	}

	str := cfg.String()
	if len(str) == 0 {
		t.Errorf("expected non-empty string representation")
	}
}
