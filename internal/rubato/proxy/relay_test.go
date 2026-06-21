package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPassThroughRelay verifies that requests are relayed unchanged to upstream.
func TestPassThroughRelay(t *testing.T) {
	// Track what was sent to upstream
	var upstreamReq *http.Request
	var upstreamBody []byte

	// Create mock upstream that records the request
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamReq = r
		var err error
		upstreamBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read upstream request body: %v", err)
		}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "test response"}},
			},
		})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "")

	// Create request with specific content
	originalBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"temperature": 0.7,
	}
	bodyBytes, _ := json.Marshal(originalBody)
	clientBody := bytes.NewBuffer(bodyBytes)

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", clientBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	// Verify request was forwarded
	if upstreamReq == nil {
		t.Fatal("request was not forwarded to upstream")
	}

	// Verify body was not mutated
	var forwardedBody map[string]interface{}
	if err := json.Unmarshal(upstreamBody, &forwardedBody); err != nil {
		t.Fatalf("failed to decode forwarded request body: %v", err)
	}

	if forwardedBody["model"] != "gpt-4" {
		t.Errorf("model was mutated: %v", forwardedBody["model"])
	}

	// Verify messages were not mutated
	messages, ok := forwardedBody["messages"].([]interface{})
	if !ok {
		t.Fatalf("messages field is not an array: %T", forwardedBody["messages"])
	}
	if len(messages) != 1 {
		t.Errorf("messages array was mutated")
	}

	// Verify response is relayed unchanged
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode upstream response: %v", err)
	}
	if resp["choices"] == nil {
		t.Errorf("response was not relayed correctly")
	}
}

// TestDeterministicClientError verifies error responses are deterministic.
func TestDeterministicClientError(t *testing.T) {
	handler := NewHandler("http://localhost:8000", "")

	// Test with invalid JSON
	invalidBody := bytes.NewBufferString(`{invalid}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", invalidBody)
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	errObj, ok := errResp["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("error field has unexpected type: %T", errResp["error"])
	}
	if errObj["type"] != "invalid_request_error" {
		t.Errorf("expected type 'invalid_request_error', got %v", errObj["type"])
	}
	if errObj["code"] != "invalid_body" {
		t.Errorf("expected code 'invalid_body', got %v", errObj["code"])
	}
}

// TestDeterministicUpstreamError verifies upstream failure responses are deterministic.
func TestDeterministicUpstreamError(t *testing.T) {
	// Create upstream that closes connection (simulating network failure)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate connection reset
		panic("connection reset")
	}))

	// Close upstream to force connection errors
	upstreamURL := upstream.URL
	upstream.Close()

	handler := NewHandler(upstreamURL, "")

	validBody := bytes.NewBufferString(`{"model": "gpt-4", "messages": []}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", validBody)
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected 502 Bad Gateway, got %d", w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to decode upstream error response: %v", err)
	}

	errObj, ok := errResp["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("error field has unexpected type: %T", errResp["error"])
	}
	if errObj["type"] != "server_error" {
		t.Errorf("expected type 'server_error', got %v", errObj["type"])
	}
	if errObj["code"] != "upstream_failure" {
		t.Errorf("expected code 'upstream_failure', got %v", errObj["code"])
	}
}

// TestHeaderPassThrough verifies that custom headers are forwarded.
func TestHeaderPassThrough(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers were forwarded
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Authorization header not forwarded")
		}
		if r.Header.Get("X-Custom") != "custom-value" {
			t.Errorf("Custom header not forwarded")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "")

	validBody := bytes.NewBufferString(`{"model": "gpt-4", "messages": []}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", validBody)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-Custom", "custom-value")
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TestConfiguredAuthorizationFallback verifies Rubato can inject upstream auth when client omits it.
func TestConfiguredAuthorizationFallback(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer configured-token" {
			t.Errorf("expected configured upstream Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "configured-token")

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"gpt-4","messages":[]}`))
	w := httptest.NewRecorder()

	handler.ChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
