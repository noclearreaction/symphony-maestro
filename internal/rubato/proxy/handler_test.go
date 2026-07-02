package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/mutate"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

// injectablePlugin is a test-only Plugin stub.
type injectablePlugin struct {
	name string
	out  string
	err  error
}

func (p *injectablePlugin) Name() string { return p.name }
func (p *injectablePlugin) Execute(_ context.Context, _ map[string]any) (string, error) {
	return p.out, p.err
}

// pluginRegistry returns a Registry loaded with the given stubs.
func pluginRegistry(stubs ...*injectablePlugin) *plugin.Registry {
	plugins := make([]plugin.Plugin, len(stubs))
	for i, s := range stubs {
		plugins[i] = s
	}
	return plugin.NewRegistry(plugins...)
}

// mutateInjector wraps a registry into a mutate.Injector.
func mutateInjector(r *plugin.Registry) *mutate.Injector {
	return mutate.NewInjector(r)
}

func TestChatCompletions_MethodNotAllowed(t *testing.T) {
	handler := NewHandler("http://localhost:8000", "", nil)

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
	handler := NewHandler("http://localhost:8000", "", nil)

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

	handler := NewHandler(upstream.URL, "", nil)

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

	handler := NewHandler(upstream.URL, "", nil)

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

	handler := NewHandler(upstream.URL, "", nil)

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

// TestChatCompletions_InjectorMutatesRequest verifies that when an injector is
// configured and the request contains a rubato anchor, the body forwarded to the
// upstream contains the runtime-state block rather than the original body.
func TestChatCompletions_InjectorMutatesRequest(t *testing.T) {
	var receivedBody []byte

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		receivedBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("upstream read body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer upstream.Close()

	stub := &injectablePlugin{name: "git_status", out: "branch: main\nstaged: 0"}
	registry := pluginRegistry(stub)
	injector := mutateInjector(registry)
	handler := NewHandler(upstream.URL, "", injector)

	reqBody := "{\"model\":\"gpt-4\",\"messages\":[" +
		"{\"role\":\"system\",\"content\":\"```rubato:anchor\\n{\\\"plugins\\\":[\\\"git_status\\\"]}\\n``` You are an assistant.\"}," +
		"{\"role\":\"user\",\"content\":\"Hello\"}" +
		"]}"

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
	if !bytes.Contains(receivedBody, []byte("```rubato:state")) {
		t.Errorf("upstream body missing runtime-state block:\n%s", receivedBody)
	}
	if !bytes.Contains(receivedBody, []byte("branch: main")) {
		t.Errorf("upstream body missing plugin output:\n%s", receivedBody)
	}
}

// TestChatCompletions_InjectorMalformedAnchor verifies a 400 is returned for
// malformed anchor declarations.
func TestChatCompletions_InjectorMalformedAnchor(t *testing.T) {
	stub := &injectablePlugin{name: "git_status", out: "branch: main"}
	registry := pluginRegistry(stub)
	injector := mutateInjector(registry)
	handler := NewHandler("http://localhost:8000", "", injector)

	reqBody := "{\"model\":\"gpt-4\",\"messages\":[" +
		"{\"role\":\"system\",\"content\":\"```rubato:anchor\\nbad json\\n```\"}," +
		"{\"role\":\"user\",\"content\":\"Hello\"}" +
		"]}"

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ChatCompletions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed anchor, got %d", w.Code)
	}
}

// TestChatCompletions_NoInjector_AnchorPassedThrough verifies that when no
// injector is configured, requests with anchor text are forwarded unchanged.
func TestChatCompletions_NoInjector_AnchorPassedThrough(t *testing.T) {
	var receivedBody []byte
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer upstream.Close()

	handler := NewHandler(upstream.URL, "", nil)

	reqBody := "{\"model\":\"gpt-4\",\"messages\":[{\"role\":\"system\",\"content\":\"```rubato:anchor\\n{\\\"plugins\\\":[\\\"git_status\\\"]}\\n```\"}" +
		",{\"role\":\"user\",\"content\":\"Hello\"}]}"
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ChatCompletions(w, req)

	if bytes.Contains(receivedBody, []byte("```rubato:state")) {
		t.Errorf("runtime-state injected without an injector configured")
	}
}
