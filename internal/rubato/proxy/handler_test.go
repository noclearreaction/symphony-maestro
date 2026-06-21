package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatCompletions_MethodNotAllowed(t *testing.T) {
	handler := NewHandler("http://localhost:8000", "")

	// Test GET method not allowed
	req := httptest.NewRequest(http.MethodGet, "/v1/chat/completions", nil)
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errResp["error"] == nil {
		t.Errorf("expected error field in response")
	}
}

func TestChatCompletions_MalformedJSON(t *testing.T) {
	handler := NewHandler("http://localhost:8000", "")

	// Test malformed JSON
	malformedBody := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", malformedBody)
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errResp["error"] == nil {
		t.Errorf("expected error field in response")
	}
}

func TestChatCompletions_ValidRequest(t *testing.T) {
	// Create a mock upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"response": "ok"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "")

	// Test valid JSON request
	validBody := bytes.NewBufferString(`{"model": "gpt-4", "messages": []}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", validBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if respBody["response"] != "ok" {
		t.Errorf("expected response 'ok', got %s", respBody["response"])
	}
}

func TestChatCompletions_UpstreamFailure(t *testing.T) {
	// Create a mock upstream server that returns an error
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "")

	validBody := bytes.NewBufferString(`{"model": "gpt-4", "messages": []}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", validBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	// Upstream error should pass through
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestChatCompletions_PassThroughHeaders(t *testing.T) {
	// Create a mock upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "")

	validBody := bytes.NewBufferString(`{"model": "gpt-4", "messages": []}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", validBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	// Check that upstream headers are passed through
	if w.Header().Get("X-Custom-Header") != "custom-value" {
		t.Errorf("expected custom header value 'custom-value', got %s", w.Header().Get("X-Custom-Header"))
	}
}
