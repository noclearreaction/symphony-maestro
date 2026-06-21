package main

import (
	"log"
	"net/http"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/config"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/proxy"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	log.Printf("Starting Rubato with %s", cfg)

	// Create the proxy handler
	handler := proxy.NewHandler(cfg.UpstreamURL, cfg.UpstreamAPIKey)

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", handler.ChatCompletions)

	// Start HTTP server
	log.Printf("Listening on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
