package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/mutate"
)

// Handler manages proxy requests to the upstream service.
type Handler struct {
	upstreamURL    string
	upstreamAPIKey string
	client         *http.Client
	injector       *mutate.Injector // nil = no injection
}

// NewHandler creates a new proxy handler.
// injector may be nil to disable plugin-based injection.
func NewHandler(upstreamURL, upstreamAPIKey string, injector *mutate.Injector) *Handler {
	return &Handler{
		upstreamURL:    upstreamURL,
		upstreamAPIKey: upstreamAPIKey,
		client:         &http.Client{},
		injector:       injector,
	}
}

// ChatCompletions handles POST /v1/chat/completions requests.
func (h *Handler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	log.Printf("proxying %s %s", r.Method, r.URL.Path)

	// Only allow POST requests
	if r.Method != http.MethodPost {
		h.respondMethodNotAllowed(w)
		return
	}

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		h.respondBadRequest(w, "failed to read request body")
		return
	}
	defer r.Body.Close()

	// Validate JSON
	var req interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("error unmarshaling JSON: %v", err)
		h.respondBadRequest(w, "invalid JSON in request body")
		return
	}
	log.Printf("request body: %s", body)

	// Apply plugin injection when configured.
	if h.injector != nil {
		mutated, err := h.injector.Apply(r.Context(), body)
		if err != nil {
			log.Printf("injection error: %v", err)
			h.respondBadRequest(w, "request injection failed: "+err.Error())
			return
		}
		body = mutated
		log.Printf("injected request body: %s", body)
	}

	// Forward to upstream
	resp, err := h.forwardRequest(r.Context(), r.URL.Path, body, r.Header)
	if err != nil {
		log.Printf("error forwarding request: %v", err)
		h.respondUpstreamFailure(w)
		return
	}
	defer resp.Body.Close()

	// Copy upstream response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy upstream status code
	w.WriteHeader(resp.StatusCode)

	// Copy upstream response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

// forwardRequest sends the request to the upstream service.
func (h *Handler) forwardRequest(ctx context.Context, path string, body []byte, header http.Header) (*http.Response, error) {
	url := h.upstreamURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Copy headers from original request
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Ensure content type is set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Prefer client Authorization header; fall back to configured upstream key.
	if req.Header.Get("Authorization") == "" && h.upstreamAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+h.upstreamAPIKey)
	}

	return h.client.Do(req)
}

// respondBadRequest returns a deterministic 400 error.
func (h *Handler) respondBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    "invalid_request_error",
			"code":    "invalid_body",
		},
	})
}

// respondMethodNotAllowed returns a deterministic 405 error.
func (h *Handler) respondMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"message": "method not allowed",
			"type":    "invalid_request_error",
			"code":    "method_not_allowed",
		},
	})
}

// respondUpstreamFailure returns a deterministic 502 error for upstream failures.
func (h *Handler) respondUpstreamFailure(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadGateway)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"message": "upstream service failed",
			"type":    "server_error",
			"code":    "upstream_failure",
		},
	})
}
